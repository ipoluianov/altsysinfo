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
	form.Panel().AddWidget(0, 0, mainForm)
	form.Show()
	form.Exec()
}
