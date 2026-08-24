package system

import (
	"bufio"
	"os"
	"runtime"
	"strconv"
	"strings"
)

func getCommonInfo() ([]DataItem, error) {
	items := make([]DataItem, 0)

	items = append(items, DataItem{Name: "Cores", Value: strconv.Itoa(runtime.NumCPU())})

	f, err := os.Open("/proc/cpuinfo")
	if err != nil {
		return nil, err
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
		return nil, err
	}

	for i, core := range cores {
		items = append(items, DataItem{Name: "Core " + strconv.Itoa(i) + " Model Name", Value: core.ModelName})
	}

	memItems, err := getMemoryInfo()
	if err != nil {
		return nil, err
	}
	items = append(items, memItems...)

	return items, nil
}

func getMemoryInfo() ([]DataItem, error) {
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		return nil, err
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
		return nil, err
	}

	totalMemGBStr := strconv.FormatFloat(float64(total)/1024/1024/1024, 'f', 2, 64) + " GB"
	availableMemGBStr := strconv.FormatFloat(float64(available)/1024/1024/1024, 'f', 2, 64) + " GB"

	items := make([]DataItem, 0)
	items = append(items, DataItem{Name: "Total Memory", Value: totalMemGBStr})
	items = append(items, DataItem{Name: "Available Memory", Value: availableMemGBStr})

	return items, nil
}
