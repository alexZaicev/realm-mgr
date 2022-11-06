#!/bin/bash

DIR="$(dirname "${BASH_SOURCE[0]}")"

if [ -z ${DATABASE_HOST+x} ]; then
  DATABASE_HOST="ccs-pg"
fi

if [ -z ${DATABASE_PORT+x} ]; then
  DATABASE_PORT="5432"
fi

if [ -z ${DATABASE_USER+x} ]; then
  DATABASE_USER="postgres"
fi

if [ -z ${DATABASE_NAME+x} ]; then
  DATABASE_NAME="realmmgr-dev-db"
fi

if ! psql -h ${DATABASE_HOST} -p ${DATABASE_PORT} -U ${DATABASE_USER} -lqt | cut -d \| -f 1 | grep -qw "${DATABASE_NAME}"; then
  createdb -e -h ${DATABASE_HOST} -p ${DATABASE_PORT} -U ${DATABASE_USER} "${DATABASE_NAME}" 'dev database for Service'
fi

# drop tables and types
psql -h ${DATABASE_HOST} -p ${DATABASE_PORT} -U ${DATABASE_USER} "${DATABASE_NAME}" -f "$DIR"/drop.sql

# install required extensions
psql -h ${DATABASE_HOST} -p ${DATABASE_PORT} -U ${DATABASE_USER} "${DATABASE_NAME}" -c 'CREATE EXTENSION IF NOT EXISTS "uuid-ossp";'

# create type, tables and insert dummy data
# drop tables and types
psql -h ${DATABASE_HOST} -p ${DATABASE_PORT} -U ${DATABASE_USER} "${DATABASE_NAME}" -f "$DIR"/create.sql