package pdf

import (
	"bytes"
	"fmt"
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

/*Layout za TAČKU 1 i 3*/
type TableLayoutT1 struct {
	TrajLabel string // "(hh:mm)" ili "(dddd:hh)"
}

func (l TableLayoutT1) RowHeight() float64 {
	return 10
}

func (l TableLayoutT1) DrawHeader(pdf *gofpdf.Fpdf) {
	left, _, _, _ := pdf.GetMargins()
	pdf.SetX(left)
	pdf.SetFont("DejaVu", "B", 7)
	pdf.SetFillColor(198, 224, 180)

	pdf.CellFormat(6, 5, "Р.", "RT", 0, "C", true, 0, "")
	pdf.CellFormat(24, 5, "Почетак – Крај", "LRT", 0, "C", true, 0, "")
	pdf.CellFormat(14, 5, "Трај.", "LRT", 0, "C", true, 0, "")
	pdf.CellFormat(60, 5, "Објекат", "LRT", 0, "", true, 0, "")
	pdf.CellFormat(12, 5, "Снага", "LRT", 0, "C", true, 0, "")
	pdf.CellFormat(30, 5, "Врста догађаја", "LRT", 0, "", true, 0, "")
	pdf.CellFormat(45, 5, "Примарни узрок дог.", "LT", 1, "", true, 0, "")

	pdf.CellFormat(6, 5, "бр", "RB", 0, "C", true, 0, "")
	pdf.CellFormat(24, 5, "", "LRB", 0, "C", true, 0, "")
	pdf.CellFormat(14, 5, l.TrajLabel, "LRB", 0, "", true, 0, "")
	pdf.CellFormat(60, 5, "", "LRB", 0, "", true, 0, "")
	pdf.CellFormat(12, 5, "MW", "LRB", 0, "C", true, 0, "")
	pdf.CellFormat(30, 5, "", "LRB", 0, "", true, 0, "")
	pdf.CellFormat(45, 5, "Временски услови", "LB", 1, "", true, 0, "")
}

func (l TableLayoutT1) DrawRow(pdf *gofpdf.Fpdf, r models.DetailRow, idx int) {
	tableRow(pdf, r, idx)
}

/***** Layout za TAČKU 2  *****/
type TableLayoutT2 struct{}

func (l TableLayoutT2) RowHeight() float64 {
	return 10
}

func (l TableLayoutT2) DrawHeader(pdf *gofpdf.Fpdf) {
	left, _, _, _ := pdf.GetMargins()
	pdf.SetX(left)
	pdf.SetFont("DejaVu", "B", 7)
	pdf.SetFillColor(198, 224, 180)

	pdf.CellFormat(6, 5, "Р.", "RT", 0, "C", true, 0, "")
	pdf.CellFormat(24, 5, "Почетак – Крај", "LRT", 0, "C", true, 0, "")
	pdf.CellFormat(14, 5, "Трај.", "LRT", 0, "C", true, 0, "")
	pdf.CellFormat(70, 5, "Објекат", "LRT", 0, "", true, 0, "")
	pdf.CellFormat(30, 5, "Врста догађаја", "LRT", 0, "", true, 0, "")
	pdf.CellFormat(50, 5, "Разлог", "LT", 1, "", true, 0, "")

	pdf.CellFormat(6, 5, "бр", "RB", 0, "C", true, 0, "")
	pdf.CellFormat(24, 5, "", "LRB", 0, "C", true, 0, "")
	pdf.CellFormat(14, 5, "hh:mm", "LRB", 0, "", true, 0, "")
	pdf.CellFormat(70, 5, "", "LRB", 0, "", true, 0, "")
	pdf.CellFormat(30, 5, "", "LRB", 0, "", true, 0, "")
	pdf.CellFormat(50, 5, "укључења / искључења", "LB", 1, "", true, 0, "")
}

func (l TableLayoutT2) DrawRow(pdf *gofpdf.Fpdf, r models.DetailRow, idx int) {
	pdf.SetFont("DejaVu", "", 7)

	objekat := fitText(pdf, r.Objekat, 70)
	vrsta := fitText(pdf, r.VrstaDogadjaja, 30)
	polje := fitText(pdf, r.Polje+" "+r.ImePolja, 70)
	razlog := fitText(pdf, r.Razlog, 50)

	// RED 1
	pdf.CellFormat(6, 5, fmt.Sprint(idx), "", 0, "C", false, 0, "")
	pdf.CellFormat(24, 5, r.Vrepoc, "", 0, "", false, 0, "")
	pdf.CellFormat(14, 5, r.Traj, "", 0, "C", false, 0, "")
	pdf.CellFormat(70, 5, objekat, "", 0, "", false, 0, "")
	pdf.CellFormat(30, 5, vrsta, "", 0, "", false, 0, "")
	pdf.CellFormat(50, 5, razlog, "", 1, "", false, 0, "")

	// RED 2
	pdf.CellFormat(6, 5, "", "", 0, "", false, 0, "")
	pdf.CellFormat(24, 5, r.Vrezav, "", 0, "", false, 0, "")
	pdf.CellFormat(14, 5, "", "", 0, "", false, 0, "")
	pdf.CellFormat(70, 5, polje, "", 0, "", false, 0, "")
	pdf.CellFormat(30, 5, "", "", 0, "", false, 0, "")
	pdf.CellFormat(50, 5, "", "", 1, "", false, 0, "")

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

func GeneratePiMMReportPDF(report *models.Report) ([]byte, error) {

	var currentTipd *models.TipdGroup
	tableHeaderPrintedOnPage := false

	currentPage := 0

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(10, 40, 10) // ⬅ veća gornja margina zbog headera
	pdf.SetAutoPageBreak(true, 15)

	// === PAGE BREAK HANDLER ===
	// pdf.SetAcceptPageBreakFunc(func() bool {
	// 	pdf.AddPage()

	// 	tableHeaderPrintedOnPage = false

	// 	if currentTipd != nil {
	// 		pdf.SetFont("DejaVu", "B", 11)
	// 		pdf.Cell(0, 8, fmt.Sprintf(
	// 			"Тачка: %s : %s",
	// 			currentTipd.Tipd,
	// 			currentTipd.Naziv,
	// 		))
	// 		pdf.Ln(10)

	// 		// ⬅⬅⬅ KLJUČNO: vrati font za normalan sadržaj
	// 		pdf.SetFont("DejaVu", "", 8)
	// 	}

	// 	return false
	// })

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
		currentTipd = &tipd
		rowIdx := 1 // redni broj unutar tipd grupe

		// TIPD HEADER (prva stranica tog TIPD-a)
		pdf.SetFont("DejaVu", "B", 11)
		pdf.Cell(0, 8, fmt.Sprintf("Тачка: %s : %s", tipd.Tipd, tipd.Naziv))
		pdf.Ln(10)

		tableHeaderPrintedOnPage = false

		layout := getTableLayout(tipd.Tipd)

		for _, day := range tipd.Days {
			for _, event := range day.Events {
				// ===== TABELA REDOVI =====
				for _, row := range event.Rows {
					row.Traj = formatTrajanje(row.Traj, tipd.Tipd)

					// 1. Ako smo već na novoj strani
					if pdf.PageNo() != currentPage {
						tableHeaderPrintedOnPage = false
						currentPage = pdf.PageNo()
					}

					// 2. PROVERI DA LI RED STANE (može da doda stranu)
					// ensureSpaceForTableRow(pdf, 10, &tableHeaderPrintedOnPage, currentTipd)
					ensureSpaceForTableRow(pdf, layout.RowHeight(), &tableHeaderPrintedOnPage, currentTipd)

					// 3. Ako je ensureSpaceForTableRow napravio novu stranu
					if pdf.PageNo() != currentPage {
						tableHeaderPrintedOnPage = false
						currentPage = pdf.PageNo()
					}

					// 4. Header tabele — TAČNO JEDNOM
					if !tableHeaderPrintedOnPage {
						// tableHeader(pdf)
						layout.DrawHeader(pdf)
						tableHeaderPrintedOnPage = true
					}

					// 5. Red
					// tableRow(pdf, row, rowIdx)
					layout.DrawRow(pdf, row, rowIdx)
					rowIdx++
				}

				pdf.Ln(2)
				/*** текс на крају сваког догађаја ***/
				pdf.SetFont("DejaVu", "", 8)
				pdf.MultiCell(0, 5, event.Tekst, "", "", false)
				drawTableBottom(pdf)
				pdf.Ln(3)

			}
		}
	}

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func ensureSpaceForTableRow(
	pdf *gofpdf.Fpdf,
	rowHeight float64,
	headerPrinted *bool,
	currentTipd *models.TipdGroup,
) {
	_, pageH := pdf.GetPageSize()
	_, _, _, bottom := pdf.GetMargins()

	if pdf.GetY()+rowHeight > pageH-bottom {

		pdf.AddPage()

		// TIPD header na vrhu nove strane
		if currentTipd != nil {
			pdf.SetFont("DejaVu", "B", 11)
			pdf.Cell(0, 8, fmt.Sprintf(
				"Тачка: %s : %s",
				currentTipd.Tipd,
				currentTipd.Naziv,
			))
			pdf.Ln(10)
			pdf.SetFont("DejaVu", "", 8)
		}

		// *headerPrinted = false
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

		pdf.SetX(left) // ✅ KRITIČNO
		pdf.SetFont("DejaVu", "B", 12)
		pdf.CellFormat(0, 6, "МЕСЕЧНИ ИЗВЕШТАЈ", "", 1, "C", false, 0, "")
		pdf.SetX(left) // ✅ KRITIČNO
		pdf.SetFont("DejaVu", "", 10)
		pdf.CellFormat(0, 5, fmt.Sprintf("за период од %s до %s", periodFrom, periodTo), "", 1, "C", false, 0, "")
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

	objekat := fitText(pdf, r.Objekat, 62)
	vrsta := fitText(pdf, r.VrstaDogadjaja, 30)
	uzrok := fitText(pdf, r.Uzrok, 45)
	vremUsl := fitText(pdf, r.VremUsl, 45)
	polje := fitText(pdf, r.Polje+" "+r.ImePolja, 62)

	// RED 1
	pdf.CellFormat(6, 5, fmt.Sprint(idx), "", 0, "C", false, 0, "")
	pdf.CellFormat(24, 5, r.Vrepoc, "", 0, "", false, 0, "")
	pdf.CellFormat(14, 5, r.Traj, "", 0, "C", false, 0, "")
	pdf.CellFormat(60, 5, objekat, "", 0, "", false, 0, "")
	pdf.CellFormat(12, 5, r.Snaga, "", 0, "C", false, 0, "")
	pdf.CellFormat(30, 5, vrsta, "", 0, "", false, 0, "")
	pdf.CellFormat(45, 5, uzrok, "", 1, "", false, 0, "")

	// RED 2
	pdf.CellFormat(6, 5, "", "", 0, "", false, 0, "")
	pdf.CellFormat(24, 5, r.Vrezav, "", 0, "", false, 0, "")
	pdf.CellFormat(14, 5, "", "", 0, "", false, 0, "")
	pdf.CellFormat(60, 5, polje, "", 0, "", false, 0, "")
	pdf.CellFormat(12, 5, "", "", 0, "", false, 0, "")
	pdf.CellFormat(30, 5, "", "", 0, "", false, 0, "")
	pdf.CellFormat(45, 5, vremUsl, "", 1, "", false, 0, "")

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

// func tableRow(pdf *gofpdf.Fpdf, r models.DetailRow, idx int) {
// 	pdf.SetFont("DejaVu", "", 7)

// 	lineHeight := 4.0

// 	// računamo maksimalan broj linija (obe vrste reda)
// 	maxLines := max(
// 		getLineCount(pdf, 45, r.Objekat),
// 		getLineCount(pdf, 50, r.Uzrok),
// 		getLineCount(pdf, 45, r.Polje+" "+r.ImePolja),
// 		getLineCount(pdf, 50, r.VremUsl),
// 	)

// 	rowHeight := float64(maxLines) * lineHeight
// 	halfHeight := rowHeight / 2

// 	x := pdf.GetX()
// 	y := pdf.GetY()

// 	// ==== OKVIRI ====
// 	widths := []float64{6, 30, 15, 45, 12, 32, 50}
// 	xi := x
// 	for _, w := range widths {
// 		pdf.Rect(xi, y, w, rowHeight, "D")
// 		xi += w
// 	}

// 	// ==== RED 1 ====
// 	pdf.SetXY(x, y)
// 	pdf.CellFormat(6, halfHeight, fmt.Sprint(idx), "", 0, "C", false, 0, "")
// 	pdf.CellFormat(30, halfHeight, r.Vrepoc, "", 0, "", false, 0, "")
// 	pdf.CellFormat(15, halfHeight, r.Traj, "", 0, "C", false, 0, "")

// 	multiCellTable(pdf, 45, lineHeight, r.Objekat, "", "L")
// 	pdf.CellFormat(12, halfHeight, r.Snaga, "", 0, "C", false, 0, "")
// 	multiCellTable(pdf, 32, lineHeight, r.VrstaDogadjaja, "", "L")
// 	multiCellTable(pdf, 50, lineHeight, r.Uzrok, "", "L")

// 	// ==== RED 2 ====
// 	pdf.SetXY(x, y+halfHeight)
// 	pdf.CellFormat(6, halfHeight, "", "", 0, "", false, 0, "")
// 	pdf.CellFormat(30, halfHeight, r.Vrezav, "", 0, "", false, 0, "")
// 	pdf.CellFormat(15, halfHeight, "", "", 0, "", false, 0, "")

// 	multiCellTable(pdf, 45, lineHeight, r.Polje+" "+r.ImePolja, "", "L")
// 	pdf.CellFormat(12, halfHeight, "", "", 0, "", false, 0, "")
// 	pdf.CellFormat(32, halfHeight, "", "", 0, "", false, 0, "")
// 	multiCellTable(pdf, 50, lineHeight, r.VremUsl, "", "L")

// 	// pomeri kursor ispod reda
// 	pdf.SetXY(x, y+rowHeight)
// }

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
