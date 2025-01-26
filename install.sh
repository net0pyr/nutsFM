#!/bin/bash

set -e

BINARY_NAME="nutsfm"

echo "Building the Go application..."
go build -o $BINARY_NAME

echo "Moving the binary to /usr/local/bin..."
sudo mv $BINARY_NAME /usr/local/bin/

echo "Making the binary executable..."
sudo chmod +x /usr/local/bin/$BINARY_NAME

echo "Installation complete. You can now run the file manager using the command: $BINARY_NAME"