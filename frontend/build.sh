#!/bin/bash

docker run --rm -it -v "$PWD":/frontend:delegated -w /frontend --entrypoint yarn node:current-alpine3.11 build --mode ${1-development}
