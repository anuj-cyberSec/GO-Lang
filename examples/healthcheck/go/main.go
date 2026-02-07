package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

type result struct {
	url      string
	status   int
	duration time.Duration
	err      error
}

func main() {
	urls := os.Args[1:]
	if len(urls) == 0 {
		urls = []string{
			"https://example.com",
			"https://httpbin.org/status/200",
			"https://httpbin.org/status/503",
			"https://httpbin.org/delay/1",
			"https://httpbin.org/delay/2",
		}
	}

	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	const workers = 4
	jobs := make(chan string)
	results := make(chan result)

	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			for url := range jobs {
				results <- checkURL(client, url)
			}
		}()
	}

	go func() {
		for _, url := range urls {
			jobs <- url
		}
		close(jobs)
		wg.Wait()
		close(results)
	}()

	for res := range results {
		if res.err != nil {
			fmt.Printf("%-35s error=%v duration=%s\n", res.url, res.err, res.duration)
			continue
		}
		fmt.Printf("%-35s status=%d duration=%s\n", res.url, res.status, res.duration)
	}
}

func checkURL(client *http.Client, url string) result {
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return result{url: url, duration: time.Since(start), err: err}
	}

	resp, err := client.Do(req)
	if err != nil {
		return result{url: url, duration: time.Since(start), err: err}
	}
	defer resp.Body.Close()

	return result{
		url:      url,
		status:   resp.StatusCode,
		duration: time.Since(start),
	}
}
