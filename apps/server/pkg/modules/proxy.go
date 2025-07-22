package modules

import (
	"io"
	"log"
	"net/http"
	"net/url"
)

func UseProxy(w http.ResponseWriter, r *http.Request) {
	targetURL := r.URL.Query().Get("url")
	if targetURL == "" {
		http.Error(w, "'url' param not found", http.StatusBadRequest)
		return
	}

	parsedURL, err := url.Parse(targetURL)
	if err != nil || !(parsedURL.Scheme == "http" || parsedURL.Scheme == "https") {
		http.Error(w, "Invalid URl, only http(s) allowed", http.StatusBadRequest)
		return
	}

	resp, err := http.Get(targetURL)
	if err != nil {
		http.Error(w, "Failed to get URL", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		log.Printf("Fout bij het kopiëren van response body: %v", err)
	}
}
