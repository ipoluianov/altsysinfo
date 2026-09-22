package system

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/v4/cpu"
)

var chassisTypes = map[string]string{
	"3": "Desktop", "4": "Low Profile Desktop", "5": "Pizza Box", "6": "Mini Tower",
	"7": "Tower", "8": "Portable", "9": "Laptop", "10": "Notebook", "11": "Hand Held",
	"13": "All In One", "14": "Sub Notebook", "15": "Space-saving", "16": "Lunch Box",
	"17": "Main Server Chassis", "23": "Rack Mount Chassis", "24": "Sealed-case PC",
	"30": "Tablet", "31": "Convertible", "32": "Detachable", "35": "Mini PC", "36": "Stick PC",
}

func addPlatformHardwareDetails(t *DetailTable) {
	dmi := func(name string) string {
		return readString(filepath.Join("/sys/class/dmi/id", name))
	}

	t.Section("Computer")
	t.AddValue("Manufacturer", dmi("sys_vendor"))
	t.AddValue("Product", dmi("product_name"))
	t.AddValue("Version", dmi("product_version"))
	t.AddValue("Family", dmi("product_family"))
	t.AddValue("Chassis", chassisTypes[dmi("chassis_type")])

	t.Section("Motherboard")
	t.AddValue("Manufacturer", dmi("board_vendor"))
	t.AddValue("Model", dmi("board_name"))
	t.AddValue("Version", dmi("board_version"))

	t.Section("BIOS / Firmware")
	t.AddValue("Vendor", dmi("bios_vendor"))
	t.AddValue("Version", dmi("bios_version"))
	t.AddValue("Release date", dmi("bios_date"))
	t.AddValue("Release", dmi("bios_release"))

	if _, err := os.Stat("/sys/firmware/efi"); err == nil {
		t.Add("Boot mode", "UEFI")
		// The first 4 bytes of an efivar are attributes, followed by the value
		data, err := os.ReadFile("/sys/firmware/efi/efivars/SecureBoot-8be4df61-93ca-11d2-aa0d-00e098032b8c")
		if err == nil && len(data) >= 5 {
			t.Add("Secure Boot", map[bool]string{true: "Enabled", false: "Disabled"}[data[4] == 1])
		}
	} else {
		t.Add("Boot mode", "Legacy BIOS")
	}
}

func addPlatformOSDetails(t *DetailTable) {
	osRelease := readKeyValueFile("/etc/os-release", "=")
	t.AddValue("Distribution", strings.Trim(osRelease["PRETTY_NAME"], `"`))
	t.AddValue("Desktop", os.Getenv("XDG_CURRENT_DESKTOP"))
	t.AddValue("Session type", os.Getenv("XDG_SESSION_TYPE"))
}

