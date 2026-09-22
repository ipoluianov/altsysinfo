package forms

import "github.com/u00io/nuiforms/ui"

type BottomWidget struct {
	ui.Widget
}

func NewBottomWidget() *BottomWidget {
	var c BottomWidget
	c.InitWidget()
	c.AddWidget(0, 0, ui.NewLabel("Bottom"))
	c.AddWidget(0, 1, ui.NewHSpacer())
	c.AddWidget(0, 2, ui.NewLabel("AltBins"))
	return &c
}
