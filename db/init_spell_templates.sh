#!/bin/bash
# Filename:    init_spell_templates.sh
# Author:      Joseph DeVictoria
# Date:        December_31_2019
# Purpose:     Simple script populate the spell_templates table in the databse with spell_templates from a csv.

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
