package urlshort

import (
	"net/http"

	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).cd
// If the path is not provided in the map, then the fallback
// http.Handler will be callinstead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if url, ok := pathsToUrls[r.URL.Path]; ok {
			http.Redirect(rw, r, url, http.StatusFound)
			return
		}
		fallback.ServeHTTP(rw, r)
	}
}

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
func MapUrlPath(data []ymlDt) map[string]string {
	urlStrings := make(map[string]string)
	for _, val := range data {
		urlStrings[val.Path] = val.Url
	}
	return urlStrings
}

func parseYml(ymlData []byte) ([]ymlDt, error) {
	yml_data := []ymlDt{}

	err := yaml.Unmarshal(ymlData, &yml_data)

	return yml_data, err
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parse_data, err := parseYml(yml)
	urlStrings := MapUrlPath(parse_data)
	if err != nil {
		return nil, err
	}
	return MapHandler(urlStrings, fallback), nil
}

type ymlDt struct {
	Path string `yaml:"path"`
	Url  string `yaml:"url"`
}
