// Copyright (c) 2026 Daniel Montanari. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/config"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/object"
)

// setupTestRepo creates a temporary git repository for testing purposes.
// It creates an initial commit, a dummy remote, and an annotated tag "v1.0.0".
// It returns the repository path and a cleanup function to be deferred.
func setupTestRepo(t *testing.T) (string, func()) {
	t.Helper()

	repoPath, err := os.MkdirTemp("", "tagger-test-repo-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	repo, err := git.PlainInit(repoPath, false)
	if err != nil {
		t.Fatalf("Failed to init repo: %v", err)
	}

	// Add a dummy remote so HaveRemote() passes during tests
	_, err = repo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{"/dev/null"}, // A dummy URL that doesn't need to be valid
	})
	if err != nil {
		t.Fatalf("Failed to create remote: %v", err)
	}

	w, err := repo.Worktree()
	if err != nil {
		t.Fatalf("Failed to get worktree: %v", err)
	}
	dummyFilePath := filepath.Join(repoPath, "README.md")
	if err := os.WriteFile(dummyFilePath, []byte("initial commit"), 0644); err != nil {
		t.Fatalf("Failed to write dummy file: %v", err)
	}
	if _, err := w.Add("README.md"); err != nil {
		t.Fatalf("Failed to add file to worktree: %v", err)
	}
	commit, err := w.Commit("initial commit", &git.CommitOptions{
		Author: &object.Signature{Name: "Tagger Test", When: time.Now()},
	})
	if err != nil {
		t.Fatalf("Failed to create initial commit: %v", err)
	}

	// Create an initial ANNOTATED tag
	_, err = repo.CreateTag("v1.0.0", commit, &git.CreateTagOptions{
		Tagger:  &object.Signature{Name: "Tagger Test", When: time.Now()},
		Message: "initial release",
	})
	if err != nil {
		t.Fatalf("Failed to create initial tag: %v", err)
	}

	cleanup := func() {
		os.RemoveAll(repoPath)
	}

	return repoPath, cleanup
}

// executeCommand runs a cobra command and captures its stdout.
func executeCommand(args ...string) (string, error) {
	// Keep old stdout to restore it later
	oldStdout := os.Stdout
	// Create a pipe to capture stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Reset flags and args for cobra
	rootCmd.SetArgs(args)
	incMajor, incMinor, incPatch, dryRun, verbose, fullOutput = false, false, false, false, false, false

	// Run the command
	cmdErr := rootCmd.Execute()

	// Close the writer and restore stdout
	w.Close()
	os.Stdout = oldStdout

	// Read the output from the pipe
	var out bytes.Buffer
	io.Copy(&out, r)

	return out.String(), cmdErr
}

func TestLastCommand(t *testing.T) {
	repoPath, cleanup := setupTestRepo(t)
	defer cleanup()

	output, err := executeCommand("last", repoPath)
	if err != nil {
		t.Fatalf("last command failed: %v", err)
	}

	expected := "v1.0.0\n"
	if output != expected {
		t.Errorf("expected output %q, got %q", expected, output)
	}
}

func TestListCommand(t *testing.T) {
	repoPath, cleanup := setupTestRepo(t)
	defer cleanup()

	output, err := executeCommand("list", repoPath)
	if err != nil {
		t.Fatalf("list command failed: %v", err)
	}

	if !strings.Contains(output, "v1.0.0") {
		t.Errorf("expected output to contain 'v1.0.0', got %q", output)
	}
}

func TestIncPatchDryRun(t *testing.T) {
	repoPath, cleanup := setupTestRepo(t)
	defer cleanup()

	output, err := executeCommand("inc", repoPath, "--patch", "--dry-run")
	if err != nil {
		t.Fatalf("inc command failed: %v", err)
	}

	expected := "v1.0.1\n"
	if !strings.HasSuffix(output, expected) {
		t.Errorf("expected output to end with %q, got %q", expected, output)
	}
}

func TestIncPatch(t *testing.T) {
	repoPath, cleanup := setupTestRepo(t)
	defer cleanup()

	// This will still print a push error to stderr because the remote is invalid,
	// but the test will pass if the tag is created locally. We can ignore the error.
	output, _ := executeCommand("inc", repoPath, "--patch")

	expectedOutput := "v1.0.1\n"
	if !strings.HasSuffix(output, expectedOutput) {
		t.Errorf("expected output to end with %q, got %q", expectedOutput, output)
	}

	// Now, verify the tag was actually created in the repo
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		t.Fatalf("Failed to open repo: %v", err)
	}

	tagrefs, err := repo.Tags()
	if err != nil {
		t.Fatalf("Failed to get tags: %v", err)
	}

	found := false
	err = tagrefs.ForEach(func(ref *plumbing.Reference) error {
		if ref.Name().Short() == "v1.0.1" {
			found = true
		}
		return nil
	})
	if err != nil {
		t.Fatalf("Error iterating tags: %v", err)
	}

	if !found {
		t.Error("expected tag 'v1.0.1' to be created, but it was not found")
	}
}