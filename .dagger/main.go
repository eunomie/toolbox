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

// Find potential bugs in the existing java code, explain them and propose alternative code to fix them
func (m *Toolbox) FindBugs(ctx context.Context) (string, error) {
	llm, err := m.asReader(ctx, "find potential bugs in the existing java code, explain them and propose alternative code to fix them")
	if err != nil {
		return "", err
	}
	return m.printLLMLastReply(ctx, llm)
}

// Ask to implement a new feature in the existing java code
func (m *Toolbox) Do(ctx context.Context, ask string) (*dagger.Directory, error) {
	llm, err := m.asEditor(ctx, ask)
	if err != nil {
		return nil, err
	}
	return llm.GetJavaWorkspace("javaWorkspace").Dir(), nil
}

// Add comments to the existing java code to improve readability and maintainability
func (m *Toolbox) AddComments(ctx context.Context) (*dagger.Directory, error) {
	llm, err := m.asEditor(ctx, "add comments to the existing java code to improve readability and maintainability")
	if err != nil {
		return nil, err
	}
	return llm.GetJavaWorkspace("javaWorkspace").Dir(), nil
}

// Refactor the existing java code to improve readability and maintainability
func (m *Toolbox) Refactor(ctx context.Context) (*dagger.Directory, error) {
	llm, err := m.asEditor(ctx, "refactor the existing java code to improve readability and maintainability")
	if err != nil {
		return nil, err
	}
	return llm.GetJavaWorkspace("javaWorkspace").Dir(), nil
}

// Bump dependencies with their latest versions
func (m *Toolbox) BumpDeps(ctx context.Context) (*dagger.Directory, error) {
	llm, err := m.asEditor(ctx, "find dependencies and their version, find the latest version and update the dependencies in the pom.xml file")
	if err != nil {
		return nil, err
	}
	return llm.GetJavaWorkspace("javaWorkspace").Dir(), nil
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
You are an expert Java programmer receiving an assignment.

- Use your workspace to complete the assignment
- Make the smallest possible changes
- Improve the code for readability
- Add comments when necessary to keep the code easy to understand
- Check that your work builds and pass tests with the 'check' tool
- Always write changes to files to the original files in the workspace
- Create the file .toolbox.md if it doesn't exist
- Keep a log of all the changes made in .toolbox.md, in markdown format: includes explanations and before and after code to document

INPUTS:

1) Your assignment is:

<assignment>
$assignment
</assignment>

2) You have access to a workspace to do your work
`, assignment)
}

func (m *Toolbox) asReader(ctx context.Context, assignment string) (*dagger.LLM, error) {
	return m.askTo(ctx, `
You are an expert Java programmer receiving an assignment.

- Use your workspace to complete the assignment
- Only read files, never make changes in the workspace

INPUTS:

1) Your assignment is:

<assignment>
$assignment
</assignment>

2) You have access to a workspace to do your work
`, assignment)
}

func (m *Toolbox) askTo(ctx context.Context, prompt, assignment string) (*dagger.LLM, error) {
	return dag.Llm().SetJavaWorkspace("javaWorkspace",
		dag.JavaWorkspace(m.Src)).
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
