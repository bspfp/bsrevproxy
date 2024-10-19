# bsrevproxy

## Usage

```bash
# bsrevproxy_linux [-f configfile] [-c]
# create sample config file
$ bsrevproxy_linux -c
# start app
$ bsrevproxy_linux
```

## Config file

```yaml
host: localhost # serve on localhost:9090
port: 9090
cert_file: ./cert.pem # http(certfile == keyfile == "") or https(with file)
key_file: ./key.pem
cors:
  allow_origin: "*"
default_redirect_url: https://www.example.com
static_dirs:
    request_path_prefix: /file
    local_path: ./static
    mime_type: text/plain
reverse_proxies:
  - request_path_prefix: /db/get
    target_url: http://localhost:10002/get
    timeoutsec: 3
redirects:
  - request_path_prefix: /abc
    target_url: https://other.example.com/def
    passsubpath: true
    passquery: true
```