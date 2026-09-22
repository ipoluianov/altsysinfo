package system

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/StackExchange/wmi"
	"github.com/shirou/gopsutil/v4/cpu"
	"golang.org/x/sys/windows/registry"
)

// wmiQuery runs a WMI query. Field type mismatches are ignored (the rest of
// the data is still loaded), and panics in the WMI layer become errors.
func wmiQuery(query string, dst any, namespace ...string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("wmi: %v", r)
		}
	}()
	if len(namespace) > 0 {
		err = wmi.QueryNamespace(query, dst, namespace[0])
	} else {
		err = wmi.Query(query, dst)
	}
	var mismatch *wmi.ErrFieldMismatch
	if errors.As(err, &mismatch) {
		return nil
	}
	return err
}

// wmiDate converts a CIM datetime ("20240115000000.000000+000") to YYYY-MM-DD.
func wmiDate(s string) string {
	if len(s) < 8 {
		return s
	}
	return s[0:4] + "-" + s[4:6] + "-" + s[6:8]
}

func registryString(path, name string) string {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, path, registry.QUERY_VALUE)
	if err != nil {
		return ""
	}
	defer key.Close()
	v, _, err := key.GetStringValue(name)
	if err != nil {
		return ""
	}
	return v
}

func registryUint(path, name string) (uint64, bool) {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, path, registry.QUERY_VALUE)
	if err != nil {
		return 0, false
	}
	defer key.Close()
	v, _, err := key.GetIntegerValue(name)
	return v, err == nil
}

func addPlatformHardwareDetails(t *DetailTable) {
	var computers []struct {
		Manufacturer string
		Model        string
		SystemFamily string
		SystemType   string
	}
	if wmiQuery("SELECT Manufacturer, Model, SystemFamily, SystemType FROM Win32_ComputerSystem", &computers) == nil {
		for _, c := range computers {
			t.Section("Computer")
			t.AddValue("Manufacturer", c.Manufacturer)
			t.AddValue("Model", c.Model)
			t.AddValue("Family", c.SystemFamily)
			t.AddValue("System type", c.SystemType)
		}
	}

	var boards []struct {
		Manufacturer string
		Product      string
		Version      string
	}
	if wmiQuery("SELECT Manufacturer, Product, Version FROM Win32_BaseBoard", &boards) == nil {
		for _, b := range boards {
			t.Section("Motherboard")
			t.AddValue("Manufacturer", b.Manufacturer)
			t.AddValue("Model", b.Product)
			t.AddValue("Version", b.Version)
		}
	}

	t.Section("BIOS / Firmware")
	var bioses []struct {
		Manufacturer       string
		SMBIOSBIOSVersion  string
		ReleaseDate        string
		SMBIOSMajorVersion uint16
		SMBIOSMinorVersion uint16
	}
	if wmiQuery("SELECT Manufacturer, SMBIOSBIOSVersion, ReleaseDate, SMBIOSMajorVersion, SMBIOSMinorVersion FROM Win32_BIOS", &bioses) == nil {
		for _, b := range bioses {
			t.AddValue("Vendor", b.Manufacturer)
			t.AddValue("Version", b.SMBIOSBIOSVersion)
			t.AddValue("Release date", wmiDate(b.ReleaseDate))
			if b.SMBIOSMajorVersion > 0 {
				t.Add("SMBIOS version", fmt.Sprintf("%d.%d", b.SMBIOSMajorVersion, b.SMBIOSMinorVersion))
			}
		}
	}

	// PEFirmwareType: 1 = BIOS, 2 = UEFI
	if v, ok := registryUint(`SYSTEM\CurrentControlSet\Control`, "PEFirmwareType"); ok {
		t.Add("Boot mode", map[bool]string{true: "UEFI", false: "Legacy BIOS"}[v == 2])
	}
	if v, ok := registryUint(`SYSTEM\CurrentControlSet\Control\SecureBoot\State`, "UEFISecureBootEnabled"); ok {
		t.Add("Secure Boot", map[bool]string{true: "Enabled", false: "Disabled"}[v == 1])
	}
}

