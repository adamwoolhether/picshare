#!/bin/bash

# Change to the directory with our code that we plan to work from
# NOTE: DO NOT INCLUDE THIS LINE IF YOU ARE USING GO MODULES!
cd "$GOPATH/src/picshare"

echo "==== Releasing picshare ===="
echo "  Deleting the local binary if it exists (so it isn't uploaded)..."
rm picshare
echo "  Done!"

echo "  Deleting existing code..."
ssh root@www.adamwoolhether.com "rm -rf /root/go/src/picshare && rm -rf /root/goo/src/picapp"
echo "  Existing code deleted successfully!"

echo "  Uploading local code..."

ssh root@www.adamwoolhether.com "mkdir -p /root/go/src/picshare/ && mkdir -p /root/app/"
rsync -avr --exclude-from='./scripts/exclude.txt' ./ \
  root@www.adamwoolhether.com:/root/go/src/picshare/
echo "  Code uploaded successfully!"

echo "  Creating picshare service"
ssh root@www.adamwoolhether.com "sudo -E sh -c 'cat > /etc/systemd/system/picshare.com.service <<EOF
[Unit]
Description=picshare.com app

[Service]
WorkingDirectory=/root/app
ExecStart=/root/app/server -prod
Restart=always
RestartSec=30

[Install]
WantedBy=multi-user.target
EOF'"

echo "  Building the code on remote server..."
ssh root@www.adamwoolhether.com 'export GOPATH=/root/go; \
  cd /root/go/src/picshare; \
  /usr/local/go/bin/go mod download; \
  /usr/local/go/bin/go build -o /root/app/server \
    $GOPATH/src/picshare/*.go'
echo "  Code built successfully!"

echo "  Moving assets..."
ssh root@www.adamwoolhether.com "cd /root/app; \
  cp -R /root/go/src/picshare/assets ."
echo "  Assets moved successfully!"

echo "  Moving views..."
ssh root@www.adamwoolhether.com "cd /root/app; \
  cp -R /root/go/src/picshare/views ."
echo "  Views moved successfully!"

echo "  Moving Caddyfile..."
ssh root@www.adamwoolhether.com "cp /root/go/src/picshare/Caddyfile /etc/caddy/Caddyfile"
echo "  Caddyfile moved successfully!"

echo "  Moving .config file..."
ssh root@www.adamwoolhether.com "cp /root/go/src/picshare/prodconfig.config /root/app/.config"
echo "  Caddyfile .config successfully!"

echo "  Restarting Caddy server..."
ssh root@www.adamwoolhether.com "sudo systemctl restart caddy"
echo "  Caddy restarted successfully!"

echo "  Restarting the server..."
ssh root@www.adamwoolhether.com "sudo systemctl reboot"
echo "  Server restarted successfully!"

echo "==== Done releasing picshare ===="