func addPlatformCPUDetails(t *DetailTable, infos []cpu.InfoStat) {
	sockets := make(map[string]struct{})
	flags := make(map[string]bool)
	for _, info := range infos {
		sockets[info.PhysicalID] = struct{}{}
		for _, f := range info.Flags {
			flags[f] = true
		}
	}
	if len(sockets) > 0 {
		t.Add("Sockets", strconv.Itoa(len(sockets)))
	}
	switch {
	case flags["vmx"]:
		t.Add("Virtualization", "Intel VT-x")
	case flags["svm"]:
		t.Add("Virtualization", "AMD-V")
	}

	const cpuBase = "/sys/devices/system/cpu"
	freq := func(name string) string {
		v := readUint(filepath.Join(cpuBase, "cpu0/cpufreq", name))
		if v == 0 {
			return ""
		}
		return fmt.Sprintf("%d MHz", v/1000)
	}
	t.Section("Frequency")
	t.AddValue("Minimum", freq("cpuinfo_min_freq"))
	t.AddValue("Maximum", freq("cpuinfo_max_freq"))
	t.AddValue("Current (CPU 0)", freq("scaling_cur_freq"))
	t.AddValue("Scaling driver", readString(filepath.Join(cpuBase, "cpu0/cpufreq/scaling_driver")))
	t.AddValue("Governor", readString(filepath.Join(cpuBase, "cpu0/cpufreq/scaling_governor")))

	// Count distinct cache instances (by the set of CPUs sharing them)
	type cacheKey struct{ level, kind, size string }
	instances := make(map[cacheKey]map[string]struct{})
	indexes, _ := filepath.Glob(filepath.Join(cpuBase, "cpu[0-9]*/cache/index[0-9]*"))
	for _, dir := range indexes {
		key := cacheKey{
			level: readString(filepath.Join(dir, "level")),
			kind:  readString(filepath.Join(dir, "type")),
			size:  readString(filepath.Join(dir, "size")),
		}
		if instances[key] == nil {
			instances[key] = make(map[string]struct{})
		}
		instances[key][readString(filepath.Join(dir, "shared_cpu_list"))] = struct{}{}
	}
	keys := make([]cacheKey, 0, len(instances))
	for k := range instances {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].level != keys[j].level {
			return keys[i].level < keys[j].level
		}
		return keys[i].kind < keys[j].kind
	})
	if len(keys) > 0 {
		t.Section("Cache")
		for _, k := range keys {
			name := "L" + k.level
			if k.kind != "Unified" {
				name += " " + k.kind
			}
			t.Add(name, fmt.Sprintf("%s x %d", k.size, len(instances[k])))
		}
	}
}

func getGraphicsDetails() ([]*DetailTable, error) {
	gpus := newDetailTable("Graphics adapters", "Card", "Model", "ID", "Driver", "Video memory")
	if list, err := getGPUs(); err == nil {
		for _, gpu := range list {
			vram := readUint(filepath.Join("/sys/class/drm", gpu.Name, "device/mem_info_vram_total"))
			vramStr := ""
			if vram > 0 {
				vramStr = formatBytes(vram)
			}
			gpus.Add(gpu.Name, gpu.Model, strings.TrimPrefix(gpu.Vendor, "0x")+":"+strings.TrimPrefix(gpu.Device, "0x"), gpu.Driver, vramStr)
		}
	}

	displays := newDetailTable("Displays", "Connector", "Manufacturer", "Model", "Resolution", "Size", "Year")
	connectors, _ := filepath.Glob("/sys/class/drm/card[0-9]*-*")
	sort.Strings(connectors)
	for _, dir := range connectors {
		if readString(filepath.Join(dir, "status")) != "connected" {
			continue
		}
		connector := filepath.Base(dir)
		connector = connector[strings.Index(connector, "-")+1:]

		resolution := ""
		if modes := strings.Fields(readString(filepath.Join(dir, "modes"))); len(modes) > 0 {
			resolution = modes[0]
		}

		edid, _ := os.ReadFile(filepath.Join(dir, "edid"))
		info := parseEDID(edid)
		displays.Add(connector, info.manufacturer, info.model, resolution, info.size, info.year)
	}

	return []*DetailTable{gpus, displays}, nil
}

type edidInfo struct {
	manufacturer string
	model        string
	size         string
	year         string
}

func parseEDID(data []byte) edidInfo {
	var info edidInfo
	if len(data) < 128 || data[0] != 0x00 || data[1] != 0xFF {
		return info
	}

	// Manufacturer ID: three 5-bit letters, big endian
	id := uint16(data[8])<<8 | uint16(data[9])
	info.manufacturer = string([]byte{
		byte('A' - 1 + (id>>10)&0x1F),
		byte('A' - 1 + (id>>5)&0x1F),
		byte('A' - 1 + id&0x1F),
	})

	if data[17] > 0 {
		info.year = strconv.Itoa(1990 + int(data[17]))
	}

	if w, h := float64(data[21]), float64(data[22]); w > 0 && h > 0 {
		info.size = fmt.Sprintf("%.1f\" (%.0fx%.0f cm)", math.Sqrt(w*w+h*h)/2.54, w, h)
	}

	// Descriptor blocks: tag 0xFC is the monitor name
	for offset := 54; offset+18 <= 126; offset += 18 {
		d := data[offset : offset+18]
		if d[0] == 0 && d[1] == 0 && d[3] == 0xFC {
			info.model = strings.TrimSpace(strings.SplitN(string(d[5:]), "\n", 2)[0])
		}
	}
	if info.model == "" {
		info.model = fmt.Sprintf("%s%04X", info.manufacturer, uint16(data[10])|uint16(data[11])<<8)
	}

	return info
}

