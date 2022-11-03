package work

import (
	"api_server/internal/infrastures/Image"
	"api_server/internal/repositories"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"
)

type Work struct {
	repo *repositories.Work
}

type CrawlerImagesResponse struct {
	Name        string
	Data        []byte
	ContentType string
}

func New(repo *repositories.Work) *Work {
	return &Work{
		repo: repo,
	}
}

//
//func (srv *Work) CrawlerImage(target string) ([]byte, error) {
//
//	targetUrl, err := url.Parse(target)
//
//	//if err != nil {
//	//	return nil, err
//	//}
//
//	t := &Image.Robot{}
//
//	imgs, err := t.Crawler(target)
//	if err != nil {
//		return nil, err
//	}
//
//	results := make(map[string][]byte)
//
//	for _, img := range imgs {
//		if strings.HasPrefix(img, "/") {
//			img = fmt.Sprintf("%s://%s/%s", targetUrl.Scheme, targetUrl.Host, img)
//		}
//		data, contentType, err := t.GetImage(img)
//		if err != nil {
//			return nil, err
//		}
//
//	}
//
//	// 1. 取得目標網站內容
//	content, err := srv.repo.GetHttpContent(target)
//	if err != nil {
//		return nil, err
//	}
//	sb := string(content)
//
//	//2. 在網站內容尋找image連結目標
//	r, err := regexp.Compile("<img.*src=\"(https:.*.jpg)\"")
//	if err != nil {
//		return nil, err
//	}
//	s := r.FindStringSubmatch(sb)
//	imgUrl := s[1]
//
//	//3. 透過img連結取得圖案
//	content, err = srv.repo.GetHttpContent(imgUrl)
//	if err != nil {
//		return nil, err
//	}
//	//os.WriteFile("picture.jpg", jpg, 0666)  存下圖檔
//	return content, nil
//}

func (srv *Work) CrawlerImage(target string) ([]byte, error) {
	// 1. 取得目標網站內容
	content, err := srv.repo.GetHttpContent(target)
	if err != nil {
		return nil, err
	}
	sb := string(content)

	//2. 在網站內容尋找image連結目標
	r, err := regexp.Compile("<img.*src=\"(https:.*.jpg)\"")
	if err != nil {
		return nil, err
	}
	s := r.FindStringSubmatch(sb)
	url := s[1]

	//3. 透過img連結取得圖案
	content, err = srv.repo.GetHttpContent(url)
	if err != nil {
		return nil, err
	}
	//os.WriteFile("picture.jpg", jpg, 0666)  存下圖檔
	return content, nil
}

func (srv *Work) CrawlerImagesAndSave(target string) ([]string, error) {
	targetUrl, err := url.Parse(target)
	var result []string
	if err != nil {
		return nil, err
	}

	dictName := fmt.Sprintf("%s/%s", "imgs", targetUrl.Host)
	err = os.MkdirAll(dictName, 0755)
	if err != nil {
		return nil, err
	}
	t := &Image.Robot{}
	imgs, err := t.Crawler(target)
	if err != nil {
		fmt.Println(err.Error())
	}
	for idx, img := range imgs {
		if strings.HasPrefix(img, "/") {
			img = fmt.Sprintf("%s://%s/%s", targetUrl.Scheme, targetUrl.Host, img)
		}
		fullName := fmt.Sprintf("%s/%d", dictName, idx)
		fullName, err = t.SaveImage(fullName, img)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		result = append(result, fullName)
	}
	return result, nil
}

func (srv *Work) CrawlerImages(target string) ([][]byte, error) {
	targetUrl, err := url.Parse(target)
	var result [][]byte
	if err != nil {
		return nil, err
	}

	dictName := fmt.Sprintf("%s/%s", "imgs", targetUrl.Host)
	err = os.MkdirAll(dictName, 0755)
	if err != nil {
		return nil, err
	}
	t := &Image.Robot{}
	imgs, err := t.Crawler(target)
	if err != nil {
		fmt.Println(err.Error())
	}

	for _, img := range imgs {
		if strings.HasPrefix(img, "/") {
			img = fmt.Sprintf("%s://%s/%s", targetUrl.Scheme, targetUrl.Host, img)
		}
		data, _, err := t.GetImage(img)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		result = append(result, data)
	}
	return result, nil
}
