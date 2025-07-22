package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tijanadmi/ddn_rdc/models"
)

type listGDRadapuMesRequest struct {
	Godina string `form:"godina" binding:"required"`
}

func (server *Server) listPGDRadapuMes(ctx *gin.Context) {
	var req listGDRadapuMesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := models.ListPGD{
		Godina: req.Godina,
	}

	pimm, err := server.store.GetPGDRadapuMes(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	rsp := pimm

	ctx.JSON(http.StatusOK, rsp)
}
