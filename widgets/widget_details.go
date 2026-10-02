package widgets

import (
	"unicode/utf8"

	"github.com/ipoluianov/altsysinfo/system"
	"github.com/ipoluianov/altsysinfo/texts"
	"github.com/ipoluianov/nui/ui"
)

// WidgetDetails shows the tables of a details category, loaded by DetailCategory.Load
type WidgetDetails struct {
	ui.Widget
}

func NewWidgetDetails(tables []*system.DetailTable, err error) *WidgetDetails {
	var c WidgetDetails
	c.InitWidget()
	c.SetXExpandable(true)
	c.SetYExpandable(true)

	if err != nil {
		c.AddWidget(0, 0, ui.NewLabel(texts.T().CannotLoad(err.Error())))
		c.AddWidget(1, 0, ui.NewVSpacer())
		return &c
	}

	row := 0
	for _, t := range tables {
		if len(tables) > 1 && t.Title != "" {
			title := ui.NewLabel(t.Title)
			title.SetFontSize(title.FontSize() + 2)
			c.AddWidget(row, 0, title)
			row++
		}
		c.AddWidget(row, 0, newDetailTableView(t))
		row++
	}
	return &c
}

func newDetailTableView(t *system.DetailTable) *ui.Table {
	tv := NewTable()
	tv.SetColumnCount(len(t.Columns))
	for i, name := range t.Columns {
		tv.SetColumnName(i, texts.T().Column(name))
	}

	if len(t.Rows) == 0 {
		tv.SetRowCount(1)
		tv.SetCellText2(0, 0, texts.T().NoData)
		tv.SetCellColor(0, 0, ColorMuted.Get())
	} else {
		tv.SetRowCount(len(t.Rows))
		for r, cells := range t.Rows {
			for col, text := range cells {
				if col < len(t.Columns) {
					tv.SetCellText2(r, col, text)
				}
			}
			if len(cells) == 1 {
				tv.SetCellColor(r, 0, ColorSection.Get())
			}
		}
	}

	// Size columns by content: about 10 px per character
	for col, name := range t.Columns {
		width := utf8.RuneCountInString(name)
		for _, cells := range t.Rows {
			if col < len(cells) {
				width = max(width, utf8.RuneCountInString(cells[col]))
			}
		}
		px := width*10 + 16
		px = max(px, 60)
		if len(t.Columns) > 2 {
			px = min(px, 520)
		} else {
			px = min(px, 900)
		}
		tv.SetColumnWidth(col, px)
	}
	return tv
}
