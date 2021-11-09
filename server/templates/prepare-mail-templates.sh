#!/bin/bash

docker run -it --rm -v $PWD:/app -w /app imolko/premailer-api ruby /app/prepare_email_templates.rb
