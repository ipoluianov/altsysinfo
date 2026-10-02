package widgets

import (
	"image/color"

	"github.com/ipoluianov/nui/ui"
)

// ThemeColors is a color of the application in the dark and the light theme
type ThemeColors struct {
	Dark, Light color.RGBA
}

// Get returns the color for the current theme
func (c ThemeColors) Get() color.RGBA {
	if ui.IsDarkTheme {
		return c.Dark
	}
	return c.Light
}

var (
	ColorLink  = ThemeColors{ui.ColorFromHex("#3d8bf2"), ui.ColorFromHex("#1d6fb5")}
	ColorMuted = ThemeColors{ui.ColorFromHex("#888888"), ui.ColorFromHex("#7a7a7a")}
	// The text of the section header rows
	ColorSection = ThemeColors{ui.ColorFromHex("#3a8ee6"), ui.ColorFromHex("#1d6fb5")}
)

const (
	// rowExtraHeight makes the rows a little taller than the theme's, so the tables breathe
	rowExtraHeight = 8
	cellPaddingX   = 8
)

// CellImageSize is the size of an image in a table cell: the row without the padding
func CellImageSize() int {
	return ui.ThemeRowHeight() + rowExtraHeight - 2*cellPaddingX
}

// NewTable creates a table in the style of the application
func NewTable() *ui.Table {
	t := ui.NewTable()
	t.SetSelectingRows(true)
	// The last column takes the rest of the width: no empty strip on the right
	t.SetStretchLastColumn(true)
	t.SetCellPadding(cellPaddingX)
	ApplyTableTheme(t)
	return t
}

// ApplyTableTheme selects the rows with the soft selection color of the theme
// instead of the accent one: the colored texts stay readable on it
func ApplyTableTheme(t *ui.Table) {
	t.SetProp("background_selected_cell", ui.ColorToHex(ui.CurrentPalette().Selection))
	t.SetRowHeight(ui.ThemeRowHeight() + rowExtraHeight)
}
