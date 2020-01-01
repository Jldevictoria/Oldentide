#!/bin/bash
# Filename:    backup.sh
# Author:      Joseph DeVictoria
# Date:        March_30_2017
# Purpose:     Simple script to backup database to csv files.

echo "Backup up database now..."

sqlite3 oldentide.db <<EOF
.headers on
.mode csv
.output accounts_backup.csv
SELECT * FROM accounts;
.quit
EOF
echo "Accounts complete."

sqlite3 oldentide.db <<EOF
.headers on
.mode csv
.output players_backup.csv
SELECT * FROM players;
.quit
EOF
echo "Players complete."

sqlite3 oldentide.db <<EOF
.headers on
.mode csv
.output npcs_backup.csv
SELECT * FROM npcs;
.quit
EOF
echo "NPCs complete."

sqlite3 oldentide.db <<EOF
.headers on
.mode csv
.output item_templates_backup.csv
SELECT * FROM item_templates;
.quit
EOF
echo "Items templates complete."

sqlite3 oldentide.db <<EOF
.headers on
.mode csv
.output spell_templates_backup.csv
SELECT * FROM spell_templates;
.quit
EOF
echo "Spell templates complete."

sqlite3 oldentide.db <<EOF
.headers on
.mode csv
.output race_templates_backup.csv
SELECT * FROM race_templates;
.quit
EOF
echo "Race templates complete."

sqlite3 oldentide.db <<EOF
.headers on
.mode csv
.output profession_templates_backup.csv
SELECT * FROM profession_templates;
.quit
EOF
echo "Profession templates complete."