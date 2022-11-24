package utils

import (
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"regexp"
)

type IImgConverter interface {
	Convert(source string) ([]byte, string, error)
}

type DataUriImgConverter struct {
}

func (*DataUriImgConverter) Convert(source string) (content []byte, contentType string, err error) {
	r, err := regexp.Compile("data:(image/.*);base64,(.*)")
	if err != nil {
		return nil, "", err
	}
	s := r.FindAllStringSubmatch(source, -1)
	if len(s) <= 0 {
		return nil, "", errors.New("source is not data uri format")
	}
	contentType = s[0][1]
	stringContent := s[0][2]
	content, err = base64.StdEncoding.DecodeString(stringContent)
	if err != nil {
		return nil, "", errors.New("base64 decode failed")
	}
	return content, contentType, nil
}

type HttpUrlImgConverter struct {
}

func (*HttpUrlImgConverter) Convert(source string) (content []byte, contentType string, err error) {
	resp, err := http.Get(source)
	if err != nil {
		return nil, "", err
	}
	content, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	contentType = resp.Header["Content-Type"][0]
	defer resp.Body.Close()
	return content, contentType, nil
}
