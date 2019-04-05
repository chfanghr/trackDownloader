#!/bin/bash

sudo mkdir /build
sudo chmod +666 /build
GO111MODULE=on go get -v -d ./...

go_build(){
GO111MODULE=on GOOS=$1 GOARCH=$2 go build -o /build/$1_$2 $3
}

go_build "linux" "amd64" $1
go_build "linux" "i386" $1
go_build "linux" "arm" $1
go_build "linux" "arm64" $1
go_build "windows" "amd64" $1
go_build "windows" "i386" $1

tar -cf /tmp/build.tar /build
mv /tmp/build.tar /build