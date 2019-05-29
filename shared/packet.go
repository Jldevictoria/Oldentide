// Filename:    packet.go
// Author:      Joseph DeVictoria
// Date:        June_16_2018
// Purpose:     Definition of packets to be used in Oldentide.

package shared

import (
	_ "fmt"
	_ "github.com/vmihailenco/msgpack"
	"net"
)

// Packet types. - These will probably need some tweaking in the future.
type Opcode uint8

const (
	EMPTY        Opcode = iota // BOTH - Ignore these.
	GENERIC      Opcode = iota // BOTH - Used for debug.
	ACK          Opcode = iota // BOTH - Acknoledge that packet was received.
	ERROR        Opcode = iota // BOTH - Pass error message.
	REQCLIST     Opcode = iota // CLIENT - Client requests character list.
	SENDCLIST    Opcode = iota // SERVER - Server sends character list.
	CREATEPLAYER Opcode = iota // CLIENT - Player sends the data for his newly created character.
	CONNECT      Opcode = iota // CLIENT - Player selects a character and connects to server.
	DISCONNECT   Opcode = iota // BOTH - Player disconnects from the server.
	SENDPLAYER   Opcode = iota // SERVER - Send all of the information to fully update the main player.
	SENDPC       Opcode = iota // SERVER - Send all the needed information to set up or uptade another PC.
	SENDNPC      Opcode = iota // SERVER - Send all the needed information to set up & update an NPC.
	MOVEPLAYER   Opcode = iota // CLIENT - Send position information based on player input/actions.
	SPENDDP      Opcode = iota // CLIENT - Send a request of dp expenditures for validation.
	TALKCMD      Opcode = iota // CLIENT - <PLACEHOLDER> Player initiated talk command (with npc).
	ATTACKCMD    Opcode = iota // CLIENT - Player initiated attack command (with npc or player).
	TRADECMD     Opcode = iota // CLIENT - Player initiated trade command (with player).
	INVITECMD    Opcode = iota // CLIENT - Player initiated invite command (with player).
	GINVITECMD   Opcode = iota // CLIENT - Player initiated guild invite command (with player).
	GKICK        Opcode = iota // CLIENT - Player initiated guild kick command (with player).
	GPROMOTE     Opcode = iota // CLIENT - Player initiated a guild officer promotion command (with player).
	GDEMOTE      Opcode = iota // CLIENT - Player initiated a guild officer demotion command (with player).
	SAYCMD       Opcode = iota // CLIENT - Player said something with /s command.
	YELLCMD      Opcode = iota // CLIENT - Player said something with /y command.
	OOCCMD       Opcode = iota // CLIENT - Player said something with /ooc command.
	HELPCMD      Opcode = iota // CLIENT - Player said something with /h command.
	PCHATCMD     Opcode = iota // CLIENT - Player seid something with /p command.
	GCHATCMD     Opcode = iota // CLIENT - Player said something with /g command.
	WHISPERCMD   Opcode = iota // CLIENT - Player said something with /w command (with player).
	RELAYSAY     Opcode = iota // SERVER - Relay a say command to proper clients.
	RELAYYELL    Opcode = iota // SERVER - Relay a yell command to proper clients.
	RELAYOOC     Opcode = iota // SERVER - Relay an ooc command to proper clients.
	RELAYHELP    Opcode = iota // SERVER - Relay a help command to proper clients.
	RELAYPCHAT   Opcode = iota // SERVER - Relay a part chat command to proper clients.
	RELAYGCHAT   Opcode = iota // SERVER - Relay a guild chat command to proper clients.
	RELAYWHISPER Opcode = iota // SERVER - Relay a whisper command to proper client.
	ACTIVATECMD  Opcode = iota // CLIENT - Player initiated game object activation command (with door/chest/switch).
	ENVUPDATE    Opcode = iota // SERVER - Send a flag to appropriate clients for an environemnt variable (door/chest/switch).
	DIALOGUETEXT Opcode = iota // SERVER - <PLACEHOLDER> Send dialog + options to client.
	DIALOGCMD    Opcode = iota // PLAYER - <PLACEHOLDER> Send dialog response from player.
	SENDITEM     Opcode = iota // SERVER - Server successfully awards a given item to a designated player.
	INITSHOP     Opcode = iota // SERVER - <PLACEHOLDER> Server has started a shop window.
	SHOPITEM     Opcode = iota // SERVER - <PLACEHOLDER> Send item to shop inventory.
	BUYITEM      Opcode = iota // CLIENT - <PLACEHOLDER> Player tries to purchase an item.
	INITLOOT     Opcode = iota // SERVER - <PLACEHOLDER> Server has started a loot window.
	LOOTITEM     Opcode = iota // SERVER - <PLACEHOLDER> Send item to loot inventory.
	TAKELOOT     Opcode = iota // CLIENT - <PLACEHOLDER> Player tried to take a loot item.
	INITTRADE    Opcode = iota // SERVER - <PLACEHOLDER> Server has started a trade window.
	OFFERITEM    Opcode = iota // CLIENT - <PLACEHOLDER> Player has offered an item for trade.
	PULLITEM     Opcode = iota // CLIENT - <PLACEHOLDER> Player has removed an item for trade.
	TRADEITEM    Opcode = iota // SERVER - <PLACEHOLDER> Server communicates item to trade window.
	REMITEM      Opcode = iota // SERVER - <PLACEHOLDER> Server removes item from trade window.
	ACCTRADE     Opcode = iota // CLIENT - <PLACEHOLDER> Player has accepted a trade.
	UNACCTRADE   Opcode = iota // CLIENT - <PLACEHOLDER> Player has unaccepted a trade..
	COMMTRADE    Opcode = iota // SERVER - <PLACEHOLDER> Communicate trade acceptance status to trade window..
	FINTRADE     Opcode = iota // SERVER - <PLACEHOLDER> Close trade window and inform player trade was accepted..
	INITCOMABT   Opcode = iota // SERVER - Initialize a combat session on client.
	ADDNPCCOMBAT Opcode = iota // SERVER - Add an NPC to a combat session on client.
	ADDPCCOMABT  Opcode = iota // SERVER - Add a PC to a combat session on client.
	REMNPCCOMBAT Opcode = iota // SERVER - Remove an NPC from a combat session on client.
	REMPCCOMBAT  Opcode = iota // SERVER - Remove a PC from a combat session on client.
	// Need to continue adding here...
	// CASTSPELL Opcode = , // CLIENT -
)

