#!/bin/bash

# temp dir: https://stackoverflow.com/a/53063602/328631
# Create a temporary directory and store its name in a variable ...
TMPDIR=$(mktemp -d)

# Bail out if the temp directory wasn't created successfully.
if [ ! -e $TMPDIR ]; then
    >&2 echo "Failed to create temp directory"
    exit 1
fi

# Make sure it gets removed even if the script exits abnormally.
trap "exit 1"           HUP INT PIPE QUIT TERM
trap 'rm -rf "$TMPDIR"' EXIT

# note \$ is for bash.
# source https://www.switch.ch/pki/manage/request/csr-openssl/
# http://apetec.com/support/GenerateSAN-CSR.htm
# https://stackoverflow.com/questions/21488845/how-can-i-generate-a-self-signed-certificate-with-subjectaltname-using-openssl
# https://stackoverflow.com/questions/6194236/openssl-certificate-version-3-with-subject-alternative-name
cat > $TMPDIR/openssl.cnf <<EOF
# OpenSSL configuration file for creating a CSR for a server certificate
# Adapt at least the FQDN and ORGNAME lines, and then run 
# openssl req -new -config myserver.cnf -keyout myserver.key -out myserver.csr
# on the command line.

# the fully qualified server (or service) name
FQDN = example.com

# the name of your organization
# (see also https://www.switch.ch/pki/participants/)
ORGNAME = Example

# subjectAltName entries: to add DNS aliases to the CSR, delete
# the '#' character in the ALTNAMES line, and change the subsequent
# 'DNS:' entries accordingly. Please note: all DNS names must
# resolve to the same IP address as the FQDN.
ALTNAMES = DNS:\$FQDN



# --- no modifications required below ---
[ req ]
default_bits = 2048
default_md = sha256
prompt = no
encrypt_key = no
distinguished_name = dn
req_extensions = req_ext

[ dn ]
C = US
O = \$ORGNAME
CN = \$FQDN

[ req_ext ]
subjectAltName = \$ALTNAMES

# subjectAltName = @alt_names

[alt_names]
DNS = \$FQDN
EOF

cat $TMPDIR/openssl.cnf
openssl req -nodes -newkey rsa:2048 -keyout example_com_key.pem -out $TMPDIR/cert.csr -config $TMPDIR/openssl.cnf -reqexts req_ext
# cat $TMPDIR/cert.csr
# openssl req -in $TMPDIR/cert.csr -noout -text
echo
echo Generating cert
openssl x509 -in $TMPDIR/cert.csr -out example_com_cert.pem -req -signkey example_com_key.pem -days 3650 -extfile $TMPDIR/openssl.cnf -extensions req_ext


echo certs created:
cp example_com_key.pem key.pem
cp example_com_cert.pem cert.pem

openssl x509 -in cert.pem -text -noout
