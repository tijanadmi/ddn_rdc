package pdf

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/tijanadmi/ddn_rdc/models"

	"github.com/jung-kurt/gofpdf"
)

/*** Interface za layout tabele, da možemo imati različite layout-e za različite izveštaje ***/
type TableLayout interface {
	DrawHeader(pdf *gofpdf.Fpdf)
	DrawRow(pdf *gofpdf.Fpdf, r models.DetailRow, idx int)
	RowHeight() float64
}

const (
	tableCellHeight      = 4.25
	tableHeaderHeight    = 2 * tableCellHeight
	tableRowHeight       = 2 * tableCellHeight
	eventTextLineHeight  = 4.0
	eventTextTopGap      = 1.0
	eventTextBottomGap   = 2.0
	tipdHeadingHeight    = 8.0
)

/*Layout za TAČKU 1 i 3*/
type TableLayoutT1 struct {
	TrajLabel string // "(hh:mm)" ili "(dddd:hh)"
}

func (l TableLayoutT1) RowHeight() float64 {
	return tableRowHeight
}

func (l TableLayoutT1) DrawHeader(pdf *gofpdf.Fpdf) {
	left, _, _, _ := pdf.GetMargins()
	pdf.SetX(left)
	pdf.SetFont("DejaVu", "B", 7)
	pdf.SetFillColor(198, 224, 180)

	pdf.CellFormat(6, tableCellHeight, "Р.", "RT", 0, "C", true, 0, "")
	pdf.CellFormat(24, tableCellHeight, "Почетак – Крај", "LRT", 0, "C", true, 0, "")
	pdf.CellFormat(14, tableCellHeight, "Трај.", "LRT", 0, "C", true, 0, "")
	pdf.CellFormat(60, tableCellHeight, "Објекат", "LRT", 0, "", true, 0, "")
	pdf.CellFormat(12, tableCellHeight, "Снага", "LRT", 0, "C", true, 0, "")
	pdf.CellFormat(30, tableCellHeight, "Врста догађаја", "LRT", 0, "", true, 0, "")
	pdf.CellFormat(45, tableCellHeight, "Примарни узрок дог.", "LT", 1, "", true, 0, "")

	pdf.CellFormat(6, tableCellHeight, "бр", "RB", 0, "C", true, 0, "")
	pdf.CellFormat(24, tableCellHeight, "", "LRB", 0, "C", true, 0, "")
	pdf.CellFormat(14, tableCellHeight, l.TrajLabel, "LRB", 0, "", true, 0, "")
	pdf.CellFormat(60, tableCellHeight, "", "LRB", 0, "", true, 0, "")
	pdf.CellFormat(12, tableCellHeight, "MW", "LRB", 0, "C", true, 0, "")
	pdf.CellFormat(30, tableCellHeight, "", "LRB", 0, "", true, 0, "")
	pdf.CellFormat(45, tableCellHeight, "Временски услови", "LB", 1, "", true, 0, "")
}

func (l TableLayoutT1) DrawRow(pdf *gofpdf.Fpdf, r models.DetailRow, idx int) {
	tableRow(pdf, r, idx)
}

/***** Layout za TAČKU 2  *****/
type TableLayoutT2 struct{}

func (l TableLayoutT2) RowHeight() float64 {
	return tableRowHeight
}

func (l TableLayoutT2) DrawHeader(pdf *gofpdf.Fpdf) {
	left, _, _, _ := pdf.GetMargins()
	pdf.SetX(left)
	pdf.SetFont("DejaVu", "B", 7)
	pdf.SetFillColor(198, 224, 180)

	pdf.CellFormat(6, tableCellHeight, "Р.", "RT", 0, "C", true, 0, "")
	pdf.CellFormat(24, tableCellHeight, "Почетак – Крај", "LRT", 0, "C", true, 0, "")
	pdf.CellFormat(14, tableCellHeight, "Трај.", "LRT", 0, "C", true, 0, "")
	pdf.CellFormat(70, tableCellHeight, "Објекат", "LRT", 0, "", true, 0, "")
	pdf.CellFormat(30, tableCellHeight, "Врста догађаја", "LRT", 0, "", true, 0, "")
	pdf.CellFormat(50, tableCellHeight, "Разлог", "LT", 1, "", true, 0, "")

	pdf.CellFormat(6, tableCellHeight, "бр", "RB", 0, "C", true, 0, "")
	pdf.CellFormat(24, tableCellHeight, "", "LRB", 0, "C", true, 0, "")
	pdf.CellFormat(14, tableCellHeight, "hh:mm", "LRB", 0, "", true, 0, "")
	pdf.CellFormat(70, tableCellHeight, "", "LRB", 0, "", true, 0, "")
	pdf.CellFormat(30, tableCellHeight, "", "LRB", 0, "", true, 0, "")
	pdf.CellFormat(50, tableCellHeight, "укључења / искључења", "LB", 1, "", true, 0, "")
}

