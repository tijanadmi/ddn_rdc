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

type listAllInteraptionOfDeliveryRequest struct {
	Mrc       string `form:"mrc" binding:"required"`
	StartDate string `form:"start_date" binding:"required"`
	EndDate   string `form:"end_date" binding:"required"`
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

type DDNInterruptionOfDeliveryPExcel struct {
	Id           string `json:"id"`
	Mrc          string `json:"rdc"`
	ObOpis       string `json:"proizvodni_objekat"`
	PoljeNaziv   string `json:"generator"`
	Vrepoc       string `json:"pocetak"`
	Vrezav       string `json:"kraj"`
	Trajanje     string `json:"trajanje"`
	VrstaPrek    string `json:"vrsta_prekida"`
	PodvrstaPrek string `json:"planiran_neplaniran"`
	Uzrok        string `json:"uzrok"`
	PoduzrokPrek string `json:"poduzrok"`
	Snaga        string `json:"snaga"`
	Opis         string `json:"opis"`
	Bi           string `json:"bi"`
}

type listDDNPAllResponse struct {
	Total   int                                `json:"total"`
	Prekidi []*DDNInterruptionOfDeliveryPExcel `json:"prekidip"`
}

func (server *Server) listAllDDNInterruptionOfDeliveryP(ctx *gin.Context) {
	var req listAllInteraptionOfDeliveryRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := models.ListInterruptionParams{
		Mrc:       req.Mrc,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Ind:       "P",
	}

	mrcs, count, err := server.store.GetAllDDNInterruptionOfDelivery(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := &listDDNPAllResponse{
		Total: count,
	}
	for _, m := range mrcs {
		rsp.Prekidi = append(rsp.Prekidi, &DDNInterruptionOfDeliveryPExcel{
			Id:           m.Id,
			Mrc:          m.Mrc,
			ObOpis:       m.ObOpis,
			PoljeNaziv:   m.PoljeNaziv,
			Vrepoc:       m.Vrepoc,
			Vrezav:       m.Vrezav,
			Trajanje:     m.Trajanje,
			VrstaPrek:    m.VrstaPrek,
			PodvrstaPrek: m.PodvrstaPrek,
			Uzrok:        m.Uzrok,
			PoduzrokPrek: m.PoduzrokPrek,
			Snaga:        m.Snaga,
			Opis:         m.Opis,
			Bi:           m.Bi,
		})
	}

	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) listDDNInterruptionOfDeliveryP(ctx *gin.Context) {
	var req listInteraptionOfDeliveryRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := models.ListInterruptionWithPaginationParams{
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

	arg := models.ListInterruptionWithPaginationParams{
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

type DDNInterruptionOfDeliveryKExcel struct {
	Id           string `json:"id"`
	Mrc          string `json:"rdc"`
	ObOpis       string `json:"proizvodni_objekat"`
	Vrepoc       string `json:"pocetka"`
	Vrezav       string `json:"zavrsetka"`
	Trajanje     string `json:"trajanje"`
	VrstaPrek    string `json:"vrsta_prekida"`
	Uzrok        string `json:"uzrok"`
	PoduzrokPrek string `json:"poduzrok"`
	Snaga        string `json:"ispala_snaga"`
	MernaMesta   string `json:"merna_mesta"`
	BrojMesta    string `json:"broj_mernih_mesta"`
	Opis         string `json:"opis"`
	Bi           string `json:"bi"`
}

type listDDNKAllResponse struct {
	Total   int                                `json:"total"`
	Prekidi []*DDNInterruptionOfDeliveryKExcel `json:"prekidik"`
}

func (server *Server) listAllDDNInterruptionOfDeliveryK(ctx *gin.Context) {
	var req listAllInteraptionOfDeliveryRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := models.ListInterruptionParams{
		Mrc:       req.Mrc,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Ind:       "K",
	}

	mrcs, count, err := server.store.GetAllDDNInterruptionOfDelivery(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := &listDDNKAllResponse{
		Total: count,
	}
	for _, m := range mrcs {
		rsp.Prekidi = append(rsp.Prekidi, &DDNInterruptionOfDeliveryKExcel{
			Id:           m.Id,
			Mrc:          m.Mrc,
			ObOpis:       m.ObOpis,
			Vrepoc:       m.Vrepoc,
			Vrezav:       m.Vrezav,
			Trajanje:     m.Trajanje,
			VrstaPrek:    m.VrstaPrek,
			Uzrok:        m.Uzrok,
			PoduzrokPrek: m.PoduzrokPrek,
			Snaga:        m.Snaga,
			MernaMesta:   m.MernaMesta,
			BrojMesta:    m.BrojMesta,
			Opis:         m.Opis,
			Bi:           m.Bi,
		})
	}

	ctx.JSON(http.StatusOK, rsp)
}
