package rest

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type HTTPClient struct {
	baseURL        string
	retryIntervals []time.Duration
}

func New(baseURL string, retryIntervals []time.Duration) HTTPClient {
	return HTTPClient{
		baseURL:        baseURL,
		retryIntervals: retryIntervals,
	}
}

func (h HTTPClient) SendPostRequest(url string, body bytes.Buffer) (*http.Response, error) {
	client := &http.Client{}
	url = h.baseURL + url
	if !strings.HasPrefix(url, "http") {
		url = fmt.Sprintf("http://%s", url)
	}
	log.Printf("Send request to url: %s", url)
	request, _ := http.NewRequest(http.MethodPost, url, &body)
	request.Header.Set("Content-Encoding", "gzip")
	request.Header.Set("Accept-Encoding", "gzip")
	request.Header.Set("Content-Type", "application/json")

	var resp *http.Response
	var err error
	for interval := range h.retryIntervals {
		resp, err = client.Do(request)

		if err == nil && resp.StatusCode == http.StatusOK {
			return resp, nil
		} else if resp != nil && resp.StatusCode == http.StatusRequestTimeout {
			log.Printf("Error to send request %s StatusCode: %d Error: %s, Sleep: %d", url, resp.StatusCode, err, interval)
			time.Sleep(time.Duration(interval))
			continue
		} else if resp != nil && err != nil {
			log.Printf("Error to send request %s StatusCode: %d Error: %s", url, resp.StatusCode, err)
			return nil, err
		} else {
			log.Printf("Error to send request %s Error: %s", url, err)
			return nil, err
		}
	}

	defer func() {
		if errBodyClose := resp.Body.Close(); errBodyClose != nil {
			err = errors.Join(err, fmt.Errorf("close body error: %w", errBodyClose))
		}
	}()
	return nil, err
}
