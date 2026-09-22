package forms

import (
	"os"
	"time"

	"github.com/ipoluianov/altsysinfo/report"
	"github.com/ipoluianov/altsysinfo/system"
	"github.com/u00io/nui/nui"
	"github.com/u00io/nuiforms/ui"
)

type TopWidget struct {
	ui.Widget
}

func NewTopWidget() *TopWidget {
	var c TopWidget
	c.InitWidget()

	btnPdfReport := ui.NewButton("PDF report")
	btnPdfReport.SetOnClick(c.savePdfReport)
	c.AddWidget(0, 0, btnPdfReport)
	c.AddWidget(0, 1, ui.NewHSpacer())
	return &c
}

func (c *TopWidget) savePdfReport() {
	info, err := system.GetInfo()
	if err != nil {
		ui.ShowMessageBox(c, "Error", "Cannot collect system information: "+err.Error())
		return
	}

	data, err := report.GeneratePDF(info)
	if err != nil {
		ui.ShowMessageBox(c, "Error", "Cannot generate report: "+err.Error())
		return
	}

	path, err := nui.SaveFileDialog(nil, nui.SaveFileDialogOptions{
		Title:           "Save PDF report",
		DefaultFileName: "sysinfo-report-" + time.Now().Format("2006-01-02") + ".pdf",
		Filters: []nui.FileDialogFilter{
			{DisplayName: "PDF files", Patterns: []string{"*.pdf"}},
		},
	})
	if err != nil {
		ui.ShowMessageBox(c, "Error", "Cannot open save dialog: "+err.Error())
		return
	}
	if path == "" {
		return
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		ui.ShowMessageBox(c, "Error", "Cannot save report: "+err.Error())
		return
	}
	ui.ShowMessageBox(c, "PDF report", "Report saved to "+path)
}
