package httpservice

import (
	"encoding/json"
	"net/http"

	"golang.org/x/net/context"
)

type codecKeys int

const decoderKey codecKeys = 0

// Decoder decodes streams it receives to the pointer v provided
type Decoder interface {
	Decode(v interface{}) (err error)
}

// ProvideJSONDecoder provides JSON encoder with a given *http.Request
// to context
func ProvideJSONDecoder(parent context.Context, r *http.Request) context.Context {
	if r == nil || r.Body == nil {
		return context.WithValue(parent, decoderKey, nil)
	}
	return context.WithValue(parent, decoderKey, json.NewDecoder(r.Body))
}

// DecoderFrom gets decoder set to the context
func DecoderFrom(ctx context.Context) (dec Decoder, ok bool) {
	dec, ok = ctx.Value(decoderKey).(Decoder)
	return
}
