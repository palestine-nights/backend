package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/palestine-nights/backend/src/db"
)

/// swagger:route PUT /categories/{id} menu updateCategory
/// Update menu category.
/// Responses:
///   200: MenuCategory
///   400: GenericError
func (server *Server) updateCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, GenericError{Error: "Invalid menu category ID, must be int"})
		return
	}

	category := db.MenuCategory{}

	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, GenericError{Error: "Invalid request payload"})
		return
	}

	category.ID = id

	// Check if ID exists.
	err = category.Update(server.DB)

	if err == nil {
		c.JSON(http.StatusOK, category)
	} else {
		c.JSON(http.StatusNotFound, GenericError{Error: err.Error()})
	}
}

/// swagger:route GET /menu/categories/{category_id} menu listMenuByCategory
/// List menu items with specified category.
/// Responses:
///   200: []MenuItem
///   500: GenericError
func (server *Server) listMenuItemsByCategory(c *gin.Context) {
	categoryID, err := strconv.ParseUint(c.Param("category_id"), 10, 64)

	menu, err := db.MenuItem.GetByCategory(db.MenuItem{}, server.DB, categoryID)

	if err == nil {
		c.JSON(http.StatusOK, menu)
	} else {
		c.JSON(http.StatusInternalServerError, GenericError{Error: err.Error()})
	}
}
