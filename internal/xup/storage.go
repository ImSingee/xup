package xup

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Storage struct {
	client *minio.Client
	bucket string
}

func Connect(c Config) (*Storage, error) {
	if err := c.Valid(); err != nil {
		return nil, err
	}
	endpoint := strings.TrimSuffix(c.Endpoint, "/")
	secure := strings.HasPrefix(endpoint, "https://")
	endpoint = strings.TrimPrefix(strings.TrimPrefix(endpoint, "https://"), "http://")
	lookup := minio.BucketLookupAuto
	if c.PathStyle {
		lookup = minio.BucketLookupPath
	}
	client, err := minio.New(endpoint, &minio.Options{Creds: credentials.NewStaticV4(c.AccessKeyID, c.SecretAccessKey, c.SessionToken), Secure: secure, Region: c.Region, BucketLookup: lookup})
	if err != nil {
		return nil, err
	}
	return &Storage{client, c.Bucket}, nil
}
func (s *Storage) Upload(ctx context.Context, local, key string) (int64, error) {
	f, err := os.Open(local)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	stat, err := f.Stat()
	if err != nil {
		return 0, err
	}
	if !stat.Mode().IsRegular() {
		return 0, errors.New("source must be a regular file")
	}
	info, err := s.client.PutObject(ctx, s.bucket, key, f, stat.Size(), minio.PutObjectOptions{})
	if err != nil {
		return 0, err
	}
	return info.Size, nil
}
func (s *Storage) Download(ctx context.Context, key, local string) (int64, error) {
	obj, err := s.client.GetObject(ctx, s.bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return 0, err
	}
	defer obj.Close()
	if _, err = obj.Stat(); err != nil {
		return 0, err
	}
	f, err := os.OpenFile(local, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return 0, fmt.Errorf("create destination (refusing to overwrite): %w", err)
	}
	ok := false
	defer func() {
		if !ok {
			os.Remove(local)
		}
	}()
	n, err := io.Copy(f, obj)
	closeErr := f.Close()
	if err != nil {
		return n, err
	}
	if closeErr != nil {
		return n, closeErr
	}
	ok = true
	return n, nil
}
func (s *Storage) Remove(ctx context.Context, key string) error {
	return s.client.RemoveObject(ctx, s.bucket, key, minio.RemoveObjectOptions{})
}
func (s *Storage) List(ctx context.Context, prefix string, output io.Writer) error {
	for item := range s.client.ListObjects(ctx, s.bucket, minio.ListObjectsOptions{Prefix: prefix, Recursive: true}) {
		if item.Err != nil {
			return item.Err
		}
		if _, err := fmt.Fprintf(output, "%d\t%s\n", item.Size, item.Key); err != nil {
			return err
		}
	}
	return nil
}
func DefaultKey(local string) string { return filepath.Base(local) }
func DefaultLocal(key string) string { return filepath.Base(strings.TrimSuffix(key, "/")) }
