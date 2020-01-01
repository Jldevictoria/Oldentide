#!/bin/bash
# Filename:    restore.sh
# Author:      Joseph DeVictoria
# Date:        March_30_2017
# Purpose:     Simple script to restore database from csv files.

echo "Restoring database now..."

sqlite3 oldentide.db <<EOF
.headers on
.mode csv
.import accounts_backup.csv accounts
.quit
EOF
echo "Accounts Restored."

sqlite3 oldentide.db <<EOF
.headers on
.mode csv
.import players_backup.csv players
.quit
EOF
echo "Players Restored."

sqlite3 oldentide.db <<EOF
.headers on
.mode csv
.import npcs_backup.csv npcs
.quit
EOF
echo "NPCs Restored."

sqlite3 oldentide.db <<EOF
.headers on
.mode csv
.import item_templates_backup.csv item_templates
.quit
EOF
echo "Items templates Restored."

sqlite3 oldentide.db <<EOF
.headers on
.mode csv
.import spell_templates_backup.csv spell_templates
.quit
EOF
echo "Spell templates Restored."

sqlite3 oldentide.db <<EOF
.headers on
.mode csv
.import race_templates_backup.csv race_templates
.quit
EOF
echo "Race templates Restored."

sqlite3 oldentide.db <<EOF
.headers on
.mode csv
.import profession_templates_backup.csv profession_templates
.quit
EOF
echo "Profession templates Restored."