func platformInterfaceExtras() map[string]interfaceExtra {
	result := make(map[string]interfaceExtra)
	dirs, _ := filepath.Glob("/sys/class/net/*")
	for _, dir := range dirs {
		var extra interfaceExtra
		name := filepath.Base(dir)

		if speed, err := strconv.Atoi(readString(filepath.Join(dir, "speed"))); err == nil && speed > 0 {
			extra.Speed = strconv.Itoa(speed) + " Mb/s"
		}

		device := filepath.Join(dir, "device")
		if _, err := os.Stat(device); err != nil {
			extra.Description = "Virtual"
		} else {
			if link, err := filepath.EvalSymlinks(filepath.Join(device, "driver")); err == nil {
				extra.Driver = filepath.Base(link)
			}
			if model, ok := PCIName(readHex16(filepath.Join(device, "vendor")), readHex16(filepath.Join(device, "device"))); ok {
				extra.Description = model
			}
		}
		if _, err := os.Stat(filepath.Join(dir, "wireless")); err == nil {
			extra.Description = strings.TrimSpace("Wireless " + extra.Description)
		}

		result[name] = extra
	}
	return result
}

func getPCIDetails() ([]*DetailTable, error) {
	t := newDetailTable("PCI devices", "Address", "Class", "Vendor", "Device", "ID", "Driver")
	dirs, err := filepath.Glob("/sys/bus/pci/devices/*")
	if err != nil {
		return nil, err
	}
	sort.Strings(dirs)
	for _, dir := range dirs {
		vendorID := readHex16(filepath.Join(dir, "vendor"))
		deviceID := readHex16(filepath.Join(dir, "device"))

		class, _ := strconv.ParseUint(strings.TrimPrefix(readString(filepath.Join(dir, "class")), "0x"), 16, 32)

		vendor, _ := PCIVendorName(vendorID)
		device, ok := PCIName(vendorID, deviceID)
		if !ok {
			device = "Unknown"
		}

		driver := ""
		if link, err := filepath.EvalSymlinks(filepath.Join(dir, "driver")); err == nil {
			driver = filepath.Base(link)
		}

		t.Add(filepath.Base(dir), PCIClassName(uint8(class>>16), uint8(class>>8)), vendor, device,
			fmt.Sprintf("%04x:%04x", vendorID, deviceID), driver)
	}
	return []*DetailTable{t}, nil
}

