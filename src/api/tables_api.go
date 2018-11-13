package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/palestine-nights/backend/src/db"
)

/*  Tables API */

func (server *Server) listTables(w http.ResponseWriter, r *http.Request) {
	table, err := db.Table.GetAll(db.Table{}, server.DB)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, table)
	}
}

func (server *Server) getTable(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.ParseUint(vars["id"], 10, 64)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid table ID, must be int")
		return
	}

	table, err := db.Table.Find(db.Table{}, server.DB, id)

	if err != nil {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("Table with id %d could not be found", id))
	} else {
		respondWithJSON(w, http.StatusOK, table)
	}
}

func (server *Server) postTable(w http.ResponseWriter, r *http.Request) {
	var table db.Table
	decoder := json.NewDecoder(r.Body)

	fmt.Print(r.Body)

	if err := decoder.Decode(&table); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Validations:
	if table.Places <= 0 {
		respondWithError(w, http.StatusBadRequest, "Places count should be more than 0")
		return
	}

	defer r.Body.Close()

	err := table.Insert(server.DB)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	} else {
		respondWithJSON(w, http.StatusCreated, table)
	}
}

func (server *Server) putTable(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.ParseUint(vars["id"], 10, 64)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid table ID, must be int")
		return
	}

	var table db.Table

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&table); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if table.Places <= 0 {
		respondWithError(w, http.StatusBadRequest, "Places count should be more than 0")
		return
	}
	defer r.Body.Close()

	table.ID = id

	// Check if ID exists.
	err = table.Update(server.DB)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, table)
	}
}

func (server *Server) deleteTable(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid table ID, must be int")
		return
	}
	err = db.Table.Destroy(db.Table{}, server.DB, id)

	// Check if ID exists.
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
	} else {
		respondWithJSON(w, http.StatusNoContent, map[string]interface{}{})
	}
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