func addPlatformOSDetails(t *DetailTable) {
	const currentVersion = `SOFTWARE\Microsoft\Windows NT\CurrentVersion`
	t.AddValue("Edition", registryString(currentVersion, "EditionID"))
	t.AddValue("Version", registryString(currentVersion, "DisplayVersion"))

	var systems []struct {
		BuildNumber string
		InstallDate string
	}
	if wmiQuery("SELECT BuildNumber, InstallDate FROM Win32_OperatingSystem", &systems) == nil {
		for _, s := range systems {
			t.AddValue("Build", s.BuildNumber)
			t.AddValue("Install date", wmiDate(s.InstallDate))
		}
	}
}

func addPlatformCPUDetails(t *DetailTable, _ []cpu.InfoStat) {
	var processors []struct {
		SocketDesignation             string
		MaxClockSpeed                 uint32
		CurrentClockSpeed             uint32
		L2CacheSize                   uint32
		L3CacheSize                   uint32
		VirtualizationFirmwareEnabled bool
	}
	err := wmiQuery("SELECT SocketDesignation, MaxClockSpeed, CurrentClockSpeed, L2CacheSize, L3CacheSize, VirtualizationFirmwareEnabled FROM Win32_Processor", &processors)
	if err != nil || len(processors) == 0 {
		return
	}
	p := processors[0]
	t.Add("Sockets", strconv.Itoa(len(processors)))
	t.AddValue("Socket", p.SocketDesignation)
	t.Add("Virtualization enabled", yesNo(p.VirtualizationFirmwareEnabled))

	t.Section("Frequency")
	if p.MaxClockSpeed > 0 {
		t.Add("Maximum", fmt.Sprintf("%d MHz", p.MaxClockSpeed))
	}
	if p.CurrentClockSpeed > 0 {
		t.Add("Current", fmt.Sprintf("%d MHz", p.CurrentClockSpeed))
	}

	t.Section("Cache")
	if p.L2CacheSize > 0 {
		t.Add("L2", fmt.Sprintf("%d KB", p.L2CacheSize))
	}
	if p.L3CacheSize > 0 {
		t.Add("L3", fmt.Sprintf("%d KB", p.L3CacheSize))
	}
}

func getGraphicsDetails() ([]*DetailTable, error) {
	gpus := newDetailTable("Graphics adapters", "Model", "ID", "Video memory", "Resolution", "Refresh rate", "Driver", "Driver date")
	var controllers []struct {
		Name                        string
		PNPDeviceID                 string
		AdapterRAM                  uint32
		DriverVersion               string
		DriverDate                  string
		CurrentHorizontalResolution uint32
		CurrentVerticalResolution   uint32
		CurrentRefreshRate          uint32
	}
	if wmiQuery("SELECT Name, PNPDeviceID, AdapterRAM, DriverVersion, DriverDate, CurrentHorizontalResolution, CurrentVerticalResolution, CurrentRefreshRate FROM Win32_VideoController", &controllers) == nil {
		for _, c := range controllers {
			ven, dev := parsePCIID(c.PNPDeviceID)
			id := ""
			if ven != "" {
				id = strings.ToLower(ven + ":" + dev)
			}
			vram := ""
			if c.AdapterRAM > 0 {
				// AdapterRAM is 32-bit, so it is capped at 4 GB
				vram = formatBytes(uint64(c.AdapterRAM))
			}
			resolution, refresh := "", ""
			if c.CurrentHorizontalResolution > 0 {
				resolution = fmt.Sprintf("%dx%d", c.CurrentHorizontalResolution, c.CurrentVerticalResolution)
			}
			if c.CurrentRefreshRate > 0 {
				refresh = fmt.Sprintf("%d Hz", c.CurrentRefreshRate)
			}
			gpus.Add(c.Name, id, vram, resolution, refresh, c.DriverVersion, wmiDate(c.DriverDate))
		}
	}

	displays := newDetailTable("Displays", "Manufacturer", "Model", "Size", "Year")
	var ids []struct {
		InstanceName      string
		ManufacturerName  []uint16
		UserFriendlyName  []uint16
		ProductCodeID     []uint16
		YearOfManufacture uint16
	}
	var params []struct {
		InstanceName           string
		MaxHorizontalImageSize uint8
		MaxVerticalImageSize   uint8
	}
	const wmiNamespace = `root\wmi`
	if wmiQuery("SELECT InstanceName, ManufacturerName, UserFriendlyName, ProductCodeID, YearOfManufacture FROM WmiMonitorID", &ids, wmiNamespace) == nil {
		sizes := make(map[string]string)
		if wmiQuery("SELECT InstanceName, MaxHorizontalImageSize, MaxVerticalImageSize FROM WmiMonitorBasicDisplayParams", &params, wmiNamespace) == nil {
			for _, p := range params {
				if w, h := float64(p.MaxHorizontalImageSize), float64(p.MaxVerticalImageSize); w > 0 && h > 0 {
					sizes[p.InstanceName] = fmt.Sprintf("%.1f\" (%.0fx%.0f cm)", math.Sqrt(w*w+h*h)/2.54, w, h)
				}
			}
		}
		for _, m := range ids {
			model := utf16Field(m.UserFriendlyName)
			if model == "" {
				model = utf16Field(m.ProductCodeID)
			}
			year := ""
			if m.YearOfManufacture > 0 {
				year = strconv.Itoa(int(m.YearOfManufacture))
			}
			displays.Add(utf16Field(m.ManufacturerName), model, sizes[m.InstanceName], year)
		}
	}

	return []*DetailTable{gpus, displays}, nil
}

