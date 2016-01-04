package testutil

import (
	"errors"
	"time"
)

//FetcherStub implements a services.UrlFetcher interface
type UrlFetcherStub struct {
	SleepTime time.Duration
}

func (sf *UrlFetcherStub) ParseUrl(newUrl string) ([]string, []string, error) {
	time.Sleep(sf.SleepTime)
	switch newUrl {
	case "http://localhost/1":
		return []string{"http://localhost/2", "http://localhost/4"}, []string{"http://localhost/lib.js"}, nil
	case "http://localhost/2":
		return []string{"http://localhost/3"}, []string{"http://localhost/lib.js"}, nil
	case "http://localhost/3":
		return []string{"http://localhost/4"}, []string{"http://localhost/lib.js"}, nil
	case "http://localhost/4":
		return []string{"http://localhost/2", "http://localhost/3", "http://localhost/5"}, []string{"http://localhost/lib.js"}, nil
	case "http://localhost/5":
		return []string{"http://localhost/6"}, []string{"http://localhost/lib.js"}, nil
	case "http://localhost/6":
		return []string{}, []string{"http://localhost/lib.js"}, nil
	default:
		return []string{}, []string{}, errors.New("Page not found")
	}
}
