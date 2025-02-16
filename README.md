# right-back-at-ya
A diagnostic HTTP server that writes everything about the request in the response

Example output:

```
curl http://localhost:8080

Method: GET
URL: /
Protocol: HTTP/1.1
Host: localhost:8080
Remote Address: [::1]:50703
Headers:
  User-Agent: [curl/8.7.1]
  Accept: [*/*]
HTTP Status code: 200
```

It also includes headers when they are present

```
curl -H "Foo: bar" http://localhost:8080

Method: GET
URL: /
Protocol: HTTP/1.1
Host: localhost:8080
Remote Address: [::1]:50704
Headers:
  Foo: [bar]
  User-Agent: [curl/8.7.1]
  Accept: [*/*]
HTTP Status code: 200
```

Requesting a path between `/1` to `/599` will return a request of
that status code. All other paths return a status code of `200`

```
curl http://localhost:8080/500

Method: GET
URL: /500
Protocol: HTTP/1.1
Host: localhost:8080
Remote Address: [::1]:50706
Headers:
  User-Agent: [curl/8.7.1]
  Accept: [*/*]
HTTP Status code: 500
```

By default rbay listens on port 8080 but that can be changed
with the environment variable `PORT`.

