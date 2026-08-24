package main

import (
	"github.com/ipoluianov/altsysinfo/forms"
	"github.com/u00io/nuiforms/ui"
)

func main() {
	form := ui.NewForm()
	form.SetTitle("SysInfo")
	form.SetSize(1100, 800)
	mainForm := forms.NewMainForm()
	form.Panel().AddWidgetOnGrid(mainForm, 0, 0)
	form.Exec()
}
