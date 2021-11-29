---
title: "gRPC with SSL/TLS"
publishdate: 2021-12-05
categories: 
    - Golang
    - Programming
tags:
  - golang
  - grpc
resources:
    - name: header
    - src: featured.jpg
---

gRPC supports [authentication](https://grpc.io/docs/guides/auth/). Adding it to your project is simple. All you have to do is configure it with just a few lines of code. One of the authentication types that gRPC supports is SSL/TLS. From the server-side, the code looks like this:

```go
creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
if err != nil {
    // handle the error - no ignore it!
}
s := grpc.NewServer(grpc.Creds(creds))
```

The client has to update the code as shown below.

```go
creds, err := credentials.NewClientTLSFromFile(certFile, "")
if err != nil {
    // handle the error - no ignore it!
}
conn, _ := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(creds))
```

But, from where take the certificate? One of ways is using the openssl.

```sh
openssl genrsa -aes256 -passout pass:gsahdg -out server.pass.key 4096
openssl rsa -passin pass:gsahdg -in server.pass.key -out server.key
openssl req -new -key server.key -out server.csr

# and generate the certificate
openssl x509 -req -sha256 -days 365 -in server.csr -signkey server.key -out server.crt
```

This is what [Internet](https://devcenter.heroku.com/articles/ssl-certificate-self) says. However, you may end up with the following error.

```
transport: x509: certificate is not valid for any names, but wanted to match localhost:8070
```

The problem is you need the certificate authority. I found a better and simpler way of solving it. Of course, we can do it using openssl but I like very simple solutions with little room for mistakes. Another solution doesn't require openssl! All you have to do is install [certstrap](https://github.com/square/certstrap). Maybe it's not the safest thing in the world but for local development, it's IMO simpler than using openssl.

To generate the certificate authority, generate the certificate and sign it you can use following commands:

```sh
certstrap init --common-name "developer20.com"
certstrap request-cert -domain localhost
certstrap sign localhost --CA developer20.com
```

Your ready to use certificate is available in file `./out/localhost.crt` and the key file can be found in `./out/localhost.key`. That's it! You should be ready to connect successfully from your gRPC client to the server.
