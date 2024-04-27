// A generated module for Ci functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

package main

import (
	"context"
)

type Ci struct{}

// GIT_COMMIT=$(shell git rev-parse --verify HEAD)
// UTC_NOW=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
//
// build-dev:
// 	go build \
// 		-ldflags="-X 'main.version=dev' -X 'main.commit=${GIT_COMMIT}' -X 'main.date=${UTC_NOW}'" \
// 		-o cdebug

func (m *Ci) Build(ctx context.Context) *File {
	return dag.Container().
		From("golang:1.22-alpine").
		WithDirectory("/app", dag.CurrentModule().Source().Directory("..")).
		WithWorkdir("/app").
		WithExec([]string{"go", "build", "-o", "cdebug"}).
		File("cdebug")
}

// Runs the e2e tests for the project.
func (m *Ci) TestE2e(ctx context.Context, docker *Service) error {
	cdebug := m.Build(ctx)

	container := m.testBase(ctx).
		WithDirectory("/app", dag.CurrentModule().Source().Directory("..")).
		WithFile("/usr/local/bin/cdebug", cdebug).
		WithWorkdir("/app").
		WithServiceBinding("docker", docker).
		WithEnvVariable("DOCKER_HOST", "tcp://docker:2375").
		WithExec([]string{"go", "test", "-v", "-count", "1", "./e2e/exec"})

	_, err := container.Sync(ctx)
	return err
}

func (m *Ci) testBase(ctx context.Context) *Container {
	return dag.Container().
		From("golang:1.22-alpine").
		WithExec([]string{"apk", "add", "--no-cache", "sudo", "docker"})
}
