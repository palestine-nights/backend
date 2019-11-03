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

/// swagger:route GET /tables tables listTables
/// List all tables.
/// Responses:
///   200: []Table
///   500: GenericError
func (server *Server) listTables(c *gin.Context) {
	table, err := db.Table.GetAll(db.Table{}, server.DB)

	if err == nil {
		c.JSON(http.StatusOK, table)
	} else {
		c.JSON(http.StatusInternalServerError, GenericError{Error: err.Error()})
	}
}

/// swagger:route GET /tables/{id} tables getTable
/// Returns table.
/// Responses:
///   200: Table
///   400: GenericError
///   404: GenericError
func (server *Server) getTable(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, GenericError{Error: "Invalid table ID, must be int"})
		return
	}

	table, err := db.Table.Find(db.Table{}, server.DB, id)

	if err == nil {
		c.JSON(http.StatusOK, table)
	} else {
		errorMsg := fmt.Sprintf("Table with id %d could not be found", id)
		c.JSON(http.StatusNotFound, GenericError{Error: errorMsg})
	}
}

/// swagger:route POST /tables tables postTable
/// Creates table.
/// Responses:
///   200: Table
///   400: GenericError
func (server *Server) postTable(c *gin.Context) {
	table := db.Table{}

	if err := c.ShouldBindJSON(&table); err != nil {
		c.JSON(http.StatusBadRequest, GenericError{Error: "Invalid request payload"})
		return
	}

	if table.Places <= 0 {
		c.JSON(http.StatusBadRequest, GenericError{Error: "Places count should be more than 0"})
		return
	}

	err := table.Insert(server.DB)

	if err == nil {
		c.JSON(http.StatusCreated, table)
	} else {
		c.JSON(http.StatusBadRequest, GenericError{Error: err.Error()})
	}
}

/// swagger:route PUT /tables/{id} tables putTable
/// Updates table.
/// Responses:
///   200: Table
///   400: GenericError
///   404: GenericError
func (server *Server) putTable(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, GenericError{Error: "Invalid table ID, must be int"})
		return
	}

	table := db.Table{}

	if err := c.ShouldBindJSON(&table); err != nil {
		c.JSON(http.StatusBadRequest, GenericError{Error: "Invalid request payload"})
		return
	}

	if table.Places <= 0 {
		c.JSON(http.StatusBadRequest, GenericError{Error: "Places count should be more than 0"})
		return
	}

	table.ID = id

	// Check if ID exists.
	err = table.Update(server.DB)
	if err != nil {
		c.JSON(http.StatusNotFound, err.Error())
	} else {
		c.JSON(http.StatusOK, table)
	}
}

// swagger:route DELETE /tables/{id} tables deleteTable
// Deletes table.
// Responses:
//   204:
//   400: GenericError
func (server *Server) deleteTable(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, GenericError{Error: "Invalid table ID, must be int"})
		return
	}

	err = db.Table.Destroy(db.Table{}, server.DB, id)

	// Check if ID exists.
	if err == nil {
		c.JSON(http.StatusNoContent, nil)
	} else {
		c.JSON(http.StatusNotFound, err.Error())
	}
}
