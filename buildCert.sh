#!/bin/bash

mkdir certs


openssl genrsa -out certs/server.key 2048

openssl ecparam -genkey -name secp384r1 -out certs/server.key

openssl req -new -x509 -sha256 -key certs/server.key -out certs/server.crt -days 3650