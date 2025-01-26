#!/bin/bash

set -e

BINARY_NAME="nutsfm"
INSTALL_DIR="/usr/local/bin"
LOG_FILE="$HOME/.nutsfm/app.log"
HISTORY_FILE="$HOME/.nutsfm/commands_history"

mkdir -p "$HOME/.nutsfm"

echo "Building the Go application..."
go build -o $BINARY_NAME

echo "Moving the binary to $INSTALL_DIR..."
sudo mv $BINARY_NAME $INSTALL_DIR/

echo "Making the binary executable..."
sudo chmod +x $INSTALL_DIR/$BINARY_NAME

echo "Creating log file..."
touch $LOG_FILE

echo "Creating history file..."
touch $HISTORY_FILE

echo "Installation complete. You can now run the file manager using the command: $BINARY_NAME"