package httpclient

import (
	"net/http"
	"net/url"
	"sync"
)

type jar struct {
	m      sync.Mutex
	perURL map[string][]*http.Cookie
}

func (j *jar) Cookies(u *url.URL) []*http.Cookie {
	j.m.Lock()
	defer j.m.Unlock()
	return j.perURL[u.Host]
}

func (j *jar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	j.m.Lock()
	defer j.m.Unlock()
	if j.perURL == nil {
		j.perURL = make(map[string][]*http.Cookie)
	}
	j.perURL[u.Host] = cookies
}