type Raw_packet struct {
	Size    int
	Client  *net.UDPAddr
	Payload []byte
}

type Opcode_packet struct {
	Opcode Opcode
}

type Empty_packet struct {
	Opcode Opcode
}

type Generic_packet struct {
	Opcode     Opcode
	Session_id int64
}

type Ack_packet struct {
	Opcode     Opcode
	Session_id int64
}

type Error_packet struct {
	Opcode     Opcode
	Session_id int64
}

type Req_clist_packet struct {
	Opcode  Opcode
	Account string
}

type Send_clist_packet struct {
	Opcode     Opcode
	Characters []string
}

type Create_player_packet struct {
	Opcode     Opcode
	Session_id int64
	Pc         Pc
}

type Connect_packet struct {
	Opcode     Opcode
	Session_id int64
	Character  string
}

type Disconnect_packet struct {
	Opcode     Opcode
	Session_id int64
}

type Send_player_packet struct {
	Opcode     Opcode
	Session_id int64
}

type Send_pc_packet struct {
	Opcode     Opcode
	Session_id int64
}

type Send_npc__packet struct {
	Opcode     Opcode
	Session_id int64
}

type Move_player_packet struct {
	Opcode     Opcode
	Session_id int64
	X          float32
	Y          float32
	Z          float32
	Direction  float32
}

type Spend_dp_packet struct {
	Opcode     Opcode
	Session_id int64
}

type Talk_cmd_packet struct {
	Opcode     Opcode
	Session_id int64
}

type Attack_cmd_packet struct {
	Opcode     Opcode
	Session_id int64
}

type Trade_cmd_packet struct {
	Opcode     Opcode
	Session_id int64
}

