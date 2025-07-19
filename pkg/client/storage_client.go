package client

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// StorageClient 存储客户端接口
type StorageClient interface {
	BaseClient
	UploadFile(ctx context.Context, bucket, objectName string, reader io.Reader, size int64) error
	DownloadFile(ctx context.Context, bucket, objectName string) (io.ReadCloser, error)
	DeleteFile(ctx context.Context, bucket, objectName string) error
	GetFileURL(ctx context.Context, bucket, objectName string, expire time.Duration) (string, error)
	ListFiles(ctx context.Context, bucket string, prefix string) ([]FileInfo, error)
}

// FileInfo 文件信息
type FileInfo struct {
	Name         string
	Size         int64
	LastModified time.Time
	ContentType  string
}

// StorageConfig 存储客户端配置
type StorageConfig struct {
	// MinIO 服务器地址
	Endpoint string
	// MinIO 访问密钥
	AccessKey string
	// MinIO 访问密钥
	SecretKey string
	// MinIO 是否使用 SSL
	UseSSL bool
	// MinIO 区域,表示存储桶的地理位置
	Region string
}

// storageClient 存储客户端实现
type storageClient struct {
	client *minio.Client
	config StorageConfig
}

// NewStorageClient 创建存储客户端
func NewStorageClient(config StorageConfig) (StorageClient, error) {
	client, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKey, config.SecretKey, ""),
		Secure: config.UseSSL,
		Region: config.Region,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create storage client: %v", err)
	}

	return &storageClient{
		client: client,
		config: config,
	}, nil
}

// UploadFile 上传文件
func (c *storageClient) UploadFile(ctx context.Context, bucket, objectName string, reader io.Reader, size int64) error {
	_, err := c.client.PutObject(ctx, bucket, objectName, reader, size, minio.PutObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to upload file: %v", err)
	}
	return nil
}

// DownloadFile 下载文件
func (c *storageClient) DownloadFile(ctx context.Context, bucket, objectName string) (io.ReadCloser, error) {
	object, err := c.client.GetObject(ctx, bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %v", err)
	}
	defer object.Close()
	return object, nil
}

// DeleteFile 删除文件
func (c *storageClient) DeleteFile(ctx context.Context, bucket, objectName string) error {
	err := c.client.RemoveObject(ctx, bucket, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file: %v", err)
	}
	return nil
}

// GetFileURL 获取文件URL
func (c *storageClient) GetFileURL(ctx context.Context, bucket, objectName string, expire time.Duration) (string, error) {
	url, err := c.client.PresignedGetObject(ctx, bucket, objectName, expire, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get file url: %v", err)
	}
	return url.String(), nil
}

// ListFiles 列出文件
func (c *storageClient) ListFiles(ctx context.Context, bucket, prefix string) ([]FileInfo, error) {
	var files []FileInfo

	for object := range c.client.ListObjects(ctx, bucket, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	}) {
		if object.Err != nil {
			// %w 保留错误信息  %v 保留错误信息
			return nil, fmt.Errorf("failed to list files: %w", object.Err)
		}
		files = append(files, FileInfo{
			Name:         object.Key,
			Size:         object.Size,
			LastModified: object.LastModified,
			ContentType:  object.ContentType,
		})
	}
	return files, nil
}

// Close 关闭客户端
func (c *storageClient) Close() error {
	return nil
}

// Ping 测试连接
func (c *storageClient) Ping(ctx context.Context) error {
	_, err := c.client.ListBuckets(ctx)
	if err != nil {
		return fmt.Errorf("failed to ping storage client: %v", err)
	}
	return nil
}
