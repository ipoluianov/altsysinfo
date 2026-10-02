package widgets

import (
	"strconv"

	"github.com/ipoluianov/altsysinfo/system"
	"github.com/ipoluianov/altsysinfo/texts"
	"github.com/ipoluianov/nui/ui"
)

type WidgetCommonInfo struct {
	ui.Widget

	lvItems *ui.Table
}

// NewWidgetCommonInfo shows the main facts of the system.GetInfo result
func NewWidgetCommonInfo(info system.Info) *WidgetCommonInfo {
	var c WidgetCommonInfo
	c.InitWidget()
	c.lvItems = NewTable()
	c.AddWidget(0, 0, c.lvItems)
	c.lvItems.SetColumnCount(2)
	c.lvItems.SetColumnWidth(0, 200)
	c.lvItems.SetColumnWidth(1, 600)
	c.lvItems.SetColumnName(0, texts.T().ColName)
	c.lvItems.SetColumnName(1, texts.T().ColValue)

	c.LoadCommonInfo(info)
	return &c
}

func (c *WidgetCommonInfo) AddRow(name, value string) {
	row := c.lvItems.RowCount()
	c.lvItems.SetRowCount(row + 1)
	c.lvItems.SetCellText2(row, 0, name)
	c.lvItems.SetCellText2(row, 1, value)
}

func (c *WidgetCommonInfo) LoadCommonInfo(info system.Info) {
	t := texts.T()

	c.lvItems.SetRowCount(0)
	c.AddRow(t.CPUModel, info.CpuInfo.ModelStr)
	c.AddRow(t.CPUCores, strconv.FormatInt(int64(info.CpuInfo.Cores), 10))

	c.AddRow(t.Category.RAM, strconv.FormatUint(info.RamInfo.Total/1024/1024, 10)+" MB")

	for i, drive := range info.Drives {
		c.AddRow(t.Drive(i+1), drive.Model+" ("+drive.Type+") - "+strconv.FormatUint(drive.Size/1024/1024/1024, 10)+" GB")
	}

	for i, gpu := range info.GPUs {
		c.AddRow(t.GPU(i+1), gpu.Model+" ("+gpu.Vendor+":"+gpu.Device+") - "+t.Driver+": "+gpu.Driver)
	}
}
