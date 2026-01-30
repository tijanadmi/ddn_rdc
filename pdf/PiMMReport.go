package pdf

import (
	"bytes"
	"fmt"

	"github.com/tijanadmi/ddn_rdc/models"

	"github.com/jung-kurt/gofpdf"
)

func GeneratePiMMReportPDF(report *models.Report) ([]byte, error) {
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.SetMargins(10, 10, 10)
	pdf.SetAutoPageBreak(true, 10)

	// === REGISTRACIJA UTF-8 FONTOVA ===
	pdf.AddUTF8Font("DejaVu", "", "assets/fonts/DejaVuSans.ttf")
	pdf.AddUTF8Font("DejaVu", "B", "assets/fonts/DejaVuSans-Bold.ttf")

	pdf.AddPage()
	pdf.SetFont("DejaVu", "", 9)

	// === NASLOV ===
	pdf.SetFont("DejaVu", "B", 14)
	pdf.Cell(0, 8, "PI MM IZVEŠTAJ")
	pdf.Ln(10)

	pdf.SetFont("DejaVu", "", 9)

	for _, tipd := range report.TipdGroups {

		// ===== TIPD HEADER =====
		pdf.SetFont("DejaVu", "B", 11)
		pdf.Cell(0, 7, fmt.Sprintf("Tip događaja: %s - %s", tipd.Tipd, tipd.Naziv))
		pdf.Ln(8)

		for _, day := range tipd.Days {

			// ===== DATUM =====
			pdf.SetFont("DejaVu", "B", 10)
			// pdf.Cell(0, 6, day.Date.Format("02.01.2006"))
			pdf.Cell(0, 6, day.Date)
			pdf.Ln(6)

			for _, event := range day.Events {

				// ===== EVENT HEADER =====
				pdf.SetFont("DejaVu", "B", 9)
				pdf.Cell(20, 6, "ID1:")
				pdf.SetFont("DejaVu", "", 9)
				pdf.Cell(30, 6, fmt.Sprint(event.ID1))
				pdf.Ln(5)

				pdf.MultiCell(0, 5, event.Tekst, "", "", false)
				pdf.Ln(2)

				// ===== TABELA HEADER =====
				pdf.SetFont("DejaVu", "B", 8)
				tableHeader(pdf)

				pdf.SetFont("DejaVu", "", 8)

				for _, row := range event.Rows {
					tableRow(pdf, row)
				}

				pdf.Ln(4)
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

func tableHeader(pdf *gofpdf.Fpdf) {
	pdf.CellFormat(25, 6, "Vreme početka", "1", 0, "C", false, 0, "")
	pdf.CellFormat(25, 6, "Vreme završetka", "1", 0, "C", false, 0, "")
	pdf.CellFormat(15, 6, "Trajanje", "1", 0, "C", false, 0, "")
	pdf.CellFormat(50, 6, "Objekat", "1", 0, "", false, 0, "")
	pdf.CellFormat(30, 6, "Polje", "1", 0, "", false, 0, "")
	pdf.CellFormat(15, 6, "Snaga", "1", 0, "C", false, 0, "")
	pdf.CellFormat(40, 6, "Vrsta događaja", "1", 0, "", false, 0, "")
	pdf.CellFormat(40, 6, "Uzrok", "1", 1, "", false, 0, "")
}

func tableRow(pdf *gofpdf.Fpdf, r models.DetailRow) {
	pdf.CellFormat(25, 6, r.Vrepoc, "1", 0, "", false, 0, "")
	pdf.CellFormat(25, 6, r.Vrezav, "1", 0, "", false, 0, "")
	pdf.CellFormat(15, 6, r.Traj, "1", 0, "C", false, 0, "")
	pdf.CellFormat(50, 6, r.Objekat, "1", 0, "", false, 0, "")
	pdf.CellFormat(30, 6, r.Polje, "1", 0, "", false, 0, "")
	pdf.CellFormat(15, 6, r.Snaga, "1", 0, "C", false, 0, "")
	pdf.CellFormat(40, 6, r.VrstaDogadjaja, "1", 0, "", false, 0, "")
	pdf.CellFormat(40, 6, r.Uzrok, "1", 1, "", false, 0, "")
}

// func formatTime(t time.Time) string {
// 	if t.IsZero() {
// 		return ""
// 	}
// 	return t.Format("02.01.2006 15:04")
// }

// func formatSnaga(s sql.NullFloat64) string {
// 	if !s.Valid {
// 		return ""
// 	}
// 	return fmt.Sprintf("%.2f", s.Float64)
// }
