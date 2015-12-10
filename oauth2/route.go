package oauth2

import (
	"github.com/gorilla/pat"
	httpservice "github.com/gourd/kit/service/http"

	"log"
)

// RoutePat adds manager's endpoint to a pat router
// TODO: remove this function and pat dependency
func RoutePat(rtr *pat.Router, base string, ep *Endpoints) {

	// tell user where the endpoints are
	log.Printf("OAuth2 authorize endpoint: %s/authorize", base)
	log.Printf("OAuth2 token endpoint:     %s/token", base)

	// bind handler with pat
	rtr.Get(base+"/authorize", ep.Auth)
	rtr.Post(base+"/authorize", ep.Auth)
	rtr.Get(base+"/token", ep.Token)
	rtr.Post(base+"/token", ep.Token)
	rtr.Get(base+"/info", ep.Info)

}

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
