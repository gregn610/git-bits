package bits

import (
	"bytes"
	"testing"
)

func TestS3RemoteName(t *testing.T) {
	s3 := &S3Remote{
		gitRemote: "origin",
	}
	
	if s3.Name() != "origin" {
		t.Errorf("Expected Name() to return 'origin', got %s", s3.Name())
	}
}

func TestChunkWriter(t *testing.T) {
	cw := &chunkWriter{
		bucketName: "test-bucket",
		key:        "test-key",
		buffer:     make([]byte, 0),
	}
	
	// Test Write
	data := []byte("test data")
	n, err := cw.Write(data)
	if err != nil {
		t.Errorf("Write failed: %v", err)
	}
	if n != len(data) {
		t.Errorf("Expected to write %d bytes, wrote %d", len(data), n)
	}
	
	// Test buffer contains data
	if !bytes.Equal(cw.buffer, data) {
		t.Errorf("Buffer doesn't contain expected data")
	}
	
	// Test multiple writes
	moreData := []byte(" more data")
	cw.Write(moreData)
	
	expected := append(data, moreData...)
	if !bytes.Equal(cw.buffer, expected) {
		t.Errorf("Buffer doesn't contain expected combined data")
	}
}

func TestNewS3RemoteValidation(t *testing.T) {
	// Test S3Remote creation - AWS SDK handles credentials automatically
	_, err := NewS3Remote(nil, "origin", "test-bucket")
	if err != nil {
		t.Errorf("NewS3Remote should not fail: %v", err)
	}
}