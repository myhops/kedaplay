package main

import "testing"

func TestServer(t *testing.T) {
	args := []string{
		"kedaplay", "server",
	}
	run(args)
}

func TestWorker(t *testing.T) {
	args := []string{
		"kedaplay", "worker",
	}
	run(args)
}
