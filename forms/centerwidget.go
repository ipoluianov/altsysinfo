package forms

import (
	"github.com/ipoluianov/altsysinfo/widgets"
	"github.com/u00io/nuiforms/ui"
)

type CenterWidget struct {
	ui.Widget
}

func NewCenterWidget() *CenterWidget {
	var c CenterWidget
	c.InitWidget()
	c.SetXExpandable(true)
	c.SetYExpandable(true)

	c.SetMode("common")

	return &c
}

func (c *CenterWidget) SetWidget(w ui.Widgeter) {
	c.RemoveAllWidgets()
	c.AddWidget(0, 0, w)
}

func (c *CenterWidget) SetMode(mode string) {
	switch mode {
	case "common":
		w := widgets.NewWidgetCommonInfo()
		c.SetWidget(w)
	case "ram":
		w := widgets.NewWidgetRamInfo()
		c.SetWidget(w)
	case "pcidev":
		w := widgets.NewWidgetPciDevInfo()
		c.SetWidget(w)
	}
}
