package widgets

import (
	"strconv"

	"github.com/ipoluianov/altsysinfo/system"
	"github.com/u00io/nuiforms/ui"
)

type WidgetCommonInfo struct {
	ui.Widget

	lvItems *ui.Table
}

func NewWidgetCommonInfo() *WidgetCommonInfo {
	var c WidgetCommonInfo
	c.InitWidget()
	c.lvItems = ui.NewTable()
	c.AddWidgetOnGrid(c.lvItems, 0, 0)
	c.lvItems.SetSelectingCell(false)
	c.lvItems.SetColumnCount(2)
	c.lvItems.SetColumnWidth(0, 200)
	c.lvItems.SetColumnWidth(1, 600)
	c.lvItems.SetColumnName(0, "Name")
	c.lvItems.SetColumnName(1, "Value")

	c.LoadCommonInfo()
	return &c
}

func (c *WidgetCommonInfo) AddRow(name, value string) {
	row := c.lvItems.RowCount()
	c.lvItems.SetRowCount(row + 1)
	c.lvItems.SetCellText2(row, 0, name)
	c.lvItems.SetCellText2(row, 1, value)
}

func (c *WidgetCommonInfo) LoadCommonInfo() {
	info, err := system.GetInfo()
	if err != nil {
		return
	}

	c.lvItems.SetRowCount(0)
	c.AddRow("CPU Model", info.CpuInfo.ModelStr)
	c.AddRow("CPU Cores", strconv.FormatInt(int64(info.CpuInfo.Cores), 10))

	c.AddRow("RAM", strconv.FormatUint(info.RamInfo.Total/1024/1024, 10)+" MB")

	for i, drive := range info.Drives {
		c.AddRow("Drive "+strconv.Itoa(i+1), drive.Model+" ("+drive.Type+") - "+strconv.FormatUint(drive.Size/1024/1024/1024, 10)+" GB")
	}

	for i, gpu := range info.GPUs {
		c.AddRow("GPU "+strconv.Itoa(i+1), gpu.Model+" ("+gpu.Vendor+":"+gpu.Device+") - Driver: "+gpu.Driver)
	}
}
