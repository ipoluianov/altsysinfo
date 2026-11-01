# Changelog

## Unreleased

The first public release.

### Added

- Categories: common information, RAM, system, CPU, storage, graphics and
  displays, network, PCI and USB devices, audio, sensors and battery.
- Any value is copied to the clipboard with one click.
- A PDF report with the CPU, memory, drives and graphics.
- A lookup of PCI vendors and devices in the built-in PCI ID database.
- System tray, always on top, keyboard shortcuts, light and dark themes.
- 31 languages.
- Only one copy runs at a time. Starting AltSysInfo again brings the
  running window to the front, even from the tray.
- Windows: a downloaded `altsysinfo.exe` installs itself into
  `%USERPROFILE%\.altbins` with the **Install** link in the status bar. It
  adds shortcuts to the Start menu and the desktop and an entry to
  "Installed apps". The link reads **Update** when an older version is
  installed, and **Uninstall** when this version or a newer one is
  installed.
- Linux: `.deb` and `.rpm` packages, and `linux-install.sh` for a
  user-level install of the x86-64 or ARM64 build, whichever fits the
  machine.
- macOS: a `.dmg` with AltSysInfo.app, signed and notarized, so it opens
  without the Gatekeeper warning.
