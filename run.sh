#!/bin/bash
### This script takes the sample input in sample.md
### and converts it to HTML
cat sample.md | go run ./cmd/gomarkdown/main.go > sample.html
