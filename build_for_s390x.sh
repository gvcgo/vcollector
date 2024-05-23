#!/bin/sh

export GOOS="linux"
export GOARCH="s390x"
go build -o ./cmd/vco/vco ./cmd/vco/
