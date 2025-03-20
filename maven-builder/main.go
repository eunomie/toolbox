package main

import (
	"context"
	"errors"

	"dagger/maven-builder/internal/dagger"
)

type MavenBuilder struct {
}

// Build the project using maven
func (m MavenBuilder) Build(ctx context.Context, dir *dagger.Directory) error {
	return m.mvn(ctx, dir, "package")
}

// Test the project using maven
func (m MavenBuilder) Test(ctx context.Context, dir *dagger.Directory) error {
	return m.mvn(ctx, dir, "test")
}

func (m MavenBuilder) mvn(ctx context.Context, dir *dagger.Directory, args ...string) error {
	check := dag.Java().Maven().WithSources(dir).
		Container().WithExec(append([]string{"mvn"}, args...), dagger.ContainerWithExecOpts{Expect: dagger.ReturnTypeAny})
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
