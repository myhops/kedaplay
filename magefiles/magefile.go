package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var Default = All

func All() {
	mg.Deps(BuildKedaplay, BuildWorker)
}

// BuildKedaplay builds a container image and pushes it to docker.io
func BuildKedaplay() error {
	env := map[string]string{
		"KO_DOCKER_REPO": "docker.io/peterzandbergen",
	}
	return sh.RunWith(env,
		"ko", "build", "./cmd/kedaplay")
}

// BuildWorker builds a container image and pushes it to docker.io
func BuildWorker() error {
	env := map[string]string{
		"KO_DOCKER_REPO": "docker.io/peterzandbergen",
	}
	return sh.RunWith(env,
		"ko", "build", "./cmd/worker")
}
