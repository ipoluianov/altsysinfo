package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

// ConfigDirectory is where the settings and the window layout are kept
func ConfigDirectory() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}
	return filepath.Join(homeDir, ".altbins", ".altsysinfo")
}

// Settings are the application options
type Settings struct {
	AlwaysOnTop bool

	// Language of the interface as a tag like "ru"; "" - the system's
	Language string `json:",omitempty"`

	// Color theme: "light"; "" - the dark one
	Theme string `json:",omitempty"`
}

var (
	settingsMtx sync.Mutex
	settings    Settings
)

func settingsPath() string {
	return filepath.Join(ConfigDirectory(), "settings.json")
}

// LoadSettings reads the settings file; missing values get the defaults
func LoadSettings() {
	var s Settings
	if bs, err := os.ReadFile(settingsPath()); err == nil {
		_ = json.Unmarshal(bs, &s)
	}
	settingsMtx.Lock()
	settings = s
	settingsMtx.Unlock()
}

// GetSettings returns the current settings; safe to call from any goroutine
func GetSettings() Settings {
	settingsMtx.Lock()
	defer settingsMtx.Unlock()
	return settings
}

// SetSettings applies and saves the settings
func SetSettings(s Settings) error {
	settingsMtx.Lock()
	settings = s
	settingsMtx.Unlock()
	return writeJSON(settingsPath(), s)
}

// WindowState is the main window layout, restored on the next start
type WindowState struct {
	X, Y          int
	Width, Height int
	Maximized     bool

	CategoriesWidth int
	Category        string
}

func windowStatePath() string {
	return filepath.Join(ConfigDirectory(), "window.json")
}

// LoadWindowState returns the saved window layout; ok is false when there is none
func LoadWindowState() (state WindowState, ok bool) {
	bs, err := os.ReadFile(windowStatePath())
	if err != nil {
		return state, false
	}
	if json.Unmarshal(bs, &state) != nil || state.Width <= 0 || state.Height <= 0 {
		return WindowState{}, false
	}
	return state, true
}

func SaveWindowState(state WindowState) error {
	return writeJSON(windowStatePath(), state)
}

func writeJSON(path string, v any) error {
	bs, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, bs, 0644)
}
