package libtspeed

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/suconghou/utilgo"
)

const reqMethod = "GET"

type writeCounter struct {
	readed   int64
	total    int64
	progress func(received int64, readed int64, total int64, start int64, end int64)
	origin   io.Reader
}

func (wc *writeCounter) Read(p []byte) (int, error) {
	n, err := wc.origin.Read(p)
	wc.readed += int64(n)
	wc.progress(wc.readed, wc.readed, wc.total, 0, wc.total)
	return n, err
}

func reqWithHost(host string, url string, method string, reqHeader http.Header, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return req, err
	}
	req.Host = host
	return req, nil
}

func benchmark(url string, thunk uint, timeout uint, transport *http.Transport) error {
	resp, err := utilgo.Dohttp(url, reqMethod, nil, nil, timeout, transport)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	total := int64(thunk)
	if resp.StatusCode >= http.StatusOK && resp.StatusCode <= http.StatusIMUsed && resp.ContentLength > 0 && resp.ContentLength >= total {
		var r io.Reader
		if total == 0 {
			r = resp.Body
			total = resp.ContentLength
		} else {
			r = io.LimitReader(resp.Body, total)
		}
		_, err = io.Copy(ioutil.Discard, &writeCounter{origin: r, total: total, progress: utilgo.ProgressBar("", "", nil, nil)})
		return err
	}
	return fmt.Errorf("%d:%s Length %d", resp.StatusCode, resp.Status, resp.ContentLength)
}

func benchmarkIP(ip string, host string, path string, https bool, thunk uint, timeout uint, transport *http.Transport) error {
	client := utilgo.NewClient(timeout, transport)
	url := fmt.Sprintf("%s://%s%s", utilgo.BoolString(https, "https", "http"), ip, path)
	req, err := reqWithHost(host, url, reqMethod, nil, nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	total := int64(thunk)
	if resp.StatusCode >= http.StatusOK && resp.StatusCode <= http.StatusIMUsed && resp.ContentLength > 0 && resp.ContentLength >= total {
		var r io.Reader
		if total == 0 {
			r = resp.Body
			total = resp.ContentLength
		} else {
			r = io.LimitReader(resp.Body, total)
		}
		_, err = io.Copy(ioutil.Discard, &writeCounter{origin: r, total: total, progress: utilgo.ProgressBar("", "", nil, nil)})
		return err
	}
	return fmt.Errorf("%d:%s Length %d", resp.StatusCode, resp.Status, resp.ContentLength)
}
