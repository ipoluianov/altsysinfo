package main

import (
	"github.com/ipoluianov/altsysinfo/app"
	"github.com/ipoluianov/altsysinfo/config"
	"github.com/ipoluianov/altsysinfo/forms"
	"github.com/ipoluianov/altsysinfo/install"
	"github.com/ipoluianov/altsysinfo/instance"
	"github.com/ipoluianov/altsysinfo/texts"
	"github.com/ipoluianov/nui/ui"
)

func main() {
	install.SetIcon(iconPNG)
	// Started from "Installed apps" to remove it
	if install.HasArg(install.UninstallArg) {
		uninstall(install.HasArg(install.QuietArg))
		return
	}

	// Taken before anything is read or written: the files belong to one copy
	inst, ok := instance.Acquire(config.ConfigDirectory())
	if !ok {
		return // the running copy shows its window instead
	}
	defer inst.Close()

	config.LoadSettings()
	texts.SetLanguage(config.GetSettings().Language)
	forms.ApplyTheme(config.GetSettings().Theme)
	ui.SetAppIcon(appIcon())
	form := ui.NewForm()
	mainForm := forms.NewMainForm()
	form.Panel().AddWidget(0, 0, mainForm)
	mainForm.UpdateTitle()
	maximized := mainForm.RestoreWindowState(form)
	form.OnClose = func() bool {
		mainForm.SaveWindowState()
		return true
	}
	forms.AddShortcuts(form)
	form.SetOnLanguageChanged(mainForm.ApplyLanguage)
	form.Show()
	inst.Serve(func() { form.Invoke(mainForm.BringToFront) })
	// The form is handled by its own goroutine once shown
	form.Invoke(func() {
		if maximized {
			form.Maximize()
		}
		form.SetAlwaysOnTop(config.GetSettings().AlwaysOnTop)
		mainForm.Activate()
		if install.HasArg(install.InstalledArg) {
			mainForm.ShowInstalled()
		}
	})
	form.Exec()
	forms.CloseTray()

	// Installed from this copy: the installed one takes over, so the lock goes first
	if install.RelaunchPending() {
		inst.Close()
		if err := install.StartInstalled(); err != nil {
			install.Inform(app.DisplayName, err.Error())
		}
	}
}
