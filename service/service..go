package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"kedaplay"
	"log/slog"
	"net/http"
	"time"

	// "kedaplay/json"
)

type contextKey string

// SLogger
var SLoggerContextKey contextKey = contextKey("sLoggerKey")

func getLogger(ctx context.Context) *slog.Logger {
	sl, ok := ctx.Value(SLoggerContextKey).(*slog.Logger)
	if !ok {
		return nil
	}
	return sl
}

const (
	firstTask = iota
	lastTask
	allTasks
)

type Service struct {
	state *kedaplay.State
	mux   *http.ServeMux
}

var (
	ErrNoTasks = errors.New("no tasks")
)

func encodeJSON(w io.Writer, obj any) {
	b, _ := json.Marshal(obj)
	w.Write(b)
}

func NewService() *Service {
	svc := &Service{
		state: &kedaplay.State{
			Tasks: []*kedaplay.Task{},
		},
	}
	svc.routes("")
	return svc
}

type status struct {
	Count int              `json:"count"`
	Tasks []*kedaplay.Task `json:"tasks"`
}

type addRequest struct {
	kedaplay.Task
}

func writeError(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(status)
	fmt.Fprintf(w, "Error: %s", err.Error())
}

func writeJSON(w http.ResponseWriter, obj any) {
	b, _ := json.Marshal(obj)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (s *Service) handleAppendTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var logger = getLogger(r.Context())
		logger.Info("handleAppendTask")

		// Get the payload from the body.
		defer r.Body.Close()
		b, err := io.ReadAll(r.Body)
		if err != nil {
			// Report error.
			return
		}

		var req addRequest
		err = json.Unmarshal(b, &req)
		if err != nil || len(req.Name) == 0 {
			// Report error.
			writeError(w, http.StatusBadRequest, fmt.Errorf("cannot unmarshal add task body: %w", err))
			return
		}
		// Add to the state.
		s.state.AddTask(&req.Task)
		// Return with success.
		w.WriteHeader(http.StatusOK)
	}
}

func (s *Service) handleRemoveTask(which int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		var t *kedaplay.Task

		var logger = getLogger(r.Context())
		logger.Info("handleRemoveTask")

		// Remove the task.
		switch which {
		case firstTask:
			t, err = s.removeFirstTask()
		default:
		}
		// Return with success.
		if err != nil {
			writeError(w, http.StatusNotFound, err)
			return
		}
		writeJSON(w, t)
	}
}

func (s *Service) removeFirstTask() (*kedaplay.Task, error) {
	t := s.state.RemoveFirst()
	if t == nil {
		return nil, ErrNoTasks
	}
	return t, nil
}

func (s *Service) handleGetTasks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var logger = getLogger(r.Context())
		logger.Info("handleGetTasks")

		// Get the tasks.
		st := status{
			Count: len(s.state.Tasks),
			Tasks: s.state.Tasks,
		}
		writeJSON(w, st)
	}
}

func (s *Service) routes(base string) {
	if s.mux != nil {
		return
	}
	s.mux = http.NewServeMux()
	s.mux.HandleFunc("POST /tasks", s.handleAppendTask())
	s.mux.HandleFunc("GET /tasks", s.handleGetTasks())
	s.mux.HandleFunc("DELETE /tasks/first", s.handleRemoveTask(firstTask))
	s.mux.HandleFunc("DELETE /tasks/all", s.handleRemoveTask(allTasks))
	s.mux.HandleFunc("/*", http.NotFound)
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func LogRequestMiddleware(next http.HandlerFunc, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.String()
		defer func(start time.Time) {
			passed := time.Since(start)
			logger.Info("handler called", slog.String("url", path), slog.Duration("duration", passed))
		}(time.Now())

		next(w, r)
	}
}
