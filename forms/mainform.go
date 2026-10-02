package forms

import (
	"github.com/ipoluianov/altsysinfo/app"
	"github.com/ipoluianov/altsysinfo/config"
	"github.com/ipoluianov/altsysinfo/texts"
	"github.com/ipoluianov/nui/ui"
)

type MainForm struct {
	ui.Widget

	splitter *ui.Splitter

	topWidget    *TopWidget
	leftWidget   *LeftWidget
	centerWidget *CenterWidget
	bottomWidget *BottomWidget

	// The width the user gave the categories
	categoriesWidth int
}

const (
	defaultWindowWidth     = 1100
	defaultWindowHeight    = 800
	categoriesInitialWidth = 240
	defaultMode            = "common"
)

var lastCreatedMainWidget *MainForm

func NewMainForm() *MainForm {
	var c MainForm
	c.InitWidget()
	lastCreatedMainWidget = &c
	c.topWidget = NewTopWidget()
	c.leftWidget = NewLeftWidget(c.SetMode)
	c.centerWidget = NewCenterWidget()
	c.bottomWidget = NewBottomWidget()

	c.AddWidget(0, 0, c.topWidget)
	// The categories keep their width when the window is resized
	c.splitter = ui.NewHSplitter()
	c.splitter.SetWidgets(c.leftWidget, c.centerWidget)
	c.categoriesWidth = categoriesInitialWidth
	c.splitter.SetFirstSize(c.categoriesWidth)
	c.splitter.SetOnSplitChanged(func() {
		c.categoriesWidth = c.splitter.FirstSize()
	})
	c.AddWidget(1, 0, c.splitter)
	c.AddWidget(2, 0, c.bottomWidget)

	c.SetPanelPadding(3)
	cache.prefetch()
	return &c
}

func (c *MainForm) SetMode(mode string) {
	c.centerWidget.SetMode(mode)
}

// Refresh loads the information of the selected category again
func (c *MainForm) Refresh() {
	c.centerWidget.Refresh()
}

// ApplySettings saves the settings and applies what they change in the window
func (c *MainForm) ApplySettings(s config.Settings) {
	languageChanged := s.Language != config.GetSettings().Language
	themeChanged := s.Theme != config.GetSettings().Theme
	if err := config.SetSettings(s); err != nil {
		ui.ShowMessageBox(c, texts.T().Error, err.Error())
	}
	if languageChanged {
		texts.SetLanguage(s.Language)
	}
	if themeChanged {
		ApplyTheme(s.Theme)
	}
	c.Form().SetAlwaysOnTop(s.AlwaysOnTop)
}

// ShowInstalled tells that this copy has just been installed and started in place of the downloaded one
func (c *MainForm) ShowInstalled() {
	ui.ShowToast(c, texts.T().Installed, ui.ToastSuccess)
}

// ApplyLanguage updates the texts that do not follow the language by themselves
func (c *MainForm) ApplyLanguage() {
	c.UpdateTitle()
	c.updateTrayMenu()
	c.leftWidget.applyLanguage()
	c.centerWidget.Rebuild()
}

// Quit closes the application from code
func (c *MainForm) Quit() {
	// Closing from code does not call Form.OnClose, so the layout is saved here
	c.SaveWindowState()
	c.Form().Close()
}

// ShowSettings opens the settings dialog
func (c *MainForm) ShowSettings() {
	c.ShowDialog(NewSettingsDialog(config.GetSettings(), c.ApplySettings))
}

// applyTheme restyles the tables: they take the theme colors when created
func (c *MainForm) applyTheme() {
	c.leftWidget.applyTheme()
	c.centerWidget.Rebuild()
	if f := c.Form(); f != nil {
		f.Update()
	}
}

func (c *MainForm) UpdateTitle() {
	c.Form().SetTitle(app.DisplayName)
	updateTrayTooltip(app.DisplayName)
}

// RestoreWindowState applies the saved window layout before the form is shown.
// Returns whether the window should be maximized once shown.
func (c *MainForm) RestoreWindowState(form *ui.Form) (maximized bool) {
	form.SetSize(defaultWindowWidth, defaultWindowHeight)
	state, ok := config.LoadWindowState()
	if !ok || !c.leftWidget.SelectMode(state.Category) {
		c.leftWidget.SelectMode(defaultMode)
	}
	if !ok {
		return false
	}
	form.SetSize(state.Width, state.Height)
	// 0,0 means the position was not known (the window manager did not report it)
	if state.X != 0 || state.Y != 0 {
		form.Move(state.X, state.Y)
	}
	if state.CategoriesWidth > 0 {
		c.categoriesWidth = state.CategoriesWidth
		c.splitter.SetFirstSize(c.categoriesWidth)
	}
	return state.Maximized
}

// SaveWindowState remembers the window layout for the next start
func (c *MainForm) SaveWindowState() {
	form := c.Form()
	state, _ := config.LoadWindowState()
	state.Maximized = form.IsMaximized()
	// The size of a maximized window is the screen size: keep the normal one
	if !state.Maximized {
		state.X, state.Y = form.Position()
		state.Width, state.Height = form.Size()
	}
	state.CategoriesWidth = c.categoriesWidth
	state.Category = c.centerWidget.Mode()
	config.SaveWindowState(state)
}

func (c *MainForm) Activate() {
	c.leftWidget.FocusTable()
}
