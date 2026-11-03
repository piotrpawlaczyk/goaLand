package api

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	stdSync "sync"
)

const cooldownSeconds = 5
const cooldown = time.Duration(cooldownSeconds) * time.Second

type RateLimitedHTTPClient struct {
	httpClient   HttpClient
	blockedUntil time.Time
	mutex        stdSync.RWMutex
}

func newRateLimitedHTTPClient(timeout time.Duration) *RateLimitedHTTPClient {
	return &RateLimitedHTTPClient{
		httpClient:   &http.Client{Timeout: timeout},
		blockedUntil: time.Now(),
	}
}

func (c *RateLimitedHTTPClient) sendQuerry(req *http.Request) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if errors.Is(err, context.DeadlineExceeded) {
		log.Print("context deadline exceeded, trying again")
		return c.httpClient.Do(req)
	}
	return resp, err
}

func (c *RateLimitedHTTPClient) Do(req *http.Request) (*http.Response, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.wait()

	resp, err := c.sendQuerry(req)
	if err != nil {
		return nil, err
	}

	c.notify(resp)

	return resp, nil
}

func (c *RateLimitedHTTPClient) wait() {
	waitTime := time.Until(c.blockedUntil)
	if waitTime > 0 {
		log.Printf("wait %s seconds to next request", waitTime.Round(time.Second))
		time.Sleep(waitTime)
	}
}

func (c *RateLimitedHTTPClient) notify(resp *http.Response) {
	available, err := strconv.Atoi(resp.Header.Get("x-requests-available-minute"))

	if err != nil {
		log.Printf("Not found number of available requests")
		return
	}
	if available > 0 {
		return
	}

	resetSeconds, err := strconv.Atoi(resp.Header.Get("X-RequestCounter-Reset"))
	if err != nil {
		log.Printf("Not found time in sec to refresh counter")
		return
	}

	c.blockedUntil = time.Now().Add(time.Duration(resetSeconds)*time.Second + cooldown)
}
