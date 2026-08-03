package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"time"
	"taskQueue/internal/database"
	"taskQueue/internal/queue"
)

type TaskRequest struct {
	Type 		string
	Payload		json.RawMessage
}

type Server struct {
	db		*database.Postgres
	queue	*queue.RedisQueue
}

func isJsonObject(message []byte) bool {
	// trim leading white space
	trimmed := bytes.TrimLeft(message, " \t\n\r")
	return trimmed[0] == '{'

}

// ongoing api gateway that serves requests
func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /tasks", handleTaskRequest)	
	mux.HandleFunc("GET /tasks/{id}", handleTaskResult)
	mux.HandleFunc("/health", handleHealth)
	// create redis client

	// create postgres connection pool
	// Create http Server instance
	server := &http.Server{
		Addr: ":8080",
		Handler: mux,
		ReadTimeout: (20 * time.Second),
		WriteTimeout: (20 * time.Second),
		IdleTimeout: (240 * time.Second),
	} 

	// listen to the port for task requests
	if err := server.ListenAndServe() ; err != nil {
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

func handleTaskRequest(w http.ResponseWriter, r *http.Request) {
	// validate the task request
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	var resp TaskRequest
	if err := dec.Decode(&resp); err != nil {
		slog.LogAttrs(
			r.Context(),
			slog.LevelWarn,
			"Failed to decode and unmarshal JSON request",
			slog.Any("error", err),
		)
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// validate the request aren't zero equivalent values
	if resp.Type == "" || len(resp.Payload) == 0 {
		slog.LogAttrs(
			r.Context(),
			slog.LevelWarn,
			"The request has a field with an empty value",
		)
		http.Error(w, "Empty field value", http.StatusBadRequest)
		return
	}

	// validate payload is JSON object
	if !isJsonObject(resp.Payload) {
		slog.LogAttrs(
			r.Context(),
			slog.LevelWarn,
			"Request payload is not a JSON object",
		)
		http.Error(w, "Request payload is not a JSON object", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

	w.Write([]byte(`{"Task_id": 123, "status": "Pending"}`))

	// validate the params
	// set status in postgres
	// push to redis queue
	// return status and task id to client
}

func handleTaskResult(w http.ResponseWriter, r *http.Request) {
	// validate the task id exists in postgres
	// validate the status of the task id in postgres
	// for a given task id, return the results stored in postgres if task completed
	// probably will recieve task id batch so return what is applicable
}

func handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	w.Write([]byte(`{"status": "UP"}`))
}