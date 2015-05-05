package gencurl

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// FromRequest generates a curl command that can be used when debugging issues
// encountered while executing requests. Be sure to capture your curl before
// you execute the request if you want to capture the post body.
func FromRequest(r *http.Request) string {
	ret := fmt.Sprintf("curl -v -X %s %s %s %s %s %s",
		r.Method,
		getHeaders(r.Header),
		ifSet(r.UserAgent(), fmt.Sprintf("--user-agent '%s'", r.UserAgent())),
		ifSet(r.Referer(), fmt.Sprintf("--referrer '%s'", r.Referer())),
		r.URL.String(),
		getRequestBody(r.Body))

	return ret
}

// FromParams is less useful than FromRequest because the structure of the
// request is likely in the http package call. Unlike a request value, we
// cannot use a response value as we do not have access to the url, method,
// or request body. If you used http.Post(), add a content-type header set to
// the bodyType parameter. If you used http.PostForm(), your content-type is set
// to "application/x-www-form-urlencoded".
func FromParams(method string, urlStr string, requestBody string, headers http.Header) string {
	ret := fmt.Sprintf("curl -v -X %s %s %s %s",
		method,
		getHeaders(headers),
		urlStr,
		ifSet(requestBody, fmt.Sprintf("-d '%s'", requestBody)))

	return ret
}

func ifSet(condition string, passThrough string) string {
	if len(condition) == 0 {
		return ""
	}
	return passThrough
}

func getRequestBody(r io.ReadCloser) string {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return ""
	}

	// copy and replace the reader
	readerCopy := ioutil.NopCloser(bytes.NewBuffer(buf))
	readerReplace := ioutil.NopCloser(bytes.NewBuffer(buf))
	r = readerReplace

	data, err := ioutil.ReadAll(readerCopy)
	if err != nil {
		return "<err reading copy>"
	}

	if len(data) == 0 {
		return ""
	}

	return fmt.Sprintf("-d '%s'", string(data))
}

func getHeaders(h http.Header) string {
	ret := ""
	for header, values := range h {
		for _, value := range values {
			ret += fmt.Sprintf(" --header '%s: %v'", header, value)
		}
	}
	return ret
}
