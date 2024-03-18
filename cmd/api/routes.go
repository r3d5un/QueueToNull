package main

import (
	"net/http"

	"github.com/justinas/alice"
)

type RouteDefinition struct {
	Path    string
	Handler http.HandlerFunc
}

type RouteDefinitionList []RouteDefinition

func (app *application) routes() http.Handler {
	app.logger.Info("creating multiplexer")
	mux := http.NewServeMux()

	standard := alice.New(app.recoverPanic)

	app.logger.Info("adding healthcheck route")
	mux.HandleFunc("GET /api/v1/healthcheck", app.healthcheckHandler)

	handlerList := RouteDefinitionList{
		{"POST /api/v1/queue/hello_world", app.postHelloWorldMessageHandler},
		{"POST /api/v1/queue/example_work_queue", app.postExampleWorkQueueHandler},
	}

	app.logger.Info("adding routes")
	for _, d := range handlerList {
		app.logger.Info("adding route", "route", d.Path)
		mux.Handle(d.Path, standard.ThenFunc(d.Handler))
	}

	app.logger.Info("creating final handler")
	handler := standard.Then(mux)

	return handler
}
