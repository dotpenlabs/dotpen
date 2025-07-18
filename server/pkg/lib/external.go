package lib

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

func UseProxy(origReq *http.Request) (*http.Response, error) {
	proxyBase := "https://proxy.bijsven.workers.dev"
	originalURL := origReq.URL.String()

	proxyURL := proxyBase + "?url=" + url.QueryEscape(originalURL)

	var body io.ReadCloser
	if origReq.Body != nil {
		body = origReq.Body
	}

	proxyReq, err := http.NewRequest(origReq.Method, proxyURL, body)
	if err != nil {
		return nil, err
	}

	if origReq.Header == nil {
		origReq.Header = make(http.Header)
	}
	proxyReq.Header = origReq.Header.Clone()

	if proxyReq.Header == nil {
		proxyReq.Header = make(http.Header)
	}

	proxyReq.Header.Set("x-auth-token", os.Getenv("API_KEY"))

	client := &http.Client{Timeout: 10 * time.Second}
	return client.Do(proxyReq)
}
