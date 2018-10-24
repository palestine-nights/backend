package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var a App

func TestMain(m *testing.M) {
	databaseUser := getEnv("DATABASE_USER", "tours_admin")
	databasePassword := getEnv("DATABASE_PASSWORD", "ladmdetouris")
	databaseName := getEnv("DATABASE_NAME", "restaurant")
	databaseHost := getEnv("DATABASE_HOST", "localhost")
	databasePort := getEnv("DATABASE_PORT", "3306")

	a = *GetApp(databaseUser, databasePassword, databaseName, databaseHost, databasePort)

	code := m.Run()
	os.Exit(code)
}

func TestGetNonExistentTable(t *testing.T) {
	req, _ := http.NewRequest("GET", "/table/1500", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Table with id 1500 could not be found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Table with id 1500 could not be found'. Got '%s'", m["error"])
	}
}

func TestGetTable1(t *testing.T) {
	req, _ := http.NewRequest("GET", "/table/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	var t1 table
	json.Unmarshal(response.Body.Bytes(), &t1)
	if t1.ID != 1 {
		t.Errorf("Expected 'id' to be 1. Got %d", t1.ID)
	}
	if t1.Places != 4 {
		t.Errorf("Expected 'places' to be 4. Got %d", t1.Places)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	return rr
}
func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
