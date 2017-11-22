package libtspeed

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/suconghou/utilgo"
)

type writeCounter struct {
	readed    int64
	total     int64
	startTime time.Time
	origin    io.Reader
}

func (wc *writeCounter) Read(p []byte) (int, error) {
	n, err := wc.origin.Read(p)
	wc.readed += int64(n)
	progress := utilgo.ProgressBar("", "", nil, nil)
	progress(wc.readed, wc.readed, wc.total, time.Since(wc.startTime).Seconds(), 0, wc.total)
	return n, err
}

func benchmark(url string, thunk uint, timeout uint, transport *http.Transport) error {
	startTime := time.Now()
	resp, err := utilgo.Dohttp(url, "GET", nil, nil, timeout, transport)
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
		_, err = io.Copy(ioutil.Discard, &writeCounter{origin: r, total: total, startTime: startTime})
		return err
	}
	return fmt.Errorf("%d:%s Length %d", resp.StatusCode, resp.Status, resp.ContentLength)
}
