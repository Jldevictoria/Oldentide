#!/bin/bash
# Filename: initialize_databse.sh
# Author:   Joseph DeVictoria
# Date:     December_31_2019
# Purpose:  Fully initialize the oldentide database based on the raw csv data provided.
# Usage:    ./initialize_database.sh <databse_name.db>

if [ $# -eq 0 ]; then
    echo "Please supply a database name i.e: \"./initialize_database.sh database_name.db\""
fi

##  -----------------------------------------  ##
##   _____ _____  _   _  ________  ___  ___    ##
##  /  ___/  __ \| | | ||  ___|  \/  | / _ \   ##
##  \ `--.| /  \/| |_| || |__ | .  . |/ /_\ \  ##
##   `--. \ |    |  _  ||  __|| |\/| ||  _  |  ##
##  /\__/ / \__/\| | | || |___| |  | || | | |  ##
##  \____/ \____/\_| |_/\____/\_|  |_/\_| |_/  ##
##                                             ##
## ------------------------------------------  ##

sqlite3 $1 < database_schema.sql

##  ----------------------------  ##
##   _   _ ______  _____  _____   ##
##  | \ | || ___ \/  __ \/  ___|  ##
##  |  \| || |_/ /| /  \/\ `--.   ##
##  | . ` ||  __/ | |     `--. \  ##
##  | |\  || |    | \__/\/\__/ /  ##
##  \_| \_/\_|     \____/\____/   ##
##                                ##
##  ----------------------------  ##

echo "Populating NPCs..."

# Make sure I am in this directory.
cd $(dirname $0)

# Import to sqlite3
# Import csv data into a temp table, then insert the data into npcs in order to auto-generate the ids
sqlite3 oldentide.db <<EOF
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
    sp,
    maxsp,
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

##  --------------------------------  ##
##  _____ _____ ________  ___ _____   ##
## |_   _|_   _|  ___|  \/  |/  ___|  ##
##   | |   | | | |__ | .  . |\ `--.   ##
##   | |   | | |  __|| |\/| | `--. \  ##
##  _| |_  | | | |___| |  | |/\__/ /  ##
##  \___/  \_/ \____/\_|  |_/\____/   ##
##                                    ##
##  --------------------------------  ##

echo "Populating Item Templates..."

# Make sure I am in this directory.
cd $(dirname $0)

# Import to sqlite3
# Import csv data into a temp table, then insert the data into npcs in order to auto-generate the ids
sqlite3 oldentide.db <<EOF
.mode csv
DELETE FROM item_templates;
DROP TABLE IF EXISTS temp_item_templates;
.import item_templates.csv temp_item_templates
INSERT INTO item_templates(
    name,
    true_name,
    lore_level,
    type,
    slot,
    icon,
    weight,
    encumbrance,
    dyeable,
    stackable,
    stack_size,
    usable,
    equipable,
    base_price,
    strength_requirement,
    constitution_requirement,
    intelligence_requirement,
    dexterity_requirement,
    skill_type_0,
    skill_requirement_0,
    skill_type_1,
    skill_requirement_1,
    skill_type_2,
    skill_requirement_2,
    skill_type_3,
    skill_requirement_3,
    skill_type_4,
    skill_requirement_4,
    description,
    true_description
) 
SELECT * FROM temp_item_templates;
DROP TABLE temp_item_templates;
.quit
EOF
echo "Item Templates Populated."

##  ---------------------------------------  ##
##   ___________ _____ _      _      _____   ##
##  /  ___| ___ \  ___| |    | |    /  ___|  ##
##  \ `--.| |_/ / |__ | |    | |    \ `--.   ##
##   `--. \  __/|  __|| |    | |     `--. \  ##
##  /\__/ / |   | |___| |____| |____/\__/ /  ##
##  \____/\_|   \____/\_____/\_____/\____/   ##
##                                           ##
##  ---------------------------------------  ##

echo "Populating Spell Templates..."

# Make sure I am in this directory.
cd $(dirname $0)

# Import to sqlite3
# Import csv data into a temp table, then insert the data into npcs in order to auto-generate the ids
sqlite3 oldentide.db <<EOF
.mode csv
DELETE FROM spell_templates;
DROP TABLE IF EXISTS temp_spell_templates;
.import spell_templates.csv temp_spell_templates
INSERT INTO spell_templates(
    spellname,
    school,
    level,
    type,
    target,
    range,
    accuracy,
    preparation_time,
    recovery_time,
    effect_1,
    effect_2,
    effect_3,
    effect_4,
    effect_5,
    description
) 
SELECT * FROM temp_spell_templates;
DROP TABLE temp_spell_templates;
.quit
EOF
echo "Spell Templates Populated."
##  --------------------------------  ##
##  ______  ___  _____  _____ _____   ##
##  | ___ \/ _ \/  __ \|  ___/  ___|  ##
##  | |_/ / /_\ \ /  \/| |__ \ `--.   ##
##  |    /|  _  | |    |  __| `--. \  ##
##  | |\ \| | | | \__/\| |___/\__/ /  ##
##  \_| \_\_| |_/\____/\____/\____/   ##
##                                    ##
##  --------------------------------  ##

echo "Populating Race Templates..."

# Make sure I am in this directory.
cd $(dirname $0)

# Import to sqlite3
# Import csv data into a temp table, then insert the data into race templates in order to auto-generate the ids
sqlite3 oldentide.db <<EOF
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

##  --------------------------------------------------------------------  ##
##  ____________ ___________ _____ _____ _____ _____ _____ _   _  _____   ##
##  | ___ \ ___ \  _  |  ___|  ___/  ___/  ___|_   _|  _  | \ | |/  ___|  ##
##  | |_/ / |_/ / | | | |_  | |__ \ `--.\ `--.  | | | | | |  \| |\ `--.   ##
##  |  __/|    /| | | |  _| |  __| `--. \`--. \ | | | | | | . ` | `--. \  ##
##  | |   | |\ \\ \_/ / |   | |___/\__/ /\__/ /_| |_\ \_/ / |\  |/\__/ /  ##
##  \_|   \_| \_|\___/\_|   \____/\____/\____/ \___/ \___/\_| \_/\____/   ##
##                                                                        ##
##  --------------------------------------------------------------------  ##

echo "Populating Profession Templates..."

# Make sure I am in this directory.
cd $(dirname $0)

# Import to sqlite3
# Import csv data into a temp table, then insert the data into profession templates in order to auto-generate the ids
sqlite3 oldentide.db <<EOF
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

echo "Database Fully Initializd! Exiting..."
