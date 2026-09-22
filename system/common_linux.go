package system

import (
	"bufio"
	"encoding/binary"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func getInfo() (info Info, err error) {
	info.CpuInfo, err = getCpuInfo()
	if err != nil {
		return
	}
	info.RamInfo, err = getRamInfo()
	if err != nil {
		return
	}
	info.Drives, err = getDrives()
	if err != nil {
		return
	}
	info.GPUs, err = getGPUs()
	if err != nil {
		return
	}
	info.RamDevices = getMemoryModules()
	return
}

func getCpuInfo() (CpuInfo, error) {
	var info CpuInfo
	f, err := os.Open("/proc/cpuinfo")
	if err != nil {
		return info, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	type Core struct {
		ModelName string
		VendorID  string
		CPUMHz    string
	}

	cores := make([]*Core, 0)
	var core Core

	for scanner.Scan() {
		line := scanner.Text()
		key, value, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)

		if key == "processor" {
			if core.ModelName != "" || core.VendorID != "" || core.CPUMHz != "" {
				cores = append(cores, &core)
			}
			core = Core{}
		}

		switch key {
		case "model name":
			core.ModelName = value
		case "vendor_id":
			core.VendorID = value
		case "cpu MHz":
			core.CPUMHz = value
		}
	}

	cores = append(cores, &core)

	if err := scanner.Err(); err != nil {
		return info, err
	}

	cpuModels := make(map[string]struct{})
	for _, core := range cores {
		cpuModels[core.ModelName] = struct{}{}
	}

	info.Cores = len(cores)
	for model := range cpuModels {
		if info.ModelStr != "" {
			info.ModelStr += "; "
		}
		info.ModelStr += model
	}

	return info, nil
}

func getRamInfo() (RamInfo, error) {
	var info RamInfo
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		return info, err
	}
	defer f.Close()

	var total, available uint64

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		key, value, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}

		fields := strings.Fields(value)
		if len(fields) == 0 {
			continue
		}

		v, err := strconv.ParseUint(fields[0], 10, 64)
		if err != nil {
			continue
		}

		// /proc/meminfo values are usually in kB
		v *= 1024

		switch key {
		case "MemTotal":
			total = v

		case "MemAvailable":
			available = v
		}
	}

	if err := scanner.Err(); err != nil {
		return info, err
	}

	info.Total = total
	info.Used = total - available
	info.Free = available

	return info, nil
}

