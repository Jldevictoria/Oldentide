#!/bin/bash
# Filename:    init_npcs.sh
# Author:      Joseph DeVictoria
# Date:        September_2_2017
# Purpose:     Simple script populate the NPC table in the databse with NPC information and nodes.

echo "Populating NPCs..."

# Make sure I am in this directory.
cd $(dirname $0)

# Import to sqlite3
# Import csv data into a temp table, then insert the data into npcs in order to auto-generate the ids
sqlite3 Oldentide.db <<EOF
.mode csv
DELETE FROM npcs;
DROP TABLE IF EXISTS temp_npcs;
.import npcs.csv temp_npcs
INSERT INTO npcs(
    firstname,
    lastname,
    guild,
    race,
    gender,
    face,
    skin,
    profession,
    alive,
    level,
    hp,
    maxhp,
    bp,
    maxbp,
    mp,
    maxmp,
    ep,
    maxep,
    strength,
    constitution,
    intelligence,
    dexterity,
    head,
    chest,
    arms,
    hands,
    legs,
    feet,
    cloak,
    righthand,
    lefthand,
    zone,
    x,
    y,
    z,
    direction
) 
SELECT * FROM temp_npcs;
DROP TABLE temp_npcs;
.quit
EOF
echo "NPCs Populated."
