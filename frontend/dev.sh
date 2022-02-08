#!/bin/bash

docker run --rm -it -v "$PWD":/frontend:delegated -w /frontend -p 3000:3000 --entrypoint yarn node:current-alpine3.11 serve