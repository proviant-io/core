package image

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"image/jpeg"
	"image/png"
	"io"
	"log"
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

	return filename, nil
}

func (gs *GcsSaver) GetImage(filePath string) (io.Reader, error){

	return gs.gcsBucketClient.getFile(filePath)
}

func (gs *GcsSaver) DeleteFile(fileName string) error {

	return nil
}

type GcsBucketClient struct {
	cl         *storage.Client
	projectID  string
	bucketName string
	location   string
}

func (c *GcsBucketClient) getFile(fileName string) (io.Reader, error) {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	log.Printf(c.location + fileName)

	return c.cl.Bucket(c.bucketName).Object(c.location + fileName).NewReader(ctx)
}

func (c *GcsBucketClient) uploadFile(img *Image, object string) error {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Upload an object with storage.Writer.
	wc := c.cl.Bucket(c.bucketName).Object(c.location + object).NewWriter(ctx)

	log.Printf("file: %s\n", object)

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
