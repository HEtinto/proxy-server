#!/bin/bash

echo "==============="
echo -n "Building..."
go build -ldflags="-s -w" -o proxy.exe
echo "Building done."
echo "==============="
