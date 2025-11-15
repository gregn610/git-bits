package bits

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestS3RemoteIntegration(t *testing.T) {
	bucket := os.Getenv("TEST_BUCKET")
	if bucket == "" {
		t.Skip("Skipping S3 integration test - no TEST_BUCKET set")
	}

	// Test S3Remote creation - AWS SDK handles credentials automatically
	s3Remote, err := NewS3Remote(nil, "origin", bucket)
	if err != nil {
		t.Fatalf("Failed to create S3Remote: %v", err)
	}

	if s3Remote.Name() != "origin" {
		t.Errorf("Expected remote name 'origin', got %s", s3Remote.Name())
	}

	// Test chunk operations
	testKey := K{0x01, 0x02, 0x03, 0x04, 0x05}
	testData := []byte("test chunk data for S3 integration")

	// Test ChunkWriter
	writer, err := s3Remote.ChunkWriter(testKey)
	if err != nil {
		t.Fatalf("Failed to get chunk writer: %v", err)
	}

	n, err := writer.Write(testData)
	if err != nil {
		t.Fatalf("Failed to write chunk data: %v", err)
	}
	if n != len(testData) {
		t.Errorf("Expected to write %d bytes, wrote %d", len(testData), n)
	}

	err = writer.Close()
	if err != nil {
		t.Fatalf("Failed to close chunk writer: %v", err)
	}

	// Test ChunkReader
	reader, err := s3Remote.ChunkReader(testKey)
	if err != nil {
		t.Fatalf("Failed to get chunk reader: %v", err)
	}
	defer reader.Close()

	// Read all data using io.ReadAll to handle EOF properly
	readData, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("Failed to read chunk data: %v", err)
	}
	if len(readData) != len(testData) {
		t.Errorf("Expected to read %d bytes, read %d", len(testData), len(readData))
	}

	if !bytes.Equal(testData, readData) {
		t.Error("Read data doesn't match written data")
	}

	// Test ListChunks
	listOutput := &bytes.Buffer{}
	err = s3Remote.ListChunks(listOutput)
	if err != nil {
		t.Fatalf("Failed to list chunks: %v", err)
	}

	// Should contain our test key
	keyHex := "0102030405000000000000000000000000000000000000000000000000000000"
	if !bytes.Contains(listOutput.Bytes(), []byte(keyHex)) {
		t.Errorf("List output should contain test key %s, got: %s", keyHex, listOutput.String())
	}
}

func TestS3RemoteErrorCases(t *testing.T) {
	// Test with invalid bucket
	s3Remote, err := NewS3Remote(nil, "origin", "invalid-bucket")
	if err != nil {
		t.Fatalf("S3Remote creation should not fail with invalid credentials: %v", err)
	}

	// Test reading non-existent chunk
	nonExistentKey := K{0xFF, 0xFF, 0xFF}
	_, err = s3Remote.ChunkReader(nonExistentKey)
	if err == nil {
		t.Error("Should fail when reading non-existent chunk")
	}
}

func TestS3ConfigurationVariations(t *testing.T) {
	// Test S3Remote creation with different remote names
	s3Remote, err := NewS3Remote(nil, "origin", "test-bucket")
	if err != nil {
		t.Fatalf("Should work: %v", err)
	}

	if s3Remote.Name() != "origin" {
		t.Error("Remote name should be preserved")
	}

	s3Remote, err = NewS3Remote(nil, "upstream", "test-bucket")
	if err != nil {
		t.Fatalf("Should work: %v", err)
	}

	if s3Remote.Name() != "upstream" {
		t.Error("Remote name should be preserved")
	}
}