func (l TableLayoutT2) DrawRow(pdf *gofpdf.Fpdf, r models.DetailRow, idx int) {
	pdf.SetFont("DejaVu", "", 7)

	objekat := fitText(pdf, r.Objekat, 70)
	vrsta := fitText(pdf, r.VrstaDogadjaja, 30)
	polje := fitText(pdf, r.Polje+" "+r.ImePolja, 70)
	razlog := fitText(pdf, r.Razlog, 50)

	// RED 1
	pdf.CellFormat(6, tableCellHeight, fmt.Sprint(idx), "", 0, "C", false, 0, "")
	pdf.CellFormat(24, tableCellHeight, r.Vrepoc, "", 0, "", false, 0, "")
	pdf.CellFormat(14, tableCellHeight, r.Traj, "", 0, "C", false, 0, "")
	pdf.CellFormat(70, tableCellHeight, objekat, "", 0, "", false, 0, "")
	pdf.CellFormat(30, tableCellHeight, vrsta, "", 0, "", false, 0, "")
	pdf.CellFormat(50, tableCellHeight, razlog, "", 1, "", false, 0, "")

	// RED 2
	pdf.CellFormat(6, tableCellHeight, "", "", 0, "", false, 0, "")
	pdf.CellFormat(24, tableCellHeight, r.Vrezav, "", 0, "", false, 0, "")
	pdf.CellFormat(14, tableCellHeight, "", "", 0, "", false, 0, "")
	pdf.CellFormat(70, tableCellHeight, polje, "", 0, "", false, 0, "")
	pdf.CellFormat(30, tableCellHeight, "", "", 0, "", false, 0, "")
	pdf.CellFormat(50, tableCellHeight, "", "", 1, "", false, 0, "")

	resetX(pdf)
}

/*********** Factory: bira layout po Tipd-u  ***********/
func getTableLayout(tipd string) TableLayout {
	switch tipd {
	case "1":
		return TableLayoutT1{TrajLabel: "hh:mm"}
	case "3":
		return TableLayoutT1{TrajLabel: "dddd:hh"}
	case "2":
		return TableLayoutT2{}
	default:
		return TableLayoutT1{TrajLabel: "hh:mm"}
	}
}

/****** Interface za renderovanje različitih TIPD grupa *******/
type TipdRenderer interface {
	Render(pdf *gofpdf.Fpdf, tipd models.TipdGroup)
}

/****** Rander za T4 **************/
type TipdRendererT4 struct{}

func (r TipdRendererT4) Render(pdf *gofpdf.Fpdf, tipd models.TipdGroup) {
	for _, day := range tipd.Days {
		if !hasVerticalSpace(pdf, 6) {
			startTipdTextPage(pdf, tipd, day.Date)
		} else {
			drawDayHeading(pdf, day.Date)
		}

		for _, event := range day.Events {
			if event.Tekst == "" {
				continue
			}

			pdf.SetFont("DejaVu", "", 9)
			lines := splitEventText(pdf, event.Tekst)
			eventHeight := float64(len(lines))*eventTextLineHeight + eventTextBottomGap

			if !hasVerticalSpace(pdf, eventHeight) {
				if eventHeight <= freshTextPageEventCapacity(pdf) ||
					!hasVerticalSpace(pdf, eventTextLineHeight) {
					startTipdTextPage(pdf, tipd, day.Date)
				}
			}

			pdf.SetTextColor(0, 0, 180)
			for _, line := range lines {
				if !hasVerticalSpace(pdf, eventTextLineHeight) {
					startTipdTextPage(pdf, tipd, day.Date)
					pdf.SetTextColor(0, 0, 180)
				}
				pdf.CellFormat(
					0, eventTextLineHeight, strings.TrimRight(line, "\r"),
					"", 1, "", false, 0, "",
				)
			}
			pdf.SetTextColor(0, 0, 0)
			if hasVerticalSpace(pdf, eventTextBottomGap) {
				pdf.Ln(eventTextBottomGap)
			}
		}

		if hasVerticalSpace(pdf, 4) {
			pdf.Ln(4)
		}
	}
}

