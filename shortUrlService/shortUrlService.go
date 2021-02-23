package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type UrlMapping struct {
	Paths map[string]string `json:"paths"`
}

func process(mapping *UrlMapping) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		path := req.URL.Path

		longPath, found := mapping.Paths[path]
		if found {
			http.Redirect(w, req, longPath, http.StatusSeeOther)
		} else {
			_, err := w.Write([]byte("Not Found"))
			if err != nil {
				fmt.Printf("Response writer error: %s", err)
			}
		}
	}
}

func getUrlMapping(path string) (*UrlMapping, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Json file open error: %s", err))
	}

	defer func() {
		err := jsonFile.Close()
		if err != nil {
			fmt.Println("Close config file error")
		}
	}()

	byteValue, err := ioutil.ReadAll(jsonFile)

	var urlMapping UrlMapping
	err = json.Unmarshal(byteValue, &urlMapping)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Json parse error: %s", err))
	}

	return &urlMapping, nil
}

func getConfigFilePath() string {
	pathPtr := flag.String("f", "data/urlMap.json", "a config file path")
	flag.Parse()
	return *pathPtr
}

func main() {
	urlMapping, err := getUrlMapping(getConfigFilePath())
	if err != nil {
		fmt.Printf("Url mapping error: %s", err)
		return
	}

	http.HandleFunc("/", process(urlMapping))
	fmt.Println("Starting server at http://localhost:8080")

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Start server error: %s", err)
		return
	}
}
