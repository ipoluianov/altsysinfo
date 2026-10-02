package forms

import (
	"github.com/ipoluianov/altsysinfo/texts"
	"github.com/ipoluianov/nui/ui"
)

// trayIcon exists only while the window is in the tray
var trayIcon *ui.TrayIcon

// HideToTray hides the window, leaving the icon in the system tray to bring
// it back. Where there is no tray the window is just minimized.
func (c *MainForm) HideToTray() {
	if trayIcon == nil {
		icon, err := ui.NewTrayIcon()
		if err != nil {
			c.Form().Minimize()
			return
		}
		trayIcon = icon
		trayIcon.SetOnClick(c.ShowFromTray)
		c.updateTrayMenu()
		c.UpdateTitle()
	}
	c.Form().Hide()
}

// ShowFromTray brings the hidden window back and removes the tray icon. The
// icon goes after the current handler, which may be the icon's own click.
func (c *MainForm) ShowFromTray() {
	c.Form().Show()
	c.Activate()
	c.Form().Invoke(CloseTray)
}

// BringToFront shows the window on top of the others, wherever it is:
// in the tray, minimized or under other windows. Called when the application is started again.
func (c *MainForm) BringToFront() {
	if trayIcon != nil {
		c.ShowFromTray()
		return
	}
	raiseWindow(c.Form())
	c.Activate()
}

// quitFromTray closes the application from the tray menu
func (c *MainForm) quitFromTray() {
	c.Quit()
}

// updateTrayMenu sets the menu again, in the current language
func (c *MainForm) updateTrayMenu() {
	if trayIcon == nil {
		return
	}
	trayIcon.SetMenu(
		ui.TrayMenuItem{Text: texts.T().TrayShow, OnClick: c.ShowFromTray},
		ui.TrayMenuItem{Separator: true},
		ui.TrayMenuItem{Text: texts.T().TrayQuit, OnClick: c.quitFromTray},
	)
}

// updateTrayTooltip shows the window title on the icon
func updateTrayTooltip(title string) {
	if trayIcon != nil {
		trayIcon.SetTooltip(title)
	}
}

// CloseTray removes the tray icon; call it when the application quits
func CloseTray() {
	if trayIcon != nil {
		trayIcon.Close()
		trayIcon = nil
	}
}
