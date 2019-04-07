// Filename:    gui.go
// Author:      Joseph DeVictoria
// Date:        Mar_27_2019
// Description: Separate out the gui components because they are so ugly.

package main

import (
	"fmt"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/text"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/window"
	"Oldentide/shared"
)
// Colors:
var	interface_style_accent_0 = math32.Color4{0.22, 0.29, 0.29, 1.0}
var interface_style_brown_0 = math32.Color4{0.51, 0.40, 0.34, 1.0}
var interface_style_brown_1 = math32.Color4{0.47, 0.35, 0.20, 1.0}
var interface_style_brown_2 = math32.Color4{0.30, 0.23, 0.18, 1.0}
var interface_style_brown_3 = math32.Color4{0.19, 0.13, 0.09, 1.0}
var interface_style_brown_4 = math32.Color4{0.31, 0.20, 0.12, 1.0}
var color_warning = math32.Color4{1.0, 1.0, 0.5, 1.0}
var color_black = math32.Color4{0, 0, 0, 1.0}

// SetupGui creates all user interface elements
func (ogs *OldentideClientGamestate) SetupGui(width, height int) {
	log.Debug("Creating GUI...")
	var err error

	// Color Styles:
	//ss := gui.StyleDefault().Slider
	//ss.

	// Layout styles:
	g_2_layout := gui.NewGridLayout(2)
	g_2_layout.SetAlignH(gui.AlignCenter)
	g_2_layout.SetAlignV(gui.AlignCenter)
	h_layout := gui.NewHBoxLayout()
	h_layout.SetAlignH(gui.AlignWidth)
	v_layout := gui.NewVBoxLayout()
	v_layout.SetAlignV(gui.AlignHeight)

	// Login menu -------------------------------
	ogs.login_menu = gui.NewPanel(600, 300)
	ogs.login_menu.SetLayout(h_layout)
	ogs.login_menu.SetColor4(&interface_style_brown_2)
	ogs.login_menu.SetBorders(3, 3, 3, 3)
	ogs.login_menu.SetBordersColor4(&color_black)
	ogs.root.Subscribe(gui.OnResize, func(evname string, ev interface{}) {
		ogs.login_menu.SetPositionX((ogs.root.Width() - ogs.login_menu.Width()) / 2)
		ogs.login_menu.SetPositionY((ogs.root.Height() - ogs.login_menu.Height()) / 2)
	})

	// Audio Control - Slider
	audio_control := CreateAudioControl(ogs.loginMusicPlayer.Gain())
	audio_control.Subscribe(gui.OnMouseUp, func (name string, ev interface{}) {
		ogs.loginMusicPlayer.SetGain(audio_control.Value())
	})
	audio_control.Subscribe(gui.OnScroll, func (name string, ev interface{}) {
		ogs.loginMusicPlayer.SetGain(audio_control.Value())
	})
	ogs.root.Add(audio_control)

	// User Dialog Box
	ogs.user_dialog = CreateUserMsgDialog(ogs.font)

	login_left := gui.NewPanel(300, 300)
	login_left.SetLayout(v_layout)

	login_left.Add(gui.NewLabelWithFont("Username:", ogs.font))
	ogs.login_username = gui.NewEdit(250, "username")
	ogs.login_username.SetFontSize(16)
	ogs.login_username.SetPosition(15, 0)
	ogs.login_username.Subscribe(gui.OnKeyDown, func(name string, ev interface{}) {
		kev := ev.(*window.KeyEvent)
		switch kev.Keycode {
		case window.KeyTab:
			ogs.root.SetKeyFocus(ogs.login_password_edit)
		case window.KeyEnter:
			fmt.Println("Enter pressed on username!!!")
			ogs.Login()
		}
	})
	login_left.Add(ogs.login_username)

	login_left.Add(gui.NewLabelWithFont("Password:", ogs.font))
	ogs.login_password_edit = gui.NewEdit(250, "password")
	ogs.login_password_edit.SetFontSize(16)
	ogs.login_password_edit.Subscribe(gui.OnChar, func(name string, ev interface{}) {
		input_char := string(ev.(*window.CharEvent).Char)
		ogs.login_password += input_char
		ogs.login_password_edit.CursorBack()
		ogs.login_password_edit.CursorInput("*")
		fmt.Println("Typed something in password box.", ogs.login_password)
	})
	ogs.login_password_edit.Subscribe(gui.OnKeyDown, func(name string, ev interface{}) {
		kev := ev.(*window.KeyEvent)
		switch kev.Keycode {
		case window.KeyBackspace:
			if pwl := len(ogs.login_password) - 1; pwl >= 0 {
				ogs.login_password = ogs.login_password[:pwl]
			}
		// TODO handle case where someone is deleting instead of backspacing.
		// TODO handle case where someone highlights the whole string and deletes it.
		case window.KeyEnter:
			fmt.Println("Enter pressed on password!!!")
			ogs.Login()
		}
	})
	login_left.Add(ogs.login_password_edit)

	login_buttons := gui.NewPanel(300, 40)
	login_buttons.SetLayout(h_layout)

	ogs.login_login_button = gui.NewButton("Login")
	ogs.login_login_button.Label.SetFontSize(16)
	ogs.login_login_button.Subscribe(gui.OnClick, func(name string, ev interface{}) {
		fmt.Println("Login button was pressed.")
		ogs.Login()
	})
	login_buttons.Add(ogs.login_login_button)

	ogs.login_quit_button = gui.NewButton("Quit")
	ogs.login_quit_button.Label.SetFontSize(16)
	ogs.login_quit_button.Subscribe(gui.OnClick, func(name string, ev interface{}) {
		ogs.Quit()
	})
	login_buttons.Add(ogs.login_quit_button)

	login_left.Add(login_buttons)

	login_right := gui.NewPanel(300, 300)
	login_right.SetLayout(v_layout)

	login_right.Add(gui.NewLabelWithFont("Server Address:", ogs.font))
	ogs.login_server_address = gui.NewEdit(250, "imp.oldentide.com")
	ogs.login_server_address.SetFontSize(16)
	ogs.login_server_address.SetText("imp.oldentide.com")
	login_right.Add(ogs.login_server_address)

	login_right.Add(gui.NewLabelWithFont("Server Web Port:", ogs.font))
	ogs.login_server_web_port = gui.NewEdit(250, "80")
	ogs.login_server_web_port.SetFontSize(16)
	ogs.login_server_web_port.SetText("80")
	login_right.Add(ogs.login_server_web_port)

	login_right.Add(gui.NewLabelWithFont("Server Game Port:", ogs.font))
	ogs.login_server_game_port = gui.NewEdit(250, "1337")
	ogs.login_server_game_port.SetFontSize(16)
	ogs.login_server_game_port.SetText("1337")
	login_right.Add(ogs.login_server_game_port)

	ogs.login_menu.Add(gui.NewPanel(25, 0))
	ogs.login_menu.Add(login_left)
	ogs.login_menu.Add(login_right)

	ogs.login_process = gui.NewPanel(400, 100)
	ogs.login_process.SetLayout(h_layout)
	ogs.login_process.SetColor4(&interface_style_brown_2)
	ogs.login_process.SetBorders(3, 3, 3, 3)
	ogs.login_process.SetBordersColor4(&color_black)
	ogs.root.Subscribe(gui.OnResize, func(evname string, ev interface{}) {
		ogs.login_process.SetPositionX((ogs.root.Width() - ogs.login_process.Width()) / 2)
		ogs.login_process.SetPositionY((ogs.root.Height() - ogs.login_process.Height()) / 2)
	})

	login_process_right := gui.NewPanel(375, 100)
	login_process_right.SetLayout(v_layout)

	login_process_right.Add(gui.NewLabelWithFont("Login Process:", ogs.font))
	ogs.login_process_slider = gui.NewHSlider(350, 35)
	ogs.login_process_slider.SetValue(0.0)
	ogs.login_process_slider.SetEnabled(false)
	ogs.login_process_slider.SetText(fmt.Sprintf("%3.0f", ogs.login_process_slider.Value()*100))
	login_process_right.Add(ogs.login_process_slider)
	ogs.login_process_status = gui.NewLabelWithFont("Waiting", ogs.font)
	login_process_right.Add(ogs.login_process_status)

	ogs.login_process.Add(gui.NewPanel(25, 0))
	ogs.login_process.Add(login_process_right)
	// ------------------------------------------

	// Character Select Menu --------------------
	ogs.cs_menu = gui.NewPanel(250, 400)
	ogs.cs_menu.SetLayout(v_layout)
	ogs.cs_menu.SetColor4(&interface_style_brown_2)
	ogs.cs_menu.SetBorders(3, 3, 3, 3)
	ogs.cs_menu.SetBordersColor4(&color_black)
	ogs.root.Subscribe(gui.OnResize, func(evname string, ev interface{}) {
		ogs.cs_menu.SetPositionX((ogs.root.Width() - ogs.cs_menu.Width()) / 2)
		ogs.cs_menu.SetPositionY((ogs.root.Height() - ogs.cs_menu.Height()) / 2)
	})

	cs_spacer_panel_top := gui.NewPanel(250, 10)
	ogs.cs_menu.Add(cs_spacer_panel_top)

	cs_list_panel := gui.NewPanel(250, 300)
	cs_list_panel.SetLayout(h_layout)

	ogs.cs_characters = gui.NewVList(200, 280)
	ogs.cs_characters.SetSingle(true)
	cs_list_panel.Add(ogs.cs_characters)

	ogs.cs_menu.Add(cs_list_panel)

	cs_spacer_panel_mid := gui.NewPanel(250, 10)
	ogs.cs_menu.Add(cs_spacer_panel_mid)

	cs_button_panel_top := gui.NewPanel(250, 50)
	cs_button_panel_top.SetLayout(h_layout)

	ogs.cs_enter_world_button = gui.NewButton("Enter World")
	ogs.cs_enter_world_button.Label.SetFontSize(16)
	ogs.cs_enter_world_button.Subscribe(gui.OnClick, func(name string, ev interface{}) {
		fmt.Println("Enter World button was pressed.")
		ogs.root.Remove(ogs.cs_menu)
		ogs.cs_menu.SetEnabled(false)
		ogs.EnterWorld()
	})
	cs_button_panel_top.Add(ogs.cs_enter_world_button)

	ogs.cs_create_character_button = gui.NewButton("Create Character")
	ogs.cs_create_character_button.Label.SetFontSize(16)
	ogs.cs_create_character_button.Subscribe(gui.OnClick, func(name string, ev interface{}) {
		fmt.Println("Create Character button was pressed.")
		ogs.root.Remove(ogs.cs_menu)
		ogs.cs_menu.SetEnabled(false)
		ogs.new_character_firstname = ""
		ogs.new_character_lastname = ""
		ogs.new_character_sex = "Male"
		ogs.new_character_race = "Human"
		ogs.new_character_hair = "Hair1"
		ogs.new_character_skin = "Skin1"
		ogs.root.Add(ogs.cc_physical_menu)
		ogs.cc_physical_menu.SetEnabled(true)
	})
	cs_button_panel_top.Add(ogs.cs_create_character_button)

	cs_button_panel_bottom := gui.NewPanel(250, 50)
	cs_button_panel_bottom.SetLayout(h_layout)

	ogs.cs_back_button = gui.NewButton("Back")
	ogs.cs_back_button.Label.SetFontSize(16)
	ogs.cs_back_button.Subscribe(gui.OnClick, func(name string, ev interface{}) {
		fmt.Println("Create Character Back button was pressed.")
		for ogs.cs_characters.ItemAt(0) != nil {
			ogs.cs_characters.RemoveAt(0)
		}
		ogs.login_password_edit.CursorEnd()
		for range ogs.login_password {
			ogs.login_password_edit.CursorBack()
		}
		ogs.login_password = ""
		ogs.root.Remove(ogs.cs_menu)
		ogs.cs_menu.SetEnabled(false)
		ogs.root.Add(ogs.login_menu)
		ogs.login_menu.SetEnabled(true)
	})
	cs_button_panel_bottom.Add(ogs.cs_back_button)

	ogs.cs_quit_button = gui.NewButton("Quit")
	ogs.cs_quit_button.Label.SetFontSize(16)
	ogs.cs_quit_button.Subscribe(gui.OnClick, func(name string, ev interface{}) {
		fmt.Println("Create Character Quit button was pressed.")
		ogs.Quit()
	})
	cs_button_panel_bottom.Add(ogs.cs_quit_button)

	ogs.cs_menu.Add(cs_button_panel_top)
	ogs.cs_menu.Add(cs_button_panel_bottom)
	// ------------------------------------------

	// Character Create Menu --------------------
	ogs.cc_physical_menu = gui.NewPanel(600, 400)
	ogs.cc_physical_menu.SetLayout(v_layout)
	ogs.cc_physical_menu.SetColor4(&interface_style_brown_2)
	ogs.cc_physical_menu.SetBorders(3, 3, 3, 3)
	ogs.cc_physical_menu.SetBordersColor4(&color_black)
	ogs.root.Subscribe(gui.OnResize, func(evname string, ev interface{}) {
		ogs.cc_physical_menu.SetPositionX((ogs.root.Width() - ogs.cc_physical_menu.Width()) / 2)
		ogs.cc_physical_menu.SetPositionY((ogs.root.Height() - ogs.cc_physical_menu.Height()) / 2)
	})

	cc_physical_header := gui.NewPanel(600, 30)
	cc_physical_header.SetLayout(h_layout)
	cc_physical_header.Add(gui.NewLabelWithFont("Character Creation (Physical)", ogs.font))
	ogs.cc_physical_menu.Add(cc_physical_header)

	cc_physical_name_bar := gui.NewPanel(600, 30)
	cc_physical_name_bar.SetLayout(h_layout)
	cc_physical_name_bar.Add(gui.NewLabelWithFont("First Name :", ogs.font))
	ogs.cc_physical_first_name = gui.NewEdit(180, "")
	ogs.cc_physical_first_name.SetFontSize(16)
	ogs.cc_physical_first_name.Subscribe(gui.OnKeyDown, func(name string, ev interface{}) {
		kev := ev.(*window.KeyEvent)
		switch kev.Keycode {
		case window.KeyTab:
			ogs.root.SetKeyFocus(ogs.cc_physical_last_name)
		}
	})
	cc_physical_name_bar.Add(ogs.cc_physical_first_name)
	ogs.cc_physical_last_name = gui.NewEdit(180, "")
	ogs.cc_physical_last_name.SetFontSize(16)
	cc_physical_name_bar.Add(ogs.cc_physical_last_name)
	cc_physical_name_bar.Add(gui.NewLabelWithFont(": Last Name", ogs.font))
	ogs.cc_physical_menu.Add(cc_physical_name_bar)

	cc_physical_main_bar := gui.NewPanel(600, 300)
	cc_physical_main_bar.SetLayout(h_layout)

	cc_physical_control_bar := gui.NewPanel(300, 300)
	cc_physical_control_bar.SetLayout(v_layout)

	cc_physical_sex_panel := gui.NewPanel(300, 30)
	cc_physical_sex_panel.SetLayout(h_layout)
	ogs.cc_physical_male_button, err = gui.NewImageButton(ogs.assets_dir + "/Interface/male.png")
	checkErr(err)
	ogs.cc_physical_male_button.Subscribe(gui.OnClick, func(name string, ev interface{}) {
		fmt.Println("Male sex selected")
		ogs.new_character_sex = "Male"
	})
	cc_physical_sex_panel.Add(ogs.cc_physical_male_button)
	ogs.cc_physical_female_button, err = gui.NewImageButton(ogs.assets_dir + "/Interface/female.png")
	checkErr(err)
	ogs.cc_physical_female_button.Subscribe(gui.OnClick, func(name string, ev interface{}) {
		fmt.Println("Female sex selected")
		ogs.new_character_sex = "Female"
	})
	cc_physical_sex_panel.Add(ogs.cc_physical_female_button)
	cc_physical_control_bar.Add(cc_physical_sex_panel)

	cc_physical_hair_panel := gui.NewPanel(300, 30)
	cc_physical_hair_panel.SetLayout(h_layout)
	ogs.cc_physical_hair_left = gui.NewButton("< Previous Hair")
	ogs.cc_physical_hair_left.Subscribe(gui.OnClick, func(name string, ev interface{}) {
		fmt.Println("Previous Hair Pressed")
		ogs.new_character_race = "-"
	})
	ogs.cc_physical_hair_right = gui.NewButton("Next Hair >")
	ogs.cc_physical_hair_right.Subscribe(gui.OnClick, func(name string, ev interface{}) {
		fmt.Println("Next Hair Pressed")
		ogs.new_character_hair += "+"
	})
	cc_physical_hair_panel.Add(ogs.cc_physical_hair_left)
	cc_physical_hair_panel.Add(ogs.cc_physical_hair_right)
	cc_physical_control_bar.Add(cc_physical_hair_panel)

	cc_physical_skin_panel := gui.NewPanel(300, 30)
	cc_physical_skin_panel.SetLayout(h_layout)
	ogs.cc_physical_skin_left = gui.NewButton("< Previous Skin")
	ogs.cc_physical_skin_left.Subscribe(gui.OnClick, func(name string, ev interface{}) {
		fmt.Println("Previous Skin Pressed")
		ogs.new_character_race = "-"
	})
	ogs.cc_physical_skin_right = gui.NewButton("Next Skin >")
	ogs.cc_physical_skin_right.Subscribe(gui.OnClick, func(name string, ev interface{}) {
		fmt.Println("Next Skin Pressed")
		ogs.new_character_hair += "+"
	})
	cc_physical_skin_panel.Add(ogs.cc_physical_skin_left)
	cc_physical_skin_panel.Add(ogs.cc_physical_skin_right)
	cc_physical_control_bar.Add(cc_physical_skin_panel)

	cc_physical_human_panel := gui.NewPanel(300, 30)
	cc_physical_human_panel.SetLayout(h_layout)
	ogs.cc_physical_human_button, err = gui.NewImageButton(ogs.assets_dir + "/Interface/human.png")
	checkErr(err)
	ogs.cc_physical_human_button.Subscribe(gui.OnClick, func(name string, ev interface{}) {
		fmt.Println("Human race selected")
		ogs.new_character_race = "Human"
	})
	cc_physical_human_panel.Add(ogs.cc_physical_human_button)
	cc_physical_control_bar.Add(cc_physical_human_panel)

	cc_physical_dwarf_elf_panel := gui.NewPanel(300, 30)
	cc_physical_dwarf_elf_panel.SetLayout(h_layout)
	ogs.cc_physical_dwarf_button, err = gui.NewImageButton(ogs.assets_dir + "/Interface/dwarf.png")
	checkErr(err)
	ogs.cc_physical_dwarf_button.Subscribe(gui.OnClick, func(name string, ev interface{}) {
		fmt.Println("Dwarf race selected")
		ogs.new_character_race = "Dwarf"
	})
	ogs.cc_physical_elf_button, err = gui.NewImageButton(ogs.assets_dir + "/Interface/elf.png")
	checkErr(err)
	ogs.cc_physical_elf_button.Subscribe(gui.OnClick, func(name string, ev interface{}) {
		fmt.Println("Elf race selected")
		ogs.new_character_race = "Elf"
	})
	cc_physical_dwarf_elf_panel.Add(ogs.cc_physical_dwarf_button)
	cc_physical_dwarf_elf_panel.Add(ogs.cc_physical_elf_button)
	cc_physical_control_bar.Add(cc_physical_dwarf_elf_panel)

	cc_physical_gnome_leshy_panel := gui.NewPanel(300, 30)
	cc_physical_gnome_leshy_panel.SetLayout(h_layout)
	ogs.cc_physical_gnome_button, err = gui.NewImageButton(ogs.assets_dir + "/Interface/gnome.png")
	checkErr(err)
	ogs.cc_physical_gnome_button.Subscribe(gui.OnClick, func(name string, ev interface{}) {
		fmt.Println("Gnome race selected")
		ogs.new_character_race = "Gnome"
	})
	ogs.cc_physical_leshy_button, err = gui.NewImageButton(ogs.assets_dir + "/Interface/leshy.png")
	checkErr(err)
	ogs.cc_physical_leshy_button.Subscribe(gui.OnClick, func(name string, ev interface{}) {
		fmt.Println("Leshy race selected")
		ogs.new_character_race = "Leshy"
	})
	cc_physical_gnome_leshy_panel.Add(ogs.cc_physical_gnome_button)
	cc_physical_gnome_leshy_panel.Add(ogs.cc_physical_leshy_button)
	cc_physical_control_bar.Add(cc_physical_gnome_leshy_panel)

	cc_physical_ogre_orc_panel := gui.NewPanel(300, 30)
	cc_physical_ogre_orc_panel.SetLayout(h_layout)
	ogs.cc_physical_ogre_button, err = gui.NewImageButton(ogs.assets_dir + "/Interface/ogre.png")
	checkErr(err)
	ogs.cc_physical_ogre_button.Subscribe(gui.OnClick, func(name string, ev interface{}) {
		fmt.Println("Ogre race selected")
		ogs.new_character_race = "Ogre"
	})
	ogs.cc_physical_orc_button, err = gui.NewImageButton(ogs.assets_dir + "/Interface/orc.png")
	checkErr(err)
	ogs.cc_physical_orc_button.Subscribe(gui.OnClick, func(name string, ev interface{}) {
		fmt.Println("Orc race selected")
		ogs.new_character_race = "Orc"
	})
	cc_physical_ogre_orc_panel.Add(ogs.cc_physical_ogre_button)
	cc_physical_ogre_orc_panel.Add(ogs.cc_physical_orc_button)
	cc_physical_control_bar.Add(cc_physical_ogre_orc_panel)

	cc_physical_main_bar.Add(cc_physical_control_bar)

	cc_physical_preview_bar := gui.NewPanel(260, 285)
	cc_physical_preview_bar.SetLayout(v_layout)
	cc_physical_preview_bar.SetColor4(&math32.Color4{0.202, 0.188, 0.179, 1.0})
	cc_physical_preview_bar.SetBorders(2, 2, 2, 2)
	cc_physical_preview_bar.SetBordersColor4(&color_black)
	cc_physical_main_bar.Add(cc_physical_preview_bar)

	ogs.cc_physical_menu.Add(cc_physical_main_bar)

	cc_physical_footer := gui.NewPanel(600, 40)
	cc_physical_footer.SetLayout(h_layout)
	ogs.cc_physical_back_button = gui.NewButton("Back")
	ogs.cc_physical_back_button.Subscribe(gui.OnClick, func(name string, ev interface{}) {
		fmt.Println("cc physical back button was pressed.")
		ogs.root.Remove(ogs.cc_physical_menu)
		ogs.cc_physical_menu.SetEnabled(false)
		ogs.root.Add(ogs.cs_menu)
		ogs.cs_menu.SetEnabled(true)
	})
	cc_physical_footer.Add(ogs.cc_physical_back_button)
	ogs.cc_physical_quit_button = gui.NewButton("Quit")
	ogs.cc_physical_quit_button.Subscribe(gui.OnClick, func(name string, ev interface{}) {
		fmt.Println("cc physical quit button was pressed.")
		ogs.Quit()
	})
	cc_physical_footer.Add(ogs.cc_physical_quit_button)
	ogs.cc_physical_next_button = gui.NewButton("Next")
	ogs.cc_physical_next_button.Subscribe(gui.OnClick, func(name string, ev interface{}) {
		fmt.Println("cc physical next button was pressed.")
		ogs.new_character_firstname = ogs.cc_physical_first_name.Text()
		ogs.new_character_lastname = ogs.cc_physical_last_name.Text()
		fmt.Println(ogs.new_character_firstname,
			ogs.new_character_lastname,
			ogs.new_character_sex,
			ogs.new_character_race,
			ogs.new_character_hair,
			ogs.new_character_skin)
		// Check that physical stats are valid
		if !ValidateCharacterPhysical(ogs) {
			return;
		}
		ogs.root.Remove(ogs.cc_physical_menu)
		ogs.cc_physical_menu.SetEnabled(false)
		ogs.root.Add(ogs.cc_points_menu)
		ogs.cc_points_menu.SetEnabled(true)
	})
	cc_physical_footer.Add(ogs.cc_physical_next_button)

	ogs.cc_physical_menu.Add(cc_physical_footer)

	ogs.cc_points_menu = gui.NewPanel(800, 450)
	ogs.cc_points_menu.SetLayout(v_layout)
	ogs.cc_points_menu.SetColor4(&interface_style_brown_2)
	ogs.cc_points_menu.SetBorders(3, 3, 3, 3)
	ogs.cc_points_menu.SetBordersColor4(&color_black)
	ogs.root.Subscribe(gui.OnResize, func(evname string, ev interface{}) {
		ogs.cc_points_menu.SetPositionX((ogs.root.Width() - ogs.cc_points_menu.Width()) / 2)
		ogs.cc_points_menu.SetPositionY((ogs.root.Height() - ogs.cc_points_menu.Height()) / 2)
	})

	cc_points_header := gui.NewPanel(800, 30)
	cc_points_header.SetLayout(h_layout)
	cc_points_header.Add(gui.NewLabelWithFont("Character Creation (Points)", ogs.font))
	ogs.cc_points_menu.Add(cc_points_header)

	cc_points_main_bar := gui.NewPanel(800, 380)
	ogs.cc_points_menu.Add(cc_points_main_bar)

	cc_points_footer := gui.NewPanel(800, 40)
	cc_points_footer.SetLayout(h_layout)
	ogs.cc_points_back_button = gui.NewButton("Back")
	ogs.cc_points_back_button.Subscribe(gui.OnClick, func(name string, ev interface{}) {
		fmt.Println("cc points back button was pressed.")
		ogs.root.Remove(ogs.cc_points_menu)
		ogs.cc_points_menu.SetEnabled(false)
		ogs.root.Add(ogs.cc_physical_menu)
		ogs.cc_physical_menu.SetEnabled(true)
	})
	cc_points_footer.Add(ogs.cc_points_back_button)
	ogs.cc_points_quit_button = gui.NewButton("Quit")
	ogs.cc_points_quit_button.Subscribe(gui.OnClick, func(name string, ev interface{}) {
		fmt.Println("cc points quit button was pressed.")
		ogs.Quit()
	})
	cc_points_footer.Add(ogs.cc_points_quit_button)
	ogs.cc_points_finish_button = gui.NewButton("Finish")
	ogs.cc_points_finish_button.Subscribe(gui.OnClick, func(name string, ev interface{}) {
		fmt.Println("cc points finish button was pressed.")
		// TODO Double check everything on the client side, then ask the server to try to create this character.
		// TODO if it doesn't work, throw an error.
		ogs.CreateCharacter(
			ogs.new_character_firstname,
			ogs.new_character_lastname,
			ogs.new_character_sex,
			ogs.new_character_race,
			ogs.new_character_skin)
		ogs.root.Remove(ogs.cc_points_menu)
		ogs.cc_points_menu.SetEnabled(false)
		ogs.root.Add(ogs.cs_menu)
		ogs.cs_menu.SetEnabled(true)
	})
	cc_points_footer.Add(ogs.cc_points_finish_button)

	ogs.cc_points_menu.Add(cc_points_footer)

	// ------------------------------------------

	// Dispatch a fake OnResize event to update all subscribed elements
	ogs.root.Dispatch(gui.OnResize, nil)

	log.Debug("Done creating GUI.")
}

func CreateUserMsgDialog(font *text.Font) *gui.Label {
	el := gui.NewLabelWithFont("", font)
	el.SetBgColor4(&color_warning)
	el.SetColor4(&color_black)
	el.SetBorders(3, 3, 3, 3)
	el.SetFontSize(32.0)
	// Set DPI ?
	// el.SetFontDPI()
	el.SetBordersColor4(&color_black)
	return el
}

// Don't let user click NEXT unless stats are valid
func ValidateCharacterPhysical(ogs *OldentideClientGamestate) bool {
	if !shared.ValidateName(ogs.new_character_firstname) {
		ogs.UserMsg("First name must be 3-20 alphabetic characters")
		return false
	}

	if !shared.ValidateName(ogs.new_character_lastname) {
		ogs.UserMsg("Last name must be 3-20 alphabetic characters")
		return false
	}

	return true
}

func CreateAudioControl(initial_gain float32) *gui.Slider {
	el := gui.NewVSlider(0, 100)
	// The width will auto-expand to fit text
	el.SetText("Volume")
	// Set the audio control to an initial volume
	el.SetValue(initial_gain)
	return el
}
