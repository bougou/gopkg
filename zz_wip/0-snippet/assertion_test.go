package snippet

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

type AssertionSource string

const (
	StatusCodeSource   AssertionSource = "status_code"
	ResponseTimeSource AssertionSource = "response_time"
	ResponseSizeSource AssertionSource = "response_size"
	HeadersSource      AssertionSource = "headers"
	CookieSource       AssertionSource = "cookie"
	JsonBodySource     AssertionSource = "json_body"
	TextBodySource     AssertionSource = "text_body"
	XMLBodySource      AssertionSource = "xml_body"
)

type HeaderSourceProperty struct {
	Key         string // header name
	Value       string // header value
	Description string // header description
}

type CookieSourceProperty struct {
}

type JsonBodySourceProperty struct {
	JsonPath string
}

type TextBodySourceProperty struct {
	Text string
}

type XMLBodySourceProperty struct {
	XPath string
}

type Comparison struct {
}

var (
	httpRequestHeaders = []string{
		"WWW-Authenticate",
		"Authorization",
		"Proxy-Authenticate",
		"Proxy-Authorization",
	}

	httpResponseHeaders = []string{}

	httpContentTypes = []string{
		"application/octet-stream",
		"application/x-www-form-urlencoded",
		"application/json",
		"text/html",
		"text/plain",
		"text/css",
	}
)

type Assertion struct {
	Source     AssertionSource `json:"source"`
	Property   string
	Comparison string
	Value      string
	ToVariable string
}

func (a *Assertion) ExtractSourceValue() string {
	return ""
}

func (a *Assertion) Assert() bool {
	return true
}

type Capture struct {
	failed        chan bool
	failed_reason string
}

func (c *Capture) Errorf(fmtStr string, args ...interface{}) {
	c.failed <- true
	c.failed_reason = fmt.Sprintf(fmtStr, args...)
}

func (c *Capture) Logf(fmt string, args ...interface{}) {
}

func Test_Ca(t1 *testing.T) {
	t := Capture{failed: make(chan bool)}

	go func() {
		// create httpexpect instance
		e := httpexpect.New(&t, "http://127.0.0.1:8081")

		// is it working?
		// e.GET("/dial/plugins").
		// 	Expect().
		// 	Status(http.StatusOK).JSON().Array().Empty()

		resp := e.GET("/dial/plugins").
			Expect().
			Status(http.StatusOK).JSON()

		e.GET("/dial/plugins").
			Expect().
			Status(http.StatusOK).JSON().Object().ContainsKey("data")

		for _, d := range resp.Path("$.data").Array().Iter() {
			d.Object()
		}

		e.GET("/dial/plugins").
			Expect().
			Status(http.StatusOK).StatusRange(httpexpect.Status3xx)

		t.failed <- false
	}()

	select {
	case c := <-t.failed:
		if c {
			fmt.Printf("failed, %s\n", t.failed_reason)
		} else {
			fmt.Println("success")
		}
	}
}

func TestFruits(t *testing.T) {
	// create httpexpect instance
	e := httpexpect.New(t, "http://127.0.0.1:8081")

	// is it working?
	// e.GET("/dial/plugins").
	// 	Expect().
	// 	Status(http.StatusOK).JSON().Array().Empty()

	resp := e.GET("/dial/plugins").
		Expect().
		Status(http.StatusOK).JSON()

	e.GET("/dial/plugins").
		Expect().
		Status(http.StatusOK).JSON().Object().ContainsKey("data")

	for _, d := range resp.Path("$.data").Array().Iter() {
		d.Object()
	}

	e.GET("/dial/plugins").
		Expect().
		Status(http.StatusOK).StatusRange(httpexpect.Status2xx)

	x := e.GET("/dial/plugins").Expect()

	t.Log(x.Status(http.StatusOK))
	t.Log(x.RoundTripTime().Raw())
	t.Log(x.Body().Length().Raw())
	// t.Log(x.Body().Raw())
}
