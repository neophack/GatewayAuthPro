#!/bin/sh

OUTPUT_PATH=${OUTPUT_PATH:-output}

apt update && apt install sudo
sudo apt-get install wget

bash < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)
gvm install go1.16.3 -B
gvm use go1.16.3
wget https://github.com/creationix/nvm/archive/v0.34.0.tar.gz
tar -zxvf v0.34.0.tar.gz
source  nvm-0.34.0/nvm.sh
nvm install v14.15.4
nvm use v14.15.4

version=1_0_0

cd frontend
yarn build
cd ./build
go-bindata -o=../../src/bindata/bindata.go  -pkg=bindata   -fs ./...

cd ../..

CGO_ENABLED=0 GOOS=linux GOARCH=arm go build ./src/main.go
mv main $OUTPUT_PATH/gatewayAuth_linux_arm_$version
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ./src/main.go
mv main $OUTPUT_PATH/gatewayAuth_linux_amd64_$version
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build ./src/main.go
mv main $OUTPUT_PATH/gatewayAuth_darwin_amd64_$version



