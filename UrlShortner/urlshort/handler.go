package urlshort

import (
	"net/http"
	"encoding/json"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.


type pathToURL struct {
	Path string
	URL string
}

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := pathsToUrls[r.URL.Path]
		if url != "" {
			http.Redirect(w, r, url, http.StatusPermanentRedirect)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

func buildMap(yamlData []pathToURL) (map[string]string) {
	builtMap := make(map[string]string)
	for _, j := range yamlData {
		builtMap[j.Path] = j.URL
	}
	return builtMap
}

func ParseJson(jsonData []byte)(pathsToURLS []pathToURL, err error) {
	err = json.Unmarshal(jsonData, &pathsToURLS)
	return
}
func JsonHandler(pathsToURLS []pathToURL, fallback http.Handler) http.HandlerFunc {
	dataMap := buildMap(pathsToURLS)
	jsonHandler := MapHandler(dataMap, fallback)
	return jsonHandler
}

/*
// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.

func parseYAML(yml []byte) (pathsToURLS []pathToURL, err error) {
	err = yml.Unmarshal(yamlData, &pathToURLS)
	return
}

func YAMLHandler(yml []byte, fallback http.Handler) (yamlHandler http.HandlerFunc, error) {
	yamlData, err := parseYAML(yml)
	if err != nil {
		return
	}
	pathMap := buildMap(yamlData)
	yamlHandler = MapHandler(pathMap, fallback)
	return
}
*/
