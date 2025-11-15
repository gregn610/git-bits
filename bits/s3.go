package bits

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Remote struct {
	gitRemote  string
	bucketName string
	client     *s3.Client
	repo       *Repository
}

func NewS3Remote(repo *Repository, remote, bucket string) (s3Remote *S3Remote, err error) {
	ctx := context.Background()
	
	// Load AWS configuration using standard AWS SDK configuration chain:
	// 1. Environment variables (AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, etc.)
	// 2. Shared credentials file (~/.aws/credentials)
	// 3. Shared config file (~/.aws/config)
	// 4. IAM roles for EC2/ECS/Lambda
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %v", err)
	}

	// Create S3 client - AWS SDK automatically handles AWS_ENDPOINT_URL
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		// Always use path-style addressing for compatibility with LocalStack and AWS S3
		o.UsePathStyle = true
	})

	return &S3Remote{
		repo:       repo,
		gitRemote:  remote,
		bucketName: bucket,
		client:     client,
	}, nil
}

func (s3 *S3Remote) Name() string {
	return s3.gitRemote
}

//ListChunks will write all chunks in the bucket to writer w
func (s *S3Remote) ListChunks(w io.Writer) (err error) {
	ctx := context.Background()
	
	paginator := s3.NewListObjectsV2Paginator(s.client, &s3.ListObjectsV2Input{
		Bucket:  aws.String(s.bucketName),
		MaxKeys: aws.Int32(500),
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("failed to list objects: %v", err)
		}

		for _, obj := range page.Contents {
			key := aws.ToString(obj.Key)
			// Only include keys that match chunk key format
			if len(key) == hex.EncodedLen(KeySize) {
				fmt.Fprintf(w, "%s\n", key)
			}
		}
	}

	return nil
}

//ChunkReader returns a file handle that the chunk with the given
//key can be read from, the user is expected to close it when finished
func (s *S3Remote) ChunkReader(k K) (rc io.ReadCloser, err error) {
	ctx := context.Background()
	
	resp, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(fmt.Sprintf("%x", k)),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %v", err)
	}
	
	return resp.Body, nil
}

// chunkWriter implements io.WriteCloser for S3 uploads
type chunkWriter struct {
	client     *s3.Client
	bucketName string
	key        string
	buffer     []byte
}

func (cw *chunkWriter) Write(p []byte) (n int, err error) {
	cw.buffer = append(cw.buffer, p...)
	return len(p), nil
}

func (cw *chunkWriter) Close() error {
	ctx := context.Background()
	_, err := cw.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(cw.bucketName),
		Key:    aws.String(cw.key),
		Body:   bytes.NewReader(cw.buffer),
	})
	return err
}

//ChunkWriter returns a file handle to which a chunk with give key
//can be written to, the user is expected to close it when finished.
func (s *S3Remote) ChunkWriter(k K) (wc io.WriteCloser, err error) {
	return &chunkWriter{
		client:     s.client,
		bucketName: s.bucketName,
		key:        fmt.Sprintf("%x", k),
		buffer:     make([]byte, 0),
	}, nil
}
