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
	"time"
)

func usage() {
	fmt.Fprintf(os.Stderr, "\nusage: go run serve.go [optional flags]\n")
	flag.PrintDefaults()
	fmt.Println()

	os.Exit(2)
}

type myHandler struct {
	certMap map[string]tls.Certificate
}

type myCert struct {
	cert      *tls.Certificate
	touchedAt time.Time
}

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
	fmt.Println()
	fmt.Println()

	// End the request
	// TODO serve from hosting directory
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
	defaultCert, err := tls.LoadX509KeyPair(certPath, privkeyPath)
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

	certMap := make(map[string]myCert)
	tlsConfig := new(tls.Config)
	tlsConfig.Certificates = []tls.Certificate{defaultCert}
	tlsConfig.GetCertificate = func(clientHello *tls.ClientHelloInfo) (*tls.Certificate, error) {

		// Load from memory
		// TODO unload untouched certificates every x minutes
		if myCert, ok := certMap[clientHello.ServerName]; ok {
			myCert.touchedAt = time.Now()
			return myCert.cert, nil
		}

		privkeyPath := filepath.Join(*certsPath, clientHello.ServerName, "privkey.pem")
		certPath := filepath.Join(*certsPath, clientHello.ServerName, "fullchain.pem")

		loadCert := func() *tls.Certificate {
			// TODO handle race condition (ask Matt)
			// the transaction is idempotent, however, so it shouldn't matter
			if _, err := os.Stat(privkeyPath); err == nil {
				fmt.Printf("Loading Certificates %s/%s/{privkey.pem,fullchain.pem}\n\n", *certsPath, clientHello.ServerName)
				cert, err := tls.LoadX509KeyPair(certPath, privkeyPath)
				if nil != err {
					return &cert
				}
				return nil
			}

			return nil
		}

		if cert := loadCert(); nil != cert {
			certMap[clientHello.ServerName] = myCert{
				cert:      cert,
				touchedAt: time.Now(),
			}
			return cert, nil
		}

		// TODO try to get cert via letsencrypt python client
		// TODO check for a hosting directory before attempting this
		/*
			cmd := exec.Command(
				"./venv/bin/letsencrypt",
				"--text",
				"--agree-eula",
				"--email", "coolaj86@gmail.com",
				"--authenticator", "standalone",
				"--domains", "www.example.com",
				"--domains", "example.com",
				"--dvsni-port", "65443",
				"auth",
			)
			err := cmd.Run()
			if nil != err {
				if cert := loadCert(); nil != cert {
					return cert, nil
				}
			}
		*/

		fmt.Fprintf(os.Stderr, "Failed to load certificates for %q.\n", clientHello.ServerName)
		fmt.Fprintf(os.Stderr, "\tTried %s/{privkey.pem,fullchain.pem}\n", filepath.Join(*certsPath, clientHello.ServerName))
		//fmt.Fprintf(os.Stderr, "\tand letsencrypt api\n")
		fmt.Fprintf(os.Stderr, "\n")
		// TODO how to prevent attack and still enable retry?
		// perhaps check DNS and hosting directory, wait 5 minutes?
		certMap[clientHello.ServerName] = myCert{
			cert:      &defaultCert,
			touchedAt: time.Now(),
		}
		return &defaultCert, nil
	}
	tlsListener := tls.NewListener(conn, tlsConfig)

	server := &http.Server{
		Addr:    addr,
		Handler: &myHandler{},
	}
	fmt.Printf("Listening on https://%s:%d\n\n", host, *port)
	server.Serve(tlsListener)
}
