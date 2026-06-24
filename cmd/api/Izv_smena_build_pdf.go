package api

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tijanadmi/ddn_rdc/models"
	"github.com/tijanadmi/ddn_rdc/pdf"
)

const (
	TipDogIskljucenje        = "2"
	TipDogObavestenjeBeleska = "B"
	TipDogObavestenje        = "O"
	TipDogObavestenjeFaks    = "F"
	TipDogTSU                = "5"
	TipDogTK                 = "6"
	TipDogSOP                = "A"
	TipDogIspad              = "1"
	TipDogPrekidP            = "P"
	TipDogRukovalac          = "D"
)

func (server *Server) getIzvSmenaPDF(ctx *gin.Context) {
	idSmene := ctx.Param("id")
	id, err := strconv.Atoi(idSmene)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// 1. build report
	report, err := server.BuildShiftReport(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// 2. generate PDF
	pdfBytes, err := pdf.GenerateShiftReportPDF(report)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// 3. return PDF
	ctx.Header("Content-Type", "application/pdf")
	ctx.Header("Content-Disposition", fmt.Sprintf(
		"inline; filename=izvestaj_smena_%d.pdf",
		id,
	))

	ctx.Data(http.StatusOK, "application/pdf", pdfBytes)
}

func (server *Server) BuildShiftReport(ctx *gin.Context, idSmene int) (*models.ShiftReport, error) {
	// 1. Učitaj smenu sa osnovnim događajima
	smena, err := server.store.GetSmenaByID(ctx, idSmene)
	if err != nil {
		return nil, err
	}

	report := &models.ShiftReport{
		Smena: *smena,
	}

	// 2. Prođi kroz sve događaje smene
	for _, dog := range smena.Dogadjaji {

		pdfDog, err := server.buildDogadjajPDF(ctx, dog)
		if err != nil {
			return nil, err
		}

		// neki tipovi događaja možda neće ulaziti u PDF
		if pdfDog != nil {
			report.Dogadjaji = append(report.Dogadjaji, *pdfDog)
		}
	}

	report.Proizvodnja, err = server.store.GetProizvodnjaPoSmeni(ctx, idSmene, report.Smena.IDTipSmena)

	return report, nil
}

func (server *Server) buildDogadjajPDF(
	ctx *gin.Context,
	dog models.Dogadjaj,
) (*models.DogadjajPDF, error) {

	pdfDog := &models.DogadjajPDF{
		RbDog:  dog.RbDog,
		Naslov: dog.Naslov,
		TipDog: dog.TipDog,
		Tip:    dog.Tip,
	}

	// fmt.Printf("tip dogadjaja je %d: %+v\n", dog.ID, dog.Tip)

	switch dog.Tip {

	case TipDogPrekidP:
		// TODO
		pp, err := server.store.GetPrekidPById(ctx, dog.ID)
		if err != nil {
			return nil, err
		}
		var detaljT567 []models.DetaljT567

		if pp.PrekidP != nil {
			for _, p := range *pp.PrekidP {
				detaljT567 = append(detaljT567, buildDetaljPrekidP(p, pp))
			}
		}
		copyCommonPDFFields(pdfDog, pp.Podnaslov, pp.UzrokTekst, pp.ManTekst)

		if pp.Posledice != nil {
			pdfDog.Posledice = *pp.Posledice
		}

		pdfDog.Detalji = detaljT567

	case TipDogRukovalac:
		// TODO
		ar, err := server.store.GetAngazovaniRukovaociById(ctx, dog.ID)
		if err != nil {
			return nil, err
		}

		pdfDog.AngazovaniRukovalac = ar.AngazovaniRukovalac

	case TipDogIspad:
		// TODO
		ispad, err := server.store.GetIspadById(ctx, dog.ID)
		if err != nil {
			return nil, err
		}

		// 1. ISPAD detalji
		detalji := buildIspadDetaljiPDF(
			func() []models.Ispad {
				if ispad.Ispad != nil {
					return *ispad.Ispad
				}
				return nil
			}(),
			ispad,
		)

		var objekti []models.ObjekatView
		// 2. COMMON manipulacije (VAŽNO)
		if len(ispad.Manipulacije) > 0 {
			objekti = buildIskljucenjeObjektiPDF(ispad)
		}

		// 3. COMMON fields
		copyCommonPDFFields(pdfDog, ispad.Podnaslov, ispad.UzrokTekst, ispad.ManTekst)

		if ispad.Posledice != nil {
			pdfDog.Posledice = *ispad.Posledice
		}

		pdfDog.Detalji = detalji
		pdfDog.Objekti = objekti

	case TipDogTSU:
		// TODO
		tsu, err := server.store.GetRadTSUById(ctx, dog.ID)
		if err != nil {
			return nil, err
		}

		var detaljT567 []models.DetaljT567

		if tsu.TSU != nil && len(*tsu.TSU) > 0 {
			for _, t := range *tsu.TSU {
				detaljT567 = append(detaljT567, buildDetaljT5(t, tsu))
			}
		}

		copyCommonPDFFields(pdfDog, tsu.Podnaslov, tsu.UzrokTekst, tsu.ManTekst)

		pdfDog.Detalji = detaljT567

	case TipDogTK:
		// TODO
		tk, err := server.store.GetRadTKById(ctx, dog.ID)
		if err != nil {
			return nil, err
		}

		var detaljT567 []models.DetaljT567

		if tk.TK != nil && len(*tk.TK) > 0 {
			for _, t := range *tk.TK {
				detaljT567 = append(detaljT567, buildDetaljT567(t, tk))
			}
		}

		copyCommonPDFFields(pdfDog, tk.Podnaslov, tk.UzrokTekst, tk.ManTekst)

		pdfDog.Detalji = detaljT567

	case TipDogSOP:
		// TODO
		sop, err := server.store.GetRadSOPById(ctx, dog.ID)
		if err != nil {
			return nil, err
		}

		var detaljT567 []models.DetaljT567

		if sop.SOP != nil && len(*sop.SOP) > 0 {
			for _, s := range *sop.SOP {
				detaljT567 = append(detaljT567, buildDetaljT7(s, sop))
			}
		}

		copyCommonPDFFields(pdfDog, sop.Podnaslov, sop.UzrokTekst, sop.ManTekst)

		pdfDog.Detalji = detaljT567

	case TipDogObavestenje:
		// TODO
		switch dog.TipObav {

		case TipDogObavestenjeBeleska:
			// B
			obav, err := server.store.GetObavBeleskaById(ctx, dog.ID)
			if err != nil {
				return nil, err
			}
			// fmt.Printf("Obavestenje za dogadjaj %d: %+v\n", dog.ID, obav.ObavBeleske)

			pdfDog.Podnaslov = obav.Podnaslov

			if obav.ObavBeleske != nil {
				pdfDog.ObavBeleske = obav.ObavBeleske
				pdfDog.ObavBeleske.TipObv = obav.ObavBeleske.TipObv
			}

		case TipDogObavestenjeFaks:
			// F
			slike, err := server.store.GetObavSlikeById(ctx, dog.ID)
			if err != nil {
				return nil, err
			}

			pdfDog.ObavSlike = slike
			/* ostaje za kasnije Obavestenje F -> Podnaslov */
		}

	case TipDogIskljucenje:
		// TODO
		isklj, err := server.store.GetIskljucenjeById(ctx, dog.ID)
		if err != nil {
			return nil, err
		}

		objekti := buildIskljucenjeObjektiPDF(isklj)

		copyCommonPDFFields(pdfDog, isklj.Podnaslov, isklj.UzrokTekst, isklj.ManTekst)

		pdfDog.Objekti = objekti
		pdfDog.Grazlog = isklj.Grazlog
		pdfDog.Razlog = isklj.Razlog

	default:
		// za sada preskačemo nepoznate tipove
		return pdfDog, nil
	}

	return pdfDog, nil
}

func copyCommonPDFFields(pdf *models.DogadjajPDF, podnaslov string, uzrok *string, man *string) {
	pdf.Podnaslov = podnaslov
	if uzrok != nil {
		pdf.UzrokTekst = *uzrok
	}
	if man != nil {
		pdf.ManTekst = *man
	}
}

// func derefInt(i *int) int {
// 	if i == nil {
// 		return 0
// 	}
// 	return *i
// }

func buildIskljucenjeObjektiPDF(isklj *models.DogadjajDetaljno) []models.ObjekatView {

	objektiMap := make(map[string][]models.ManipView)

	if isklj.Manipulacije == nil {
		return nil
	}

	for _, m := range isklj.Manipulacije {

		// osnovna rečenica
		rec := m.Manipulacija

		if m.TekstMan != nil && *m.TekstMan != "" {
			rec += " " + *m.TekstMan
		}
		if m.Ees != nil && *m.Ees != "" {
			rec += " " + *m.Ees
		}
		if m.TekstEes != nil && *m.TekstEes != "" {
			rec += " " + *m.TekstEes
		}
		if m.DvTrafoPolje != nil && *m.DvTrafoPolje != "" && *m.DvTrafoPolje != " " {
			rec += " " + *m.DvTrafoPolje
		}
		if m.Napomena != nil && *m.Napomena != "" {
			rec += "\n" + *m.Napomena
		}

		// dopuna logika
		dopuna := ""

		switch {
		case m.StatusMan == "1":
			dopuna = "Stor."

		case isklj.Dopuna != nil &&
			*isklj.Dopuna == "2":

			if m.DopunaMan != nil &&
				*m.DopunaMan != "1" &&
				isklj.DatumDopune != nil &&
				!isklj.DatumDopune.Equal(isklj.DatumSmene) {

				dopuna = "Dop."
			}
		}

		mv := models.ManipView{
			DopunaDaNe:  dopuna,
			Vrepoc:      m.Vrepoc,
			Vrezav:      derefString(m.Vrezav),
			RecenicaMan: rec,
			Rb:          *m.Rb,
		}

		objektiMap[m.Objekat] = append(objektiMap[m.Objekat], mv)
	}

	// map -> slice + sort
	var objekti []models.ObjekatView

	for naziv, stavke := range objektiMap {

		sort.Slice(stavke, func(i, j int) bool {
			return stavke[i].Rb < stavke[j].Rb
		})

		objekti = append(objekti, models.ObjekatView{
			Naziv:  naziv,
			Stavke: stavke,
		})
	}

	return objekti
}

func buildIspadDetaljiPDF(ispadList []models.Ispad, parent *models.DogadjajDetaljno) []models.DetaljT567 {

	if len(ispadList) == 0 {
		return nil
	}

	var detalji []models.DetaljT567

	for _, ispad := range ispadList {
		detalj := buildDetaljT1(ispad, parent)
		detalji = append(detalji, detalj)
	}

	return detalji
}
