package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func usage() {
	fmt.Fprintf(os.Stderr, "\nusage: go run serve.go [optional flags]\n")
	flag.PrintDefaults()
	fmt.Println()

	os.Exit(2)
}

type myHandler struct{}

func (m *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Print debug info
	fmt.Println(r.Host)
	fmt.Println(r.Method)
	fmt.Println(r.RequestURI)
	fmt.Println(r.URL) // has many keys, such as Query
	for k, v := range r.Header {
		fmt.Println(k, v)
	}
	fmt.Println(r.Body)

	// End the request
	fmt.Fprintf(w, "Hi there, %s %q? Wow!\n\nWith Love,\n\t%s", r.Method, r.URL.Path[1:], r.Host)
}

func main() {
	flag.Usage = usage

	port := flag.Uint("port", 443, "https port")
	certsPath := flag.String("letsencrypt-path", "/etc/letsencrypt/live", "path at which an 'xyz.example.com' containing 'fullchain.pem' and 'privkey.pem' can be found")
	defaultHost := flag.String("default-hostname", "localhost.daplie.com", "the default folder to find certificates to use when no matches are found")

	flag.Parse()

	host := strings.ToLower(*defaultHost)
	// See https://groups.google.com/a/letsencrypt.org/forum/#!topic/ca-dev/l1Dd6jzWeu8
	/*
		if strings.HasPrefix("www.", host) {
			fmt.Println("TODO: 'www.' prefixed certs should be obtained for every 'example.com' domain.")
		}
		host = strings.TrimPrefix("www.", host)
	*/

	fmt.Printf("Loading Certificates %s/%s/{privkey.pem,fullchain.pem}\n", *certsPath, *defaultHost)
	privkeyPath := filepath.Join(*certsPath, *defaultHost, "privkey.pem")
	certPath := filepath.Join(*certsPath, *defaultHost, "fullchain.pem")
	cert, err := tls.LoadX509KeyPair(certPath, privkeyPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't load default certificates: %s\n", err)
		os.Exit(1)
	}

	addr := ":" + strconv.Itoa(int(*port))

	conn, err := net.Listen("tcp", addr)
	if nil != err {
		fmt.Fprintf(os.Stderr, "Couldn't bind to TCP socket %q: %s\n", addr, err)
		os.Exit(1)
	}

	tlsConfig := new(tls.Config)
	tlsConfig.Certificates = []tls.Certificate{cert}
	tlsConfig.GetCertificate = func(clientHello *tls.ClientHelloInfo) (*tls.Certificate, error) {
		return &cert, nil
	}
	tlsListener := tls.NewListener(conn, tlsConfig)

	server := &http.Server{
		Addr:    addr,
		Handler: &myHandler{},
	}
	fmt.Printf("Listening on https://%s:%d\n", host, *port)
	server.Serve(tlsListener)
}
