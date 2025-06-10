package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tijanadmi/ddn_rdc/models"
)

type listInteraptionOfDeliveryRequest struct {
	Mrc       string `form:"mrc" binding:"required"`
	StartDate string `form:"start_date" binding:"required"`
	EndDate   string `form:"end_date" binding:"required"`
	PageID    int32  `form:"page_id" binding:"required,min=1"`
	PageSize  int32  `form:"page_size" binding:"required,min=5,max=100"`
}

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

// Definišemo strukturu za response
type listDDNPResponse struct {
	Total   int                                 `json:"total"`
	Prekidi []*models.DDNInterruptionOfDelivery `json:"prekidip"`
}

func (server *Server) listDDNInterruptionOfDeliveryP(ctx *gin.Context) {
	var req listInteraptionOfDeliveryRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := models.ListInterruptionParams{
		Mrc:       req.Mrc,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Ind:       "P",
		Limit:     req.PageSize,
		Offset:    (req.PageID - 1) * req.PageSize,
	}

	mrcs, count, err := server.store.GetDDNInterruptionOfDelivery(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	rsp := listDDNPResponse{
		Total:   count,
		Prekidi: mrcs,
	}

	ctx.JSON(http.StatusOK, rsp)
}

// Definišemo strukturu za response
type listDDNKResponse struct {
	Total   int                                 `json:"total"`
	Prekidi []*models.DDNInterruptionOfDelivery `json:"prekidik"`
}

func (server *Server) listDDNInterruptionOfDeliveryK(ctx *gin.Context) {
	var req listInteraptionOfDeliveryRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := models.ListInterruptionParams{
		Mrc:       req.Mrc,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Ind:       "K",
		Limit:     req.PageSize,
		Offset:    (req.PageID - 1) * req.PageSize,
	}

	mrcs, count, err := server.store.GetDDNInterruptionOfDelivery(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := listDDNKResponse{
		Total:   count,
		Prekidi: mrcs,
	}

	ctx.JSON(http.StatusOK, rsp)
}
