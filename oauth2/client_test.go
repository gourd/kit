package oauth2_test

import (
	"github.com/gourd/kit/oauth2"

	"math/rand"
	"testing"

	"github.com/RangelReale/osin"
)

func dummyNewClient(redirectUri string) *oauth2.Client {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	randSeq := func(n int) string {
		b := make([]rune, n)
		for i := range b {
			b[i] = letters[rand.Intn(len(letters))]
		}
		return string(b)
	}

	return &oauth2.Client{
		ID:          randSeq(10),
		Secret:      randSeq(10),
		RedirectURI: redirectUri,
		UserID:      "",
	}
}

func TestClient(t *testing.T) {
	var c osin.Client = &oauth2.Client{}
	_ = c
}
