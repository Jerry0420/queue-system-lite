#!/bin/sh

go install github.com/codegangsta/gin@latest
gin --port 3001 --appPort 8000 run main.go