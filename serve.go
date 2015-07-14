package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func usage() {
	fmt.Fprintf(os.Stderr, "\nusage: go run serve.go [optional flags]\n")
	flag.PrintDefaults()
	fmt.Println()

	os.Exit(2)
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
	fmt.Printf("Listening on https://%s:%d\n", host, *port)
}
