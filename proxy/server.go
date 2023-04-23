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
	// 启动协程转发数据
	go transfer(destConn, clientConn)
	go transfer(clientConn, destConn)

}

func handleHttp(w http.ResponseWriter, r *http.Request) {

	if strings.Contains(r.URL.String(), "http://www.lib.scut.edu.cn/2016/1025/c8738a127507/page.htm") {
		u, _ := url.Parse("http://www.lib.scut.edu.cn/2016/1025/c8738a127508/page.htm")
		r.URL = u
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
