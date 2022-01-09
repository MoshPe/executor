#!bin/bash
set -eu

#set -x ; debugging

cd ~
echo "you are now in $PWD"

if [ ! -d ".docker/" ] 
then
    echo "Directory ./docker/ does not exist"
    echo "Creating the directory"
    mkdir .docker
fi

cd .docker/
echo "type in your certificate password (characters are not echoed)"
read -p '>' -s PASSWORD 

echo "Type in the number of days til the key will expire - max 365"
read -p '>' DAYS

# For client authentication, create a client key and certificate signing request
# Generate a client key with 2048-bit security
openssl genrsa -out key.pem 4096
# Process the key as a client key.
openssl req -subj '/CN=client' -new -key key.pem -out client.csr

# To make the key suitable for client authentication, create an extensions config file:
sh -c 'echo "extendedKeyUsage = clientAuth" > extfile-client.cnf'

# Sign the (public) key with your password for a period of one year
openssl x509 -req -days $DAYS -sha256 -in client.csr -CA ca.pem -CAkey ca-key.pem -passin "pass:$PASSWORD" -CAcreateserial -out cert.pem -extfile extfile-client.cnf

echo "Removing unnecessary files i.e. client.csr extfile.cnf server.csr"
rm ca.srl client.csr extfile.cnf server.csr

echo "Changing the permissions to readonly by root for the server files."
# To make them only readable by you: 
chmod 0400 ca-key.pem key.pem server-key.pem

echo "Changing the permissions of the client files to read-only by everyone"
# Certificates can be world-readable, but you might want to remove write access to prevent accidental damage
# these are all x509 certificates aka public key certificates
# X.509 certificates are used in many Internet protocols, including TLS/SSL, which is the basis for HTTPS.
chmod 0444 ca.pem server-cert.pem cert.pem
