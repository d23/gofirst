#!/bin/bash
CGO_ENABLED=1
/usr/bin/go build -a -ldflags '-s' gofirst.go
