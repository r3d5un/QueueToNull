package main

import (
	"net/http"

	"github.com/r3d5un/QueueToNull/internal/loggingutils"
)

type HelloWorldResponse struct {
	Message string `json:"message"`
}

func (app *application) postHelloWorldMessageHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := loggingutils.LoggerFromContext(ctx)

	logger.Info("publishing message")
	err := app.queues.HelloWorldQueue.PublishHelloWorld()
	if err != nil {
		logger.Error("unable to publish message", "error", err)
		app.serverErrorResponse(w, r, err)
		return
	}
	logger.Info("message published")

	logger.Info("writing response")
	err = app.writeJSON(
		w,
		http.StatusOK,
		HelloWorldResponse{Message: "'Hello, World!' posted"},
		nil,
	)
	if err != nil {
		logger.Error("unable to write response", "error", err)
		app.serverErrorResponse(w, r, err)
		return
	}
}
