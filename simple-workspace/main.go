package main

import (
	"context"

	"dagger/simple-workspace/internal/dagger"
)

type SimpleWorkspace struct {
	Sources *dagger.Directory
}

func New(sources *dagger.Directory) SimpleWorkspace {
	return SimpleWorkspace{
		Sources: sources,
	}
}

// Read a file
func (w SimpleWorkspace) Read(ctx context.Context, path string) (string, error) {
	return w.Sources.File(path).Contents(ctx)
}

// Write a file
func (w SimpleWorkspace) Write(path, content string) SimpleWorkspace {
	w.Sources = w.Sources.WithNewFile(path, content)
	return w
}

// Copy an entire directory into the workspace
func (w SimpleWorkspace) CopyDir(
	// The target path
	path string,
	// The directory to copy at the target path.
	// Existing content is overwritten at the file granularity.
	dir *dagger.Directory,
) SimpleWorkspace {
	w.Sources = w.Sources.WithDirectory(path, dir)
	return w
}

// Remove a file from the workspace
func (w SimpleWorkspace) Rm(path string) SimpleWorkspace {
	w.Sources = w.Sources.WithoutFile(path)
	return w
}

// Remove a directory from the workspace
func (w SimpleWorkspace) RmDir(path string) SimpleWorkspace {
	w.Sources = w.Sources.WithoutDirectory(path)
	return w
}

// List the contents of a directory in the workspace
func (w SimpleWorkspace) ListDir(
	ctx context.Context,
	// Path of the target directory
	// +optional
	// +default="/"
	path string,
) ([]string, error) {
	return w.Sources.Directory(path).Entries(ctx)
}

// Walk all files in the workspace (optionally filtered by a glob pattern), and return their path.
func (w SimpleWorkspace) Walk(
	ctx context.Context,
	// A glob pattern to filter files. Only matching files will be included.
	// The glob format is the same as Dockerfile/buildkit
	// +optional
	// +default="**"
	pattern string,
) ([]string, error) {
	return w.Sources.Glob(ctx, pattern)
}
