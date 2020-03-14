package main

import (
	"context"
	"net/http"
)

const _FLASH_CONTEXT_KEY = "FlashContext"

type FlashContext struct {
	Flashs   []string
	Warnings []string
	Errors   []string
}

func AddFlashContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fc := &FlashContext{}
		c := context.WithValue(r.Context(), _FLASH_CONTEXT_KEY, fc)
		q := r.WithContext(c)
		next.ServeHTTP(w, q)
	})
}

func GetFlashContext(r *http.Request) *FlashContext {
	// on errors, we don't have flashes; in fact the middleware
	// isn't even called on errors; see
	// https://github.com/gorilla/mux/issues/416 for more details

	if fc := r.Context().Value(_FLASH_CONTEXT_KEY); fc == nil {
		return &FlashContext{}
	} else {
		return fc.(*FlashContext)
	}
}

func Flash(r *http.Request, s string) {
	fc := GetFlashContext(r)
	fc.Flashs = append(fc.Flashs, s)
}

func FlashWarning(r *http.Request, s string) {
	fc := GetFlashContext(r)
	fc.Warnings = append(fc.Warnings, s)
}

func FlashError(r *http.Request, s string) {
	fc := GetFlashContext(r)
	fc.Errors = append(fc.Errors, s)
}
