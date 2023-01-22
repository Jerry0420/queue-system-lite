#!/bin/sh

dburl=postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@$POSTGRES_HOST:$POSTGRES_PORT/$POSTGRES_DB?sslmode=disable
CURRENTDIR=$(dirname "$0")

option="${1}"
case ${option} in
    up) echo up db
      migrate -path $CURRENTDIR/migrations -database $dburl up
      ;;
    down) echo down db
      migrate -path $CURRENTDIR/migrations -database $dburl down
      ;;
    create) echo create db $2
      migrate create -ext sql -dir $CURRENTDIR/migrations -seq $2
      ;;
    *)  echo 'Unknown!'
      exit 0
    ;;
esac

# migrate -path /__w/queue-system/queue-system/scripts/migration_tools/migrations -database postgres://root:root@db_test:5432/queue_system?sslmode=disable down