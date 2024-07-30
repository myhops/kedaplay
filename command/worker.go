package command

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/myhops/kedaplay"
)

type workerCmd struct {
	logger   *slog.Logger
	resource string
	sleep    int
}

type WorkerConfig struct {
	Resource string
	Sleep    int
	Logger   *slog.Logger
}

var (
	ErrNotFound = errors.New("no tasks found")
)

func NewWorkerCmd(cfg *WorkerConfig) *workerCmd {
	cmd := &workerCmd{}
	cmd.resource = cfg.Resource
	cmd.logger = cfg.Logger.With(slog.String("resource", cfg.Resource))
	cmd.sleep = cfg.Sleep
	return cmd
}

func (c *workerCmd) getTask() (*kedaplay.Task, error) {
	c.logger.Info("starting process task")
	// req, err := http.NewRequestWithContext(ctx, http.MethodDelete, opts.Resource, nil)
	req, err := http.NewRequestWithContext(context.Background(), http.MethodDelete, c.resource, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating http request: %w", err)
	}
	c.logger.Info("issuing request")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error issuing request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}
	// Unmarshal the task.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading body: %w", err)
	}
	var task kedaplay.Task
	err = json.Unmarshal(body, &task)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling ")
	}
	return &task, nil
}

func (c *workerCmd) processTask(ctx context.Context) error {
	task, err := c.getTask()
	if err != nil {
		return err
	}
	// wait for the processing time
	c.logger.Info("starting task", slog.String("task", task.Name), slog.Int("processingTime", task.ProcessingTime))
	select {
	case <-time.After(time.Duration(task.ProcessingTime * int(time.Second))):
		log.Printf("processed task %s, time: %ds", task.Name, task.ProcessingTime)
	case <-ctx.Done():
		log.Printf("processed task %s canceled", task.Name)
	}

	return nil
}

func (c *workerCmd) processPendingTasks(ctx context.Context) error {
	var err error
	for {
		err = c.processTask(ctx)
		if err != nil {
			break
		}
	}
	if !errors.Is(err, ErrNotFound) {
		c.logger.Info("error processing task", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (c *workerCmd) Run(ctx context.Context) error {
	c.logger.Info("starting work", slog.String("resource", c.resource), slog.Int("sleep", c.sleep))

	for {
		err := c.processPendingTasks(ctx)
		if err != nil {
			c.logger.Error("processPendingTask returned error", slog.String("error", err.Error()))
		}
		// All task processed or error.
		sd := time.Second * time.Duration(c.sleep)
		c.logger.Info("all done", slog.Duration("sleep_duration", sd))
		select {
		case <-ctx.Done():
			c.logger.Info("canceled", slog.String("error", ctx.Err().Error()))
			return nil
		case <-time.After(sd):
			c.logger.Info("looking for new tasks")
		}
	}
}
