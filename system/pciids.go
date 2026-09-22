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

var (
	pciExtraOnce sync.Once
	pciVendors   map[uint16]string
	pciClasses   map[uint16]string // class<<8 | subclass; subclass 0xff is unused
)

// PCIVendorName returns the vendor name for a PCI vendor ID.
func PCIVendorName(vendorID uint16) (string, bool) {
	pciExtraOnce.Do(func() {
		pciVendors, pciClasses = parsePCIVendorsAndClasses(pciIDs)
	})
	name, ok := pciVendors[vendorID]
	return name, ok
}

// PCIClassName returns the name of a PCI device class/subclass.
func PCIClassName(class, subclass uint8) string {
	pciExtraOnce.Do(func() {
		pciVendors, pciClasses = parsePCIVendorsAndClasses(pciIDs)
	})
	if name, ok := pciClasses[uint16(class)<<8|uint16(subclass)]; ok {
		return name
	}
	if name, ok := pciClasses[uint16(class)<<8|0xff]; ok {
		return name
	}
	return fmt.Sprintf("Class %02x%02x", class, subclass)
}

func parsePCIVendorsAndClasses(data string) (map[uint16]string, map[uint16]string) {
	vendors := make(map[uint16]string)
	classes := make(map[uint16]string)
	var class uint64
	inClass := false

	scanner := bufio.NewScanner(strings.NewReader(data))
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || line[0] == '#' {
			continue
		}

		switch {
		case strings.HasPrefix(line, "C "):
			fields := strings.Fields(line)
			if len(fields) < 3 {
				inClass = false
				continue
			}
			var err error
			class, err = strconv.ParseUint(fields[1], 16, 8)
			inClass = err == nil
			if inClass {
				classes[uint16(class)<<8|0xff] = strings.TrimSpace(line[len("C ")+len(fields[1]):])
			}

		case line[0] != '\t':
			inClass = false
			fields := strings.Fields(line)
			if id, ok := parsePCIHexID(fields); ok {
				vendors[id] = strings.TrimSpace(line[len(fields[0]):])
			}

		case inClass && !strings.HasPrefix(line, "\t\t"):
			fields := strings.Fields(line)
			if len(fields) < 2 {
				continue
			}
			subclass, err := strconv.ParseUint(fields[0], 16, 8)
			if err != nil {
				continue
			}
			classes[uint16(class)<<8|uint16(subclass)] = strings.TrimSpace(line[1+len(fields[0]):])
		}
	}

	return vendors, classes
}
