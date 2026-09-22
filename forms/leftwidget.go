package forms

import (
	"github.com/ipoluianov/altsysinfo/system"
	"github.com/u00io/nuiforms/ui"
)

type LeftWidget struct {
	ui.Widget

	lvItems *ui.Table
}

type categoryItem struct {
	mode string
	name string
}

func NewLeftWidget(onModeChanged func(mode string)) *LeftWidget {
	var c LeftWidget
	c.InitWidget()
	c.lvItems = ui.NewTable()
	c.AddWidget(0, 0, c.lvItems)
	c.SetMaxWidth(250)

	items := []categoryItem{
		{mode: "common", name: "Common"},
		{mode: "ram", name: "RAM"},
	}
	for _, category := range system.DetailCategories() {
		items = append(items, categoryItem{mode: category.ID, name: category.Name})
	}
	items = append(items, categoryItem{mode: "pcidev", name: "PCI Vendor/Device"})

	c.lvItems.SetRowCount(len(items))
	c.lvItems.SetColumnCount(1)
	c.lvItems.SetColumnWidth(0, 240)
	c.lvItems.SetAllowScroll(false, true)
	c.lvItems.SetColumnName(0, "Categories")

	for i, item := range items {
		c.lvItems.SetCellText2(i, 0, item.name)
	}

	c.lvItems.SetOnSelectionChanged(func(row, col int) {
		if row >= 0 && row < len(items) {
			onModeChanged(items[row].mode)
		}
	})
	return &c
}
