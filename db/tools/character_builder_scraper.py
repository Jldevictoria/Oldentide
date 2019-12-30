#!/bin/python
# This script will parse and reformat all of the race and profession modifiers in the RoE files.

# Modules.
import argparse
import os
import pprint
import re
import sys
import time

# Set up command line arguments.
parser = argparse.ArgumentParser(description='This script will parse and reformat all of the profession and class modifiers in the RoE files.')
parser.add_argument('crsm_dir', type=str, help='Full path of the Crsm folder in the RoE game directory.')
args = parser.parse_args()

# Initialize path and files.
race_modifiers_file = 'race_templates.csv'
profession_modifiers_file = 'profession_templates.csv'

# Initialize modifiers data structures.
skills = (
	"Hppl",
	"Mppl",
	"Strength",
	"Constitution",
	"Intelligence",
	"Dexterity",
	"Axe",
	"Dagger",
	"Hammer",
	"Polearm",
	"Spear",
	"Staff",
	"Sword",
	"Unarmed",
	"Archery",
	"Crossbow",
	"Sling",
	"Thrown",
	"Armor",
	"Dual Weapon",
	"Shield",
	"Bardic",
	"Conjuring",
	"Druidic",
	"Illusion",
	"Necromancy",
	"Shamanic",
	"Sorcery",
	"Summoning",
	"Spellcraft",
	"Focus",
	"Alchemy",
	"Armorsmithing",
	"Calligraphy",
	"Enchanting",
	"Fletching",
	"Lapidary",
	"Tailoring",
	"Weaponsmithing",
	"Herbalism",
	"Hunting",
	"Mining",
	"Bargaining",
	"Camping",
	"First Aid",
	"Lore",
	"Pick Locks",
	"Scouting",
	"Search",
	"Stealth",
	"Traps",
	"Aeolandis",
	"Hieroform",
	"High Gundis",
	"Old Praxic",
	"Praxic",
	"Runic",
	"Skill_1_multi",
	"Skill_1_names",
	"Skill_1_value",
	"Skill_2_multi",
	"Skill_2_names",
	"Skill_2_value",
	"Skill_3_multi",
	"Skill_3_names",
	"Skill_3_value",
	"Skill_4_multi",
	"Skill_4_names",
	"Skill_4_value",
	"Skill_5_multi",
	"Skill_5_names",
	"Skill_5_value",
	"Description"
)

headers = (
	"hppl",
	"mppl",
	"strength_mod",
	"constitution_mod",
	"intelligence_mod",
	"dexterity_mod",
	"axe_mod",
    "dagger_mod",
    "hammer_mod",
    "polearm_mod",
    "spear_mod",
    "staff_mod",
    "sword_mod",
    "unarmed_mod",
    "archery_mod",
    "crossbow_mod",
    "sling_mod",
    "thrown_mod",
    "armor_mod",
    "dualweapon_mod",
    "shield_mod",
    "bardic_mod",
    "conjuring_mod",
    "druidic_mod",
    "illusion_mod",
    "necromancy_mod",
    "shamanic_mod",
    "sorcery_mod",
    "summoning_mod",
    "spellcraft_mod",
    "focus_mod",
    "alchemy_mod",
    "armorsmithing_mod",
    "calligraphy_mod",
    "enchanting_mod",
    "fletching_mod",
    "lapidary_mod",
    "tailoring_mod",
    "weaponsmithing_mod",
    "herbalism_mod",
    "hunting_mod",
    "mining_mod",
    "bargaining_mod",
    "camping_mod",
    "firstaid_mod",
    "lore_mod",
    "picklocks_mod",
    "scouting_mod",
    "search_mod",
    "stealth_mod",
    "traps_mod",
    "aeolandis_mod",
    "hieroform_mod",
    "highgundis_mod",
    "oldpraxic_mod",
    "praxic_mod",
    "runic_mod",
	"skill_1_multi",
	"skill_1_names",
	"skill_1_value",
	"skill_2_multi",
	"skill_2_names",
	"skill_2_value",
	"skill_3_multi",
	"skill_3_names",
	"skill_3_value",
	"skill_4_multi",
	"skill_4_names",
	"skill_4_value",
	"skill_5_multi",
	"skill_5_names",
	"skill_5_value",
	"description"
)

races = dict()
professions = dict()

# Parse race modifiers and store in data structure.
races_path = args.crsm_dir + "\\rsm"
print("Parsing character race modifiers in \"" + races_path + "\"")
for filename in os.listdir(races_path):
	race_name = filename.split('.')[0]
	race_dict = dict()
	race_file_path = races_path + "\\" + filename
	race_data = open(race_file_path).readlines()
	for line in race_data:
		skill, mod = line.rstrip("\n").rsplit(None,1)
		if skill in skills:
			race_dict[skill] = float(mod)
	if race_name == "Dwarf":
		race_dict["Strength"] = 5.0
		race_dict["Constitution"] = 15.0
		race_dict["Intelligence"] = -5.0
		race_dict["Dexterity"] = -15.0
	elif race_name == "Elf":
		race_dict["Strength"] = -5.0
		race_dict["Constitution"] = -15.0
		race_dict["Intelligence"] = 5.0
		race_dict["Dexterity"] = 15.0
	elif race_name == "Gnome":
		race_dict["Strength"] = -10.0
		race_dict["Constitution"] = -10.0
		race_dict["Intelligence"] = 10.0
		race_dict["Dexterity"] = 10.0
	elif race_name == "Human":
		race_dict["Strength"] = 0.0
		race_dict["Constitution"] = 0.0
		race_dict["Intelligence"] = 0.0
		race_dict["Dexterity"] = 0.0
	elif race_name == "Leshy":
		race_dict["Strength"] = -15.0
		race_dict["Constitution"] = -5.0
		race_dict["Intelligence"] = 20.0
		race_dict["Dexterity"] = 0.0
	elif race_name == "Ogre":
		race_dict["Strength"] = 20.0
		race_dict["Constitution"] = 5.0
		race_dict["Intelligence"] = -15.0
		race_dict["Dexterity"] = -10.0
	elif race_name == "Orc":
		race_dict["Strength"] = 5.0
		race_dict["Constitution"] = 0.0
		race_dict["Intelligence"] = -10.0
		race_dict["Dexterity"] = 5.0
	races[race_name] = race_dict		

