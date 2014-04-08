// Copyright 2011 Gary Burd
// Copyright 2013 Beego Web authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package models

import (
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/astaxie/beego"
)

var userAgent = "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/29.0.1541.0 Safari/537.36"

var (
	dialTimeout    = flag.Duration("dial_timeout", 10*time.Second, "Timeout for dialing an HTTP connection.")
	requestTimeout = flag.Duration("request_timeout", 20*time.Second, "Time out for roundtripping an HTTP request.")
)

func timeoutDial(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, *dialTimeout)
}

type transport struct {
	t http.Transport
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	timer := time.AfterFunc(*requestTimeout, func() {
		t.t.CancelRequest(req)
		beego.Warn("Canceled request for %s", req.URL)
	})
	defer timer.Stop()
	resp, err := t.t.RoundTrip(req)
	return resp, err
}

var (
	httpTransport = &transport{t: http.Transport{Dial: timeoutDial, ResponseHeaderTimeout: *requestTimeout / 2}}
	httpClient    = &http.Client{Transport: httpTransport}
)

func getHttpJson(url string, v interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		err = json.NewDecoder(resp.Body).Decode(v)
		if _, ok := err.(*json.SyntaxError); ok {
			return errors.New("JSON syntax error at " + url)
		}
		return nil
	}
	return errors.New("can't get infomation")
}

func getFiles(files []*rawFile) error {
	ch := make(chan error, len(files))
	for i := range files {
		go func(i int) {
			req, err := http.NewRequest("GET", files[i].rawURL, nil)
			if err != nil {
				ch <- err
				return
			}
			req.Header.Set("User-Agent", userAgent)
			resp, err := httpClient.Do(req)
			if err != nil {
				ch <- err
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode == 200 {
				p, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					ch <- err
					return
				}
				files[i].data = p
			}
			ch <- nil
		}(i)
	}
	for _ = range files {
		if err := <-ch; err != nil {
			return err
		}
	}
	return nil
}
