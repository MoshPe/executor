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

echo "Type in the server name you’ll use to connect to the Docker server"
read -p '>' SERVER

echo "Type in the number of days til the key will expire - max 365"
read -p '>' DAYS

# 256bit AES (Advanced Encryption Standard) is the encryption cipher which is used for generating certificate authority (CA) with 2048-bit security.
openssl genrsa -aes256 -passout pass:$PASSWORD -out ca-key.pem 4096 

# Sign the the previously created CA key with your password and address for a period of one year.
# i.e. generating a self-signed certificate for CA
# X.509 is a standard that defines the format of public key certificates, with fixed size 256-bit (32-byte) hash

openssl req -new -x509 -days $DAYS -key ca-key.pem -passin pass:$PASSWORD -sha256 -out ca.pem -subj "/C=TR/ST=./L=./O=./CN=$SERVER" 

# Generating a server key with 2048-bit security
openssl genrsa -out server-key.pem 4096

# Generating a certificate signing request (CSR) for the the server key with the name of your host.
openssl req -new -key server-key.pem -subj "/CN=$SERVER"  -out server.csr

# Sign the key with your password for a period of one year
# i.e. generating a self-signed certificate for the key

echo "Type in the IP for the TLS to connect, Example: 1.2.3.4"
read '>' IP1

echo "Type in the second IP if needed, otherwise type '.':"
read '>' IP2

if [[ IP2 = '.' ]]; then
    sh -c 'echo "subjectAltName = DNS:$SERVER,IP:$IP1" >> extfile.cnf'
else
    sh -c 'echo "subjectAltName = DNS:$SERVER,IP:$IP1,IP:$IP2" >> extfile.cnf'
fi

sh -c 'echo "extendedKeyUsage = serverAuth" >> extfile.cnf'

openssl x509 -req -days $DAYS -sha256 -in server.csr -CA ca.pem -CAkey ca-key.pem -passin "pass:$PASSWORD" -CAcreateserial -out server-cert.pem -extfile extfile.cnf
