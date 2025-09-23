package storage

import (
	"net/url"
	"sync"
)

type URLStorage struct {
	visited map[string]bool
	mutex   sync.RWMutex
}

func NewURLStorage() *URLStorage {
	return &URLStorage{
		visited: make(map[string]bool),
	}
}

func (u *URLStorage) Has(urlStr string) bool {
	norm := u.NormalizeURL(urlStr)

	u.mutex.RLock()
	defer u.mutex.RUnlock()
	return u.visited[norm]
}

func (u *URLStorage) Add(urlStr string) {
	norm := u.NormalizeURL(urlStr)

	u.mutex.Lock()
	defer u.mutex.Unlock()
	u.visited[norm] = true
}

func (s *URLStorage) NormalizeURL(urlStr string) string {
	urlNorm, err := url.Parse(urlStr)
	if err != nil {
		return urlStr
	}

	urlNorm.Fragment = ""

	if urlNorm.Host != "" {
		return urlNorm.Host + urlNorm.Path
	}

	return urlNorm.String()
}
