
### Google へリダイレクト

```console
$ docker run --rm -p 10000:10000 -p 9901:9901 -v $(pwd)/envoy/01_google-redirect.yaml:/etc/envoy/envoy.yaml envoyproxy/envoy:v1.14.1
```

### built-in filters

[Network filters](https://www.envoyproxy.io/docs/envoy/latest/configuration/listeners/network_filters/network_filters#config-network-filters)

### HTTPS
TLS 終端、HTTPリダイレクト。  

`tls_context` は廃止されたっぽくてエラーになる。  
https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/security/ssl

```console
$ docker run --rm -p 8080:8080 -p 6443:6443 -p 9901:9901 -v $(pwd)/envoy/02_https.yaml:/etc/envoy/envoy.yaml -v $(pwd)/envoy/certs:/etc/certs envoyproxy/envoy:v1.14.1
```

http://localhost:9901/certs で確認できる。  

```console
$ curl -i -H "Host: example.com" http://localhost:8080
```

