package widgets

import (
	"strconv"

	"github.com/ipoluianov/altsysinfo/system"
	"github.com/u00io/nuiforms/ui"
)

type WidgetRamInfo struct {
	ui.Widget

	lvItems *ui.Table
}

func NewWidgetRamInfo() *WidgetRamInfo {
	var c WidgetRamInfo
	c.InitWidget()
	c.lvItems = ui.NewTable()
	c.AddWidget(0, 0, c.lvItems)
	c.lvItems.SetSelectingRows(true)

	c.lvItems.SetColumnCount(3)
	c.lvItems.SetColumnWidth(0, 300)
	c.lvItems.SetColumnWidth(1, 200)
	c.lvItems.SetColumnWidth(2, 200)
	c.lvItems.SetColumnName(0, "Model")
	c.lvItems.SetColumnName(1, "Size")
	c.lvItems.SetColumnName(2, "Speed")

	c.LoadRamInfo()
	return &c
}

func (c *WidgetRamInfo) LoadRamInfo() {
	info, err := system.GetInfo()
	if err != nil {
		return
	}

	c.lvItems.SetRowCount(0)
	for i, ram := range info.RamDevices {
		c.lvItems.SetRowCount(i + 1)
		c.lvItems.SetCellText2(i, 0, ram.Model)
		c.lvItems.SetCellText2(i, 1, strconv.FormatUint(ram.Size/1024/1024, 10)+" MB")
		c.lvItems.SetCellText2(i, 2, strconv.FormatUint(ram.Speed, 10)+" MHz")
	}
}
