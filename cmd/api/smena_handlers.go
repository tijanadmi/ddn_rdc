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

func buildRecenica1(isp models.Ispad) string {

	var kon string
	var vSnaga string

	// KON logika
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

	case isp.TipOb == "1" && isp.VrDogSif != "72":
		kon = " ispada "

	case isp.TipOb == "7":
		kon = " kvar na DVxKV "

	case isp.VrDogSif == "72":
		kon = " konzum u mraku "
	}

	if isp.Snaga != nil && *isp.Snaga != "" {
		vSnaga = " ispala snaga " + *isp.Snaga
	}

	// helper values
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

func buildZastitaOpis(isp models.Ispad) string {
	var glavna1, glavna2, glavna3 string
	var rezervna, rezervna2 string
	var prekidac, neelektricna string
	var sab1, sab2 string
	var jps, jps2 string
	var recenica string

	// ---------------- GLAVNA 1 ----------------
	if hasAny(isp.ZDsdfGl1, isp.ZKvarGl1, isp.ZPrstGl1, isp.ZZMSPGl1,
		isp.ZUzmsGl1, isp.ZRapuGl1, isp.ZLokkGl1,
		isp.IdZTelePocGl1, isp.IdZTeleKrajGl1) {

		glavna1 = "Glavna zaštita 1\n"
		if isp.ZDsdfGl1 != nil {
			glavna1 += val(isp.ZDsdfGl1) + "   "
		}
		if isp.ZKvarGl1 != nil {
			glavna1 += "Kvar u:" + val(isp.ZKvarGl1) + "   "
		}
		if isp.ZPrstGl1 != nil {
			glavna1 += val(isp.ZPrstGl1) + "   "
		}
		if isp.ZZMSPGl1 != nil {
			glavna1 += val(isp.ZZMSPGl1) + "   "
		}
		if isp.ZUzmsGl1 != nil {
			glavna1 += val(isp.ZUzmsGl1) + "   "
		}
		if isp.ZRapuGl1 != nil {
			glavna1 += "APU:" + val(isp.ZRapuGl1) + "   "
		}
		if isp.IdZTelePocGl1 != nil {
			glavna1 += val(isp.IdZTelePocGl1) + "   "
		}
		if isp.IdZTeleKrajGl1 != nil {
			glavna1 += val(isp.IdZTeleKrajGl1) + "   "
		}
		if isp.ZLokkGl1 != nil {
			glavna1 += "Lokator kvara:" + val(isp.ZLokkGl1) + "km "
		}
	}

	// ---------------- GLAVNA 2 ----------------
	if hasAny(isp.ZDsdfGl2, isp.ZKvarGl2, isp.ZPrstGl2, isp.ZZMSPGl2,
		isp.ZUzmsGl2, isp.ZRapuGl2, isp.ZLokkGl2,
		isp.IdZTelePocGl2, isp.IdZTeleKrajGl2) {

		if (isp.TipOb == "DV" || isp.TipOb == "TD" || isp.TipOb == "KB" || isp.TipOb == "TK") && isp.Napon != "400" {
			glavna2 = "Glavna zaštita sa funkcijom jedinice polja\n"
		} else {
			glavna2 = "Glavna zaštita 2\n"
		}

		if isp.ZDsdfGl2 != nil {
			glavna2 += val(isp.ZDsdfGl2) + "   "
		}
		if isp.ZKvarGl2 != nil {
			glavna2 += "Kvar u:" + val(isp.ZKvarGl2) + "   "
		}
		if isp.ZPrstGl2 != nil {
			glavna2 += val(isp.ZPrstGl2) + "   "
		}
		if isp.ZZMSPGl2 != nil {
			glavna2 += val(isp.ZZMSPGl2) + "   "
		}
		if isp.ZUzmsGl2 != nil {
			glavna2 += val(isp.ZUzmsGl2) + "   "
		}
		if isp.ZRapuGl2 != nil {
			glavna2 += val(isp.ZRapuGl2) + "   "
		}
		if isp.IdZTelePocGl2 != nil {
			glavna2 += val(isp.IdZTelePocGl2) + "   "
		}
		if isp.IdZTeleKrajGl2 != nil {
			glavna2 += val(isp.IdZTeleKrajGl2) + "   "
		}
		if isp.ZLokkGl2 != nil {
			glavna2 += "Lokator kvara:" + val(isp.ZLokkGl2) + "km "
		}
	}

	// ---------------- GLAVNA 3 ----------------
	if hasAny(isp.ZDsdfGl3, isp.ZKvarGl3, isp.ZPrstGl3, isp.ZZMSPGl3,
		isp.ZUzmsGl3, isp.ZRapuGl3, isp.ZLokkGl3,
		isp.IdZTelePocGl3, isp.IdZTeleKrajGl3) {

		glavna3 = "Glavna zaštita sa funkcijom jedinice polja\n"

		if isp.ZDsdfGl3 != nil {
			glavna3 += val(isp.ZDsdfGl3) + "   "
		}
		if isp.ZKvarGl3 != nil {
			glavna3 += "Kvar u:" + val(isp.ZKvarGl3) + "   "
		}
		if isp.ZPrstGl3 != nil {
			glavna3 += val(isp.ZPrstGl3) + "   "
		}
		if isp.ZZMSPGl3 != nil {
			glavna3 += val(isp.ZZMSPGl3) + "   "
		}
		if isp.ZUzmsGl3 != nil {
			glavna3 += val(isp.ZUzmsGl3) + "   "
		}
		if isp.ZRapuGl3 != nil {
			glavna3 += val(isp.ZRapuGl3) + "   "
		}
		if isp.IdZTelePocGl3 != nil {
			glavna3 += val(isp.IdZTelePocGl3) + "   "
		}
		if isp.IdZTeleKrajGl3 != nil {
			glavna3 += val(isp.IdZTeleKrajGl3) + "   "
		}
		if isp.ZLokkGl3 != nil {
			glavna3 += "Lokator kvara:" + val(isp.ZLokkGl3) + "km "
		}
	}

	// ---------------- REZERVNA ----------------
	if hasAny(isp.ZDisRez, isp.ZKvarRez, isp.ZPrstRez, isp.ZZMSPRez) {
		rezervna = "Dopunska (rezervna) zaštita 1\n"
		if isp.ZDisRez != nil {
			rezervna += val(isp.ZDisRez) + "   "
		}
		if isp.ZKvarRez != nil {
			rezervna += "Kvar u:" + val(isp.ZKvarRez) + "   "
		}
		if isp.ZPrstRez != nil {
			rezervna += val(isp.ZPrstRez) + "   "
		}
		if isp.ZZMSPRez != nil {
			rezervna += val(isp.ZZMSPRez) + "   "
		}
	}

	// ---------------- REZERVNA 2 ----------------
	if hasAny(isp.ZDisRez2, isp.ZKvarRez2, isp.ZPrstRez2, isp.ZZMSPRez2) {
		rezervna2 = "Dopunska (rezervna) zaštita 2\n"
		if isp.ZDisRez2 != nil {
			rezervna2 += val(isp.ZDisRez2) + "   "
		}
		if isp.ZKvarRez2 != nil {
			rezervna2 += "Kvar u:" + val(isp.ZKvarRez2) + "   "
		}
		if isp.ZPrstRez2 != nil {
			rezervna2 += val(isp.ZPrstRez2) + "   "
		}
		if isp.ZZMSPRez2 != nil {
			rezervna2 += val(isp.ZZMSPRez2) + "   "
		}
	}

	// ---------------- PREKIDAC ----------------
	if hasAny(isp.ZPrekVn, isp.ZPrekNn) {
		if isp.ZPrekVn != nil {
			if isp.Fup == "01" || isp.Fup == "18" || isp.Fup == "22" {
				prekidac += "Prekidač VN strana:" + val(isp.ZPrekVn) + "   "
			} else {
				prekidac += val(isp.ZPrekVn) + "   "
			}
		}
		if isp.ZPrekNn != nil {
			prekidac += "Prekidač SN strana:" + val(isp.ZPrekNn) + "   "
		}
	}

	// ---------------- NEELEKTRICNA ----------------
	if hasAny(isp.ZNel1, isp.ZNel2, isp.ZNel3) {
		if isp.ZNel1 != nil {
			neelektricna += val(isp.ZNel1) + "   "
		}
		if isp.ZNel2 != nil {
			neelektricna += val(isp.ZNel2) + "   "
		}
		if isp.ZNel3 != nil {
			neelektricna += val(isp.ZNel3) + "   "
		}
	}

	// ---------------- SABIRNICE ----------------
	if hasAny(isp.ZSabzSab, isp.ZOtprSab, isp.ZKvarGl1) {
		sab1 = "Sabirnička zaštita 1\n"
		if isp.ZSabzSab != nil {
			sab1 += val(isp.ZSabzSab) + "   "
		}
		if isp.ZOtprSab != nil {
			sab1 += "Zas. od otkaza prek.:" + val(isp.ZOtprSab) + "   "
		}
		if isp.ZKvarGl1 != nil {
			sab1 += "Kvar u:" + val(isp.ZKvarGl1) + "   "
		}
	}

	if hasAny(isp.ZSabzSab2, isp.ZOtprSab2, isp.ZKvarGl2) {
		sab2 = "Sabirnička zaštita 2\n"
		if isp.ZSabzSab2 != nil {
			sab2 += val(isp.ZSabzSab2) + "   "
		}
		if isp.ZOtprSab2 != nil {
			sab2 += "Zas. od otkaza prek.:" + val(isp.ZOtprSab2) + "   "
		}
		if isp.ZKvarGl2 != nil {
			sab2 += "Kvar u:" + val(isp.ZKvarGl2) + "   "
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

	opis := zastita
	if isp.Opis != "" {
		if opis != "" {
			opis += "\n"
		}
		opis += isp.Opis
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
