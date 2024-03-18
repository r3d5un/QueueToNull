package main

import (
	"encoding/json"
	"net/http"

	"github.com/r3d5un/QueueToNull/internal/loggingutils"
)

type ExampleWorkQueuePostBody struct {
	Message string `json:"message"`
}

type ExampleWorkQueueResponse struct {
	Message string `json:"message"`
}

func (app *application) postExampleWorkQueueHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := loggingutils.LoggerFromContext(ctx)

	logger.Info("unmarshalling message")
	var b ExampleWorkQueuePostBody
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		logger.Error("unable to decode HTTP request body", "body", r.Body)
		app.serverErrorResponse(w, r, err)
		return
	}
	logger.Info("decoded request body", "message", "b")

	logger.Info("publishing message")
	err = app.queues.ExampleWorkQueue.PublishMessage(b.Message)
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
		ExampleWorkQueueResponse{Message: "message published"},
		nil,
	)
}
