package system

import (
	"bufio"
	_ "embed"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
)

//go:embed pci.ids
var pciIDs string

type pciDeviceID struct {
	vendor uint16
	device uint16
}

var (
	pciDevicesOnce sync.Once
	pciDevices     map[pciDeviceID]PciInfo
)

// PCIName returns the device name for a PCI vendor/device pair.
// IDs are the 16-bit hexadecimal values exposed by PCI configuration space.
func PCIName(vendorID, deviceID uint16) (string, bool) {
	pciDevicesOnce.Do(func() {
		pciDevices = parsePCIIDs(pciIDs)
	})

	info, ok := pciDevices[pciDeviceID{vendor: vendorID, device: deviceID}]
	return info.DeviceName, ok
}

type PciInfo struct {
	VendorID   uint16
	DeviceID   uint16
	VendorName string
	DeviceName string
}

func (c *PciInfo) VendorIDStr() string {
	return fmt.Sprintf("%04x", c.VendorID)
}

func (c *PciInfo) DeviceIDStr() string {
	return fmt.Sprintf("%04x", c.DeviceID)
}

func PCIList() []PciInfo {
	pciDevicesOnce.Do(func() {
		pciDevices = parsePCIIDs(pciIDs)
	})

	var list []PciInfo
	for _, info := range pciDevices {
		list = append(list, info)
	}
	sort.Slice(list, func(i, j int) bool {
		if list[i].VendorID != list[j].VendorID {
			return list[i].VendorID < list[j].VendorID
		}
		return list[i].DeviceID < list[j].DeviceID
	})
	return list
}

func parsePCIIDs(data string) map[pciDeviceID]PciInfo {
	devices := make(map[pciDeviceID]PciInfo)
	var vendorID uint16
	var vendorName string
	var haveVendor bool

	scanner := bufio.NewScanner(strings.NewReader(data))
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || line[0] == '#' {
			continue
		}

		switch {
		case line[0] != '\t':
			// Vendor records begin with exactly four hexadecimal digits. This
			// also excludes class records, which begin with "C ".
			fields := strings.Fields(line)
			id, ok := parsePCIHexID(fields)
			vendorID, haveVendor = id, ok
			if ok {
				vendorName = strings.TrimSpace(line[len(fields[0]):])
			} else {
				vendorName = ""
			}

		case haveVendor && !strings.HasPrefix(line, "\t\t"):
			fields := strings.Fields(line)
			deviceID, ok := parsePCIHexID(fields)
			if !ok {
				continue
			}

			name := strings.TrimSpace(line[1+len(fields[0]):])
			if name != "" {
				devices[pciDeviceID{vendor: vendorID, device: deviceID}] = PciInfo{
					VendorID:   vendorID,
					DeviceID:   deviceID,
					VendorName: vendorName,
					DeviceName: name,
				}
			}
		}
	}

	return devices
}

func parsePCIHexID(fields []string) (uint16, bool) {
	if len(fields) < 2 || len(fields[0]) != 4 {
		return 0, false
	}

	id, err := strconv.ParseUint(fields[0], 16, 16)
	if err != nil {
		return 0, false
	}

	return uint16(id), true
}
