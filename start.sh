#!/bin/sh


set -e

echo "run db migration"

source /app/app.env

echo "HIHIHI  $DB_SOURCE"

/app/migrate -path /app/migration --database "$DB_SOURCE" -verbose up

echo "start app"

echo "HIHIHI  $@"

exec "$@"


