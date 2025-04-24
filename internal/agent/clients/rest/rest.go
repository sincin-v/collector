package rest

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type HTTPClient struct {
	baseURL string
}

func New(baseURL string) HTTPClient {
	return HTTPClient{baseURL: baseURL}
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
	resp, err := client.Do(request)
	defer func() {
		if errBodyClose := resp.Body.Close(); errBodyClose != nil {
			log.Printf("ERROR SEND DATA !!!!")
			err = errors.Join(err, fmt.Errorf("close body error: %w", errBodyClose))
		}
	}()

	if err != nil {
		log.Printf("Error to send request %s Error: %s", url, err)
		return nil, err
	}else if resp.StatusCode != http.StatusOK {
		log.Printf("Error to send request %s StatusCode: %d", url, resp.StatusCode)
		return nil, err
	}

	return resp, nil
}
