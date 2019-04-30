#!/bin/bash

sudo mkdir /build
sudo chmod +666 /build
GO111MODULE=on go get -d ./...

go_build(){
local sf=""
for "windows" in $1 ;do sf= ".exe ";done
GO111MODULE=on GOOS=$1 GOARCH=$2 go build -o /build/$1_$2$sf $3
}

go_build "linux" "amd64" $1
go_build "linux" "386" $1
go_build "linux" "arm" $1
go_build "linux" "arm64" $1
go_build "windows" "amd64" $1
go_build "windows" "386" $1
go_build "darwin" "amd64" $1
go_build "darwin" "386" $1

#for i in /build/windows*;do mv $i $i.exe; done

#tar -cf /tmp/build.tar /build
#mv /tmp/build.tar /build
