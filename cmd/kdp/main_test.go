package main

import "testing"

func TestServer(t *testing.T) {
	args := []string{
		"kedaplay", "server", "--logformat=text", "--port=10000",
	}
	run(args)
}

func TestWorker(t *testing.T) {
	args := []string{
		"kedaplay", "worker",
	}
	run(args)
}
