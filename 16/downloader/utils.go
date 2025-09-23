package downloader

import (
	"log"
	"net/http"
	"net/url"
)

type RobotsTxt struct {
	rules map[string][]string
}

func NewRobotsTxt() *RobotsTxt {
	return &RobotsTxt{
		rules: make(map[string][]string),
	}
}

func (r *RobotsTxt) Load(client *http.Client, baseURL *url.URL) error {
	robotsURL := baseURL.Scheme + "://" + baseURL.Host + "/robots.txt"

	resp, err := client.Get(robotsURL)
	if err != nil {
		return err
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Println("Error closing response body")
		}
	}()

	return nil
}

func (r *RobotsTxt) Allowed(urlStr string) bool {
	return true
}
