package controllers

import (
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/cleitonmarx/webcrawler/services"
	"github.com/dimfeld/httptreemux"
)

type MainController struct {
	fetcher services.UrlFetcher
	logger  *log.Logger
}

func (mc *MainController) GetHandler(responseWriter http.ResponseWriter, request *http.Request, parameters map[string]string) {
	sendResponseMessage(responseWriter, 200, "You reached the webcrawler server. V:0.1.1")
}

func (mc *MainController) CrawlerHandler(responseWriter http.ResponseWriter, request *http.Request, parameters map[string]string) {
	request.ParseForm()
	urlParam := request.Form.Get("url")
	depthParam := request.Form.Get("depth")
	timeoutParam := request.Form.Get("timeout")

	newUrl, err := validateUrlParam(urlParam)
	if err != nil {
		sendResponseMessage(responseWriter, 400, err.Error())
		return
	}

	crawlerDepth, err := validateDepthParam(depthParam)
	if err != nil {
		sendResponseMessage(responseWriter, 400, err.Error())
		return
	}

	timeout, err := validateTimeoutParam(timeoutParam)
	if err != nil {
		sendResponseMessage(responseWriter, 400, err.Error())
		return
	}

	crawler := services.NewCrawler(mc.fetcher, crawlerDepth, mc.logger)
	sitemap := crawler.GenerateSitemap(newUrl, timeout)
	mc.logger.Printf("Finished | Total url fetched: %d | URL: %s\n", len(sitemap.List), urlParam)
	xmlOutput, _ := xml.Marshal(sitemap)
	responseWriter.Header().Add("Content-Type", "text/xml")
	responseWriter.Write([]byte(xml.Header))
	responseWriter.Write(xmlOutput)

}

//validateUrlParam validate the url parameter
func validateUrlParam(urlParam string) (*url.URL, error) {
	if len(urlParam) == 0 {
		return nil, errors.New("The parameter url is missing")
	}
	newUrl, err := url.Parse(urlParam)
	if err != nil {
		return nil, errors.New("The parameter url isn't a valid URL")
	}
	return newUrl, nil
}

//validateUrlParam validate the depth parameter
func validateDepthParam(depthParam string) (int, error) {
	//default crawler depth
	crawlerDepth := 5
	if len(depthParam) > 0 {
		temp, err := strconv.Atoi(depthParam)
		if err != nil {
			return 0, errors.New("The parameter depth isn't a valid integer")
		}
		if temp > 0 && temp <= 10 {
			crawlerDepth = temp
		}
	}

	return crawlerDepth, nil
}

//validateUrlParam validate the timeout parameter
func validateTimeoutParam(timeoutParam string) (time.Duration, error) {
	//default crawler timeout
	timeout := 5 * time.Minute
	if len(timeoutParam) > 0 {
		temp, err := time.ParseDuration(timeoutParam)
		if err != nil {
			return 0, errors.New("The parameter timeout isn't a valid time. Examples: 5m, 10s")
		}
		if temp > 1*time.Millisecond && temp <= 30*time.Minute {
			timeout = temp
		}
	}
	return timeout, nil
}

func (mc *MainController) NotFoundHandler(responseWriter http.ResponseWriter, request *http.Request) {
	sendResponseMessage(responseWriter, 404, "Not Found")
}

//MethodNotAllowedHandler handles requests with path exists but the method isn't allowed
func (mc *MainController) MethodNotAllowedHandler(
	responseWriter http.ResponseWriter,
	request *http.Request,
	methods map[string]httptreemux.HandlerFunc,
) {
	arrayMethods := make([]string, 0, len(methods))
	for key := range methods {
		arrayMethods = append(arrayMethods, key)
	}
	var message = fmt.Sprintf("Method not allowed. Use this methods: %s", strings.Join(arrayMethods, ","))
	sendResponseMessage(responseWriter, 405, message)
}

func sendResponseMessage(responseWriter http.ResponseWriter, status int, message string) {
	responseWriter.Header().Add("Content-Type", "application/json")
	responseWriter.WriteHeader(status)
	fmt.Fprintf(responseWriter, fmt.Sprintf(`{"status":%d, "message": "%s"}`, status, message))
}

func NewMainController(urlFetcher services.UrlFetcher, logger *log.Logger) MainController {
	return MainController{urlFetcher, logger}
}
