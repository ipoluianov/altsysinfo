package system

import (
	"fmt"
	"net"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	xcpu "golang.org/x/sys/cpu"
)

// DetailTable is one table on a details page.
// A row with a single cell is a section header.
type DetailTable struct {
	Title   string
	Columns []string
	Rows    [][]string
}

func newDetailTable(title string, columns ...string) *DetailTable {
	return &DetailTable{Title: title, Columns: columns}
}

func newKeyValueTable(title string) *DetailTable {
	return newDetailTable(title, "Name", "Value")
}

// Section adds a section header row.
func (t *DetailTable) Section(title string) {
	t.Rows = append(t.Rows, []string{title})
}

// Add adds a row.
func (t *DetailTable) Add(cells ...string) {
	for i := range cells {
		cells[i] = strings.TrimSpace(cells[i])
	}
	t.Rows = append(t.Rows, cells)
}

// AddValue adds a name/value row; empty values are skipped.
func (t *DetailTable) AddValue(name, value string) {
	if strings.TrimSpace(value) == "" {
		return
	}
	t.Add(name, value)
}

// DetailCategory is a page with detailed information.
type DetailCategory struct {
	ID   string
	Name string
	Load func() ([]*DetailTable, error)
}

func DetailCategories() []DetailCategory {
	return []DetailCategory{
		{ID: "system", Name: "System", Load: getSystemDetails},
		{ID: "cpu", Name: "CPU", Load: getCPUDetails},
		{ID: "storage", Name: "Storage", Load: getStorageDetails},
		{ID: "graphics", Name: "Graphics & Displays", Load: getGraphicsDetails},
		{ID: "network", Name: "Network", Load: getNetworkDetails},
		{ID: "pci", Name: "PCI Devices", Load: getPCIDetails},
		{ID: "usb", Name: "USB Devices", Load: getUSBDetails},
		{ID: "audio", Name: "Audio", Load: getAudioDetails},
		{ID: "sensors", Name: "Sensors", Load: getSensorsDetails},
		{ID: "battery", Name: "Battery", Load: getBatteryDetails},
	}
}

func getSystemDetails() ([]*DetailTable, error) {
	t := newKeyValueTable("System")
	addPlatformHardwareDetails(t)

	t.Section("Operating System")
	if h, err := host.Info(); err == nil {
		t.AddValue("Name", strings.TrimSpace(h.Platform+" "+h.PlatformVersion))
		t.AddValue("Family", h.PlatformFamily)
		t.AddValue("Kernel", h.KernelVersion)
		t.AddValue("Architecture", h.KernelArch)
		t.AddValue("Host name", h.Hostname)
		t.AddValue("Boot time", time.Unix(int64(h.BootTime), 0).Format("2006-01-02 15:04:05"))
		t.AddValue("Uptime", formatDuration(time.Duration(h.Uptime)*time.Second))
		if h.VirtualizationRole == "guest" {
			t.AddValue("Virtual machine", h.VirtualizationSystem)
		}
	}
	addPlatformOSDetails(t)

	return []*DetailTable{t}, nil
}

func getCPUDetails() ([]*DetailTable, error) {
	t := newKeyValueTable("CPU")
	t.Section("Processor")

	infos, err := cpu.Info()
	if err == nil && len(infos) > 0 {
		c := infos[0]
		t.AddValue("Model", c.ModelName)
		t.AddValue("Vendor", c.VendorID)
		if c.Family != "" {
			t.AddValue("Family / Model / Stepping", fmt.Sprintf("%s / %s / %d", c.Family, c.Model, c.Stepping))
		}
		t.AddValue("Microcode", c.Microcode)
		if c.Mhz > 0 {
			t.AddValue("Frequency", fmt.Sprintf("%.0f MHz", c.Mhz))
		}
	}
	t.AddValue("Architecture", runtime.GOARCH)
	if n, err := cpu.Counts(false); err == nil && n > 0 {
		t.AddValue("Physical cores", strconv.Itoa(n))
	}
	if n, err := cpu.Counts(true); err == nil && n > 0 {
		t.AddValue("Logical processors", strconv.Itoa(n))
	}

	addPlatformCPUDetails(t, infos)

	if features := cpuFeatures(); len(features) > 0 {
		t.Section("Instruction set extensions")
		t.Add("Supported", strings.Join(features, ", "))
	}

	return []*DetailTable{t}, nil
}

