#!/bin/bash
# Filename:    init_race_templates.sh
# Author:      Joseph DeVictoria
# Date:        March_18_2019
# Purpose:     Simple script populate the race templates table in the databse with race information.

echo "Populating Race Templates..."

# Make sure I am in this directory.
cd $(dirname $0)

# Import to sqlite3
# Import csv data into a temp table, then insert the data into race templates in order to auto-generate the ids
sqlite3 Oldentide.db <<EOF
.mode csv
DELETE FROM race_templates;
DROP TABLE IF EXISTS temp_race_templates;
.import race_templates.csv temp_race_templates
INSERT INTO race_templates(
    race,
    strength_mod,
    constitution_mod,
    intelligence_mod,
    dexterity_mod,
    axe_mod,
    dagger_mod,
    unarmed_mod,
    hammer_mod,
    polearm_mod,
    spear_mod,
    staff_mod,
    sword_mod,
    archery_mod,
    crossbow_mod,
    sling_mod,
    thrown_mod,
    armor_mod,
    dualweapon_mod,
    shield_mod,
    bardic_mod,
    conjuring_mod,
    druidic_mod,
    illusion_mod,
    necromancy_mod,
    sorcery_mod,
    shamanic_mod,
    spellcraft_mod,
    summoning_mod,
    focus_mod,
    armorsmithing_mod,
    tailoring_mod,
    fletching_mod,
    weaponsmithing_mod,
    alchemy_mod,
    lapidary_mod,
    calligraphy_mod,
    enchanting_mod,
    herbalism_mod,
    hunting_mod,
    mining_mod,
    bargaining_mod,
    camping_mod,
    firstaid_mod,
    lore_mod,
    picklocks_mod,
    scouting_mod,
    search_mod,
    stealth_mod,
    traps_mod,
    aeolandis_mod,
    hieroform_mod,
    highgundis_mod,
    oldpraxic_mod,
    praxic_mod,
    runic_mod,
    description
) 
SELECT * FROM temp_race_templates;
DROP TABLE temp_race_templates;
.quit
EOF
echo "Race Templates Populated."
