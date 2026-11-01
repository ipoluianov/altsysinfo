# Contributing

Bug reports and ideas are welcome in
[issues](https://github.com/ipoluianov/altsysinfo/issues). Please include
the AltSysInfo version (shown in About) and your operating system.

## Building

You need Go, see the version in [go.mod](go.mod). No cgo is required, so
any platform builds for any other.

```sh
go build .
go test ./...
```

The UI library is [nui](https://github.com/ipoluianov/nui). To work on both
at once, uncomment the `replace` line at the end of `go.mod`.

## Code layout

- `forms/` - the windows, dialogs and the status bar
- `widgets/` - the views of the categories
- `texts/` - all the texts, one file per group of languages
- `system/` - collecting the information, per operating system
- `report/` - the PDF report
- `config/` - settings
- `install/` - installing a downloaded copy on Windows
- `instance/` - keeps a single running copy
- `build/` - release builds and packages, `scripts/` - the Linux installer

## Pull requests

- Keep a change to one thing, and describe what it does and why.
- Every text goes into all the languages in `texts/`.
  `TestAllTranslated` fails on a missing one. A machine translation is
  fine, a native speaker can fix it later.
- Say on which systems you tried it.
