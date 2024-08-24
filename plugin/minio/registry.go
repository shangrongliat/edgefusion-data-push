package minio

import (
	"bytes"
	"context"
	"log"
	"net/http"

	"edgefusion-data-push/plugin/config"
	"edgefusion-data-push/plugin/logs"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Minio interface {
	PutNetObject(url, bucket, objectName string) error
	PutFileObject(url, bucket, objectName string) (*minio.UploadInfo, error)
	PutStreamObject(bucket, objectName string, data []byte) error
}

type MinioClient struct {
	client *minio.Client
}

func NewMinioService(cfg *config.Config) (Minio, error) {
	// 初使化 minio client对象。
	client, err := minio.New(cfg.Minio.EndPoint, &minio.Options{
		Creds: credentials.NewStaticV4(cfg.Minio.AccessKey, cfg.Minio.SecretKey, ""),
	})
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	return &MinioClient{
		client: client,
	}, nil
}

func (m *MinioClient) PutNetObject(url, bucket, objectName string) error {
	// 打开网络文件
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	// 上传网络文件到MinIO
	_, err = m.client.PutObject(context.Background(), bucket, objectName, resp.Body, -1, minio.PutObjectOptions{})
	if err != nil {
		logs.L().Error("Failed to put object", logs.Error(err))
		return err
	}
	return nil
}

func (m *MinioClient) PutFileObject(url, bucket, objectName string) (*minio.UploadInfo, error) {
	// 上传本地文件到MinIO
	uploadInfo, err := m.client.FPutObject(context.Background(), bucket, objectName, url, minio.PutObjectOptions{})
	if err != nil {
		logs.L().Error("Failed to put object", logs.Error(err))
		return nil, err
	}
	return &uploadInfo, nil
}
func (m *MinioClient) PutStreamObject(bucket, objectName string, data []byte) error {
	// 上传本地文件到MinIO
	_, err := m.client.PutObject(context.Background(), bucket, objectName, bytes.NewReader(data), -1, minio.PutObjectOptions{})
	if err != nil {
		logs.L().Error("Failed to put object", logs.Error(err))
		return err
	}
	return nil
}
