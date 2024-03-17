package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	app.logger.Info("creating multiplexer")
	mux := http.NewServeMux()

	standard := alice.New(app.recoverPanic)

	app.logger.Info("adding healthcheck route")
	mux.HandleFunc("GET /api/v1/healthcheck", app.healthcheckHandler)

	app.logger.Info("creating final handler")
	handler := standard.Then(mux)

	return handler
}
