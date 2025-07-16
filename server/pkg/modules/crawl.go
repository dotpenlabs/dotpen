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

var c = cache.New(10*time.Minute, 30*time.Minute)

func UseCrawl(w http.ResponseWriter, r *http.Request, app core.App) {
	w.Header().Set("Content-Type", "application/json")
	defer func() {
		if rec := recover(); rec != nil {
			app.Logger().Warn("GET /api/crawl: Panic", "error", fmt.Sprint(rec))
			http.Error(w, `{"error":"Internal error"}`, 500)
		}
	}()

	u := r.URL.Query().Get("url")
	if u == "" {
		http.Error(w, `{"error":"URL required"}`, 400)
		return
	}
	if !strings.HasPrefix(u, "http") {
		u = "https://" + u
	}
	if v, f := c.Get(u); f {
		_ = json.NewEncoder(w).Encode(v)
		return
	}

	m, err := (func(u string) (*MetaData, error) {
		pu, err := url.Parse(u)
		if err != nil || pu.Host == "" {
			return nil, fmt.Errorf("invalid url")
		}
		h := strings.ToLower(strings.TrimPrefix(pu.Host, "www."))

		switch h {
		case "twitter.com", "x.com", "t.co":
			app.Logger().Debug("GET /api/crawl: Twitter crawler", "url", u)
			return UseTwitter(u)
		case "youtube.com", "m.youtube.com", "youtu.be":
			app.Logger().Debug("GET /api/crawl: YouTube crawler", "url", u)
			return UseYouTube(u)
		case "tiktok.com", "www.tiktok.com":
			app.Logger().Debug("GET /api/crawl: TikTok crawler", "url", u)
			return UseTikTok(u)
		default:
			app.Logger().Debug("GET /api/crawl: Default crawler", "url", u)
			return UseDefault(u)
		}
	})(u)

	// try last time with default (if used an specific crawler)
	if err != nil {
		app.Logger().Warn("GET /api/crawl: Overwriting to default crawler", "url", u, "error", err.Error())

		m, err = UseDefault(u)
		if err == nil {
			m.URL = u
			c.Set(u, m, cache.DefaultExpiration)
			_ = json.NewEncoder(w).Encode(m)
			return
		} else {
			app.Logger().Warn("GET /api/crawl: Crawl error (default)", "url", u, "error", err.Error())
			urlwoh, err := url.Parse(u)
			if err != nil {
				http.Error(w, `{"error":"failed"}`, 500)
				return
			}

			if m == nil {
				m = &MetaData{}
			}
			m.Title = urlwoh.Hostname()
			m.Favicon = urlwoh.Host + "/favicon.ico"
		}
	}

	m.URL = u
	c.Set(u, m, cache.DefaultExpiration)
	_ = json.NewEncoder(w).Encode(m)
}

func UseDefault(u string) (*MetaData, error) {
	req, _ := http.NewRequest("GET", u, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Connection", "keep-alive")

	cl := &http.Client{Timeout: 10 * time.Second}
	r, err := cl.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	if r.StatusCode >= 400 {
		b, _ := io.ReadAll(r.Body)
		return nil, fmt.Errorf("HTTP %d: %.100s", r.StatusCode, b)
	}

	b, _ := io.ReadAll(r.Body)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(b)))
	if err != nil {
		return nil, err
	}

	m := &MetaData{}
	set := func(f *string, v string) {
		if *f == "" && v != "" {
			*f = v
		}
	}

	doc.Find("meta").Each(func(_ int, s *goquery.Selection) {
		if n, ok := s.Attr("name"); ok {
			v, _ := s.Attr("content")
			switch strings.ToLower(n) {
			case "title":
				set(&m.Title, v)
			case "description":
				set(&m.Description, v)
			case "author":
				set(&m.Author, v)
			}
		}
		if p, ok := s.Attr("property"); ok {
			v, _ := s.Attr("content")
			switch strings.ToLower(p) {
			case "og:title":
				set(&m.Title, v)
			case "og:description":
				set(&m.Description, v)
			case "og:image":
				set(&m.Image, v)
			case "og:url":
				set(&m.URL, v)
			}
		}
	})

	doc.Find("link").EachWithBreak(func(_ int, s *goquery.Selection) bool {
		if r, _ := s.Attr("rel"); strings.Contains(strings.ToLower(r), "icon") {
			href, _ := s.Attr("href")
			if u2, err := url.Parse(href); err == nil {
				if !u2.IsAbs() {
					base, _ := url.Parse(u)
					href = base.ResolveReference(u2).String()
				}
				m.Favicon = href
			}
			return false
		}
		return true
	})

	doc.Find("script[type='application/ld+json']").Each(func(_ int, s *goquery.Selection) {
		var d interface{}
		if err := json.Unmarshal([]byte(s.Text()), &d); err != nil {
			return
		}
		var list []map[string]interface{}
		switch x := d.(type) {
		case map[string]interface{}:
			list = append(list, x)
		case []interface{}:
			for _, i := range x {
				if o, ok := i.(map[string]interface{}); ok {
					list = append(list, o)
				}
			}
		}
		for _, o := range list {
			if m.Title == "" {
				if v, ok := o["name"].(string); ok {
					m.Title = v
				}
			}
			if m.Description == "" {
				if v, ok := o["description"].(string); ok {
					m.Description = v
				}
			}
			if m.Image == "" {
				switch img := o["image"].(type) {
				case string:
					m.Image = img
				case map[string]interface{}:
					if v, ok := img["url"].(string); ok {
						m.Image = v
					}
				}
			}
			if m.Author == "" {
				switch a := o["author"].(type) {
				case string:
					m.Author = a
				case map[string]interface{}:
					if v, ok := a["name"].(string); ok {
						m.Author = v
					}
				}
			}
			if m.Date == "" {
				if v, ok := o["datePublished"].(string); ok {
					m.Date = v
				}
			}
		}
	})

	if m.Title == "" {
		m.Title = strings.TrimSpace(doc.Find("title").First().Text())
	}

	return m, nil
}

