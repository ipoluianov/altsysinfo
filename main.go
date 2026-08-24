package main

import (
	"github.com/ipoluianov/alsysinfo/forms"
	"github.com/u00io/nuiforms/ui"
)

func main() {
	form := ui.NewForm()
	mainForm := forms.NewMainForm()
	form.Panel().AddWidgetOnGrid(mainForm, 0, 0)
	form.Exec()
}
