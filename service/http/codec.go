package httpservice

import (
	"encoding/json"
	"net/http"

	"github.com/gourd/kit/context"

	"golang.org/x/net/context"
)

type codecKeys int

const (
	decoderKey codecKeys = iota
	partialDecoderKey
)

// DecoderProvider provides decoder
type DecoderProvider func(r *http.Request) Decoder

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
	return WithDecoder(parent, func(r *http.Request) Decoder {
		return json.NewDecoder(r.Body)
	})
}

// WithDecoder adds a decoder to context so you can latter retrieve with
// DecoderFrom(context)
func WithDecoder(parent context.Context, provider DecoderProvider) context.Context {
	return context.WithValue(parent, decoderKey, provider)
}

// DecoderFrom gets decoder set to the context
func DecoderFrom(ctx context.Context) (dec Decoder, ok bool) {
	var pro DecoderProvider
	r := gourdctx.HTTPRequest(ctx)
	if r == nil {
		ok = false
		return
	}
	if pro, ok = ctx.Value(decoderKey).(DecoderProvider); ok {
		dec = pro(r)
		return
	}
	return
}

// WithPartialDecoder adds a decoder to context so you can latter retrieve with
// DecoderFrom(context)
func WithPartialDecoder(parent context.Context, provider DecoderProvider) context.Context {
	return context.WithValue(parent, decoderKey, provider)
}

// PartialDecoderFrom gets decoder set to the context
func PartialDecoderFrom(ctx context.Context) (dec Decoder, ok bool) {
	var pro DecoderProvider
	r := gourdctx.HTTPRequest(ctx)
	if r == nil {
		ok = false
		return
	}
	if pro, ok = ctx.Value(partialDecoderKey).(DecoderProvider); ok {
		dec = pro(r)
		return
	}
	return DecoderFrom(ctx)
}
