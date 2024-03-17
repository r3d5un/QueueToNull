package main

import (
	"log/slog"
	"net/http"

	"github.com/r3d5un/QueueToNull/internal/loggingutils"
)

type HealthCheckMessage struct {
	Status string `json:"status"`
}

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := loggingutils.LoggerFromContext(ctx)

	healthCheckMessage := HealthCheckMessage{
		Status: "available",
	}

	logger.InfoContext(ctx, "writing response", "response", healthCheckMessage)
	err := app.writeJSON(w, http.StatusOK, healthCheckMessage, nil)
	if err != nil {
		slog.ErrorContext(ctx, "error writing response", "error", err)
		app.serverErrorResponse(w, r, err)
		return
	}
}
