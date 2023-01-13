package utils

import (
	"golang.org/x/xerrors"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strings"
	"time"
)

func Post(url, reqBody string, retry int) (res []byte, err error) {
	for i := 0; i <= retry; i++ {
		// 首次访问为立即访问，若访问失败，则计算出一个随机值作为间隔，再次访问
		if i > 0 {
			// wait = i^2+[0-9]的一个随机数
			// 其中RandInt为[0,MaxInt64)的整数
			wait := math.Pow(float64(i), 2) + float64(RandInt()%10)
			log.Printf("retry after %f seconds\n", wait)
			time.Sleep(time.Duration(time.Duration(wait) * time.Second))
		}
		res, err = post(url, reqBody)
		if err == nil {
			return res, nil
		}
	}
	return nil, xerrors.Errorf("failed to fetch %s: %w", err)
}

func post(url string, reqBody string) ([]byte, error) {
	reqData := strings.NewReader(reqBody)
	req, err := http.NewRequest("POST", url, reqData)
	if err != nil {
		return nil, xerrors.Errorf("new request fail:%w\n", err)
	}

	req.Header.Add("Content-Type", "application/json")
	cli := http.Client{
		Timeout: time.Second * 5,
	}

	res, err := cli.Do(req)
	if err != nil {
		return nil, xerrors.Errorf("request fail:%w\n", err)
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			panic(err)
		}
	}(res.Body)

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, xerrors.Errorf("read response's body fail:%w\n", err)
	}
	return resBody, err
}
