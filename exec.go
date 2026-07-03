package coder

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/togo-framework/autopilot"
	"github.com/togo-framework/providers"
	"github.com/togo-framework/togo"
)

// coderExec implements autopilot.Executor — it provisions a fresh Coder workspace
// per issue and tears it down afterwards.
type coderExec struct {
	bin      string
	template string
	repoDir  string
	keep     bool
}

func newExec(k *togo.Kernel) *coderExec {
	return &coderExec{
		bin:      providers.Value(k, providers.CapExecute, Name, "bin", "coder", false),
		template: providers.Value(k, providers.CapExecute, Name, "template", "", false),
		repoDir:  providers.Value(k, providers.CapExecute, Name, "repo_dir", "~/workspace", false),
		keep:     providers.Value(k, providers.CapExecute, Name, "keep", "", false) == "1",
	}
}

// Acquire creates a workspace from the configured template and returns an Env
// that runs commands inside it. release() deletes the workspace (unless CODER_KEEP).
func (c *coderExec) Acquire(ctx context.Context, issue autopilot.Issue) (autopilot.Env, func(), error) {
	if c.template == "" {
		return nil, nil, fmt.Errorf("coder: CODER_TEMPLATE is required (the workspace template to provision)")
	}
	name := "autopilot-" + shortID(issue.ID)
	if out, err := c.run(ctx, "create", name, "--template", c.template, "-y"); err != nil {
		return nil, nil, fmt.Errorf("coder create %s: %v: %s", name, err, strings.TrimSpace(out))
	}
	release := func() {
		if c.keep {
			return
		}
		_, _ = c.run(context.Background(), "delete", name, "-y")
	}
	return coderEnv{bin: c.bin, ws: name, dir: c.repoDir}, release, nil
}

func (c *coderExec) run(ctx context.Context, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, c.bin, args...)
	var b bytes.Buffer
	cmd.Stdout, cmd.Stderr = &b, &b
	err := cmd.Run()
	return b.String(), err
}

// coderEnv implements autopilot.Env by running commands inside a Coder workspace
// via `coder ssh <ws> -- bash -lc 'cd <dir> && <cmd>'`.
type coderEnv struct {
	bin string
	ws  string
	dir string
}

func (e coderEnv) Dir() string { return e.dir }

func (e coderEnv) Run(ctx context.Context, name string, args ...string) (string, error) {
	remote := "cd " + shellQuote(e.dir) + " && " + shellJoin(append([]string{name}, args...))
	cmd := exec.CommandContext(ctx, e.bin, "ssh", e.ws, "--", "bash", "-lc", remote)
	var out, errb bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &errb
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("coder ssh %s: %v: %s", e.ws, err, strings.TrimSpace(errb.String()))
	}
	return out.String(), nil
}

func shortID(id string) string {
	if len(id) > 8 {
		return id[:8]
	}
	return id
}

// shellQuote single-quotes a string for safe embedding in a remote bash -lc.
func shellQuote(s string) string { return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'" }

func shellJoin(parts []string) string {
	q := make([]string, len(parts))
	for i, p := range parts {
		q[i] = shellQuote(p)
	}
	return strings.Join(q, " ")
}
