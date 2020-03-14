package main

import (
	"context"
	"net/http"
)

const _STATUS_CONTEXT_KEY = "StatusContext"
const _STATUS_UNSET = -1

type StatusContext struct {
	status int
}

func (sc StatusContext) Status() int {
	if sc.status == _STATUS_UNSET {
		return http.StatusOK
	} else {
		return sc.status
	}
}

func AddStatusContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Context().Value(_STATUS_CONTEXT_KEY) == nil {
			sc := &StatusContext{
				status: _STATUS_UNSET,
			}
			c := context.WithValue(r.Context(), _STATUS_CONTEXT_KEY, sc)
			r = r.WithContext(c)
		}

		next.ServeHTTP(w, r)
	})
}

func GetStatusContext(r *http.Request) *StatusContext {
	return r.Context().Value(_STATUS_CONTEXT_KEY).(*StatusContext)
}

func Status(r *http.Request, status int) {
	if sc := GetStatusContext(r); sc.status == _STATUS_UNSET {
		sc.status = status
	}
}
