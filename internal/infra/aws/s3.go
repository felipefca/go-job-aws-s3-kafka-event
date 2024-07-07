package aws

import (
	"context"
	"fmt"
	"go-job-aws-s3-kafka-event/internal/appctx"
	"io"
	"regexp"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Storage interface {
	Reader(ctx context.Context, key string) (io.ReadCloser, error)
	ReaderListObjects(ctx context.Context, prefix string) (*s3.ListObjectsV2Output, error)
	CopyObject(ctx context.Context, fileName string, destination string, key string) error
	MoveFile(ctx context.Context, fileName string, destination string) error
}

type storage struct {
	client     *s3.Client
	buckerName string
	folder     string
}

func NewStorage(client *s3.Client, bucketName string, folder string) Storage {
	return &storage{
		client:     client,
		buckerName: bucketName,
		folder:     folder,
	}
}

func (s *storage) Reader(ctx context.Context, key string) (io.ReadCloser, error) {
	out, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &s.buckerName,
		Key:    &key,
	})
	if err != nil {
		return nil, err
	}

	return out.Body, err
}

func (s *storage) ReaderListObjects(ctx context.Context, prefix string) (*s3.ListObjectsV2Output, error) {
	out, err := s.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket:    &s.buckerName,
		Prefix:    &prefix,
		Delimiter: aws.String("/"),
	})
	if err != nil {
		return nil, err
	}

	return out, err
}

func (s *storage) CopyObject(ctx context.Context, fileName string, destination string, key string) error {
	copySource := fmt.Sprintf("%s/%s", s.buckerName, key)
	destinationKey := fmt.Sprintf("%s%s", destination, fileName)

	input := &s3.CopyObjectInput{
		Bucket:     &s.buckerName,
		CopySource: &copySource,
		Key:        &destinationKey,
	}

	_, err := s.client.CopyObject(ctx, input)
	if err != nil {
		return err
	}

	return nil
}

func (s *storage) MoveFile(ctx context.Context, fileName string, destination string) error {
	logger := appctx.FromContext(ctx)
	logger.Info("Moving file to desstination folder")

	inputRegex := regexp.MustCompile(s.folder + fileName)

	resp, err := s.ReaderListObjects(ctx, s.folder)
	if err != nil {
		return err
	}

	for _, obj := range resp.Contents {
		if !inputRegex.MatchString(*obj.Key) {
			continue
		}

		err := s.CopyObject(ctx, fileName, destination, *obj.Key)
		if err != nil {
			return err
		}

		_, err = s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: &s.buckerName,
			Key:    obj.Key,
		})
		if err != nil {
			return err
		}

		logger.Info("File moved successfully: " + fileName)
	}

	return nil
}
