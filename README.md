# gencurl
gencurl generates a curl command based on an http.Request to be used for logging and debugging

```go
// create an http request
data := []byte(`{"key":"value"}`)
req, err := http.NewRequest("POST", "http://www.example.com", bytes.NewBuffer(data))
if err != nil {
    // handle err
}
req.Header.Add("X-Custom", "custom data")

curl := gencurl.FromRequest(req)

// later, execute the request. On error, you can print curl to replicate and debug an issue
```

The generated curl command for this example would be:
`curl -v -X POST --header 'X-Custom: custom data'   http://www.example.com -d '{"key":"value"}'`

With this, you can test integrations and dig deeper. I suggest placing the generated curl in every error handling case dealing with an http request.
