package urlshort

import (
	"fmt"
	"net/http"

	"encoding/json"
	"gopkg.in/yaml.v3"
)

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(pathsToUrls)
		if path, ok := pathsToUrls[r.URL.Path]; ok {
			http.Redirect(w, r, path, http.StatusFound)
		}
		fallback.ServeHTTP(w, r)
	}
}

func YAMLHandler(yaml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYAML(yaml)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedYaml)
	return MapHandler(pathMap, fallback), nil
}

func JSONHandler(json []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedJson, err := parseJSON(json)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedJson)
	return MapHandler(pathMap, fallback), nil
}

type itemPath struct {
	URL  string
	Path string
}

func buildMap(items []itemPath) map[string]string {
	pathsToUrls := make(map[string]string, len(items))
	for _, item := range items {
		pathsToUrls[item.Path] = item.URL
	}
	return pathsToUrls
}

func parseYAML(data []byte) ([]itemPath, error) {
	var paths []itemPath
	err := yaml.Unmarshal(data, &paths)
	if err != nil {
		return nil, err
	}
	return paths, nil
}

func parseJSON(data []byte) ([]itemPath, error) {
	var paths []itemPath
	err := json.Unmarshal(data, &paths)
	if err != nil {
		return nil, err
	}
	return paths, nil
}
