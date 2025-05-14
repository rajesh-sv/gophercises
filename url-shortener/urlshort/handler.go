package urlshort

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"gopkg.in/yaml.v3"
)

type pathUrl struct {
	Path string `yaml:"path" json:"path"`
	URL  string `yaml:"url" json:"url"`
}

func buildMap(pathUrls []pathUrl) map[string]string {
	pathToUrls := make(map[string]string)
	for _, pathUrl := range pathUrls {
		pathToUrls[pathUrl.Path] = pathUrl.URL
	}
	return pathToUrls
}

func parseYAML(yamlData []byte) ([]pathUrl, error) {
	var pathUrls []pathUrl
	err := yaml.Unmarshal(yamlData, &pathUrls)
	if err != nil {
		return nil, err
	}
	return pathUrls, nil
}

func parseJSON(jsonData []byte) ([]pathUrl, error) {
	var pathUrls []pathUrl
	err := json.Unmarshal(jsonData, &pathUrls)
	if err != nil {
		return nil, err
	}
	return pathUrls, nil
}

func SqliteHandler(db *sql.DB, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		row := db.QueryRow("SELECT * FROM mappings WHERE path=?", path)
		var pathUrl pathUrl
		if err := row.Scan(&pathUrl.Path, &pathUrl.URL); err == nil {
			http.Redirect(w, r, pathUrl.URL, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if url, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, url, http.StatusFound)
		}
		fallback.ServeHTTP(w, r)
	}
}

func YAMLHandler(yamlData []byte, fallback http.Handler) (http.HandlerFunc, error) {
	pathUrls, err := parseYAML(yamlData)
	if err != nil {
		return nil, err
	}

	pathToUrls := buildMap(pathUrls)
	return MapHandler(pathToUrls, fallback), nil
}

func JSONHandler(jsonData []byte, fallback http.Handler) (http.HandlerFunc, error) {
	pathUrls, err := parseJSON(jsonData)
	if err != nil {
		return nil, err
	}

	pathToUrls := buildMap(pathUrls)
	return MapHandler(pathToUrls, fallback), nil
}
