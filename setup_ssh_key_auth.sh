#!/bin/bash
set -e

if [ "$#" -ne 3 ]; then
    echo "Usage: $0 <username> <server_ip> <email>"
    exit 1
fi

USERNAME="$1"
SERVER_IP="$2"
EMAIL="$3"

# Generate an SSH key pair if it doesn't exist
if [ ! -f "$HOME/.ssh/id_rsa" ]; then
    ssh-keygen -t rsa -b 4096 -C "$EMAIL"
fi

# Copy the public key to the remote server
ssh-copy-id "$USERNAME@$SERVER_IP"

# Test SSH key-based authentication
echo "Testing SSH key-based authentication:"
ssh "$USERNAME@$SERVER_IP" "echo 'SSH key-based authentication is successfully set up!'"
