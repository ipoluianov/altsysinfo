package system

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/StackExchange/wmi"
	"github.com/shirou/gopsutil/v4/mem"
	"golang.org/x/sys/windows/registry"
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
	info.RamDevices, err = getMemoryModules()
	if err != nil {
		return
	}
	return
}

func getCpuInfo() (CpuInfo, error) {
	var info CpuInfo

	info.Cores = runtime.NumCPU()

	key, err := registry.OpenKey(
		registry.LOCAL_MACHINE,
		`HARDWARE\DESCRIPTION\System\CentralProcessor`,
		registry.READ,
	)
	if err != nil {
		return info, err
	}
	defer key.Close()

	subKeys, err := key.ReadSubKeyNames(-1)
	if err != nil {
		return info, err
	}

	models := make(map[string]struct{})

	for _, subKeyName := range subKeys {
		subKey, err := registry.OpenKey(
			key,
			subKeyName,
			registry.READ,
		)
		if err != nil {
			continue
		}

		model, _, err := subKey.GetStringValue("ProcessorNameString")
		subKey.Close()

		if err == nil && model != "" {
			models[model] = struct{}{}
		}
	}

	for model := range models {
		if info.ModelStr != "" {
			info.ModelStr += "; "
		}
		info.ModelStr += model
	}

	return info, nil
}

func getRamInfo() (RamInfo, error) {
	var info RamInfo
	m, err := mem.VirtualMemory()
	if err != nil {
		return info, err
	}

	info.Total = m.Total
	info.Used = m.Used
	info.Free = m.Free

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
	type physicalDisk struct {
		FriendlyName string
		SerialNumber string
		MediaType    uint16
		BusType      uint16
		Size         uint64
	}

	var src []physicalDisk

	query := `
		SELECT FriendlyName, SerialNumber, MediaType, BusType, Size
		FROM MSFT_PhysicalDisk
	`

	err := wmi.QueryNamespace(
		query,
		&src,
		`root\Microsoft\Windows\Storage`,
	)
	if err != nil {
		return nil, err
	}

	drives := make([]DriveInfo, 0, len(src))

	/*
		0 = Unspecified
		3 = HDD
		4 = SSD
		5 = SCM
	*/

	diskTypes := map[uint16]string{
		0: "Unspecified",
		3: "HDD",
		4: "SSD",
		5: "SCM",
	}

	for _, d := range src {

		diskType, ok := diskTypes[d.MediaType]
		if !ok {
			diskType = "Unknown"
		}

		drives = append(drives, DriveInfo{
			Name:   d.FriendlyName,
			Model:  d.FriendlyName,
			Vendor: "",
			Type:   diskType,
			Size:   d.Size,
		})
	}

	return drives, nil
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
	type win32VideoController struct {
		Name          string
		PNPDeviceID   string
		DriverVersion string
		AdapterRAM    uint32
	}

	var src []win32VideoController

	query := `
		SELECT Name, PNPDeviceID, DriverVersion, AdapterRAM
		FROM Win32_VideoController
	`

	if err := wmi.Query(query, &src); err != nil {
		return nil, err
	}

	result := make([]GPUInfo, 0, len(src))

	for _, gpu := range src {
		ven, dev := parsePCIID(gpu.PNPDeviceID)
		result = append(result, GPUInfo{
			Name:   strings.TrimSpace(gpu.Name),
			Vendor: ven,
			Device: dev,
			Driver: strings.TrimSpace(gpu.DriverVersion),
			Model:  strings.TrimSpace(gpu.Name),
		})
	}

	return result, nil
}

func parsePCIID(pnpDeviceID string) (vendorID, deviceID string) {
	parts := strings.Split(pnpDeviceID, `\`)
	if len(parts) < 2 {
		return "", ""
	}

	for _, part := range strings.Split(parts[1], "&") {
		switch {
		case strings.HasPrefix(part, "VEN_"):
			vendorID = strings.TrimPrefix(part, "VEN_")

		case strings.HasPrefix(part, "DEV_"):
			deviceID = strings.TrimPrefix(part, "DEV_")
		}
	}

	return
}

func getMemoryModules() ([]RamDeviceInfo, error) {
	type Win32_PhysicalMemory struct {
		BankLabel            string
		DeviceLocator        string
		Manufacturer         string
		PartNumber           string
		SerialNumber         string
		Capacity             uint64
		Speed                uint32
		ConfiguredClockSpeed uint32
		MemoryType           uint16
		SMBIOSMemoryType     uint32
		FormFactor           uint16
	}

	var dst []Win32_PhysicalMemory

	query := wmi.CreateQuery(&dst, "")
	err := wmi.Query(query, &dst)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	result := make([]RamDeviceInfo, 0, len(dst))
	for _, mem := range dst {
		result = append(result, RamDeviceInfo{
			Model: mem.PartNumber,
			Size:  mem.Capacity,
			Speed: uint64(mem.ConfiguredClockSpeed),
		})
	}

	return result, nil
}
