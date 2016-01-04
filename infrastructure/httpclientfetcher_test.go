package infrastructure

import (
	"log"
	"net/url"
	"testing"
)

func TestParseUrl_DigitalOceanWebsite_RetreiveLinksAndAssets(t *testing.T) {
	parser := NewHttpClientFetcher()
	links, dep, _ := parser.ParseUrl("http://www.digitalocean.com")
	log.Println("--Links--")
	log.Println(links)
	log.Println("--Deps--")
	log.Println(dep)
	if len(links) == 0 {
		t.Error("No links found")
	}
	if len(dep) == 0 {
		t.Error("No assets found")
	}
}

func TestNormalizeLink_DenomalizedLinks_NormalizedLinks(t *testing.T) {
	mainUrl, _ := url.Parse("http://www.teste.com")
	normalizedLink := normalizeLink(mainUrl, "//www.teste.com/link")
	log.Println(normalizedLink)
	normalizedLink = normalizeLink(mainUrl, "/test/")
	log.Println(normalizedLink)
}
