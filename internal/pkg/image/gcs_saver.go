package image

import (
	"bytes"
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"image/jpeg"
	"image/png"
	"io"
	"path"
	"path/filepath"
	"time"
)

func NewGcsSaver(gcsClient *storage.Client, bucketName, projectId, uploadPath string) Saver {

	gcsBucketClient := &GcsBucketClient{
		cl:         gcsClient,
		bucketName: bucketName,
		projectID:  projectId,
		location:   uploadPath,
	}

	return &GcsSaver{
		gcsBucketClient: gcsBucketClient,
	}
}

type GcsSaver struct {
	gcsBucketClient *GcsBucketClient
}

func (gs *GcsSaver) SaveBase64(base64 string) (string, error) {

	img, err := decodeFromBase64(base64)

	if err != nil {
		return "", err
	}

	filename := generateFileName(img.mimeType)

	err = gs.gcsBucketClient.uploadFile(img, filename)
	if err != nil {
		return "", err
	}

	// generate full filepath

	return path.Join(gs.gcsBucketClient.location, filename), nil
}

func (gs *GcsSaver) GetImage(filePath string) (*bytes.Buffer, string, error) {
	return gs.gcsBucketClient.getFile(filePath)
}

func (gs *GcsSaver) DeleteFile(fileName string) error {
	return gs.gcsBucketClient.deleteFile(fileName)
}

type GcsBucketClient struct {
	cl         *storage.Client
	projectID  string
	bucketName string
	location   string
}

func (c *GcsBucketClient) deleteFile(fileName string) error {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	err := c.cl.Bucket(c.bucketName).Object(path.Join(c.location, fileName)).Delete(ctx)

	return err
}

func (c *GcsBucketClient) getFile(fileName string) (*bytes.Buffer, string, error) {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	r, err := c.cl.Bucket(c.bucketName).Object(path.Join(c.location, fileName)).NewReader(ctx)

	if err != nil {
		return nil, "", err
	}

	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, r)

	if err != nil {
		return nil,"", err
	}

	mime := fmt.Sprintf("image/%s", filepath.Ext(fileName)[1:])

	return buf, mime, nil
}

func (c *GcsBucketClient) uploadFile(img *Image, fileName string) error {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	wc := c.cl.Bucket(c.bucketName).Object(path.Join(c.location, fileName)).NewWriter(ctx)

	switch img.mimeType {
	case "png":
		err := png.Encode(wc, img.img)
		if err != nil {
			return err
		}
	case "jpeg":
		err := jpeg.Encode(wc, img.img, &jpeg.Options{
			Quality: 100,
		})
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsuported file type %s", img.mimeType)
	}

	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}

	return nil
}
