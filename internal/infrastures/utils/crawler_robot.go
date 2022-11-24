package utils

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type Robot struct {
}

func (*Robot) Crawler(target string) ([]string, error) {
	resp, err := http.Get(target)
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(resp.Body)
	content := string(data)

	r, err := regexp.Compile("<img.*?src=[\"|'](.*?)[\"|']")

	if err != nil {
		return nil, err
	}

	s := r.FindAllStringSubmatch(content, -1)

	imgMap := make(map[string]bool)
	var imgs []string
	for _, target := range s {
		if imgMap[target[1]] {
			continue
		}
		imgMap[target[1]] = true
		imgs = append(imgs, target[1])
	}
	return imgs, nil
}

func (s *Robot) GetImage(source string) ([]byte, string, error) {
	var imageConverter IImgConverter

	if strings.Contains(source, "http") {
		imageConverter = &HttpUrlImgConverter{}

	} else if strings.Contains(source, "data:image") {
		imageConverter = &DataUriImgConverter{}

	} else {
		return nil, "", errors.New(fmt.Sprintf("not support this source: %s", source))
	}

	return imageConverter.Convert(source)
}

func (s *Robot) SaveImage(fullName string, source string) (string, error) {
	data, contentType, err := s.GetImage(source)

	ext, err := s.contentTypeConvert(contentType)

	if err != nil {
		return "", err
	}

	fullName = fmt.Sprintf("%s.%s", fullName, ext)
	if err := os.WriteFile(fullName, data, 0666); err != nil {
		return "", nil
	}
	return fullName, nil
}

func (*Robot) contentTypeConvert(contentType string) (string, error) {
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
		return "", errors.New("Not support this content type: " + contentType)
	}
	return ext, nil
}

func (*Robot) generateFileName(source string) (string, error) {
	reg, err := regexp.Compile("\\W*")
	if err != nil {
		return "", err
	}
	var fileName string
	sourceLen := len(source)
	if sourceLen < 6 {
		fileName = source
	} else {
		fileName = source[len(source)-6 : len(source)-1]
	}
	k := reg.ReplaceAllString(fileName, "")
	return k, nil
}