// utf16Field converts a zero-padded WMI character array to a string.
func utf16Field(chars []uint16) string {
	var sb strings.Builder
	for _, c := range chars {
		if c == 0 {
			break
		}
		sb.WriteRune(rune(c))
	}
	return strings.TrimSpace(sb.String())
}

func platformInterfaceExtras() map[string]interfaceExtra {
	result := make(map[string]interfaceExtra)
	var adapters []struct {
		NetConnectionID string
		Name            string
		Speed           uint64
		ServiceName     string
	}
	if wmiQuery("SELECT NetConnectionID, Name, Speed, ServiceName FROM Win32_NetworkAdapter WHERE NetConnectionID IS NOT NULL", &adapters) != nil {
		return result
	}
	for _, a := range adapters {
		extra := interfaceExtra{Description: a.Name, Driver: a.ServiceName}
		// Speed is in bit/s; disconnected adapters report a huge placeholder value
		if a.Speed > 0 && a.Speed < 1<<62 {
			extra.Speed = strconv.FormatUint(a.Speed/1000000, 10) + " Mb/s"
		}
		result[a.NetConnectionID] = extra
	}
	return result
}

type pnpEntity struct {
	Name         string
	Manufacturer string
	PNPClass     string
	Service      string
	DeviceID     string
	Status       string
}

func queryPnPEntities(prefix string) ([]pnpEntity, error) {
	var entities []pnpEntity
	query := `SELECT Name, Manufacturer, PNPClass, Service, DeviceID, Status FROM Win32_PnPEntity WHERE DeviceID LIKE '` + prefix + `\\%'`
	if err := wmiQuery(query, &entities); err != nil {
		return nil, err
	}
	sort.Slice(entities, func(i, j int) bool { return entities[i].DeviceID < entities[j].DeviceID })
	return entities, nil
}

func getPCIDetails() ([]*DetailTable, error) {
	entities, err := queryPnPEntities("PCI")
	if err != nil {
		return nil, err
	}
	t := newDetailTable("PCI devices", "Class", "Vendor", "Device", "ID", "Driver", "Status")
	for _, e := range entities {
		ven, dev := parsePCIID(e.DeviceID)
		t.Add(e.PNPClass, e.Manufacturer, e.Name, strings.ToLower(ven+":"+dev), e.Service, e.Status)
	}
	return []*DetailTable{t}, nil
}

func getUSBDetails() ([]*DetailTable, error) {
	entities, err := queryPnPEntities("USB")
	if err != nil {
		return nil, err
	}
	t := newDetailTable("USB devices", "Class", "Manufacturer", "Name", "ID", "Driver", "Status")
	for _, e := range entities {
		t.Add(e.PNPClass, e.Manufacturer, e.Name, parseUSBID(e.DeviceID), e.Service, e.Status)
	}
	return []*DetailTable{t}, nil
}

