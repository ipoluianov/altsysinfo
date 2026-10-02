package forms

import "github.com/ipoluianov/nui/ui"

// The application hotkeys; "Mod" is Cmd on macOS and Ctrl elsewhere
const (
	shortcutRefresh   = "F5"
	shortcutPdfReport = "Mod+S"
	shortcutTray      = "Mod+T"
	shortcutHelp      = "F1"
)

// AddShortcuts sets the application hotkeys on the main window. They work
// whatever widget has the focus, the filter text boxes keep their own keys.
func AddShortcuts(form *ui.Form) {
	form.AddShortcut(shortcutRefresh, func() { lastCreatedTopWidget.onBtnRefresh() })
	form.AddShortcut(shortcutPdfReport, func() { lastCreatedTopWidget.savePdfReport() })
	form.AddShortcut(shortcutTray, func() { lastCreatedTopWidget.onBtnTray() })
	form.AddShortcut(shortcutHelp, func() { openDocs(lastCreatedMainWidget, "help_f1") })
	form.AddShortcut("Ctrl+M", func() {
		if form.IsMaximized() {
			form.Restore()
		} else {
			form.Maximize()
		}
	})
	form.AddShortcut("Ctrl+H", form.Minimize)
	form.AddShortcut("Alt+X", lastCreatedMainWidget.Quit)
}

// withShortcut adds the hotkey to a tooltip, e.g. "Refresh (F5)"
func withShortcut(text string, shortcut string) string {
	return text + " (" + ui.MustParseShortcut(shortcut).String() + ")"
}
