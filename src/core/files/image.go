package core_files

import (
	"bytes"
	"errors"
	"fmt"
	"image/jpeg"
	"image/png"
	"os"

	"github.com/chai2010/webp"
)

func SaveImage(data []byte, filename string) error {
	imgBuffer := bytes.NewBuffer(data)

	img, err := jpeg.Decode(imgBuffer)
	if err != nil {
		imgBuffer.Reset()
		imgBuffer.Write(data)
		img, err = png.Decode(imgBuffer)
	}
	if err != nil {
		imgBuffer.Reset()
		imgBuffer.Write(data)
		img, err = webp.Decode(imgBuffer)
	}
	if err != nil {
		return errors.New("unknown image format")
	}

	file, err := os.Create(fmt.Sprintf("%s/%s.webp", os.Getenv("IMAGES_PATH"), filename))
	if err != nil {
		return err
	}

	defer file.Close()

	err = webp.Encode(file, img, &webp.Options{
		Quality:  100,
		Lossless: true,
		Exact:    true,
	})

	return err
}

func GetImage(filename string) ([]byte, error) {
	return os.ReadFile(fmt.Sprintf("%s/%s.webp", os.Getenv("IMAGES_PATH"), filename))
}

func DeleteImage(filename string) error {
	return os.Remove(filename)
}
