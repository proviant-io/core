package image

import (
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path"
	"strings"
)

// generate uniq name for image
// save image to storage

type Image struct {
	img      image.Image
	mimeType string
}

func DecodeFromBase64(b64 string) (*Image, error) {
	var err error

	commaIndex := strings.Index(b64, ",")

	imageType := strings.TrimSuffix(b64[5:commaIndex], ";base64")
	base64Image := b64[commaIndex+1:]

	imageDecoder := base64.NewDecoder(base64.StdEncoding, strings.NewReader(base64Image))

	var img image.Image
	var mimeType string

	switch imageType {
	case "image/png":
		img, err = png.Decode(imageDecoder)
		mimeType = "png"

		if err != nil {
			return nil, err
		}

	case "image/jpeg":
		img, err = jpeg.Decode(imageDecoder)
		mimeType = "jpeg"

		if err != nil {
			return nil, err
		}
	}
	return &Image{
		img:      img,
		mimeType: mimeType,
	}, nil
}

type Saver interface {
	Save(img image.Image) (string, error)
}

type LocalSaver struct {
	location string
}

func (ls *LocalSaver) generateFileName(mimeType string) string {
	fileName := uuid.New().String()

	return path.Join(ls.location, fmt.Sprintf("%s.%s", fileName, mimeType))
}

func (ls *LocalSaver) Save(img Image) (string, error) {

	fileName := ls.generateFileName(img.mimeType)

	f, err := os.Create(fileName)
	if err != nil {
		// Handle error
	}
	defer f.Close()

	switch img.mimeType {
	case "png":
		err = png.Encode(f, img.img)
		if err != nil {
			return "", err
		}
	case "jpg":
		err = jpeg.Encode(f, img.img, &jpeg.Options{
			Quality: 100,
		})
		if err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("unsuported file type")
	}

	return fileName, nil
}
