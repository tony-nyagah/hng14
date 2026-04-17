#!/usr/bin/env bash
set -e

APP_DIR=/home/ubuntu/simple-api

mkdir -p $APP_DIR
chmod +x $APP_DIR/simple-api

sudo cp $APP_DIR/simple-api.service /etc/systemd/system/simple-api.service
sudo systemctl daemon-reload
sudo systemctl enable simple-api
sudo systemctl restart simple-api
sudo systemctl status simple-api --no-pager

echo systemd done

sudo cp $APP_DIR/simple-api.nginx.conf /etc/nginx/sites-available/simple-api
sudo ln -sf /etc/nginx/sites-available/simple-api /etc/nginx/sites-enabled/simple-api
sudo rm -f /etc/nginx/sites-enabled/default

sudo nginx -t
sudo systemctl reload nginx
echo nginx done
