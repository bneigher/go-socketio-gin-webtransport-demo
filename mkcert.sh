#!/bin/sh

CERT_DIR="./cert"
CERT_FILE="$CERT_DIR/cert.pem"
KEY_FILE="$CERT_DIR/cert.key"
EXT_FILE="$CERT_DIR/san.cnf"

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
    mkdir -p "$CERT_DIR"

    echo "Generating CA key and certificate:"
    openssl req -x509 -sha256 -nodes -days 3650 -newkey rsa:2048 \
      -keyout $CERT_DIR/ca.key -out $CERT_DIR/ca.pem \
      -subj "/O=quic-go Certificate Authority/"

    echo "Generating CSR"
    openssl req -out $CERT_DIR/cert.csr -new -newkey rsa:2048 -nodes -keyout $CERT_DIR/cert.key \
      -subj "/O=quic-go/"

    # Create a temporary SAN configuration file
    cat > "$EXT_FILE" <<EOF
subjectAltName=DNS:localhost
EOF

    echo "Sign certificate with Subject Alternative Names (SAN):"
    openssl x509 -req -sha256 -days 3650 -in $CERT_DIR/cert.csr -out $CERT_DIR/cert.pem \
      -CA $CERT_DIR/ca.pem -CAkey $CERT_DIR/ca.key -CAcreateserial -extfile "$EXT_FILE"

    # debug output the certificate
    openssl x509 -noout -text -in $CERT_DIR/cert.pem

    # Clean up unnecessary files
    rm $CERT_DIR/ca.key $CERT_DIR/cert.csr $EXT_FILE
else
    echo "Certificates found. Skipping generation."
fi