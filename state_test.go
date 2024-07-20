package kedaplay

import (
	"reflect"
	"testing"
)

func TestState_RemoveFirst(t *testing.T) {
	type fields struct {
		Tasks []*Task
	}
	tests := []struct {
		name   string
		fields fields
		want   *Task
	}{
		// TODO: Add test cases.
		{
			name: "one",
			fields: fields{
				Tasks: []*Task{
					{
						Name:           "one",
						ProcessingTime: 1,
					},
				},
			},
			want: &Task{
				Name:           "one",
				ProcessingTime: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &State{
				Tasks: tt.fields.Tasks,
			}
			if got := s.RemoveFirst(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("State.RemoveFirst() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemaining(t *testing.T) {
	state := &State{
		Tasks: []*Task{
			{
				Name:           "one",
				ProcessingTime: 1,
			},
			{
				Name:           "two",
				ProcessingTime: 2,
			},
		},
	}
	s := state.RemoveFirst()
	t.Logf("%s", s.Name)

}
