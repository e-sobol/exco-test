package fetcher

import (
	"context"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	defaultTimeoutMs = 800
	maxConcurrentReq = 10
)

func handleFetchUrlRequest(request UrlFetchRequest) []UrlFetchResult {
	globalTimeout := UnwrapPointerOrDefault(request.ExecutionTimeout, defaultTimeoutMs)
	responses := make([]UrlFetchResult, len(request.UrlRequests))

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(globalTimeout)*time.Millisecond)
	defer cancel()

	sem := make(chan struct{}, maxConcurrentReq)
	var wg sync.WaitGroup

	for i, urlRequest := range request.UrlRequests {
		wg.Add(1)
		go func() {
			defer wg.Done()

			select {
			case sem <- struct{}{}:
				defer func() { <-sem }()
				timeout := UnwrapPointerOrDefault(urlRequest.Timeout, globalTimeout)
				if timeout > globalTimeout {
					timeout = globalTimeout
				}
				reqCtx, reqCancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Millisecond)
				defer reqCancel()

				fetchResponse, err := fetchUrl(reqCtx, urlRequest.Url, timeout, urlRequest.Headers)
				if err != nil {
					log.Print(err.Error())
				}
				responses[i] = fetchResponse
			case <-ctx.Done():
				return
			}
		}()
	}
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-ctx.Done():
		log.Println("reached global execution timeout")
		cleanupAbortedResponses(responses, request.UrlRequests)
	}
	return responses
}

func fetchUrl(ctx context.Context, url string, timeout int, headers map[string]string) (UrlFetchResult, error) {
	client := &http.Client{Timeout: time.Duration(timeout) * time.Millisecond}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return UrlFetchResult{Url: url, Code: http.StatusInternalServerError, Error: err.Error()}, err
	}
	for key, value := range headers {
		req.Header.Add(key, value)
	}
	resp, err := client.Do(req)
	if err != nil {
		return UrlFetchResult{Url: url, Code: http.StatusInternalServerError, Error: err.Error()}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return UrlFetchResult{Url: url, Code: resp.StatusCode}, nil
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return UrlFetchResult{Url: url, Code: http.StatusInternalServerError, Error: err.Error()}, err
	}
	return UrlFetchResult{Url: url, Code: resp.StatusCode, Payload: string(body)}, nil
}

func cleanupAbortedResponses(responses []UrlFetchResult, requests []UrlRequest) {
	for i := range responses {
		if responses[i].Code == 0 && responses[i].Url == "" {
			responses[i].Url = requests[i].Url
			responses[i].Error = "Request aborted by timeout"
		}
	}
}
