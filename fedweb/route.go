package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"net/url"
)

const _ROUTER_CONTEXT_KEY = "RouterContext"

type RouterContext struct {
	Router *mux.Router
}

func AddRouterContext(router *mux.Router) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Context().Value(_ROUTER_CONTEXT_KEY) == nil {
				fc := &RouterContext{
					Router: router,
				}
				c := context.WithValue(r.Context(), _ROUTER_CONTEXT_KEY, fc)
				r = r.WithContext(c)
			}

			next.ServeHTTP(w, r)
		})
	}
}

func Router(r *http.Request) *mux.Router {
	return r.Context().Value(_ROUTER_CONTEXT_KEY).(*RouterContext).Router
}

func RouterHandler(r *http.Request, path string) (http.Handler, error) {
	// get the path as an URL

	log.Printf("path=%v", path)

	address, err := url.Parse(path)
	if err != nil {
		return nil, errors.Wrap(err, "bad address")
	}

	// construct a fake request we will pass to our router

	dummy := http.Request{
		Method: "GET",
		Host:   r.Host,
		URL:    address,
	}

	log.Printf("dummy=%v", dummy)

	var match mux.RouteMatch
	if ok := Router(r).Match(&dummy, &match); !ok {
		return nil, errors.New("no matching rule")
	}

	s, err := match.Route.GetPathTemplate()
	log.Printf("s=%v err=%v", s, err)

	// return values

	return match.Handler, nil
}
