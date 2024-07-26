#!/bin/sh
# wait-for-it.sh

set -e

host="$1"
shift
cmd="$@"

until PGPASSWORD=$DB_PASSWORD psql -h "$host" -U "postgres" -c '\q'; do
  >&2 echo "[wait-for-it.sh] Postgres is unavailable - sleeping"
  sleep 1
done

>&2 echo "[wait-for-it.sh] Postgres is up - running application"
exec $cmd
