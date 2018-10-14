package main

import (
	"bytes"
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"

	"golang.org/x/crypto/acme/autocert"
)

var (
	tlsDir   = flag.String("tls", "cert", "tls certificates dir")
	hostname = flag.String("hostname", "example.com", "hostname")

	httpsPort = flag.Int("https_port", 443, "https port")
	httpPort  = flag.Int("http_port", 80, "http port")

	backendHost = flag.String("backend", "localhost:8080", "reverse proxy backend")
)

func main() {
	flag.Parse()

	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(*hostname),
		Cache:      autocert.DirCache(*tlsDir),
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world"))
	})

	director := func(req *http.Request) {
		url := *req.URL
		url.Scheme = "http"
		url.Host = *backendHost

		var buffer []byte
		var err error
		if req.Body != nil {
			buffer, err = ioutil.ReadAll(req.Body)
			if err != nil {
				log.Fatal(err.Error())
			}
		}
		req2, err := http.NewRequest(req.Method, url.String(), bytes.NewBuffer(buffer))
		if err != nil {
			log.Printf("failed to make new request", err)
			return
		}
		req2.Header = req.Header
		*req = *req2
	}
	rp := &httputil.ReverseProxy{Director: director}

	server := &http.Server{
		Addr: fmt.Sprintf(":%d", *httpsPort),
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
		},
		Handler: rp,
	}

	go http.ListenAndServe(fmt.Sprintf(":%d", *httpPort), certManager.HTTPHandler(nil))

	err := server.ListenAndServeTLS("", "") //key and cert are comming from Let's Encrypt
	log.Println(err)
}
