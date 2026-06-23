package pdf

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/tijanadmi/ddn_rdc/models"

	"github.com/jung-kurt/gofpdf"
)

func GenerateShiftReportPDF(
	report *models.ShiftReport,
) ([]byte, error) {

	pdf := gofpdf.New("P", "mm", "A4", "")

	dispeceri := extractDispeceri(&report.Smena)
	// topMargin := calcTopMargin(len(dispeceri))

	headerBottom := calcHeaderBottomY(10, len(dispeceri))
	topMargin := headerBottom + 8

	pdf.SetMargins(15, topMargin, 15)
	pdf.SetAutoPageBreak(true, 15)

	pdf.AddUTF8Font("DejaVu", "", "assets/fonts/DejaVuSans.ttf")
	pdf.AddUTF8Font("DejaVu", "B", "assets/fonts/DejaVuSans-Bold.ttf")
	pdf.AddUTF8Font("DejaVu", "BI", "assets/fonts/DejaVuSans-BoldOblique.ttf")

	registerShiftHeader(pdf, &report.Smena, dispeceri)

	pdf.AddPage()
	// fmt.Println("TOP:", pdf.GetY())

	_, top, _, _ := pdf.GetMargins()
	pdf.SetY(top)

	for _, dog := range report.Dogadjaji {

		renderDogadjaj(pdf, &dog)

		pdf.Ln(3)
	}

	ensureSpace(pdf, 50)
	renderShiftFooter(pdf, &report.Smena)

	var buf bytes.Buffer

	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func registerShiftHeader(
	pdf *gofpdf.Fpdf,
	smena *models.Smena,
	dispeceri []string,
) {

	const startY = 10.0
	const lineH = 6.0

	pdf.SetHeaderFunc(func() {

		left, _, right, _ := pdf.GetMargins()
		pageW, _ := pdf.GetPageSize()

		// ===== LEVO =====

		pdf.SetFont("DejaVu", "B", 10)
		pdf.SetXY(left, startY)

		pdf.CellFormat(60, lineH, smena.IdBroj, "", 1, "L", false, 0, "")
		pdf.SetX(left)
		pdf.CellFormat(60, lineH, smena.DatDnevStr, "", 0, "L", false, 0, "")

		// ===== SREDINA =====

		centerX := pageW/2 - 40
		pdf.SetXY(centerX, startY)

		pdf.SetFont("DejaVu", "B", 10)
		pdf.CellFormat(80, lineH, "DEŽURNI:", "", 1, "C", false, 0, "")

		pdf.SetFont("DejaVu", "B", 10) // <- promena (bez italika)

		for _, d := range dispeceri {
			pdf.SetX(centerX)
			pdf.CellFormat(80, lineH, d, "", 1, "C", false, 0, "")
		}

		// ===== DESNO =====

		rightX := pageW - right - 35
		pdf.SetXY(rightX, startY)

		pdf.SetFont("DejaVu", "B", 10)
		pdf.CellFormat(35, lineH, fmt.Sprintf("%d/{nb}", pdf.PageNo()), "", 1, "R", false, 0, "")
		pdf.SetX(rightX)
		pdf.CellFormat(35, lineH, smena.Dan, "", 1, "R", false, 0, "")
		pdf.SetX(rightX)
		pdf.CellFormat(35, lineH, smena.IntSmena, "", 1, "R", false, 0, "")

		// ===== DONJA LINIJA =====

		headerBottom := calcHeaderBottomY(startY, len(dispeceri))

		left, top, right, _ := pdf.GetMargins()

		fmt.Println(
			"Page:",
			pdf.PageNo(),
			"headerBottom:",
			headerBottom,
			"topMargin:",
			top,
		)

		pdf.Line(
			left,
			headerBottom,
			pageW-right,
			headerBottom,
		)
		fmt.Println(
			"HEADER Page:",
			pdf.PageNo(),
			"Y after header:",
			pdf.GetY(),
		)
		pdf.SetY(top)
		fmt.Println(
			"after setY top - HEADER Page:",
			pdf.PageNo(),
			"Y after header:",
			pdf.GetY(),
		)
	})

	pdf.AliasNbPages("")
}

func renderDogadjaj(
	pdf *gofpdf.Fpdf,
	dog *models.DogadjajPDF,
) {

	renderDogadjajHeader(pdf, dog)

	// fmt.Println("dog.Tip:", dog.Tip)

	switch dog.Tip {

	case "O":
		// fmt.Println("dog.ObavBeleske:", dog.ObavBeleske)
		if dog.ObavBeleske != nil && dog.ObavBeleske.TipObv == "B" {
			renderObavBeleska(pdf, dog.ObavBeleske)
		}
		if dog.ObavSlike != nil {
			renderObavSlika(pdf, dog)
		}

	case "2":
		renderIskljucenje(pdf, dog)
	case "5": // TSU
		renderT567(pdf, dog)

	case "6": // TK
		renderT567(pdf, dog)

	case "A": // SOP
		renderT567(pdf, dog)
	case "P": // Prekid proizvodnje
		renderT567(pdf, dog)

	case "D": // Angazovani rukovalac
		if dog.AngazovaniRukovalac != nil {
			renderAngazovaniRukovalac(pdf, dog)
		}
	case "1": // Ispad
		renderIspad(pdf, dog)
	case "7": // Ispad
		renderIspad(pdf, dog)

		// kasnije:
		// case "A": angazovani
		// case "S": slike
		// case "I": iskljucenje
		// case "D": detalji
	}
}

// func renderDogadjajHeader(
// 	pdf *gofpdf.Fpdf,
// 	dog *models.DogadjajPDF,
// ) {

// 	const lineH = 6.0

// 	pdf.SetFont("DejaVu", "B", 10)

// 	left := pdf.GetX()
// 	startY := pdf.GetY()

// 	// ===== 1. Rb (fiksno 6 karaktera) =====
// 	rb := fmt.Sprintf("%-6s", dog.RbDog)

// 	pdf.SetXY(left, startY)
// 	pdf.CellFormat(18, lineH, rb, "", 0, "L", false, 0, "")

// 	// ===== 2. Naslov (indented + wrap) =====
// 	titleX := left + 18 + 3 // 3 blanka

// 	pdf.SetXY(titleX, startY)

// 	pdf.MultiCell(
// 		0,
// 		lineH,
// 		dog.Naslov,
// 		"",
// 		"L",
// 		false,
// 	)

// 	// Y nakon naslova
// 	endY := pdf.GetY()

// 	// pomeramo Y za sledeći blok
// 	pdf.SetY(endY + 2)
// }

func renderDogadjajHeader(pdf *gofpdf.Fpdf, dog *models.DogadjajPDF) {

	const lineH = 6.0

	left := pdf.GetX()
	startY := pdf.GetY()

	bodyX := left

	// =========================
	// 1. IZRAČUN VISINE BLOKA
	// =========================

	linesNaslov := pdf.SplitLines([]byte(dog.Naslov), 160)
	linesPod := pdf.SplitLines([]byte(dog.Podnaslov), 160)

	blockH := float64(len(linesNaslov))*lineH + 2

	if strings.TrimSpace(dog.Podnaslov) != "" {
		blockH += float64(len(linesPod))*lineH + 2
	}

	// + margin za Rb liniju
	blockH += lineH

	// =========================
	// 2. PAGE BREAK PROVERA
	// =========================

	_, pageH := pdf.GetPageSize()
	_, _, _, bottom := pdf.GetMargins()

	if pdf.GetY()+blockH > pageH-bottom {
		pdf.AddPage()
		startY = pdf.GetY()
	}

	// =========================
	// 3. CRTANJE BLOKA
	// =========================

	pdf.SetFont("DejaVu", "B", 10)

	rb := fmt.Sprintf("%-6s", dog.RbDog)

	pdf.SetXY(bodyX, startY)
	pdf.CellFormat(18, lineH, rb, "", 0, "L", false, 0, "")

	titleX := bodyX + 21

	// NASLOV
	pdf.SetXY(titleX, startY)
	pdf.MultiCell(0, lineH, dog.Naslov, "", "L", false)

	y := pdf.GetY()

	// PODNASLOV
	if strings.TrimSpace(dog.Podnaslov) != "" {
		pdf.SetFont("DejaVu", "B", 10)
		pdf.SetXY(titleX, y)
		pdf.MultiCell(0, lineH, dog.Podnaslov, "", "L", false)
		y = pdf.GetY()
	}

	// kraj bloka
	pdf.SetY(y + 2)
}

func renderShiftFooter(pdf *gofpdf.Fpdf, smena *models.Smena) {

	renderShiftComment(pdf, smena)

	renderShiftSignatures(pdf, smena)
}

func renderShiftComment(pdf *gofpdf.Fpdf, smena *models.Smena) {

	if strings.TrimSpace(smena.KomentZat) == "" {
		return
	}

	const lineH = 5.0

	pdf.Ln(10)

	pdf.SetFont("DejaVu", "B", 10)
	pdf.Cell(0, 6, "KOMENTAR KOD PRIMOPREDAJE SMENE:")
	pdf.Ln(6)

	pdf.SetFont("DejaVu", "B", 8)

	pdf.MultiCell(
		0,
		lineH,
		smena.KomentZat,
		"",
		"L",
		false,
	)

	pdf.Ln(10)
}

func estimateSignatureHeight(s *models.Smena) float64 {

	linesLeft := 0
	for _, v := range []string{
		s.PredaoDisp1,
		s.PredaoDisp2,
		s.PredaoDisp3,
	} {
		if strings.TrimSpace(v) != "" {
			linesLeft++
		}
	}

	linesRight := 0
	for _, v := range []string{
		s.PrimDisp1,
		s.PrimDisp2,
		s.PrimDisp3,
	} {
		if strings.TrimSpace(v) != "" {
			linesRight++
		}
	}

	maxLines := linesLeft
	if linesRight > maxLines {
		maxLines = linesRight
	}

	// header + spacing + lines
	return 20 + float64(maxLines)*5 + 10
}

func renderShiftSignatures(pdf *gofpdf.Fpdf, smena *models.Smena) {

	// 🔒 KLJUČ: rezerviši prostor pre crtanja
	need := estimateSignatureHeight(smena)
	ensureSpace(pdf, need)

	startY := pdf.GetY()

	leftX := float64(15)
	rightX := float64(110)

	lineH := 5.0
	signGap := 10.0
	// lineWidth := 60.0

	// =====================
	// LEVO - PREDAO
	// =====================
	pdf.SetXY(leftX, startY)

	pdf.SetFont("DejaVu", "BI", 10)
	pdf.Cell(80, 6, "SMENU PREDAO:")
	pdf.Ln(8)

	for _, s := range []string{
		smena.PredaoDisp1,
		smena.PredaoDisp2,
		smena.PredaoDisp3,
	} {
		if s == "" {
			continue
		}
		pdf.SetX(leftX)
		pdf.CellFormat(80, lineH, s, "", 1, "L", false, 0, "")

		// linija za potpis
		// y := pdf.GetY()

		// pdf.Line(leftX, y, leftX+lineWidth, y)

		pdf.Ln(signGap)
	}

	// =====================
	// DESNO - PRIMIO
	// =====================
	pdf.SetXY(rightX, startY)

	pdf.SetFont("DejaVu", "BI", 10)
	pdf.Cell(80, 6, "SMENU PRIMIO:")
	pdf.Ln(8)

	for _, s := range []string{
		smena.PrimDisp1,
		smena.PrimDisp2,
		smena.PrimDisp3,
	} {
		if s == "" {
			continue
		}
		// pdf.SetX(rightX)
		// pdf.Cell(80, 5, s)
		// pdf.Ln(12)
		pdf.SetX(rightX)
		pdf.CellFormat(80, lineH, s, "", 1, "L", false, 0, "")

		// linija za potpis
		// y := pdf.GetY()

		// pdf.Line(rightX, y, rightX+lineWidth, y)
		pdf.Ln(signGap)
	}

	// =====================
	// KRAJ BLOCKA
	// =====================
	pdf.Ln(5)
}

func ensureSpace(pdf *gofpdf.Fpdf, need float64) {
	_, pageH := pdf.GetPageSize()
	_, _, _, bottom := pdf.GetMargins()

	if pdf.GetY()+need > pageH-bottom {
		pdf.AddPage()
	}
}

func extractDispeceri(smena *models.Smena) []string {
	valid := make([]string, 0, 4)

	for _, d := range []string{
		smena.DezDisp1Ime,
		smena.DezDisp2Ime,
		smena.DezDisp3Ime,
		smena.DezDisp4Ime,
	} {
		if strings.TrimSpace(d) != "" {
			valid = append(valid, d)
		}
	}

	return valid
}

func calcHeaderBottomY(startY float64, numDisp int) float64 {
	const lineH = 6.0

	return startY +
		lineH + // "DEŽURNI:"
		float64(numDisp)*lineH +
		3.0 // padding
}

// func calcTopMargin(numDisp int) float64 {
// 	const startY = 10.0
// 	const lineH = 6.0

// 	return startY +
// 		lineH +
// 		float64(numDisp)*lineH +
// 		8.0 // dodatni spacing ispod headera
// }

func renderObavBeleska(
	pdf *gofpdf.Fpdf,
	ob *models.ObavBeleska,
) {
	const lineH = 5.0

	left := pdf.GetX()

	dopWidth := 18.0
	bodyX := left + 21

	startY := pdf.GetY()

	// ===== LEVA KOLONA (DOPUNA) =====

	leftEndY := startY

	if strings.TrimSpace(ob.Dopuna) != "" {

		pdf.SetFont("DejaVu", "", 9)

		pdf.SetXY(left, startY)

		pdf.MultiCell(
			dopWidth,
			lineH,
			"Dopuna: "+ob.Dopuna,
			"",
			"L",
			false,
		)

		leftEndY = pdf.GetY()
	}

	// ===== DESNA KOLONA =====

	pdf.SetXY(bodyX, startY)

	if strings.TrimSpace(ob.TekstObv) != "" {

		pdf.SetFont("DejaVu", "", 9)

		pdf.MultiCell(
			0,
			lineH,
			ob.TekstObv,
			"",
			"L",
			false,
		)
	}

	rightEndY := pdf.GetY()

	if strings.TrimSpace(ob.Napomena) != "" {

		pdf.SetFont("DejaVu", "BI", 9)

		pdf.SetXY(bodyX, rightEndY)

		pdf.MultiCell(
			0,
			lineH,
			ob.Napomena,
			"",
			"L",
			false,
		)

		rightEndY = pdf.GetY()
	}

	// ===== NASTAVI ISPOD DUŽE KOLONE =====

	if leftEndY > rightEndY {
		pdf.SetY(leftEndY + 2)
	} else {
		pdf.SetY(rightEndY + 2)
	}
}

func renderIskljucenje(pdf *gofpdf.Fpdf, dog *models.DogadjajPDF) {
	const lineH = 5.0

	left := pdf.GetX()
	startY := pdf.GetY()

	bodyX := left + 21 // isti indent kao ObavBeleska

	// ===== NASLOV =====
	pdf.SetFont("DejaVu", "", 9)
	pdf.SetXY(bodyX, startY)

	pdf.MultiCell(0, lineH, dog.Grazlog+" / "+dog.Razlog, "", "L", false)

	y := pdf.GetY()

	// ===== UZROK TEKST =====
	if strings.TrimSpace(dog.UzrokTekst) != "" {
		pdf.SetFont("DejaVu", "", 9)
		pdf.SetXY(bodyX, y)

		pdf.MultiCell(0, lineH, dog.UzrokTekst, "", "L", false)
		y = pdf.GetY()
	}

	pdf.Ln(2)

	// ===== OBJEKTI =====
	for _, obj := range dog.Objekti {

		// naziv objekta
		pdf.SetFont("DejaVu", "B", 9)
		pdf.SetX(bodyX)
		pdf.Cell(0, lineH, obj.Naziv)
		pdf.Ln(lineH)

		// ===== STAVKE (GRID 4 KOLONE) =====
		for _, s := range obj.Stavke {

			startRowY := pdf.GetY()

			col1X := bodyX
			col2X := bodyX + 10
			col3X := bodyX + 25
			col4X := bodyX + 45

			// --- kolona 1 (dopuna) ---
			pdf.SetFont("DejaVu", "", 9)
			pdf.SetXY(col1X, startRowY)
			pdf.CellFormat(10, lineH, s.DopunaDaNe, "", 0, "L", false, 0, "")

			// --- kolona 2 ---
			pdf.SetXY(col2X, startRowY)
			pdf.SetFont("DejaVu", "", 9)
			pdf.CellFormat(10, lineH, s.Vrepoc, "", 0, "L", false, 0, "")

			// --- kolona 3 ---
			pdf.SetXY(col3X, startRowY)
			pdf.SetFont("DejaVu", "", 9)
			pdf.CellFormat(20, lineH, "-  "+s.Vrezav, "", 0, "L", false, 0, "")

			// --- kolona 4 (tekst koji se lomi) ---
			pdf.SetXY(col4X, startRowY)
			pdf.SetFont("DejaVu", "", 9)

			pdf.MultiCell(
				0,
				lineH,
				s.RecenicaMan,
				"",
				"L",
				false,
			)

			// ključ: poravnanje reda
			rowEndY := pdf.GetY()
			pdf.SetY(rowEndY + 1)
		}

		pdf.Ln(2)
	}

	// ===== MAN TEKST (kao frontend box) =====
	if strings.TrimSpace(dog.ManTekst) != "" {

		pdf.Ln(2)

		pdf.SetFont("DejaVu", "", 9)
		// pdf.SetTextColor(80, 80, 80)

		pdf.SetX(bodyX)
		pdf.MultiCell(
			0,
			lineH,
			dog.ManTekst,
			"",
			"L",
			false,
		)

		pdf.SetTextColor(0, 0, 0)
	}
}

func renderT567(
	pdf *gofpdf.Fpdf,
	dog *models.DogadjajPDF,
) {

	const lineH = 5.0

	left := pdf.GetX()
	bodyX := left + 21

	// =====================
	// DETALJI T5/T6/T7
	// =====================

	for _, d := range dog.Detalji {

		rowStartY := pdf.GetY()

		// ----- DOPUNA -----

		if strings.TrimSpace(d.DopunaDaNe) != "" {

			pdf.SetFont("DejaVu", "", 8)

			pdf.SetXY(left, rowStartY)

			pdf.MultiCell(
				18,
				lineH,
				d.DopunaDaNe,
				"",
				"L",
				false,
			)
		}

		// ----- RECENICA 1 -----

		pdf.SetXY(bodyX, rowStartY)

		pdf.SetFont("DejaVu", "B", 9)

		pdf.MultiCell(
			0,
			lineH,
			d.Recenica1,
			"",
			"L",
			false,
		)

		rowEndY := pdf.GetY()

		// ----- RECENICA 2 -----

		if strings.TrimSpace(d.Recenica2) != "" {

			pdf.SetXY(bodyX, rowEndY)

			pdf.SetFont("DejaVu", "", 9)

			pdf.MultiCell(
				0,
				lineH,
				d.Recenica2,
				"",
				"L",
				false,
			)

			rowEndY = pdf.GetY()
		}

		// ----- OPIS -----

		if strings.TrimSpace(d.Opis) != "" {

			pdf.SetXY(bodyX, rowEndY)

			pdf.SetFont("DejaVu", "", 9)

			pdf.MultiCell(
				0,
				lineH,
				d.Opis,
				"",
				"L",
				false,
			)

			rowEndY = pdf.GetY()
		}

		pdf.SetY(rowEndY + 1)
	}

	// =====================
	// UZROK
	// =====================

	if strings.TrimSpace(dog.UzrokTekst) != "" {

		pdf.Ln(2)

		pdf.SetX(bodyX)

		pdf.SetFont("DejaVu", "", 9)

		pdf.MultiCell(
			0,
			lineH,
			dog.UzrokTekst,
			"",
			"L",
			false,
		)
	}

	// =====================
	// MAN TEKST
	// =====================

	if strings.TrimSpace(dog.ManTekst) != "" {

		pdf.Ln(2)

		pdf.SetX(bodyX)

		pdf.SetFont("DejaVu", "", 9)

		pdf.MultiCell(
			0,
			lineH,
			dog.ManTekst,
			"",
			"L",
			false,
		)
	}
}

func renderAngazovaniRukovalac(
	pdf *gofpdf.Fpdf,
	dog *models.DogadjajPDF,
) {
	const lineH = 5.0

	if dog.AngazovaniRukovalac == nil {
		return
	}

	ar := dog.AngazovaniRukovalac

	left := pdf.GetX()
	bodyX := left + 21

	pdf.SetX(bodyX)

	// ===== PODACI =====

	pdf.SetFont("DejaVu", "", 9)

	if ar.VremeNaloga != nil {

		pdf.CellFormat(
			0,
			lineH,
			"Dana: "+ar.VremeNaloga.Format("02.01.2006 15:04"),
			"",
			1,
			"L",
			false,
			0,
			"",
		)
	}

	if ar.ImeNaloga != nil && *ar.ImeNaloga != "" {

		pdf.SetX(bodyX)
		pdf.CellFormat(
			0,
			lineH,
			"Nalog izdat: "+*ar.ImeNaloga,
			"",
			1,
			"L",
			false,
			0,
			"",
		)
	}

	if ar.VremeDolaska != nil {

		pdf.SetX(bodyX)
		pdf.CellFormat(
			0,
			lineH,
			"Vreme dolaska: "+ar.VremeDolaska.Format("02.01.2006 15:04"),
			"",
			1,
			"L",
			false,
			0,
			"",
		)
	}

	if ar.VremeOdlaska != nil {

		pdf.SetX(bodyX)
		pdf.CellFormat(
			0,
			lineH,
			"Vreme odlaska: "+ar.VremeOdlaska.Format("02.01.2006 15:04"),
			"",
			1,
			"L",
			false,
			0,
			"",
		)
	}

	if ar.Rukovalac != nil && *ar.Rukovalac != "" {

		pdf.SetX(bodyX)
		pdf.CellFormat(
			0,
			lineH,
			"Rukovalac: "+*ar.Rukovalac,
			"",
			1,
			"L",
			false,
			0,
			"",
		)
	}

	if ar.Objekat != nil && *ar.Objekat != "" {

		pdf.SetX(bodyX)
		pdf.CellFormat(
			0,
			lineH,
			"Objekat: "+*ar.Objekat,
			"",
			1,
			"L",
			false,
			0,
			"",
		)
	}

	// ===== OPIS =====

	if ar.Opis != nil && strings.TrimSpace(*ar.Opis) != "" {

		pdf.Ln(3)

		pdf.SetX(bodyX)

		pdf.MultiCell(
			0,
			lineH,
			*ar.Opis,
			"",
			"L",
			false,
		)
	}

	pdf.Ln(3)
}

func renderIspad(pdf *gofpdf.Fpdf, dog *models.DogadjajPDF) {
	const lineH = 5.0

	left := pdf.GetX()
	startY := pdf.GetY()

	bodyX := left + 21

	// =========================
	// HEADER
	// =========================
	pdf.SetFont("DejaVu", "", 9)
	pdf.SetXY(bodyX, startY)

	pdf.Ln(2)

	ensureSpace(pdf, float64(len(dog.Detalji))*10+20)

	// =========================
	// 1. UZROK
	// =========================
	pdf.SetFont("DejaVu", "B", 9)

	// naslov
	pdf.SetX(bodyX)
	pdf.Cell(0, lineH, "UZROK POREMEĆAJA I HRONOLOGIJA")
	pdf.Ln(lineH + 2)

	// detalji (grid 2 kolone)
	for _, d := range dog.Detalji {
		ensureSpace(pdf, 15)
		startY := pdf.GetY()

		// layout pozicije
		dopX := bodyX
		textX := bodyX + 8 // Recenica1 uvucena
		subX := bodyX + 10 // Recenica2 i Opis još malo

		// ===== DOPUNA =====
		pdf.SetFont("DejaVu", "", 9)
		pdf.SetXY(dopX, startY)
		pdf.CellFormat(8, lineH, d.DopunaDaNe, "", 0, "L", false, 0, "")

		// ===== RECENICA 1 (glavna) =====
		pdf.SetFont("DejaVu", "B", 9)
		pdf.SetXY(textX, startY)
		pdf.MultiCell(0, lineH, d.Recenica1, "", "L", false)

		y := pdf.GetY()

		// ===== RECENICA 2 =====
		if strings.TrimSpace(d.Recenica2) != "" {
			pdf.SetFont("DejaVu", "", 9)
			pdf.SetXY(subX, y)
			pdf.MultiCell(0, lineH, d.Recenica2, "", "L", false)
			y = pdf.GetY()
		}

		// ===== OPIS =====
		if strings.TrimSpace(d.Opis) != "" {
			pdf.SetFont("DejaVu", "", 9)
			pdf.SetXY(subX, y)
			pdf.MultiCell(0, lineH, d.Opis, "", "L", false)
		}

		pdf.SetY(pdf.GetY() + 1)
	}

	// =========================
	// UZROK TEKST
	// =========================
	if strings.TrimSpace(dog.UzrokTekst) != "" {
		pdf.Ln(1)
		pdf.SetX(bodyX)
		pdf.MultiCell(0, lineH, dog.UzrokTekst, "", "L", false)
	}

	pdf.Ln(2)

	// =========================
	// 2. SANIRANJE
	// =========================
	ensureSpace(pdf, float64(len(dog.Objekti))*20+30)
	if len(dog.Objekti) > 0 {
		pdf.SetFont("DejaVu", "B", 9)
		pdf.SetX(bodyX)
		pdf.Cell(0, lineH, "SANIRANJE POREMEĆAJA")
		pdf.Ln(lineH + 1)
	}

	for _, obj := range dog.Objekti {

		pdf.SetFont("DejaVu", "B", 9)
		pdf.SetX(bodyX)
		pdf.Cell(0, lineH, obj.Naziv)
		pdf.Ln(lineH)

		for _, s := range obj.Stavke {

			startRowY := pdf.GetY()

			col1X := bodyX
			col2X := bodyX + 10
			col3X := bodyX + 25
			col4X := bodyX + 45

			pdf.SetFont("DejaVu", "", 9)

			pdf.SetXY(col1X, startRowY)
			pdf.CellFormat(10, lineH, s.DopunaDaNe, "", 0, "L", false, 0, "")

			pdf.SetXY(col2X, startRowY)
			pdf.CellFormat(12, lineH, s.Vrepoc, "", 0, "L", false, 0, "")

			pdf.SetXY(col3X, startRowY)
			pdf.CellFormat(20, lineH, "- "+s.Vrezav, "", 0, "L", false, 0, "")

			pdf.SetXY(col4X, startRowY)
			pdf.MultiCell(0, lineH, s.RecenicaMan, "", "L", false)

			pdf.SetY(pdf.GetY() + 1)
		}
	}

	// =========================
	// 3. MAN TEKST
	// =========================
	if strings.TrimSpace(dog.ManTekst) != "" {
		pdf.Ln(2)
		pdf.SetX(bodyX)
		pdf.MultiCell(0, lineH, dog.ManTekst, "", "L", false)
	}

	// =========================
	// 4. POSLEDICE
	// =========================
	if strings.TrimSpace(dog.Posledice) != "" {
		ensureSpace(pdf, float64(len(strings.Split(dog.Posledice, "\n")))*6+10)
		pdf.Ln(2)
		pdf.SetX(bodyX)

		pdf.SetFont("DejaVu", "B", 9)
		pdf.Cell(0, lineH, "POSLEDICE")
		pdf.Ln(lineH)

		pdf.SetFont("DejaVu", "", 9)
		pdf.SetX(bodyX)
		pdf.MultiCell(0, lineH, dog.Posledice, "", "L", false)
	}
}

func renderObavSlika(pdf *gofpdf.Fpdf, dog *models.DogadjajPDF) {
	const lineH = 5.0

	if len(dog.ObavSlike) == 0 {
		return
	}

	// =========================
	// SORTIRANJE PO RB
	// =========================
	sort.Slice(dog.ObavSlike, func(i, j int) bool {
		return dog.ObavSlike[i].RB < dog.ObavSlike[j].RB
	})

	left := pdf.GetX()
	bodyX := left + 21

	pageW, _ := pdf.GetPageSize()
	// _, _, _, bottom := pdf.GetMargins()

	maxWidth := pageW - bodyX - 20

	pdf.SetFont("DejaVu", "B", 9)
	pdf.SetX(bodyX)
	pdf.Cell(0, lineH, "OBAVEŠTENJA - SLIKE")
	pdf.Ln(lineH + 2)

	// =========================
	// SLIKE LOOP
	// =========================
	for i, s := range dog.ObavSlike {

		// estimate visine slike (fiksni blok ~60)
		ensureSpace(pdf, 70)

		// naslov slike
		pdf.SetFont("DejaVu", "", 9)
		pdf.SetX(bodyX)
		pdf.Cell(0, lineH, fmt.Sprintf("Slika %d", i+1))
		pdf.Ln(lineH + 1)

		// decode base64
		imgBytes, err := base64.StdEncoding.DecodeString(s.Base64)
		if err != nil {
			continue
		}

		imgOpt := gofpdf.ImageOptions{
			ImageType: s.Format, // "jpg", "png"
			ReadDpi:   true,
		}

		// register image
		imgName := fmt.Sprintf("img_%d_%d", i, time.Now().UnixNano())
		pdf.RegisterImageOptionsReader(imgName, imgOpt, bytes.NewReader(imgBytes))

		info := pdf.GetImageInfo(imgName)

		// scale to fit width
		w := maxWidth
		h := info.Height() * (w / info.Width())

		x := bodyX

		pdf.ImageOptions(
			imgName,
			x,
			pdf.GetY(),
			w,
			h,
			false,
			imgOpt,
			0,
			"",
		)

		pdf.Ln(h + 3)
	}
}
