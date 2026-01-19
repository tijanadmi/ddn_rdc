package api

import (
	"fmt"
	"strings"

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

	allowedOrigins := server.config.AllowedOrigins
	var origins []string
	if allowedOrigins != "" {
		origins = strings.Split(allowedOrigins, ",")
		for i := range origins {
			origins[i] = strings.TrimSpace(origins[i])
			fmt.Println("Allowed Origin:", origins[i])
		}
	} else {
		origins = []string{"http://localhost:5173"}
		fmt.Println("Allowed Origin: http://localhost:5173")
	}
	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins: origins,
		// AllowOrigins: []string{"http://192.168.36.188", "http://192.168.29.68:5173", "http://localhost:5173", "http://192.168.72.147:4000"},
		// AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	router.POST("/users/login", server.loginUser)
	router.POST("/users/logout", server.logoutUser)
	router.POST("/tokens/renew_access", server.renewAccessToken)
	// router.POST("/users/get_user_by_token", server.GetUserByToken)
	// router.POST("/tokens/renew_access", server.renewAccessToken)
	// router.GET("/mrc/:id", server.getMrcById)
	// router.GET("/mrc", server.listMrcs)
	// router.GET("/tipprek/:id", server.getSTipPrekById)
	// router.GET("/tipprek", server.listTipPrek)
	// router.GET("/vrprek/:id", server.getSVrPrekById)
	// router.GET("/vrprek", server.listVrPrek)
	// router.GET("/uzrokprek/:id", server.getSUzrokPrekById)
	// router.GET("/uzrokprek", server.listUzrokPrek)
	// router.GET("/poduzrokprek/:id", server.getSPoduzrokPrekById)
	// router.GET("/poduzrokprek", server.listPoduzrokPrek)
	// router.GET("/mernamesta/:id", server.getSMernaMestaById)
	// router.GET("/mernamesta", server.listMernaMesta)
	// router.GET("/obj/:id", server.getObjId)
	// router.GET("/objtsrp", server.listObjTSRP)
	// router.GET("/objheteve", server.listObjHETEVE)
	// router.GET("/poljage", server.listPoljaGE)
	// router.GET("/poljage/:id", server.getPoljeGEById)

	// router.GET("/interruptionofdelivery/:id", server.getDDNInterruptionOfDeliveryById)
	// router.GET("/interruptionofproduction", server.listDDNInterruptionOfDeliveryPByPage)
	// router.GET("/interruptionofproduction_all", server.listAllDDNInterruptionOfDeliveryP)
	// router.GET("/interruptionofproduction_excel", server.listExcelDDNInterruptionOfDeliveryP)
	// router.GET("/interruptionofusers", server.listDDNInterruptionOfDeliveryKByPage)
	// router.GET("/interruptionofusers_all", server.listAllDDNInterruptionOfDeliveryK)
	// router.GET("/interruptionofusers_excel", server.listExcelDDNInterruptionOfDeliveryK)

	// router.GET("/mesecni", server.listPiMM)
	// router.GET("/mesecnip", server.listPiMMByPage)
	// router.GET("/mesecnit4", server.listPiMMT4)
	// router.GET("/mesecnit4p", server.listPiMMT4ByPage)

	// router.GET("/dnevni", server.listPiDD)
	// router.GET("/dnevnip", server.listPiDDByPage)
	// router.GET("/dnevnit4p", server.listPiDDT4)

	// router.GET("/pogonski", server.listPiPI)
	// router.GET("/pogonskip", server.listPiPIByPage)
	// router.GET("/pogonskit4", server.listPiPIT4)

	// router.GET("/radapu_mes", server.listPGDRadapuMes)
	// router.GET("/dapua", server.listPGDDapuA)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.GET("/users/get_user_by_token", server.GetUserByToken)
	// authRoutes.POST("/tokens/renew_access", server.renewAccessToken)
	authRoutes.GET("/mrc/:id", server.getMrcById)
	authRoutes.GET("/mrc", server.listMrcs)
	authRoutes.GET("/mrciu", server.listMrcsForInsert)
	authRoutes.GET("/tipprek/:id", server.getSTipPrekById)
	authRoutes.GET("/tipprek", server.listTipPrek)
	authRoutes.GET("/vrprek/:id", server.getSVrPrekById)
	authRoutes.GET("/vrprek", server.listVrPrek)
	authRoutes.GET("/podvrprek", server.listPodVrPrek)
	authRoutes.GET("/uzrokprek/:id", server.getSUzrokPrekById)
	authRoutes.GET("/uzrokprek", server.listUzrokPrek)
	authRoutes.GET("/poduzrokprek/:id", server.getSPoduzrokPrekById)
	authRoutes.GET("/poduzrokprek", server.listPoduzrokPrek)
	authRoutes.GET("/mernamesta/:id", server.getSMernaMestaById)
	authRoutes.GET("/mernamesta", server.listMernaMesta)
	authRoutes.GET("/obj/:id", server.getObjId)
	authRoutes.GET("/objtsrp", server.listObjTSRP)
	authRoutes.GET("/objheteve", server.listObjHETEVE)
	authRoutes.GET("/poljage", server.listPoljaGE)
	authRoutes.GET("/poljage/:id", server.getPoljeGEById)

	authRoutes.GET("/interruptionofdelivery/:id", server.getDDNInterruptionOfDeliveryById)
	authRoutes.GET("/interruptionofproduction", server.listDDNInterruptionOfDeliveryPByPage)
	authRoutes.GET("/interruptionofproduction_all", server.listAllDDNInterruptionOfDeliveryP)
	authRoutes.GET("/interruptionofproduction_excel", server.listExcelDDNInterruptionOfDeliveryP)

	authRoutes.POST("/interruptionofdelivery", server.CreateDDNPrekidPr)
	authRoutes.PUT("/interruptionofdelivery/:id/:version", server.UpdateDDNPrekidPr)
	authRoutes.PUT("/interruptionofdelivery/:id/:version/bi", server.UpdateDDNPrekidIspBI)

	authRoutes.GET("/interruptionofusers", server.listDDNInterruptionOfDeliveryKByPage)
	authRoutes.GET("/interruptionofusers_all", server.listAllDDNInterruptionOfDeliveryK)
	authRoutes.GET("/interruptionofusers_excel", server.listExcelDDNInterruptionOfDeliveryK)

	authRoutes.GET("/mesecni", server.listPiMM)
	authRoutes.GET("/mesecnip", server.listPiMMByPage)
	authRoutes.GET("/mesecnit4", server.listPiMMT4)
	authRoutes.GET("/mesecnit4p", server.listPiMMT4ByPage)

	authRoutes.GET("/dnevni", server.listPiDD)
	authRoutes.GET("/dnevnip", server.listPiDDByPage)
	authRoutes.GET("/dnevnit4p", server.listPiDDT4)

	authRoutes.GET("/pogonski", server.listPiPI)
	authRoutes.GET("/pogonskip", server.listPiPIByPage)
	authRoutes.GET("/pogonskit4", server.listPiPIT4)

	authRoutes.GET("/radapu_mes", server.listPGDRadapuMes)
	authRoutes.GET("/dapua", server.listPGDDapuA)

	authRoutes.POST("/createinterruptionofproduction", server.CreateDDNPrekidPr)
	authRoutes.PUT("/interruptionofproduction//:id/:version", server.UpdateDDNPrekidPr)
	authRoutes.PUT("/interruptionofdelivery/bi/:id/:version", server.UpdateDDNPrekidIspBI)
	authRoutes.DELETE("/interruptionofdelivery/:id/:version", server.deleteDDNInterruptionOfDelivery)

	authRoutes.POST("/createinterruptionofusers", server.CreateDDNPrekidK)
	authRoutes.PUT("/interruptionofusers/:id/:version", server.UpdateDDNPrekidK)

	server.router = router
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
