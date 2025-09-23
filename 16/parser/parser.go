package parser

import (
	"net/url"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
)

func ExtractLinks(content []byte, baseURL *url.URL) []string {
	var links []string
	visited := make(map[string]bool)

	doc, err := html.Parse(strings.NewReader(string(content)))
	if err != nil {
		return links
	}

	var extract func(*html.Node)
	extract = func(n *html.Node) {
		if n.Type == html.ElementNode {
			var attr string
			switch n.Data {
			case "a", "link":
				attr = "href"
			case "img", "script", "iframe":
				attr = "src"
			case "source":
				attr = "srcset"
			}

			if attr != "" {
				for _, a := range n.Attr {
					if a.Key == attr {
						absoluteURL := resolveURL(a.Val, baseURL)
						if absoluteURL != "" && !visited[absoluteURL] {
							links = append(links, absoluteURL)
							visited[absoluteURL] = true
						}
					}
				}
			}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				extract(c)
			}
		}
	}
	extract(doc)
	return links
}

func resolveURL(link string, baseURL *url.URL) string {
	if strings.HasPrefix(link, "javascript:") || strings.HasPrefix(link, "mailto:") || strings.HasPrefix(link, "#") {
		return ""
	}

	parsed, err := url.Parse(link)
	if err != nil {
		return ""
	}

	resolved := baseURL.ResolveReference(parsed)
	return resolved.String()
}

func IsPageURL(urlStr string) bool {
	u, err := url.Parse(urlStr)
	if err != nil {
		return false
	}

	path := u.Path
	ext := strings.ToLower(filepath.Ext(path))

	resourceExts := map[string]bool{
		".css": true, ".js": true, ".png": true, ".jpg": true, ".jpeg": true,
		".gif": true, ".svg": true, ".ico": true, ".woff": true, ".woff2": true,
		".ttf": true, ".eot": true, ".pdf": true, ".zip": true, ".rar": true,
		".tar": true, ".gz": true, ".mp3": true, ".mp4": true, ".avi": true,
		".mov": true, ".wav": true,
	}

	return !resourceExts[ext] && !strings.Contains(path, ".")
}

func filepathExt(path string) string {
	for i := len(path) - 1; i >= 0 && path[i] != '/'; i-- {
		if path[i] == '.' {
			return path[i:]
		}
	}

	return ""
}
