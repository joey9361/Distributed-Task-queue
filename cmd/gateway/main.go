package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
)

type TaskRequest struct {

}

// ongoing api gateway that serves requests
func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	mux := http.NewServeMux()
	mux.HandleFunc("/task-request", handleTaskRequest)	
	// create redis client

	// create postgres connection pool

	// listen to the port for task requests
	if err := http.ListenAndServe(":8080", mux) ; err != nil {
		slog.LogAttrs(
			context.Background(),
			slog.LevelError,
			"Web server shut down unexpectedly",
			slog.Any("error", err),
			slog.String("port", "8080"),
		)
		os.Exit(1)
	}
	
	// 

}

func handleTaskRequest(w http.ResponseWriter, params *http.Request) {
	// validate the task request
	// validate the params
	// push to redis queue
	// set status in postgres
}

// func 