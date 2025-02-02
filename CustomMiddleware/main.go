package main

import (
	"fmt"
	"net/http"
	"time"
)

type CustomRoundtriper struct {
	Proxied    http.RoundTripper
	ApiKey     string
	MaxRetries int
}

func (c *CustomRoundtriper) RoundTrip(req *http.Request) (*http.Response, error) {
	var resp *http.Response;
	var err error;
	for i := 0; i < c.MaxRetries; i++ {
		start := time.Now()
		req.Header.Set("Authorization", "Bearer "+c.ApiKey)
		resp, err = c.Proxied.RoundTrip(req)
		if err == nil && resp.StatusCode < 500 {
			fmt.Printf("Received response from %s in %v\n", req.URL, time.Since(start))
			return resp, err
		}
		fmt.Printf("Request to %s failed (attempt %d/%d), retrying...after %d seconds \n",req.URL.String() ,i+1, c.MaxRetries, 1 << i)
		time.Sleep(time.Duration(1 << i) * time.Second) // Exponential Backoff
	}
	return resp, err
}

func main() {
	client := http.Client{
		Transport: &CustomRoundtriper{
			Proxied: http.DefaultTransport, 
			ApiKey: "Some Random Shit",
			MaxRetries: 3,
		},
	}
	go func() {
		client.Get("https://www.google.com")
	}()
	client.Get("https://www.gogle.com")
	client.Get("https://www.google.com")
	client.Get("https://www.google.com")
	client.Get("https://www.google.com")
	client.Get("https://www.google.com")
}
