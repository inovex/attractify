#!/bin/bash

docker run --rm -it -v "$PWD":/frontend:delegated -w /frontend -p 8080:8080 --entrypoint yarn node:current-alpine3.11 serve