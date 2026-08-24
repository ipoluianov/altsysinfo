package system

import "testing"

func TestParsePCIIDs(t *testing.T) {
	data := `# comment
1234  Vendor
	0001  First device
		1234 0002  Subsystem
	0002  Second device
C 01  Mass storage controller
	06  SATA controller
abcd  Another vendor
	00ff  Another device
`

	devices := parsePCIIDs(data)
	tests := []struct {
		vendor uint16
		device uint16
		want   string
	}{
		{0x1234, 0x0001, "First device"},
		{0x1234, 0x0002, "Second device"},
		{0xabcd, 0x00ff, "Another device"},
	}

	for _, test := range tests {
		got, ok := devices[pciDeviceID{vendor: test.vendor, device: test.device}]
		if !ok || got.DeviceName != test.want {
			t.Errorf("lookup %04x:%04x = %q, %v; want %q, true", test.vendor, test.device, got.DeviceName, ok, test.want)
		}
	}

	got := devices[pciDeviceID{vendor: 0xabcd, device: 0x00ff}]
	if got.VendorName != "Another vendor" {
		t.Errorf("VendorName = %q; want %q", got.VendorName, "Another vendor")
	}

	if _, ok := devices[pciDeviceID{vendor: 0x1234, device: 0x0002}]; !ok {
		t.Fatal("device following a subsystem was not parsed")
	}
	if _, ok := devices[pciDeviceID{vendor: 0x0001, device: 0x0006}]; ok {
		t.Fatal("PCI class entry was parsed as a device")
	}
}

func TestPCIName(t *testing.T) {
	name, ok := PCIName(0x0010, 0x8139)
	if !ok || name != "AT-2500TX V3 Ethernet" {
		t.Fatalf("PCIName() = %q, %v; want %q, true", name, ok, "AT-2500TX V3 Ethernet")
	}

	if name, ok := PCIName(0xffff, 0xffff); ok || name != "" {
		t.Fatalf("unknown PCIName() = %q, %v; want empty, false", name, ok)
	}
}
