package api

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tijanadmi/ddn_rdc/models"
)

func (server *Server) listOOpenShifts(ctx *gin.Context) {

	// fmt.Println("usao u handler listOpenShifts")

	// 1. Pozovi funkciju iz store-a da dobiješ otvorene smene
	smene, err := server.store.GetOtvoreneSmene(ctx)
	if err != nil {
		fmt.Printf("Greška prilikom dobijanja otvorenih smena: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// fmt.Printf("Dobijene smene: %+v\n", smene)

	// 2. Vrati rezultat
	ctx.JSON(http.StatusOK, gin.H{
		"data": smene,
	})
}

type listClosedShiftsByPageRequest struct {
	Mrc       string `form:"mrc" binding:"required"`
	StartDate string `form:"start_date" binding:"required"`
	EndDate   string `form:"end_date" binding:"required"`
	PageID    int32  `form:"page_id" binding:"required,min=1"`
	PageSize  int32  `form:"page_size" binding:"required,min=5,max=100"`
}

type listClosedShiftsByPageResponse struct {
	Total int            `json:"total"`
	Smene []models.Smena `json:"smene"`
}

func (server *Server) listClosedShiftsByPage(ctx *gin.Context) {

	// fmt.Println("usao u handler listOpenShifts")

	var req listClosedShiftsByPageRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := models.ListShiftsWithPaginationParams{
		Mrc:       req.Mrc,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Limit:     req.PageSize,
		Offset:    (req.PageID - 1) * req.PageSize,
	}

	// 1. Pozovi funkciju iz store-a da dobiješ otvorene smene
	smene, count, err := server.store.GetZatvoreneSmene(ctx, arg)
	if err != nil {
		fmt.Printf("Greška prilikom dobijanja otvorenih smena: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// fmt.Printf("Dobijene smene: %+v\n", smene)
	rsp := listClosedShiftsByPageResponse{
		Total: count,
		Smene: smene,
	}

	ctx.JSON(http.StatusOK, rsp)
}

type ManipView struct {
	DopunaDaNe  string `json:"dopuna_da_ne"`
	Vrepoc      string `json:"vrepoc"`
	Vrezav      string `json:"vrezav"`
	RecenicaMan string `json:"recenica_man"`
	Rb          int    `json:"-"`
}

type ObjekatView struct {
	Naziv  string      `json:"naziv"`
	Stavke []ManipView `json:"stavke"`
}

func (server *Server) getIskljucenje(ctx *gin.Context) {
	// 1. Dohvati ID iz URL parametra
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Nevažeći ID"})
		return
	}

	// 2. Pozovi funkciju iz store-a
	dogadjaj, err := server.store.GetIskljucenjeById(ctx, id)
	if err != nil {
		fmt.Printf("Greška prilikom dobijanja događaja: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 4. Grupisanje manipulacija po objektu i sortiranje po Rb
	objektiMap := make(map[string][]ManipView)
	for _, m := range *dogadjaj.Manipulacije {
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
		if m.DvTrafoPolje != nil && *m.DvTrafoPolje != " " && *m.DvTrafoPolje != "" {
			rec += " " + *m.DvTrafoPolje
		}
		if m.Napomena != nil && *m.Napomena != "" {
			rec += "\n" + *m.Napomena
		}

		dopuna := ""

		switch {
		case m.StatusMan == "1":
			dopuna = "Stor."

		case dogadjaj.Dopuna != nil && *dogadjaj.Dopuna == "2":
			if m.DopunaMan != nil && *m.DopunaMan != "1" &&
				dogadjaj.DatumDopune != nil &&
				!dogadjaj.DatumDopune.Equal(dogadjaj.DatumSmene) {

				dopuna = "Dop."
			}
		}

		mv := ManipView{
			DopunaDaNe:  dopuna,
			Vrepoc:      m.Vrepoc,
			Vrezav:      derefString(m.Vrezav),
			RecenicaMan: rec,
			Rb:          m.Rb, // dodajemo Rb za kasnije sortiranje
		}

		objektiMap[m.Objekat] = append(objektiMap[m.Objekat], mv)
	}

	// 5. Transform map u slice i sortiraj po Rb
	var objekti []ObjekatView
	for naziv, stavke := range objektiMap {
		sort.Slice(stavke, func(i, j int) bool {
			return stavke[i].Rb < stavke[j].Rb
		})
		objekti = append(objekti, ObjekatView{
			Naziv:  naziv,
			Stavke: stavke,
		})
	}

	// 6. Kreiraj finalni JSON za frontend
	ctx.JSON(http.StatusOK, gin.H{
		"rb_dog":      dogadjaj.RbDog,
		"naslov":      dogadjaj.Naslov,
		"podnaslov":   dogadjaj.Podnaslov,
		"grazlog":     dogadjaj.Grazlog,
		"razlog":      dogadjaj.Razlog,
		"uzrok_tekst": dogadjaj.UzrokTekst,
		"man_tekst":   dogadjaj.ManTekst,
		"objekti":     objekti,
	})
}

// helper funkcija
func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func (server *Server) getObavBeleska(ctx *gin.Context) {
	//  Dohvati ID iz URL parametra
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Nevažeći ID"})
		return
	}

	//  Pozovi funkciju iz store-a
	dogadjaj, err := server.store.GetObavBeleskaById(ctx, id)
	if err != nil {
		fmt.Printf("Greška prilikom dobijanja događaja: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var obav models.ObavBeleska

	if dogadjaj.ObavBeleske != nil {
		obav = *dogadjaj.ObavBeleske
	} else {
		obav = models.ObavBeleska{}
	}

	//  Kreiraj finalni JSON za frontend
	ctx.JSON(http.StatusOK, gin.H{
		"rb_dog":       dogadjaj.RbDog,
		"naslov":       dogadjaj.Naslov,
		"podnaslov":    dogadjaj.Podnaslov,
		"obav_beleska": obav,
	})
}

// func buildDopuna(tk models.TK, d *models.DogadjajDetaljno) string {
// 	switch {
// 	case tk.Status != nil && *tk.Status == "1":
// 		return "Stor."

// 	case d.Dopuna != nil && *d.Dopuna == "2":
// 		if tk.Dopuna == nil || *tk.Dopuna != "1" {
// 			if d.DatumDopune != nil && !d.DatumDopune.Equal(d.DatumSmene) {
// 				return "Dop."
// 			}
// 		}
// 	}
// 	return ""
// }

func buildDopuna(status *string, dopunaStavke *string, d *models.DogadjajDetaljno) string {
	switch {
	// Stor.
	case status != nil && *status == "1":
		return "Stor."

	// Dop.
	case d.Dopuna != nil && *d.Dopuna == "2":
		if dopunaStavke == nil || *dopunaStavke != "1" {
			if d.DatumDopune != nil && !d.DatumDopune.Equal(d.DatumSmene) {
				return "Dop."
			}
		}
	}

	return ""
}

func buildDetaljT567(tk models.TK, d *models.DogadjajDetaljno) models.DetaljT567 {
	var rec1 string

	// vreme
	if tk.Vrezav == nil || *tk.Vrezav == "" {
		rec1 = "U " + tk.Vrepoc
	} else {
		rec1 = "Od " + tk.Vrepoc + " do " + *tk.Vrezav
	}

	// objekti (ispravljen edge case!)
	if tk.ObID2 == nil || tk.ObjekatNaziv2 == nil || *tk.ObjekatNaziv2 == "" {
		rec1 += " u " + tk.ObjekatNaziv
	} else {
		rec1 += " od " + tk.ObjekatNaziv + " do " + *tk.ObjekatNaziv2
	}

	rec1 += "       Vrsta događaja:   " + tk.VrstaDog

	rec2 := "Vrsta TK opreme:  " + tk.Vropr

	opis := ""
	if tk.Opis != nil {
		opis = *tk.Opis
	}

	dopuna := buildDopuna(tk.Status, tk.Dopuna, d)

	return models.DetaljT567{
		DopunaDaNe: dopuna,
		Recenica1:  rec1,
		Recenica2:  rec2,
		Opis:       opis,
	}
}

func (server *Server) getRadTK(ctx *gin.Context) {
	// 1. Dohvati ID iz URL parametra
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Nevažeći ID"})
		return
	}

	// 2. Pozovi funkciju iz store-a
	dogadjaj, err := server.store.GetRadTKById(ctx, id)
	if err != nil {
		fmt.Printf("Greška prilikom dobijanja događaja: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var detaljT567 []models.DetaljT567

	if dogadjaj.TK != nil && len(*dogadjaj.TK) > 0 {
		for _, tk := range *dogadjaj.TK {
			detalj := buildDetaljT567(tk, dogadjaj)
			detaljT567 = append(detaljT567, detalj)
		}
	}

	// 6. Kreiraj finalni JSON za frontend
	ctx.JSON(http.StatusOK, gin.H{
		"rb_dog":      dogadjaj.RbDog,
		"naslov":      dogadjaj.Naslov,
		"podnaslov":   dogadjaj.Podnaslov,
		"uzrok_tekst": dogadjaj.UzrokTekst,
		"man_tekst":   dogadjaj.ManTekst,
		"detaljT567":  detaljT567,
	})
}

func buildDetaljT5(tsu models.TSU, d *models.DogadjajDetaljno) models.DetaljT567 {
	var rec1 string

	// vreme
	if tsu.Vrezav == nil || *tsu.Vrezav == "" {
		rec1 = "U " + tsu.Vrepoc
	} else {
		rec1 = "Od " + tsu.Vrepoc + " do " + *tsu.Vrezav
	}

	rec1 += " u " + tsu.ObjekatNaziv

	rec1 += "       Vrsta događaja:   " + tsu.VrstaDog

	rec2 := "Vrsta TSU opreme:  " + tsu.Vropr

	opis := ""
	if tsu.Opis != nil {
		opis = *tsu.Opis
	}

	dopuna := buildDopuna(tsu.Status, tsu.Dopuna, d)

	return models.DetaljT567{
		DopunaDaNe: dopuna,
		Recenica1:  rec1,
		Recenica2:  rec2,
		Opis:       opis,
	}
}

func (server *Server) getRadTSU(ctx *gin.Context) {
	// 1. Dohvati ID iz URL parametra
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Nevažeći ID"})
		return
	}

	// 2. Pozovi funkciju iz store-a
	dogadjaj, err := server.store.GetRadTSUById(ctx, id)
	if err != nil {
		fmt.Printf("Greška prilikom dobijanja događaja: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var detaljT567 []models.DetaljT567

	if dogadjaj.TSU != nil && len(*dogadjaj.TSU) > 0 {
		for _, tsu := range *dogadjaj.TSU {
			detalj := buildDetaljT5(tsu, dogadjaj)
			detaljT567 = append(detaljT567, detalj)
		}
	}

	// 6. Kreiraj finalni JSON za frontend
	ctx.JSON(http.StatusOK, gin.H{
		"rb_dog":      dogadjaj.RbDog,
		"naslov":      dogadjaj.Naslov,
		"podnaslov":   dogadjaj.Podnaslov,
		"uzrok_tekst": dogadjaj.UzrokTekst,
		"man_tekst":   dogadjaj.ManTekst,
		"detaljT567":  detaljT567,
	})
}

func buildDetaljT7(sop models.SOP, d *models.DogadjajDetaljno) models.DetaljT567 {
	var rec1 string

	// vreme
	if sop.Vrezav == nil || *sop.Vrezav == "" {
		rec1 = "U " + sop.Vrepoc
	} else {
		rec1 = "Od " + sop.Vrepoc + " do " + *sop.Vrezav
	}

	rec1 += " u " + sop.ObjekatNaziv

	rec1 += "       Vrsta događaja:   " + sop.VrstaDog

	rec2 := "Vrsta uređaja za sopstvenu potrošnju:  " + sop.NazSop
	if sop.RBrSop != "" {
		rec2 += "    Broj:  " + sop.RBrSop
	}

	opis := ""
	if sop.Opis != nil {
		opis = *sop.Opis
	}

	dopuna := buildDopuna(sop.Status, sop.Dopuna, d)

	return models.DetaljT567{
		DopunaDaNe: dopuna,
		Recenica1:  rec1,
		Recenica2:  rec2,
		Opis:       opis,
	}
}

func (server *Server) getRadSOP(ctx *gin.Context) {
	// 1. Dohvati ID iz URL parametra
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Nevažeći ID"})
		return
	}

	// 2. Pozovi funkciju iz store-a
	dogadjaj, err := server.store.GetRadSOPById(ctx, id)
	if err != nil {
		fmt.Printf("Greška prilikom dobijanja događaja: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var detaljT567 []models.DetaljT567

	if dogadjaj.SOP != nil && len(*dogadjaj.SOP) > 0 {
		for _, sop := range *dogadjaj.SOP {
			detalj := buildDetaljT7(sop, dogadjaj)
			detaljT567 = append(detaljT567, detalj)
		}
	}

	// 6. Kreiraj finalni JSON za frontend
	ctx.JSON(http.StatusOK, gin.H{
		"rb_dog":      dogadjaj.RbDog,
		"naslov":      dogadjaj.Naslov,
		"podnaslov":   dogadjaj.Podnaslov,
		"uzrok_tekst": dogadjaj.UzrokTekst,
		"man_tekst":   dogadjaj.ManTekst,
		"detaljT567":  detaljT567,
	})
}
