#!/bin/bash
set -e

if [ $# -eq 0 ] ; then
    echo "No version title provided"
    exit 1
fi

echo "Compiling for Linux..."
export GOOS=linux GOARCH=amd64
go build -o Anaxim

zip_name_linux="Anaxim-$1-linux.zip"
echo "Zipping to $zip_name_linux..."
zip -r "$zip_name_linux" Anaxim Maps/

echo "Compiling for Windows..."
export GOOS=windows GOARCH=amd64 CGO_ENABLED=1
# Using zig because mingw sucks, and so does windows
export CC="zig cc -target x86_64-windows"
export CXX="zig c++ -target x86_64-windows"
go build -o Anaxim.exe

zip_name_win="Anaxim-$1-windows.zip"
echo "Zipping to $zip_name_win..."
zip -r "$zip_name_win" Anaxim.exe Maps/
