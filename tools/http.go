package tools

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

// GetHTTPRequest http request get for monitoring
func GetHTTPRequest(url string) ([]byte, error) {
	client := &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: time.Second * 2,
				//KeepAlive: 10 * time.Second,
			}).Dial,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Request error[ %s ]", err.Error())
	}

	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("error http get request[ %s ][ %s ]", url, err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("failed to get request[ %d ][ %s ]", resp.StatusCode, resp.Body)
	}

	bytes, _ := ioutil.ReadAll(resp.Body)
	return bytes, nil
}
