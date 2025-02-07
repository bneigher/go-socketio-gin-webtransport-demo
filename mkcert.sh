#!/bin/sh

CERT_DIR="./cert"
CERT_FILE="$CERT_DIR/cert.pem"
KEY_FILE="$CERT_DIR/cert.key"

# Check if openssl is installed
if ! command -v openssl &> /dev/null; then
    echo "openssl is not installed. Please install openssl to generate certificates"
    exit 1
fi

# Function to check if file exists
file_not_exists() {
    [ ! -f "$1" ]
}

# Generate certs if they don't exist
if file_not_exists "$CERT_FILE" || file_not_exists "$KEY_FILE"; then
    echo "Certificates not found. Generating new certificates..."
    openssl req -x509 -sha256 -nodes -days 3650 -newkey rsa:2048 \
      -keyout $CERT_DIR/ca.key -out $CERT_DIR/ca.pem \
      -subj "/O=quic-go Certificate Authority/"

    echo "Generating CSR"
    openssl req -out $CERT_DIR/cert.csr -new -newkey rsa:2048 -nodes -keyout $CERT_DIR/cert.key \
      -subj "/O=quic-go/"

    echo "Sign certificate:"
    openssl x509 -req -sha256 -days 3650 -in $CERT_DIR/cert.csr  -out $CERT_DIR/cert.pem \
      -CA $CERT_DIR/ca.pem -CAkey $CERT_DIR/ca.key -CAcreateserial
      # -extfile <(printf "subjectAltName=DNS:localhost")

    # debug output the certificate
    openssl x509 -noout -text -in $CERT_DIR/cert.pem

    # we don't need the CA key, the serial number and the CSR any more
    rm $CERT_DIR/ca.key $CERT_DIR/cert.csr .srl
else
    echo "Certificates found. Skipping generation."
fi