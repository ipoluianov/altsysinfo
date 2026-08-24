package widgets

import "github.com/u00io/nuiforms/ui"

type WidgetRamInfo struct {
	ui.Widget
}

func NewWidgetRamInfo() *WidgetRamInfo {
	var c WidgetRamInfo
	c.InitWidget()
	c.AddWidgetOnGrid(ui.NewLabel("RAM"), 0, 0)
	return &c
}
