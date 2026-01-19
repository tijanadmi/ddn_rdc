package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tijanadmi/ddn_rdc/models"
	"github.com/tijanadmi/ddn_rdc/token"
)

type listInteraptionOfDeliveryByPageRequest struct {
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

func (server *Server) listExcelDDNInterruptionOfDeliveryP(ctx *gin.Context) {
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

	rsp := listDDNPResponse{
		Total:   count,
		Prekidi: mrcs,
	}

	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) listDDNInterruptionOfDeliveryPByPage(ctx *gin.Context) {
	var req listInteraptionOfDeliveryByPageRequest
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

	mrcs, count, err := server.store.GetDDNInterruptionOfDeliveryByPage(ctx, arg)
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

func (server *Server) listDDNInterruptionOfDeliveryKByPage(ctx *gin.Context) {
	var req listInteraptionOfDeliveryByPageRequest
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

	mrcs, count, err := server.store.GetDDNInterruptionOfDeliveryByPage(ctx, arg)
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

func (server *Server) listExcelDDNInterruptionOfDeliveryK(ctx *gin.Context) {
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

	rsp := listDDNKResponse{
		Total:   count,
		Prekidi: mrcs,
	}

	ctx.JSON(http.StatusOK, rsp)
}

/***************** Create Prekid Proizvodnje ***************/

type createDDNPrekidIspRequest struct {
	IdSMrc            int    `json:"id_s_mrc" binding:"required"`
	IdTipob           int    `json:"id_tipob" binding:"required"`
	ObId              int    `json:"ob_id" binding:"required"`
	Vrepoc            string `json:"vrepoc" binding:"required"`
	Vrezav            string `json:"vrezav"`
	IdSVrPrek         int    `json:"id_s_vr_prek" binding:"required"`
	IdSUzrokPrek      int    `json:"id_s_uzrok_prek"`
	Snaga             string `json:"snaga"`
	Opis              string `json:"opis"`
	P2TrafId          int    `json:"p2_traf_id"`
	IdSPoduzrokPrek   int    `json:"id_s_poduzrok_prek"`
	IdTipObjektaNdc   string `json:"id_tip_objekta_ndc"`
	IdTipDogadjajaNdc string `json:"id_tip_dogadjaja_ndc"`
}

type createDDNPrekidPrRequest struct {
	IdSMrc            int    `json:"id_s_mrc" binding:"required"`
	IdTipob           int    `json:"id_tipob" binding:"required"`
	ObId              int    `json:"ob_id" binding:"required"`
	Vrepoc            string `json:"vrepoc" binding:"required"`
	Vrezav            string `json:"vrezav"`
	IdSVrPrek         int    `json:"id_s_vr_prek" binding:"required"`
	IdSUzrokPrek      int    `json:"id_s_uzrok_prek"`
	Snaga             string `json:"snaga"`
	Opis              string `json:"opis"`
	P2TrafId          int    `json:"p2_traf_id"`
	IdSPoduzrokPrek   int    `json:"id_s_poduzrok_prek"`
	IdTipObjektaNdc   string `json:"id_tip_objekta_ndc"`
	IdTipDogadjajaNdc string `json:"id_tip_dogadjaja_ndc"`
}

func parseDateTime(value string) (time.Time, error) {
	layout := "02.01.2006 15:04"
	return time.Parse(layout, value)
}

func validateDDNPrekidPrInput(req createDDNPrekidPrRequest) error {
	vrepoc, err := parseDateTime(req.Vrepoc)
	if err != nil {
		return fmt.Errorf("неисправан формат за време почетка (dd.mm.yyyy hh:mi)")
	}

	var vrezav time.Time
	if req.Vrezav != "" {
		vrezav, err = parseDateTime(req.Vrezav)
		if err != nil {
			return fmt.Errorf("неисправан формат за време завршетка (dd.mm.yyyy hh:mi)")
		}
		if vrezav.Before(vrepoc) {
			return fmt.Errorf("време завршетка не може бити пре времена почетка")
		}
	}

	if req.IdSVrPrek != 1 && req.IdSUzrokPrek == 0 {
		return fmt.Errorf("узрок прекида је обавезан ако врста прекида није 1")
	}

	if req.IdSUzrokPrek == 1 && req.IdSPoduzrokPrek == 0 {
		return fmt.Errorf("подузрок прекида је обавезан када је узрок прекида 1")
	}

	return nil
}

func (server *Server) CreateDDNPrekidPr(ctx *gin.Context) {

	var req createDDNPrekidPrRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Validacija
	err := validateDDNPrekidPrInput(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := models.CreateDDNInterruptionOfDeliveryPParams{
		IdSMrc:            req.IdSMrc,
		IdTipob:           req.IdTipob,
		ObId:              req.ObId,
		Vrepoc:            req.Vrepoc,
		Vrezav:            req.Vrezav,
		IdSVrPrek:         req.IdSVrPrek,
		IdSUzrokPrek:      req.IdSUzrokPrek,
		Snaga:             req.Snaga,
		Opis:              req.Opis,
		KorUneo:           payload.Username,
		P2TrafId:          req.P2TrafId,
		IdSPoduzrokPrek:   req.IdSPoduzrokPrek,
		IdTipObjektaNdc:   req.IdTipObjektaNdc,
		IdTipDogadjajaNdc: req.IdTipDogadjajaNdc,
	}

	id, err := server.store.InsertDDNInterruptionOfDeliveryP(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	pr, err := server.store.GetDDNInterruptionOfDeliveryById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, pr)
}

/***************** End Create Prekid Proizvodnje ***************/
/***************** Start Create Prekid Korisnika ***************/

type createDDNPrekidKRequest struct {
	IdSMrc          int    `json:"id_s_mrc" binding:"required"`
	IdTipob         int    `json:"id_tipob" binding:"required"`
	ObId            int    `json:"ob_id" binding:"required"`
	Vrepoc          string `json:"vrepoc" binding:"required"`
	Vrezav          string `json:"vrezav"`
	IdSVrPrek       int    `json:"id_s_vr_prek" binding:"required"`
	IdSUzrokPrek    int    `json:"id_s_uzrok_prek"`
	Snaga           string `json:"snaga"`
	Opis            string `json:"opis"`
	IdSPoduzrokPrek int    `json:"id_s_poduzrok_prek"`
	IdSMernaMesta   int    `json:"id_s_merna_mesta"`
	BrojMesta       string `json:"broj_mesta"`
}

func validateDDNPrekidKInput(req createDDNPrekidKRequest) error {
	vrepoc, err := parseDateTime(req.Vrepoc)
	if err != nil {
		return fmt.Errorf("неисправан формат за време почетка (dd.mm.yyyy hh:mi)")
	}

	var vrezav time.Time
	if req.Vrezav != "" {
		vrezav, err = parseDateTime(req.Vrezav)
		if err != nil {
			return fmt.Errorf("неисправан формат за време завршетка (dd.mm.yyyy hh:mi)")
		}
		if vrezav.Before(vrepoc) {
			return fmt.Errorf("време завршетка не може бити пре времена почетка")
		}
	}

	if req.IdSVrPrek != 1 && req.IdSUzrokPrek == 0 {
		return fmt.Errorf("узрок прекида је обавезан ако врста прекида није 1")
	}

	if req.IdSUzrokPrek == 1 && req.IdSPoduzrokPrek == 0 {
		return fmt.Errorf("подузрок прекида је обавезан када је узрок прекида 1")
	}

	return nil
}

func (server *Server) CreateDDNPrekidK(ctx *gin.Context) {

	var req createDDNPrekidKRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Validacija
	err := validateDDNPrekidKInput(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := models.CreateDDNInterruptionOfDeliveryKParams{
		IdSMrc:          req.IdSMrc,
		IdTipob:         req.IdTipob,
		ObId:            req.ObId,
		Vrepoc:          req.Vrepoc,
		Vrezav:          req.Vrezav,
		IdSVrPrek:       req.IdSVrPrek,
		IdSUzrokPrek:    req.IdSUzrokPrek,
		Snaga:           req.Snaga,
		Opis:            req.Opis,
		KorUneo:         payload.Username,
		IdSPoduzrokPrek: req.IdSPoduzrokPrek,
		IdSMernaMesta:   req.IdSMernaMesta,
		BrojMesta:       req.BrojMesta,
	}

	id, err := server.store.InsertDDNInterruptionOfDeliveryK(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	pr, err := server.store.GetDDNInterruptionOfDeliveryById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, pr)
}

/***************** End Create Prekid Korisnika ***************/

/***************** Start Update Prekid Proizvodnje ***************/
type updateDDNPrekidPrRequest struct {
	IdSMrc            int    `json:"id_s_mrc" binding:"required"`
	IdTipob           int    `json:"id_tipob" binding:"required"`
	ObId              int    `json:"ob_id" binding:"required"`
	Vrepoc            string `json:"vrepoc" binding:"required"`
	Vrezav            string `json:"vrezav"`
	IdSVrPrek         int    `json:"id_s_vr_prek" binding:"required"`
	IdSUzrokPrek      int    `json:"id_s_uzrok_prek"`
	Snaga             string `json:"snaga"`
	Opis              string `json:"opis"`
	P2TrafId          int    `json:"p2_traf_id"`
	IdSPoduzrokPrek   int    `json:"id_s_poduzrok_prek"`
	IdTipDogadjajaNdc string `json:"id_tip_dogadjaja_ndc"`
	IdTipObjektaNdc   string `json:"id_tip_objekta_ndc"`
}

func (server *Server) UpdateDDNPrekidPr(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("neispravan id")))
		return
	}

	version, err := strconv.Atoi(ctx.Param("version"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("neispravna verzija")))
		return
	}

	var req updateDDNPrekidPrRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	createReq := createDDNPrekidPrRequest{
		IdSMrc:            req.IdSMrc,
		IdTipob:           req.IdTipob,
		ObId:              req.ObId,
		Vrepoc:            req.Vrepoc,
		Vrezav:            req.Vrezav,
		IdSVrPrek:         req.IdSVrPrek,
		IdSUzrokPrek:      req.IdSUzrokPrek,
		Snaga:             req.Snaga,
		Opis:              req.Opis,
		P2TrafId:          req.P2TrafId,
		IdSPoduzrokPrek:   req.IdSPoduzrokPrek,
		IdTipObjektaNdc:   req.IdTipObjektaNdc,
		IdTipDogadjajaNdc: req.IdTipDogadjajaNdc,
	}
	// Validacija – ista kao kod create
	if err := validateDDNPrekidPrInput(createReq); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := models.CreateDDNInterruptionOfDeliveryPParams{
		IdSMrc:            req.IdSMrc,
		IdTipob:           req.IdTipob,
		ObId:              req.ObId,
		Vrepoc:            req.Vrepoc,
		Vrezav:            req.Vrezav,
		IdSVrPrek:         req.IdSVrPrek,
		IdSUzrokPrek:      req.IdSUzrokPrek,
		Snaga:             req.Snaga,
		Opis:              req.Opis,
		KorUneo:           payload.Username,
		P2TrafId:          req.P2TrafId,
		IdSPoduzrokPrek:   req.IdSPoduzrokPrek,
		IdTipObjektaNdc:   req.IdTipObjektaNdc,
		IdTipDogadjajaNdc: req.IdTipDogadjajaNdc,
	}

	err = server.store.UpdateDDNInterruptionOfDeliveryP(ctx, id, version, arg)
	if err != nil {
		ctx.JSON(http.StatusConflict, errorResponse(err))
		return
	}

	pr, err := server.store.GetDDNInterruptionOfDeliveryById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, pr)
}

