package Image

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type IImgConverter interface {
	Convert(source string) ([]byte, string, error)
}

type DataUri struct {
}

func (*DataUri) Convert(source string) ([]byte, string, error) {
	r, err := regexp.Compile("data:(image/.*);base64,(.*)")
	if err != nil {
		return nil, "", err
	}
	s := r.FindAllStringSubmatch(source, -1)
	if len(s) <= 0 {
		return nil, "", errors.New("source is not data uri format")
	}
	contentType := s[0][1]
	data := s[0][2]
	content, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, "", errors.New("base64 decode failed")
	}
	return content, contentType, nil
}

type HttpUrl struct {
}

func (*HttpUrl) Convert(source string) ([]byte, string, error) {
	resp, err := http.Get(source)
	if err != nil {
		return nil, "", err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	contentType := resp.Header["Content-Type"][0]
	defer resp.Body.Close()
	return data, contentType, nil
}

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

func (s *Robot) SaveImage(fullName string, source string) error {
	var imageConverter IImgConverter

	if strings.Contains(source, "http") {
		imageConverter = &HttpUrl{}

	} else if strings.Contains(source, "data:image") {
		imageConverter = &DataUri{}

	} else {
		return errors.New(fmt.Sprintf("not support this source: %s", source))
	}

	//fileName, err := s.generateFileName(source)
	//if err != nil {
	//	return err
	//}

	data, contentType, err := imageConverter.Convert(source)

	ext, err := s.contentTypeConvert(contentType)

	if err != nil {
		return err
	}

	return os.WriteFile(fmt.Sprintf("%s.%s", fullName, ext), data, 0666)
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
