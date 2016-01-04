package datatypes

import (
	"sort"
	"sync"
)

//SitemapList represents a thread-safe list to keep the sitemap itens
type SitemapList struct {
	List map[string]*SitemapItem
	lock *sync.Mutex
}

//Search a specific item by location
func (smt *SitemapList) Search(location string) *SitemapItem {
	smt.lock.Lock()
	defer smt.lock.Unlock()
	if node, ok := smt.List[location]; ok {
		return node
	}
	return nil
}

//Add a new item
func (smt *SitemapList) Add(item *SitemapItem) {
	smt.lock.Lock()
	defer smt.lock.Unlock()
	smt.List[item.Location] = item
}

//CreateSitemap generates a Sitemap xml structure
func (sl *SitemapList) CreateSitemap() Sitemap {
	sl.lock.Lock()
	defer sl.lock.Unlock()
	arrayList := make(SiteMapItemList, len(sl.List))
	i := 0
	for _, v := range sl.List {
		arrayList[i] = v
		i++
	}
	sort.Sort(arrayList)
	return Sitemap{
		List:    arrayList,
		Xmlns:   "http://www.sitemaps.org/schemas/sitemap/0.9",
		XmlnsDO: "http://www.digitalocean.com/schemas/webcrawler",
	}
}

//NewSitemapList creates a new SitemapList instance
func NewSitemapList() *SitemapList {
	list := make(map[string]*SitemapItem)
	return &SitemapList{
		List: list,
		lock: &sync.Mutex{},
	}
}