/*************** Rander za tacke 1,2,3 (koje imaju tabelu) *****************/

type TipdRendererTable struct {
	Layout TableLayout
}

func (r TipdRendererTable) Render(pdf *gofpdf.Fpdf, tipd models.TipdGroup) {
	tableHeaderPrintedOnPage := false
	currentPage := pdf.PageNo()
	rowIdx := 1

	for _, day := range tipd.Days {
		for _, event := range day.Events {
			pdf.SetFont("DejaVu", "", 7)

			if pdf.PageNo() != currentPage {
				tableHeaderPrintedOnPage = false
				currentPage = pdf.PageNo()
			}

			eventHeight := measureTableEventHeight(pdf, r.Layout, event)
			requiredHeight := eventHeight
			if !tableHeaderPrintedOnPage {
				requiredHeight += tableHeaderHeight
			}

			if !hasVerticalSpace(pdf, requiredHeight) {
				if eventHeight <= freshTablePageEventCapacity(pdf) {
					startTipdTablePage(pdf, r.Layout, tipd)
					tableHeaderPrintedOnPage = true
					currentPage = pdf.PageNo()
				} else {
					if tableHeaderPrintedOnPage ||
						!hasVerticalSpace(pdf, tableHeaderHeight+r.Layout.RowHeight()) {
						startTipdTablePage(pdf, r.Layout, tipd)
					} else {
						r.Layout.DrawHeader(pdf)
						pdf.SetFont("DejaVu", "", 7)
					}
					tableHeaderPrintedOnPage = true
					currentPage = pdf.PageNo()

					renderOversizedTableEvent(pdf, r.Layout, tipd, event, &rowIdx)
					currentPage = pdf.PageNo()
					continue
				}
			}

			if !tableHeaderPrintedOnPage {
				r.Layout.DrawHeader(pdf)
				tableHeaderPrintedOnPage = true
			}

			for _, row := range event.Rows {
				row.Traj = formatTrajanje(row.Traj, tipd.Tipd)
				r.Layout.DrawRow(pdf, row, rowIdx)
				rowIdx++
			}

			drawTableEventText(pdf, event.Tekst)
		}
	}
}

/*********** Factory: bira render po tackama ***********/
func getTipdRenderer(tipd string) TipdRenderer {
	switch tipd {
	case "4":
		return TipdRendererT4{}
	case "1", "2", "3":
		return TipdRendererTable{
			Layout: getTableLayout(tipd),
		}
	default:
		return TipdRendererTable{
			Layout: getTableLayout(tipd),
		}
	}
}