func UseTwitter(u string) (*MetaData, error) {
	pu, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	if strings.Contains(pu.Path, "/status/") {
		r, err := http.Get("https://publish.twitter.com/oembed?url=" + url.QueryEscape(u))
		if err != nil {
			return nil, err
		}
		defer r.Body.Close()

		var d struct {
			HTML, AuthorName string
		}
		if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
			return nil, err
		}

		txt := ""
		if m := regexp.MustCompile(`<p[^>]*>(.*?)</p>`).FindStringSubmatch(d.HTML); len(m) > 1 {
			txt = regexp.MustCompile(`<.*?>`).ReplaceAllString(m[1], "")
		}

		title := d.AuthorName + "'s tweet"
		if txt != "" {
			title = d.AuthorName + ": " + txt
		}

		return &MetaData{
			Title:  title,
			Author: d.AuthorName,
		}, nil
	}

	m, err := UseDefault(u)
	if err != nil || (m.Title == "" && m.Author == "" && m.Description == "") {
		parts := strings.Split(strings.Trim(pu.Path, "/"), "/")
		if len(parts) > 0 {
			username := parts[0]
			return &MetaData{
				Title:   "@" + username + " on X",
				Favicon: "https://unavatar.io/twitter/" + username,
			}, nil
		}
		return nil, fmt.Errorf("invalid X profile URL")
	}

	return m, nil
}

func UseYouTube(u string) (*MetaData, error) {
	m, err := UseDefault(u)
	if err != nil {
		return nil, err
	}

	pu, _ := url.Parse(u)
	h := strings.ToLower(strings.TrimPrefix(pu.Host, "www."))

	if m.Image == "" {
		id := ""
		if h == "youtu.be" {
			id = strings.Trim(pu.Path, "/")
		} else if pu.Query().Has("v") {
			id = pu.Query().Get("v")
		} else {
			parts := strings.Split(pu.Path, "/")
			for i := range parts {
				if parts[i] == "embed" && i+1 < len(parts) {
					id = parts[i+1]
					break
				}
			}
		}
		if id != "" {
			m.Image = "https://img.youtube.com/vi/" + id + "/maxresdefault.jpg"
		}
	}

	if m.Title == "" {
		if m.Author != "" {
			m.Title = m.Author + "'s video"
		} else {
			m.Title = "YouTube Video"
		}
	}

	return m, nil
}

func UseTikTok(u string) (*MetaData, error) {
	r, err := http.Get("https://www.tiktok.com/oembed?url=" + url.QueryEscape(u))
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	var d struct {
		AuthorName   string `json:"author_name"`
		Title        string `json:"title"`
		ThumbnailURL string `json:"thumbnail_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		return nil, err
	}

	return &MetaData{
		Title:  d.Title,
		Author: d.AuthorName,
		Image:  d.ThumbnailURL,
		URL:    u,
	}, nil
}
