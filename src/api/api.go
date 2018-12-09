package api

import (
	"database/sql"
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"        //
	_ "github.com/jinzhu/gorm/dialects/mysql" //
	"github.com/jmoiron/sqlx"
	"github.com/palestine-nights/backend/src/db"
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
func GetServer(user, password, database, host, port string) *Server {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, password, host, port, database)

	DB := db.Initialize(connectionString)

	router := gin.Default()
	router.Use(cors.Default())
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
		tablesRouter.POST("", server.postTable)
		tablesRouter.GET("", server.listTables)
		tablesRouter.GET("/:id", server.getTable)
		tablesRouter.PUT("/:id", server.putTable)
		tablesRouter.DELETE("/:id", server.deleteTable)
	}

	reservationsRouter := server.Router.Group("/reservations")
	{
		reservationsRouter.POST("", server.postReservation)
		reservationsRouter.GET("", server.getReservations)
		reservationsRouter.GET("/:id", server.getReservation)
		reservationsRouter.POST("/approve/:id", server.approveReservation)
		reservationsRouter.POST("/cancel/:id", server.cancelReservation)
	}

	menuRouter := server.Router.Group("/menu")
	{
		menuRouter.POST("", server.postMenuItem)
		menuRouter.GET("", server.listMenu)
		menuRouter.GET("/:id", server.getMenuItem)
		menuRouter.PUT("/:id", server.putMenuItem)
		menuRouter.DELETE("/:id", server.deleteMenuItem)
	}

	categoriesRouter := server.Router.Group("/categories")
	{
		categoriesRouter.GET("", server.getAllCategories)
		categoriesRouter.PUT("/:id", server.updateCategory)
		categoriesRouter.GET("/:category_id", server.listMenuItemsByCategory)
	}
}
