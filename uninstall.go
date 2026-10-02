package main

import (
	"os"

	"github.com/ipoluianov/altsysinfo/app"
	"github.com/ipoluianov/altsysinfo/config"
	"github.com/ipoluianov/altsysinfo/install"
	"github.com/ipoluianov/altsysinfo/instance"
	"github.com/ipoluianov/altsysinfo/texts"
)

// uninstall removes the installed application, leaving its settings. It has
// no window: the questions are system message boxes. quiet removes it without
// asking or reporting.
func uninstall(quiet bool) {
	// Only the language is needed; nothing is written
	config.LoadSettings()
	texts.SetLanguage(config.GetSettings().Language)
	t := texts.T()

	// The running copy shows its window, so it is clear what to close
	inst, ok := instance.Acquire(config.ConfigDirectory())
	if !ok {
		if !quiet {
			install.Inform(app.DisplayName, t.UninstallRunning)
		}
		os.Exit(1)
	}
	defer inst.Close()

	if !quiet && !install.Confirm(app.DisplayName, t.UninstallAsk(config.ConfigDirectory())) {
		return
	}
	if err := install.Uninstall(); err != nil {
		if !quiet {
			install.Inform(app.DisplayName, t.Error+": "+err.Error())
		}
		inst.Close()
		os.Exit(1)
	}
	if !quiet {
		install.Inform(app.DisplayName, t.Uninstalled)
	}
}
