package kedaplay

import "slices"

type State struct {
	Tasks []*Task `json:"tasks"`
}

type Task struct {
	Name           string `json:"name"`
	ProcessingTime int    `json:"processingTime"`
}

func (s *State) AddTask(task *Task) error {
	s.Tasks = append(s.Tasks, task)
	return nil
}

func (s *State) RemoveFirst() *Task {
	l := len(s.Tasks)
	if l == 0 {
		return nil
	}
	removed := s.Tasks[0]
	s.Tasks = slices.Delete(s.Tasks, 0, 1)
	return removed
}

func (s *State) Append(t *Task) {
	s.Tasks	= append(s.Tasks, t)
}

