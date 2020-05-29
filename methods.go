package httpclient

import (
	"bytes"
	"net/http"
)

//TODO
//Add OPTIONS, PUT,...
func (c *Client) Get(url string, header map[string]string) (r Response) {
	c.preReq()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Response{Success: false}
	}
	if c.DisableUrlEncode {
		req.URL.Opaque = url[len(req.URL.Scheme)+len(`://`)+len(req.URL.Host):]
	}
	req.Header.Set("User-Agent", userAgent)
	for key, value := range header {
		req.Header.Set(key, value)
	}

	return c.req(req)
}

func (c *Client) Post(url string, header map[string]string, data []byte) (r Response) {
	c.preReq()

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return Response{Success: false}
	}
	if c.DisableUrlEncode {
		req.URL.Opaque = url[len(req.URL.Scheme)+len(`://`)+len(req.URL.Host):]
	}
	req.Header.Set("User-Agent", userAgent)
	for key, value := range header {
		req.Header.Set(key, value)
	}

	return c.req(req)
}