/***************** End Update Prekid Proizvodnje ***************/

/***************** Start Update Prekid Korisnika ***************/

type updateDDNPrekidKRequest struct {
	IdSMrc          int    `json:"id_s_mrc" binding:"required"`
	IdTipob         int    `json:"id_tipob" binding:"required"`
	ObId            int    `json:"ob_id" binding:"required"`
	Vrepoc          string `json:"vrepoc" binding:"required"`
	Vrezav          string `json:"vrezav"`
	IdSVrPrek       int    `json:"id_s_vr_prek" binding:"required"`
	IdSUzrokPrek    int    `json:"id_s_uzrok_prek"`
	Snaga           string `json:"snaga"`
	Opis            string `json:"opis"`
	IdSPoduzrokPrek int    `json:"id_s_poduzrok_prek"`
	IdSMernaMesta   int    `json:"id_s_merna_mesta"`
	BrojMesta       string `json:"broj_mesta"`
}

func (server *Server) UpdateDDNPrekidK(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("neispravan id")))
		return
	}

	version, err := strconv.Atoi(ctx.Param("version"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("neispravna verzija")))
		return
	}

	var req updateDDNPrekidKRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	createReq := createDDNPrekidKRequest{
		IdSMrc:          req.IdSMrc,
		IdTipob:         req.IdTipob,
		ObId:            req.ObId,
		Vrepoc:          req.Vrepoc,
		Vrezav:          req.Vrezav,
		IdSVrPrek:       req.IdSVrPrek,
		IdSUzrokPrek:    req.IdSUzrokPrek,
		Snaga:           req.Snaga,
		Opis:            req.Opis,
		IdSPoduzrokPrek: req.IdSPoduzrokPrek,
		IdSMernaMesta:   req.IdSMernaMesta,
		BrojMesta:       req.BrojMesta,
	}
	// Validacija – ista kao kod create
	if err := validateDDNPrekidKInput(createReq); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := models.CreateDDNInterruptionOfDeliveryKParams{
		IdSMrc:          req.IdSMrc,
		IdTipob:         req.IdTipob,
		ObId:            req.ObId,
		Vrepoc:          req.Vrepoc,
		Vrezav:          req.Vrezav,
		IdSVrPrek:       req.IdSVrPrek,
		IdSUzrokPrek:    req.IdSUzrokPrek,
		Snaga:           req.Snaga,
		Opis:            req.Opis,
		KorUneo:         payload.Username,
		IdSPoduzrokPrek: req.IdSPoduzrokPrek,
		IdSMernaMesta:   req.IdSMernaMesta,
		BrojMesta:       req.BrojMesta,
	}

	err = server.store.UpdateDDNInterruptionOfDeliveryK(ctx, id, version, arg)
	if err != nil {
		ctx.JSON(http.StatusConflict, errorResponse(err))
		return
	}

	pr, err := server.store.GetDDNInterruptionOfDeliveryById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, pr)
}

