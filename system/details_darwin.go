package system

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
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

// platformCPUFeatures reads the ARM features from sysctl: golang.org/x/sys/cpu
// sees only NEON and FP on macOS.
func platformCPUFeatures() []string {
	if runtime.GOARCH != "arm64" {
		return nil
	}
	var result []string
	for _, f := range []struct{ name, key string }{
		{"NEON", "hw.optional.neon"},
		{"FP", "hw.optional.floatingpoint"},
		{"AES", "hw.optional.arm.FEAT_AES"},
		{"PMULL", "hw.optional.arm.FEAT_PMULL"},
		{"SHA1", "hw.optional.arm.FEAT_SHA1"},
		{"SHA2", "hw.optional.arm.FEAT_SHA256"},
		{"SHA3", "hw.optional.arm.FEAT_SHA3"},
		{"SHA512", "hw.optional.arm.FEAT_SHA512"},
		{"CRC32", "hw.optional.arm.FEAT_CRC32"},
		{"LSE Atomics", "hw.optional.arm.FEAT_LSE"},
		{"DotProd", "hw.optional.arm.FEAT_DotProd"},
		{"FP16", "hw.optional.arm.FEAT_FP16"},
		{"BF16", "hw.optional.arm.FEAT_BF16"},
		{"I8MM", "hw.optional.arm.FEAT_I8MM"},
		{"SME", "hw.optional.arm.FEAT_SME"},
		{"SME2", "hw.optional.arm.FEAT_SME2"},
	} {
		if v, err := unix.SysctlUint32(f.key); err == nil && v != 0 {
			result = append(result, f.name)
		}
	}
	return result
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
