package image

import (
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"strings"
)

type Image struct {
	img      image.Image
	mimeType string
}

type Saver interface {
	SaveBase64(base64 string) (string, error)
	DeleteFile(fileName string) error
	GetImage(filename string) (io.Reader, error)
}

func generateFileName(mimeType string) string {
	fileName := uuid.New().String()

	return fmt.Sprintf("%s.%s", fileName, mimeType)
}

func decodeFromBase64(b64 string) (*Image, error) {
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

func isBase64ImageValidSize(base64 string) error {

	imgSize := calcBase64OrigLength(base64)
	imgMax := 8 * 1024 * 1024 * 10

	if imgSize > imgMax {
		return fmt.Errorf("image size is %d and limit is %d", imgSize, imgMax)
	}

	return nil
}

func calcBase64OrigLength(base64 string) int {

	l := len(base64)

	// count how many trailing '=' there are (if any)
	eq := 0
	if l >= 2 {
		if base64[l-1] == '=' {
			eq++
		}
		if base64[l-2] == '=' {
			eq++
		}

		l -= eq
	}

	// basically:
	// eq == 0 :	bits-wasted = 0
	// eq == 1 :	bits-wasted = 2
	// eq == 2 :	bits-wasted = 4

	// so orig length ==  (l*6 - eq*2) / 8

	return (l*3 - eq) / 4
}
