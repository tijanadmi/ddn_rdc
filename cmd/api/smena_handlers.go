package api

import (
	"fmt"
	"net/http"

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
