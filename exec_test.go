package coder

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/togo-framework/togo"
)

func TestShellQuoteJoin(t *testing.T) {
	if got := shellQuote("a b"); got != "'a b'" {
		t.Fatalf("quote: %q", got)
	}
	if got := shellQuote("it's"); got != `'it'\''s'` {
		t.Fatalf("quote apostrophe: %q", got)
	}
	if got := shellJoin([]string{"git", "status", "--porcelain"}); got != "'git' 'status' '--porcelain'" {
		t.Fatalf("join: %q", got)
	}
}

func TestNewExecReadsConfig(t *testing.T) {
	t.Setenv("CODER_TEMPLATE", "omni-dev")
	t.Setenv("CODER_REPO_DIR", "/home/coder/workspace/app")
	c := newExec(togo.New())
	if c.template != "omni-dev" || c.repoDir != "/home/coder/workspace/app" {
		t.Fatalf("config not read: %+v", c)
	}
}

// TestEnvRunConstructsSSH uses a fake `coder` binary that echoes its args, to
// verify Run builds `coder ssh <ws> -- bash -lc 'cd <dir> && <cmd>'`.
func TestEnvRunConstructsSSH(t *testing.T) {
	dir := t.TempDir()
	fake := filepath.Join(dir, "coder")
	if err := os.WriteFile(fake, []byte("#!/bin/sh\nprintf '%s\\n' \"$@\"\n"), 0755); err != nil {
		t.Fatal(err)
	}
	e := coderEnv{bin: fake, ws: "autopilot-abc", dir: "/repo"}
	out, err := e.Run(context.Background(), "git", "status", "--porcelain")
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 6 || lines[0] != "ssh" || lines[1] != "autopilot-abc" || lines[3] != "bash" || lines[4] != "-lc" {
		t.Fatalf("unexpected ssh argv: %q", lines)
	}
	if !strings.Contains(lines[5], "cd '/repo' &&") || !strings.Contains(lines[5], "'git' 'status' '--porcelain'") {
		t.Fatalf("remote cmd wrong: %q", lines[5])
	}
}
