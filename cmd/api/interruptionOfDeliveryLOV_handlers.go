package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tijanadmi/ddn_rdc/models"
)

func (server *Server) getMrcById(ctx *gin.Context) {
	id := ctx.Param("id")
	mrcID, err := strconv.Atoi(id)
	if err != nil {
		
		ctx.JSON(http.StatusUnprocessableEntity, errorResponse(err))
		return
	}

	mrc, err := server.store.GetMrcById(ctx, mrcID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, mrc)
}

func (server *Server) getSTipPrekById(ctx *gin.Context) {
	id := ctx.Param("id")
	mrcID, err := strconv.Atoi(id)
	if err != nil {
		
		ctx.JSON(http.StatusUnprocessableEntity, errorResponse(err))
		return
	}

	mrc, err := server.store.GetSTipPrekById(ctx, mrcID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, mrc)
}

func (server *Server) getSVrPrekById(ctx *gin.Context) {
	id := ctx.Param("id")
	mrcID, err := strconv.Atoi(id)
	if err != nil {
		
		ctx.JSON(http.StatusUnprocessableEntity, errorResponse(err))
		return
	}

	mrc, err := server.store.GetSVrPrekById(ctx, mrcID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, mrc)
}

func (server *Server) getSUzrokPrekById(ctx *gin.Context) {
	id := ctx.Param("id")
	mrcID, err := strconv.Atoi(id)
	if err != nil {
		
		ctx.JSON(http.StatusUnprocessableEntity, errorResponse(err))
		return
	}

	mrc, err := server.store.GetSUzrokPrekById(ctx, mrcID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, mrc)
}

func (server *Server) getSPoduzrokPrekById(ctx *gin.Context) {
	id := ctx.Param("id")
	mrcID, err := strconv.Atoi(id)
	if err != nil {
		
		ctx.JSON(http.StatusUnprocessableEntity, errorResponse(err))
		return
	}

	mrc, err := server.store.GetSPoduzrokPrekById(ctx, mrcID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, mrc)
}

func (server *Server) getSMernaMestaById(ctx *gin.Context) {
	id := ctx.Param("id")
	mrcID, err := strconv.Atoi(id)
	if err != nil {
		
		ctx.JSON(http.StatusUnprocessableEntity, errorResponse(err))
		return
	}

	mrc, err := server.store.GetSMernaMestaById(ctx, mrcID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, mrc)
}

type listMrcRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=100"`
}

func (server *Server) listMrcs(ctx *gin.Context) {
	var req listMrcRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := models.ListLimitOffsetParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	mrcs, err := server.store.GetSMrc(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, mrcs)
}

func (server *Server) listTipPrek(ctx *gin.Context) {
	var req listMrcRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := models.ListLimitOffsetParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	mrcs, err := server.store.GetSTipPrek(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, mrcs)
}

func (server *Server) listVrPrek(ctx *gin.Context) {
	var req listMrcRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := models.ListLimitOffsetParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	mrcs, err := server.store.GetSVrPrek(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, mrcs)
}

func (server *Server) listUzrokPrek(ctx *gin.Context) {
	var req listMrcRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := models.ListLimitOffsetParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	mrcs, err := server.store.GetSUzrokPrek(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, mrcs)
}

func (server *Server) listPoduzrokPrek(ctx *gin.Context) {
	var req listMrcRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := models.ListLimitOffsetParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	mrcs, err := server.store.GetSPoduzrokPrek(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, mrcs)
}

func (server *Server) listMernaMesta(ctx *gin.Context) {
	var req listMrcRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := models.ListLimitOffsetParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	mrcs, err := server.store.GetSMernaMesta(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, mrcs)
}