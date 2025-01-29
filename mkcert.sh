#!/bin/sh

CERT_DIR="./cert"
CERT_FILE="$CERT_DIR/cert.pem"
KEY_FILE="$CERT_DIR/cert.key"

# Check if mkcert is installed
if ! command -v mkcert &> /dev/null; then
    echo "mkcert is not installed. Please install mkcert to generate certificates: https://github.com/FiloSottile/mkcert"
    exit 1
fi

# Function to check if file exists
file_not_exists() {
    [ ! -f "$1" ]
}

# Generate certs if they don't exist
if file_not_exists "$CERT_FILE" || file_not_exists "$KEY_FILE"; then
    echo "Certificates not found. Generating new certificates..."
    mkcert -install
    mkdir -p "$CERT_DIR"
    mkcert -cert-file "$CERT_FILE" -key-file "$KEY_FILE" localhost
else
    echo "Certificates found. Skipping generation."
fi