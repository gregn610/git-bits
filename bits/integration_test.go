package bits

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestSplitCombineIntegration(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "git-bits-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Initialize git repo
	if err := initGitRepo(tmpDir); err != nil {
		t.Skip("Git not available, skipping integration test")
	}

	repo, err := NewRepository(tmpDir, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Test data
	testData := []byte("Hello, World! This is test data for git-bits integration testing.")
	
	// Test Split
	input := bytes.NewReader(testData)
	output := &bytes.Buffer{}
	
	err = repo.Split(input, output)
	if err != nil {
		t.Fatalf("Split failed: %v", err)
	}

	keys := output.String()
	if !strings.Contains(keys, "--- to use this file decode it with the 'git-bits' extension ---") {
		t.Error("Split output should contain header")
	}

	// Test Combine
	keysInput := strings.NewReader(keys)
	combinedOutput := &bytes.Buffer{}
	
	err = repo.Combine(keysInput, combinedOutput)
	if err != nil {
		t.Fatalf("Combine failed: %v", err)
	}

	if !bytes.Equal(testData, combinedOutput.Bytes()) {
		t.Error("Combined data doesn't match original")
	}
}

func TestRepositoryPath(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "git-bits-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	if err := initGitRepo(tmpDir); err != nil {
		t.Skip("Git not available")
	}

	repo, err := NewRepository(tmpDir, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Test Path function
	k := K{0x01, 0x02, 0x03}
	
	// Test without mkdir
	path, err := repo.Path(k, false)
	if err != nil {
		t.Fatal(err)
	}
	
	if !strings.HasSuffix(path, "0000000000000000000000000000000000000000000000000000000000") {
		t.Errorf("Path doesn't match expected pattern, got: %s", path)
	}

	// Test with mkdir
	path, err = repo.Path(k, true)
	if err != nil {
		t.Fatal(err)
	}
	
	// Check directory was created
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Error("Directory should have been created")
	}
}

func TestForEachKeys(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "git-bits-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	if err := initGitRepo(tmpDir); err != nil {
		t.Skip("Git not available")
	}

	repo, err := NewRepository(tmpDir, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create test key stream
	keyStream := strings.NewReader(
		"--- to use this file decode it with the 'git-bits' extension ---\n" +
		"0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef\n" +
		"fedcba9876543210fedcba9876543210fedcba9876543210fedcba9876543210\n" +
		"----------------------- end of chunks --------------------------\n")

	var processedKeys []K
	err = repo.ForEach(keyStream, func(k K) error {
		processedKeys = append(processedKeys, k)
		return nil
	})

	if err != nil {
		t.Fatal(err)
	}

	if len(processedKeys) != 2 {
		t.Errorf("Expected 2 keys, got %d", len(processedKeys))
	}
}

func TestLocalStore(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "git-bits-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	if err := initGitRepo(tmpDir); err != nil {
		t.Skip("Git not available")
	}

	repo, err := NewRepository(tmpDir, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Test LocalStore creation
	store, err := repo.LocalStore()
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	// Verify store is functional
	if store == nil {
		t.Error("LocalStore should not be nil")
	}

	// Test that database file was created
	dbPath := filepath.Join(repo.chunkDir, "a.chunks")
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Error("Database file should have been created")
	}
}

func TestErrorCases(t *testing.T) {
	// Test NewRepository with invalid directory
	_, err := NewRepository("/nonexistent/directory", nil)
	if err == nil {
		t.Error("Should fail with nonexistent directory")
	}

	// Test with non-git directory
	tmpDir, err := os.MkdirTemp("", "git-bits-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	_, err = NewRepository(tmpDir, nil)
	if err == nil {
		t.Error("Should fail with non-git directory")
	}
}

func TestChunkOperations(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "git-bits-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	if err := initGitRepo(tmpDir); err != nil {
		t.Skip("Git not available")
	}

	repo, err := NewRepository(tmpDir, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Test empty input
	emptyInput := strings.NewReader("")
	output := &bytes.Buffer{}
	
	err = repo.Split(emptyInput, output)
	if err != nil {
		t.Fatal(err)
	}

	// Should still have header/footer
	result := output.String()
	if !strings.Contains(result, "--- to use this file decode it with the 'git-bits' extension ---") {
		t.Error("Empty split should still have header")
	}

	// Test already chunked file (should pass through)
	chunkedInput := strings.NewReader(result)
	passthroughOutput := &bytes.Buffer{}
	
	err = repo.Split(chunkedInput, passthroughOutput)
	if err != nil {
		t.Fatal(err)
	}

	if result != passthroughOutput.String() {
		t.Error("Already chunked file should pass through unchanged")
	}
}

// Helper function to initialize git repo
func initGitRepo(dir string) error {
	// Try to initialize git repo
	cmd := []string{"git", "init"}
	if err := runCommand(dir, cmd...); err != nil {
		return err
	}
	
	// Configure git user (required for commits)
	if err := runCommand(dir, "git", "config", "user.name", "Test User"); err != nil {
		return err
	}
	
	if err := runCommand(dir, "git", "config", "user.email", "test@example.com"); err != nil {
		return err
	}
	
	return nil
}

func runCommand(dir string, cmd ...string) error {
	// Simple command runner for tests
	exec := exec.Command(cmd[0], cmd[1:]...)
	exec.Dir = dir
	return exec.Run()
}