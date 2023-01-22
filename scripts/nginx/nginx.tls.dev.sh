#!/bin/bash

CURRENTDIR=$(dirname "$0")

openssl genrsa -out $CURRENTDIR/dev_tls/nginx.key 2048
openssl req -new -x509 -key $CURRENTDIR/dev_tls/nginx.key \
    -addext "subjectAltName = DNS:queue-system.vip" \
    -out $CURRENTDIR/dev_tls/nginx.crt \
    -subj /C=TW/ST=Taipei/L=Taipei/O=Jerry0420/CN=queue-system.vip \
    -days 3650