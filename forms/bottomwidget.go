package forms

import (
	"os"
	"runtime"

	"github.com/ipoluianov/altsysinfo/app"
	"github.com/ipoluianov/altsysinfo/install"
	"github.com/ipoluianov/altsysinfo/texts"
	"github.com/ipoluianov/altsysinfo/widgets"
	"github.com/ipoluianov/nui/ui"
)

// BottomWidget is the status bar: the computer the information is about and the links
type BottomWidget struct {
	ui.Widget
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
	// A downloaded copy offers to install itself
	if install.Available() {
		links = append(links, newLinkLabel(func() string { return texts.T().Install }, c.onInstall))
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

// onInstall copies the application to ~/.altbins and registers it, then
// quits for the installed copy to start (see main)
func (c *BottomWidget) onInstall() {
	t := texts.T()
	ui.ShowQuestionMessageBoxOKCancel(c, t.Install, t.InstallAsk(install.Dir()), func() {
		if err := install.Install(); err != nil {
			ui.ShowMessageBox(c, t.Error, t.InstallFailed(err.Error()))
			return
		}
		install.RelaunchAfterExit()
		lastCreatedMainWidget.Quit()
	}, nil)
}
