package libtspeed

import (
	"bufio"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/suconghou/utilgo"
)

// Log to print
var Log = log.New(os.Stdout, "", 0)

// Run http speed test
func Run(r io.Reader, thunk uint, timeout uint, transport *http.Transport) error {
	thunk = thunk * 1024
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		url := scanner.Text()
		if utilgo.IsURL(url) {
			Log.Print(url)
			err := benchmark(url, thunk, timeout, transport)
			if err != nil {
				Log.Print(err)
			} else {
				Log.Print("")
			}
		}
	}
	return nil
}
