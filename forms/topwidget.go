package forms

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/ipoluianov/altsysinfo/config"
	"github.com/ipoluianov/altsysinfo/report"
	"github.com/ipoluianov/altsysinfo/texts"
	"github.com/ipoluianov/nui/ui"
)

type TopWidget struct {
	ui.Widget

	btnRefresh   *ui.ToolButton
	btnPdfReport *ui.ToolButton
	btnTray      *ui.ToolButton
}

var lastCreatedTopWidget *TopWidget

func NewTopWidget() *TopWidget {
	var c TopWidget
	c.InitWidget()
	c.SetPanelPadding(0)

	c.btnRefresh = ui.NewToolButton(nil, "", c.onBtnRefresh)
	setIcon("refresh", c.btnRefresh.SetImage)
	c.btnRefresh.SetTooltipFunc(func() string { return withShortcut(texts.T().ToolRefresh, shortcutRefresh) })

	c.btnPdfReport = ui.NewToolButton(nil, "", c.savePdfReport)
	setIcon("export", c.btnPdfReport.SetImage)
	c.btnPdfReport.SetTooltipFunc(func() string { return withShortcut(texts.T().ToolPdfReport, shortcutPdfReport) })

	c.btnTray = ui.NewToolButton(nil, "", c.onBtnTray)
	setIcon("tray", c.btnTray.SetImage)
	c.btnTray.SetTooltipFunc(func() string { return withShortcut(texts.T().ToolTray, shortcutTray) })

	// Flat icons: the toolbar stays light, a button shows its shape under the mouse
	for _, btn := range []*ui.ToolButton{c.btnRefresh, c.btnPdfReport, c.btnTray} {
		btn.SetFlat(true)
	}

	c.AddWidget(0, 0, c.btnRefresh)
	c.AddWidget(0, 1, c.btnPdfReport)
	c.AddWidget(0, 2, ui.NewHSpacer())
	c.AddWidget(0, 3, c.btnTray)

	lastCreatedTopWidget = &c
	return &c
}

func (c *TopWidget) onBtnRefresh() {
	lastCreatedMainWidget.Refresh()
}

func (c *TopWidget) onBtnTray() {
	lastCreatedMainWidget.HideToTray()
}

// savePdfReport makes the report in the background from the loaded information, then asks where to save it
func (c *TopWidget) savePdfReport() {
	e := cache.get(infoKey)
	go func() {
		<-e.done
		var data []byte
		err := e.res.err
		if err == nil {
			data, err = report.GeneratePDF(e.res.info)
		}
		ui.Invoke(func() { c.askSavePdfReport(data, err, e.res.err != nil) })
	}()
}

func (c *TopWidget) askSavePdfReport(data []byte, err error, collectFailed bool) {
	t := texts.T()
	if err != nil && collectFailed {
		ui.ShowMessageBox(c, t.Error, t.CannotCollect(err.Error()))
		return
	}
	if err != nil {
		ui.ShowMessageBox(c, t.Error, t.CannotGenerate(err.Error()))
		return
	}

	name := "sysinfo-report-" + time.Now().Format("2006-01-02") + ".pdf"
	opts := ui.SaveFileDialogOptions{
		Title:           t.SaveReportTitle,
		DefaultFileName: name,
		Filters: []ui.FileDialogFilter{
			{DisplayName: t.PDFFiles, Patterns: []string{"*.pdf"}},
		},
	}
	// The system dialog is a separate window: it would open below a window kept on top
	form := c.Form()
	onTop := config.GetSettings().AlwaysOnTop
	if onTop {
		form.SetAlwaysOnTop(false)
	}
	form.ShowSaveFileDialog(opts, func(path string, err error) {
		if onTop {
			form.SetAlwaysOnTop(true)
		}
		// No system dialog (e.g. Linux without zenity or kdialog): save to the home folder
		noDialog := errors.Is(err, ui.ErrNoFileDialog)
		if err != nil && !noDialog {
			ui.ShowMessageBox(c, t.Error, t.CannotOpenSaveDialog(err.Error()))
			return
		}
		if noDialog {
			home, _ := os.UserHomeDir()
			path = filepath.Join(home, name)
		}
		if path == "" {
			return // cancelled
		}

		if err := os.WriteFile(path, data, 0644); err != nil {
			ui.ShowMessageBox(c, t.Error, t.CannotSave(err.Error()))
			return
		}
		ui.ShowToast(c, t.SavedTo(path), ui.ToastSuccess)
	})
}
