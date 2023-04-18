package proxy

import (
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var logger = log.New(os.Stderr, "proxy:", log.Llongfile|log.LstdFlags)

func Serve(w http.ResponseWriter, r *http.Request) {
	logger.Println(r.Method, r.Host, r.URL)
	if r.Method == http.MethodConnect {
		handleHttps(w, r)
	} else {
		logger.Println("is http")
		handleHttp(w, r)
	}

}

func handleHttps(w http.ResponseWriter, r *http.Request) {

	destConn, err := net.DialTimeout("tcp", r.Host, 60*time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)

	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}

	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}

	go transfer(destConn, clientConn)
	go transfer(clientConn, destConn)

}

func handleHttp(w http.ResponseWriter, r *http.Request) {

	if strings.Contains(r.URL.String(), "110.65.10.252/cxxl/Index.aspx") {
		u, _ := url.Parse("http://110.65.10.252/cxxl/ShowNews.aspx?NewsNo=5747F92DECB987AD")
		r.URL = u
		//r.URL.Path = "/cxxl/ShowNews.aspx"
		//r.URL.Query().Set("NewsNo", "E03104625029843B")
		//r.RequestURI = r.URL.String() + "?NewsNo=E03104625029843B"
	}
	logger.Println(r.URL, r.RequestURI)
	resp, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	defer resp.Body.Close()

	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)

}

func transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()
	io.Copy(destination, source)
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