func getUSBDetails() ([]*DetailTable, error) {
	t := newDetailTable("USB devices", "Bus/Device", "ID", "Manufacturer", "Product", "Speed", "USB version")
	dirs, err := filepath.Glob("/sys/bus/usb/devices/*")
	if err != nil {
		return nil, err
	}

	type usbDevice struct {
		bus, dev int
		cells    []string
	}
	var devices []usbDevice
	usbIDs := loadUSBIDs()

	for _, dir := range dirs {
		if strings.Contains(filepath.Base(dir), ":") { // interfaces
			continue
		}
		vendorID := readString(filepath.Join(dir, "idVendor"))
		productID := readString(filepath.Join(dir, "idProduct"))
		if vendorID == "" {
			continue
		}

		manufacturer := readString(filepath.Join(dir, "manufacturer"))
		product := readString(filepath.Join(dir, "product"))
		if manufacturer == "" {
			manufacturer = usbIDs[vendorID]
		}
		if product == "" {
			product = usbIDs[vendorID+":"+productID]
		}

		speed := readString(filepath.Join(dir, "speed"))
		if speed != "" {
			speed += " Mb/s"
		}

		d := usbDevice{
			bus: int(readUint(filepath.Join(dir, "busnum"))),
			dev: int(readUint(filepath.Join(dir, "devnum"))),
		}
		d.cells = []string{
			fmt.Sprintf("%03d/%03d", d.bus, d.dev),
			vendorID + ":" + productID,
			manufacturer, product, speed,
			readString(filepath.Join(dir, "version")),
		}
		devices = append(devices, d)
	}

	sort.Slice(devices, func(i, j int) bool {
		if devices[i].bus != devices[j].bus {
			return devices[i].bus < devices[j].bus
		}
		return devices[i].dev < devices[j].dev
	})
	for _, d := range devices {
		t.Add(d.cells...)
	}
	return []*DetailTable{t}, nil
}

// loadUSBIDs reads the system usb.ids database if it is installed.
// Keys are "vvvv" for vendors and "vvvv:pppp" for products.
func loadUSBIDs() map[string]string {
	result := make(map[string]string)
	for _, path := range []string{"/usr/share/hwdata/usb.ids", "/usr/share/misc/usb.ids", "/usr/share/usb.ids"} {
		f, err := os.Open(path)
		if err != nil {
			continue
		}
		defer f.Close()

		vendor := ""
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			if line == "" || line[0] == '#' || strings.HasPrefix(line, "\t\t") {
				continue
			}
			fields := strings.Fields(line)
			if len(fields) < 2 || len(fields[0]) != 4 {
				if line[0] != '\t' {
					vendor = ""
				}
				continue
			}
			if line[0] != '\t' {
				// Vendor list ends where the class and other sections begin
				if _, err := strconv.ParseUint(fields[0], 16, 16); err != nil {
					break
				}
				vendor = fields[0]
				result[vendor] = strings.TrimSpace(line[4:])
			} else if vendor != "" {
				result[vendor+":"+fields[0]] = strings.TrimSpace(line[5:])
			}
		}
		break
	}
	return result
}

func getAudioDetails() ([]*DetailTable, error) {
	t := newDetailTable("Sound cards", "#", "ID", "Driver", "Name", "Codec")
	f, err := os.Open("/proc/asound/cards")
	if err != nil {
		return []*DetailTable{t}, nil
	}
	defer f.Close()

	// Format: " 0 [PCH            ]: HDA-Intel - HDA Intel PCH"
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		index, rest, ok := strings.Cut(strings.TrimSpace(line), " ")
		if !ok || !strings.HasPrefix(strings.TrimSpace(rest), "[") {
			continue
		}
		id, rest, _ := strings.Cut(strings.TrimSpace(rest)[1:], "]:")
		driver, name, _ := strings.Cut(rest, " - ")

		var codecs []string
		codecFiles, _ := filepath.Glob("/proc/asound/card" + index + "/codec#*")
		for _, file := range codecFiles {
			codec := readKeyValueFile(file, ":")["Codec"]
			if codec != "" {
				codecs = append(codecs, codec)
			}
		}

		t.Add(index, id, driver, name, strings.Join(codecs, ", "))
	}
	return []*DetailTable{t}, nil
}

