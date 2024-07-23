package command

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"kedaplay"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type workerCmd struct {
	logger *slog.Logger
}

type WorkerOptions struct {
	Resource     string
	Sleep        int
}

var (
	ErrNotFound = errors.New("no tasks found")
)

func NewWorkerCmd() *workerCmd {
	return &workerCmd{}
}

func (c *workerCmd) processTask(ctx context.Context, opts *WorkerOptions) error {
	log.Print("starting process task")
	// req, err := http.NewRequestWithContext(ctx, http.MethodDelete, opts.Resource, nil)
	req, err := http.NewRequestWithContext(context.Background(), http.MethodDelete, opts.Resource, nil)
	if err != nil {
		return fmt.Errorf("error creating http request: %w", err)
	}
	c.logger.Info("issuing request")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error issuing request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return ErrNotFound
	}
	// Unmarshal the task.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading body: %w", err)
	}
	var task kedaplay.Task
	err = json.Unmarshal(body, &task)
	if err != nil {
		return fmt.Errorf("error unmarshalling ")
	}
	// wait for the processing time
	select {
	case <-time.After(time.Duration(task.ProcessingTime * int(time.Second))):
		log.Printf("processed task %s, time: %ds", task.Name, task.ProcessingTime)
	case <-ctx.Done():
		log.Printf("processed task %s canceled", task.Name)
	}

	return nil
}

func (c *workerCmd) processPendingTasks(ctx context.Context, opts *WorkerOptions) error {
	var err error
	for {
		err = c.processTask(ctx, opts)
		if err != nil {
			break
		}
	}
	if !errors.Is(err, ErrNotFound) {
		c.logger.Info("")
		return err
	}
	return nil
}

func (c *workerCmd) Run(ctx context.Context, args []string, logger *slog.Logger) error {
	resource := "http://localhost:8080/tasks/first"
	if r := os.Getenv("KDPW_RESOURCE"); r != "" {
		resource = r
	}
	
	log.Print("starting run")
	opts := &WorkerOptions{
		Resource: resource,
		Sleep:    5,
	}
	c.logger = logger.With(slog.String("component", "worker"))
	c.logger.Info("starting work", slog.String("resource", opts.Resource), slog.Int("sleep", opts.Sleep))

	for {
		err := c.processPendingTasks(ctx, opts)
		if err != nil {
			c.logger.Error("processPendingTask returned error", slog.String("error", err.Error()))
		}
		// All task processed or error.
		sd := time.Second * time.Duration(opts.Sleep)
		log.Printf("all done, sleeping for %s", sd.String())
		select {
		case <-ctx.Done():
			log.Printf("canceled, err: %s", ctx.Err().Error())
			return nil
		case <-time.After(sd):
			log.Printf("looking for new tasks")
		}
	}
}
