#!/bin/sh

apk --no-cache add curl
mkdir /app
mkdir /app/migrate
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.0/migrate.linux-amd64.tar.gz | tar xvz -C /app/migrate
mv /app/migrate/migrate /usr/bin/
rm -r /app/migrate