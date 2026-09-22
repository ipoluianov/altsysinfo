package forms

import "github.com/u00io/nuiforms/ui"

type LeftWidget struct {
	ui.Widget

	lvItems *ui.Table
}

func NewLeftWidget(onModeChanged func(mode string)) *LeftWidget {
	var c LeftWidget
	c.InitWidget()
	c.lvItems = ui.NewTable()
	c.AddWidget(0, 0, c.lvItems)
	c.SetMaxWidth(250)

	c.lvItems.SetRowCount(10)
	c.lvItems.SetColumnCount(1)
	c.lvItems.SetColumnWidth(0, 240)
	c.lvItems.SetAllowScroll(false, true)
	c.lvItems.SetColumnName(0, "Categories")

	c.lvItems.SetCellText2(0, 0, "Common")
	c.lvItems.SetCellText2(1, 0, "RAM")
	c.lvItems.SetCellText2(2, 0, "PCI Vendor/Device")

	c.lvItems.SetOnSelectionChanged(func(row, col int) {
		switch row {
		case 0:
			onModeChanged("common")
		case 1:
			onModeChanged("ram")
		case 2:
			onModeChanged("pcidev")
		}

	})
	return &c
}
