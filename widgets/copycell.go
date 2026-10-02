package widgets

import (
	"image/color"

	"github.com/ipoluianov/altsysinfo/texts"
	"github.com/ipoluianov/nui/ui"
)

// The copy button drawn after a value: two overlapping squares
const (
	copyIconSquare = 9 // the side of each square
	copyIconShift  = 3 // how far the back square is shifted up and right
	copyIconSize   = copyIconSquare + copyIconShift
	copyIconGap    = 8 // between the text and the button
	// The button takes the clicks a little around itself, it is small
	copyIconHitMargin = 4
)

// AddCopyButtons draws a copy button after the values of the column; a click
// on it puts the value into the clipboard. Rows with one cell (the section
// headers) and empty values get no button.
func AddCopyButtons(tv *ui.Table, col int) {
	for row := 0; row < tv.RowCount(); row++ {
		if tv.GetCellText2(row, col) == "" {
			continue
		}
		tv.SetCellOnDraw(row, col, func(cnv *ui.Canvas) {
			drawCopyCell(cnv, tv, row, col)
		})
	}
	tv.SetOnCellMouseDown(func(button ui.MouseButton, row int, c int, x int, y int, mods ui.KeyModifiers) {
		if button != ui.MouseButtonLeft || c != col || row < 0 {
			return
		}
		text := tv.GetCellText2(row, col)
		if text == "" {
			return
		}
		iconX, _ := copyCellLayout(tv, col, text)
		inCellX := x - columnOffset(tv, col)
		if inCellX < iconX-copyIconHitMargin || inCellX > iconX+copyIconSize+copyIconHitMargin {
			return
		}
		ui.ClipboardSetText(text)
		ui.ShowToast(tv, texts.T().Copied, ui.ToastSuccess)
	})
}

// copyCellLayout returns where the button is in the cell and how wide the text may be:
// right after the text, or at the right edge when the text does not fit
func copyCellLayout(tv *ui.Table, col int, text string) (iconX int, textWidth int) {
	pad := tv.CellPadding()
	maxText := columnVisibleWidth(tv, col) - pad*2 - copyIconGap - copyIconSize
	textWidth, _, err := ui.MeasureText(tv.FontFamily(), tv.FontSize(), text)
	if err != nil || textWidth > maxText {
		textWidth = max(maxText, 0)
	}
	return pad + textWidth + copyIconGap, textWidth
}

// drawCopyCell draws the value, as the table does, and the copy button after it
func drawCopyCell(cnv *ui.Canvas, tv *ui.Table, row int, col int) {
	text := tv.GetCellText2(row, col)
	pad := tv.CellPadding()
	height := tv.RowHeight()
	iconX, textWidth := copyCellLayout(tv, col, text)

	selected := tv.IsRowSelected(row)
	var textColor color.Color = tv.ForegroundColor()
	iconColor := ColorMuted.Get()
	if selected {
		textColor = ui.CurrentPalette().HighlightedText
		iconColor = ui.CurrentPalette().HighlightedText
	}

	cnv.SetHAlign(ui.HAlignLeft)
	cnv.SetVAlign(ui.VAlignCenter)
	cnv.SetFontFamily(tv.FontFamily())
	cnv.SetFontSize(tv.FontSize())
	cnv.SetColor(textColor)
	cnv.DrawText(pad, pad, textWidth, height-pad*2, text)

	drawCopyIcon(cnv, iconX, (height-copyIconSize)/2, iconColor)
}

// drawCopyIcon draws the classic copy icon: a square in front of another one.
// It is drawn with one-pixel rectangles, so it stays sharp at any screen scale.
func drawCopyIcon(cnv *ui.Canvas, x int, y int, col color.Color) {
	s, d := copyIconSquare, copyIconShift
	// The back square: only what the front one does not cover
	cnv.FillRect(x+d, y, s, 1, col)
	cnv.FillRect(x+d+s-1, y, 1, s, col)
	cnv.FillRect(x+d, y, 1, d, col)
	cnv.FillRect(x+s, y+s-1, d, 1, col)
	// The front square
	cnv.FillRect(x, y+d, s, 1, col)
	cnv.FillRect(x, y+d+s-1, s, 1, col)
	cnv.FillRect(x, y+d, 1, s, col)
	cnv.FillRect(x+s-1, y+d, 1, s, col)
}

// columnOffset is where the column starts, as the table places it
func columnOffset(tv *ui.Table, col int) int {
	offset := 0
	for i := 0; i < col; i++ {
		offset += tv.ColumnWidth(i)
	}
	return offset
}

// columnVisibleWidth is the width of the column; the stretched last one
// takes the rest of the table
func columnVisibleWidth(tv *ui.Table, col int) int {
	width := tv.ColumnWidth(col)
	if col == tv.ColumnCount()-1 {
		width = max(width, tv.Width()-columnOffset(tv, col))
	}
	return width
}
