#!/usr/bin/bash

GOOS=linux go build -ldflags="-s -w" -o build/deepsea-api local/deepsea/api
ls -slath build