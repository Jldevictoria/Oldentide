#!/bin/bash
# Filename:    restore.sh
# Author:      Joseph DeVictoria
# Date:        March_30_2017
# Purpose:     Simple script to restore database from csv files.

echo "Restoring database now..."

sqlite3 Oldentide.db <<EOF
.headers on
.mode csv
.import accounts_backup.csv accounts
.quit
EOF
echo "Accounts Restored."

sqlite3 Oldentide.db <<EOF
.headers on
.mode csv
.import players_backup.csv players
.quit
EOF
echo "Players Restored."

sqlite3 Oldentide.db <<EOF
.headers on
.mode csv
.import npcs_backup.csv npcs
.quit
EOF
echo "NPCs Restored."

sqlite3 Oldentide.db <<EOF
.headers on
.mode csv
.import items_backup.csv items
.quit
EOF
echo "Items Restored."
