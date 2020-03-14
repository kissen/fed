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
		if r.Context().Value(_FLASH_CONTEXT_KEY) == nil {
			fc := &FlashContext{}
			c := context.WithValue(r.Context(), _FLASH_CONTEXT_KEY, fc)
			r = r.WithContext(c)
		}

		next.ServeHTTP(w, r)
	})
}

func GetFlashContext(r *http.Request) *FlashContext {
	return r.Context().Value(_FLASH_CONTEXT_KEY).(*FlashContext)
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
