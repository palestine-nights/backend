package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/palestine-nights/backend/src/db"
)

/*  Menu API */

func (server *Server) listMenu(w http.ResponseWriter, r *http.Request) {
	menu, err := db.MenuItem.GetAll(db.MenuItem{}, server.DB)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, menu)
	}
}

func (server *Server) listMenuByCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	category := vars["category"]

	menu, err := db.MenuItem.GetByCategory(db.MenuItem{}, server.DB, category)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, menu)
	}
}

func (server *Server) getMenuItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.ParseUint(vars["id"], 10, 64)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid menu item ID, must be int")
		return
	}

	menuItem, err := db.MenuItem.Find(db.MenuItem{}, server.DB, id)

	if err != nil {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("Menu item with id %d could not be found", id))
	} else {
		respondWithJSON(w, http.StatusOK, menuItem)
	}
}

func (server *Server) postMenuItem(w http.ResponseWriter, r *http.Request) {
	menuItem := db.MenuItem{}
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&menuItem); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validations
	if len(menuItem.Name) == 0 {
		respondWithError(w, http.StatusBadRequest, "Name should not be empty")
		return
	}
	if len(menuItem.Category) == 0 {
		respondWithError(w, http.StatusBadRequest, "Category should not be empty")
		return
	}
	if menuItem.Price <= 0 {
		respondWithError(w, http.StatusBadRequest, "Price should be more than 0")
		return
	}

	defer r.Body.Close()

	err := menuItem.Insert(server.DB)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	} else {
		respondWithJSON(w, http.StatusCreated, menuItem)
	}
}

func (server *Server) putMenuItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.ParseUint(vars["id"], 10, 64)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid menu item ID, must be int")
		return
	}

	var menuItem db.MenuItem

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&menuItem); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validations
	if len(menuItem.Name) == 0 {
		respondWithError(w, http.StatusBadRequest, "Name should not be empty")
		return
	}
	if len(menuItem.Category) == 0 {
		respondWithError(w, http.StatusBadRequest, "Category should not be empty")
		return
	}
	if menuItem.Price <= 0 {
		respondWithError(w, http.StatusBadRequest, "Price should be more than 0")
		return
	}
	defer r.Body.Close()

	menuItem.ID = id

	// Check if ID exists.
	err = menuItem.Update(server.DB)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, menuItem)
	}
}

func (server *Server) deleteMenuItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid menu item ID, must be int")
		return
	}
	err = db.MenuItem.Destroy(db.MenuItem{}, server.DB, id)

	// Check if ID exists.
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
	} else {
		respondWithJSON(w, http.StatusNoContent, map[string]interface{}{})
	}
}

func (server *Server) getAllCategories(w http.ResponseWriter, r *http.Request) {

	categories, err := db.MenuItem.GetCategories(db.MenuItem{}, server.DB)

	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, categories)
	}
}
