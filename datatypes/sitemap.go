package datatypes

import "time"

//Sitemap represents the main structure to render a sitemap xml file.
type Sitemap struct {
	XMLName []struct{}      `xml:"urlset"`
	List    SiteMapItemList `xml:"url"`
	Xmlns   string          `xml:"xmlns,attr"`
	XmlnsDO string          `xml:"xmlns:webcrawler,attr"`
}

//SitemapItem represents a SitemapList item
type SitemapItem struct {
	Location     string    `xml:"loc"`
	Links        []string  `xml:"webcrawler:links>webcrawler:link"`
	Assets       []string  `xml:"webcrawler:assets>webcrawler:asset"`
	ErrorInFetch bool      `xml:"webcrawler:error"`
	ErrorMessage string    `xml:"webcrawler:errormessage,omitempty"`
	Updated      time.Time `xml:"webcrawler:updated"`
}

//SiteMapItemList implements sort.Interface
type SiteMapItemList []*SitemapItem

func (s SiteMapItemList) Len() int {
	return len(s)
}

func (s SiteMapItemList) Less(i, j int) bool {
	return s[i].Updated.Before(s[j].Updated)
}

func (s SiteMapItemList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
