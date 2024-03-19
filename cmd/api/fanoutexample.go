package main

import (
	"encoding/json"
	"net/http"

	"github.com/r3d5un/QueueToNull/internal/loggingutils"
)

type FanoutExamplePostBody struct {
	Message string `json:"message"`
}

type FanoutExampleResponse struct {
	Message string `json:"message"`
}

func (app *application) postFanoutExampleHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := loggingutils.LoggerFromContext(ctx)

	logger.Info("unmarshalling message")
	var b FanoutExamplePostBody
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		logger.Error("unable to decode HTTP request body", "body", r.Body)
		app.serverErrorResponse(w, r, err)
		return
	}
	logger.Info("decoded request body", "message", "b")

	logger.Info("publishing message")
	err = app.queues.FanoutExample.PublishMessage(b.Message)
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
		FanoutExampleResponse{Message: "message published"},
		nil,
	)
}
