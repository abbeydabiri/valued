#!/bin/bash
GOOS=linux GOARCH=386 go build  -o Valued.elf -ldflags "-s -w" Valued.go && upx Valued.elf && mv Valued.elf app/.
