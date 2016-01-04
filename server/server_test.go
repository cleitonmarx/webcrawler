package server

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/cleitonmarx/webcrawler/datatypes"
	"github.com/cleitonmarx/webcrawler/testutil"
)

func TestMainAdress_Return200OK(t *testing.T) {
	input := InputParam{t, "GET", "", nil, 0}
	expected := ExpectedParam{http.StatusOK, `{"status":200, "message": "You reached the webcrawler server. V:0.1.1"}`}
	execServerRequestTest(input, expected)
}

func TestNotFoundAddress_Return404NotFound(t *testing.T) {
	input := InputParam{t, "GET", "wrongpath", nil, 0}
	expected := ExpectedParam{http.StatusNotFound, `{"status":404, "message": "Not Found"}`}
	execServerRequestTest(input, expected)
}

func TestUseWrongVerb_Return405NotFound(t *testing.T) {
	input := InputParam{t, "POST", "", &url.Values{}, 0}
	expected := ExpectedParam{http.StatusMethodNotAllowed, `{"status":405, "message": "Method not allowed. Use this methods: GET"}`}
	execServerRequestTest(input, expected)
}

func TestCrawlerRequest_WithDepth6_Return200OKWith6Itens(t *testing.T) {
	input := InputParam{t, "POST", "crawler", &url.Values{"url": {"http://localhost/1"}, "depth": {"6"}}, 0}
	expected := ExpectedParam{http.StatusOK, ""}
	responseMessage := execServerRequestTest(input, expected)

	sitemap := datatypes.Sitemap{}
	err := xml.Unmarshal(responseMessage, &sitemap)
	handleTestError(t, err)

	if sitemap.List.Len() != 6 {
		t.Errorf("Expected 6, actual %d", sitemap.List.Len())
	}
}

func TestCrawlerRequest_WithDepth1_Return200OKWith3Itens(t *testing.T) {
	input := InputParam{t, "POST", "crawler", &url.Values{"url": {"http://localhost/1"}, "depth": {"1"}}, 0}
	expected := ExpectedParam{http.StatusOK, ""}
	responseMessage := execServerRequestTest(input, expected)

	sitemap := datatypes.Sitemap{}
	err := xml.Unmarshal(responseMessage, &sitemap)
	handleTestError(t, err)

	if sitemap.List.Len() != 3 {
		t.Errorf("Expected 3, actual %d", sitemap.List.Len())
	}
}

func TestCrawlerRequest_WithoutParameters_Return400BadRequest(t *testing.T) {
	input := InputParam{t, "POST", "crawler", &url.Values{}, 0}
	expected := ExpectedParam{http.StatusBadRequest, `{"status":400, "message": "The parameter url is missing"}`}
	execServerRequestTest(input, expected)

}

func TestCrawlerRequest_WithInvalidURL_Return400BadRequest(t *testing.T) {
	input := InputParam{t, "POST", "crawler", &url.Values{"url": {"$ˆ&%$&ˆ%#$#ˆ&%"}}, 0}
	expected := ExpectedParam{http.StatusBadRequest, `{"status":400, "message": "The parameter url isn't a valid URL"}`}
	execServerRequestTest(input, expected)
}

func TestCrawlerRequest_WithInvalidDepth_Return400BadRequest(t *testing.T) {
	input := InputParam{t, "POST", "crawler", &url.Values{"url": {"http://localhost/1"}, "depth": {"aaa"}}, 0}
	expected := ExpectedParam{http.StatusBadRequest, `{"status":400, "message": "The parameter depth isn't a valid integer"}`}
	execServerRequestTest(input, expected)
}

func TestCrawlerRequest_WithInvalidTimeout_Return400BadRequest(t *testing.T) {
	input := InputParam{t, "POST", "crawler", &url.Values{"url": {"http://localhost/1"}, "timeout": {"1ml"}}, 0}
	expected := ExpectedParam{http.StatusBadRequest, `{"status":400, "message": "The parameter timeout isn't a valid time. Examples: 5m, 10s"}`}
	execServerRequestTest(input, expected)
}

type InputParam struct {
	t            *testing.T
	method       string
	path         string
	parameters   *url.Values
	fetcherDelay time.Duration
}

type ExpectedParam struct {
	status  int
	message string
}

func execServerRequestTest(input InputParam, expected ExpectedParam) []byte {
	testServer := createNewTestServer(input.fetcherDelay)
	defer testServer.Close()

	var (
		response *http.Response
		err      error
	)
	if input.method == "POST" {
		response, err = http.PostForm(fmt.Sprintf("%s/%s", testServer.URL, input.path), *input.parameters)
	} else {
		response, err = http.Get(fmt.Sprintf("%s/%s", testServer.URL, input.path))
	}

	handleTestError(input.t, err)
	if response.StatusCode != expected.status {
		input.t.Errorf("Expected %d, actual %d", expected.status, response.StatusCode)
	}
	defer response.Body.Close()

	message, err := ioutil.ReadAll(response.Body)
	handleTestError(input.t, err)
	if len(expected.message) > 0 && expected.message != string(message) {
		input.t.Errorf("Expected %s, actual %s", expected.message, string(message))
	}

	return message
}

func createNewTestServer(fetcherDelay time.Duration) *httptest.Server {
	stubFetcher := &testutil.UrlFetcherStub{SleepTime: fetcherDelay}
	configRepoStub := &testutil.ConfigRepoStub{}
	crawlerServer, _ := New(configRepoStub, stubFetcher)
	crawlerServer.Init()
	testServer := httptest.NewServer(crawlerServer.HttpServer)
	return testServer
}

func handleTestError(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}
