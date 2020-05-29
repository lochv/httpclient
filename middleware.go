package httpclient

import (
	"crypto/tls"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

//var dNSPool = []string{"1.1.1.1", "1.0.0.1", "9.9.9.9", "8.8.8.8", "8.8.4.4", "208.67.222.222", "208.67.222.222"}

type redirects struct {
	Count       int
	Urls        []string
	StatusCodes []int
	Sizes       []int
}

type middleware struct {
	transport http.Transport
	redirects *redirects
	analyzer  bool
	maxRetry  int
	readTimeout time.Duration
}

func (m middleware) appendRedirect(url string, statuscode int, size int) {
	m.redirects.Count += 1
	m.redirects.Urls = append(m.redirects.Urls, url)
	m.redirects.StatusCodes = append(m.redirects.StatusCodes, statuscode)
	m.redirects.Sizes = append(m.redirects.Sizes, size)
}

func (m middleware) analyze(r *Response) {

	//TODO detect some security bugs

	return
}

func (m middleware) RoundTrip(req *http.Request) (resp *http.Response, err error) {

	transport := m.transport
	//proxyUrl, _ := url.Parse("http://127.0.0.1:8080")
	//transport.Proxy = http.ProxyURL(proxyUrl)
	transport.MaxIdleConns = 2000
	transport.MaxIdleConnsPerHost = 1000
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true, Renegotiation: tls.RenegotiateOnceAsClient}
	transport.DisableKeepAlives = true
	transport.DialContext = (&net.Dialer{
		Timeout:   10 * time.Second,
		KeepAlive: 1 * time.Second,

		//TODO fix memory leak
		//Resolver: &net.Resolver{
		//PreferGo: true,
		//Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
		//	d := net.Dialer{
		//		Timeout: time.Second * 10,
		//	}
		//	return d.DialContext(ctx, "udp", "8.8.8.8:53")
		//for {
		//	ctx,_ = context.WithDeadline(context.Background(), time.Now().Local().Add(time.Hour * time.Duration(0) +
		//		time.Minute * time.Duration(0) +
		//		time.Second * time.Duration(10)))
		//	conn, err := d.DialContext(ctx, "udp", DNSPool[rand.Intn(len(DNSPool))]+":53") //memory leak
		//	if err != nil {
		//		fmt.Println(err.Error())
		//		time.Sleep(time.Second)
		//		continue
		//	}
		//	return conn, nil
		//}
		//},
		//},
	}).DialContext

	transport.TLSHandshakeTimeout = 10 * time.Second
	transport.IdleConnTimeout = 30 * time.Second
	transport.ExpectContinueTimeout = 5 * time.Second
	transport.ResponseHeaderTimeout = 10 * time.Second
	transport.DisableCompression = true
	var retryed = 0
	for {
		resp, err = transport.RoundTrip(req)
		if err != nil {
			if retryed == m.maxRetry {
				return
			}
			retryed += 1
			continue
		} else {
			break
		}
	}

	//prevent redirect to other domain...
	if req.Response != nil && req.Response.Request != nil {
		if !match(*req.URL, *req.Response.Request.URL) {
			resp.Body.Close()
			return resp, errors.New(redirectOutMessage)
		}
	}

	if resp.StatusCode > 299 && resp.StatusCode < 400 {
		var body []byte
		readDone := make(chan int)
		go func() {
			body, _ = ioutil.ReadAll(io.LimitReader(resp.Body, maxBodySize))
			readDone <- 1
		}()

		select {
		case <-time.After(readTimeout):
			//deadline
		case <-readDone:
		}
		resp.Body.Close()
		size := len(body)
		m.appendRedirect(req.URL.String(), resp.StatusCode, size)
		resp.Body.Close()
	}

	return
}
