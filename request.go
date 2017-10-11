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
	if err != nil {
		return n, err
	}
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
	if resp.StatusCode >= 200 && resp.StatusCode <= 209 && resp.ContentLength >= total {
		r := &writeCounter{origin: io.LimitReader(resp.Body, total), total: total, startTime: startTime}
		_, err = io.Copy(ioutil.Discard, r)
		return err
	}
	return fmt.Errorf("%d:%s Length %d", resp.StatusCode, resp.Status, resp.ContentLength)
}
