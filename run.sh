#!/bin/bash

echo "==============="
echo -n "Building..."
go build -o proxy-server.exe
echo "Building done"
echo "Running..."
echo "==============="
./proxy-server.exe