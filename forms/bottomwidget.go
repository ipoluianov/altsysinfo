package forms

import (
	"os"
	"runtime"

	"github.com/ipoluianov/altsysinfo/app"
	"github.com/ipoluianov/altsysinfo/config"
	"github.com/ipoluianov/altsysinfo/install"
	"github.com/ipoluianov/altsysinfo/texts"
	"github.com/ipoluianov/altsysinfo/widgets"
	"github.com/ipoluianov/nui/ui"
)

// BottomWidget is the status bar: the computer the information is about and the links
type BottomWidget struct {
	ui.Widget

	// What the install link does: install, update or uninstall
	installStatus install.Status
}

func NewBottomWidget() *BottomWidget {
	var c BottomWidget
	c.InitWidget()
	c.SetPanelPadding(6)

	c.AddWidget(0, 0, newMutedLabel(machineSummary()))
	c.AddWidget(0, 1, ui.NewHSpacer())
	links := []*ui.Label{
		newLinkLabel(func() string { return texts.T().Settings }, func() { lastCreatedMainWidget.ShowSettings() }),
		newLinkLabel(func() string { return texts.T().Help }, func() { openDocs(&c, "help") }),
	}
	// A downloaded copy offers to install or update itself, the installed one to be removed
	c.installStatus = install.CurrentStatus()
	if c.installStatus != install.StatusNone {
		links = append(links, newLinkLabel(c.installText, c.onInstallLink))
	}
	links = append(links, newLinkLabel(func() string { return texts.T().About }, c.onAbout))
	for i, lbl := range links {
		if i > 0 {
			space := ui.NewSpace()
			space.SetSize(12, 0)
			c.AddWidget(0, 1+i*2, space)
		}
		c.AddWidget(0, 2+i*2, lbl)
	}
	return &c
}

// machineSummary is e.g. "my-host  ·  darwin/arm64"
func machineSummary() string {
	s := runtime.GOOS + "/" + runtime.GOARCH
	if host, err := os.Hostname(); err == nil && host != "" {
		s = host + "  ·  " + s
	}
	return s
}

// newMutedLabel creates a label in the secondary text color of the theme
func newMutedLabel(text string) *ui.Label {
	lbl := ui.NewLabel(text)
	lbl.SetForegroundColor(widgets.ColorMuted.Get())
	mutedLabels = append(mutedLabels, lbl)
	return lbl
}

// newLinkLabel creates a hyperlink-style label with the text from text() that calls onClick on left click.
// It is underlined only under the mouse, so a row of links stays quiet.
func newLinkLabel(text func() string, onClick func()) *ui.Label {
	lbl := ui.NewLabel("")
	lbl.SetTextFunc(text)
	lbl.SetForegroundColor(widgets.ColorLink.Get())
	lbl.SetOnMouseEnter(func() { lbl.SetUnderline(true) })
	lbl.SetOnMouseLeave(func() { lbl.SetUnderline(false) })
	linkLabels = append(linkLabels, lbl)
	lbl.SetMouseCursor(ui.MouseCursorPointer)
	lbl.SetOnMouseDown(func(button ui.MouseButton, x int, y int, mods ui.KeyModifiers) bool {
		if button != ui.MouseButtonLeft {
			return false
		}
		onClick()
		return true
	})
	return lbl
}

// openDocs opens the docs on the site; campaign tells which place in the app the visit came from
func openDocs(parent ui.Widgeter, campaign string) {
	if err := app.OpenSiteURL(app.DocsURL, campaign); err != nil {
		ui.ShowMessageBox(parent, texts.T().Error, err.Error())
	}
}

func (c *BottomWidget) onAbout() {
	c.ShowDialog(NewAboutDialog())
}

// installText is the text of the install link for what it does now
func (c *BottomWidget) installText() string {
	switch c.installStatus {
	case install.StatusUpdate:
		return texts.T().Update
	case install.StatusUninstall:
		return texts.T().Uninstall
	}
	return texts.T().Install
}

func (c *BottomWidget) onInstallLink() {
	if c.installStatus == install.StatusUninstall {
		c.onUninstall()
		return
	}
	c.onInstall()
}

// onInstall copies the application to ~/.altbins and registers it, then
// quits for the installed copy to start (see main). Over an older installed
// version it is the same: that one is replaced.
func (c *BottomWidget) onInstall() {
	t := texts.T()
	ui.ShowQuestionMessageBoxOKCancel(c, c.installText(), t.InstallAsk(install.Dir()), func() {
		if err := install.Install(); err != nil {
			ui.ShowMessageBox(c, t.Error, t.InstallFailed(err.Error()))
			return
		}
		install.RelaunchAfterExit()
		lastCreatedMainWidget.Quit()
	}, nil)
}

// onUninstall removes the installed copy; the settings stay. When that is
// this copy, it quits and its binary is deleted once it has.
func (c *BottomWidget) onUninstall() {
	t := texts.T()
	ui.ShowQuestionMessageBoxOKCancel(c, t.Uninstall, t.UninstallAsk(config.ConfigDirectory()), func() {
		installed := install.IsInstalledCopy()
		if err := install.Uninstall(); err != nil {
			ui.ShowMessageBox(c, t.Error, err.Error())
			return
		}
		if installed {
			lastCreatedMainWidget.Quit()
			return
		}
		c.installStatus = install.CurrentStatus()
		ui.ShowToast(c, t.Uninstalled, ui.ToastSuccess)
	}, nil)
}
