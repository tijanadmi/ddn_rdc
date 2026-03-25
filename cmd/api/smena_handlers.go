package api

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
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
