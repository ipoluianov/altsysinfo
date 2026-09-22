package report

import (
	"bytes"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/go-pdf/fpdf"
	"github.com/ipoluianov/altsysinfo/system"
)

const (
	pageWidth = 180.0 // A4 width minus 15 mm margins
	rowHeight = 7.0
)

// GeneratePDF builds a PDF report describing the system.
func GeneratePDF(info system.Info) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(15, 15, 15)
	pdf.SetAutoPageBreak(true, 15)
	tr := pdf.UnicodeTranslatorFromDescriptor("")

	generated := time.Now().Format("2006-01-02 15:04")
	pdf.SetFooterFunc(func() {
		pdf.SetY(-12)
		pdf.SetFont("Helvetica", "", 8)
		pdf.SetTextColor(128, 128, 128)
		pdf.CellFormat(pageWidth/2, 5, "SysInfo report - "+generated, "", 0, "L", false, 0, "")
		pdf.CellFormat(pageWidth/2, 5, "Page "+strconv.Itoa(pdf.PageNo()), "", 0, "R", false, 0, "")
	})

	pdf.AddPage()

	pdf.SetFont("Helvetica", "B", 20)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(pageWidth, 10, "System Report", "", 1, "L", false, 0, "")
	pdf.SetFont("Helvetica", "", 10)
	pdf.SetTextColor(100, 100, 100)
	pdf.CellFormat(pageWidth, 6, "Generated: "+generated, "", 1, "L", false, 0, "")
	pdf.Ln(4)

	section := func(title string) {
		pdf.Ln(3)
		pdf.SetFont("Helvetica", "B", 13)
		pdf.SetTextColor(0, 0, 0)
		pdf.CellFormat(pageWidth, 8, title, "B", 1, "L", false, 0, "")
		pdf.Ln(2)
	}

	table := func(headers []string, widths []float64, rows [][]string) {
		if len(rows) == 0 {
			pdf.SetFont("Helvetica", "I", 10)
			pdf.SetTextColor(128, 128, 128)
			pdf.CellFormat(pageWidth, rowHeight, "No data available", "", 1, "L", false, 0, "")
			return
		}
		pdf.SetFont("Helvetica", "B", 10)
		pdf.SetTextColor(0, 0, 0)
		pdf.SetFillColor(230, 230, 230)
		for i, h := range headers {
			pdf.CellFormat(widths[i], rowHeight, h, "1", 0, "L", true, 0, "")
		}
		pdf.Ln(-1)
		pdf.SetFont("Helvetica", "", 10)
		for _, row := range rows {
			for i, v := range row {
				pdf.CellFormat(widths[i], rowHeight, fitText(pdf, tr(v), widths[i]-2), "1", 0, "L", false, 0, "")
			}
			pdf.Ln(-1)
		}
	}

	hostname, _ := os.Hostname()

	section("Common")
	table([]string{"Name", "Value"}, []float64{50, 130}, [][]string{
		{"Host name", hostname},
		{"OS / Arch", runtime.GOOS + " / " + runtime.GOARCH},
		{"CPU Model", info.CpuInfo.ModelStr},
		{"CPU Cores", strconv.Itoa(info.CpuInfo.Cores)},
		{"RAM Total", formatMB(info.RamInfo.Total)},
		{"RAM Used", formatMB(info.RamInfo.Used)},
		{"RAM Free", formatMB(info.RamInfo.Free)},
	})

	section("Memory Modules")
	var ramRows [][]string
	for i, ram := range info.RamDevices {
		ramRows = append(ramRows, []string{
			strconv.Itoa(i + 1),
			ram.Model,
			formatMB(ram.Size),
			strconv.FormatUint(ram.Speed, 10) + " MT/s",
		})
	}
	table([]string{"#", "Model", "Size", "Speed"}, []float64{10, 90, 40, 40}, ramRows)

	section("Drives")
	var driveRows [][]string
	for _, drive := range info.Drives {
		driveRows = append(driveRows, []string{
			drive.Name,
			drive.Model,
			drive.Type,
			strconv.FormatUint(drive.Size/1024/1024/1024, 10) + " GB",
		})
	}
	table([]string{"Name", "Model", "Type", "Size"}, []float64{25, 105, 20, 30}, driveRows)

	section("Graphics")
	var gpuRows [][]string
	for _, gpu := range info.GPUs {
		gpuRows = append(gpuRows, []string{gpu.Model, gpu.Vendor + ":" + gpu.Device, gpu.Driver})
	}
	table([]string{"Model", "ID", "Driver"}, []float64{125, 30, 25}, gpuRows)

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func formatMB(bytes uint64) string {
	return strconv.FormatUint(bytes/1024/1024, 10) + " MB"
}

// fitText truncates text with "..." so that it fits into width.
func fitText(pdf *fpdf.Fpdf, text string, width float64) string {
	if pdf.GetStringWidth(text) <= width {
		return text
	}
	for len(text) > 0 && pdf.GetStringWidth(text+"...") > width {
		text = text[:len(text)-1]
	}
	return text + "..."
}
