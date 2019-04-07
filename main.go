// Filename:    main.go
// Author:      Joseph DeVictoria
// Date:        Mar_7_2019
// Description: Attempt to make Oldentide in a lower level engine.
//              Ack. Daniel Salvadori for bulk of base code.

package main

import (
	"flag"
	"fmt"
	"github.com/g3n/engine/audio"
	"github.com/g3n/engine/audio/al"
	"github.com/g3n/engine/audio/vorbis"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/camera/control"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/light"
	_ "github.com/g3n/engine/loader/obj"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/text"
	"github.com/g3n/engine/util/logger"
	"github.com/g3n/engine/window"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
	"Oldentide/shared"
	"net"
	"github.com/vmihailenco/msgpack"
)

var log *logger.Logger

// Client Game State Unique States
const (
	LOGIN_SCREEN                uint8 = iota
	CHARACTER_SELECT            uint8 = iota
	CHARACTER_CREATION_PHYSICAL uint8 = iota
	CHARACTER_CREATION_POINTS   uint8 = iota
	LOADING                     uint8 = iota
	IN_WORLD                    uint8 = iota
)

type OldentideClientGamestate struct {
	// g3n window and engine components
	wmgr              window.IWindowManager
	win               window.IWindow
	gs                *gls.GLS
	renderer          *renderer.Renderer
	scene             *core.Node
	camera            *camera.Perspective
	last_camera_angle *math32.Vector3
	orbit_control     *control.OrbitControl
	stepDelta         *math32.Vector2
	levelScene        *core.Node

	// Sound & Music players
	loginMusicPlayer *audio.Player
	sound_device     *al.Device
	audio_context    *al.Context

	// Fonts
	font *text.Font

	// Gui components
	root                        *gui.Root
	user_dialog                 *gui.Label
	login_menu                  *gui.Panel
	login_username              *gui.Edit
	login_password_edit         *gui.Edit
	login_server_address        *gui.Edit
	login_server_web_port       *gui.Edit
	login_server_game_port      *gui.Edit
	login_login_button          *gui.Button
	login_quit_button           *gui.Button
	login_process               *gui.Panel
	login_process_slider        *gui.Slider
	login_process_status        *gui.Label
	cs_menu                     *gui.Panel
	cs_characters               *gui.List
	cs_enter_world_button       *gui.Button
	cs_create_character_button  *gui.Button
	cs_back_button              *gui.Button
	cs_quit_button              *gui.Button
	cc_physical_menu            *gui.Panel
	cc_physical_first_name      *gui.Edit
	cc_physical_last_name       *gui.Edit
	cc_physical_model_panel     *gui.Panel
	cc_physical_male_button     *gui.ImageButton
	cc_physical_female_button   *gui.ImageButton
	cc_physical_hair_left       *gui.Button
	cc_physical_hair_right      *gui.Button
	cc_physical_skin_left       *gui.Button
	cc_physical_skin_right      *gui.Button
	cc_physical_dwarf_button    *gui.ImageButton
	cc_physical_elf_button      *gui.ImageButton
	cc_physical_gnome_button    *gui.ImageButton
	cc_physical_human_button    *gui.ImageButton
	cc_physical_leshy_button    *gui.ImageButton
	cc_physical_ogre_button     *gui.ImageButton
	cc_physical_orc_button      *gui.ImageButton
	cc_physical_back_button     *gui.Button
	cc_physical_next_button     *gui.Button
	cc_physical_quit_button     *gui.Button
	cc_points_menu              *gui.Panel
	cc_physical_profession_list *gui.List
	cc_points_attribute_panel   *gui.Panel
	cc_points_skill_panel       *gui.Panel
	cc_points_option_panel      *gui.Panel
	cc_points_back_button       *gui.Button
	cc_points_finish_button     *gui.Button
	cc_points_quit_button       *gui.Button

	// Game specific components
	assets_dir                 string
	cursor                     int
	login_password             string
	session_id                 int
	client_game_state          uint8
	new_character_firstname    string
	new_character_lastname     string
	new_character_sex          string
	new_character_race         string
	new_character_hair         string
	new_character_skin         string
	new_character_profession   string
	new_character_elective_one string
	new_character_elective_two string
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func (ogs *OldentideClientGamestate) UserMsg(msg string) {
	fmt.Println(msg)
	ogs.user_dialog.SetText(msg)

	// Center message horizontally
	if ogs.user_dialog.Width() < ogs.root.Width() {
		ogs.user_dialog.SetPositionX((ogs.root.Width() - ogs.user_dialog.Width()) / 2)
	} else {
		fmt.Println("UserMsg is too large to fit on the screen.")
		ogs.user_dialog.SetPositionX(0)
	}

	ogs.user_dialog.SetEnabled(true)
	ogs.root.Add(ogs.user_dialog)

	// TODO: Multiple messages will interfere with each other, so have new
	// msg cancel old message timeout
	go func() {
		time.Sleep(5 * time.Second)
		if ogs.user_dialog.Enabled() {
			ogs.user_dialog.SetEnabled(false)
			ogs.root.Remove(ogs.user_dialog)
		}
	}()
}

func (ogs *OldentideClientGamestate) UpdateLoginStatus(pct float32, login_step string) {
	ogs.login_process_slider.SetValue(pct)
	ogs.login_process_slider.SetText(fmt.Sprintf("%3.0f", ogs.login_process_slider.Value()*100))
	ogs.login_process_status.SetText(login_step)
}

func (ogs *OldentideClientGamestate) Login() {
	if ogs.client_game_state != LOGIN_SCREEN {
		return
	}

	username := ogs.login_username.Text()

	if !shared.ValidateUsername(username) {
		ogs.UserMsg("Username must be 3-30 alphanumeric characters")
		return
	}

	ogs.root.Remove(ogs.login_menu)
	ogs.login_menu.SetEnabled(false)
	ogs.root.Add(ogs.login_process)
	ogs.login_process.SetEnabled(true)

	steps := 5
	step := 0

	ogs.UpdateLoginStatus(float32(step)/float32(steps), "Extracting Credentials")
	step += 1
	fmt.Println("Username: ", username)
	fmt.Println("Password: ", ogs.login_password)
	fmt.Println("Server Address: ", ogs.login_server_address.Text())
	fmt.Println("Server Web Port: ", ogs.login_server_web_port.Text())
	fmt.Println("Server Game Port: ", ogs.login_server_game_port.Text())

	ogs.UpdateLoginStatus(float32(step)/float32(steps), "Checking Login Status")
	step += 1
	addrport := ogs.login_server_address.Text() + ":" + ogs.login_server_web_port.Text()
	login_server_page := "http://" + addrport + "/login"
	resp, err := http.PostForm(login_server_page, url.Values{"username": {username}, "password": {ogs.login_password}})
	if err != nil {
		ogs.root.Add(ogs.login_menu)
		ogs.login_menu.SetEnabled(true)
	}

	if resp == nil {
		msg := fmt.Sprintf("No response from %s", addrport)
		ogs.UserMsg(msg)
		return
	}

	fmt.Println(resp)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	shared.IfErrPrintErr(err)

	if resp.StatusCode != 200 {
		ogs.UserMsg("Remote login  for user `"+username+"` failed: "+string(body))

		ogs.root.Remove(ogs.login_process)
		ogs.login_process.SetEnabled(false)
		ogs.root.Add(ogs.login_menu)
		ogs.login_menu.SetEnabled(true)
		return
	}

	// fmt.Println("Body", body)
	// fmt.Println("Body str", string(body))
	session_id, err := strconv.Atoi(string(body))
	shared.IfErrPrintErr(err)

	fmt.Println("Login was a success! Session ID = ", session_id)
	ogs.session_id = session_id


	ogs.UpdateLoginStatus(float32(step)/float32(steps), "Saving Account Information")
	step += 1

	ogs.UpdateLoginStatus(float32(step)/float32(steps), "Loading Character List")
	step += 1

	ogs.root.Remove(ogs.login_menu) // Debug.
	ogs.login_menu.SetEnabled(false)
	ogs.root.Remove(ogs.login_process)
	ogs.login_process.SetEnabled(false)
	ogs.root.Add(ogs.cs_menu)
	ogs.cs_menu.SetEnabled(true)

	// Get all the characters for the account
	characters := ogs.GetCharacters(username)
	for _, character := range characters {
		name := fmt.Sprintf("%s %s", character.Firstname, character.Lastname)
		ogs.cs_characters.Add(gui.NewLabelWithFont(name, ogs.font))
	}

	stop := false
	for stop == false {
		if ogs.cs_characters.Selected() != nil {
			//fmt.Println(ogs.cs_characters.Selected().item.Text())
			stop = true
			/*if (ogs.cs_characters.Selected()[0].Text() == "Character_4") {
				stop = true
			}*/
		}
	}
}

func (ogs *OldentideClientGamestate) SendPacket(packet interface{}) {
	// Set up server connection.
	// Create udp socket description struct.
	addr := ogs.login_server_address.Text()
	port := ogs.login_server_game_port.Text()
	addrport := addr+":"+port
	fmt.Println("Sending packet to", addrport)

	server_connection, err := net.Dial("udp", addrport)
	defer server_connection.Close()
	shared.CheckErr(err)

	pac, err := msgpack.Marshal(packet)
	shared.CheckErr(err)

	server_connection.Write(pac)
}

func (ogs *OldentideClientGamestate) CreateCharacter(
	fname string,
	lname string,
	sex string,
	race string,
	skin string) {
	fmt.Println("Creating Character named", fname, lname)

	player := shared.Make_player(
		fname,
		lname,
		sex,
		race,
		skin)
	packet := shared.Create_player_packet{Opcode: shared.CREATEPLAYER, Pc: player}
	ogs.SendPacket(&packet)
}

func (ogs *OldentideClientGamestate) GetCharacters(account string) []shared.Pc {
	// TODO: Convert this to HTTP request instead of UDP?
	fmt.Println("Querying all characters for account", account)
	packet := shared.Req_clist_packet{Opcode: shared.REQCLIST, Account: account}
	ogs.SendPacket(&packet)

	// TODO: Block on packet for character list. In the meantime, bogus data
	var players []shared.Pc
	players = append(players, shared.Test_make_player("Michael"))
	players = append(players, shared.Test_make_player("Joseph"))
	return players
}

func (ogs *OldentideClientGamestate) EnterWorld() {
	fmt.Println("Entering World")
	ogs.client_game_state = LOADING
	ogs.orbit_control.Enabled = true
	ogs.client_game_state = IN_WORLD

	// Create a blue torus and add it to the scene
	geom := geometry.NewTorus(1, 0.4, 12, 32, math32.Pi*2)
	mat := material.NewPhong(math32.NewColor("DarkBlue"))
	torusMesh := graphic.NewMesh(geom, mat)
	ogs.scene.Add(torusMesh)

	// Add a point light to the scene
	pointLight := light.NewPoint(&math32.Color{1, 1, 1}, 5.0)
	pointLight.SetPosition(1, 0, 2)
	ogs.scene.Add(pointLight)

	// Add an axis helper to the scene
	axis := graphic.NewAxisHelper(0.5)
	ogs.scene.Add(axis)
}

// Quit saves the user data and quits the game
func (ogs *OldentideClientGamestate) Quit() {
	fmt.Println("Bye!")
	// Stop audio.
	if ogs.loginMusicPlayer.State() > 0 {
		ogs.loginMusicPlayer.Stop()
	}
	// Destroy the audio context.
	al.DestroyContext(ogs.audio_context)
	// Close the sound device.
	err := al.CloseDevice(ogs.sound_device)
	checkErr(err)
	// Close the window.
	ogs.win.SetShouldClose(true)
}

// LoadSkybox loads a space skybox and adds it to the scene
func (ogs *OldentideClientGamestate) LoadSkyBox(skybox_name string, file_format string) {
	log.Debug("Creating Skybox...")

	// Load skybox textures.
	skyboxDir := (ogs.assets_dir + "/Textures/Skyboxes/" + skybox_name + "/")
	skyboxData := graphic.SkyboxData{skyboxDir, file_format, [6]string{"px", "nx", "py", "ny", "pz", "nz"}}
	skybox, err := graphic.NewSkybox(skyboxData)
	checkErr(err)

	// Add skybox to scene.
	skybox.SetRenderOrder(-1) // The skybox should always be rendered first
	ogs.scene.Add(skybox)

	log.Debug("Done creating skybox")
}

// Attemplts to load audio libraries (from .dll).
func (ogs *OldentideClientGamestate) LoadAudioLibs() error {
	// Open default audio device
	var err error
	ogs.sound_device, err = al.OpenDevice("")
	if ogs.sound_device == nil {
		return fmt.Errorf("Error: %s opening OpenAL default device", err)
	}

	// Create audio context
	ogs.audio_context, err = al.CreateContext(ogs.sound_device, nil)
	if err != nil {
		return fmt.Errorf("Error creating audio context:%s", err)
	}

	// Make the context the current one
	err = al.MakeContextCurrent(ogs.audio_context)
	if err != nil {
		return fmt.Errorf("Error setting audio context current:%s", err)
	}

	log.Debug("%s version: %s", al.GetString(al.Vendor), al.GetString(al.Version))
	log.Debug("%s", vorbis.VersionString())
	return nil
}

// LoadAudio loads music and sound effects
func (ogs *OldentideClientGamestate) LoadAudio() {
	log.Debug("Load Audio")

	// Create listener and add it to the current camera
	listener := audio.NewListener()
	cdir := ogs.camera.Direction()
	listener.SetDirectionVec(&cdir)
	ogs.camera.GetCamera().Add(listener)

	// Helper function to create player and handle errors
	createPlayer := func(fname string) *audio.Player {
		log.Debug("Loading " + fname)
		p, err := audio.NewPlayer(fname)
		if err != nil {
			log.Error("Failed to create player for: %v", fname)
		}
		return p
	}

	// Load each music file you want to use into memory here. (Be a hog, they are small files!)
	ogs.loginMusicPlayer = createPlayer(ogs.assets_dir + "/Music/Komiku__End_of_the_trip.ogg")
	ogs.loginMusicPlayer.SetGain(0.05)
	ogs.loginMusicPlayer.SetLooping(true)
}

// Update updates the current level if any
func (ogs *OldentideClientGamestate) Update(timeDelta float64) {
	return
}

// ToggleFullScreen toggles whether is game is fullscreen or windowed
func (ogs *OldentideClientGamestate) ToggleFullScreen() {
	log.Debug("Toggle FullScreen")
	ogs.win.SetFullScreen(!ogs.win.FullScreen())
}

// onKey handles keyboard events for the game
func (ogs *OldentideClientGamestate) onKey(evname string, ev interface{}) {
	if ogs.client_game_state == IN_WORLD {
		kev := ev.(*window.KeyEvent)
		switch kev.Keycode {
		case window.KeyEscape:
			fmt.Println("Escape Key Pressed.")
		case window.KeyP:
			fmt.Println("P Key Pressed.") // Bring up party menu.
		case window.KeyC:
			fmt.Println("C Key Pressed.") // Bring up character menu.
		case window.KeyI:
			fmt.Println("I Key Pressed.") // Bring up inventory menu.
		case window.KeySlash:
			fmt.Println("/ Key Pressed.") // Focus the chat window and write "/".
		case window.KeyM:
			fmt.Println("M Key Pressed.") // Bring up the map.
		case window.KeyB:
			fmt.Println("B Key Pressed.") // Bring up the spellbook (spells and manuevers).
		}
	}
}

// onMouse handles mouse events for the game
func (ogs *OldentideClientGamestate) onMouse(evname string, ev interface{}) {
	mev := ev.(*window.MouseEvent)
	fmt.Println(mev)
}

// onCursor handles cursor movement for the game
func (ogs *OldentideClientGamestate) onCursor(evname string, ev interface{}) {
	// Calculate direction of potential movement based on camera angle
	var dir math32.Vector3
	ogs.camera.WorldDirection(&dir)
	ogs.stepDelta.Set(0, 0)

	if math32.Abs(dir.Z) > math32.Abs(dir.X) {
		if dir.Z > 0 {
			ogs.stepDelta.Y = 1
		} else {
			ogs.stepDelta.Y = -1
		}
	} else {
		if dir.X > 0 {
			ogs.stepDelta.X = 1
		} else {
			ogs.stepDelta.X = -1
		}
	}
}

// RenderFrame renders a frame of the scene with the GUI overlaid
func (ogs *OldentideClientGamestate) RenderFrame() {

	// Process GUI timers
	ogs.root.TimerManager.ProcessTimers()

	// Render the scene/gui using the specified camera
	rendered, err := ogs.renderer.Render(ogs.camera)
	checkErr(err)

	// Check I/O events
	ogs.wmgr.PollEvents()

	// Update window if necessary
	if rendered {
		ogs.win.SwapBuffers()
	}
}

func main() {
	runtime.LockOSThread() // OpenGL functions must be executed in the same thread as window.New()

	// Parse command line flags
	showLog := flag.Bool("debug", false, "display the debug log")
	flag.Parse()

	// Create logger
	log = logger.New("Oldentide", nil)
	log.AddWriter(logger.NewConsole(false))
	log.SetFormat(logger.FTIME | logger.FMICROS)
	if *showLog == true {
		log.SetLevel(logger.DEBUG)
	} else {
		log.SetLevel(logger.INFO)
	}
	log.Info("Initializing Oldentide")

	// Create Client Gamestate struct
	ogs := new(OldentideClientGamestate)

	// Manually scan the $GOPATH directories to find the data directory
	rawPaths := os.Getenv("GOPATH")
	paths := strings.Split(rawPaths, ":")
	for _, j := range paths {
		// Checks Assets path
		path := filepath.Join(j, "src/Oldentide/Assets")
		if _, err := os.Stat(path); err == nil {
			ogs.assets_dir = path
		}
	}

	// Get the window manager
	var err error
	ogs.wmgr, err = window.Manager("glfw")
	checkErr(err)

	// Create window and OpenGL context
	ogs.win, err = ogs.wmgr.CreateWindow(1280, 720, "Oldentide", false)
	checkErr(err)

	// Set window icon.
	// ogs.wmgr.

	// Create OpenGL state
	ogs.gs, err = gls.New()
	checkErr(err)

	// Set Cursor to Oldentide Cursor
	ogs.cursor, err = ogs.wmgr.CreateCursor(ogs.assets_dir+"/Interface/OldentideCursor30.png", 0, 0)
	if err != nil {
		fmt.Println("Error creating cursor.")
		//log.Log().Fatal("Error creating cursor: %s", err)
	}
	ogs.win.SetCustomCursor(ogs.cursor)

	// Build Font
	ogs.font, err = text.NewFont(ogs.assets_dir + "/Fonts/Sitka.ttf")
	checkErr(err)

	// Speed up a bit by not checking OpenGL errors
	ogs.gs.SetCheckErrors(false)

	// Sets window background color
	ogs.gs.ClearColor(0.1, 0.1, 0.1, 1.0)

	// Sets the OpenGL viewport size the same as the window size
	// This normally should be updated if the window is resized.
	width, height := ogs.win.Size()
	ogs.gs.Viewport(0, 0, int32(width), int32(height))

	// Creates GUI root panel
	ogs.root = gui.NewRoot(ogs.gs, ogs.win)
	ogs.root.SetSize(float32(width), float32(height))

	// Subscribe to window resize events. When the window is resized:
	// - Update the viewport size
	// - Update the root panel size
	// - Update the camera aspect ratio
	ogs.win.Subscribe(window.OnWindowSize, func(evname string, ev interface{}) {
		width, height := ogs.win.Size()
		ogs.gs.Viewport(0, 0, int32(width), int32(height))
		ogs.root.SetSize(float32(width), float32(height))
		aspect := float32(width) / float32(height)
		ogs.camera.SetAspect(aspect)
	})

	// Subscribe window to events
	ogs.win.Subscribe(window.OnKeyDown, ogs.onKey)
	ogs.win.Subscribe(window.OnMouseUp, ogs.onMouse)
	ogs.win.Subscribe(window.OnMouseDown, ogs.onMouse)

	// Creates a renderer and adds default shaders
	ogs.renderer = renderer.NewRenderer(ogs.gs)
	err = ogs.renderer.AddDefaultShaders()
	checkErr(err)
	ogs.renderer.SetGui(ogs.root)

	// Adds a perspective camera to the scene
	// The camera aspect ratio should be updated if the window is resized.
	aspect := float32(width) / float32(height)
	ogs.camera = camera.NewPerspective(65, aspect, 0.01, 1000)
	ogs.camera.SetPosition(0, 4, 5)
	ogs.camera.LookAt(&math32.Vector3{0, 5, 0})

	// Create orbit control and set limits
	ogs.orbit_control = control.NewOrbitControl(ogs.camera, ogs.win)
	ogs.orbit_control.Enabled = false
	ogs.orbit_control.EnablePan = false
	ogs.orbit_control.MaxPolarAngle = 2 * math32.Pi / 3
	ogs.orbit_control.MinDistance = 5
	ogs.orbit_control.MaxDistance = 15

	// Create main scene and child levelScene
	ogs.scene = core.NewNode()
	ogs.levelScene = core.NewNode()
	ogs.scene.Add(ogs.camera)
	ogs.scene.Add(ogs.levelScene)
	ogs.stepDelta = math32.NewVector2(0, 0)
	ogs.renderer.SetScene(ogs.scene)

	// Add white ambient light to the scene
	ambLight := light.NewAmbient(&math32.Color{1.0, 1.0, 1.0}, 0.4)
	ogs.scene.Add(ambLight)

	// Try to open audio libraries
	err = ogs.LoadAudioLibs()
	if err != nil {
		log.Error("%s", err)
	} else {
		ogs.LoadAudio()
		// Queue the music!
		ogs.loginMusicPlayer.Play()
	}

	ogs.SetupGui(width, height)
	ogs.RenderFrame()

	ogs.LoadSkyBox("Blue_Clouds", "jpg")
	//ogs.LoadSkyBox("Sunny_High_Plains", "jpg")

	ogs.win.Subscribe(window.OnCursor, ogs.onCursor)

	// Done Loading - Show login window.
	ogs.client_game_state = LOGIN_SCREEN
	ogs.root.Add(ogs.login_menu)

	now := time.Now()
	newNow := time.Now()
	log.Info("Starting Render Loop")

	// Start the render loop
	for !ogs.win.ShouldClose() {

		newNow = time.Now()
		timeDelta := now.Sub(newNow)
		now = newNow
		//ogs.last_camera_angle.Copy(*ogs.camera.Target())
		//ogs.last_camera_angle.X += 1
		if ogs.client_game_state == LOGIN_SCREEN ||
			ogs.client_game_state == CHARACTER_SELECT ||
			ogs.client_game_state == CHARACTER_CREATION_PHYSICAL ||
			ogs.client_game_state == CHARACTER_CREATION_POINTS ||
			ogs.client_game_state == LOADING {
			lca := ogs.camera.Target()
			//fmt.Println(lca)
			lca.X += 0.001
			ogs.camera.LookAt(&lca)
		}
		ogs.Update(timeDelta.Seconds())
		ogs.RenderFrame()
	}
}