// parseUSBID extracts "vvvv:pppp" from "USB\VID_046D&PID_C52B\...".
func parseUSBID(deviceID string) string {
	var vid, pid string
	parts := strings.Split(deviceID, `\`)
	if len(parts) < 2 {
		return ""
	}
	for _, part := range strings.Split(parts[1], "&") {
		switch {
		case strings.HasPrefix(part, "VID_"):
			vid = strings.TrimPrefix(part, "VID_")
		case strings.HasPrefix(part, "PID_"):
			pid = strings.TrimPrefix(part, "PID_")
		}
	}
	if vid == "" {
		return ""
	}
	return strings.ToLower(vid + ":" + pid)
}

func getAudioDetails() ([]*DetailTable, error) {
	var devices []struct {
		Name         string
		Manufacturer string
		Status       string
	}
	if err := wmiQuery("SELECT Name, Manufacturer, Status FROM Win32_SoundDevice", &devices); err != nil {
		return nil, err
	}
	t := newDetailTable("Sound devices", "Name", "Manufacturer", "Status")
	for _, d := range devices {
		t.Add(d.Name, d.Manufacturer, d.Status)
	}
	return []*DetailTable{t}, nil
}

var batteryStatuses = map[uint16]string{
	1: "Discharging", 2: "On AC power", 3: "Fully charged", 4: "Low", 5: "Critical",
	6: "Charging", 7: "Charging (high)", 8: "Charging (low)", 9: "Charging (critical)",
	10: "Undefined", 11: "Partially charged",
}

var batteryChemistries = map[uint16]string{
	3: "Lead Acid", 4: "NiCd", 5: "NiMH", 6: "Li-ion", 7: "Zinc air", 8: "Li-polymer",
}

func getBatteryDetails() ([]*DetailTable, error) {
	t := newKeyValueTable("Batteries")

	var batteries []struct {
		Name                     string
		EstimatedChargeRemaining uint16
		EstimatedRunTime         uint32
		BatteryStatus            uint16
		Chemistry                uint16
	}
	if err := wmiQuery("SELECT Name, EstimatedChargeRemaining, EstimatedRunTime, BatteryStatus, Chemistry FROM Win32_Battery", &batteries); err != nil {
		return nil, err
	}

	// Capacity data lives in root\wmi and may be unavailable on some systems
	const wmiNamespace = `root\wmi`
	var static []struct {
		DesignedCapacity uint32
		ManufactureName  string
	}
	var full []struct{ FullChargedCapacity uint32 }
	var cycles []struct{ CycleCount uint32 }
	_ = wmiQuery("SELECT DesignedCapacity, ManufactureName FROM BatteryStaticData", &static, wmiNamespace)
	_ = wmiQuery("SELECT FullChargedCapacity FROM BatteryFullChargedCapacity", &full, wmiNamespace)
	_ = wmiQuery("SELECT CycleCount FROM BatteryCycleCount", &cycles, wmiNamespace)

	for i, b := range batteries {
		t.Section(b.Name)
		t.AddValue("Chemistry", batteryChemistries[b.Chemistry])
		t.AddValue("Status", batteryStatuses[b.BatteryStatus])
		t.Add("Charge", fmt.Sprintf("%d %%", b.EstimatedChargeRemaining))
		// 71582788 minutes means "on AC power / unknown"
		if b.EstimatedRunTime > 0 && b.EstimatedRunTime < 71582788 {
			t.Add("Estimated run time", fmt.Sprintf("%dh %dm", b.EstimatedRunTime/60, b.EstimatedRunTime%60))
		}

		var design, fullCap float64
		if i < len(static) {
			t.AddValue("Manufacturer", static[i].ManufactureName)
			design = float64(static[i].DesignedCapacity)
		}
		if i < len(full) {
			fullCap = float64(full[i].FullChargedCapacity)
		}
		// Capacities are in mWh
		if fullCap > 0 {
			t.Add("Full capacity", fmt.Sprintf("%.2f Wh", fullCap/1000))
		}
		if design > 0 {
			t.Add("Design capacity", fmt.Sprintf("%.2f Wh", design/1000))
		}
		if fullCap > 0 && design > 0 {
			t.Add("Health", fmt.Sprintf("%.0f %% (wear %.0f %%)", fullCap/design*100, math.Max(0, 100-fullCap/design*100)))
		}
		if i < len(cycles) && cycles[i].CycleCount > 0 {
			t.Add("Cycle count", strconv.FormatUint(uint64(cycles[i].CycleCount), 10))
		}
	}
	return []*DetailTable{t}, nil
}
