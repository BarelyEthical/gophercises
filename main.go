package main

import (
	"fmt"
	"net/http"
	"github.com/gophercises/urlShortner/urlshort"
	"log"
	"os"
	"bytes"
)

func getJsonFile(fileName string) []byte {
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal("could not open file %s", fileName)
	}
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(f)
	if err != nil {
		log.Fatal("could not read file %s", fileName)
	}
	return buf.Bytes()
}

func main() {
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	/*
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}*/
	jsonFile := getJsonFile("paths.json")
	jsonData, _ := urlshort.ParseJson(jsonFile)
	jHandler := urlshort.JsonHandler(jsonData, mux)
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", jHandler)
	/*
	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yaml := `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`
	yamlHandler, err := urlshort.YAMLHandler([]byte(yaml), mapHandler)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", yamlHandler)
	*/
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
