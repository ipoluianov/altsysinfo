// Package texts holds the texts of the application in all its languages
package texts

import (
	"github.com/ipoluianov/nui/ui"
	"github.com/ipoluianov/nui/ui/i18n"
)

// Strings are all the texts of the application. A misspelled field is a
// compile error; a field a language leaves empty falls back to English
// (texts_test.go checks that none is left).
type Strings struct {
	Error string

	// Main window
	ToolRefresh   string
	ToolPdfReport string
	ToolTray      string
	TrayShow      string
	TrayQuit      string
	Settings      string
	Help          string
	About         string

	// List of the categories
	Categories string
	Category   CategoryStrings

	// Tables
	ColName       string
	ColValue      string
	ColModel      string
	ColSize       string
	ColSpeed      string
	ColVendorID   string
	ColDeviceID   string
	ColVendorName string
	ColDeviceName string
	// The filter of the PCI vendors and devices
	VendorDeviceName string
	// Rows of the common information
	CPUModel string
	CPUCores string
	Drive    func(n int) string
	GPU      func(n int) string
	Driver   string
	NoData   string
	Loading  string
	// A toast after a value is copied to the clipboard
	Copied     string
	CannotLoad func(err string) string

	// PDF report
	SaveReportTitle      string
	PDFFiles             string
	CannotCollect        func(err string) string
	CannotGenerate       func(err string) string
	CannotOpenSaveDialog func(err string) string
	CannotSave           func(err string) string
	SavedTo              func(path string) string

	// Settings dialog
	AlwaysOnTop    string
	Language       string
	LanguageSystem string
	Theme          string
	ThemeDark      string
	ThemeLight     string

	// About dialog
	AboutTitle   func(name string) string
	Version      string
	Author       string
	License      string
	VisitWebsite string
	Close        string

	// Installing a downloaded copy to ~/.altbins (Windows) and removing it
	Install          string
	InstallAsk       func(dir string) string
	InstallFailed    func(err string) string
	Installed        string // a toast in the installed copy
	UninstallAsk     func(dir string) string
	UninstallRunning string
	Uninstalled      string
}

// CategoryStrings are the names of the categories in the list
type CategoryStrings struct {
	Common          string
	RAM             string
	System          string
	CPU             string
	Storage         string
	Graphics        string
	Network         string
	PCI             string
	USB             string
	Audio           string
	Sensors         string
	Battery         string
	PCIVendorDevice string
}

// Name returns the name of the category by its ID; an unknown one is returned as is
func (c *CategoryStrings) Name(id string) string {
	names := map[string]string{
		"common":   c.Common,
		"ram":      c.RAM,
		"system":   c.System,
		"cpu":      c.CPU,
		"storage":  c.Storage,
		"graphics": c.Graphics,
		"network":  c.Network,
		"pci":      c.PCI,
		"usb":      c.USB,
		"audio":    c.Audio,
		"sensors":  c.Sensors,
		"battery":  c.Battery,
		"pcidev":   c.PCIVendorDevice,
	}
	if name, ok := names[id]; ok {
		return name
	}
	return id
}

// Column translates a column name of the details tables, which come in English
func (s *Strings) Column(name string) string {
	switch name {
	case "Name":
		return s.ColName
	case "Value":
		return s.ColValue
	case "Model":
		return s.ColModel
	case "Size":
		return s.ColSize
	case "Speed":
		return s.ColSpeed
	}
	return name
}

// Languages are offered in the settings, by tag and own name
var Languages = []struct{ Tag, Name string }{
	{"en", "English"},
	{"ru", "Русский"},
	{"pl", "Polski"},
	{"sr", "Српски"},
	{"de", "Deutsch"},
	{"fr", "Français"},
	{"es", "Español"},
	{"it", "Italiano"},
	{"pt", "Português"},
	{"zh", "中文"},
	{"ja", "日本語"},
	{"ko", "한국어"},
}

var catalog = i18n.NewCatalog(en, map[string]Strings{
	"ru": ru, "pl": pl, "sr": sr, "de": de, "fr": fr, "es": es, "it": it, "pt": pt, "zh": zh, "ja": ja, "ko": ko,
})

// T returns the texts in the language of the application
func T() *Strings {
	return catalog.Get(ui.Language())
}

// The library translates its own buttons (OK, Cancel...) to Russian and Chinese only
func init() {
	for lang, s := range map[string]ui.UIStrings{
		"pl": {OK: "OK", Cancel: "Anuluj", Yes: "Tak", No: "Nie"},
		"sr": {OK: "У реду", Cancel: "Откажи", Yes: "Да", No: "Не"},
		"de": {OK: "OK", Cancel: "Abbrechen", Yes: "Ja", No: "Nein"},
		"fr": {OK: "OK", Cancel: "Annuler", Yes: "Oui", No: "Non"},
		"es": {OK: "Aceptar", Cancel: "Cancelar", Yes: "Sí", No: "No"},
		"it": {OK: "OK", Cancel: "Annulla", Yes: "Sì", No: "No"},
		"pt": {OK: "OK", Cancel: "Cancelar", Yes: "Sim", No: "Não"},
		"ja": {OK: "OK", Cancel: "キャンセル", Yes: "はい", No: "いいえ"},
		"ko": {OK: "확인", Cancel: "취소", Yes: "예", No: "아니요"},
	} {
		ui.RegisterUIStrings(lang, s)
	}
}

// SetLanguage switches the application to the language of the settings, "" - the system's
func SetLanguage(lang string) {
	if lang == "" {
		lang = ui.SystemLanguage()
	}
	ui.SetLanguage(lang)
}
