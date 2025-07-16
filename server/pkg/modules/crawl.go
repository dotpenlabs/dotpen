package modules

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/patrickmn/go-cache"
	"github.com/pocketbase/pocketbase/core"
)

type MetaData struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Image       string `json:"image,omitempty"`
	Favicon     string `json:"favicon,omitempty"`
	URL         string `json:"url,omitempty"`
	Author      string `json:"author,omitempty"`
	Date        string `json:"date,omitempty"`
}

var metadataCache = cache.New(10*time.Minute, 30*time.Minute)

func UseCrawl(w http.ResponseWriter, r *http.Request, app core.App) {
	w.Header().Set("Content-Type", "application/json")

	defer func() {
		if rec := recover(); rec != nil {
			app.Logger().Warn("GET /api/crawl: Recovery triggered", "error", "panic occurred")
			http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		}
	}()

	targetURL := r.URL.Query().Get("url")
	if targetURL == "" {
		http.Error(w, `{"error": "URL query parameter is required"}`, http.StatusBadRequest)
		return
	}

	if !strings.HasPrefix(targetURL, "http://") && !strings.HasPrefix(targetURL, "https://") {
		targetURL = "https://" + targetURL
	}

	if cached, found := metadataCache.Get(targetURL); found {
		_ = json.NewEncoder(w).Encode(cached)
		return
	}

	parsedURL, err := url.Parse(targetURL)
	if err != nil || parsedURL.Host == "" {
		http.Error(w, `{"error": "Invalid URL"}`, http.StatusBadRequest)
		return
	}

	var meta *MetaData
	if isTwitterURL(parsedURL.Host) {
		meta, err = fetchTwitterMeta(targetURL)
	} else {
		meta, err = fetchGenericMeta(targetURL)
	}

	if err != nil {
		if strings.HasPrefix(targetURL, "https://") {
			altURL := "http://" + strings.TrimPrefix(targetURL, "https://")
			parsedURL, err2 := url.Parse(altURL)
			if err2 == nil && parsedURL.Host != "" {
				if isTwitterURL(parsedURL.Host) {
					meta, err = fetchTwitterMeta(altURL)
				} else {
					meta, err = fetchGenericMeta(altURL)
				}
				if err == nil {
					meta.URL = altURL
					metadataCache.Set(targetURL, meta, cache.DefaultExpiration)
					_ = json.NewEncoder(w).Encode(meta)
					return
				}
			}
		}

		app.Logger().Warn("GET /api/crawl: Failed to fetch or parse", "error", err.Error())
		http.Error(w, `{"error": "Failed to fetch or parse the URL"}`, http.StatusInternalServerError)
		return
	}

	meta.URL = targetURL
	metadataCache.Set(targetURL, meta, cache.DefaultExpiration)
	_ = json.NewEncoder(w).Encode(meta)
}

func isTwitterURL(host string) bool {
	host = strings.ToLower(strings.TrimPrefix(host, "www."))
	return host == "twitter.com" || host == "t.co" || host == "x.com"
}

func fetchTwitterMeta(tweetURL string) (*MetaData, error) {
	endpoint := "https://publish.twitter.com/oembed?url=" + url.QueryEscape(tweetURL)
	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected content-type %s; body: %.100s", contentType, body)
	}

	var result struct {
		HTML       string `json:"html"`
		AuthorName string `json:"author_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to decode Twitter JSON: %v (body: %.100s)", err, body)
	}

	re := regexp.MustCompile(`<p[^>]*>(.*?)</p>`)
	matches := re.FindStringSubmatch(result.HTML)
	text := ""
	if len(matches) > 1 {
		text = stripTags(matches[1])
	}

	title := fmt.Sprintf("%s: %s", result.AuthorName, text)
	if text == "" {
		title = fmt.Sprintf("%s's tweet", result.AuthorName)
	}

	return &MetaData{
		Title:  title,
		Author: result.AuthorName,
	}, nil
}

func fetchGenericMeta(targetURL string) (*MetaData, error) {
	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP status %d, body: %.200s", resp.StatusCode, body)
	}

	htmlBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(htmlBytes)))
	if err != nil {
		return nil, err
	}

	meta := &MetaData{}

	// Meta tags
	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		if name, exists := s.Attr("name"); exists {
			content, _ := s.Attr("content")
			switch strings.ToLower(name) {
			case "title":
				if meta.Title == "" {
					meta.Title = content
				}
			case "description":
				if meta.Description == "" {
					meta.Description = content
				}
			case "author":
				if meta.Author == "" {
					meta.Author = content
				}
			}
		}
		if property, exists := s.Attr("property"); exists {
			content, _ := s.Attr("content")
			switch strings.ToLower(property) {
			case "og:title":
				if meta.Title == "" {
					meta.Title = content
				}
			case "og:description":
				if meta.Description == "" {
					meta.Description = content
				}
			case "og:image":
				if meta.Image == "" {
					meta.Image = content
				}
			case "og:url":
				if meta.URL == "" {
					meta.URL = content
				}
			}
		}
	})

	// Favicon extraction
	if meta.Favicon == "" {
		var faviconHref string
		doc.Find("link").EachWithBreak(func(i int, s *goquery.Selection) bool {
			if rel, _ := s.Attr("rel"); strings.Contains(strings.ToLower(rel), "icon") {
				faviconHref, _ = s.Attr("href")
				return false // break
			}
			return true // continue
		})
		if faviconHref != "" {
			faviconURL, err := url.Parse(faviconHref)
			if err == nil && !faviconURL.IsAbs() {
				base, _ := url.Parse(targetURL)
				faviconHref = base.ResolveReference(faviconURL).String()
			}
			meta.Favicon = faviconHref
		}
	}

	// JSON-LD parsing
	doc.Find("script[type='application/ld+json']").Each(func(i int, s *goquery.Selection) {
		jsonText := s.Text()
		var data interface{}
		if err := json.Unmarshal([]byte(jsonText), &data); err != nil {
			return
		}
		switch v := data.(type) {
		case []interface{}:
			for _, item := range v {
				parseJSONLD(item, meta)
			}
		case map[string]interface{}:
			parseJSONLD(v, meta)
		}
	})

	// Title fallback
	if meta.Title == "" {
		meta.Title = strings.TrimSpace(doc.Find("title").Text())
	}

	return meta, nil
}

func parseJSONLD(data interface{}, meta *MetaData) {
	obj, ok := data.(map[string]interface{})
	if !ok {
		return
	}

	// Titel
	if meta.Title == "" {
		if name, ok := obj["name"].(string); ok {
			meta.Title = name
		}
	}

	// Beschrijving
	if meta.Description == "" {
		if desc, ok := obj["description"].(string); ok {
			meta.Description = desc
		}
	}

	// Afbeelding
	if meta.Image == "" {
		switch img := obj["image"].(type) {
		case string:
			meta.Image = img
		case map[string]interface{}:
			if urlStr, ok := img["url"].(string); ok {
				meta.Image = urlStr
			}
		}
	}

	// Auteur
	if meta.Author == "" {
		switch author := obj["author"].(type) {
		case string:
			meta.Author = author
		case map[string]interface{}:
			if name, ok := author["name"].(string); ok {
				meta.Author = name
			}
		}
	}

	// Publicatiedatum
	if meta.Date == "" {
		if date, ok := obj["datePublished"].(string); ok {
			meta.Date = date
		}
	}
}

func stripTags(html string) string {
	re := regexp.MustCompile(`<.*?>`)
	return re.ReplaceAllString(html, "")
}
