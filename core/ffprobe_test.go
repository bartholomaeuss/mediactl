package core

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWalkFFProbeProcessesFilesRecursively(t *testing.T) {
	root := t.TempDir()

	makeFile(t, filepath.Join(root, "movie.mkv"))
	makeFile(t, filepath.Join(root, "ignore.txt"))

	sub := filepath.Join(root, "nested", "deeper")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	makeFile(t, filepath.Join(sub, "clip.mkv"))

	var calls []string
	spyRunner := func(input, output string) error {
		calls = append(calls, input+"->"+output)
		return os.WriteFile(output, []byte("{}"), 0o644)
	}

	var out bytes.Buffer
	processed, err := WalkFFProbe(root, ".mkv", &out, spyRunner)
	if err != nil {
		t.Fatalf("WalkFFProbe returned error: %v", err)
	}
	if processed != 2 {
		t.Fatalf("expected 2 processed files, got %d", processed)
	}

	want := map[string]bool{
		filepath.Join(root, "movie.mkv") + "->" + filepath.Join(root, "movie.mkv.json"): true,
		filepath.Join(sub, "clip.mkv") + "->" + filepath.Join(sub, "clip.mkv.json"):     true,
	}

	if len(calls) != len(want) {
		t.Fatalf("unexpected number of ffprobe calls: got %d want %d", len(calls), len(want))
	}

	for _, c := range calls {
		if !want[c] {
			t.Fatalf("unexpected ffprobe invocation: %s", c)
		}
	}

	output := out.String()
	if !strings.Contains(output, "movie.mkv.json") || !strings.Contains(output, "clip.mkv.json") {
		t.Fatalf("expected output to mention written files, got %q", output)
	}
}

func TestWalkFFProbeAvoidsOverwritingSidecars(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "movie.mkv")
	makeFile(t, target)

	existing := target + ".json"
	makeFile(t, existing)
	makeFile(t, target+".1.json")

	var outputs []string
	spyRunner := func(input, output string) error {
		outputs = append(outputs, output)
		return os.WriteFile(output, []byte("{}"), 0o644)
	}

	processed, err := WalkFFProbe(root, ".mkv", nil, spyRunner)
	if err != nil {
		t.Fatalf("WalkFFProbe returned error: %v", err)
	}
	if processed != 1 {
		t.Fatalf("expected 1 processed file, got %d", processed)
	}

	if len(outputs) != 1 {
		t.Fatalf("expected 1 output file, got %d", len(outputs))
	}

	expected := target + ".2.json"
	if outputs[0] != expected {
		t.Fatalf("expected output %s, got %s", expected, outputs[0])
	}
}

func TestWalkFFProbePropagatesRunnerErrors(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "bad.mkv")
	makeFile(t, target)

	failingRunner := func(input, output string) error {
		return fmt.Errorf("runner fail on %s", input)
	}

	processed, err := WalkFFProbe(root, ".mkv", nil, failingRunner)
	if err == nil || !strings.Contains(err.Error(), "runner fail") {
		t.Fatalf("expected runner error, got %v", err)
	}
	if processed != 0 {
		t.Fatalf("expected processed count 0 on error, got %d", processed)
	}
}

func makeFile(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(path, []byte("dummy"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
}
