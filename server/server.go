package server

import (
	"encoding/json"
	"io"
	"kedaplay"
	"net/http"
)

type Server struct {
	state *kedaplay.State
	mux *http.ServeMux
}

func encodeJSON(w io.Writer, obj any) {
	b, _ := json.Marshal(obj)
	w.Write(b)
}

type status struct {
	Count int             `json:"count"`
	Tasks []*kedaplay.Task `json:"tasks"`
}

func (s *Server) handleAppendTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the payload from the body.
		// Add to the state.
		// Return with success.
	}
}

func (s *Server) handleRemoveTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Remove the task.
		// Return error if no tasks was available.
		// Return with success.
	}
}

func (s *Server) handleGetTasks() http.HandlerFunc {
	
}

func (s *Server) routes() {
	if s.mux == nil {
		s.mux = http.NewServeMux()
	}
	s.mux.HandleFunc("POST /tasks", s.handleAppendTask())
	s.mux.HandleFunc("GET /tasks", s.)
}
