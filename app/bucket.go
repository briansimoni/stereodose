package app

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/google/go-cloud/blob"
	"github.com/google/go-cloud/blob/s3blob"
	"github.com/pkg/errors"
)

func setupBucket(cloud, bucketName string) (*blob.Bucket, error) {
	ctx := context.Background()
	// Open a connection to the bucket.
	var (
		b   *blob.Bucket
		err error
	)
	switch cloud {
	case "gcp":
		b, err = setupGCP(ctx, bucketName)
	case "aws":
		// AWS is handled below in the next code sample.
		b, err = setupAWS(ctx, bucketName)
	default:
		return nil, errors.Errorf("Failed to recognize cloud. Want gcp or aws, got: %s", cloud)
	}
	if err != nil {
		return nil, err
	}
	return b, nil
}

func setupAWS(ctx context.Context, bucketName string) (*blob.Bucket, error) {
	c := &aws.Config{
		// Either hard-code the region or use AWS_REGION.
		Region: aws.String("us-east-1"),
		// credentials.NewEnvCredentials assumes two environment variables are
		// present:
		// 1. AWS_ACCESS_KEY_ID, and
		// 2. AWS_SECRET_ACCESS_KEY.
		Credentials: credentials.NewEnvCredentials(),
	}
	s := session.Must(session.NewSession(c))
	return s3blob.OpenBucket(ctx, s, bucketName)
}

func setupGCP(ctx context.Context, bucketName string) (*blob.Bucket, error) {
	return nil, errors.New("Google bucket storage is not yet implemented")
}
