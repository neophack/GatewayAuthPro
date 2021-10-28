#!/bin/sh

cd frontend
yarn build
cd ./build
go-bindata -o=../../src/bindata/bindata.go  -pkg=bindata   -fs ./...

cd ../..

CGO_ENABLED=0 GOOS=linux GOARCH=arm go build ./src/main.go
mv main gatewayAuth_linux_arm
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ./src/main.go
mv main gatewayAuth_linux_amd64
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build ./src/main.go
mv main gatewayAuth_darwin_amd64