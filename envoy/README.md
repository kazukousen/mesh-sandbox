
### Google へリダイレクト

```console
$ docker run --rm -p 10000:10000 -p 9901:9901 -v $(pwd)/envoy/01_google-redirect.yaml:/etc/envoy/envoy.yaml envoyproxy/envoy:v1.14.1
```

