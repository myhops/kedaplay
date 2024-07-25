package service

import (
	"bytes"
	"encoding/json"
	"io"
	"kedaplay"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestEncode(t *testing.T) {
	srv := &Service{
		state: &kedaplay.State{
			Tasks: []*kedaplay.Task{
				{
					Name: "one",
				},
			},
		},
	}
	s := &status{
		Count: 2,
		Tasks: srv.state.Tasks,
	}
	b := &bytes.Buffer{}
	encodeJSON(b, s)
	t.Log(b.String())
}

func TestService(t *testing.T) {
	svcCfg := Config{}
	svc := NewService(&svcCfg, slog.Default())
	r := httptest.NewRequest("GET", "/tasks", nil)
	w := httptest.NewRecorder()
	svc.mux.ServeHTTP(w, r)
	body, _ := io.ReadAll(w.Body)
	t.Logf("%s", string(body))
	t.Logf("%s", w.Result().Status)
}

func TestAddTask(t *testing.T) {
	var r *http.Request
	var w *httptest.ResponseRecorder
	var body []byte
	var task *kedaplay.Task

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	svc := NewService(&Config{}, logger)
	// Get the tasks.
	r = httptest.NewRequest("GET", "/tasks", nil)
	w = httptest.NewRecorder()
	svc.mux.ServeHTTP(w, r)
	body, _ = io.ReadAll(w.Body)
	t.Logf("Get tasks: %s", w.Result().Status)
	t.Logf("%s", string(body))

	task = &kedaplay.Task{
		Name:           "first task",
		ProcessingTime: 10,
	}
	body, _ = json.Marshal(task)
	bodyR := bytes.NewReader(body)
	r = httptest.NewRequest("POST", "/tasks", bodyR)
	w = httptest.NewRecorder()
	svc.mux.ServeHTTP(w, r)
	t.Logf("Add task: status: %d", w.Code)
	body, _ = io.ReadAll(w.Body)
	t.Logf("%s", string(body))

	// Get the tasks.
	r = httptest.NewRequest("GET", "/tasks", nil)
	w = httptest.NewRecorder()
	svc.mux.ServeHTTP(w, r)
	body, _ = io.ReadAll(w.Body)
	t.Logf("Get tasks: %s", w.Result().Status)
	t.Logf("%s", string(body))

	// Delete the task
	r = httptest.NewRequest("DELETE", "/tasks/first", nil)
	w = httptest.NewRecorder()
	svc.mux.ServeHTTP(w, r)
	body, _ = io.ReadAll(w.Body)
	t.Logf("delete task: %s", w.Result().Status)
	t.Logf("%s", string(body))

	// Get the tasks.
	r = httptest.NewRequest("GET", "/tasks", nil)
	w = httptest.NewRecorder()
	svc.mux.ServeHTTP(w, r)
	body, _ = io.ReadAll(w.Body)
	t.Logf("Get tasks: %s", w.Result().Status)
	t.Logf("%s", string(body))

	// Delete the task
	r = httptest.NewRequest("DELETE", "/tasks/first", nil)
	w = httptest.NewRecorder()
	svc.mux.ServeHTTP(w, r)
	body, _ = io.ReadAll(w.Body)
	t.Logf("delete task: %s", w.Result().Status)
	t.Logf("%s", string(body))

	// add bad task
	body = []byte("no json")
	bodyR = bytes.NewReader(body)
	r = httptest.NewRequest("POST", "/tasks", bodyR)
	w = httptest.NewRecorder()
	svc.mux.ServeHTTP(w, r)
	t.Logf("Add task: status: %d", w.Code)
	body, _ = io.ReadAll(w.Body)
	t.Logf("%s", string(body))

}
