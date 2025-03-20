package main

import (
	"context"

	"dagger/toolbox/internal/dagger"
)

type Toolbox struct {
	// +private
	Src *dagger.Directory
}

// Explain what the code does and why it does it in the most concise way possible
func (m *Toolbox) Explain(ctx context.Context) (string, error) {
	llm, err := m.asReader(ctx, "explain what the code does and why it does it in the most concise way possible")
	if err != nil {
		return "", err
	}
	return m.printLLMLastReply(ctx, llm)
}

// Find potential bugs in the existing code, explain them and propose alternative code to fix them
func (m *Toolbox) FindBugs(ctx context.Context) (string, error) {
	llm, err := m.asReader(ctx, "find potential bugs in the existing code, explain them and propose alternative code to fix them")
	if err != nil {
		return "", err
	}
	return m.printLLMLastReply(ctx, llm)
}

// Ask to implement a new feature in the existing code
func (m *Toolbox) Do(ctx context.Context, ask string) (*dagger.Directory, error) {
	llm, err := m.asEditor(ctx, ask)
	if err != nil {
		return nil, err
	}
	return llm.GetSimpleWorkspace("sources").Sources(), nil
}

// Add comments to the existing code to improve readability and maintainability
func (m *Toolbox) AddComments(ctx context.Context) (*dagger.Directory, error) {
	llm, err := m.asEditor(ctx, "add comments to the existing code to improve readability and maintainability")
	if err != nil {
		return nil, err
	}
	return llm.GetSimpleWorkspace("sources").Sources(), nil
}

// Refactor the existing code to improve readability and maintainability
func (m *Toolbox) Refactor(ctx context.Context) (*dagger.Directory, error) {
	llm, err := m.asEditor(ctx, "refactor the existing code to improve readability and maintainability")
	if err != nil {
		return nil, err
	}
	return llm.GetSimpleWorkspace("sources").Sources(), nil
}

// Bump dependencies with their latest versions
func (m *Toolbox) BumpDeps(ctx context.Context) (*dagger.Directory, error) {
	llm, err := m.asEditor(ctx, "find dependencies and their version, find the latest version and update them")
	if err != nil {
		return nil, err
	}
	return llm.GetSimpleWorkspace("sources").Sources(), nil
}

func (m *Toolbox) printLLMLastReply(ctx context.Context, llm *dagger.LLM) (string, error) {
	out, err := llm.LastReply(ctx)
	if err != nil {
		return "", err
	}
	return dag.Glow().DisplayMarkdown(ctx, out)
}

func (m *Toolbox) asEditor(ctx context.Context, assignment string) (*dagger.LLM, error) {
	return m.askTo(ctx, `
You are an expert programmer receiving an assignment.

- Use your workspace to complete the assignment
- Make the smallest possible changes
- Improve the code for readability
- Add comments when necessary to keep the code easy to understand
- Based on the language, use the appropriate builder to build and test the code: you have access to a maven builder and a go builder
- Always write changes to files to the original files in the workspace
- Create the file .toolbox.md if it doesn't exist
- Keep a log of all the changes made in .toolbox.md, in markdown format: includes explanations and before and after code to document

INPUTS:

1) Your assignment is:

<assignment>
$assignment
</assignment>

2) You have access to a workspace to do your work.
Inside this workspace you can read, write, rm files or walk the directory tree.
You can also run the tools copyDir, rmDir and listDir about directory operations.
And you can build and test.
`, assignment)
}

func (m *Toolbox) asReader(ctx context.Context, assignment string) (*dagger.LLM, error) {
	return m.askTo(ctx, `
You are an expert programmer receiving an assignment.

- Use your workspace to complete the assignment
- Only read files, never make changes in the workspace

INPUTS:

1) Your assignment is:

<assignment>
$assignment
</assignment>

2) You have access to a simple workspace containing the sources to do your work
Inside this workspace you can read files or walk the directory tree.
`, assignment)
}

func (m *Toolbox) askTo(ctx context.Context, prompt, assignment string) (*dagger.LLM, error) {
	return dag.Llm().
		SetSimpleWorkspace("sources", dag.SimpleWorkspace(m.Src)).
		SetMavenBuilder("maven-builder", dag.MavenBuilder()).
		SetGoBuilder("go-builder", dag.GoBuilder()).
		WithPromptVar("assignment", assignment).
		WithPrompt(prompt).Sync(ctx)
}

func New(
	// +defaultPath="."
	src *dagger.Directory) *Toolbox {
	return &Toolbox{
		Src: src,
	}
}
