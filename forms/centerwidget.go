package forms

import (
	"github.com/ipoluianov/alsysinfo/system"
	"github.com/u00io/nuiforms/ui"
)

type CenterWidget struct {
	ui.Widget

	lvItems *ui.Table
}

func NewCenterWidget() *CenterWidget {
	var c CenterWidget
	c.InitWidget()
	c.SetXExpandable(true)
	c.SetYExpandable(true)

	c.lvItems = ui.NewTable()
	c.AddWidgetOnGrid(c.lvItems, 0, 0)

	c.LoadCommonInfo()
	c.lvItems.SetColumnCount(2)
	c.lvItems.SetColumnWidth(0, 200)
	c.lvItems.SetColumnWidth(1, 600)
	c.lvItems.SetColumnName(0, "Name")
	c.lvItems.SetColumnName(1, "Value")

	c.AddTimer(500, c.updateTimer)

	return &c
}

func (c *CenterWidget) LoadCommonInfo() {
	items, err := system.GetCommonInfo()
	if err != nil {
		return
	}
	c.lvItems.SetRowCount(len(items))
	for i, item := range items {
		c.lvItems.SetCellText2(i, 0, item.Name)
		c.lvItems.SetCellText2(i, 1, item.Value)
	}
}

func (c *CenterWidget) updateTimer() {
	c.LoadCommonInfo()
}
