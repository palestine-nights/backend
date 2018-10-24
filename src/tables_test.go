package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"
)

var a App

func TestMain(m *testing.M) {
	databaseUser := getEnv("DATABASE_USER", "tours_admin")
	databasePassword := getEnv("DATABASE_PASSWORD", "ladmdetouris")
	databaseName := getEnv("DATABASE_NAME", "restaurant")
	databaseHost := getEnv("DATABASE_HOST", "localhost")
	databasePort := getEnv("DATABASE_PORT", "3306")

	a = *GetApp(databaseUser, databasePassword, databaseName, databaseHost, databasePort)
	a.DB.DropTable(table{})
	a.DB.AutoMigrate(table{})
	rand.Seed(time.Now().UnixNano())

	code := m.Run()
	os.Exit(code)
}

// Create table and check it's existence
func TestCreateTable(t *testing.T) {
	// Create
	places := 1 + rand.Int() % 10
	desc := "description" + strconv.FormatInt(rand.Int63() % 100, 10)
	str := fmt.Sprintf("{\"places\":%d,\"description\":\"%s\"}", places, desc)
	payload := []byte(str)
	req, _ := http.NewRequest("POST", "/table", bytes.NewBuffer(payload))
	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)
	var t1 table
	json.Unmarshal(response.Body.Bytes(), &t1)
	if t1.Places != places {
		t.Errorf("Expected 'places' to be %d. Got %d", places, t1.Places)
	}
	if t1.Description != desc {
		t.Errorf("Expected 'description' to be %s. Got %s", desc, t1.Description)
	}

	// Check
	url := "/table/" + strconv.FormatInt(int64(t1.ID), 10)
	req, _ = http.NewRequest("GET", url, nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	var t2 table
	json.Unmarshal(response.Body.Bytes(), &t2)
	if t2.ID != t1.ID {
		t.Errorf("Expected 'id' to be %d. Got %d", t1.ID, t2.ID)
	}
	if t2.Places != t1.Places {
		t.Errorf("Expected 'places' to be %d. Got %d", t1.Places, t2.Places)
	}
	if t2.Description != t1.Description {
		t.Errorf("Expected 'description' to be %s. Got %s", t1.Description, t2.Description)
	}
}

func TestGetNonExistentTable(t *testing.T) {
	req, _ := http.NewRequest("GET", "/table/0", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Table with id 0 could not be found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Table with id 1500 could not be found'. Got '%s'", m["error"])
	}
}

// Incorrect json (string instead of int)
func TestCreateIncorrectFields(t *testing.T) {
	payload := []byte(`{"places":"string","description":"d"}`)
	req, _ := http.NewRequest("POST", "/table", bytes.NewBuffer(payload))
	response := executeRequest(req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)
	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Invalid request payload" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Invalid request payload'. Got '%s'", m["error"])
	}
}

// No 'places' field in request (also if places <= 0)
func TestCreateNoPlaces(t *testing.T) {
	payload := []byte(`{"description":"d"}`)
	req, _ := http.NewRequest("POST", "/table", bytes.NewBuffer(payload))
	response := executeRequest(req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)
	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Places count should be more than 0" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Places count should be more than 0'. Got '%s'", m["error"])
	}
}

// Create new table, even if id is set in payload
func TestCreateWithExistingId(t *testing.T) {
	payload := []byte(`{"id":1,"places":3,"description":"test creation"}`)
	req, _ := http.NewRequest("POST", "/table", bytes.NewBuffer(payload))
	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)
	var t1 table
	json.Unmarshal(response.Body.Bytes(), &t1)
	if t1.ID == 1 {
		t.Errorf("Expected 'id' to be not 1. Got %d", t1.ID)
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
