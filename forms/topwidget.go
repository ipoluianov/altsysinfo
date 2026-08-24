package forms

import "github.com/u00io/nuiforms/ui"

type TopWidget struct {
	ui.Widget
}

func NewTopWidget() *TopWidget {
	var c TopWidget
	c.InitWidget()

	c.AddWidgetOnGrid(ui.NewLabel("TOP"), 0, 0)
	return &c
}
