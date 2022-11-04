#!/bin/bash
DIR="$(dirname "${BASH_SOURCE[0]}")"

DB_NAME=realmmgr-dev-db

if ! psql -h ccs-pg -U postgres -lqt | cut -d \| -f 1 | grep -qw "$DB_NAME"; then
  createdb -e -h ccs-pg -U postgres "$DB_NAME" 'dev database for Service'
fi

# drop tables and types
psql -h ccs-pg -U postgres "$DB_NAME" -f "$DIR"/drop.sql

# install required extensions
psql -h ccs-pg -U postgres "$DB_NAME" -c 'CREATE EXTENSION IF NOT EXISTS "uuid-ossp";'

# create type, tables and insert dummy data
# drop tables and types
psql -h ccs-pg -U postgres "$DB_NAME" -f "$DIR"/create.sql
psql -h ccs-pg -U postgres "$DB_NAME" -f "$DIR"/populate_dummy_data.sql