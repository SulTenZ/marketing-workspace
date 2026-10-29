// Package handler adalah adapter Vercel; tidak membuat listener atau worker loop.
package handler

import (
	"net/http"
	"os"
	"sync"

	"unggulmart/marketing-workspace/apps/api/internal/app"
)

var runtimeHandler = sync.OnceValue(func() http.Handler {
	return app.FromEnvironment(os.LookupEnv, os.Stdout)
})

func Handler(w http.ResponseWriter, r *http.Request) {
	runtimeHandler().ServeHTTP(w, r)
}
