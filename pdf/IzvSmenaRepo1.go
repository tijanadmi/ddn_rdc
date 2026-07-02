package pdf

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/tijanadmi/ddn_rdc/models"

	"github.com/jung-kurt/gofpdf"
)

type PDFBlock struct {
	Height float64
	Draw   func(pdf *gofpdf.Fpdf)
}

func RenderBlocks(pdf *gofpdf.Fpdf, blocks []PDFBlock) {
	for _, b := range blocks {
		// ensureSpace(pdf, b.Height)
		if b.Height > 0 {
			ensureSpace(pdf, 20) // samo mali safety buffer
		}
		b.Draw(pdf)
	}
}

func textHeight(pdf *gofpdf.Fpdf, text string, width float64, lineH float64) float64 {
	if strings.TrimSpace(text) == "" {
		return 0
	}
	lines := pdf.SplitLines([]byte(text), width)
	return float64(len(lines)) * lineH
}

func maxFloat(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

// func ensureSpace(pdf *gofpdf.Fpdf, need float64) {
// 	_, pageH := pdf.GetPageSize()
// 	_, _, _, bottom := pdf.GetMargins()

// 	if pdf.GetY()+need > pageH-bottom {
// 		pdf.AddPage()
// 	}
// }

func buildDogadjajBlocks(pdf *gofpdf.Fpdf, dog *models.DogadjajPDF) []PDFBlock {

	var blocks []PDFBlock

	blocks = append(blocks, buildHeaderBlock(pdf, dog))

	switch dog.Tip {

	case "O":
		if dog.ObavBeleske != nil {
			blocks = append(blocks, buildObavBeleskaBlock(pdf, dog.ObavBeleske))
		}
		if dog.ObavSlike != nil {
			blocks = append(blocks, buildObavSlikaBlock(pdf, dog))
		}

	case "2":
		blocks = append(blocks, buildIskljucenjeBlock(pdf, dog))

	case "5", "6", "A", "P":
		blocks = append(blocks, buildT567Block(pdf, dog))

	case "D":
		if dog.AngazovaniRukovalac != nil {
			blocks = append(blocks, buildAngazovaniBlock(pdf, dog))
		}

	case "1", "7":
		blocks = append(blocks, buildIspadBlock(pdf, dog))
	}

	return blocks
}

func buildHeaderBlock(pdf *gofpdf.Fpdf, dog *models.DogadjajPDF) PDFBlock {

	const lineH = 6.0

	left := pdf.GetX()

	bodyX := left
	titleX := left + 21

	linesNaslov := pdf.SplitLines([]byte(dog.Naslov), 160)
	linesPod := pdf.SplitLines([]byte(dog.Podnaslov), 160)

	h := float64(len(linesNaslov))*lineH + 2
	h += float64(len(linesPod)) * lineH
	h += lineH

	return PDFBlock{
		Height: h,
		Draw: func(pdf *gofpdf.Fpdf) {

			startY := pdf.GetY()

			pdf.SetFont("DejaVu", "B", 10)

			pdf.SetXY(bodyX, startY)
			pdf.CellFormat(18, lineH, dog.RbDog, "", 0, "L", false, 0, "")

			pdf.SetXY(titleX, startY)
			pdf.MultiCell(0, lineH, dog.Naslov, "", "L", false)

			y := pdf.GetY()

			if dog.Podnaslov != "" {
				pdf.SetXY(titleX, y)
				pdf.MultiCell(0, lineH, dog.Podnaslov, "", "L", false)
				y = pdf.GetY()
			}

			pdf.SetY(y + 2)
		},
	}
}

func buildObavBeleskaBlock(pdf *gofpdf.Fpdf, ob *models.ObavBeleska) PDFBlock {

	const lineH = 5.0

	left := pdf.GetX()

	dopWidth := 18.0
	bodyX := left + 21

	pageW, _ := pdf.GetPageSize()
	_, _, right, _ := pdf.GetMargins()

	bodyWidth := pageW - right - bodyX

	// ======================
	// HEIGHT
	// ======================

	leftH := 0.0
	if strings.TrimSpace(ob.Dopuna) != "" {
		leftH = estimateMultiCellHeight(
			pdf,
			"Dopuna: "+ob.Dopuna,
			dopWidth,
			lineH,
		)
	}

	rightH := 0.0

	if strings.TrimSpace(ob.TekstObv) != "" {
		rightH += estimateMultiCellHeight(
			pdf,
			ob.TekstObv,
			bodyWidth,
			lineH,
		)
	}

	if strings.TrimSpace(ob.Napomena) != "" {
		rightH += estimateMultiCellHeight(
			pdf,
			ob.Napomena,
			bodyWidth,
			lineH,
		)
	}

	h := maxFloat(leftH, rightH) + 2

	// ======================
	// DRAW
	// ======================

	return PDFBlock{

		Height: h,

		Draw: func(pdf *gofpdf.Fpdf) {

			startY := pdf.GetY()

			leftEndY := startY

			// ------------------
			// Dopuna
			// ------------------

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

			// ------------------
			// Tekst
			// ------------------

			rightEndY := startY

			if strings.TrimSpace(ob.TekstObv) != "" {

				pdf.SetFont("DejaVu", "", 9)

				pdf.SetXY(bodyX, startY)

				pdf.MultiCell(
					bodyWidth,
					lineH,
					ob.TekstObv,
					"",
					"L",
					false,
				)

				rightEndY = pdf.GetY()
			}

			// ------------------
			// Napomena
			// ------------------

			if strings.TrimSpace(ob.Napomena) != "" {

				pdf.SetFont("DejaVu", "BI", 9)

				pdf.SetXY(bodyX, rightEndY)

				pdf.MultiCell(
					bodyWidth,
					lineH,
					ob.Napomena,
					"",
					"L",
					false,
				)

				rightEndY = pdf.GetY()
			}

			if leftEndY > rightEndY {
				pdf.SetY(leftEndY + 2)
			} else {
				pdf.SetY(rightEndY + 2)
			}
		},
	}
}

func buildObavSlikaBlock(pdf *gofpdf.Fpdf, dog *models.DogadjajPDF) PDFBlock {

	const lineH = 5.0

	if len(dog.ObavSlike) == 0 {
		return PDFBlock{
			Height: 0,
			Draw:   func(pdf *gofpdf.Fpdf) {},
		}
	}

	pageW, _ := pdf.GetPageSize()

	left := pdf.GetX()
	bodyX := left + 21

	maxWidth := pageW - bodyX - 20

	type imgItem struct {
		name string
		w    float64
		h    float64
		rb   int
		opt  gofpdf.ImageOptions
	}

	items := make([]imgItem, 0, len(dog.ObavSlike))

	h := 0.0

	// =========================
	// BUILD + MEASURE PHASE
	// =========================

	for i, s := range dog.ObavSlike {

		imgBytes, err := base64.StdEncoding.DecodeString(s.Base64)
		if err != nil {
			continue
		}

		imgOpt := gofpdf.ImageOptions{
			ImageType: s.Format,
			ReadDpi:   true,
		}

		imgName := fmt.Sprintf("img_%d_%d", i, time.Now().UnixNano())

		pdf.RegisterImageOptionsReader(
			imgName,
			imgOpt,
			bytes.NewReader(imgBytes),
		)

		info := pdf.GetImageInfo(imgName)

		w := maxWidth
		imgH := info.Height() * (w / info.Width())

		// naslov + slika + spacing
		h += lineH + imgH + 5

		items = append(items, imgItem{
			name: imgName,
			w:    w,
			h:    imgH,
			rb:   s.RB,
			opt:  imgOpt,
		})
	}

	// safety padding
	h += 2

	// =========================
	// DRAW PHASE
	// =========================

	return PDFBlock{
		Height: h,

		Draw: func(pdf *gofpdf.Fpdf) {

			for i, it := range items {

				ensureSpace(pdf, it.h+lineH+5)

				pdf.SetX(bodyX)
				pdf.SetFont("DejaVu", "", 9)
				pdf.Cell(0, lineH, fmt.Sprintf("Slika %d", i+1))
				pdf.Ln(lineH + 1)

				x := bodyX
				y := pdf.GetY()

				pdf.ImageOptions(
					it.name,
					x,
					y,
					it.w,
					it.h,
					false,
					it.opt,
					0,
					"",
				)

				pdf.SetY(y + it.h + 3)
			}
		},
	}
}

func buildIskljucenjeBlock(pdf *gofpdf.Fpdf, dog *models.DogadjajPDF) PDFBlock {

	const lineH = 5.0

	left := pdf.GetX()
	bodyX := left + 21

	// =========================
	// 1. HEIGHT CALCULATION
	// =========================

	h := 0.0

	// --- naslov ---
	title := dog.Grazlog + " / " + dog.Razlog

	h += estimateMultiCellHeight(
		pdf,
		title,
		160,
		lineH,
	)

	// --- uzrok ---
	if strings.TrimSpace(dog.UzrokTekst) != "" {
		h += estimateMultiCellHeight(pdf, dog.UzrokTekst, 160, lineH)
	}

	h += 2 // spacing

	// --- objekti + stavke ---
	for _, obj := range dog.Objekti {

		h += lineH + 2

		for _, s := range obj.Stavke {

			text := s.RecenicaMan
			rowH := estimateMultiCellHeight(pdf, text, 145, lineH)

			h += rowH + 1
		}

		h += 2
	}

	// --- MAN tekst ---
	if strings.TrimSpace(dog.ManTekst) != "" {
		h += 2
		h += estimateMultiCellHeight(pdf, dog.ManTekst, 160, lineH)
	}

	// safety padding
	h += 5

	// =========================
	// 2. DRAW FUNCTION
	// =========================

	return PDFBlock{
		Height: h,

		Draw: func(pdf *gofpdf.Fpdf) {

			startY := pdf.GetY()

			// =====================
			// NASLOV (GRAZLOG / RAZLOG)
			// =====================

			pdf.SetFont("DejaVu", "", 9)
			pdf.SetXY(bodyX, startY)
			pdf.MultiCell(
				0,
				lineH,
				dog.Grazlog+" / "+dog.Razlog,
				"",
				"L",
				false,
			)

			y := pdf.GetY()

			// =====================
			// UZROK
			// =====================

			if strings.TrimSpace(dog.UzrokTekst) != "" {
				pdf.SetFont("DejaVu", "", 9)
				pdf.SetXY(bodyX, y)

				pdf.MultiCell(0, lineH, dog.UzrokTekst, "", "L", false)
				y = pdf.GetY()
			}

			pdf.SetY(y + 2)

			// =====================
			// OBJEKTI + GRID
			// =====================

			for _, obj := range dog.Objekti {

				pdf.SetFont("DejaVu", "B", 9)
				pdf.SetX(bodyX)
				pdf.Cell(0, lineH, obj.Naziv)
				pdf.Ln(lineH)

				for _, s := range obj.Stavke {

					startRowY := pdf.GetY()

					col2X := bodyX + 10
					col3X := bodyX + 25
					col4X := bodyX + 45

					// kolona 1
					pdf.SetFont("DejaVu", "", 8)
					pdf.SetXY(left+5, startRowY)
					pdf.CellFormat(10, lineH, s.DopunaDaNe, "", 0, "L", false, 0, "")

					// kolona 2
					pdf.SetXY(col2X, startRowY)
					pdf.SetFont("DejaVu", "", 9)
					pdf.CellFormat(10, lineH, s.Vrepoc, "", 0, "L", false, 0, "")

					// kolona 3
					pdf.SetXY(col3X, startRowY)
					pdf.SetFont("DejaVu", "", 9)
					pdf.CellFormat(20, lineH, "- "+s.Vrezav, "", 0, "L", false, 0, "")

					// kolona 4 (MultiCell)
					pdf.SetXY(col4X, startRowY)
					pdf.SetFont("DejaVu", "", 9)

					pdf.MultiCell(0, lineH, s.RecenicaMan, "", "L", false)

					pdf.SetY(pdf.GetY() + 1)
				}

				pdf.Ln(2)
			}

			// =====================
			// MAN TEKST
			// =====================

			if strings.TrimSpace(dog.ManTekst) != "" {

				pdf.Ln(2)

				pdf.SetFont("DejaVu", "", 9)
				pdf.SetX(bodyX)

				pdf.MultiCell(
					0,
					lineH,
					dog.ManTekst,
					"",
					"L",
					false,
				)
			}
		},
	}
}

func buildT567Block(pdf *gofpdf.Fpdf, dog *models.DogadjajPDF) PDFBlock {

	const lineH = 5.0

	left := pdf.GetX()
	bodyX := left + 21

	// =========================
	// 1. HEIGHT CALCULATION
	// =========================

	h := 0.0

	// svaki detalj = potencijalno više linija
	for _, d := range dog.Detalji {

		rowText := d.Recenica1

		if strings.TrimSpace(d.Recenica2) != "" {
			rowText += "\n" + d.Recenica2
		}

		if strings.TrimSpace(d.Opis) != "" {
			rowText += "\n" + d.Opis
		}

		// dopuna + tekst zajedno ne utiču na širinu iste kolone,
		// ali utiču na visinu reda
		rowH := estimateMultiCellHeight(pdf, rowText, 145, lineH)

		// + baseline padding
		h += rowH + 1
	}

	// uzrok
	if strings.TrimSpace(dog.UzrokTekst) != "" {
		h += estimateMultiCellHeight(pdf, dog.UzrokTekst, 160, lineH)
	}

	// man tekst
	if strings.TrimSpace(dog.ManTekst) != "" {
		h += estimateMultiCellHeight(pdf, dog.ManTekst, 160, lineH)
	}

	// safety spacing
	h += 5

	// =========================
	// 2. DRAW FUNCTION
	// =========================

	return PDFBlock{
		Height: h,

		Draw: func(pdf *gofpdf.Fpdf) {

			// startY := pdf.GetY()

			// =====================
			// DETALJI
			// =====================

			for _, d := range dog.Detalji {

				rowStartY := pdf.GetY()

				// layout X pozicije
				dopX := left + 5
				body := bodyX
				sub := bodyX + 3

				// ----- DOPUNA -----
				if strings.TrimSpace(d.DopunaDaNe) != "" {
					pdf.SetFont("DejaVu", "", 8)
					pdf.SetXY(dopX, rowStartY)
					pdf.CellFormat(18, lineH, d.DopunaDaNe, "", 0, "L", false, 0, "")
				}

				y := rowStartY

				// ----- RECENICA 1 -----
				pdf.SetFont("DejaVu", "B", 9)
				pdf.SetXY(body, y)
				pdf.MultiCell(0, lineH, d.Recenica1, "", "L", false)
				y = pdf.GetY()

				// ----- RECENICA 2 -----
				if strings.TrimSpace(d.Recenica2) != "" {
					pdf.SetFont("DejaVu", "", 9)
					pdf.SetXY(sub, y)
					pdf.MultiCell(0, lineH, d.Recenica2, "", "L", false)
					y = pdf.GetY()
				}

				// ----- OPIS -----
				if strings.TrimSpace(d.Opis) != "" {
					pdf.SetFont("DejaVu", "", 9)
					pdf.SetXY(sub, y)
					pdf.MultiCell(0, lineH, d.Opis, "", "L", false)
					y = pdf.GetY()
				}

				// next row
				pdf.SetY(y + 1)
			}

			// =====================
			// UZROK
			// =====================

			if strings.TrimSpace(dog.UzrokTekst) != "" {
				pdf.Ln(2)
				pdf.SetX(bodyX)

				pdf.SetFont("DejaVu", "", 9)
				pdf.MultiCell(0, lineH, dog.UzrokTekst, "", "L", false)
			}

			// =====================
			// MAN TEKST
			// =====================

			if strings.TrimSpace(dog.ManTekst) != "" {
				pdf.Ln(2)
				pdf.SetX(bodyX)

				pdf.SetFont("DejaVu", "", 9)
				pdf.MultiCell(0, lineH, dog.ManTekst, "", "L", false)
			}

			// final spacing handled by engine
			pdf.Ln(2)
		},
	}
}

func buildAngazovaniBlock(pdf *gofpdf.Fpdf, dog *models.DogadjajPDF) PDFBlock {

	const lineH = 5.0

	if dog.AngazovaniRukovalac == nil {
		return PDFBlock{Height: 0, Draw: func(pdf *gofpdf.Fpdf) {}}
	}

	ar := dog.AngazovaniRukovalac

	bodyX := pdf.GetX() + 21

	// =========================
	// 1. HEIGHT CALCULATION
	// =========================

	h := 0.0

	if ar.VremeNaloga != nil {
		h += lineH
	}

	if ar.ImeNaloga != nil && *ar.ImeNaloga != "" {
		h += lineH
	}

	if ar.VremeDolaska != nil {
		h += lineH
	}

	if ar.VremeOdlaska != nil {
		h += lineH
	}

	if ar.Rukovalac != nil && *ar.Rukovalac != "" {
		h += lineH
	}

	if ar.Objekat != nil && *ar.Objekat != "" {
		h += lineH
	}

	// opis (MultiCell)
	if ar.Opis != nil && strings.TrimSpace(*ar.Opis) != "" {
		h += estimateMultiCellHeight(pdf, *ar.Opis, 160, lineH)
	}

	// spacing
	h += 6

	// =========================
	// 2. DRAW FUNCTION
	// =========================

	return PDFBlock{
		Height: h,

		Draw: func(pdf *gofpdf.Fpdf) {

			startY := pdf.GetY()

			pdf.SetFont("DejaVu", "", 9)
			pdf.SetX(bodyX)

			// ===== PODACI =====

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

			// final spacing handled by engine
			pdf.SetY(startY + h)
		},
	}
}

func buildIspadBlock(pdf *gofpdf.Fpdf, dog *models.DogadjajPDF) PDFBlock {

	const lineH = 5.0

	left := pdf.GetX()
	bodyX := left + 21

	// =========================
	// 1. HEIGHT CALCULATION
	// =========================

	h := 0.0

	// padding top
	h += 2

	// ===== UZROK + HRONOLOGIJA =====
	h += lineH // naslov

	for _, d := range dog.Detalji {

		text := d.Recenica1

		if strings.TrimSpace(d.Recenica2) != "" {
			text += "\n" + d.Recenica2
		}

		if strings.TrimSpace(d.Opis) != "" {
			text += "\n" + d.Opis
		}

		h += estimateMultiCellHeight(pdf, text, 145, lineH) + 1
	}

	// uzrok tekst
	if strings.TrimSpace(dog.UzrokTekst) != "" {
		h += estimateMultiCellHeight(pdf, dog.UzrokTekst, 160, lineH)
	}

	h += 3

	// ===== SANIRANJE =====
	if len(dog.Objekti) > 0 {
		h += lineH // header "SANIRANJE"
	}

	for _, obj := range dog.Objekti {

		h += lineH // naziv objekta

		for _, s := range obj.Stavke {

			rowText := s.RecenicaMan
			rowH := estimateMultiCellHeight(pdf, rowText, 145, lineH)

			h += rowH + 1
		}
	}

	// ===== MAN TEKST =====
	if strings.TrimSpace(dog.ManTekst) != "" {
		h += estimateMultiCellHeight(pdf, dog.ManTekst, 160, lineH)
	}

	// ===== POSLEDICE =====
	if strings.TrimSpace(dog.Posledice) != "" {
		h += lineH // naslov
		h += estimateMultiCellHeight(pdf, dog.Posledice, 160, lineH)
	}

	// safety padding
	h += 6

	// =========================
	// 2. DRAW FUNCTION
	// =========================

	return PDFBlock{
		Height: h,

		Draw: func(pdf *gofpdf.Fpdf) {

			startY := pdf.GetY()

			// =====================
			// HEADER SPACING
			// =====================
			pdf.SetY(startY + 2)

			// =====================
			// UZROK + HRONOLOGIJA
			// =====================

			pdf.SetFont("DejaVu", "B", 9)
			pdf.SetX(bodyX)
			pdf.Cell(0, lineH, "UZROK POREMEĆAJA I HRONOLOGIJA")
			pdf.Ln(lineH + 2)

			for _, d := range dog.Detalji {

				rowStartY := pdf.GetY()

				dopX := left + 5
				textX := bodyX
				subX := bodyX + 3

				// ===== DOPUNA =====
				if strings.TrimSpace(d.DopunaDaNe) != "" {
					pdf.SetFont("DejaVu", "", 8)
					pdf.SetXY(dopX, rowStartY)
					pdf.CellFormat(8, lineH, d.DopunaDaNe, "", 0, "L", false, 0, "")
				}

				y := rowStartY

				// ===== RECENICA 1 =====
				pdf.SetFont("DejaVu", "B", 9)
				pdf.SetXY(textX, y)
				pdf.MultiCell(0, lineH, d.Recenica1, "", "L", false)
				y = pdf.GetY()

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
					y = pdf.GetY()
				}

				pdf.SetY(y + 1)
			}

			// =====================
			// UZROK TEKST
			// =====================

			if strings.TrimSpace(dog.UzrokTekst) != "" {
				pdf.Ln(1)
				pdf.SetX(bodyX)
				pdf.SetFont("DejaVu", "", 9)
				pdf.MultiCell(0, lineH, dog.UzrokTekst, "", "L", false)
			}

			pdf.Ln(3)

			// =====================
			// SANIRANJE
			// =====================

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

					col1X := left + 5
					col2X := bodyX + 10
					col3X := bodyX + 25
					col4X := bodyX + 45

					pdf.SetFont("DejaVu", "", 8)

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

			// =====================
			// MAN TEKST
			// =====================

			if strings.TrimSpace(dog.ManTekst) != "" {
				pdf.Ln(2)
				pdf.SetX(bodyX)
				pdf.SetFont("DejaVu", "", 9)
				pdf.MultiCell(0, lineH, dog.ManTekst, "", "L", false)
			}

			// =====================
			// POSLEDICE
			// =====================

			if strings.TrimSpace(dog.Posledice) != "" {

				pdf.Ln(2)

				pdf.SetX(bodyX)
				pdf.SetFont("DejaVu", "B", 9)
				pdf.Cell(0, lineH, "POSLEDICE")
				pdf.Ln(lineH)

				pdf.SetFont("DejaVu", "", 9)
				pdf.SetX(bodyX)
				pdf.MultiCell(0, lineH, dog.Posledice, "", "L", false)
			}
		},
	}
}
