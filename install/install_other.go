//go:build !windows

package install

import "errors"

var errNotSupported = errors.New("installation is supported on Windows only")

// Available tells whether the Install button makes sense: never here
func Available() bool {
	return false
}

func Install() error {
	return errNotSupported
}

func Uninstall() error {
	return errNotSupported
}

func StartInstalled() error {
	return errNotSupported
}

// Confirm and Inform are the uninstaller's dialogs, which has no window of its own
func Confirm(title string, text string) bool {
	return false
}

func Inform(title string, text string) {}
