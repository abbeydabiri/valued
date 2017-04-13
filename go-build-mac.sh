#!/bin/bash
# go build  -o Valued.app -ldflags "-s -w" Valued.go && upx Valued.app && mv Valued.app app/.

go build  -o Valued.app -ldflags "-s -w" Valued.go && mv Valued.app app/.