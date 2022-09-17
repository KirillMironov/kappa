package httputil

import (
	"gopkg.in/yaml.v2"
	"io"
)

func StructFromBodyYAML[T any](body io.ReadCloser) (obj T, _ error) {
	decoder := yaml.NewDecoder(body)
	return obj, decoder.Decode(&obj)
}
