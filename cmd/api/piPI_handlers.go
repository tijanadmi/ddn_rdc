package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tijanadmi/ddn_rdc/models"
)

/*type listPiDDT4Response struct {
	Total       int              `json:"total"`
	DogadjajiT4 []*models.PiDDT4 `json:"dogadjaji"`
}*/

func (server *Server) listPiPI(ctx *gin.Context) {
	var req listPiDDRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := models.ListPiDDParams{
		Datizv: req.Datizv,
		Tipd:   req.TipD,
		IdSMrc: req.IdSMrc,
	}

	pidd, count, err := server.store.GetPiPIByParams(ctx, arg)
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

func (server *Server) listPiPIT4(ctx *gin.Context) {
	var req listPiDDT4Request
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := models.ListPiDDT4Params{
		Datizv: req.Datizv,
		IdSMrc: req.IdSMrc,
	}

	pidd, count, err := server.store.GetPiPIT4ByParams(ctx, arg)
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

func (server *Server) listPiPIByPage(ctx *gin.Context) {
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

	pidd, count, err := server.store.GetPiPIByParamsByPage(ctx, arg)
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
