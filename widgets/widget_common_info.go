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

func (c *WidgetCommonInfo) LoadCommonInfo() {
	info, err := system.GetInfo()
	if err != nil {
		return
	}

	c.lvItems.SetRowCount(1 + 1 + len(info.Drives) + len(info.GPUs))

	c.lvItems.SetCellText2(0, 0, "CPU")
	c.lvItems.SetCellText2(0, 1, strconv.FormatInt(int64(info.CpuInfo.Cores), 10)+" x "+info.CpuInfo.ModelStr)

	c.lvItems.SetCellText2(1, 0, "RAM")
	c.lvItems.SetCellText2(1, 1, strconv.FormatUint(info.RamInfo.Total/1024/1024, 10)+" MB")

	for i, drive := range info.Drives {
		c.lvItems.SetCellText2(2+i, 0, "Drive")
		c.lvItems.SetCellText2(2+i, 1, drive.Model+" ("+drive.Type+") - "+strconv.FormatUint(drive.Size/1024/1024/1024, 10)+" GB")
	}

	for i, gpu := range info.GPUs {
		c.lvItems.SetCellText2(2+len(info.Drives)+i, 0, "GPU")
		c.lvItems.SetCellText2(2+len(info.Drives)+i, 1, gpu.Model+" ("+gpu.Vendor+":"+gpu.Device+") - Driver: "+gpu.Driver)
	}
}
