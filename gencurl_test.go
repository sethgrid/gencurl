package gencurl

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestFromRequest(t *testing.T) {
	urlStr := "http://example.com"
	data := []byte(`{"key":"value"}`)
	body := bytes.NewBuffer(data)
	method := "POST"
	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		t.Fatal(err)
	}

	headerContentType := "Content-Type"
	contentType := "application/json"

	headerXCustom := "X-Custom"
	xCustom1 := `{"json":"data"}`
	xCustom2 := "more data"

	req.Header.Add(headerContentType, contentType)
	req.Header.Add(headerXCustom, xCustom1)
	req.Header.Add(headerXCustom, xCustom2)

	curl := FromRequest(req)
	t.Log("Generated Curl: " + curl)

	// be sure to capture your generated curl from your request before
	// executing the request if you want to capture the post body as the
	// execution of the request will drain the reader for the post body

	/*
		c := http.Client{}
		resp, err := c.Do(req)
		if err != nil {
			t.Fatalf("unable to process http request - %s", err)
		}
		defer resp.Body.Close()
	*/

	if want := fmt.Sprintf("-X %s", method); !strings.Contains(curl, want) {
		t.Errorf("missing ", want)
	}
	if want := fmt.Sprintf("--header '%s: %s'", headerContentType, contentType); !strings.Contains(curl, want) {
		t.Errorf("missing ", want)
	}
	if want := fmt.Sprintf("--header '%s: %s'", headerXCustom, xCustom1); !strings.Contains(curl, want) {
		t.Errorf("missing ", want)
	}
	if want := fmt.Sprintf("--header '%s: %s'", headerXCustom, xCustom2); !strings.Contains(curl, want) {
		t.Errorf("missing ", want)
	}
	if want := fmt.Sprintf("-d '%s'", string(data)); !strings.Contains(curl, want) {
		t.Errorf("missing ", want)
	}

	// Check the body was not emptied
	bytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		t.Errorf("expected no errors reading body, got %s", err)
	}
	if len(bytes) == 0 {
		t.Errorf("expected body to not be drained.")
	}
}

func TestFromParams(t *testing.T) {
	urlStr := "http://www.example.com"
	data := url.Values{"key": {"value"}}
	_, err := http.PostForm(urlStr, data)
	if err != nil {
		t.Fatal(err)
	}
	curl := FromParams("POST", urlStr, data.Encode(), http.Header{})
	t.Log(curl)

	if want := fmt.Sprintf("-X POST"); !strings.Contains(curl, want) {
		t.Errorf("missing %s", want)
	}
}

func TestFromParamsWithNoDataNoHeaders(t *testing.T) {
	urlStr := "http://www.example.com"
	curl := FromParams("GET", urlStr, "", nil)
	t.Log(curl)

	if want := fmt.Sprintf("-X GET"); !strings.Contains(curl, want) {
		t.Errorf("missing %s", want)
	}
}

func TestFromParamsWithHeaders(t *testing.T) {
	urlStr := "http://www.example.com"
	data := url.Values{"key": {"value"}}
	_, err := http.PostForm(urlStr, data)
	if err != nil {
		t.Fatal(err)
	}
	curl := FromParams("POST", urlStr, data.Encode(), http.Header{"Content-Type": []string{"application/json"}})
	t.Log(curl)

	if want := fmt.Sprintf("-X POST"); !strings.Contains(curl, want) {
		t.Errorf("missing %s", want)
	}
	if want := fmt.Sprintf("--header 'Content-Type: application/json'"); !strings.Contains(curl, want) {
		t.Errorf("missing %s", want)
	}
}
