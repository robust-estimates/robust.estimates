package s3

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Storage struct {
	Client *minio.Client
	Bucket string
}

func New(endpoint, key, secret, bucket string, secure bool) (*Storage, error) {
	c, e := minio.New(endpoint, &minio.Options{Creds: credentials.NewStaticV4(key, secret, ""), Secure: secure})

	return &Storage{c, bucket}, e
}

func (s *Storage) Put(ctx context.Context, key string, r io.Reader, contentType string) (int64, error) {
	v, e := s.Client.PutObject(ctx, s.Bucket, key, r, -1, minio.PutObjectOptions{ContentType: contentType, PartSize: 5 * 1024 * 1024})

	return v.Size, e
}

func (s *Storage) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	o, e := s.Client.GetObject(ctx, s.Bucket, key, minio.GetObjectOptions{})
	if e != nil {
		return nil, e
	}

	if _, e = o.Stat(); e != nil {
		o.Close()
		return nil, e
	}

	return o, nil
}

func (s *Storage) Delete(ctx context.Context, key string) error {
	return s.Client.RemoveObject(ctx, s.Bucket, key, minio.RemoveObjectOptions{})
}
