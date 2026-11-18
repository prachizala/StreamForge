package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/prachizala/StreamForge/services/gateway/internal/queue"
)

var dispatcher = queue.NewDispatcher(4)

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"service":   "streamforge-gateway",
		"status":    "healthy",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"runtime":   "Go 1.22+ Concurrent Worker Pool",
	})
}

func handleJobs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodPost {
		var payload struct {
			Filename string `json:"filename"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || payload.Filename == "" {
			payload.Filename = "uploaded_media.mp4"
		}

		job := dispatcher.CreateJob(payload.Filename)
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(job)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}

func handleJobStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := strings.TrimPrefix(r.URL.Path, "/api/jobs/")
	if id == "" {
		http.Error(w, `{"error":"missing job id"}`, http.StatusBadRequest)
		return
	}

	job, exists := dispatcher.GetJob(id)
	if !exists {
		http.Error(w, `{"error":"job not found"}`, http.StatusNotFound)
		return
	}

	_ = json.NewEncoder(w).Encode(job)
}

func handleJobStream(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/stream/")
	if id == "" {
		http.Error(w, "missing job id", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	ticker := time.NewTicker(40 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-r.Context().Done():
			return
		case <-ticker.C:
			job, exists := dispatcher.GetJob(id)
			if !exists {
				return
			}

			data, _ := json.Marshal(job)
			_, _ = fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()

			if job.Status == queue.StatusCompleted || job.Status == queue.StatusFailed {
				return
			}
		}
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", corsMiddleware(handleHealth))
	mux.HandleFunc("/api/jobs", corsMiddleware(handleJobs))
	mux.HandleFunc("/api/jobs/", corsMiddleware(handleJobStatus))
	mux.HandleFunc("/api/stream/", corsMiddleware(handleJobStream))

	port := ":8080"
	log.Printf("StreamForge Go Gateway listening on %s...", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("Server exited with error: %v", err)
	}
}
