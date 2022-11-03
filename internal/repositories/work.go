package repositories

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Work struct {
}

func NewWork() *Work {
	return &Work{}
}

func (*Work) GetHttpContent(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}
	_ = resp.Body.Close()
	return body, nil
}

func (*Work) SaveImage(data []byte, fileName string, contentType string) error {
	var ext string
	switch contentType {
	case "image/jpeg":
		ext = "jpeg"
	case "image/png":
		ext = "png"
	default:
		ext = ""
	}
	if ext == "" {
		return errors.New(fmt.Sprintf("not support this content type: %s", contentType))
	}
	fullName := fmt.Sprintf("%s/%s.%s", "imgs", fileName, ext)

	if err := os.WriteFile(fullName, data, 0666); err != nil {
		return err
	}
	return nil
}
