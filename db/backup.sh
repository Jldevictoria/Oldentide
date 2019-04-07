#!/bin/bash
# Filename:    backup.sh
# Author:      Joseph DeVictoria
# Date:        March_30_2017
# Purpose:     Simple script to backup database to csv files.

echo "Backup up database now..."

sqlite3 Oldentide.db <<EOF
.headers on
.mode csv
.output accounts_backup.csv
SELECT * FROM accounts;
.quit
EOF
echo "Accounts complete."

sqlite3 Oldentide.db <<EOF
.headers on
.mode csv
.output players_backup.csv
SELECT * FROM players;
.quit
EOF
echo "Players complete."

sqlite3 Oldentide.db <<EOF
.headers on
.mode csv
.output npcs_backup.csv
SELECT * FROM npcs;
.quit
EOF
echo "NPCs complete."

sqlite3 Oldentide.db <<EOF
.headers on
.mode csv
.output items_backup.csv
SELECT * FROM items;
.quit
EOF
echo "Items complete."

