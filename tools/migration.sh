#!/bin/bash

# Simple script to migrate the old python database to the new one for Go.

OLDDB=""
NEWDB=""
OPTIND=1
while getopts "o:n:h?" opt; do
    case "$opt" in
    o)
        OLDDB="$OPTARG"
        ;;
    n)
        NEWDB="$OPTARG"
        ;;
    h|?)
        echo "Usage options:"
        echo -e "-h, -?\t\tShow this help and exit"
        echo -e "-o <file>\tOld database file to dump data from (required)"
        echo -e "-n <file>\tNew database file to insert data into (required)"
        exit 0
        ;;
    esac
done

if [[ ! -f "$OLDDB" || ! -f "$NEWDB" ]]; then
    echo "Missing database files!"
    exit 1
fi


TEMPFILE=$(tempfile)
read -d "" TMP << EOF
PRAGMA foreign_keys = ON;

DROP TABLE IF EXISTS "server";
DROP TABLE IF EXISTS "server_history";

CREATE TABLE "servers" ("id" integer primary key autoincrement,"last_updated" datetime,"title" varchar(64)  UNIQUE,"game_url" varchar(255),"site_url" varchar(255),"players_current" integer,"players_avg" integer,"players_min" integer,"players_max" integer,"players_mon" integer,"players_tue" integer,"players_wed" integer,"players_thu" integer,"players_fri" integer,"players_sat" integer,"players_sun" integer );
CREATE TABLE "server_populations" ("id" integer primary key autoincrement,"timestamp" datetime,"players" integer,"server_id" integer REFERENCES servers(id) ON DELETE CASCADE ON UPDATE CASCADE );
CREATE INDEX idx_server_populations_server_id ON "server_populations"("server_id") ;

BEGIN TRANSACTION;
EOF
echo "$TMP" > "$TEMPFILE"

read -d "" TMP << EOF
.mode insert
select
    id,
    last_updated,
    title,
    game_url,
    site_url,
    players_current,
    players_avg,
    players_min,
    players_max,
    0,0,0,0,0,0,0
from gameservers_server;
EOF
echo "$TMP" | sqlite3 "$OLDDB" | sed 's/table/servers/g' >> "$TEMPFILE"

read -d "" TMP << EOF
.mode insert
select
    id,
    created,
    players,
    server_id
from gameservers_serverhistory;
EOF
echo "$TMP" | sqlite3 "$OLDDB" | sed 's/table/server_populations/g' >> "$TEMPFILE"
echo "COMMIT;" >> "$TEMPFILE"


sqlite3 "$NEWDB" < "$TEMPFILE"
rm "$TEMPFILE"
