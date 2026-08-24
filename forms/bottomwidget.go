package forms

import "github.com/u00io/nuiforms/ui"

type BottomWidget struct {
	ui.Widget
}

func NewBottomWidget() *BottomWidget {
	var c BottomWidget
	c.InitWidget()
	c.AddWidgetOnGrid(ui.NewLabel("Bottom"), 0, 0)
	c.AddWidgetOnGrid(ui.NewHSpacer(), 0, 1)
	c.AddWidgetOnGrid(ui.NewLabel("AltBins"), 0, 2)
	return &c
}
