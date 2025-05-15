package main

import (
	"encoding/json"
	"flag"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"
)

type storyArcOptions struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

type storyArc struct {
	Title   string            `json:"title"`
	Story   []string          `json:"story"`
	Options []storyArcOptions `json:"options"`
}

func main() {
	jsonFilePath := flag.String("file", "gopher.json", "filepath to json file containing story arc")
	flag.Parse()
	storyArcsMap := parseJSON(jsonFilePath)

	tmpl, err := template.ParseFiles("story.html")
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimSpace(r.URL.Path)
		storyArcKey := string(path[1:])
		storyArc, ok := storyArcsMap[storyArcKey]
		if !ok {
			http.Redirect(w, r, "/intro", http.StatusFound)
		}
		tmpl.Execute(w, storyArc)
	})
	http.ListenAndServe(":8080", nil)
}

func parseJSON(jsonFilePath *string) map[string]storyArc {
	f, err := os.Open(*jsonFilePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	m := make(map[string]storyArc)
	dec := json.NewDecoder(f)
	if err := dec.Decode(&m); err != nil && err != io.EOF {
		panic(err)
	}

	return m
}
