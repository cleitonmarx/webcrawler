package infrastructure

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/cleitonmarx/webcrawler/services"
)

//HttpClientFetcher implements a services.UrlFetcher
type HttpClientFetcher struct {
	httpClient *http.Client
}

//ParseUrl parse an url in links and assets
func (hcf *HttpClientFetcher) ParseUrl(newUrl string) ([]string, []string, error) {
	var (
		links, dependencies []string
		htmlDoc             *goquery.Document
		response            *http.Response
		err                 error
	)

	retry(3, 100*time.Millisecond, func() bool {
		response, err = hcf.httpClient.Get(newUrl)
		if err != nil {
			return false
		}
		htmlDoc, err = goquery.NewDocumentFromResponse(response)
		return err != nil
	})

	if err != nil {
		return []string{}, []string{}, err
	}

	dependencies = extractImages(htmlDoc)
	dependencies = append(dependencies, extractScripts(htmlDoc)...)
	dependencies = append(dependencies, extractCSSFiles(htmlDoc)...)
	links = extractLinks(htmlDoc)

	return links, dependencies, err

}

//retry try execute a routine for n times
func retry(numberRetries int, sleepInterval time.Duration, execFunc func() bool) {
	executionPassed := false
	executionCount := 0
	for !executionPassed && executionCount < numberRetries {
		executionPassed = execFunc()
		if !executionPassed {
			executionCount++
			time.Sleep(sleepInterval)
		}
	}
}

//extractImages returns all document images
func extractImages(htmlDoc *goquery.Document) []string {
	images := []string{}
	htmlDoc.Find("img").Each(func(i int, selected *goquery.Selection) {
		if source, ok := selected.Attr("src"); ok && len(source) > 0 {
			images = append(images, normalizeLink(htmlDoc.Url, source))
		}
	})
	return images
}

//extractCSSFiles returns all document css
func extractCSSFiles(htmlDoc *goquery.Document) []string {
	css := []string{}
	htmlDoc.Find("link").Each(func(i int, selected *goquery.Selection) {
		if href, ok := selected.Attr("href"); ok && len(href) > 0 {
			if typeLink, ok := selected.Attr("type"); ok && strings.Contains(typeLink, "css") {
				css = append(css, normalizeLink(htmlDoc.Url, href))
			}
		}
	})
	return css
}

//extractScripts returns all document js references
func extractScripts(htmlDoc *goquery.Document) []string {
	scripts := []string{}
	htmlDoc.Find("script").Each(func(i int, selected *goquery.Selection) {
		if source, ok := selected.Attr("src"); ok && len(source) > 0 {
			scripts = append(scripts, normalizeLink(htmlDoc.Url, source))
		}
	})
	return scripts
}

//extractLinks returns all document links
func extractLinks(htmlDoc *goquery.Document) []string {
	linkMap := make(map[string]bool)
	htmlDoc.Find("a").Each(func(i int, selected *goquery.Selection) {
		if href, ok := selected.Attr("href"); ok && len(href) > 0 {
			normLink := normalizeLink(htmlDoc.Url, href)
			if isAllowedLink(normLink) {
				if _, ok := linkMap[normLink]; !ok {
					linkMap[normLink] = true
				}
			}

		}
	})

	links := make([]string, 0, len(linkMap))
	for key := range linkMap {
		links = append(links, key)
	}

	return links
}

//isAllowedLink verify if a link is a valid link
func isAllowedLink(link string) bool {
	return !(strings.HasPrefix(link, "javascript") || strings.HasSuffix(link, ".zip") ||
		strings.HasSuffix(link, ".rar") || strings.HasSuffix(link, ".jpg") || strings.HasSuffix(link, ".gif") ||
		strings.HasSuffix(link, ".png") || strings.HasSuffix(link, ".eps"))
}

//normalizeLink normalize a link path
func normalizeLink(mainUrl *url.URL, link string) string {
	normalizedLink := link
	if strings.HasPrefix(link, "//") {
		normalizedLink = fmt.Sprintf("%s:%s", mainUrl.Scheme, link)
	} else if strings.HasPrefix(link, "/") {
		normalizedLink = fmt.Sprintf("%s://%s%s", mainUrl.Scheme, mainUrl.Host, link)
	} else if strings.HasPrefix(link, "#") {
		normalizedLink = fmt.Sprintf("%s://%s%s", mainUrl.Scheme, mainUrl.Host, link)
	}
	return normalizedLink

}

//NewHttpClientFetcher creates a new HttpClientFetcher
func NewHttpClientFetcher() services.UrlFetcher {
	timeout := 5 * time.Second
	httpClient := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(netw, addr, timeout)
				if err != nil {
					return nil, err
				}
				conn.SetDeadline(time.Now().Add(timeout))
				return conn, nil
			},
		},
	}

	return &HttpClientFetcher{
		httpClient: httpClient,
	}
}
