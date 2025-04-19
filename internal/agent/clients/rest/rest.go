package rest

import (
	"bytes"
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
	url = h.baseURL + url
	if !strings.HasPrefix(url, "http") {
		url = fmt.Sprintf("http://%s", url)
	}
	log.Printf("Send request to url: %s", url)
	resp, err := http.Post(url, "application/json", &body)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Printf("Error to send request %s Error: %s", url, err)
		return nil, err
	}
	defer resp.Body.Close()
	return resp, nil
}
