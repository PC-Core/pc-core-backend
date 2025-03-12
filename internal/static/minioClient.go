package static

import (
	"context"

	"github.com/PC-Core/pc-core-backend/internal/errors"
	"github.com/PC-Core/pc-core-backend/internal/static/serrors"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIOClient struct {
	client *minio.Client
	bucket string
}

func NewMinIOClient(ep string, access string, secret string, secure bool, bucket string) (*MinIOClient, errors.PCCError) {
	var err error
	client, err := minio.New(ep, &minio.Options{
		Creds:  credentials.NewStaticV4(access, secret, ""),
		Secure: secure,
	})

	if err != nil {
		return nil, serrors.MinIOCast(err)
	}

	return &MinIOClient{
		client,
		bucket,
	}, nil
}

func (c *MinIOClient) UploadFiles(files []StaticFile) ([]string, errors.PCCError) {
	locs := make([]string, 0, len(files))

	for _, file := range files {

		info, err := c.client.PutObject(context.Background(), c.bucket, file.Name, file.File, -1, minio.PutObjectOptions{
			ContentType: file.ContentType,
		})

		if err != nil {
			return nil, serrors.MinIOCast(err)
		}

		locs = append(locs, info.Location)
	}

	return locs, nil
}

func (c *MinIOClient) BucketExists() (bool, errors.PCCError) {
	exist, err := c.client.BucketExists(context.Background(), c.bucket)

	if err != nil {
		return exist, serrors.MinIOCast(err)
	}

	return exist, nil
}
