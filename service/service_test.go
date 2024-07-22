package service

import (
	"bytes"
	"encoding/json"
	"io"
	"kedaplay"
	"net/http"
	"net/http/httptest"
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
	t.Error(b.String())
}

func TestService(t *testing.T) {
	svc := NewService()
	r := httptest.NewRequest("GET", "/tasks", nil)
	w := httptest.NewRecorder()
	svc.mux.ServeHTTP(w, r)
	body, _ := io.ReadAll(w.Body)
	t.Logf("%s", string(body))
	t.Errorf("%s", w.Result().Status)
}

func TestAddTask(t *testing.T) {
	var r *http.Request
	var w *httptest.ResponseRecorder
	var body []byte
	var task *kedaplay.Task

	svc := NewService()
	// Get the tasks.
	r = httptest.NewRequest("GET", "/tasks", nil)
	w = httptest.NewRecorder()
	svc.mux.ServeHTTP(w, r)
	body, _ = io.ReadAll(w.Body)
	t.Errorf("Get tasks: %s", w.Result().Status)
	t.Logf("%s", string(body))

	task = &kedaplay.Task{
		Name: "first task",
		ProcessingTime: 10,
	}
	body, _ = json.Marshal(task)
	bodyR := bytes.NewReader(body)
	r = httptest.NewRequest("POST", "/tasks", bodyR)
	w = httptest.NewRecorder()
	svc.mux.ServeHTTP(w, r)
	t.Errorf("Add task: status: %d", w.Code)
	body, _ = io.ReadAll(w.Body)
	t.Logf("%s", string(body))

	// Get the tasks.
	r = httptest.NewRequest("GET", "/tasks", nil)
	w = httptest.NewRecorder()
	svc.mux.ServeHTTP(w, r)
	body, _ = io.ReadAll(w.Body)
	t.Errorf("Get tasks: %s", w.Result().Status)
	t.Logf("%s", string(body))

	// Delete the task
	r = httptest.NewRequest("DELETE", "/tasks/first", nil)
	w = httptest.NewRecorder()
	svc.mux.ServeHTTP(w, r)
	body, _ = io.ReadAll(w.Body)
	t.Errorf("delete task: %s", w.Result().Status)
	t.Logf("%s", string(body))
	
	// Get the tasks.
	r = httptest.NewRequest("GET", "/tasks", nil)
	w = httptest.NewRecorder()
	svc.mux.ServeHTTP(w, r)
	body, _ = io.ReadAll(w.Body)
	t.Errorf("Get tasks: %s", w.Result().Status)
	t.Logf("%s", string(body))
	
	// Delete the task
	r = httptest.NewRequest("DELETE", "/tasks/first", nil)
	w = httptest.NewRecorder()
	svc.mux.ServeHTTP(w, r)
	body, _ = io.ReadAll(w.Body)
	t.Errorf("delete task: %s", w.Result().Status)
	t.Logf("%s", string(body))

	// add bad task
	body = []byte("no json")
	bodyR = bytes.NewReader(body)
	r = httptest.NewRequest("POST", "/tasks", bodyR)
	w = httptest.NewRecorder()
	svc.mux.ServeHTTP(w, r)
	t.Errorf("Add task: status: %d", w.Code)
	body, _ = io.ReadAll(w.Body)
	t.Logf("%s", string(body))

}

