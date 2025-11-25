#!/bin/bash

# This is tested with
#$ curl --version
# curl 8.7.1 (x86_64-apple-darwin24.0) libcurl/8.7.1 (SecureTransport) LibreSSL/3.3.6 zlib/1.2.12 nghttp2/1.
# 64.0
# Release-Date: 2024-03-27
# Protocols: dict file ftp ftps gopher gophers http https imap imaps ipfs ipns ldap ldaps mqtt pop3 pop3s rt
# sp smb smbs smtp smtps telnet tftp
# Features: alt-svc AsynchDNS GSS-API HSTS HTTP2 HTTPS-proxy IPv6 Kerberos Largefile libz MultiSSL NTLM SPNE
# GO SSL threadsafe UnixSockets
#
# $ awk --version
# awk version 20200816
#
# BSD and GNU versions of these tools may behave differently from one another.

# Create an order and extract ID
RESPONSE=$(curl -s -X POST localhost:8080/api/order \
  -H "Content-Type: application/json" \
  -d '{
    "couponCode": "",
    "items": [
      {"productId": "1", "quantity": 2}
    ]
  }')

echo $RESPONSE
echo "================"
ORDER_ID=$(echo "$RESPONSE"| awk -F'"' '/"ID":/ {print $4; exit}')

echo "Created order: $ORDER_ID"
echo "================"

# Fetch the order (and show headers)
curl -i localhost:8080/api/order/$ORDER_ID
