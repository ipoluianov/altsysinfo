package system

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/shirou/gopsutil/v4/cpu"
	"golang.org/x/sys/unix"
)

func addPlatformHardwareDetails(t *DetailTable) {
	addProfilerSections(t, systemProfiler("SPHardwareDataType"))
}

func addPlatformOSDetails(t *DetailTable) {
	for _, item := range systemProfiler("SPSoftwareDataType") {
		addProfilerItem(t, item)
	}
}

func addPlatformCPUDetails(t *DetailTable, _ []cpu.InfoStat) {
	if n, err := unix.SysctlUint32("hw.packages"); err == nil && n > 0 {
		t.Add("Sockets", fmt.Sprint(n))
	}

	// Apple Silicon: performance levels (P-cores / E-cores)
	for level := 0; level < 4; level++ {
		prefix := fmt.Sprintf("hw.perflevel%d.", level)
		name, err := unix.Sysctl(prefix + "name")
		if err != nil {
			break
		}
		cores, _ := unix.SysctlUint32(prefix + "physicalcpu")
		t.Add(name+" cores", fmt.Sprint(cores))
	}

	if hz, err := unix.SysctlUint64("hw.cpufrequency_max"); err == nil && hz > 0 {
		t.Section("Frequency")
		t.Add("Maximum", fmt.Sprintf("%d MHz", hz/1000000))
	}

	t.Section("Cache")
	for _, c := range []struct{ name, key string }{
		{"L1 Instruction", "hw.l1icachesize"},
		{"L1 Data", "hw.l1dcachesize"},
		{"L2", "hw.l2cachesize"},
		{"L3", "hw.l3cachesize"},
	} {
		if size, err := unix.SysctlUint64(c.key); err == nil && size > 0 {
			t.Add(c.name, formatBytes(size))
		}
	}
}

func getGraphicsDetails() ([]*DetailTable, error) {
	return profilerTable("Graphics & Displays", "SPDisplaysDataType"), nil
}

// platformInterfaceExtras maps BSD device names to hardware port names
// using "networksetup -listallhardwareports".
func platformInterfaceExtras() map[string]interfaceExtra {
	result := make(map[string]interfaceExtra)
	out, err := exec.Command("networksetup", "-listallhardwareports").Output()
	if err != nil {
		return result
	}
	port := ""
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		key, value, ok := strings.Cut(scanner.Text(), ":")
		if !ok {
			continue
		}
		switch strings.TrimSpace(key) {
		case "Hardware Port":
			port = strings.TrimSpace(value)
		case "Device":
			result[strings.TrimSpace(value)] = interfaceExtra{Description: port}
		}
	}
	return result
}

func getPCIDetails() ([]*DetailTable, error) {
	return profilerTable("PCI devices", "SPPCIDataType"), nil
}

func getUSBDetails() ([]*DetailTable, error) {
	// SPUSBHostDataType replaces SPUSBDataType in newer macOS versions
	return profilerTable("USB devices", "SPUSBDataType", "SPUSBHostDataType"), nil
}

func getAudioDetails() ([]*DetailTable, error) {
	return profilerTable("Sound devices", "SPAudioDataType"), nil
}

func getBatteryDetails() ([]*DetailTable, error) {
	return profilerTable("Power", "SPPowerDataType"), nil
}
