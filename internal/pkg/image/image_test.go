package image

import (
	"github.com/stretchr/testify/assert"
	"image/png"
	"io/ioutil"
	"os"
	"testing"
)

func TestDecodeFromBase64(t *testing.T) {

	expected, err := os.Stat("./test-assets/img.png")
	assert.NoError(t, err)

	base64, err := ioutil.ReadFile("./test-assets/base64.txt")
	assert.NoError(t, err)

	actual, err := decodeFromBase64(string(base64))
	assert.NoError(t, err)

	file, err := ioutil.TempFile("/tmp", "test-image.*.png")
	assert.NoError(t, err)
	defer os.Remove(file.Name())

	err = png.Encode(file, actual.img)
	assert.NoError(t, err)

	_, err = ioutil.ReadFile(file.Name())
	assert.NoError(t, err)

	fileInfo, err := os.Stat(file.Name())
	assert.NoError(t, err)

	assert.Greater(t, fileInfo.Size(), int64(0))
	assert.Less(t, fileInfo.Size(), expected.Size())
}

func TestIsBase64ImageValidSize(t *testing.T){

	base64, err := ioutil.ReadFile("./test-assets/base64.txt")
	assert.NoError(t, err)

	assert.NoError(t, isBase64ImageValidSize(string(base64)))

}