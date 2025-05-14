package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/rajesh-sv/gophercises/url-shortener/urlshort"
	_ "modernc.org/sqlite"
)

func main() {
	mux := defaultMux()

	// Build the SqliteHandler using mux as the fallback
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	dsnURI := filepath.Join(dir, "url.db")
	db, err := sql.Open("sqlite", dsnURI)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	fmt.Println("Connected to sqlite database!")
	sqliteHandler := urlshort.SqliteHandler(db, mux)

	// Build the MapHandler using the sqliteHandler as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, sqliteHandler)

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

	// Build the JSONHandler using the yamlHandler as the
	// fallback
	json := `
	[{"path":"/neovim","url":"https://neovim.io/"},{"path":"/ghostty","url":"https://ghostty.org/"}]
	`
	jsonHandler, err := urlshort.JSONHandler([]byte(json), yamlHandler)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", jsonHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
