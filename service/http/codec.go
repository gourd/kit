package httpservice

import (
	"encoding/json"
	"net/http"

	"golang.org/x/net/context"
)

type codecKeys int

const (
	decoderKey codecKeys = iota
	partialDecoderKey
)

// Decoder decodes streams it receives to the pointer v provided
type Decoder interface {
	Decode(v interface{}) (err error)
}

// ProvideJSONDecoder provides JSON encoder with a given *http.Request
// to context
func ProvideJSONDecoder(parent context.Context, r *http.Request) context.Context {
	if r == nil || r.Body == nil {
		return WithDecoder(parent, nil)
	}
	return WithDecoder(parent, json.NewDecoder(r.Body))
}

// WithDecoder adds a decoder to context so you can latter retrieve with
// DecoderFrom(context)
func WithDecoder(parent context.Context, decoder Decoder) context.Context {
	return context.WithValue(parent, decoderKey, decoder)
}

// DecoderFrom gets decoder set to the context
func DecoderFrom(ctx context.Context) (dec Decoder, ok bool) {
	dec, ok = ctx.Value(decoderKey).(Decoder)
	return
}

// WithPartialDecoder adds a decoder to context so you can latter retrieve with
// DecoderFrom(context)
func WithPartialDecoder(parent context.Context, decoder Decoder) context.Context {
	return context.WithValue(parent, decoderKey, decoder)
}

// PartialDecoderFrom gets decoder set to the context
func PartialDecoderFrom(ctx context.Context) (dec Decoder, ok bool) {
	if dec, ok = ctx.Value(partialDecoderKey).(Decoder); ok {
		return
	}
	return DecoderFrom(ctx)
}