type Invite_cmd_packet struct {
	Opcode     Opcode
	Session_id int64
}

type Guild_invite_cmd_packet struct {
	Opcode     Opcode
	Session_id int64
}

type Guild_kick_cmd_packet struct {
	Opcode     Opcode
	Session_id int64
}

type Guild_promote_cmd_packet struct {
	Opcode     Opcode
	Session_id int64
}

type Guild_demote_cmd_packet struct {
	Opcode     Opcode
	Session_id int64
}

type Say_cmd_packet struct {
	Opcode     Opcode
	Session_id int64
	Text       string
}

type Yell_cmd_packet struct {
	Opcode     Opcode
	Session_id int64
	Text       string
}

type Ooc_cmd_packet struct {
	Opcode     Opcode
	Session_id int64
	Text       string
}

type Help_cmd_packet struct {
	Opcode     Opcode
	Session_id int64
	Text       string
}

type Pchat_cmd_packet struct {
	Opcode     Opcode
	Session_id int64
	Text       string
}

type Gchat_cmd_packet struct {
	Opcode     Opcode
	Session_id int64
	Text       string
}

type Whisper_cmd_packet struct {
	Opcode     Opcode
	Session_id int64
	Target     string
	Text       string
}

type Relay_say_packet struct {
	Opcode     Opcode
	Session_id int64
}

type Relay_yell_packet struct {
	Opcode     Opcode
	Session_id int64
}

type Relay_ooc_packet struct {
	Opcode     Opcode
	Session_id int64
}

type Relay_help_packet struct {
	Opcode     Opcode
	Session_id int64
}

type Relay_party_chat_packet struct {
	Opcode     Opcode
	Session_id int64
}

type Relay_guild_chat_packet struct {
	Opcode     Opcode
	Session_id int64
}

type Relay_whisper_packet struct {
	Opcode     Opcode
	Session_id int64
}

type Activate_cmd_packet struct {
	Opcode     Opcode
	Session_id int64
}

type Environment_update_packet struct {
	Opcode     Opcode
	Session_id int64
}

type Dialogue_text_packet struct {
	Opcode     Opcode
	Session_id int64
}

type Dialogue_cmd_packet struct {
	Opcode     Opcode
	Session_id int64
}

type Send_item_packet struct {
	Opcode     Opcode
	Session_id int64
}

type Init_shop_packet struct {
	Opcode     Opcode
	Session_id int64
}

type Shop_item_packet struct {
	Opcode     Opcode
	Session_id int64
}

type Buy_item_packet struct {
	Opcode     Opcode
	Session_id int64
}

type Init_loot_packet struct {
	Opcode     Opcode
	Session_id int64
}

type Loot_item_packet struct {
	Opcode     Opcode
	Session_id int64
}

type Take_loot_packet struct {
	Opcode     Opcode
	Session_id int64
}

type Init_trade_packet struct {
	Opcode     Opcode
	Session_id int64
}

type Offer_item_packet struct {
	Opcode     Opcode
	Session_id int64
}

type Pull_item_packet struct {
	Opcode     Opcode
	Session_id int64
}

type Trade_item_packet struct {
	Opcode     Opcode
	Session_id int64
}

type Remove_item_packet struct {
	Opcode     Opcode
	Session_id int64
}

type Accept_trade_packet struct {
	Opcode     Opcode
	Session_id int64
}

type Unaccept_trade_packet struct {
	Opcode     Opcode
	Session_id int64
}

type Communicate_trade_packet struct {
	Opcode     Opcode
	Session_id int64
}

type Finalize_trade_packet struct {
	Opcode     Opcode
	Session_id int64
}

type Init_combat_packet struct {
	Opcode     Opcode
	Session_id int64
}

type Add_npc_combat_packet struct {
	Opcode     Opcode
	Session_id int64
}

type Add_pc_combat_packet struct {
	Opcode     Opcode
	Session_id int64
}

type Remove_npc_combat_packet struct {
	Opcode     Opcode
	Session_id int64
}

type Remove_pc_combat_packet struct {
	Opcode     Opcode
	Session_id int64
}
