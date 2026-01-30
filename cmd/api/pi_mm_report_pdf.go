package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tijanadmi/ddn_rdc/models"
	"github.com/tijanadmi/ddn_rdc/pdf"
)

type piMMReportPDFRequest struct {
	StartDate string `form:"start_date" binding:"required"`
	EndDate   string `form:"end_date" binding:"required"`
	Tipd      int    `form:"tipd" binding:"required"`
	Kom       string `form:"komisija" binding:"required"`
}

func (server *Server) getPiMMReportPDF(ctx *gin.Context) {
	var req piMMReportPDFRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := models.ListPiMMByParam{
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Tipd:      req.Tipd,
		Kom:       req.Kom,
	}

	// 1️⃣ DB → grupisani report
	report, err := server.store.GetPiMMReportByParams(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// 2️⃣ Report → PDF bytes
	pdfBytes, err := pdf.GeneratePiMMReportPDF(report)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// 3️⃣ HEADERS
	ctx.Header("Content-Type", "application/pdf")
	ctx.Header("Content-Disposition", "inline; filename=pi_mm_report.pdf")
	ctx.Header("Content-Length", strconv.Itoa(len(pdfBytes)))

	ctx.Data(http.StatusOK, "application/pdf", pdfBytes)
}
