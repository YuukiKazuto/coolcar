package oss

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"
)

type Service struct {
	bucket *oss.Bucket
}

func (s *Service) SignURL(method oss.HTTPMethod, path string, timeout int64) (string, error) {
	url, err := s.bucket.SignURL(path, method, timeout)
	if err != nil {
		return "", err
	}
	return url, nil
}

func (s *Service) Get(path string) (io.ReadCloser, error) {
	body, err := s.bucket.GetObject(path)
	return body, err
}

func NewService(endpoint, accessKeyID, accessKeySecret string) (*Service, error) {
	client, err := oss.New(endpoint, accessKeyID, accessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("cannot creates a new client: %v", err)
	}

	bucket, err := client.Bucket("coolcar-kh")
	if err != nil {
		return nil, fmt.Errorf("cannot gets the bucket instanc: %v", err)
	}

	return &Service{
		bucket: bucket,
	}, nil
}
