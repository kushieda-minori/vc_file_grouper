#!/bin/bash

#remove extra slashes from GOPATH
GOPATH=$(echo $GOPATH | tr -s /)
GOPATH=${GOPATH%/}

GOOS=darwin GOARCH=amd64 go build -gcflags -trimpath=$GOPATH -asmflags -trimpath=$GOPATH -o vc_file_grouper_OSX
GOOS=linux GOARCH=amd64 go build -gcflags -trimpath=$GOPATH -asmflags -trimpath=$GOPATH -o vc_file_grouper_Linux64
GOOS=linux GOARCH=386 go build -gcflags -trimpath=$GOPATH -asmflags -trimpath=$GOPATH -o vc_file_grouper_Linux32
GOOS=linux GOARCH=arm GOARM=7 go build -gcflags -trimpath=$GOPATH -asmflags -trimpath=$GOPATH -o vc_file_grouper_LinuxArm32v7l
GOOS=linux GOARCH=arm64 GOARM=7 go build -gcflags -trimpath=$GOPATH -asmflags -trimpath=$GOPATH -o vc_file_grouper_LinuxArm64v7l
GOOS=windows GOARCH=amd64 go build -gcflags -trimpath=$GOPATH -asmflags -trimpath=$GOPATH -o vc_file_grouper_Win64.exe
GOOS=windows GOARCH=386 go build -gcflags -trimpath=$GOPATH -asmflags -trimpath=$GOPATH -o vc_file_grouper_Win32.exe
GOOS=freebsd GOARCH=386 go build -gcflags -trimpath=$GOPATH -asmflags -trimpath=$GOPATH -o vc_file_grouper_freebsd32
GOOS=freebsd GOARCH=amd64 go build -gcflags -trimpath=$GOPATH -asmflags -trimpath=$GOPATH -o vc_file_grouper_freebsd64