# Parse profession modifiers and store in a data structure.
professions_path = args.crsm_dir + "\\csm"
print("Parsing character profession modifiers in \"" + professions_path + "\"")
for filename in os.listdir(professions_path):
	profession_name = filename.split('.')[0]
	profession_dict = dict()
	profession_file_path = professions_path + "\\" + filename
	profession_data = open(profession_file_path).readlines()
	for line in profession_data:
		if line == profession_data[0]:
			hppl, mppl = line.split("  ")
			profession_dict['Hppl'] = hppl.rstrip("\n")
			profession_dict['Mppl'] = mppl.rstrip("\n")
		else:
			skill, mod = line.rstrip("\n").rsplit(None,1)
			if skill in skills:
				profession_dict[skill] = float(mod)
	professions[profession_name] = profession_dict

# Parse profession base info and store in a data structure.
professions_base_path = args.crsm_dir + "\\crq"
print("Parsing character profession base info in \"" + professions_base_path + "\"")
for filename in os.listdir(professions_base_path):
	profession_base_name = filename.split('.')[0]
	profession_base_file_path = professions_base_path + "\\" + filename
	profession_base_data = open(profession_base_file_path).readlines()
	# State variables for parsing class skill options.
	optional = False
	base_mod = 0
	skill_options = []
	skill_option_index = 1
	for line in profession_base_data:
		if (len(line.split()) > 1):
			skill, mod = line.rstrip("\n").rsplit(None,1)
			if skill == 'STR':
				professions[profession_base_name]['Strength'] = float(mod)
			elif skill == 'DEX':
				professions[profession_base_name]['Dexterity'] = float(mod)
			elif skill == 'CON':
				professions[profession_base_name]['Constitution'] = float(mod)
			elif skill == 'INT':
				professions[profession_base_name]['Intelligence'] = float(mod)
			elif skill in skills:
				sos = "Skill_" + str(skill_option_index) + "_multi"
				professions[profession_base_name][sos] = "0"
				sos = "Skill_" + str(skill_option_index) + "_names"
				professions[profession_base_name][sos] = skill
				sos = "Skill_" + str(skill_option_index) + "_value"
				professions[profession_base_name][sos] = mod
				skill_option_index += 1
			elif skill == "*":
				sos = "Skill_" + str(skill_option_index) + "_multi"
				if mod == "0":
					professions[profession_base_name][sos] = "0"
				else:
					professions[profession_base_name][sos] = "1"
				sos = "Skill_" + str(skill_option_index) + "_value"
				professions[profession_base_name][sos] = mod
				skill_option_index += 1
		else:
			skill = line.rstrip("\n")
			if skill in skills:
				sos = "Skill_" + str(skill_option_index) + "_names"
				if sos in professions[profession_base_name]:
					professions[profession_base_name][sos] += ("-" + skill)
				else:
					professions[profession_base_name][sos] = skill

# Parse profession information and store in a data structure.
info_path = args.crsm_dir + "\\info"
print("Parsing character information in \"" + info_path + "\"")
for filename in os.listdir(info_path):
	info_name = filename.split('.')[0]
	info_file_name = info_path + "\\" + filename
	info_data = open(info_file_name).readlines()
	if (info_data[-1] == "\n"):
		info_data = info_data[0:-1]
	if info_name in professions:
		professions[info_name]['Description'] = info_data[-1].rstrip("\n").replace(",", " -")
	elif info_name in races:
		races[info_name]['Description'] = info_data[-1].rstrip("\n").replace(",", " -")

#pprint.pprint(races)
#pprint.pprint(professions)

# Write parsed race data to a csv file.
race_mod_file = open(race_modifiers_file, "w+")
race_mod_file.write("race")
for header_key in headers:
	if not "skill_" in header_key and not "hppl" in header_key and not "mppl" in header_key:
		race_mod_file.write("," + header_key)
race_mod_file.write("\n")
for race in sorted(races):
	race_mod_file.write(race)
	for skill in skills:
		if not "Skill_" in skill and not "Hppl" in skill and not "Mppl" in skill:
			if skill in races[race]:
				race_mod_file.write("," + str(races[race][skill]))
			else:
				race_mod_file.write(",0.0")
	race_mod_file.write("\n")

# Write parsed profession data to a csv file.
profession_mod_file = open(profession_modifiers_file, "w+")
profession_mod_file.write("profession")
for header_key in headers:
	profession_mod_file.write("," + header_key)
profession_mod_file.write("\n")
for profession in sorted(professions):
	profession_mod_file.write(profession)
	for skill in skills:
		if skill in professions[profession]:
			profession_mod_file.write("," + str(professions[profession][skill]))
		else:
			profession_mod_file.write(",0.0")
	profession_mod_file.write("\n")