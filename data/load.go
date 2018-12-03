package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/palestine-nights/backend/src/db"
)

type menuItems struct {
	Items []db.MenuItem `json:"menu"`
}

func main() {

	files, err := filepath.Glob("./data/menu/*.json")

	for _, file := range files {
		if err != nil {
			log.Fatal(err)
		}

		jsonFile, err := os.Open(file)

		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("Successfully Opened" + file)

		defer jsonFile.Close()

		byteValue, _ := ioutil.ReadAll(jsonFile)

		var menu menuItems

		json.Unmarshal(byteValue, &menu)

		requestURL := "https://api.palestinenights.com/menu"

		fmt.Print(menu.Items)

		for _, item := range menu.Items {
			j, _ := json.Marshal(item)
			body := bytes.NewBuffer(j)
			rsp, err := http.Post(requestURL, "application/x-www-form-urlencoded", body)

			if err != nil {
				panic(err)
			}

			defer rsp.Body.Close()

			_, err = ioutil.ReadAll(rsp.Body)

			if err != nil {
				panic(err)
			}

			fmt.Println("--- Done ---")
		}
	}
}
