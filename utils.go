package httpclient

import (
	"net/url"
	"strings"
)

//http://google.com match http://google.com:80
//http://google.com match http://www.google.com:80
//http://google.com not match https://google.com
func match(url1 url.URL, url2 url.URL) bool {

	mapPort := map[string]string{"http": "80", "https": "443"}
	var (
		host1 string
		port1 string
		host2 string
		port2 string
	)

	host2Array1 := strings.Split(url1.Host, ":")
	host1 = host2Array1[0]
	if !strings.HasPrefix(host1, "www.") {
		host1 = "www." + host1
	}
	if len(host2Array1) == 1 {
		port1 = mapPort[url1.Scheme]
	} else {
		port1 = host2Array1[1]
	}

	host2Array2 := strings.Split(url2.Host, ":")
	host2 = host2Array2[0]
	if !strings.HasPrefix(host2, "www.") {
		host2 = "www." + host2
	}
	if len(host2Array2) == 1 {
		port2 = mapPort[url2.Scheme]
	} else {
		port2 = host2Array2[1]
	}

	if host1 == host2 && port1 == port2 {
		return true
	}
	return false
}
