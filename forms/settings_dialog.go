package forms

import (
	"github.com/ipoluianov/altsysinfo/config"
	"github.com/ipoluianov/altsysinfo/texts"
	"github.com/ipoluianov/nui/ui"
)

// SettingsDialog edits the application options (config.Settings)
type SettingsDialog struct {
	ui.DialogContent

	chkOnTop *ui.Checkbox

	// The first item is the system language, then the languages list
	cmbLanguage *ui.ComboBox
	cmbTheme    *ui.ComboBox

	btnOK     *ui.Button
	btnCancel *ui.Button

	onAccept func(config.Settings)
}

func NewSettingsDialog(settings config.Settings, onAccept func(config.Settings)) *SettingsDialog {
	var c SettingsDialog
	c.InitWidget()
	c.onAccept = onAccept
	t := texts.T()

	checks := ui.NewPanel()
	c.AddWidget(0, 0, checks)
	c.AddWidget(1, 0, ui.NewVSpacer())
	buttons := ui.NewPanel()
	c.AddWidget(2, 0, buttons)

	c.chkOnTop = ui.NewCheckbox(t.AlwaysOnTop)
	c.chkOnTop.SetChecked(settings.AlwaysOnTop)
	c.chkOnTop.SetXExpandable(true)
	checks.AddWidget(0, 0, c.chkOnTop)

	// Label and combo box rows, the combo boxes aligned in a column
	choices := ui.NewPanel()
	choices.SetPanelPadding(0)
	checks.AddWidget(1, 0, choices)
	choices.AddWidget(0, 0, ui.NewLabel(t.Language))
	c.cmbLanguage = ui.NewComboBox()
	c.cmbLanguage.AddItem(t.LanguageSystem, "")
	c.cmbLanguage.SetSelectedIndex(0)
	for i, l := range texts.Languages {
		c.cmbLanguage.AddItem(l.Name, l.Tag)
		if l.Tag == settings.Language {
			c.cmbLanguage.SetSelectedIndex(i + 1)
		}
	}
	choices.AddWidget(0, 1, c.cmbLanguage)
	choices.AddWidget(0, 2, ui.NewHSpacer())

	choices.AddWidget(1, 0, ui.NewLabel(t.Theme))
	c.cmbTheme = ui.NewComboBox()
	c.cmbTheme.AddItem(t.ThemeDark, themeDark)
	c.cmbTheme.AddItem(t.ThemeLight, themeLight)
	c.cmbTheme.SetSelectedIndex(0)
	if settings.Theme == themeLight {
		c.cmbTheme.SetSelectedIndex(1)
	}
	choices.AddWidget(1, 1, c.cmbTheme)

	c.btnOK = ui.NewButton(ui.UIText().OK)
	c.btnOK.SetOnClick(c.Accept)
	c.btnCancel = ui.NewButton(ui.UIText().Cancel)
	c.btnCancel.SetOnClick(c.Cancel)
	buttons.AddWidget(0, 0, ui.NewHSpacer())
	buttons.AddWidget(0, 1, c.btnOK)
	buttons.AddWidget(0, 2, c.btnCancel)

	c.OnDialogShow = func() {
		c.Form().SetTitle(t.Settings)
		c.Form().SetSize(420, 240)
		c.Form().MoveToCenterOfParent()
		c.Form().SetAcceptButton(c.btnOK)
		c.Form().SetCancelButton(c.btnCancel)
	}
	return &c
}

func (c *SettingsDialog) Accept() {
	s := config.GetSettings()
	s.AlwaysOnTop = c.chkOnTop.Checked()
	s.Language, _ = c.cmbLanguage.SelectedItemData().(string)
	s.Theme, _ = c.cmbTheme.SelectedItemData().(string)
	// The callback changes the window below: it runs on its goroutine
	if c.onAccept != nil {
		c.RunInParent(func() { c.onAccept(s) })
	}
	c.Form().Close()
}

func (c *SettingsDialog) Cancel() {
	c.Form().Close()
}
