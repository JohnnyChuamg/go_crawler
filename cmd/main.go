package main

import (
	"api_server/internal/infrastures/Image"
	"fmt"
	"net/url"
	"os"
	"strings"
)

func main() {
	//const _target = "https://www.deviantart.com/wlop"
	//const _target = "https://tw.yahoo.com"
	const _target = "https://www.taiwan.net.tw/m1.aspx?sNo=0012076"
	targetUrl, err := url.Parse(_target)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	dictName := fmt.Sprintf("%s/%s", "imgs", targetUrl.Host)
	err = os.MkdirAll(dictName, 0755)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	t := &Image.Robot{}
	imgs, err := t.Crawler(_target)
	if err != nil {
		fmt.Println(err.Error())
	}
	for idx, img := range imgs {
		if strings.HasPrefix(img, "/") {
			img = fmt.Sprintf("%s://%s/%s", targetUrl.Scheme, targetUrl.Host, img)
		}
		if err := t.SaveImage(fmt.Sprintf("%s/%d", dictName, idx), img); err != nil {
			fmt.Println(err.Error())
		}
	}
	return
	//
	//defer func() {
	//	if r := recover(); r != nil {
	//		err, ok := r.(error)
	//		if !ok {
	//			err = fmt.Errorf("unknown error: %v", err)
	//		}
	//		log.Fatal().Msgf("%v", err)
	//		time.Sleep(3 * time.Second)
	//	}
	//}()
	//if err := server.NewServer(); err != nil {
	//	log.Fatal().Msgf("%v", err)
	//}
}
