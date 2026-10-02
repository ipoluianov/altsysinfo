package forms

import (
	"github.com/ipoluianov/altsysinfo/widgets"
	"github.com/ipoluianov/nui/ui"
)

// The color themes offered in the settings (config.Settings.Theme)
const (
	themeDark  = ""
	themeLight = "light"
)

// linkLabels are recolored when the theme changes, see newLinkLabel
var linkLabels []*ui.Label

// mutedLabels are recolored when the theme changes, see newMutedLabel
var mutedLabels []*ui.Label

// ApplyTheme switches the application to the theme of the settings.
// The tables are styled when created, so the main form rebuilds them.
func ApplyTheme(theme string) {
	if theme == themeLight {
		ui.ApplyLightTheme()
	} else {
		ui.ApplyDarkTheme()
	}
	for _, lbl := range linkLabels {
		lbl.SetForegroundColor(widgets.ColorLink.Get())
	}
	for _, lbl := range mutedLabels {
		lbl.SetForegroundColor(widgets.ColorMuted.Get())
	}
	for _, setIcon := range themedIcons {
		setIcon()
	}
	if lastCreatedMainWidget != nil {
		lastCreatedMainWidget.applyTheme()
	}
}
