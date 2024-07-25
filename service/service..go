package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"kedaplay"
	"log/slog"
	"net/http"
	"sync"
	"time"
	// "kedaplay/json"
)

type contextKey string

const (
	firstTask = iota
	lastTask
	allTasks
)

type Service struct {
	state      *kedaplay.State
	mux        *http.ServeMux
	mutex      sync.Mutex
	writeError func(w http.ResponseWriter, status int, err error)
	logger     *slog.Logger
}

var (
	ErrNoTasks = errors.New("no tasks")
)

type errorResponseFormat int

const (
	ErrorResponseJSON = iota
	ErrorResponseText
)

type Config struct {
	BaseUrl             string
	ErrorResponseFormat errorResponseFormat
}

func encodeJSON(w io.Writer, obj any) {
	b, _ := json.Marshal(obj)
	w.Write(b)
}

func NewService(cfg *Config, logger *slog.Logger) *Service {
	svc := &Service{
		state: &kedaplay.State{
			Tasks: []*kedaplay.Task{},
		},
		logger: logger.With(slog.String("package", "service")),
	}
	switch cfg.ErrorResponseFormat {
	case ErrorResponseJSON:
		svc.writeError = writeErrorJSON
	case ErrorResponseText:
		svc.writeError = writeErrorText
	}
	svc.routes(cfg.BaseUrl)
	return svc
}

type status struct {
	Count int              `json:"count"`
	Tasks []*kedaplay.Task `json:"tasks"`
}

type addRequest struct {
	kedaplay.Task
}

type errorResponse struct {
	Error string `json:"error"`
}

func writeErrorJSON(w http.ResponseWriter, status int, err error) {
	eo := errorResponse{
		Error: err.Error(),
	}
	b, _ := json.Marshal(eo)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(b)
}

func writeErrorText(w http.ResponseWriter, status int, err error) {
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
		s.logger.Info("handleAppendTask")

		// Get the payload from the body.
		b, err := io.ReadAll(r.Body)
		if err != nil {
			// Report error.
			return
		}

		var req addRequest
		err = json.Unmarshal(b, &req)
		if err != nil || len(req.Name) == 0 {
			// Report error.
			s.writeError(w, http.StatusBadRequest, fmt.Errorf("cannot unmarshal add task body: %w", err))
			return
		}
		// Add to the state.
		func() {
			s.mutex.Lock()
			defer s.mutex.Unlock()
			s.state.AddTask(&req.Task)
		}()
		// Return with success.
		w.WriteHeader(http.StatusOK)
	}
}

func (s *Service) handleRemoveTask(which int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		var t *kedaplay.Task

		s.logger.Info("handleRemoveTask")

		// Remove the task.
		switch which {
		case firstTask:
			t, err = s.removeFirstTask()
		default:
		}
		if err != nil {
			s.writeError(w, http.StatusNotFound, err)
			return
		}
		// Return with success.
		writeJSON(w, t)
	}
}

func (s *Service) removeFirstTask() (*kedaplay.Task, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	t := s.state.RemoveFirst()
	if t == nil {
		return nil, ErrNoTasks
	}
	return t, nil
}

func (s *Service) handleGetTasks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var logger = s.logger

		logger.Info("handleGetTasks")

		// Get the tasks.
		var tasks []*kedaplay.Task
		func() {
			s.mutex.Lock()
			defer s.mutex.Unlock()
			tasks = s.state.GetTasks()
		}()
		st := status{
			Count: len(tasks),
			Tasks: tasks,
		}
		writeJSON(w, st)
	}
}

func getPattern(method string, baseUrl string, path string) string {
	var w bytes.Buffer
	if method != "" {
		fmt.Fprintf(&w, "%s ", method)
	}
	fmt.Fprintf(&w, "%s%s", baseUrl, path)
	return w.String()
}

func (s *Service) routes(base string) {
	if s.mux != nil {
		return
	}
	s.mux = http.NewServeMux()
	s.mux.HandleFunc(getPattern("POST", base, "/tasks"), s.handleAppendTask())
	s.mux.HandleFunc(getPattern("GET", base, "/tasks"), s.handleGetTasks())
	s.mux.HandleFunc(getPattern("DELETE", base, "/tasks/first"), s.handleRemoveTask(firstTask))
	s.mux.HandleFunc(getPattern("DELETE", base, "/tasks/all"), s.handleRemoveTask(allTasks))
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
