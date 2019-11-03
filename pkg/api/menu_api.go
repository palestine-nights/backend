package api

import (
	"fmt"
	"net/http"
	"strconv"
)

import (
	"github.com/gin-gonic/gin"
	"github.com/palestine-nights/backend/pkg/db"
)

/// swagger:route GET /menu menu listMenu
/// List all menu items.
/// Responses:
///   200: []MenuItem
///   500: GenericError
func (server *Server) listMenu(c *gin.Context) {
	menu, err := db.MenuItem.GetAll(db.MenuItem{}, server.DB)

	if err == nil {
		c.JSON(http.StatusOK, menu)
	} else {
		c.JSON(http.StatusInternalServerError, GenericError{Error: err.Error()})
	}
}

/// swagger:route GET /menu/{id} menu getMenuItem
/// Returns menu item.
/// Responses:
///   200: MenuItem
///   404: GenericError
func (server *Server) getMenuItem(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusNotFound, GenericError{Error: "Invalid menu item ID, must be int"})
		return
	}

	menuItem, err := db.MenuItem.Find(db.MenuItem{}, server.DB, id)

	if err == nil {
		c.JSON(http.StatusOK, menuItem)
	} else {
		errorMsg := fmt.Sprintf("Menu item with id %d could not be found", id)
		c.JSON(http.StatusNotFound, GenericError{Error: errorMsg})
	}
}

/// swagger:route POST /menu menu postMenuItem
/// Create menu item.
/// Responses:
///   201: MenuItem
///   400: GenericError
func (server *Server) postMenuItem(c *gin.Context) {
	menuItem := db.MenuItem{}

	if err := c.BindJSON(&menuItem); err != nil {
		c.JSON(http.StatusBadRequest, GenericError{Error: "Invalid request payload"})
		return
	}

	// Validations
	if len(menuItem.Name) == 0 {
		c.JSON(http.StatusBadRequest, GenericError{Error: "Name should not be empty"})
		return
	}

	if menuItem.Price <= 0 {
		c.JSON(http.StatusBadRequest, GenericError{Error: "Price should be more than 0"})
		return
	}

	err := menuItem.Insert(server.DB)

	if err == nil {
		c.JSON(http.StatusCreated, menuItem)
	} else {
		c.JSON(http.StatusBadRequest, GenericError{Error: err.Error()})
	}
}

/// swagger:route PUT /menu/{id} menu putMenuItem
/// Update menu item.
/// Responses:
///   200: MenuItem
///   400: GenericError
func (server *Server) putMenuItem(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, GenericError{Error: "Invalid menu item ID, must be int"})
		return
	}

	menuItem := db.MenuItem{}

	if err := c.ShouldBindJSON(&menuItem); err != nil {
		c.JSON(http.StatusBadRequest, GenericError{Error: "Invalid request payload"})
		return
	}

	// Validations
	if len(menuItem.Name) == 0 {
		c.JSON(http.StatusBadRequest, GenericError{Error: "Name should not be empty"})
		return
	}

	if menuItem.Price <= 0 {
		c.JSON(http.StatusBadRequest, GenericError{Error: "Price should be more than 0"})
		return
	}

	menuItem.ID = id

	// Check if ID exists.
	err = menuItem.Update(server.DB)
	if err == nil {
		c.JSON(http.StatusOK, menuItem)
	} else {
		c.JSON(http.StatusNotFound, GenericError{Error: err.Error()})
	}
}

/// swagger:route DELETE /menu/{id} menu deleteMenuItem
/// Delete menu item.
/// Responses:
///   204:
///   404: GenericError
func (server *Server) deleteMenuItem(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, GenericError{Error: "Invalid menu item ID, must be int"})
		return
	}

	err = db.MenuItem.Destroy(db.MenuItem{}, server.DB, id)

	// Check if ID exists.
	if err == nil {
		c.JSON(http.StatusNoContent, nil)
	} else {
		c.JSON(http.StatusNotFound, GenericError{Error: err.Error()})
	}
}
