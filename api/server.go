package api

import (
	"github.com/gin-gonic/gin"

	db "github.com/broemp/broempSignal/db/sqlc"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

// NewServer creates a new HTTP server and setups routing
func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/users", server.createUser)
	router.GET("/users/:id", server.getUser)
	router.GET("/users/", server.listUser)

	// AFK endpoint
	router.POST("/afk/create", server.createAFK)
	router.GET("/afk/list/:id", server.listAFK)
	router.GET("/afk/count/:id", server.countAFK)

	server.router = router
	return server
}

// Start running the HTTP server on address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
