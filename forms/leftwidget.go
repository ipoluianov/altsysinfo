package forms

import "github.com/u00io/nuiforms/ui"

type LeftWidget struct {
	ui.Widget

	lvItems *ui.Table
}

func NewLeftWidget() *LeftWidget {
	var c LeftWidget
	c.InitWidget()
	c.lvItems = ui.NewTable()
	c.AddWidgetOnGrid(c.lvItems, 0, 0)
	c.SetMaxWidth(300)

	c.lvItems.SetRowCount(10)
	c.lvItems.SetColumnCount(1)
	c.lvItems.SetColumnWidth(0, 290)
	c.lvItems.SetAllowScroll(false, true)
	c.lvItems.SetColumnName(0, "Categories")

	c.lvItems.SetCellText2(0, 0, "Processor")
	c.lvItems.SetCellText2(1, 0, "Memory")
	c.lvItems.SetCellText2(2, 0, "Storage")
	c.lvItems.SetCellText2(3, 0, "Network")
	c.lvItems.SetCellText2(4, 0, "Display")
	c.lvItems.SetCellText2(5, 0, "Audio")
	c.lvItems.SetCellText2(6, 0, "Input Devices")
	c.lvItems.SetCellText2(7, 0, "Output Devices")
	c.lvItems.SetCellText2(8, 0, "Power Management")
	c.lvItems.SetCellText2(9, 0, "Security")
	return &c
}
