package main

import (
	"flag"
	"fmt"
	"net/http"

	"example.com/urlshort"
)

var fileName *string

func init() {
	fileName = flag.String("fname", "urlpaths.yml", "Yml File containing URL and Paths")
	flag.Parse()
}

func main() {
	mux := defaultMux()
	ymlFileContaint := urlshort.ReadYmlFile(*fileName)
	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/google": "https://google.com",
		"/yahoo":  "https://yahoo.com",
	}

	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yamlHandler, err := urlshort.YAMLHandler(ymlFileContaint, mapHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", yamlHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
