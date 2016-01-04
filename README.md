# webcrawler
[![Build Status](http://jenkins.cleitonmarx.svc.tutum.io/buildStatus/icon?job=webcrawler_integration&build=3)](http://jenkins.cleitonmarx.svc.tutum.io/job/webcrawler_integration/3/)  
Microservice that crawl webpages and generate a sitemap with links and pages assets. 

###Build
`go get github.com/cleitonmarx/webcrawler`  
`cd $GOPATH/src/github.com/cleitonmarx/webcrawler`  
make sure godep is installed, `go get github.com/tools/godep` and then build with  
`godep restore`  
`godep go build -a ./cmd/...`  
