package gencurl

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestSimpleCurl(t *testing.T) {
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

	// be sure to capture your generated curl from your request before
	// executing the request if you want to capture the post body as the
	// execution of the request will drain the reader for the post body

	// example executing of the request
	/*
		c := http.Client{}
		resp, err := c.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
	*/

	t.Log("Generated Curl: " + curl)

	if want := fmt.Sprintf("-X %s", method); !strings.Contains(curl, want) {
		t.Errorf("missing " + want)
	}
	if want := fmt.Sprintf("--header '%s: %s'", headerContentType, contentType); !strings.Contains(curl, want) {
		t.Errorf("missing " + want)
	}
	if want := fmt.Sprintf("--header '%s: %s'", headerXCustom, xCustom1); !strings.Contains(curl, want) {
		t.Errorf("missing " + want)
	}
	if want := fmt.Sprintf("--header '%s: %s'", headerXCustom, xCustom2); !strings.Contains(curl, want) {
		t.Errorf("missing " + want)
	}
	if want := fmt.Sprintf("-d '%s'", string(data)); !strings.Contains(curl, want) {
		t.Errorf("missing " + want)
	}
}
