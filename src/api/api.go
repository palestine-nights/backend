package api

import (
	"database/sql"
)

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql" // TODO: Find why blank import.
	"github.com/jmoiron/sqlx"
)

// GenericError error model.
//
// swagger:model
type GenericError struct {
	// Error massage.
	Error string `json:"error"`
}

// Server is composition of router and DB instances.
// swagger:ignore
type Server struct {
	Router *gin.Engine
	DB     *sqlx.DB
	DBConn *sql.Conn
}

// GetServer returns server instance.
func GetServer(DB *sqlx.DB) *Server {
	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true

	config.AddAllowHeaders("Authorization")
	router.Use(cors.New(config))

	server := Server{Router: router, DB: DB}

	server.initializeRouter()

	return &server
}

// ListenAndServe server.
func (server *Server) ListenAndServe() {
	server.Router.Run()
}

func (server *Server) initializeRouter() {

	server.Router.StaticFile("/", "./html/home.html")

	tablesRouter := server.Router.Group("/tables")
	{
		tablesRouter.GET("", server.listTables)
		tablesRouter.GET("/:id", server.getTable)

		tablesRouter.Use(AuthMiddleware)
		{
			tablesRouter.POST("", server.postTable)
			tablesRouter.PUT("/:id", server.putTable)
			tablesRouter.DELETE("/:id", server.deleteTable)
		}
	}

	reservationsRouter := server.Router.Group("/reservations")
	{
		reservationsRouter.GET("", server.getReservations)
		reservationsRouter.GET("/:id", server.getReservation)
		reservationsRouter.POST("", server.postReservation)

		reservationsRouter.Use(AuthMiddleware)
		{
			reservationsRouter.POST("/approve/:id", server.approveReservation)
			reservationsRouter.POST("/cancel/:id", server.cancelReservation)
		}
	}

	menuRouter := server.Router.Group("/menu")
	{
		menuRouter.GET("", server.listMenu)
		menuRouter.GET("/:id", server.getMenuItem)

		menuRouter.Use(AuthMiddleware)
		{
			menuRouter.POST("", server.postMenuItem)
			menuRouter.PUT("/:id", server.putMenuItem)
			menuRouter.DELETE("/:id", server.deleteMenuItem)
		}
	}

	categoriesRouter := server.Router.Group("/categories")
	{
		categoriesRouter.GET("", server.getAllCategories)
		categoriesRouter.GET("/:category_id", server.listMenuItemsByCategory)

		categoriesRouter.Use(AuthMiddleware)
		{
			categoriesRouter.POST("", server.postCategory)
			categoriesRouter.PUT("/:id", server.updateCategory)
		}
	}
}
