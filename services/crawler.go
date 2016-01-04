package services

import (
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/cleitonmarx/webcrawler/datatypes"
)

//UrlFetcher is an interface for URL Fetching
type UrlFetcher interface {
	ParseUrl(newUrl string) ([]string, []string, error)
}

//Crawler represents a recursive url fetcher
type Crawler struct {
	mainURL  *url.URL
	LinkList *datatypes.SitemapList
	fetcher  UrlFetcher
	depth    int
	logger   *log.Logger
	timeout  time.Time
}

//GenerateSitemap fetch a url and generate a sitemap datastructure
func (c *Crawler) GenerateSitemap(urlToFetch *url.URL, timeout time.Duration) datatypes.Sitemap {
	c.mainURL = urlToFetch
	wg := &sync.WaitGroup{}
	wg.Add(1)
	c.timeout = time.Now().Add(timeout)
	c.fetchLink(urlToFetch.String(), c.depth, wg)
	wg.Wait()
	return c.LinkList.CreateSitemap()
}

//fetchLink is a recursive and concurrent function to a url and your links
func (c *Crawler) fetchLink(linkToFetch string, depth int, waitGroup *sync.WaitGroup) {
	defer func() {
		waitGroup.Done()
		c.logger.Printf("Done | Depth: %d | link: %s\n", depth, linkToFetch)
	}()

	if time.Now().After(c.timeout) {
		c.logger.Printf("Timeout | Depth: %d | link: %s\n", depth, linkToFetch)
		return
	}

	c.logger.Printf("Fetching | Depth: %d | link: %s\n", depth, linkToFetch)

	links, dependencies, err := c.fetcher.ParseUrl(linkToFetch)
	newItem := &datatypes.SitemapItem{
		Location: linkToFetch,
		Updated:  time.Now(),
	}
	c.LinkList.Add(newItem)

	if err != nil {
		newItem.ErrorInFetch = true
		newItem.ErrorMessage = err.Error()
		return
	}

	newItem.Links = links
	newItem.Assets = dependencies

	if depth <= 0 {
		return
	}

	for _, link := range links {
		if time.Now().After(c.timeout) {
			c.logger.Printf("Timeout | Depth: %d | link: %s\n", depth, linkToFetch)
			return
		}
		if foundLink := c.LinkList.Search(link); foundLink == nil && c.allowedLink(link) {
			waitGroup.Add(1)
			go c.fetchLink(link, depth-1, waitGroup)
			time.Sleep(100 * time.Millisecond)
		} else {
			c.logger.Printf("Skipped | Depth: %d - %s\n", depth, link)
		}
	}
}

//allowedLink allows just links with same domain, skipping subdomains or diferent hosts
func (c *Crawler) allowedLink(link string) bool {
	linkURL, err := url.Parse(link)
	if err == nil && linkURL.Host == c.mainURL.Host {
		return true
	}
	return false
}

//NewCrawler creates a new Crawler instance
func NewCrawler(urlFetcher UrlFetcher, depth int, logger *log.Logger) *Crawler {
	sitemapList := datatypes.NewSitemapList()

	return &Crawler{
		LinkList: sitemapList,
		fetcher:  urlFetcher,
		depth:    depth,
		logger:   logger,
	}
}
