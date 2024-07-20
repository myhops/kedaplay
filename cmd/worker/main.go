package main

import (
	"context"
	"time"
)

type state struct {
	Tasks int
}

type repo interface {
	Get(context.Context) (*state, error)
	Put(context.Context, *state) error
}

func run(ctx context.Context, name string) error {
	// time.NewTimer()
	select {
	case <-ctx.Done():

	case <-time.After(time.Hour):

	}
	return nil
}

func main() {

}