func GeneratePiMMReportPDF(report *models.Report) ([]byte, error) {

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(10, 37, 10)
	pdf.SetAutoPageBreak(true, 10)

	// === REGISTRACIJA UTF-8 FONTOVA ===
	pdf.AddUTF8Font("DejaVu", "", "assets/fonts/DejaVuSans.ttf")
	pdf.AddUTF8Font("DejaVu", "B", "assets/fonts/DejaVuSans-Bold.ttf")

	registerPiMMHeader(
		pdf,
		report.StartDate, // npr "01.01.2026"
		report.EndDate,   // npr "31.01.2026"
	)

	pdf.AddPage()
	pdf.SetFont("DejaVu", "", 8)

	for _, tipd := range report.TipdGroups {
		if !hasVerticalSpace(pdf, measureTipdStartHeight(pdf, tipd)) {
			pdf.AddPage()
		}
		drawTipdHeading(pdf, tipd)

		renderer := getTipdRenderer(tipd.Tipd)
		renderer.Render(pdf, tipd)
	}

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func pageBottomY(pdf *gofpdf.Fpdf) float64 {
	_, pageH := pdf.GetPageSize()
	_, _, _, bottom := pdf.GetMargins()
	return pageH - bottom
}

func hasVerticalSpace(pdf *gofpdf.Fpdf, height float64) bool {
	return pdf.GetY()+height <= pageBottomY(pdf)
}

func drawTipdHeading(pdf *gofpdf.Fpdf, tipd models.TipdGroup) {
	pdf.SetFont("DejaVu", "B", 11)
	pdf.SetTextColor(0, 128, 0)
	pdf.Cell(0, 8, fmt.Sprintf("Тачка: %s : %s", tipd.Tipd, tipd.Naziv))
	pdf.SetTextColor(0, 0, 0)
	pdf.Ln(tipdHeadingHeight)
}

func startTipdTablePage(pdf *gofpdf.Fpdf, layout TableLayout, tipd models.TipdGroup) {
	pdf.AddPage()
	drawTipdHeading(pdf, tipd)
	layout.DrawHeader(pdf)
	pdf.SetFont("DejaVu", "", 7)
}

func drawDayHeading(pdf *gofpdf.Fpdf, date string) {
	pdf.SetFont("DejaVu", "B", 9)
	pdf.SetTextColor(0, 128, 0)
	pdf.Cell(0, 6, "Дан: "+date)
	pdf.SetTextColor(0, 0, 0)
	pdf.Ln(6)
	pdf.SetFont("DejaVu", "", 9)
}

func startTipdTextPage(pdf *gofpdf.Fpdf, tipd models.TipdGroup, date string) {
	pdf.AddPage()
	drawTipdHeading(pdf, tipd)
	drawDayHeading(pdf, date)
}

func freshTablePageEventCapacity(pdf *gofpdf.Fpdf) float64 {
	_, top, _, _ := pdf.GetMargins()
	return pageBottomY(pdf) - top - tipdHeadingHeight - tableHeaderHeight
}

func freshTextPageEventCapacity(pdf *gofpdf.Fpdf) float64 {
	_, top, _, _ := pdf.GetMargins()
	return pageBottomY(pdf) - top - tipdHeadingHeight - 6
}

func firstTipdEvent(tipd models.TipdGroup) (models.EventGroup, bool) {
	for _, day := range tipd.Days {
		if len(day.Events) > 0 {
			return day.Events[0], true
		}
	}
	return models.EventGroup{}, false
}

func measureTipdStartHeight(pdf *gofpdf.Fpdf, tipd models.TipdGroup) float64 {
	event, ok := firstTipdEvent(tipd)
	if !ok {
		return tipdHeadingHeight
	}

	if tipd.Tipd == "4" {
		pdf.SetFont("DejaVu", "", 9)
		eventHeight := float64(len(splitEventText(pdf, event.Tekst)))*eventTextLineHeight + eventTextBottomGap
		if eventHeight > freshTextPageEventCapacity(pdf) {
			eventHeight = eventTextLineHeight
		}
		return tipdHeadingHeight + 6 + eventHeight
	}

	layout := getTableLayout(tipd.Tipd)
	pdf.SetFont("DejaVu", "", 7)
	eventHeight := measureTableEventHeight(pdf, layout, event)
	if eventHeight > freshTablePageEventCapacity(pdf) {
		eventHeight = layout.RowHeight()
	}
	return tipdHeadingHeight + tableHeaderHeight + eventHeight
}

func tableTextWidth(pdf *gofpdf.Fpdf) float64 {
	pageW, _ := pdf.GetPageSize()
	left, _, right, _ := pdf.GetMargins()
	return pageW - left - right
}

func splitEventText(pdf *gofpdf.Fpdf, text string) []string {
	if text == "" {
		return nil
	}
	return pdf.SplitText(text, tableTextWidth(pdf))
}

func measureTableEventHeight(
	pdf *gofpdf.Fpdf,
	layout TableLayout,
	event models.EventGroup,
) float64 {
	height := float64(len(event.Rows))*layout.RowHeight() + eventTextTopGap
	if event.Tekst != "" {
		height += float64(len(splitEventText(pdf, event.Tekst))) * eventTextLineHeight
		height += eventTextBottomGap
	}
	return height
}

func drawTableEventText(pdf *gofpdf.Fpdf, text string) {
	pdf.Ln(eventTextTopGap)
	if text == "" {
		return
	}

	pdf.MultiCell(0, eventTextLineHeight, text, "", "", false)
	drawTableBottom(pdf)
	pdf.Ln(eventTextBottomGap)
}

func renderOversizedTableEvent(
	pdf *gofpdf.Fpdf,
	layout TableLayout,
	tipd models.TipdGroup,
	event models.EventGroup,
	rowIdx *int,
) {
	for _, row := range event.Rows {
		if !hasVerticalSpace(pdf, layout.RowHeight()) {
			startTipdTablePage(pdf, layout, tipd)
		}

		row.Traj = formatTrajanje(row.Traj, tipd.Tipd)
		layout.DrawRow(pdf, row, *rowIdx)
		(*rowIdx)++
	}

	pdf.Ln(eventTextTopGap)
	lines := splitEventText(pdf, event.Tekst)
	for _, line := range lines {
		if !hasVerticalSpace(pdf, eventTextLineHeight) {
			startTipdTablePage(pdf, layout, tipd)
		}
		pdf.CellFormat(
			0, eventTextLineHeight, strings.TrimRight(line, "\r"),
			"", 1, "", false, 0, "",
		)
	}

	if event.Tekst != "" {
		drawTableBottom(pdf)
		if hasVerticalSpace(pdf, eventTextBottomGap) {
			pdf.Ln(eventTextBottomGap)
		}
	}
}

func registerPiMMHeader(
	pdf *gofpdf.Fpdf,
	periodFrom string,
	periodTo string,
) {
	pdf.SetHeaderFunc(func() {

		left, _, right, _ := pdf.GetMargins()
		pageW, _ := pdf.GetPageSize()

		rightBlockWidth := 20.0
		rightX := pageW - right - rightBlockWidth

		pdf.SetFont("DejaVu", "", 8)

		// ===== LEVI DEO HEADERA =====
		pdf.SetY(8)
		pdf.SetX(left)
		pdf.CellFormat(0, 4, `А.Д. "Електромрежа Србије"`, "", 0, "L", false, 0, "")

		// ===== DESNI BLOK (LEVO PORAVNAT) =====
		pdf.SetXY(rightX, 8)
		pdf.CellFormat(rightBlockWidth, 4, fmt.Sprintf("Страна: %d/{nb}", pdf.PageNo()), "", 1, "L", false, 0, "")
		pdf.SetX(rightX)
		pdf.CellFormat(rightBlockWidth, 4, time.Now().Format("02.01.2006."), "", 1, "L", false, 0, "")
		pdf.SetX(rightX)
		pdf.CellFormat(rightBlockWidth, 4, "PI_MM", "", 1, "L", false, 0, "")

		// ===== NASLOV =====
		// pdf.Ln(1)

		pdf.SetTextColor(192, 0, 0) // crvena EMS stil

		pdf.SetX(left) // ✅ KRITIČNO
		pdf.SetFont("DejaVu", "B", 12)
		pdf.CellFormat(0, 6, "МЕСЕЧНИ ИЗВЕШТАЈ", "", 1, "C", false, 0, "")
		pdf.SetX(left) // ✅ KRITIČNO
		pdf.SetFont("DejaVu", "", 10)
		pdf.CellFormat(0, 5, fmt.Sprintf("за период од %s до %s", periodFrom, periodTo), "", 1, "C", false, 0, "")
		pdf.SetTextColor(0, 0, 0)
		pdf.Ln(5)
	})

	pdf.AliasNbPages("")
}

func tableHeader(pdf *gofpdf.Fpdf) {
	left, _, _, _ := pdf.GetMargins()
	pdf.SetX(left)

	pdf.SetFont("DejaVu", "B", 7)

	pdf.SetFillColor(198, 224, 180) // svetlo zelena (kao u Excelu)
	pdf.SetTextColor(0, 0, 0)       // crni tekst
	pdf.SetDrawColor(0, 0, 0)       // crne ivice

	// PRVI RED
	pdf.CellFormat(6, 5, "Р.", "RT", 0, "C", true, 0, "")
	pdf.CellFormat(24, 5, "Почетак – Крај", "LRT", 0, "C", true, 0, "")
	pdf.CellFormat(14, 5, "Трај.", "LRT", 0, "C", true, 0, "")
	pdf.CellFormat(60, 5, "Објекат", "LRT", 0, "", true, 0, "")
	pdf.CellFormat(12, 5, "Снага", "LRT", 0, "C", true, 0, "")
	pdf.CellFormat(30, 5, "Врста догађаја", "LRT", 0, "", true, 0, "")
	pdf.CellFormat(45, 5, "Примарни узрок дог.", "LT", 1, "", true, 0, "")

	// DRUGI RED
	pdf.CellFormat(6, 5, "бр", "RB", 0, "C", true, 0, "")
	pdf.CellFormat(24, 5, "", "LRB", 0, "C", true, 0, "")
	pdf.CellFormat(14, 5, "hh:mm", "LRB", 0, "C", true, 0, "")
	pdf.CellFormat(60, 5, "", "LRB", 0, "", true, 0, "")
	pdf.CellFormat(12, 5, "MW", "LRB", 0, "C", true, 0, "")
	pdf.CellFormat(30, 5, "", "LRB", 0, "", true, 0, "")
	pdf.CellFormat(45, 5, "Временски услови", "LB", 1, "", true, 0, "")

	pdf.SetFillColor(255, 255, 255)

	resetX(pdf) // ✅ KRITIČNO

}

func tableRow(pdf *gofpdf.Fpdf, r models.DetailRow, idx int) {
	pdf.SetFont("DejaVu", "", 7)
	// pdf.SetTextColor(0, 0, 180)

	objekat := fitText(pdf, r.Objekat, 62)
	vrsta := fitText(pdf, r.VrstaDogadjaja, 30)
	uzrok := fitText(pdf, r.Uzrok, 45)
	vremUsl := fitText(pdf, r.VremUsl, 45)
	polje := fitText(pdf, r.Polje+" "+r.ImePolja, 62)

	// RED 1
	pdf.CellFormat(6, tableCellHeight, fmt.Sprint(idx), "", 0, "C", false, 0, "")
	pdf.CellFormat(24, tableCellHeight, r.Vrepoc, "", 0, "", false, 0, "")
	pdf.CellFormat(14, tableCellHeight, r.Traj, "", 0, "C", false, 0, "")
	pdf.CellFormat(60, tableCellHeight, objekat, "", 0, "", false, 0, "")
	pdf.CellFormat(12, tableCellHeight, r.Snaga, "", 0, "C", false, 0, "")
	pdf.CellFormat(30, tableCellHeight, vrsta, "", 0, "", false, 0, "")
	pdf.CellFormat(45, tableCellHeight, uzrok, "", 1, "", false, 0, "")

	// RED 2
	pdf.CellFormat(6, tableCellHeight, "", "", 0, "", false, 0, "")
	pdf.CellFormat(24, tableCellHeight, r.Vrezav, "", 0, "", false, 0, "")
	pdf.CellFormat(14, tableCellHeight, "", "", 0, "", false, 0, "")
	pdf.CellFormat(60, tableCellHeight, polje, "", 0, "", false, 0, "")
	pdf.CellFormat(12, tableCellHeight, "", "", 0, "", false, 0, "")
	pdf.CellFormat(30, tableCellHeight, "", "", 0, "", false, 0, "")
	pdf.CellFormat(45, tableCellHeight, vremUsl, "", 1, "", false, 0, "")

	// pdf.SetTextColor(0, 0, 0)

	resetX(pdf)

}

func drawTableBottom(pdf *gofpdf.Fpdf) {
	left, _, _, _ := pdf.GetMargins()
	pdf.SetX(left)

	pdf.CellFormat(6, 0, "", "T", 0, "", false, 0, "")
	pdf.CellFormat(25, 0, "", "T", 0, "", false, 0, "")
	pdf.CellFormat(10, 0, "", "T", 0, "", false, 0, "")
	pdf.CellFormat(62, 0, "", "T", 0, "", false, 0, "")
	pdf.CellFormat(12, 0, "", "T", 0, "", false, 0, "")
	pdf.CellFormat(30, 0, "", "T", 0, "", false, 0, "")
	pdf.CellFormat(45, 0, "", "T", 1, "", false, 0, "")

	resetX(pdf)
}

func multiCellTable(
	pdf *gofpdf.Fpdf,
	w, lineHeight float64,
	txt, border, align string,
) {
	x := pdf.GetX()
	y := pdf.GetY()

	pdf.MultiCell(w, lineHeight, txt, border, align, false)
	pdf.SetXY(x+w, y)
}

func getLineCount(pdf *gofpdf.Fpdf, w float64, text string) int {
	lines := pdf.SplitLines([]byte(text), w)
	return len(lines)
}
func max(nums ...int) int {
	m := nums[0]
	for _, n := range nums {
		if n > m {
			m = n
		}
	}
	return m
}
