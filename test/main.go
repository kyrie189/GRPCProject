//package main
//
//import (
//	"bytes"
//	"fmt"
//	"net/http"
//	"sync"
//	"time"
//)
//
//func callAPI(url string, wg *sync.WaitGroup, ch chan<- string,payload string) {
//	defer wg.Done()
//	start := time.Now()
//	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJleHAiOjE2OTY1NjE4MzYsImlzcyI6IjM4Mzg0LVNlYXJjaEVuZ2luZSJ9._hGmPF5W9DBfR6WDv5O6APGztZNY-qgpXuS8Eto6XuY" // 替换为您的身份验证令牌
//	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
//	if err != nil {
//		fmt.Printf("Error creating request: %s\n", err.Error())
//		return
//	}
//	//req.Header.Set("Content-Type", "application/json")
//	req.Header.Set("Authorization", token)
//
//	client := http.Client{}
//	resp, err := client.Do(req)
//	if err != nil {
//		fmt.Printf("Error making POST request: %s\n", err.Error())
//		return
//	}
//	defer resp.Body.Close()
//
//	duration := time.Since(start)
//	ch <- fmt.Sprintf("API response time: %s", duration.String())
//}
//
//func main() {
//	url := "http://127.0.0.1:4000/api/v1/lottery/draw" // 替换为您的 API URL
//
//	concurrency := 5000000                // 并发请求数量
//
//	results := make(chan string, concurrency)
//	var wg sync.WaitGroup
//	for i := 0; i < concurrency; i++ {
//		wg.Add(1)
//		payload := fmt.Sprintf("{\"id\": %d}",i)
//		go callAPI(url, &wg, results,payload)
//	}
//	go func() {
//		wg.Wait()
//		close(results)
//	}()
//
//	for result := range results {
//		fmt.Println(result)
//	}
//}


package main

import (
"bytes"
"fmt"
"net/http"
"sync"
"time"
)

func callAPI(url string, wg *sync.WaitGroup, ch chan<- string, payload string, client *http.Client) {
	defer wg.Done()
	//start := time.Now()
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJleHAiOjE2OTY2NDk0MTgsImlzcyI6IjM4Mzg0LVNlYXJjaEVuZ2luZSJ9.LRGWryKfUZ999465UIJ3TuBT6pKi-FEWS_zMFYNe0NA" // 替换为您的身份验证令牌
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		ch <- fmt.Sprintf("Error creating request: %s", err.Error())
		return
	}
	req.Header.Set("Authorization", token)

	resp, err := client.Do(req)
	if err != nil {
		ch <- fmt.Sprintf("Error making POST request: %s", err.Error())
		return
	}
	defer resp.Body.Close()

	//duration := time.Since(start)

	//ch <- fmt.Sprintf("API response time: %s", duration.String())
}

func main() {
	url := "http://127.0.0.1:4000/api/v1/lottery/draw" // 替换为您的 API URL

	concurrency := 100000              // 适当降低并发请求数量
	maxIdleConnections := 1000        // 最大闲置连接数
	requestTimeout := 3000 * time.Second // 请求超时时间

	results := make(chan string, concurrency)
	var wg sync.WaitGroup

	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        maxIdleConnections,
			MaxIdleConnsPerHost: maxIdleConnections,
		},
		Timeout: requestTimeout,
	}

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		payload := fmt.Sprintf("{\"id\": %d}", i)
		go callAPI(url, &wg, results, payload, client)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		fmt.Println(result)
	}
}