golang-https-example
====================

A TLS /  SSL enabled WebServer in Go, done right - includes valid https certificates.

Install
=======

Install the server and some certificates

```bash
# Clone this repo
git clone git@github.com:coolaj86/golang-https-example.git
pushd golang-https-example

# Clone some valid dummy certificates
git clone git@github.com:Daplie/localhost.daplie.com-certificates.git \
  ./etc/letsencrypt/live/localhost.daplie.com/
```

Test
====

Run the server

```bash
# Run the Code
go run serve.go --port 8443 --letsencrypt-dir=./etc/letsencrypt/live/
```

View it in your browser

<https://localhost.daplie.com:8443>

Test it with `openssl`

```bash
openssl s_client -showcerts \
  -connect localhost:8443 \
  -servername localhost.daplie.com \
  -CAfile ./etc/letsencrypt/live/localhost.daplie.com/certs/ca/root.pem
```

Test it with `curl`

```bash
# should work
curl https://localhost.daplie.com:8443

# if the Root CA isn't in your bundle
curl https://localhost.daplie.com:8443 \
  --cacert=./etc/letsencrypt/live/localhost.daplie.com/certs/ca/root.pem
```
