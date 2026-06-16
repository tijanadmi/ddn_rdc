package api

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

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
	for _, m := range dogadjaj.Manipulacije {
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
			Rb:          *m.Rb, // dodajemo Rb za kasnije sortiranje
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

/************** HELPER FUNKCIJE ZA GRADNJU TEKSTA U DETAJLU DOGADJAJA tipa ISPAD *****************/

// func buildRecenica1(isp models.Ispad) string {

// 	var kon string
// 	var vSnaga string
// 	fmt.Println("ispad", isp)

// 	fmt.Println("vrDogSif i SmPk:", isp.VrDogSif, isp.SmPk)
// 	// KON logika
// 	switch {
// 	case isp.VrDogSif == "7" || isp.VrDogSif == "71":
// 		if isp.SmPk != "ostaje u pogonu" {
// 			kon = " je "
// 		}
// 		if isp.SmPk != "" {
// 			kon += isp.SmPk + " "
// 		} else {
// 			kon += "isključen "
// 		}

// 	case isp.TipOb == "1" && isp.VrDogSif != "72":
// 		kon = " ispada "

// 	case isp.TipOb == "7":
// 		kon = " kvar na DVxKV "

// 	case isp.VrDogSif == "72":
// 		kon = " konzum u mraku "
// 	}
// 	fmt.Println("kon:", kon)

// 	if isp.Snaga != nil && *isp.Snaga != "" {
// 		vSnaga = " ispala snaga " + *isp.Snaga
// 	}

// 	// helper values
// 	vrepoc := isp.Vrepoc
// 	vrezav := isp.Vrezav
// 	obj := isp.Objekat
// 	polje := isp.DvTrafoPolje

// 	// CASE 1: VREZAV IS NULL
// 	if vrezav == "" {

// 		if obj == "" {

// 			if isp.VrDogSif == "72" && isp.Snaga != nil {
// 				return "U " + vrepoc + kon + polje + vSnaga
// 			}
// 			return "U " + vrepoc + kon + polje
// 		}

// 		if polje == "" || polje == " " {
// 			if isp.VrDogSif == "72" && isp.Snaga != nil {
// 				return "U " + vrepoc + kon + "u " + obj + vSnaga
// 			}
// 			return "U " + vrepoc + kon + "u " + obj
// 		}

// 		if isp.VrDogSif == "72" && isp.Snaga != nil {
// 			return "U " + vrepoc + kon + polje + " u " + obj + vSnaga
// 		}

// 		return "U " + vrepoc + kon + polje + " u " + obj
// 	}

// 	// CASE 2: VREZAV NOT NULL
// 	if obj == "" {

// 		if isp.VrDogSif == "72" && isp.Snaga != nil {
// 			return "Od " + vrepoc + " do " + vrezav + kon + polje + vSnaga
// 		}
// 		return "Od " + vrepoc + " do " + vrezav + kon + polje
// 	}

// 	if polje == "" || polje == " " {
// 		if isp.VrDogSif == "72" && isp.Snaga != nil {
// 			return "Od " + vrepoc + " do " + vrezav + kon + "u " + obj + vSnaga
// 		}
// 		return "Od " + vrepoc + " do " + vrezav + kon + "u " + obj
// 	}

// 	if isp.VrDogSif == "72" && isp.Snaga != nil {
// 		return "Od " + vrepoc + " do " + vrezav + kon + polje + " u " + obj + vSnaga
// 	}

// 	return "Od " + vrepoc + " do " + vrezav + kon + polje + " u " + obj
// }

func buildRecenica1(isp models.Ispad) string {

	var kon string
	var vSnaga string

	// KON logika (identično PL/SQL)
	switch {
	case isp.VrDogSif == "7" || isp.VrDogSif == "71":
		if isp.SmPk != "ostaje u pogonu" {
			kon = " je "
		}
		if isp.SmPk != "" {
			kon += isp.SmPk + " "
		} else {
			kon += "isključen "
		}

	case isp.TipDog == "1" && isp.VrDogSif != "72":
		kon = " ispada "

	case isp.TipDog == "7":
		kon = " kvar na DVxKV "

	case isp.VrDogSif == "72":
		kon = " konzum u mraku "
	}

	// SNAGA (PL/SQL: v_snaga := ' ispala snaga '||:SNAGA;)
	if isp.Snaga != nil && *isp.Snaga != "" {
		vSnaga = " ispala snaga " + *isp.Snaga
	}

	vrepoc := isp.Vrepoc
	vrezav := isp.Vrezav
	obj := isp.Objekat
	polje := isp.DvTrafoPolje

	// CASE 1: VREZAV IS NULL
	if vrezav == "" {

		if obj == "" {

			if isp.VrDogSif == "72" && isp.Snaga != nil {
				return "U " + vrepoc + kon + polje + vSnaga
			}
			return "U " + vrepoc + kon + polje
		}

		if polje == "" || polje == " " {
			if isp.VrDogSif == "72" && isp.Snaga != nil {
				return "U " + vrepoc + kon + "u " + obj + vSnaga
			}
			return "U " + vrepoc + kon + "u " + obj
		}

		if isp.VrDogSif == "72" && isp.Snaga != nil {
			return "U " + vrepoc + kon + polje + " u " + obj + vSnaga
		}

		return "U " + vrepoc + kon + polje + " u " + obj
	}

	// CASE 2: VREZAV NOT NULL
	if obj == "" {

		if isp.VrDogSif == "72" && isp.Snaga != nil {
			return "Od " + vrepoc + " do " + vrezav + kon + polje + vSnaga
		}
		return "Od " + vrepoc + " do " + vrezav + kon + polje
	}

	if polje == "" || polje == " " {
		if isp.VrDogSif == "72" && isp.Snaga != nil {
			return "Od " + vrepoc + " do " + vrezav + kon + "u " + obj + vSnaga
		}
		return "Od " + vrepoc + " do " + vrezav + kon + "u " + obj
	}

	if isp.VrDogSif == "72" && isp.Snaga != nil {
		return "Od " + vrepoc + " do " + vrezav + kon + polje + " u " + obj + vSnaga
	}

	return "Od " + vrepoc + " do " + vrezav + kon + polje + " u " + obj
}

func buildRecenica2(isp models.Ispad) string {

	if isp.VrDogSif == "72" {
		return " "
	}

	rec := "Vrsta događaja: " + isp.VrstaDog + ";  Uzrok:  " + isp.Uzrok1

	if isp.RadApu != "" && isp.RadApu != "Nepotreban" {
		rec += "\nRad APU-a: " + isp.RadApu
	}

	if isp.VremUsl != "" {
		if isp.RadApu != "" && isp.RadApu != "Nepotreban" {
			rec += ";   Vremenski uslovi: " + isp.VremUsl
		} else {
			rec += "\nVremenski uslovi: " + isp.VremUsl
		}
	}

	return rec
}

func val(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func hasAny(values ...*string) bool {
	for _, v := range values {
		if v != nil && *v != "" {
			return true
		}
	}
	return false
}

func clean(s *string) string {
	if s == nil {
		return ""
	}
	return strings.TrimSpace(*s)
}

func isNotEmpty(s *string) bool {
	return s != nil && strings.TrimSpace(*s) != ""
}

func buildZastitaOpis(isp models.Ispad) string {
	var glavna1, glavna2, glavna3 string
	var rezervna, rezervna2 string
	var prekidac, neelektricna string
	var sab1, sab2 string
	var jps, jps2 string
	var recenica string

	// fmt.Println("ispad buildZastitaOpis", val(isp.ZZMSPGl1))

	// ---------------- GLAVNA 1 ----------------
	if hasAny(isp.ZDsdfGl1, isp.ZKvarGl1, isp.ZPrstGl1, isp.ZZMSPGl1,
		isp.ZUzmsGl1, isp.ZRapuGl1, isp.IdZTelePocGl1, isp.IdZTeleKrajGl1, isp.ZLokkGl1) {

		glavna1 = "Glavna zaštita 1\n"

		if v := clean(isp.ZDsdfGl1); v != "" {
			glavna1 += v + "   "
		}

		if v := clean(isp.ZKvarGl1); v != "" {
			glavna1 += "Kvar u:" + v + "   "
		}

		if v := clean(isp.ZPrstGl1); v != "" {
			glavna1 += v + "   "
		}

		if v := clean(isp.ZZMSPGl1); v != "" {
			glavna1 += v + "   "
		}

		if v := clean(isp.ZUzmsGl1); v != "" {
			glavna1 += v + "   "
		}

		if v := clean(isp.ZRapuGl1); v != "" {
			glavna1 += "APU:" + v + "   "
		}

		if v := clean(isp.IdZTelePocGl1); v != "" {
			glavna1 += v + "   "
		}

		if v := clean(isp.IdZTeleKrajGl1); v != "" {
			glavna1 += v + "   "
		}

		if v := clean(isp.ZLokkGl1); v != "" {
			glavna1 += "Lokator kvara:" + v + "km "
		}
	}

	// ---------------- GLAVNA 2 ----------------

	// ---------------- PROVERA DA LI ULAZI ----------------
	if hasAny(
		isp.ZDsdfGl2, isp.ZKvarGl2, isp.ZPrstGl2, isp.ZZMSPGl2,
		isp.ZUzmsGl2, isp.ZRapuGl2, isp.ZLokkGl2,
		isp.IdZTelePocGl2, isp.IdZTeleKrajGl2,
	) {

		// ---------------- HEADER (PL/SQL TIP_OB / NAPON LOGIKA) ----------------
		if (isp.TipOb == "DV" || isp.TipOb == "TD" || isp.TipOb == "KB" || isp.TipOb == "TK") &&
			isp.Napon != "400" {
			glavna2 = "Glavna zaštita sa funkcijom jedinice polja\n"
		} else {
			glavna2 = "Glavna zaštita 2\n"
		}

		// ---------------- DSDF ----------------
		if isNotEmpty(isp.ZDsdfGl2) {
			glavna2 += val(isp.ZDsdfGl2) + "   "
		}

		// ---------------- KVAR ----------------
		if isNotEmpty(isp.ZKvarGl2) {
			glavna2 += "Kvar u:" + val(isp.ZKvarGl2) + "   "
		}

		// ---------------- PRST ----------------
		if isNotEmpty(isp.ZPrstGl2) {
			glavna2 += val(isp.ZPrstGl2) + "   "
		}

		// ---------------- ZMSP ----------------
		if isNotEmpty(isp.ZZMSPGl2) {
			glavna2 += val(isp.ZZMSPGl2) + "   "
		}

		// ---------------- UZMS ----------------
		if isNotEmpty(isp.ZUzmsGl2) {
			glavna2 += val(isp.ZUzmsGl2) + "   "
		}

		// ---------------- RAPU ----------------
		if isNotEmpty(isp.ZRapuGl2) {
			glavna2 += val(isp.ZRapuGl2) + "   "
		}

		// ---------------- TELE POC ----------------
		if isNotEmpty(isp.IdZTelePocGl2) {
			glavna2 += val(isp.IdZTelePocGl2) + "   "
		}

		// ---------------- TELE KRAJ ----------------
		if isNotEmpty(isp.IdZTeleKrajGl2) {
			glavna2 += val(isp.IdZTeleKrajGl2) + "   "
		}

		// ---------------- LOKATOR ----------------
		if isNotEmpty(isp.ZLokkGl2) {
			glavna2 += "Lokator kvara:" + val(isp.ZLokkGl2) + "km "
		}
	}
	// ---------------- GLAVNA 3 ----------------
	if hasAny(
		isp.ZDsdfGl3, isp.ZKvarGl3, isp.ZPrstGl3, isp.ZZMSPGl3,
		isp.ZUzmsGl3, isp.ZRapuGl3, isp.ZLokkGl3,
		isp.IdZTelePocGl3, isp.IdZTeleKrajGl3,
	) {

		glavna3 = "Glavna zaštita sa funkcijom jedinice polja\n"

		if isNotEmpty(isp.ZDsdfGl3) {
			glavna3 += val(isp.ZDsdfGl3) + "   "
		}
		if isNotEmpty(isp.ZKvarGl3) {
			glavna3 += "Kvar u:" + val(isp.ZKvarGl3) + "   "
		}
		if isNotEmpty(isp.ZPrstGl3) {
			glavna3 += val(isp.ZPrstGl3) + "   "
		}
		if isNotEmpty(isp.ZZMSPGl3) {
			glavna3 += val(isp.ZZMSPGl3) + "   "
		}
		if isNotEmpty(isp.ZUzmsGl3) {
			glavna3 += val(isp.ZUzmsGl3) + "   "
		}
		if isNotEmpty(isp.ZRapuGl3) {
			glavna3 += val(isp.ZRapuGl3) + "   "
		}
		if isNotEmpty(isp.IdZTelePocGl3) {
			glavna3 += val(isp.IdZTelePocGl3) + "   "
		}
		if isNotEmpty(isp.IdZTeleKrajGl3) {
			glavna3 += val(isp.IdZTeleKrajGl3) + "   "
		}
		if isNotEmpty(isp.ZLokkGl3) {
			glavna3 += "Lokator kvara:" + val(isp.ZLokkGl3) + "km "
		}
	}

	// ---------------- REZERVNA ----------------
	if hasAny(isp.ZDisRez, isp.ZKvarRez, isp.ZPrstRez, isp.ZZMSPRez) {

		rezervna = "Dopunska (rezervna) zaštita 1\n"

		if isNotEmpty(isp.ZDisRez) {
			rezervna += val(isp.ZDisRez) + "   "
		}

		if isNotEmpty(isp.ZKvarRez) {
			rezervna += "Kvar u:" + val(isp.ZKvarRez) + "   "
		}

		if isNotEmpty(isp.ZPrstRez) {
			rezervna += val(isp.ZPrstRez) + "   "
		}

		if isNotEmpty(isp.ZZMSPRez) {
			rezervna += val(isp.ZZMSPRez) + "   "
		}
	}
	// ---------------- REZERVNA 2 ----------------
	if hasAny(isp.ZDisRez2, isp.ZKvarRez2, isp.ZPrstRez2, isp.ZZMSPRez2) {

		rezervna2 = "Dopunska (rezervna) zaštita 2\n"

		if isNotEmpty(isp.ZDisRez2) {
			rezervna2 += val(isp.ZDisRez2) + "   "
		}

		if isNotEmpty(isp.ZKvarRez2) {
			rezervna2 += "Kvar u:" + val(isp.ZKvarRez2) + "   "
		}

		if isNotEmpty(isp.ZPrstRez2) {
			rezervna2 += val(isp.ZPrstRez2) + "   "
		}

		if isNotEmpty(isp.ZZMSPRez2) {
			rezervna2 += val(isp.ZZMSPRez2) + "   "
		}
	}

	// ---------------- PREKIDAC ----------------

	if hasAny(isp.ZPrekVn, isp.ZPrekNn) {

		if isNotEmpty(isp.ZPrekVn) {
			if isp.Fup == "01" || isp.Fup == "18" || isp.Fup == "22" {
				prekidac += "Prekidač VN strana:" + strings.TrimSpace(*isp.ZPrekVn) + "   "
			} else {
				prekidac += strings.TrimSpace(*isp.ZPrekVn) + "   "
			}
		}

		if isNotEmpty(isp.ZPrekNn) {
			prekidac += "Prekidač SN strana:" + strings.TrimSpace(*isp.ZPrekNn) + "   "
		}
	}

	// ---------------- NEELEKTRICNA ----------------
	if hasAny(isp.ZNel1, isp.ZNel2, isp.ZNel3) {

		if isNotEmpty(isp.ZNel1) {
			neelektricna += strings.TrimSpace(*isp.ZNel1) + "   "
		}
		if isNotEmpty(isp.ZNel2) {
			neelektricna += strings.TrimSpace(*isp.ZNel2) + "   "
		}
		if isNotEmpty(isp.ZNel3) {
			neelektricna += strings.TrimSpace(*isp.ZNel3) + "   "
		}
	}

	// ---------------- SABIRNICE ----------------
	if hasAny(isp.ZSabzSab, isp.ZOtprSab, isp.ZKvarGl1) {
		sab1 = "Sabirnička zaštita 1\n"

		if isNotEmpty(isp.ZSabzSab) {
			sab1 += strings.TrimSpace(*isp.ZSabzSab) + "   "
		}

		if isNotEmpty(isp.ZOtprSab) {
			sab1 += "Zas. od otkaza prek.:" + strings.TrimSpace(*isp.ZOtprSab) + "   "
		}

		if isNotEmpty(isp.ZKvarGl1) {
			sab1 += "Kvar u:" + strings.TrimSpace(*isp.ZKvarGl1) + "   "
		}
	}

	if hasAny(isp.ZSabzSab2, isp.ZOtprSab2, isp.ZKvarGl2) {
		sab2 = "Sabirnička zaštita 2\n"

		if isNotEmpty(isp.ZSabzSab2) {
			sab2 += strings.TrimSpace(*isp.ZSabzSab2) + "   "
		}

		if isNotEmpty(isp.ZOtprSab2) {
			sab2 += "Zas. od otkaza prek.:" + strings.TrimSpace(*isp.ZOtprSab2) + "   "
		}

		if isNotEmpty(isp.ZKvarGl2) {
			sab2 += "Kvar u:" + strings.TrimSpace(*isp.ZKvarGl2) + "   "
		}
	}

	// ---------------- SKLAPANJE ----------------
	add := func(s string) {
		if s == "" {
			return
		}
		if recenica != "" {
			recenica += "\n"
		}
		recenica += s
	}

	add(glavna1)
	add(glavna2)
	add(glavna3)
	add(rezervna)
	add(rezervna2)
	add(prekidac)
	add(neelektricna)

	// ⚠️ specijalna logika za sabirnice
	if sab1 != "" && isp.Fup == "00" {
		recenica = sab1
	}
	if sab2 != "" && isp.Fup == "00" {
		if sab1 != "" {
			recenica += "\n" + sab2
		} else {
			recenica = sab2
		}
	}

	// ---------------- JPS ----------------
	if hasAny(isp.ZJpsVn, isp.ZJpsNn) {
		if isp.ZJpsVn != nil {
			if isp.Fup == "01" || isp.Fup == "18" || isp.Fup == "22" {
				jps += "Jedinica polja sabirničke zaštite 1 VN strana:" + val(isp.ZJpsVn) + "   "
			} else {
				jps += "Jedinica polja sabirničke zaštite 1:" + val(isp.ZJpsVn) + "   "
			}
		}
		if isp.ZJpsNn != nil {
			jps += "Jedinica polja sabirničke zaštite 1 SN strana:" + val(isp.ZJpsNn) + "   "
		}
		add(jps)
	}

	if hasAny(isp.ZJpsVn2, isp.ZJpsNn2) {
		if isp.ZJpsVn2 != nil {
			if isp.Fup == "01" || isp.Fup == "18" || isp.Fup == "22" {
				jps2 += "Jedinica polja sabirničke zaštite 2 VN strana:" + val(isp.ZJpsVn2) + "   "
			} else {
				jps2 += "Jedinica polja sabirničke zaštite 2:" + val(isp.ZJpsVn2) + "   "
			}
		}
		if isp.ZJpsNn2 != nil {
			jps2 += "Jedinica polja sabirničke zaštite 2 SN strana:" + val(isp.ZJpsNn2) + "   "
		}
		add(jps2)
	}

	return recenica
}

func buildDetaljT1(isp models.Ispad, dog *models.DogadjajDetaljno) models.DetaljT567 {

	rec1 := buildRecenica1(isp)
	rec2 := buildRecenica2(isp)
	zastita := buildZastitaOpis(isp)

	var opis string

	if zastita == "" {
		opis = isp.Opis
	} else {
		opis = "RAD ZAŠTITE:\n" + zastita

		if isp.Opis != "" {
			opis += "\n" + isp.Opis
		}
	}

	dopuna := ""
	if isp.StatusIspkv1 == "1" {
		dopuna = "Stor."
	} else if dog.Dopuna != nil && *dog.Dopuna == "2" {
		if isp.DopunaIspkv1 != nil && *isp.DopunaIspkv1 != "1" &&
			dog.DatumDopune != nil &&
			!dog.DatumDopune.Equal(dog.DatumSmene) {

			dopuna = "Dop."
		}
	}

	return models.DetaljT567{
		DopunaDaNe: dopuna,
		Recenica1:  rec1,
		Recenica2:  rec2,
		Opis:       opis,
	}
}

func (server *Server) getIspad(ctx *gin.Context) {
	// 1. Dohvati ID iz URL parametra
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Nevažeći ID"})
		return
	}

	// 2. Pozovi funkciju iz store-a
	dogadjaj, err := server.store.GetIspadById(ctx, id)
	if err != nil {
		fmt.Printf("Greška prilikom dobijanja događaja: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	/*** deo koji se odnosi na ispade ***/

	var detaljT567 []models.DetaljT567

	if dogadjaj.Ispad != nil && len(*dogadjaj.Ispad) > 0 {
		for _, ispad := range *dogadjaj.Ispad {
			detalj := buildDetaljT1(ispad, dogadjaj)
			detaljT567 = append(detaljT567, detalj)
		}
	}

	/**** kraj dela koji se odnosi na ispade ***/

	// 4. Grupisanje manipulacija po objektu i sortiranje po Rb
	objektiMap := make(map[string][]ManipView)
	if dogadjaj.Manipulacije != nil {
		for _, m := range dogadjaj.Manipulacije {
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
				Rb:          *m.Rb, // dodajemo Rb za kasnije sortiranje
			}

			objektiMap[m.Objekat] = append(objektiMap[m.Objekat], mv)
		}
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
		"uzrok_tekst": dogadjaj.UzrokTekst,
		"man_tekst":   dogadjaj.ManTekst,
		"posledice":   dogadjaj.Posledice,
		"detaljT567":  detaljT567,
		"objekti":     objekti,
	})
}

/************ Prekid proizvodnje ***********/

func buildDetaljPrekidP(tk models.PrekidP, d *models.DogadjajDetaljno) models.DetaljT567 {
	var rec1 string

	// ---------------- VREME ----------------
	if tk.Vrezav == nil || *tk.Vrezav == "" {
		rec1 = "U " + tk.Vrepoc
	} else {
		rec1 = "Od " + tk.Vrepoc + " do " + *tk.Vrezav
	}

	// ---------------- NIL ZAŠTITA ----------------
	generator := ""
	if tk.Generator != nil {
		generator = *tk.Generator
	}

	vrPrek := ""
	if tk.VrPrek != nil {
		vrPrek = *tk.VrPrek
	}

	objekat := ""
	if tk.Objekat != nil {
		objekat = *tk.Objekat
	}

	snaga := ""
	if tk.Snaga != nil {
		snaga = *tk.Snaga
	}

	// ---------------- RECENICA 1 ----------------
	rec1 += " " + vrPrek +
		" prekid proizvodnje u " + objekat +
		" , " + generator +
		" ispala snaga " + snaga + " MW."

	// ---------------- RECENICA 2 ----------------
	rec2 := ""
	if tk.TipPrek != nil && *tk.TipPrek != "" {
		rec2 = "Tip prekida: " + *tk.TipPrek + "."
	}

	// ---------------- OPIS ----------------
	opis := ""
	if tk.Opis != nil {
		opis = *tk.Opis
	}

	// ---------------- DOPUNA ----------------
	dopuna := buildDopuna(tk.Status, tk.Dopuna, d)

	return models.DetaljT567{
		DopunaDaNe: dopuna,
		Recenica1:  rec1,
		Recenica2:  rec2,
		Opis:       opis,
	}
}

func (server *Server) getPrekidP(ctx *gin.Context) {
	// 1. Dohvati ID iz URL parametra
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Nevažeći ID"})
		return
	}

	// 2. Pozovi funkciju iz store-a
	dogadjaj, err := server.store.GetPrekidPById(ctx, id)
	if err != nil {
		fmt.Printf("Greška prilikom dobijanja događaja: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var detaljT567 []models.DetaljT567

	if dogadjaj.PrekidP != nil && len(*dogadjaj.PrekidP) > 0 {
		for _, tk := range *dogadjaj.PrekidP {
			detalj := buildDetaljPrekidP(tk, dogadjaj)
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
		"posledice":   dogadjaj.Posledice,
		"detaljT567":  detaljT567,
	})
}

func (server *Server) getObavSlike(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Nevažeći ID"})
		return
	}

	images, err := server.store.GetObavSlikeById(ctx, id)
	if err != nil {
		fmt.Printf("Greška prilikom dobijanja slika: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":     id,
		"images": images,
	})
}

func (server *Server) getAngazovaniRukovaoci(ctx *gin.Context) {
	// 1. Dohvati ID iz URL parametra
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Nevažeći ID"})
		return
	}

	// 2. Pozovi funkciju iz store-a
	dogadjaj, err := server.store.GetAngazovaniRukovaociById(ctx, id)
	if err != nil {
		fmt.Printf("Greška prilikom dobijanja događaja: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if dogadjaj == nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Događaj nije pronađen",
		})
		return
	}

	// 6. Kreiraj finalni JSON za frontend
	ctx.JSON(http.StatusOK, gin.H{
		"rb_dog":               dogadjaj.RbDog,
		"naslov":               dogadjaj.Naslov,
		"podnaslov":            dogadjaj.Podnaslov,
		"angazovani_rukovaoci": dogadjaj.AngazovaniRukovalac,
	})
}
