package httputil

import (
	"encoding/json"
	"io"
)

func StructFromBodyJSON[T any](body io.ReadCloser) (obj T, _ error) {
	decoder := json.NewDecoder(body)
	return obj, decoder.Decode(&obj)
}
