#!/bin/bash
GOOS=darwin GOARCH=amd64 go build -o vc_file_grouper_OSX
GOOS=linux GOARCH=amd64 go build -o vc_file_grouper_Linux64
GOOS=linux GOARCH=386 go build -o vc_file_grouper_Linux32
GOOS=windows GOARCH=amd64 go build -o vc_file_grouper_Win64.exe
GOOS=windows GOARCH=386 go build -o vc_file_grouper_Win32.exe
