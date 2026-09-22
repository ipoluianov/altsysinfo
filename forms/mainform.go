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
	c.leftWidget = NewLeftWidget(c.SetMode)
	c.centerWidget = NewCenterWidget()
	c.bottomWidget = NewBottomWidget()

	c.AddWidget(0, 0, c.topWidget)
	c.panelCenter = ui.NewPanel()
	c.panelCenter.AddWidget(0, 0, c.leftWidget)
	c.panelCenter.AddWidget(0, 1, c.centerWidget)
	c.AddWidget(1, 0, c.panelCenter)
	c.AddWidget(2, 0, c.bottomWidget)
	return &c
}

func (c *MainForm) SetMode(mode string) {
	c.centerWidget.SetMode(mode)
}
