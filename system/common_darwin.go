package system

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
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
	info.Cores = runtime.NumCPU()
	infos, err := cpu.Info()
	if err != nil {
		return info, err
	}
	if len(infos) > 0 {
		info.ModelStr = infos[0].ModelName
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
	info.Used = m.Total - m.Available
	info.Free = m.Available
	return info, nil
}

// getDrives groups the volumes reported by system_profiler by physical drive.
// The size is the largest container on the drive, which is close to the
// drive capacity.
func getDrives() ([]DriveInfo, error) {
	byName := make(map[string]*DriveInfo)
	var names []string
	for _, volume := range systemProfiler("SPStorageDataType") {
		drive, _ := volume["physical_drive"].(map[string]any)
		name := profilerString(drive["device_name"])
		if name == "" {
			continue
		}
		d, ok := byName[name]
		if !ok {
			d = &DriveInfo{Name: name, Model: name, Type: strings.ToUpper(profilerString(drive["medium_type"]))}
			byName[name] = d
			names = append(names, name)
		}
		if size := profilerUint(volume["size_in_bytes"]); size > d.Size {
			d.Size = size
		}
		d.Removable = profilerString(drive["is_internal_disk"]) == "no"
	}

	result := make([]DriveInfo, 0, len(names))
	for _, name := range names {
		result = append(result, *byName[name])
	}
	return result, nil
}

func getGPUs() ([]GPUInfo, error) {
	var result []GPUInfo
	for _, item := range systemProfiler("SPDisplaysDataType") {
		model := profilerString(item["sppci_model"])
		if model == "" {
			model = profilerString(item["_name"])
		}
		result = append(result, GPUInfo{
			Name:   model,
			Model:  model,
			Vendor: profilerString(item["spdisplays_vendor-id"]),
			Device: profilerString(item["spdisplays_device-id"]),
		})
	}
	return result, nil
}

func getMemoryModules() []RamDeviceInfo {
	var result []RamDeviceInfo
	for _, item := range systemProfiler("SPMemoryDataType") {
		// Apple Silicon: unified memory is reported as a single item
		if size := profilerString(item["SPMemoryDataType"]); size != "" {
			result = append(result, RamDeviceInfo{
				Model: strings.TrimSpace(profilerString(item["dimm_manufacturer"]) + " " + profilerString(item["dimm_type"])),
				Size:  parseSizeString(size),
			})
			continue
		}
		// Intel: memory slots
		slots, _ := item["_items"].([]any)
		for _, s := range slots {
			slot, ok := s.(map[string]any)
			if !ok {
				continue
			}
			size := parseSizeString(profilerString(slot["dimm_size"]))
			if size == 0 {
				continue
			}
			model := profilerString(slot["dimm_part_number"])
			if model == "" || model == "-" {
				model = profilerString(slot["dimm_type"])
			}
			speed, _ := strconv.ParseUint(strings.Fields(profilerString(slot["dimm_speed"]) + " 0")[0], 10, 64)
			result = append(result, RamDeviceInfo{Model: model, Size: size, Speed: speed})
		}
	}
	return result
}

// parseSizeString parses sizes like "16 GB".
func parseSizeString(s string) uint64 {
	fields := strings.Fields(s)
	if len(fields) != 2 {
		return 0
	}
	v, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return 0
	}
	multipliers := map[string]float64{"KB": 1 << 10, "MB": 1 << 20, "GB": 1 << 30, "TB": 1 << 40}
	return uint64(v * multipliers[strings.ToUpper(fields[1])])
}

// systemProfiler returns the items of a system_profiler data type.
func systemProfiler(dataType string) []map[string]any {
	out, err := exec.Command("system_profiler", "-json", dataType).Output()
	if err != nil {
		return nil
	}
	var result map[string][]map[string]any
	if err := json.Unmarshal(out, &result); err != nil {
		return nil
	}
	return result[dataType]
}

func profilerString(v any) string {
	switch v := v.(type) {
	case nil:
		return ""
	case string:
		return v
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		return yesNo(v)
	default:
		return fmt.Sprint(v)
	}
}

func profilerUint(v any) uint64 {
	switch v := v.(type) {
	case float64:
		return uint64(v)
	case string:
		n, _ := strconv.ParseUint(v, 10, 64)
		return n
	}
	return 0
}

var (
	// Data type prefixes are dropped ("spdisplays_vram" -> "vram"); for other
	// "sp" identifiers only "sp" is dropped ("spbattery_information").
	profilerPrefix     = regexp.MustCompile(`^(sppci_vendor|sppci|spdisplays|sppower|spusb|spaudio|spnvme|spsata|spstorage|spnetwork|spmemory|sphardware|spsoftware)_|^sp([a-z0-9]+_)`)
	profilerIdentifier = regexp.MustCompile(`^[a-z0-9]+(_[a-z0-9]+)+$`)
)

// profilerName turns identifiers like "spdisplays_vram" into "Vram".
func profilerName(s string) string {
	s = profilerPrefix.ReplaceAllString(s, "$2")
	s = strings.TrimSuffix(s, "_in_bytes")
	s = strings.NewReplacer("_", " ", "-", " ").Replace(s)
	r := []rune(s)
	if len(r) > 0 {
		r[0] = unicode.ToUpper(r[0])
	}
	return string(r)
}

// profilerValue formats a value; enum-like identifiers are made readable.
func profilerValue(key string, v any) string {
	if strings.HasSuffix(key, "_in_bytes") || strings.HasSuffix(key, "_bytes") {
		if n := profilerUint(v); n > 0 {
			return formatBytes(n)
		}
	}
	s := profilerString(v)
	if profilerPrefix.MatchString(s) || profilerIdentifier.MatchString(s) {
		return profilerName(s)
	}
	return s
}

// addProfilerItem adds the values of a system_profiler item to a table.
// Nested item lists become separate sections.
func addProfilerItem(t *DetailTable, item map[string]any) {
	keys := make([]string, 0, len(item))
	for k := range item {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var children []map[string]any
	for _, k := range keys {
		if strings.HasPrefix(k, "_") && k != "_items" {
			continue
		}
		switch v := item[k].(type) {
		case []any:
			for _, e := range v {
				if m, ok := e.(map[string]any); ok {
					children = append(children, m)
				}
			}
		case map[string]any:
			nestedKeys := make([]string, 0, len(v))
			for nk := range v {
				nestedKeys = append(nestedKeys, nk)
			}
			sort.Strings(nestedKeys)
			for _, nk := range nestedKeys {
				if _, isMap := v[nk].(map[string]any); !isMap {
					t.AddValue(profilerName(nk), profilerValue(nk, v[nk]))
				}
			}
		default:
			t.AddValue(profilerName(k), profilerValue(k, v))
		}
	}

	addProfilerSections(t, children)
}

// addProfilerSections adds each item as a section titled by its name.
func addProfilerSections(t *DetailTable, items []map[string]any) {
	for _, item := range items {
		title := profilerValue("_name", item["_name"])
		if strings.HasPrefix(title, "kHW_") { // internal names of built-in devices
			title = profilerString(item["sppci_model"])
		}
		if title != "" {
			t.Section(title)
		}
		addProfilerItem(t, item)
	}
}

// profilerTable builds a key/value table from system_profiler data types;
// the first data type that returns items is used.
func profilerTable(title string, dataTypes ...string) []*DetailTable {
	t := newKeyValueTable(title)
	for _, dataType := range dataTypes {
		if items := systemProfiler(dataType); len(items) > 0 {
			addProfilerSections(t, items)
			break
		}
	}
	return []*DetailTable{t}
}
