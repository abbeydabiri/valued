#!/bin/bash
GOOS=windows GOARCH=386 go build  -o App.exe -ldflags "-s -w" && upx App.exe && mv App.exe app/.
