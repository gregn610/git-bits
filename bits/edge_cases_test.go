package bits

import (
	"bytes"
	"crypto/rand"
	"io"
	"os"
	"strings"
	"testing"
)

func TestLargeFileChunking(t *testing.T) {
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

	// Test with larger data (1MB)
	largeData := make([]byte, 1024*1024)
	if _, err := rand.Read(largeData); err != nil {
		t.Fatal(err)
	}

	input := bytes.NewReader(largeData)
	output := &bytes.Buffer{}
	
	err = repo.Split(input, output)
	if err != nil {
		t.Fatal(err)
	}

	// Should produce multiple chunks
	keys := output.String()
	lines := strings.Split(strings.TrimSpace(keys), "\n")
	
	// Filter out header/footer
	keyLines := 0
	for _, line := range lines {
		if len(line) == 64 { // Hex encoded 32-byte key
			keyLines++
		}
	}
	
	if keyLines == 0 {
		t.Error("Large file should produce at least one chunk")
	}

	// Test combine
	keysInput := strings.NewReader(keys)
	combinedOutput := &bytes.Buffer{}
	
	err = repo.Combine(keysInput, combinedOutput)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(largeData, combinedOutput.Bytes()) {
		t.Error("Large file combine failed")
	}
}

func TestInvalidKeyFormats(t *testing.T) {
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

	// Test invalid hex key
	invalidKeys := strings.NewReader(
		"--- to use this file decode it with the 'git-bits' extension ---\n" +
		"invalid_hex_key\n" +
		"----------------------- end of chunks --------------------------\n")

	output := &bytes.Buffer{}
	err = repo.Combine(invalidKeys, output)
	if err == nil {
		t.Error("Should fail with invalid hex key")
	}

	// Test wrong key length
	wrongLengthKeys := strings.NewReader(
		"--- to use this file decode it with the 'git-bits' extension ---\n" +
		"0123456789abcdef\n" + // Too short
		"----------------------- end of chunks --------------------------\n")

	output.Reset()
	err = repo.Combine(wrongLengthKeys, output)
	if err == nil {
		t.Error("Should fail with wrong key length")
	}
}

func TestEmptyAndCornerCases(t *testing.T) {
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

	// Test single byte
	singleByte := []byte{0x42}
	input := bytes.NewReader(singleByte)
	output := &bytes.Buffer{}
	
	err = repo.Split(input, output)
	if err != nil {
		t.Fatal(err)
	}

	keysInput := strings.NewReader(output.String())
	combinedOutput := &bytes.Buffer{}
	
	err = repo.Combine(keysInput, combinedOutput)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(singleByte, combinedOutput.Bytes()) {
		t.Error("Single byte round-trip failed")
	}

	// Test binary data with null bytes
	binaryData := []byte{0x00, 0xFF, 0x00, 0xFF, 0x42, 0x00}
	input = bytes.NewReader(binaryData)
	output.Reset()
	
	err = repo.Split(input, output)
	if err != nil {
		t.Fatal(err)
	}

	keysInput = strings.NewReader(output.String())
	combinedOutput.Reset()
	
	err = repo.Combine(keysInput, combinedOutput)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(binaryData, combinedOutput.Bytes()) {
		t.Error("Binary data round-trip failed")
	}
}

func TestForEachErrorHandling(t *testing.T) {
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

	// Test ForEach with error in callback
	keyStream := strings.NewReader(
		"--- to use this file decode it with the 'git-bits' extension ---\n" +
		"0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef\n" +
		"----------------------- end of chunks --------------------------\n")

	testError := io.ErrUnexpectedEOF
	err = repo.ForEach(keyStream, func(k K) error {
		return testError
	})

	if err == nil {
		t.Error("Expected an error, got nil")
	} else if !strings.Contains(err.Error(), testError.Error()) {
		t.Errorf("Expected error containing %q, got %v", testError.Error(), err)
	}
}

func TestKeyOperations(t *testing.T) {
	// Test Key type operations
	k1 := K{}
	k2 := K{}
	
	// Keys should be equal when zero
	if k1 != k2 {
		t.Error("Zero keys should be equal")
	}

	// Test key with data
	k1[0] = 0x42
	if k1 == k2 {
		t.Error("Modified key should not equal zero key")
	}

	// Test KeyOp structure
	kop := KeyOp{
		Op:      PushOp,
		K:       k1,
		Skipped: true,
		CopyN:   1024,
	}

	if kop.Op != PushOp {
		t.Error("KeyOp operation not set correctly")
	}
	if !kop.Skipped {
		t.Error("KeyOp skipped flag not set correctly")
	}
	if kop.CopyN != 1024 {
		t.Error("KeyOp copy count not set correctly")
	}
}

func TestChunkBufferSize(t *testing.T) {
	// Test that ChunkBufferSize is reasonable
	if ChunkBufferSize <= 0 {
		t.Error("ChunkBufferSize should be positive")
	}
	
	if ChunkBufferSize < 1024 {
		t.Error("ChunkBufferSize should be at least 1KB")
	}
}

func TestRemoteChunkConstant(t *testing.T) {
	// Test RemoteChunk constant
	if RemoteChunk == nil {
		t.Error("RemoteChunk should not be nil")
	}
	
	if len(RemoteChunk) != 0 {
		t.Error("RemoteChunk should be empty slice")
	}
}