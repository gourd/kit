package oauth2

import (
	httpservice "github.com/gourd/kit/service/http"

	"log"
)

// Route adds manager's endpoint to a router with httpservice.RouterFunc
func Route(rfn httpservice.RouterFunc, base string, ep *Endpoints) (err error) {

	// tell user where the endpoints are
	log.Printf("OAuth2 authorize endpoint:   %s/authorize", base)
	log.Printf("OAuth2 token endpoint:       %s/token", base)
	log.Printf("OAuth2 information endpoint: %s/info", base)

	// bind handler with router function
	if err = rfn(base+"/authorize", []string{"GET", "POST"}, ep.Auth); err != nil {
		return
	}
	if err = rfn(base+"/token", []string{"GET", "POST"}, ep.Token); err != nil {
		return
	}
	if err = rfn(base+"/info", []string{"GET"}, ep.Info); err != nil {
		return
	}
	return
}
