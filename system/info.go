package system

type Info struct {
	CpuInfo CpuInfo
	RamInfo RamInfo
	Drives  []DriveInfo
	GPUs    []GPUInfo
}

type CpuInfo struct {
	Cores    int
	ModelStr string
}

type RamInfo struct {
	Total uint64
	Used  uint64
	Free  uint64
}

type DriveInfo struct {
	Name       string
	Model      string
	Vendor     string
	Type       string
	Size       uint64 // bytes
	Removable  bool
	Rotational bool
}

type GPUInfo struct {
	Name   string
	Vendor string
	Device string
	Driver string
	Model  string
}
