package widgets

import (
	"sort"
	"strings"

	"github.com/ipoluianov/altsysinfo/system"
	"github.com/u00io/nuiforms/ui"
)

type WidgetPciDevInfo struct {
	ui.Widget

	panelFilter *ui.Panel
	txtVendor   *ui.TextBox
	txtDevice   *ui.TextBox
	txtName     *ui.TextBox

	lvItems *ui.Table
}

func NewWidgetPciDevInfo() *WidgetPciDevInfo {
	var c WidgetPciDevInfo
	c.InitWidget()
	c.panelFilter = ui.NewPanel()

	c.panelFilter.AddWidget(0, 0, ui.NewLabel("VEN"))
	c.txtVendor = ui.NewTextBox()
	c.txtVendor.SetHint("Vendor ID")
	c.txtVendor.SetMinWidth(100)
	c.txtVendor.SetXExpandable(false)
	c.panelFilter.AddWidget(1, 0, c.txtVendor)

	c.panelFilter.AddWidget(0, 1, ui.NewLabel("DEV"))
	c.txtDevice = ui.NewTextBox()
	c.txtDevice.SetHint("Device ID")
	c.txtDevice.SetMinWidth(100)
	c.txtDevice.SetXExpandable(false)
	c.panelFilter.AddWidget(1, 1, c.txtDevice)

	c.panelFilter.AddWidget(0, 2, ui.NewLabel("Vendor/Device Name"))
	c.txtName = ui.NewTextBox()
	c.txtName.SetHint("Vendor/Device Name")
	c.txtName.SetMinWidth(100)
	c.txtName.SetXExpandable(true)
	c.panelFilter.AddWidget(1, 2, c.txtName)

	c.panelFilter.SetXExpandable(true)
	c.panelFilter.SetYExpandable(false)

	c.panelFilter.AddWidget(0, 2, ui.NewHSpacer())

	c.txtDevice.SetOnTextChanged(func() {
		c.LoadPciDevInfo()
	})
	c.txtVendor.SetOnTextChanged(func() {
		c.LoadPciDevInfo()
	})
	c.txtName.SetOnTextChanged(func() {
		c.LoadPciDevInfo()
	})

	c.lvItems = ui.NewTable()
	c.AddWidget(0, 0, c.panelFilter)
	c.AddWidget(1, 0, c.lvItems)
	c.SetXExpandable(true)
	c.SetYExpandable(true)

	c.lvItems.SetSelectingRows(true)

	c.lvItems.SetRowCount(10)
	c.lvItems.SetColumnCount(4)
	c.lvItems.SetColumnWidth(0, 70)
	c.lvItems.SetColumnWidth(1, 70)
	c.lvItems.SetColumnWidth(2, 300)
	c.lvItems.SetColumnWidth(3, 600)
	c.lvItems.SetAllowScroll(false, true)
	c.lvItems.SetColumnName(0, "Vendor ID")
	c.lvItems.SetColumnName(1, "Device ID")
	c.lvItems.SetColumnName(2, "Vendor Name")
	c.lvItems.SetColumnName(3, "Device Name")

	c.LoadPciDevInfo()
	return &c
}

func (c *WidgetPciDevInfo) LoadPciDevInfo() {
	list := system.PCIList()
	filtered := make([]system.PciInfo, 0)
	filterVendor := c.txtVendor.Text()
	filterDevice := c.txtDevice.Text()
	filterName := c.txtName.Text()
	filterVendor = strings.ToLower(filterVendor)
	filterDevice = strings.ToLower(filterDevice)
	filterName = strings.ToLower(filterName)
	filterVendor = strings.TrimSpace(filterVendor)
	filterDevice = strings.TrimSpace(filterDevice)
	filterName = strings.TrimSpace(filterName)

	for _, item := range list {
		if filterVendor != "" && !strings.Contains(strings.ToLower(item.VendorIDStr()), filterVendor) {
			continue
		}
		if filterDevice != "" && !strings.Contains(strings.ToLower(item.DeviceIDStr()), filterDevice) {
			continue
		}
		if filterName != "" && !strings.Contains(strings.ToLower(item.VendorName), filterName) && !strings.Contains(strings.ToLower(item.DeviceName), filterName) {
			continue
		}
		filtered = append(filtered, item)
	}

	sort.Slice(filtered, func(i, j int) bool {
		if filtered[i].VendorID != filtered[j].VendorID {
			return filtered[i].VendorID < filtered[j].VendorID
		}
		return filtered[i].DeviceID < filtered[j].DeviceID
	})

	c.lvItems.SetRowCount(len(filtered))
	for i, item := range filtered {
		c.lvItems.SetCellText2(i, 0, item.VendorIDStr())
		c.lvItems.SetCellText2(i, 1, item.DeviceIDStr())
		c.lvItems.SetCellText2(i, 2, item.VendorName)
		c.lvItems.SetCellText2(i, 3, item.DeviceName)
	}
}
