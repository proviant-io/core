package image

import (
	"fmt"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path"
)

func NewLocalSaver(location string) Saver {
	return &LocalSaver{
		location: location,
	}
}

type LocalSaver struct {
	location string
}

func (ls *LocalSaver) SaveBase64(base64 string) (string, error) {

	err := isBase64ImageValidSize(base64)

	if err != nil {
		return "", err
	}

	img, err := decodeFromBase64(base64)

	if err != nil {
		return "", fmt.Errorf("failed to parse image: %s", err.Error())
	}

	return ls.persist(*img)
}

func (ls *LocalSaver) GetImage(filename string) (io.Reader, error){

	return nil, nil
}

func (ls *LocalSaver) DeleteFile(fileName string) error {
	fullPath := path.Join(ls.location, fileName)

	return os.Remove(fullPath)
}

func (ls *LocalSaver) generateFileName(mimeType string) string {
	return path.Join(ls.location, generateFileName(mimeType))
}

func (ls *LocalSaver) persist(img Image) (string, error) {

	if _, err := os.Stat(ls.location); os.IsNotExist(err) {
		err := os.Mkdir(ls.location, 0644)
		if err != nil {
			return "", err
		}
	}

	fileName := ls.generateFileName(img.mimeType)

	f, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	defer f.Close()

	switch img.mimeType {
	case "png":
		err = png.Encode(f, img.img)
		if err != nil {
			return "", err
		}
	case "jpeg":
		err = jpeg.Encode(f, img.img, &jpeg.Options{
			Quality: 100,
		})
		if err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("unsuported file type %s", img.mimeType)
	}

	return fileName, nil
}
