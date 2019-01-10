<!--
title: 'client-cert-gen in go'
description: 'Generate client certificate from root CA, and write MAC address to UID subject of certificate'
layout: Doc
framework: v1
platform: AWS
language: nodeJS
authorName: 'Tony Lee'
-->
# client-cert-gen
Generate client certificate from root CA, and write MAC address to UID subject of certificate.

## Build
```bash
go build
```

## Prerequisites
- root CA filename must be `rootCA.crt` and `rootCA.key`.
- Only support PKCS#1 now.

## Run
```bash
$ client-cert-gen --mac 00:00:00:00:00:03

mac = 00:00:00:00:00:03
[2019-01-10 11:26:46] Creating 00-00-00-00-00-03/ folder ...
[2019-01-10 11:26:47] Creating 00-00-00-00-00-03/client.crt ...
[2019-01-10 11:26:47] Creating 00-00-00-00-00-03/client.crt.md5 ...
[2019-01-10 11:26:47] Creating 00-00-00-00-00-03/client.key ...
[2019-01-10 11:26:47] Creating 00-00-00-00-00-03/client.key.md5 ...
[2019-01-10 11:26:47] Done.
```

## Generated certificate
```bash
$ openssl x509 -in 00-00-00-00-00-03/client.crt -text -noout

Certificate:
    Data:
        Version: 3 (0x2)
        Serial Number: 2 (0x2)
    Signature Algorithm: sha256WithRSAEncryption
        Issuer: C=AU, ST=Some-State, O=Internet Widgits Pty Ltd
        Validity
            Not Before: Jan 10 03:26:47 2019 GMT
            Not After : Jan  7 03:26:47 2029 GMT
        Subject: O=abc, CN=*.example.com/UID=00:00:00:00:00:03
        Subject Public Key Info:
            Public Key Algorithm: rsaEncryption
                Public-Key: (2048 bit)
                Modulus:
                    00:b3:15:...
                    ...12:dd
                Exponent: 65537 (0x10001)
        X509v3 extensions:
            X509v3 Key Usage: critical
                Digital Signature
            X509v3 Extended Key Usage:
                TLS Web Client Authentication
            X509v3 Authority Key Identifier:
```