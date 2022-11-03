package work

import (
	"api_server/internal/repositories"
	"regexp"
)

type Work struct {
	repo *repositories.Work
}

func New(repo *repositories.Work) *Work {
	return &Work{
		repo: repo,
	}
}

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