func getSensorsDetails() ([]*DetailTable, error) {
	t := newDetailTable("Sensors", "Chip", "Sensor", "Value")
	dirs, _ := filepath.Glob("/sys/class/hwmon/hwmon*")
	sort.Strings(dirs)

	type sensorKind struct {
		prefix string
		format func(v float64) string
	}
	kinds := []sensorKind{
		{"temp", func(v float64) string { return fmt.Sprintf("%.1f °C", v/1000) }},
		{"fan", func(v float64) string { return fmt.Sprintf("%.0f RPM", v) }},
		{"in", func(v float64) string { return fmt.Sprintf("%.3f V", v/1000) }},
		{"power", func(v float64) string { return fmt.Sprintf("%.1f W", v/1000000) }},
	}

	for _, dir := range dirs {
		chip := readString(filepath.Join(dir, "name"))
		for _, kind := range kinds {
			inputs, _ := filepath.Glob(filepath.Join(dir, kind.prefix+"[0-9]*_input"))
			sort.Strings(inputs)
			for _, input := range inputs {
				value, err := strconv.ParseFloat(readString(input), 64)
				if err != nil {
					continue
				}
				base := strings.TrimSuffix(filepath.Base(input), "_input")
				label := readString(filepath.Join(dir, base+"_label"))
				if label == "" {
					label = base
				}
				t.Add(chip, label, kind.format(value))
			}
		}
	}
	return []*DetailTable{t}, nil
}

func getBatteryDetails() ([]*DetailTable, error) {
	t := newKeyValueTable("Batteries")
	dirs, _ := filepath.Glob("/sys/class/power_supply/*")
	sort.Strings(dirs)

	for _, dir := range dirs {
		props := readKeyValueFile(filepath.Join(dir, "uevent"), "=")
		p := func(name string) string { return props["POWER_SUPPLY_"+name] }
		num := func(name string) float64 {
			v, _ := strconv.ParseFloat(p(name), 64)
			return v
		}

		switch p("TYPE") {
		case "Mains":
			t.Section("AC adapter " + filepath.Base(dir))
			t.Add("Online", yesNo(p("ONLINE") == "1"))
			continue
		case "Battery":
		default:
			continue
		}

		title := strings.TrimSpace(p("MANUFACTURER") + " " + p("MODEL_NAME"))
		if title == "" {
			title = filepath.Base(dir)
		}
		t.Section(title)
		if p("SCOPE") == "Device" {
			t.Add("Type", "Peripheral device")
		}
		t.AddValue("Technology", p("TECHNOLOGY"))
		t.AddValue("Status", p("STATUS"))
		if p("CAPACITY") != "" {
			t.Add("Charge", p("CAPACITY")+" %")
		}
		t.AddValue("Charge level", p("CAPACITY_LEVEL"))

		// Energy is reported in µWh, charge in µAh
		full, design, unit := num("ENERGY_FULL"), num("ENERGY_FULL_DESIGN"), "Wh"
		if full == 0 {
			full, design, unit = num("CHARGE_FULL"), num("CHARGE_FULL_DESIGN"), "Ah"
		}
		if full > 0 {
			t.Add("Full capacity", fmt.Sprintf("%.2f %s", full/1e6, unit))
		}
		if design > 0 {
			t.Add("Design capacity", fmt.Sprintf("%.2f %s", design/1e6, unit))
		}
		if full > 0 && design > 0 {
			t.Add("Health", fmt.Sprintf("%.0f %% (wear %.0f %%)", full/design*100, math.Max(0, 100-full/design*100)))
		}
		if p("CYCLE_COUNT") != "" && p("CYCLE_COUNT") != "0" {
			t.Add("Cycle count", p("CYCLE_COUNT"))
		}
		if v := num("VOLTAGE_NOW"); v > 0 {
			t.Add("Voltage", fmt.Sprintf("%.2f V", v/1e6))
		}
		t.AddValue("Serial number", p("SERIAL_NUMBER"))
	}
	return []*DetailTable{t}, nil
}

// readKeyValueFile parses "key<sep>value" lines.
func readKeyValueFile(path, sep string) map[string]string {
	result := make(map[string]string)
	data, err := os.ReadFile(path)
	if err != nil {
		return result
	}
	for _, line := range strings.Split(string(data), "\n") {
		key, value, ok := strings.Cut(line, sep)
		if ok {
			result[strings.TrimSpace(key)] = strings.TrimSpace(value)
		}
	}
	return result
}
