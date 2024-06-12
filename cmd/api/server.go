package api

import (
	"fmt"

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
func NewServer( config util.Config, store db.DatabaseRepo) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
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

	
	router.POST("/users/login", server.loginUser)
	router.POST("/tokens/renew_access", server.renewAccessToken)
	

	//authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	
	//authRoutes.GET("/accounts/:id", server.getAccount)
	

	server.router = router
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}