# webcrawler
[![Build Status](http://jenkins.cleitonmarx.svc.tutum.io/buildStatus/icon?job=webcrawler_integration)](http://jenkins.cleitonmarx.svc.tutum.io/job/webcrawler_integration)
Microservice that crawl webpages and generate a sitemap with links and pages assets.

###Build
`go get github.com/cleitonmarx/webcrawler`  
`cd $GOPATH/src/github.com/cleitonmarx/webcrawler`  
make sure godep is installed, `go get github.com/tools/godep` and then build with  
`godep restore`  
`godep go build -a ./cmd/...`  

###Build Docker image  
`go get github.com/cleitonmarx/webcrawler`  
`cd $GOPATH/src/github.com/cleitonmarx/webcrawler`  
`docker build -t="webcrawler" .`  

###Run Docker image
`docker run -p 3333:3333 -d --name webcrawler webcrawler`

##How to use
####Get current version:
`curl -X GET http://127.0.0.1:3333/`

####Crawling a website:
`curl -X POST http://127.0.0.1:3333/crawler -d "url=http://www.digitalocean.com&depth=3&timeout=3s"`  

#####Parameters:
url - mandatory (e.g. http://www.digitalocean.com, http://www.londondrugs.com/shop/electronics)  
depth - optional - default: 5  
timeout - optional - default: 5m (e.g. 1s to one second, 4m / four minutes)   
