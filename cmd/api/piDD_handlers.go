package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tijanadmi/ddn_rdc/models"
)

type listPiDDRequest struct {
	Datizv string `form:"datizv" binding:"required"`
	TipD   int    `form:"tipd" binding:"required,min=1"`
	IdSMrc string `form:"id_s_mrc" binding:"required"`
	Fup    string `form:"fup" binding:"required"`
}

type listPiDDByPageRequest struct {
	Datizv   string `form:"datizv" binding:"required"`
	TipD     int    `form:"tipd" binding:"required,min=1"`
	Fup      string `form:"fup" binding:"required"`
	IdSMrc   string `form:"id_s_mrc" binding:"required"`
	PageID   int32  `form:"page_id" binding:"required,min=1"`
	PageSize int32  `form:"page_size" binding:"required,min=5,max=100"`
}

type listPiDDT4Request struct {
	Datizv string `form:"datizv" binding:"required"`
	IdSMrc string `form:"id_s_mrc" binding:"required"`
}

// Definišemo strukturu za response
type listPiDDResponse struct {
	Total     int            `json:"total"`
	Dogadjaji []*models.PiDD `json:"dogadjaji"`
}

/*type listPiDDT4Response struct {
	Total       int              `json:"total"`
	DogadjajiT4 []*models.PiDDT4 `json:"dogadjaji"`
}*/

func (server *Server) listPiDD(ctx *gin.Context) {
	var req listPiDDRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := models.ListPiDDParams{
		Datizv: req.Datizv,
		Tipd:   req.TipD,
		IdSMrc: req.IdSMrc,
		Fup:       req.Fup,
	}

	pidd, count, err := server.store.GetPiDDByParams(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	rsp := listPiDDResponse{
		Total:     count,
		Dogadjaji: pidd,
	}

	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) listPiDDT4(ctx *gin.Context) {
	var req listPiDDT4Request
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := models.ListPiDDT4Params{
		Datizv: req.Datizv,
		IdSMrc: req.IdSMrc,
	}

	pidd, count, err := server.store.GetPiDDT4ByParams(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	rsp := listPiMmT4Response{
		Total:       count,
		DogadjajiT4: pidd,
	}

	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) listPiDDByPage(ctx *gin.Context) {
	var req listPiDDByPageRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := models.ListPiDDParamsByPage{
		Datizv: req.Datizv,
		Tipd:   req.TipD,
		Fup:    req.Fup,
		IdSMrc: req.IdSMrc,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	pidd, count, err := server.store.GetPiDDByParamsByPage(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	rsp := listPiDDResponse{
		Total:     count,
		Dogadjaji: pidd,
	}

	ctx.JSON(http.StatusOK, rsp)
}
