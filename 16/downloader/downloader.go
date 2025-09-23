package downloader

import (
	"L2/16/parser"
	"L2/16/storage"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Config struct {
	BaseURL       string
	OutputDir     string
	MaxDepth      int
	Concurrency   int
	Timeout       time.Duration
	UserAgent     string
	RespectRobots bool
}

type Downloader struct {
	config     *Config
	client     *http.Client
	queue      chan *downloadTask
	visited    *storage.URLStorage
	wg         sync.WaitGroup
	statsMutex sync.Mutex
	robotsTxt  *RobotsTxt
}

type downloadTask struct {
	URL    string
	Depth  int
	IsPage bool
}

func NewDownloader(config *Config) *Downloader {
	client := &http.Client{
		Timeout: config.Timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return nil
		},
	}

	return &Downloader{
		config:    config,
		client:    client,
		queue:     make(chan *downloadTask, 1000),
		visited:   storage.NewURLStorage(),
		robotsTxt: NewRobotsTxt(),
	}
}

func (d *Downloader) Start() error {
	baseURL, err := url.Parse(d.config.BaseURL)
	if err != nil {
		return err
	}

	if d.config.RespectRobots {
		err = d.robotsTxt.Load(d.client, baseURL)
		if err != nil {
			return err
		}
	}

	for i := 0; i < d.config.Concurrency; i++ {
		d.wg.Add(1)
		go d.worker()
	}

	d.queue <- &downloadTask{
		URL:    d.config.BaseURL,
		Depth:  0,
		IsPage: true,
	}

	close(d.queue)
	d.wg.Wait()

	return nil
}

func (d *Downloader) worker() {
	defer d.wg.Done()

	for task := range d.queue {
		d.processTask(task)
	}
}

func (d *Downloader) processTask(task *downloadTask) {
	if d.visited.Has(task.URL) {
		return
	}
	d.visited.Add(task.URL)

	if d.config.RespectRobots && !d.robotsTxt.Allowed(task.URL) {
		log.Printf("Skipping %s (disallowed by robots.txt)", task.URL)
		return
	}

	if task.Depth > d.config.MaxDepth {
		return
	}

	content, contentType, err := d.downloadContent(task.URL)
	if err != nil {
		log.Printf("Failed to download %s: %v", task.URL, err)
		return
	}

	err = d.saveContent(task.URL, content)
	if err != nil {
		log.Printf("Failed to save %s: %v", task.URL, err)
		return
	}

	if task.IsPage && strings.Contains(contentType, "text/html") {
		d.extractLinks(task.URL, content, task.Depth)
	}
}

func (d *Downloader) downloadContent(urlStr string) ([]byte, string, error) {
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("User-Agent", d.config.UserAgent)

	resp, err := d.client.Do(req)
	if err != nil {
		return nil, "", err
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Printf("Error closing response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	return content, resp.Header.Get("Content-Type"), nil
}

func (d *Downloader) saveContent(urlStr string, content []byte) error {
	parseURL, err := url.Parse(urlStr)
	if err != nil {
		return err
	}

	localPath := d.urlToLocalPath(parseURL)
	fullPath := filepath.Join(d.config.OutputDir, localPath)

	dir := filepath.Dir(fullPath)
	if err = os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	if err = os.WriteFile(fullPath, content, 0644); err != nil {
		return err
	}

	return nil
}

func (d *Downloader) urlToLocalPath(u *url.URL) string {
	path := u.Path
	if path == "" || strings.HasSuffix(path, "/") {
		path += "index.html"
	}

	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}

	if path == "" {
		path = "index.html"
	}

	if u.Host != "" {
		path = filepath.Join(u.Host, path)
	}

	return filepath.Clean(path)
}

func (d *Downloader) extractLinks(baseURL string, content []byte, currentDepth int) {
	base, err := url.Parse(baseURL)
	if err != nil {
		log.Printf("Failed to parse base URL: %v", err)
		return
	}

	links := parser.ExtractLinks(content, base)
	for _, link := range links {
		if !d.isSameDomain(link) {
			continue
		}

		d.queue <- &downloadTask{
			URL:    link,
			Depth:  currentDepth + 1,
			IsPage: parser.IsPageURL(link),
		}
	}
}

func (d *Downloader) isSameDomain(link string) bool {
	linkURL, err := url.Parse(link)
	if err != nil {
		return false
	}

	baseURL, err := url.Parse(d.config.BaseURL)
	if err != nil {
		return false
	}

	return linkURL.Host == baseURL.Host || linkURL.Host == ""
}
