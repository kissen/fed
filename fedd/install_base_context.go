package main

import (
	"context"
	"gitlab.cs.fau.de/kissen/fed/fedd/ap"
	"gitlab.cs.fau.de/kissen/fed/fedd/db"
	"net/http"
)

func InstallBaseContext(store db.FedStorage) func(http.Handler) http.Handler {
	config := Config()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := ap.WithFedContext(context.Background())
			fc := ap.FromContext(ctx)

			fc.Scheme = &config.Base.Scheme
			fc.Host = &config.Base.Host
			fc.BasePath = &config.Base.Path
			fc.Storage = store

			q := r.WithContext(ctx)
			next.ServeHTTP(w, q)
		})
	}
}
