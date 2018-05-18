package main

import (
	"crypto/tls"
	"net/http"
	"flag"
	"log"

	"golang.org/x/crypto/acme/autocert"
)


var (
	tlsDir = flag.String("tls", "cert","tls certificates dir")
	hostname = flag.String("hostname", "","hostname")
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

	server := &http.Server{
		Addr: ":443",
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
		},
	}

	go http.ListenAndServe(":80", certManager.HTTPHandler(nil))

	err := server.ListenAndServeTLS("", "") //key and cert are comming from Let's Encrypt
	log.Println(err)
}
