#!/bin/bash

echo "==============="
echo -n "Building..."
go build -ldflags="-s -w" -o proxy-server.exe
echo "Building done."
echo "==============="
