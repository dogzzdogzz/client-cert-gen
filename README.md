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
client-cert-gen --mac 00:00:00:00:00:03

mac = 00:00:00:00:00:03
[2019-01-10 11:26:46] Creating 00-00-00-00-00-03/ folder ...
[2019-01-10 11:26:47] Creating 00-00-00-00-00-03/client.crt ...
[2019-01-10 11:26:47] Creating 00-00-00-00-00-03/client.crt.md5 ...
[2019-01-10 11:26:47] Creating 00-00-00-00-00-03/client.key ...
[2019-01-10 11:26:47] Creating 00-00-00-00-00-03/client.key.md5 ...
[2019-01-10 11:26:47] Done.
```
