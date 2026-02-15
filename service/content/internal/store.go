package content

import (
	"bytes"
	"context"
	"io"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Store struct {
	client     *minio.Client
	bucketName string
}

func NewStore(endpoint, accessKey, secretKey, bucket string, useSSL bool) (*Store, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		return nil, err
	}
	if !exists {
		if err := client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, err
		}
	}

	return &Store{
		client:     client,
		bucketName: bucket,
	}, nil
}

func (s *Store) SaveContent(ctx context.Context, body string) (string, error) {
	id := uuid.NewString()
	reader := bytes.NewReader([]byte(body))

	_, err := s.client.PutObject(ctx, s.bucketName, id, reader, int64(reader.Len()), minio.PutObjectOptions{
		ContentType: "text/plain",
	})
	if err != nil {
		return "", err
	}
	return id, nil
}

func (s *Store) GetContent(ctx context.Context, id string) (string, error) {
	obj, err := s.client.GetObject(ctx, s.bucketName, id, minio.GetObjectOptions{})
	if err != nil {
		return "", err
	}
	defer obj.Close()

	b, err := io.ReadAll(obj)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
