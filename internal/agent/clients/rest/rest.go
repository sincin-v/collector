package rest

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type HttpClient struct {
	baseURL string
}

func New(baseURL string) HttpClient {
	return HttpClient{baseURL: baseURL}
}

func (h HttpClient) SendPostRequest(url string) (*http.Response, error) {
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