/***************** End Update Prekid Korisnika ***************/

/***************** Start Update Prekid Proizvodnje BI ***************/
type updateDDNPrekidIspBIRequest struct {
	BI int `json:"bi" binding:"oneof=0 1"`
}

func (server *Server) UpdateDDNPrekidIspBI(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("neispravan id")))
		return
	}

	version, err := strconv.Atoi(ctx.Param("version"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("neispravna verzija")))
		return
	}

	var req updateDDNPrekidIspBIRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = server.store.UpdateDDNInterruptionOfDeliveryBI(
		ctx,
		id,
		version,
		req.BI,
	)
	if err != nil {
		ctx.JSON(http.StatusConflict, errorResponse(err))
		return
	}

	pr, err := server.store.GetDDNInterruptionOfDeliveryById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, pr)
}

/***************** End Update Prekid Proizvodnje BI ***************/

/***** Delete Prekid Isporuke **********/

func (server *Server) deleteDDNInterruptionOfDelivery(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("neispravan id")))
		return
	}

	version, err := strconv.Atoi(ctx.Param("version"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("neispravna verzija")))
		return
	}

	err = server.store.DeleteDDNInterruptionOfDelivery(ctx, id, version)
	if err != nil {
		// optimistic locking → pogrešna verzija
		ctx.JSON(http.StatusConflict, errorResponse(err))
		return
	}

	ctx.Status(http.StatusNoContent)
}

/***** Delete Prekid Isporuke **********/
