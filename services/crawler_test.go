package services

import (
	"log"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/cleitonmarx/webcrawler/testutil"
)

func TestGenerateSiteMap_Depth6_Result6links(t *testing.T) {
	newURL, _ := url.Parse("http://localhost/1")
	crawlerService := NewCrawler(&testutil.UrlFetcherStub{SleepTime: 0 * time.Millisecond}, 6, log.New(os.Stdout, "", 0))
	sitemap := crawlerService.GenerateSitemap(newURL, 1*time.Minute)
	lenghtList := len(sitemap.List)
	if lenghtList != 6 {
		t.Errorf("Expected 6 links, actual %d", lenghtList)
	}

}

func TestGenerateSiteMap_Depth1_Result3links(t *testing.T) {
	newURL, _ := url.Parse("http://localhost/1") //2, 4
	crawlerService := NewCrawler(&testutil.UrlFetcherStub{SleepTime: 0 * time.Millisecond}, 1, log.New(os.Stdout, "", 0))
	sitemap := crawlerService.GenerateSitemap(newURL, 1*time.Minute)
	lenghtList := len(sitemap.List)
	if lenghtList != 3 {
		t.Errorf("Expected 3 links, actual %d", lenghtList)
	}
}

func TestGenerateSiteMap_Timeout1s_Result1links(t *testing.T) {
	newURL, _ := url.Parse("http://localhost/1") //1
	crawlerService := NewCrawler(&testutil.UrlFetcherStub{SleepTime: 1 * time.Second}, 1, log.New(os.Stdout, "", 0))
	sitemap := crawlerService.GenerateSitemap(newURL, 1*time.Second)
	lenghtList := len(sitemap.List)
	if lenghtList != 1 {
		t.Errorf("Expected 1 links, actual %d", lenghtList)
	}
}
