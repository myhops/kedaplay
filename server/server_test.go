package server

import (
	"bytes"
	"kedaplay"
	"testing"
)

func TestEncode(t *testing.T) {
	srv := &Server{
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
