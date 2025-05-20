package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tijanadmi/ddn_rdc/models"
)

type listPiMmRequest struct {
	StartDate string `form:"start_date" binding:"required"`
	EndDate   string `form:"end_date" binding:"required"`
	TipD      int    `form:"tipd" binding:"required,min=1"`
}

// Definišemo strukturu za response
type listPiMmResponse struct {
	Total     int            `json:"total"`
	Dogadjaji []*models.PiMM `json:"dogadjaji"`
}

func (server *Server) listPiMM(ctx *gin.Context) {
	var req listPiMmRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := models.ListPiMMParams{
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Tipd:      req.TipD,
	}

	pimm, count, err := server.store.GetPiMMByParams(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	rsp := listPiMmResponse{
		Total:     count,
		Dogadjaji: pimm,
	}

	ctx.JSON(http.StatusOK, rsp)
}
