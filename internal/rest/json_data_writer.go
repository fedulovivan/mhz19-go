package rest

import (
	"net/http"

	// "github.com/goccy/go-json"
	"encoding/json"
)

// same as JSONDataWriter from ozzo-routing/content/type.go
// but "encoding/json" is replaced with "github.com/goccy/go-json"

type JSONDataWriterCustom struct{}

func (w *JSONDataWriterCustom) SetHeader(res http.ResponseWriter) {
	res.Header().Set("Content-Type", "application/json")
}
func (w *JSONDataWriterCustom) Write(res http.ResponseWriter, data interface{}) (err error) {
	enc := json.NewEncoder(res)
	enc.SetEscapeHTML(false)
	return enc.Encode(data)
}
