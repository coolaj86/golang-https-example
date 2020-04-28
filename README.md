# Moved

See <https://git.coolaj86.com/coolaj86/golang-https-example>

# golang-https-example

A TLS / SSL enabled WebServer in Go, done right - includes valid https certificates.

# Install

Install the server and some certificates

```bash
# Clone this repo
git clone ssh://gitea@git.coolaj86.com:22042/coolaj86/golang-https-example.git
pushd golang-https-example

# Clone some valid dummy certificates
git clone git@example.com:example/localhost.example.com-certificates.git \
  ./etc/letsencrypt/live/localhost.rootprojects.org/
```

# Test

Run the server

```bash
# Run the Code
go run serve.go \
  --port 8443 \
  --letsencrypt-path=./etc/letsencrypt/live/
```

View it in your browser

<https://localhost.rootprojects.org:8443>

Test it with `openssl`

```bash
openssl s_client -showcerts \
  -connect localhost:8443 \
  -servername localhost.rootprojects.org \
  -CAfile ./etc/letsencrypt/live/localhost.rootprojects.org/certs/ca/root.pem
```

Test it with `curl`

```bash
# should work
curl https://localhost.rootprojects.org:8443

# if the Root CA isn't in your bundle
curl https://localhost.rootprojects.org:8443 \
  --cacert=./etc/letsencrypt/live/localhost.rootprojects.org/certs/ca/root.pem
```
