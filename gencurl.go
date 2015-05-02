package gencurl

import (
	"bytes"
	"fmt"
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
		getRequestBody(r))

	return ret
}

func ifSet(condition string, passThrough string) string {
	if len(condition) == 0 {
		return ""
	}
	return passThrough
}

func getRequestBody(r *http.Request) string {
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return ""
	}

	// copy and replace the reader
	readerCopy := ioutil.NopCloser(bytes.NewBuffer(buf))
	readerReplace := ioutil.NopCloser(bytes.NewBuffer(buf))
	r.Body = readerReplace

	data, err := ioutil.ReadAll(readerCopy)
	if err != nil {
		return "<err reading copy>"
	}

	if len(data) == 0 {
		return ""
	}

	return fmt.Sprintf(" -d '%s'", string(data))
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
