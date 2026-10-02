package forms

import (
	"github.com/ipoluianov/altsysinfo/texts"
	"github.com/ipoluianov/altsysinfo/widgets"
	"github.com/ipoluianov/nui/ui"
)

// CenterWidget shows the information of the selected category
type CenterWidget struct {
	ui.Widget

	mode string
	// The information the page waits for; nil when shown
	pending *cacheEntry
}

func NewCenterWidget() *CenterWidget {
	var c CenterWidget
	c.InitWidget()
	c.SetPanelPadding(0)
	c.SetXExpandable(true)
	c.SetYExpandable(true)
	return &c
}

func (c *CenterWidget) SetWidget(w ui.Widgeter) {
	c.RemoveAllWidgets()
	c.AddWidget(0, 0, w)
}

func (c *CenterWidget) Mode() string {
	return c.mode
}

// Refresh loads the information of the category again
func (c *CenterWidget) Refresh() {
	if key, ok := cacheKey(c.mode); ok {
		cache.drop(key)
	}
	c.SetMode(c.mode)
}

// Rebuild shows the category again with the information already loaded,
// e.g. in the colors of a new theme
func (c *CenterWidget) Rebuild() {
	c.SetMode(c.mode)
}

// SetMode shows the category; while its information is loading, a note says so
func (c *CenterWidget) SetMode(mode string) {
	c.mode = mode
	c.pending = nil
	key, ok := cacheKey(mode)
	if !ok {
		c.show(mode, loadResult{})
		return
	}
	e := cache.get(key)
	if e.ready() {
		c.show(mode, e.res)
		return
	}
	c.pending = e
	c.showLoading()
	go func() {
		<-e.done
		ui.Invoke(func() {
			// The user may have gone to another category meanwhile
			if c.pending == e && c.mode == mode {
				c.pending = nil
				c.show(mode, e.res)
			}
		})
	}()
}

func (c *CenterWidget) showLoading() {
	panel := ui.NewPanel()
	lbl := ui.NewLabel(texts.T().Loading)
	lbl.SetForegroundColor(widgets.ColorMuted.Get())
	panel.AddWidget(0, 0, lbl)
	panel.AddWidget(1, 0, ui.NewVSpacer())
	c.SetWidget(panel)
}

func (c *CenterWidget) show(mode string, res loadResult) {
	switch mode {
	case "common", "ram":
		if res.err != nil {
			c.SetWidget(widgets.NewWidgetDetails(nil, res.err))
		} else if mode == "common" {
			c.SetWidget(widgets.NewWidgetCommonInfo(res.info))
		} else {
			c.SetWidget(widgets.NewWidgetRamInfo(res.info))
		}
	case "pcidev":
		c.SetWidget(widgets.NewWidgetPciDevInfo())
	default:
		c.SetWidget(widgets.NewWidgetDetails(res.tables, res.err))
	}
	if f := c.Form(); f != nil {
		f.UpdateLayout()
		f.Update()
	}
}