func readString(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

func readUint(path string) uint64 {
	s := readString(path)
	v, _ := strconv.ParseUint(s, 10, 64)
	return v
}

func getDrives() ([]DriveInfo, error) {
	entries, err := os.ReadDir("/sys/block")
	if err != nil {
		return nil, err
	}

	var result []DriveInfo

	for _, entry := range entries {
		name := entry.Name()

		if strings.HasPrefix(name, "loop") ||
			strings.HasPrefix(name, "ram") ||
			strings.HasPrefix(name, "zram") {
			continue
		}

		base := filepath.Join("/sys/block", name)

		sectors := readUint(filepath.Join(base, "size")) // number of 512-byte sectors

		var tp string
		tp = "HDD"
		if ok, _ := isSSD(name); ok {
			tp = "SSD"
		}

		drive := DriveInfo{
			Name:       name,
			Model:      readString(filepath.Join(base, "device/model")),
			Vendor:     readString(filepath.Join(base, "device/vendor")),
			Size:       sectors * 512,
			Removable:  readUint(filepath.Join(base, "removable")) != 0,
			Rotational: readUint(filepath.Join(base, "queue/rotational")) != 0,
			Type:       tp,
		}

		if isOpticalDrive(name) {
			continue
		}

		result = append(result, drive)
	}

	return result, nil
}

func isOpticalDrive(name string) bool {
	data, err := os.ReadFile("/sys/block/" + name + "/device/type")
	if err != nil {
		return false
	}

	return strings.TrimSpace(string(data)) == "5"
}

func isSSD(name string) (bool, error) {
	data, err := os.ReadFile("/sys/block/" + name + "/queue/rotational")
	if err != nil {
		return false, err
	}

	return strings.TrimSpace(string(data)) == "0", nil
}

type GPU struct {
	Name     string
	VendorID uint16
	DeviceID uint16
	Driver   string
}

func readHex16(path string) uint16 {
	s := readString(path)
	s = strings.TrimPrefix(s, "0x")

	v, err := strconv.ParseUint(s, 16, 16)
	if err != nil {
		return 0
	}

	return uint16(v)
}

func getGPUs() ([]GPUInfo, error) {
	entries, err := os.ReadDir("/sys/class/drm")
	if err != nil {
		return nil, err
	}

	var result []GPUInfo

	for _, entry := range entries {
		name := entry.Name()

		if !strings.HasPrefix(name, "card") {
			continue
		}

		index := strings.TrimPrefix(name, "card")
		if _, err := strconv.Atoi(index); err != nil {
			continue
		}

		base := filepath.Join("/sys/class/drm", name, "device")

		vendorPath := filepath.Join(base, "vendor")
		devicePath := filepath.Join(base, "device")
		if _, err := os.Stat(vendorPath); err != nil {
			continue
		}

		gpu := GPUInfo{
			Name:   name,
			Vendor: readString(vendorPath),
			Device: readString(devicePath),
		}

		driverLink := filepath.Join(base, "driver")
		if path, err := filepath.EvalSymlinks(driverLink); err == nil {
			gpu.Driver = filepath.Base(path)
		}

		model, ok := PCIName(readHex16(vendorPath), readHex16(devicePath))
		if ok {
			gpu.Model = model
		} else {
			gpu.Model = "Unknown" + " (" + gpu.Vendor + ":" + gpu.Device + ")"
		}

		result = append(result, gpu)
	}

	return result, nil
}

// getMemoryModules returns installed memory modules (SMBIOS type 17).
// The udev database is readable without root; the raw SMBIOS table
// is used as a fallback (requires root).
func getMemoryModules() []RamDeviceInfo {
	if modules := getMemoryModulesUdev(); len(modules) > 0 {
		return modules
	}
	return getMemoryModulesSMBIOS()
}

func getMemoryModulesUdev() []RamDeviceInfo {
	f, err := os.Open("/run/udev/data/+dmi:id")
	if err != nil {
		return nil
	}
	defer f.Close()

	devices := make(map[int]map[string]string)
	maxIndex := -1

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line, ok := strings.CutPrefix(scanner.Text(), "E:MEMORY_DEVICE_")
		if !ok {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		indexStr, field, ok := strings.Cut(key, "_")
		if !ok {
			continue
		}
		index, err := strconv.Atoi(indexStr)
		if err != nil {
			continue
		}
		if devices[index] == nil {
			devices[index] = make(map[string]string)
		}
		devices[index][field] = value
		maxIndex = max(maxIndex, index)
	}

	var result []RamDeviceInfo
	for i := 0; i <= maxIndex; i++ {
		dev := devices[i]
		if dev == nil || dev["PRESENT"] == "0" {
			continue
		}
		size, _ := strconv.ParseUint(dev["SIZE"], 10, 64)
		if size == 0 {
			continue
		}
		speed, _ := strconv.ParseUint(dev["CONFIGURED_SPEED_MTS"], 10, 64)
		if speed == 0 {
			speed, _ = strconv.ParseUint(dev["SPEED_MTS"], 10, 64)
		}
		result = append(result, RamDeviceInfo{
			Model: dev["PART_NUMBER"],
			Size:  size,
			Speed: speed,
		})
	}
	return result
}

func getMemoryModulesSMBIOS() []RamDeviceInfo {
	data, err := os.ReadFile("/sys/firmware/dmi/tables/DMI")
	if err != nil {
		return nil
	}

	var result []RamDeviceInfo

	for len(data) >= 4 {
		structType := data[0]
		length := int(data[1])
		if length < 4 || length > len(data) {
			break
		}
		formatted := data[:length]

		// Strings follow the formatted area and end with a double null
		end := length
		for end+1 < len(data) && (data[end] != 0 || data[end+1] != 0) {
			end++
		}
		strs := strings.Split(string(data[length:end]), "\x00")
		data = data[min(end+2, len(data)):]

		if structType == 127 { // end of table
			break
		}
		if structType != 17 || length < 0x1B {
			continue
		}

		getString := func(offset int) string {
			idx := int(formatted[offset])
			if idx == 0 || idx > len(strs) {
				return ""
			}
			return strings.TrimSpace(strs[idx-1])
		}
		word := func(offset int) uint64 {
			if offset+2 > length {
				return 0
			}
			return uint64(binary.LittleEndian.Uint16(formatted[offset:]))
		}
		dword := func(offset int) uint64 {
			if offset+4 > length {
				return 0
			}
			return uint64(binary.LittleEndian.Uint32(formatted[offset:]))
		}

		var size uint64
		switch rawSize := word(0x0C); {
		case rawSize == 0 || rawSize == 0xFFFF: // not installed / unknown
			continue
		case rawSize == 0x7FFF:
			size = (dword(0x1C) & 0x7FFFFFFF) * 1024 * 1024
		case rawSize&0x8000 != 0:
			size = (rawSize & 0x7FFF) * 1024
		default:
			size = rawSize * 1024 * 1024
		}

		speed := word(0x20)
		if speed == 0xFFFF {
			speed = dword(0x58)
		}
		if speed == 0 {
			speed = word(0x15)
			if speed == 0xFFFF {
				speed = dword(0x54)
			}
		}

		result = append(result, RamDeviceInfo{
			Model: getString(0x1A),
			Size:  size,
			Speed: speed,
		})
	}

	return result
}
