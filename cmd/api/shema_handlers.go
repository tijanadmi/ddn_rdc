package api

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func (server *Server) getShemeByOrg(ctx *gin.Context) {

	id := ctx.Param("id")
	orgID, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("neispravan idOrg")))
		return
	}

	if orgID == 0 {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("idOrg je obavezan")))
		return
	}

	// ---------------- JEDNOPOLNE ----------------
	jednopolne, err := server.store.GetShemeByOrg(ctx, orgID, "JS")
	if err != nil {
		fmt.Printf("Greška pri čitanju jednopolnih šema: %v\n", err)

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// ---------------- DISPOZICIONE ----------------
	dispozicione, err := server.store.GetShemeByOrg(ctx, orgID, "DS")
	if err != nil {
		fmt.Printf("Greška pri čitanju dispozicionih šema: %v\n", err)

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// ---------------- RESPONSE ----------------
	ctx.JSON(http.StatusOK, gin.H{
		"jednopolne":   jednopolne,
		"dispozicione": dispozicione,
	})
}

func (server *Server) getShemaPDF(ctx *gin.Context) {

	fmt.Println("Dohvatanje šeme PDF...")
	idStr := ctx.Param("id")

	id, err := strconv.Atoi(idStr)
	fmt.Printf("ID šeme: %d\n", id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("neispravan ID")))
		return
	}

	// ---------------- PUTANJA + IME ----------------
	putanja, imeDok, err := server.store.GetShemaPutanjaByID(ctx, id)
	fmt.Println(putanja)
	if err != nil {
		fmt.Printf("Greška pri dohvatanju fajla: %v\n", err)

		ctx.JSON(http.StatusNotFound, errorResponse(errors.New("fajl nije pronađen")))
		return
	}

	baseUNC := `\sheme`

	var basePath string

	if os.Getenv("ENVIRONMENT") == "production" {
		basePath = "/mnt/sheme"
	} else {
		// DEV: uzmi root do \sheme iz pune UNC putanje
		idx := strings.Index(putanja, baseUNC)

		if idx != -1 {
			basePath = putanja[:idx+len(baseUNC)]
		} else {
			// fallback (ako format nije očekivan)
			basePath = baseUNC
		}
	}

	// izvuci deo posle \sheme
	idx := strings.Index(putanja, baseUNC)
	if idx == -1 {
		return
	}

	relative := putanja[idx+len(baseUNC):]

	// očisti početni backslash
	relative = strings.TrimLeft(relative, `\`)

	// Windows -> Linux
	relative = strings.ReplaceAll(relative, `\`, "/")

	fullPath := filepath.Join(basePath, relative)

	fmt.Println("Nova putanja:", fullPath)

	// ---------------- HEADERS ----------------
	ctx.Header("Content-Type", "application/pdf")

	// inline prikaz + ime fajla
	ctx.Header(
		"Content-Disposition",
		fmt.Sprintf(`inline; filename="%s"`, imeDok),
	)

	// ---------------- SERVIRANJE ----------------
	ctx.File(fullPath)
}
