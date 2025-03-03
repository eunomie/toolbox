package main

import (
	"context"
	"errors"

	"dagger/java-workspace/internal/dagger"
)

type JavaWorkspace struct {
	// +private
	Src *dagger.Directory
}

func New(src *dagger.Directory) JavaWorkspace {
	return JavaWorkspace{
		Src: src,
	}
}

// Read a file
func (w JavaWorkspace) Read(ctx context.Context, path string) (string, error) {
	return w.Src.File(path).Contents(ctx)
}

// Write a file
func (w JavaWorkspace) Write(path, content string) JavaWorkspace {
	w.Src = w.Src.WithNewFile(path, content)
	return w
}

// Copy an entire directory into the workspace
func (w JavaWorkspace) CopyDir(
	// The target path
	path string,
	// The directory to copy at the target path.
	// Existing content is overwritten at the file granularity.
	dir *dagger.Directory,
) JavaWorkspace {
	w.Src = w.Src.WithDirectory(path, dir)
	return w
}

// Remove a file from the workspace
func (w JavaWorkspace) Rm(path string) JavaWorkspace {
	w.Src = w.Src.WithoutFile(path)
	return w
}

// Remove a directory from the workspace
func (w JavaWorkspace) RmDir(path string) JavaWorkspace {
	w.Src = w.Src.WithoutDirectory(path)
	return w
}

// List the contents of a directory in the workspace
func (w JavaWorkspace) ListDir(
	ctx context.Context,
	// Path of the target directory
	// +optional
	// +default="/"
	path string,
) ([]string, error) {
	return w.Src.Directory(path).Entries(ctx)
}

// Walk all files in the workspace (optionally filtered by a glob pattern), and return their path.
func (w JavaWorkspace) Walk(
	ctx context.Context,
	// A glob pattern to filter files. Only matching files will be included.
	// The glob format is the same as Dockerfile/buildkit
	// +optional
	// +default="**"
	pattern string,
) ([]string, error) {
	return w.Src.Glob(ctx, pattern)
}

// Build the code at the current directory in the workspace
func (w JavaWorkspace) Build(ctx context.Context) error {
	return w.mvn(ctx, "package")
}

// Run tests
func (w JavaWorkspace) Test(ctx context.Context) error {
	return w.mvn(ctx, "test")
}

func (w JavaWorkspace) mvn(ctx context.Context, args ...string) error {
	check := dag.Java().Maven().WithSources(w.Src).Container().WithExec(append([]string{"mvn"}, args...), dagger.ContainerWithExecOpts{Expect: dagger.ReturnTypeAny})
	code, err := check.ExitCode(ctx)
	if err != nil {
		return err
	}
	if code == 0 {
		return nil
	}
	stderr, err := check.Stderr(ctx)
	if err != nil {
		return err
	}
	return errors.New(stderr)
}

func (w JavaWorkspace) Dir() *dagger.Directory {
	return w.Src
}

func base() *dagger.Container {
	digest := "sha256:56fa17d2a7e7f168a043a2712e63aed1f8543aeafdcee47c58dcffe38ed51099"
	return dag.
		Container().
		From("docker.io/library/alpine:latest@" + digest)
}
