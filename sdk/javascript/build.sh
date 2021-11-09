#!/bin/bash

docker run --rm -v "$PWD":/app -w /app -e NODE_ENV=$1 -e NODE_OPTIONS=--openssl-legacy-provider node:latest yarn build --mode ${1-development}

# Added
#       NODE_OPTIONS=--openssl-legacy-provider 
# since node switched to another hashing provider
# Reference: https://github.com/webpack/webpack/issues/14532