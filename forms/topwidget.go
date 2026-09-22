package forms

import "github.com/u00io/nuiforms/ui"

type TopWidget struct {
	ui.Widget
}

func NewTopWidget() *TopWidget {
	var c TopWidget
	c.InitWidget()

	c.AddWidget(0, 0, ui.NewLabel("TOP"))
	return &c
}