func cpuFeatures() []string {
	var result []string
	add := func(ok bool, name string) {
		if ok {
			result = append(result, name)
		}
	}
	switch runtime.GOARCH {
	case "amd64", "386":
		x := xcpu.X86
		add(x.HasSSE2, "SSE2")
		add(x.HasSSE3, "SSE3")
		add(x.HasSSSE3, "SSSE3")
		add(x.HasSSE41, "SSE4.1")
		add(x.HasSSE42, "SSE4.2")
		add(x.HasPOPCNT, "POPCNT")
		add(x.HasAVX, "AVX")
		add(x.HasAVX2, "AVX2")
		add(x.HasFMA, "FMA3")
		add(x.HasAVX512F, "AVX-512F")
		add(x.HasAVX512BW, "AVX-512BW")
		add(x.HasAVX512VL, "AVX-512VL")
		add(x.HasAVX512VNNI, "AVX-512VNNI")
		add(x.HasBMI1, "BMI1")
		add(x.HasBMI2, "BMI2")
		add(x.HasAES, "AES-NI")
		add(x.HasPCLMULQDQ, "PCLMULQDQ")
		add(x.HasRDRAND, "RDRAND")
		add(x.HasRDSEED, "RDSEED")
		add(x.HasADX, "ADX")
		add(x.HasERMS, "ERMS")
	case "arm64":
		a := xcpu.ARM64
		add(a.HasASIMD, "NEON")
		add(a.HasFP, "FP")
		add(a.HasAES, "AES")
		add(a.HasPMULL, "PMULL")
		add(a.HasSHA1, "SHA1")
		add(a.HasSHA2, "SHA2")
		add(a.HasSHA3, "SHA3")
		add(a.HasSHA512, "SHA512")
		add(a.HasCRC32, "CRC32")
		add(a.HasATOMICS, "LSE Atomics")
		add(a.HasASIMDDP, "DotProd")
		add(a.HasSVE, "SVE")
		add(a.HasSVE2, "SVE2")
	}
	return result
}

func getStorageDetails() ([]*DetailTable, error) {
	drives := newDetailTable("Physical drives", "Name", "Model", "Type", "Size")
	if list, err := getDrives(); err == nil {
		for _, d := range list {
			model := d.Model
			if d.Vendor != "" && d.Vendor != "ATA" && !strings.HasPrefix(model, d.Vendor) {
				model = d.Vendor + " " + model
			}
			drives.Add(d.Name, model, d.Type, formatBytes(d.Size))
		}
	}

	volumes := newDetailTable("Volumes", "Device", "Mount point", "File system", "Size", "Used", "Free", "Used %")
	if parts, err := disk.Partitions(false); err == nil {
		sort.Slice(parts, func(i, j int) bool { return parts[i].Mountpoint < parts[j].Mountpoint })
		for _, p := range parts {
			if p.Fstype == "squashfs" { // snap packages and other images
				continue
			}
			u, err := disk.Usage(p.Mountpoint)
			if err != nil || u.Total == 0 {
				continue
			}
			volumes.Add(p.Device, p.Mountpoint, p.Fstype,
				formatBytes(u.Total), formatBytes(u.Used), formatBytes(u.Free),
				fmt.Sprintf("%.0f%%", u.UsedPercent))
		}
	}

	return []*DetailTable{drives, volumes}, nil
}

// interfaceExtra holds OS-specific information about a network interface.
type interfaceExtra struct {
	Description string
	Speed       string
	Driver      string
}

func getNetworkDetails() ([]*DetailTable, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	extras := platformInterfaceExtras()

	t := newDetailTable("Network interfaces", "Interface", "Description", "Status", "MAC", "IPv4", "IPv6", "Speed", "Driver")
	for _, iface := range ifaces {
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		var ipv4, ipv6 []string
		if addrs, err := iface.Addrs(); err == nil {
			for _, addr := range addrs {
				ipNet, ok := addr.(*net.IPNet)
				if !ok {
					continue
				}
				if ipNet.IP.To4() != nil {
					ipv4 = append(ipv4, ipNet.String())
				} else {
					ipv6 = append(ipv6, ipNet.String())
				}
			}
		}
		status := "Down"
		if iface.Flags&net.FlagUp != 0 {
			status = "Up"
		}
		extra := extras[iface.Name]
		t.Add(iface.Name, extra.Description, status, iface.HardwareAddr.String(),
			strings.Join(ipv4, ", "), strings.Join(ipv6, ", "), extra.Speed, extra.Driver)
	}
	return []*DetailTable{t}, nil
}

func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return strconv.FormatUint(bytes, 10) + " B"
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func formatDuration(d time.Duration) string {
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	}
	return fmt.Sprintf("%dh %dm", hours, minutes)
}

func yesNo(v bool) string {
	if v {
		return "Yes"
	}
	return "No"
}
