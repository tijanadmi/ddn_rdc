package api

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	db "github.com/tijanadmi/ddn_rdc/repository"
	"github.com/tijanadmi/ddn_rdc/token"
	"github.com/tijanadmi/ddn_rdc/util"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	config     util.Config
	store      db.DatabaseRepo
	tokenMaker token.Maker
	router     *gin.Engine
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(config util.Config, store db.DatabaseRepo) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}
	/*if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}*/

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	/*if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}*/

	// Configure CORS
	router.Use(cors.New(cors.Config{
		//AllowOrigins:     []string{"http://localhost:3000"},
		//AllowOrigins:     []string{"http://192.168.36.188", "http://192.168.36.197"},
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	router.POST("/users/login", server.loginUser)
	router.POST("/users/get_user_by_token", server.GetUserByToken)
	router.POST("/tokens/renew_access", server.renewAccessToken)
	router.GET("/mrc/:id", server.getMrcById)
	router.GET("/mrc", server.listMrcs)
	router.GET("/tipprek/:id", server.getSTipPrekById)
	router.GET("/tipprek", server.listTipPrek)
	router.GET("/vrprek/:id", server.getSVrPrekById)
	router.GET("/vrprek", server.listVrPrek)
	router.GET("/uzrokprek/:id", server.getSUzrokPrekById)
	router.GET("/uzrokprek", server.listUzrokPrek)
	router.GET("/poduzrokprek/:id", server.getSPoduzrokPrekById)
	router.GET("/poduzrokprek", server.listPoduzrokPrek)
	router.GET("/mernamesta/:id", server.getSMernaMestaById)
	router.GET("/mernamesta", server.listMernaMesta)
	router.GET("/obj/:id", server.getObjId)
	router.GET("/objtsrp", server.listObjTSRP)
	router.GET("/objheteve", server.listObjHETEVE)
	router.GET("/poljage", server.listPoljaGE)
	router.GET("/poljage/:id", server.getPoljeGEById)

	router.GET("/interruptionofdelivery/:id", server.getDDNInterruptionOfDeliveryById)
	router.GET("/interruptionofproduction", server.listDDNInterruptionOfDeliveryPByPage)
	router.GET("/interruptionofproduction_all", server.listAllDDNInterruptionOfDeliveryP)
	router.GET("/interruptionofproduction_excel", server.listExcelDDNInterruptionOfDeliveryP)
	router.GET("/interruptionofusers", server.listDDNInterruptionOfDeliveryKByPage)
	router.GET("/interruptionofusers_all", server.listAllDDNInterruptionOfDeliveryK)
	router.GET("/interruptionofusers_excel", server.listExcelDDNInterruptionOfDeliveryK)

	router.GET("/mesecni", server.listPiMM)
	router.GET("/mesecnip", server.listPiMMByPage)
	router.GET("/mesecnit4", server.listPiMMT4)
	router.GET("/mesecnit4p", server.listPiMMT4ByPage)

	router.GET("/dnevni", server.listPiDD)
	router.GET("/dnevnip", server.listPiDDByPage)
	router.GET("/dnevnit4p", server.listPiDDT4)

	router.GET("/pogonski", server.listPiPI)
	router.GET("/pogonskip", server.listPiPIByPage)
	router.GET("/pogonskit4", server.listPiPIT4)

	router.GET("/radapu_mes", server.listPGDRadapuMes)
	router.GET("/dapua", server.listPGDDapuA)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.POST("/createinterruptionofproduction", server.CreateDDNPrekidIsp)

	server.router = router
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
