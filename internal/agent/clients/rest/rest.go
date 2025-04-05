package rest

import (
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

func (h HTTPClient) SendPostRequest(url string) (*http.Response, error) {
	url = h.baseURL + url
	if !strings.HasPrefix(url, "http") {
		url = fmt.Sprintf("http://%s", url)
	}
	log.Printf("Send request to url: %s", url)
	resp, err := http.Post(url, "text/plain", nil)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Fatalf("Error to send request %s", url)
		return nil, err
	}
	defer resp.Body.Close()
	return resp, nil
}
