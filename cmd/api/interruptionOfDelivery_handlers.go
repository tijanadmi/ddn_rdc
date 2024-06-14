package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tijanadmi/ddn_rdc/models"
)

func (server *Server) getDDNInterruptionOfDeliveryById(ctx *gin.Context) {
	id := ctx.Param("id")
	ID, err := strconv.Atoi(id)
	if err != nil {
		
		ctx.JSON(http.StatusUnprocessableEntity, errorResponse(err))
		return
	}

	mrc, err := server.store.GetDDNInterruptionOfDeliveryById(ctx, ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, mrc)
}

func (server *Server) listDDNInterruptionOfDeliveryP(ctx *gin.Context) {
	var req listMrcRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := models.ListInterruptionParams{
		Ind: "P",
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	mrcs, err := server.store.GetDDNInterruptionOfDelivery(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, mrcs)
}

func (server *Server) listDDNInterruptionOfDeliveryK(ctx *gin.Context) {
	var req listMrcRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := models.ListInterruptionParams{
		Ind: "K",
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	mrcs, err := server.store.GetDDNInterruptionOfDelivery(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, mrcs)
}
