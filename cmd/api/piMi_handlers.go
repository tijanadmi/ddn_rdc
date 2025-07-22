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
	Fup       string `form:"fup" binding:"required"`
}

type listPiMmByPageRequest struct {
	StartDate string `form:"start_date" binding:"required"`
	EndDate   string `form:"end_date" binding:"required"`
	TipD      int    `form:"tipd" binding:"required,min=1"`
	Fup       string `form:"fup" binding:"required"`
	PageID    int32  `form:"page_id" binding:"required,min=1"`
	PageSize  int32  `form:"page_size" binding:"required,min=5,max=100"`
}

type listPiMmT4Request struct {
	StartDate string `form:"start_date" binding:"required"`
	EndDate   string `form:"end_date" binding:"required"`
}

type listPiMmT4ByPageRequest struct {
	StartDate string `form:"start_date" binding:"required"`
	EndDate   string `form:"end_date" binding:"required"`
	PageID    int32  `form:"page_id" binding:"required,min=1"`
	PageSize  int32  `form:"page_size" binding:"required,min=5,max=100"`
}

// Definišemo strukturu za response
type listPiMmResponse struct {
	Total     int            `json:"total"`
	Dogadjaji []*models.PiMM `json:"dogadjaji"`
}

type listPiMmT4Response struct {
	Total       int              `json:"total"`
	DogadjajiT4 []*models.PiMMT4 `json:"dogadjaji"`
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
		Fup:       req.Fup,
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

func (server *Server) listPiMMT4(ctx *gin.Context) {
	var req listPiMmT4Request
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := models.ListPiMMT4Params{
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	}

	pimm, count, err := server.store.GetPiMMT4ByParams(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	rsp := listPiMmT4Response{
		Total:       count,
		DogadjajiT4: pimm,
	}

	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) listPiMMT4ByPage(ctx *gin.Context) {
	var req listPiMmT4ByPageRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := models.ListPiMMT4ParamsByPage{
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Limit:     req.PageSize,
		Offset:    (req.PageID - 1) * req.PageSize,
	}

	pimm, count, err := server.store.GetPiMMT4ByParamsByPage(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	rsp := listPiMmT4Response{
		Total:       count,
		DogadjajiT4: pimm,
	}

	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) listPiMMByPage(ctx *gin.Context) {
	var req listPiMmByPageRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := models.ListPiMMParamsByPage{
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Tipd:      req.TipD,
		Fup:       req.Fup,
		Limit:     req.PageSize,
		Offset:    (req.PageID - 1) * req.PageSize,
	}

	pimm, count, err := server.store.GetPiMMByParamsByPage(ctx, arg)
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

// func (server *Server) listPiMMByPageT4(ctx *gin.Context) {
// 	var req listPiMmT4Request
// 	if err := ctx.ShouldBindQuery(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		return
// 	}

// 	arg := models.ListPiMMT4Params{
// 		StartDate: req.StartDate,
// 		EndDate:   req.EndDate,
// 		Limit:     req.PageSize,
// 		Offset:    (req.PageID - 1) * req.PageSize,
// 	}

// 	pimm, count, err := server.store.GetPiMMT4ByParams(ctx, arg)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}
// 	rsp := listPiMmT4Response{
// 		Total:       count,
// 		DogadjajiT4: pimm,
// 	}

// 	ctx.JSON(http.StatusOK, rsp)
// }
