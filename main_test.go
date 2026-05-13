package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func goEnvValue(t *testing.T, key string) string {
	t.Helper()
	out, err := exec.Command("go", "env", key).Output()
	if err != nil {
		t.Fatalf("go env %s error = %v", key, err)
	}
	return strings.TrimSpace(string(out))
}

func runQuickTodo(t *testing.T, configDir string, args ...string) string {
	t.Helper()
	cmdArgs := append([]string{"run", "."}, args...)
	cmd := exec.Command("go", cmdArgs...)
	cmd.Env = append(os.Environ(),
		"HOME="+configDir,
		"XDG_CONFIG_HOME="+configDir,
		"GOMODCACHE="+goEnvValue(t, "GOMODCACHE"),
		"GOCACHE="+goEnvValue(t, "GOCACHE"),
	)
	if runtime.GOOS == "windows" {
		cmd.Env = append(cmd.Env, "APPDATA="+configDir)
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go run . %v error = %v\n%s", args, err, string(out))
	}
	return string(out)
}

func TestCLIEndToEndTodoFlow(t *testing.T) {
	configDir := t.TempDir()

	output := runQuickTodo(t, configDir, "add", "write README")
	if !strings.Contains(output, "write README") {
		t.Fatalf("add output missing todo: %s", output)
	}

	output = runQuickTodo(t, configDir, "add", "ship release", "-d")
	if !strings.Contains(output, "ship release") || !strings.Contains(output, "✅") {
		t.Fatalf("add done output unexpected: %s", output)
	}

	output = runQuickTodo(t, configDir, "list")
	if !strings.Contains(output, "write README") || !strings.Contains(output, "ship release") {
		t.Fatalf("list output missing todos: %s", output)
	}

	output = runQuickTodo(t, configDir, "list", "-d")
	if !strings.Contains(output, "ship release") || strings.Contains(output, "write README") {
		t.Fatalf("list done filter output unexpected: %s", output)
	}

	output = runQuickTodo(t, configDir, "list", "-f", "readme", "-i")
	if !strings.Contains(output, "write README") || strings.Contains(output, "ship release") {
		t.Fatalf("list content filter output unexpected: %s", output)
	}

	output = runQuickTodo(t, configDir, "mod", "0", "update README")
	if !strings.Contains(output, "modified") || !strings.Contains(output, "update README") {
		t.Fatalf("mod output unexpected: %s", output)
	}

	output = runQuickTodo(t, configDir, "mod", "0", " - urgent", "--append")
	if !strings.Contains(output, "update README - urgent") {
		t.Fatalf("mod append output unexpected: %s", output)
	}

	output = runQuickTodo(t, configDir, "mod", "0", "NOW: ", "--insert")
	if !strings.Contains(output, "NOW: update README - urgent") {
		t.Fatalf("mod insert output unexpected: %s", output)
	}

	output = runQuickTodo(t, configDir, "mod", "0", "-d", "-p", "2")
	if !strings.Contains(output, "modified") || !strings.Contains(output, "✅") || !strings.Contains(output, "⭐⭐") {
		t.Fatalf("mod done priority output unexpected: %s", output)
	}

	output = runQuickTodo(t, configDir, "rm", "1")
	if !strings.Contains(output, "removed") || !strings.Contains(output, "ship release") || !strings.Contains(output, "last") {
		t.Fatalf("rm output unexpected: %s", output)
	}

	dataFile := filepath.Join(configDir, "Library", "Application Support", "QuickTodo", "todos", "default.json")
	if runtime.GOOS != "darwin" {
		dataFile = filepath.Join(configDir, "QuickTodo", "todos", "default.json")
	}
	if _, err := os.Stat(dataFile); err != nil {
		t.Fatalf("expected data file at %s: %v", dataFile, err)
	}
}

func TestCLIRootCommandShortcutsAndPlaceholders(t *testing.T) {
	configDir := t.TempDir()

	output := runQuickTodo(t, configDir, "root task")
	if !strings.Contains(output, "root task") {
		t.Fatalf("root add output unexpected: %s", output)
	}

	output = runQuickTodo(t, configDir)
	if !strings.Contains(output, "root task") {
		t.Fatalf("root list output unexpected: %s", output)
	}

	output = runQuickTodo(t, configDir, "done", "0")
	if !strings.Contains(output, "done todo: 0") {
		t.Fatalf("done placeholder output unexpected: %s", output)
	}

	output = runQuickTodo(t, configDir, "clear")
	if !strings.Contains(output, "clear all todos") {
		t.Fatalf("clear placeholder output unexpected: %s", output)
	}
}
