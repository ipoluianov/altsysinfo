// Package install puts the application into ~/.altbins and registers it in
// the system, so a downloaded copy can install itself with one click and be
// removed the usual way. Only Windows is supported: Linux has its installer
// script and packages, macOS its .dmg.
package install

import (
	"os"
	"path/filepath"
	"slices"
)

const (
	// appName is the file name of the installed binary
	appName = "altsysinfo"

	// UninstallArg starts the application as its own uninstaller
	UninstallArg = "--uninstall"
	// QuietArg removes it without asking (with UninstallArg)
	QuietArg = "--quiet"
	// InstalledArg tells the installed copy that it has just been installed
	InstalledArg = "--installed"
)

// iconPNG is the application icon, set by SetIcon
var iconPNG []byte

// relaunch is set once the application is installed: the installed copy is
// started when this one quits
var relaunch bool

// SetIcon sets the icon of the shortcuts and of the entry in the list of apps
func SetIcon(png []byte) {
	iconPNG = png
}

// HasArg tells whether the application was started with the argument
func HasArg(arg string) bool {
	return slices.Contains(os.Args[1:], arg)
}

// Dir is where the application is installed, shared by all altbins utilities
func Dir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}
	return filepath.Join(homeDir, ".altbins")
}

// RelaunchAfterExit asks to start the installed copy when this one quits
func RelaunchAfterExit() {
	relaunch = true
}

// RelaunchPending tells whether the installed copy must be started on exit
func RelaunchPending() bool {
	return relaunch
}
