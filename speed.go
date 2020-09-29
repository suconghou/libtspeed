package libtspeed

import (
	"bufio"
	"io"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/suconghou/utilgo"
)

// Log to print
var Log = log.New(os.Stdout, "", 0)

// Run http speed test
func Run(r io.Reader, chunk uint, timeout uint, transport *http.Transport) error {
	chunk = chunk * 1024
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		url := scanner.Text()
		if utilgo.IsURL(url, true) {
			Log.Print(url)
			err := benchmark(url, chunk, timeout, transport)
			if err != nil {
				Log.Print(err)
			} else {
				Log.Print("")
			}
		}
	}
	return nil
}

// RunHost do ip speed test
func RunHost(r io.Reader, host string, path string, https bool, chunk uint, timeout uint, transport *http.Transport) error {
	chunk = chunk * 1024
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		ip := scanner.Text()
		if net.ParseIP(ip) != nil {
			Log.Print(ip)
			err := benchmarkIP(ip, host, path, https, chunk, timeout, transport)
			if err != nil {
				Log.Print(err)
			} else {
				Log.Print("")
			}
		}
	}
	return nil
}
