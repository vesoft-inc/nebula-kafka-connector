package source

import (
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var _ Source = (*s3Source)(nil)

type (
	s3Source struct {
		c    *Config
		obj  io.ReadCloser
		size int64
	}
)

func openS3File(c *Config) (*s3Source, error) {
	creds := credentials.NewStaticCredentials(c.S3.AccessKey, c.S3.SecretKey, c.S3.Token)

	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(c.S3.Region),
		Endpoint:         aws.String(c.S3.Endpoint),
		S3ForcePathStyle: aws.Bool(true),
		Credentials:      creds,
	})
	if err != nil {
		return nil, fmt.Errorf("create s3 session faild: %w", err)
	}

	svc := s3.New(sess)

	resp, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(c.S3.Bucket),
		Key:    aws.String(strings.TrimLeft(c.S3.Key, "/")),
	})
	if err != nil {
		return nil, err
	}

	return &s3Source{
		c:    c,
		obj:  resp.Body,
		size: *resp.ContentLength,
	}, nil
}

func (s *s3Source) Config() *Config {
	return s.c
}

func (s *s3Source) Size() (int64, error) {
	return s.size, nil
}

func (s *s3Source) Read(p []byte) (int, error) {
	return s.obj.Read(p)
}

func (s *s3Source) Close() error {
	return s.obj.Close()
}
