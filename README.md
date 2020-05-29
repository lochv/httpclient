# httpclient

-------------------
Example:

```
package main

import (
	"fmt"
	"github.com/lochv/httpclient"
	"time"
)

func main() {
	c := httpclient.Client{
		Session:          false, // use session
		Following:        false, //follow redirects
		DisableUrlEncode: true,
		Analyze:          false,
		Retry:            0, //retry if request fail
		ReadTimeout:      30 * time.Second,
		MaxBodySize:      0, // bytes
	}
	r := c.Get("https://golang.org/#/api/liferay", nil)
	fmt.Println(r.StatusCode)

	c = httpclient.Client{
		Session:          false, // use session
		Following:        false, //follow redirects
		DisableUrlEncode: false,
		Analyze:          false,
		Retry:            0, //retry if request fail
		ReadTimeout:      30 * time.Second,
		MaxBodySize:      0, // bytes
	}
	r = c.Get("https://golang.org/#/api/liferay", nil)
	fmt.Println(r.StatusCode)
}
```