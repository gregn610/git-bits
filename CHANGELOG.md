# Changelog

## Unreleased

### Update Go to 1.24 and modernize dependencies
- Update Go version to 1.24
- Upgrade AWS SDK to v2 (github.com/aws/aws-sdk-go-v2)
- Update bbolt to maintained fork
- Add Cobra CLI framework for better command structure
- Modernize all dependencies for better security and performance

### Add LocalStack integration testing with Docker
- Add Docker Compose setup for LocalStack S3 testing
- Create containerized test environment
- Add Makefile targets: docker-test, localstack-up, localstack-down
- Include comprehensive S3 integration tests
- Support both local development and CI/CD testing

## Released

### 0.3.2
 - As forked
 