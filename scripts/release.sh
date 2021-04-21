#!/bin/bash

# Change to the directory with our code that we plan to work from
# NOTE: DO NOT INCLUDE THIS LINE IF YOU ARE USING GO MODULES!
cd "$GOPATH/src/piacapp"

echo "==== Releasing picapp ===="
echo "  Deleting the local binary if it exists (so it isn't uploaded)..."
rm picapp
echo "  Done!"

echo "  Deleting existing code..."
ssh root@test.adamwoolhether.com "rm -rf /root/go/src/picapp"
echo "  Code deleted successfully!"

echo "  Uploading code..."
# The \ at the end of the line tells bash that our
# command isn't done and wraps to the next line.
ssh root@test.adamwoolhether.com "mkdir -p /root/go/src/picapp/"
rsync -avr --exclude-from='./scripts/exclude.txt' ./ \
  root@test.adamwoolhether.com:/root/go/src/picapp/
echo "  Code uploaded successfully!"

#echo "  Go getting deps..."
#ssh root@test.adamwoolhether.com "export GOPATH=/root/go; \
#  /usr/local/go/bin/go get golang.org/x/crypto/bcrypt"
#ssh root@test.adamwoolhether.com "export GOPATH=/root/go; \
#  /usr/local/go/bin/go get github.com/gorilla/mux"
#ssh root@test.adamwoolhether.com "export GOPATH=/root/go; \
#  /usr/local/go/bin/go get github.com/gorilla/schema"
#ssh root@test.adamwoolhether.com "export GOPATH=/root/go; \
#  /usr/local/go/bin/go get github.com/lib/pq"
#ssh root@test.adamwoolhether.com "export GOPATH=/root/go; \
#  /usr/local/go/bin/go get github.com/jinzhu/gorm"
#ssh root@test.adamwoolhether.com "export GOPATH=/root/go; \
#  /usr/local/go/bin/go get github.com/gorilla/csrf"
#ssh root@test.adamwoolhether.com "export GOPATH=/root/go; \
#  /usr/local/go/bin/go get gopkg.in/mailgun/mailgun-go.v1"

echo "  Building the code on remote server..."
ssh root@test.adamwoolhether.com 'export GOPATH=/root/go; \
  cd /root/go/src/picapp; \
  /usr/local/go/bin/go mod download; \
  /usr/local/go/bin/go build -o /root/app/server \
    $GOPATH/src/picapp/*.go'
echo "  Code built successfully!"

echo "  Moving assets..."
ssh root@test.adamwoolhether.com "cd /root/app; \
  cp -R /root/go/src/picapp/assets ."
echo "  Assets moved successfully!"

echo "  Moving views..."
ssh root@test.adamwoolhether.com "cd /root/app; \
  cp -R /root/go/src/picapp/views ."
echo "  Views moved successfully!"

echo "  Moving Caddyfile..."
ssh root@test.adamwoolhether.com "cp /root/go/src/picapp/Caddyfile /etc/caddy/Caddyfile"
echo "  Caddyfile moved successfully!"

echo "  Copying the binary..."
ssh root@test.adamwoolhether.com "cp server /root/app"
echo "  Binary copied successfully!"

echo "  Restarting the server..."
ssh root@test.adamwoolhether.com "sudo systemctl restart picshare.com"
echo "  Server restarted successfully!"

echo "  Restarting Caddy server..."
ssh root@test.adamwoolhether.com "sudo systemctl restart caddy"
echo "  Caddy restarted successfully!"

echo "==== Done releasing picapp ===="