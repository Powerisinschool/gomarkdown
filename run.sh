#!/bin/bash
### This script takes the sample input in test_input.md
### and converts it to HTML
cat test_input.md | go run ./cmd/gomarkdown/main.go
