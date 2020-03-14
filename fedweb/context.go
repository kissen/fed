package main

import (
	"context"
	"net/http"
)

const _FEDWEB_CONTEXT_KEY = "FedWebContext"

type FedWebContext struct {
	Flashs   []string
	Warnings []string
	Errors   []string
	Status   int
}

func AddContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Context().Value(_FEDWEB_CONTEXT_KEY) == nil {
			fc := &FedWebContext{
				Status: http.StatusOK,
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
