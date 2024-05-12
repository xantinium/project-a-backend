#!/bin/bash

binaryName="platform-binary"

rm $binaryName
go build -o $binaryName main.go
