package forms

import "github.com/u00io/nuiforms/ui"

type MainForm struct {
	ui.Widget

	panelCenter *ui.Panel

	topWidget    *TopWidget
	leftWidget   *LeftWidget
	centerWidget *CenterWidget
	bottomWidget *BottomWidget
}

func NewMainForm() *MainForm {
	var c MainForm
	c.InitWidget()
	c.topWidget = NewTopWidget()
	c.leftWidget = NewLeftWidget()
	c.centerWidget = NewCenterWidget()
	c.bottomWidget = NewBottomWidget()

	c.AddWidgetOnGrid(c.topWidget, 0, 0)
	c.panelCenter = ui.NewPanel()
	c.panelCenter.AddWidgetOnGrid(c.leftWidget, 0, 0)
	c.panelCenter.AddWidgetOnGrid(c.centerWidget, 0, 1)
	c.AddWidgetOnGrid(c.panelCenter, 1, 0)
	c.AddWidgetOnGrid(c.bottomWidget, 2, 0)
	return &c
}
