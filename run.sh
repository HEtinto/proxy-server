#!/bin/bash

echo "==============="
echo -n "Building..."
go build -o main.exe
echo "Building done"
echo "Running..."
echo "==============="
./main.exe