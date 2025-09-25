package utilities

import (
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/briandowns/spinner"
)

func StartSpinner(message string) func() {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " " + message
	s.Start()
	return func() {
		s.Stop()
		fmt.Println(message + " done.")
	}
}

func TestServerLatency(host string) (time.Duration, error) {
	if !strings.HasPrefix(host, "http://") && !strings.HasPrefix(host, "https://") {
		host = "http://" + host
	}
	u, err := url.Parse(host)
	if err != nil {
		return 0, err
	}

	var totalLatency time.Duration
	attempts := 3

	for i := 0; i < attempts; i++ {
		start := time.Now()
		conn, err := net.DialTimeout("tcp", u.Host, 30*time.Second)
		if err != nil {
			return 0, err
		}
		conn.Close()
		latency := time.Since(start)
		totalLatency += latency
		time.Sleep(100 * time.Millisecond)
	}

	return totalLatency / time.Duration(attempts), nil
}

func NormalizedURL(rawurl string) string {
	if !strings.HasPrefix(rawurl, "http://") && !strings.HasPrefix(rawurl, "https://") {
		rawurl = "https://" + strings.TrimPrefix(rawurl, "//")
	}

	if !strings.HasSuffix(rawurl, "/") {
		rawurl += "/"
	}

	return rawurl
}
