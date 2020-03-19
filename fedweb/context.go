package main

import (
	"context"
	"log"
	"gitlab.cs.fau.de/kissen/fed/fedweb/fedclient"
	"net/http"
)

const _FEDWEB_CONTEXT_KEY = "FedWebContext"

type FedWebContext struct {
	// The name of the tab that should be highlighted in the
	// navigation bar. If empty, not tab will be highlighted.
	Selected string

	// The title that should be used.
	Title string

	// Flashes to display on top of the page. Might be nil.
	Flashs []string

	// Warning (yellow) flashes to display on the top of the page.
	// Might be nil.
	Warnings []string

	// Error (red) flashes to display on the top of the page.
	// Might be nil.
	Errors []string

	// The HTTP status code to reply with. After being set the
	// first time with function Status(), it will not change anymore.
	// This means a handler can (1) set the status and then (2)
	// just call another handler to take care of the request w/o
	// changing the HTTP status code.
	//
	// Initialized to 200.
	Status int

	// Currently logged in user for this session. Might be nil.
	Client fedclient.FedClient
}

func AddContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Context().Value(_FEDWEB_CONTEXT_KEY) == nil {
			dbgClient, err := fedclient.New("http://localhost:9999/ap/alice")
			if err != nil {
				log.Println("creating client failed:", err)
				dbgClient = nil
			}

			fc := &FedWebContext{
				Status: http.StatusOK,
				Client: dbgClient,
			}
			c := context.WithValue(r.Context(), _FEDWEB_CONTEXT_KEY, fc)
			r = r.WithContext(c)
		}

		next.ServeHTTP(w, r)
	})
}

func Context(r *http.Request) *FedWebContext {
	return r.Context().Value(_FEDWEB_CONTEXT_KEY).(*FedWebContext)
}

func Selected(r *http.Request, tab string) {
	if fc := Context(r); len(fc.Selected) == 0 {
		fc.Selected = tab
	}
}

func Title(r *http.Request, title string) {
	if fc := Context(r); len(fc.Title) == 0 {
		fc.Title = title
	}
}

func Flash(r *http.Request, s string) {
	fc := Context(r)
	fc.Flashs = append(fc.Flashs, s)
}

func FlashWarning(r *http.Request, s string) {
	fc := Context(r)
	fc.Warnings = append(fc.Warnings, s)
}

func FlashError(r *http.Request, s string) {
	fc := Context(r)
	fc.Errors = append(fc.Errors, s)
}

func Status(r *http.Request, status int) {
	if fc := Context(r); fc.Status == 0 {
		fc.Status = status
	}
}
