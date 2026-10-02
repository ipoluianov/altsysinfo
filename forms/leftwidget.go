package forms

import (
	"image"

	"github.com/ipoluianov/altsysinfo/system"
	"github.com/ipoluianov/altsysinfo/texts"
	"github.com/ipoluianov/altsysinfo/widgets"
	"github.com/ipoluianov/nui/ui"
)

// LeftWidget is the list of the information categories
type LeftWidget struct {
	ui.Widget

	lvItems *ui.Table
	items   []categoryItem
}

type categoryItem struct {
	mode string
}

func NewLeftWidget(onModeChanged func(mode string)) *LeftWidget {
	var c LeftWidget
	c.InitWidget()
	c.SetPanelPadding(0)
	c.SetMinWidth(150)
	c.lvItems = widgets.NewTable()
	c.AddWidget(0, 0, c.lvItems)

	c.items = []categoryItem{{mode: "common"}, {mode: "ram"}}
	for _, category := range system.DetailCategories() {
		c.items = append(c.items, categoryItem{mode: category.ID})
	}
	c.items = append(c.items, categoryItem{mode: "pcidev"})

	c.lvItems.SetRowCount(len(c.items))
	c.lvItems.SetColumnCount(1)
	c.lvItems.SetAllowScroll(false, true)
	c.applyLanguage()
	// forms/icons/cat-<mode>.svg
	for i, item := range c.items {
		setIcon("cat-"+item.mode, func(img image.Image) {
			c.lvItems.SetCellImage(i, 0, img, widgets.CellImageSize())
		})
	}

	c.lvItems.SetOnSelectionChanged(func(row, col int) {
		if row >= 0 && row < len(c.items) {
			onModeChanged(c.items[row].mode)
		}
	})
	return &c
}

// SelectMode selects the category; returns false when there is no such one
func (c *LeftWidget) SelectMode(mode string) bool {
	for i, item := range c.items {
		if item.mode == mode {
			c.lvItems.SetCurrentCell2(i, 0)
			return true
		}
	}
	return false
}

// applyLanguage names the categories in the language of the application
func (c *LeftWidget) applyLanguage() {
	t := texts.T()
	c.lvItems.SetColumnName(0, t.Categories)
	for i, item := range c.items {
		c.lvItems.SetCellText2(i, 0, t.Category.Name(item.mode))
	}
}

func (c *LeftWidget) applyTheme() {
	widgets.ApplyTableTheme(c.lvItems)
}

func (c *LeftWidget) FocusTable() {
	c.lvItems.Focus()
}
