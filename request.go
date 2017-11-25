package libtspeed

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/suconghou/utilgo"
)

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

func benchmark(url string, thunk uint, timeout uint, transport *http.Transport) error {
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
		_, err = io.Copy(ioutil.Discard, &writeCounter{origin: r, total: total, progress: utilgo.ProgressBar("", "", nil, nil)})
		return err
	}
	return fmt.Errorf("%d:%s Length %d", resp.StatusCode, resp.Status, resp.ContentLength)
}
