#!/bin/bash
rm -rf build/*

env GOOS=windows GOARCH=amd64 go build -o build/hide_windows_amd64
env GOOS=windows GOARCH=386 go build -o build/hide_windows_386
env GOOS=darwin GOARCH=arm64 go build -o build/hide_darwin_arm64
env GOOS=darwin GOARCH=amd64 go build -o build/hide_darwin_amd64
env GOOS=linux GOARCH=amd64 go build -o build/hide_linux_amd64
env GOOS=linux GOARCH=arm64 go build -o build/hide_linux_arm64
env GOOS=linux GOARCH=arm go build -o build/hide_linux_arm
env GOOS=linux GOARCH=386 go build -o build/hide_linux_386