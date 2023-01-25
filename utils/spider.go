package utils

import (
	"encoding/json"
	"golang.org/x/xerrors"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strings"
	"time"
)

//// FetchConcurrently 根据不同的reqBody并发爬取url
//func FetchConcurrently(url string, reqBodys []string, concurrency, wait, retry int) (responses [][]byte, err error) {
//	reqChan := make(chan string, len(reqBodys)) // 请求体的chan
//	resChan := make(chan []byte, len(reqBodys)) // 响应体的chan
//	errChan := make(chan error, len(reqBodys))  // 错误的chan
//	defer close(reqChan)
//	defer close(resChan)
//	defer close(errChan)
//
//	// 将请求体传入chan中
//	// TODO 为什么开并发
//	go func() {
//		for _, reqBody := range reqBodys {
//			reqChan <- reqBody
//		}
//	}()
//
//	// 开启并发任务
//	tasks := GenWorkers(concurrency, wait)
//
//	for range reqBodys {
//		tasks <- func() {
//			reqBody := <-reqChan
//			res, err := Post(url, reqBody, retry) //获取chan里的请求体，开始爬取
//			// 如果爬取不成功，则返回错误
//			if err != nil {
//				errChan <- err
//				return
//			}
//			// 如果爬取成功，返回响应体
//			resChan <- res
//		}
//	}
//
//	// 处理结果
//	var errs []error
//	timeout := time.After(10 * 60 * time.Second)
//	for range reqBodys {
//		select {
//		case res := <-resChan:
//			responses = append(responses, res)
//		case err := <-errChan:
//			errs = append(errs, err)
//		case <-timeout:
//			return nil, xerrors.New("Timeout Fetching URL")
//		}
//	}
//	if 0 < len(errs) {
//		return responses, fmt.Errorf("%s", errs)
//
//	}
//	return responses, nil
//}

// GenWorkers 将一个函数传入通道中，开个goroutine来执行
func GenWorkers(num, wait int) chan<- func() {
	tasks := make(chan func())
	for i := 0; i < num; i++ {
		go func() {
			for f := range tasks {
				f()
				time.Sleep(time.Duration(wait) * time.Second)
			}
		}()
	}
	return tasks
}

func Fetch(method, url string, reqBody any, retry int) (res []byte, err error) {
	for i := 0; i <= retry; i++ {
		// 首次访问为立即访问，若访问失败，则计算出一个随机值作为间隔，再次访问
		if i > 0 {
			// wait = i^2+[0-9]的一个随机数
			// 其中RandInt为[0,MaxInt64)的整数
			wait := math.Pow(float64(i), 2) + float64(RandInt()%10)
			log.Printf("retry after %f seconds\n", wait)
			time.Sleep(time.Duration(time.Duration(wait) * time.Second))
		}
		res, err = fetch(method, url, reqBody)
		if err == nil {
			return res, nil
		}
	}
	return nil, xerrors.Errorf("failed to fetch %s: %w", url, err)
}

func fetch(method, url string, reqBody any) ([]byte, error) {
	var (
		req *http.Request
		err error
	)
	if method == "POST" {
		reqData, err := json.Marshal(reqBody)
		if err != nil {
			return nil, xerrors.Errorf("fail to marshal reqBody:%w\n", err)
		}
		req, err = http.NewRequest("POST", url, strings.NewReader(string(reqData)))
		if err != nil {
			return nil, xerrors.Errorf("new request fail:%w\n", err)
		}
		req.Header.Add("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, xerrors.Errorf("new request fail:%w\n", err)
		}
	}
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
