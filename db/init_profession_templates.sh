#!/bin/bash
# Filename:    init_profession_templates.sh
# Author:      Joseph DeVictoria
# Date:        March_18_2019
# Purpose:     Simple script populate the profession templates table in the databse with profession information.

echo "Populating Profession Templates..."

# Make sure I am in this directory.
cd $(dirname $0)

# Import to sqlite3
# Import csv data into a temp table, then insert the data into profession templates in order to auto-generate the ids
sqlite3 Oldentide.db <<EOF
.mode csv
DELETE FROM profession_templates;
DROP TABLE IF EXISTS temp_profession_templates;
.import profession_templates.csv temp_profession_templates
INSERT INTO profession_templates(
    profession,
    hppl,
    mppl,
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
    skill_1_multi,
    skill_1_names,
    skill_1_value,
    skill_2_multi,
    skill_2_names,
    skill_2_value,
    skill_3_multi,
    skill_3_names,
    skill_3_value,
    skill_4_multi,
    skill_4_names,
    skill_4_value,
    skill_5_multi,
    skill_5_names,
    skill_5_value,
    description
) 
SELECT * FROM temp_profession_templates;
DROP TABLE temp_profession_templates;
.quit
EOF
echo "Profession Templates Populated."
