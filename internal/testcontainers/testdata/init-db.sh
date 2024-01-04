#!/bin/bash
set -e

# Creates a custom table named testdb for the TestContainers PoC.
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE TABLE IF NOT EXISTS testdb (name varchar(255));
EOSQL
