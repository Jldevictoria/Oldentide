// Filename:    packet.go
// Author:      Joseph DeVictoria
// Date:        June_16_2018
// Purpose:     Definition of packets to be used in Oldentide.

package shared

import (
	"net"
)

// Opcode is an eight bit unsigned integer that represents 2^8 different Packet types.
type Opcode uint8

// Packet type opcode, used to deserialize which type of packet is being sent or received from the network.
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

// RawPacket is a representation of the simplest UDP packet. The Payload field is what will be deserialized as a given Packet in the rest of packet.go.
type RawPacket struct {
	Size    int
	Client  *net.UDPAddr
	Payload []byte
}

// OpcodePacket represents a packet containing nothing but a simple opcode (ANY).
type OpcodePacket struct {
	Opcode Opcode
}

// EmptyPacket is a simple packet containing nothing by an opcode (EMPTY). BOTH - Ignore these.
type EmptyPacket struct {
	Opcode Opcode
}

// GenericPacket is a simple packet containing only an opcode (GENERIC) and a SessionID. BOTH - Used for debug.
type GenericPacket struct {
	Opcode    Opcode
	SessionID uint64
}

// AckPacket is the packet with Opcode (ACK) used for acknowledging reception of a packet.
type AckPacket struct {
	Opcode    Opcode
	SessionID uint64
}

// ErrorPacket is the packet with Opcode (ERROR) used to transmit knowledge of an error.
type ErrorPacket struct {
	Opcode    Opcode
	SessionID uint64
}

// ReqClistPacket is the packet with Opcode (REQCLIST)
type ReqClistPacket struct {
	Opcode    Opcode
	SessionID uint64
	Account   string
}

// SendClistPacket  is the packet with Opcode (SENDCLIST)
type SendClistPacket struct {
	Opcode     Opcode
	Characters []string
}

// CreatePlayerPacket is the packet with Opcode (CREATEPLAYER)
type CreatePlayerPacket struct {
	Opcode    Opcode
	SessionID uint64
	Pc        Pc
}

// ConnectPacket is the packet with Opcode (CONNECT)
type ConnectPacket struct {
	Opcode    Opcode
	SessionID uint64
	Firstname string
}

// DisconnectPacket is the packet with Opcode (DISCONNECT)
type DisconnectPacket struct {
	Opcode    Opcode
	SessionID uint64
}

// SendPlayerPacket is the packet with Opcode (SENDPLAYER)
type SendPlayerPacket struct {
	Opcode    Opcode
	SessionID uint64
}

// SendPcPacket is the packet with Opcode (SENDPC)
type SendPcPacket struct {
	Opcode    Opcode
	SessionID uint64
}

// SendNpcPacket is the packet with Opcode (SENDNPC)
type SendNpcPacket struct {
	Opcode    Opcode
	SessionID uint64
}

// MovePlayerPacket is the packet with Opcode (MOVEPLAYER)
type MovePlayerPacket struct {
	Opcode    Opcode
	SessionID uint64
	X         float32
	Y         float32
	Z         float32
	Direction float32
}

// SpendDpPacket is the packet with Opcode (SPENDDP)
type SpendDpPacket struct {
	Opcode    Opcode
	SessionID uint64
}

// TalkCmdPacket is the packet with Opcode (TALKCMD)
type TalkCmdPacket struct {
	Opcode    Opcode
	SessionID uint64
}

// AttackCmdPacket is the packet with Opcode (ATTACKCMD)
type AttackCmdPacket struct {
	Opcode    Opcode
	SessionID uint64
}

// TradeCmdPacket is the packet with Opcode (TRADECMD)
type TradeCmdPacket struct {
	Opcode    Opcode
	SessionID uint64
}

// InviteCmdPacket is the packet with Opcode (INVITECMD)
type InviteCmdPacket struct {
	Opcode    Opcode
	SessionID uint64
}

// GuildInviteCmdPacket is the packet with Opcode (GINVITECMD)
type GuildInviteCmdPacket struct {
	Opcode    Opcode
	SessionID uint64
}

// GuildKickCmdPacket is the packet with Opcode (GKICK)
type GuildKickCmdPacket struct {
	Opcode    Opcode
	SessionID uint64
}

// GuildPromoteCmdPacket is the packet with Opcode (GPROMOTE)
type GuildPromoteCmdPacket struct {
	Opcode    Opcode
	SessionID uint64
}

// GuildDemoteCmdPacket is the packet with Opcode (GDEMOTE)
type GuildDemoteCmdPacket struct {
	Opcode    Opcode
	SessionID uint64
}

// SayCmdPacket is the packet with Opcode (SAYCMD)
type SayCmdPacket struct {
	Opcode    Opcode
	SessionID uint64
	Message   string
}

// YellCmdPacket is the packet with Opcode (YELLCMD)
type YellCmdPacket struct {
	Opcode    Opcode
	SessionID uint64
	Message   string
}

// OocCmdPacket is the packet with Opcode (OOCCMD)
type OocCmdPacket struct {
	Opcode    Opcode
	SessionID uint64
	Message   string
}

// HelpCmdPacket is the packet with Opcode (HELPCMD)
type HelpCmdPacket struct {
	Opcode    Opcode
	SessionID uint64
	Message   string
}

// PchatCmdPacket is the packet with Opcode (PCHATCMD)
type PchatCmdPacket struct {
	Opcode    Opcode
	SessionID uint64
	Message   string
}

// GchatCmdPacket is the packet with Opcode (GCHATCMD)
type GchatCmdPacket struct {
	Opcode    Opcode
	SessionID uint64
	Message   string
}

// WhisperCmdPacket is the packet with Opcode (WHISPERCMD)
type WhisperCmdPacket struct {
	Opcode    Opcode
	SessionID uint64
	Target    string
	Message   string
}

// RelaySayPacket is the packet with Opcode (RELAYSAY)
type RelaySayPacket struct {
	Opcode    Opcode
	SessionID uint64
}

// RelayYellPacket is the packet with Opcode (RELAYYELL)
type RelayYellPacket struct {
	Opcode    Opcode
	SessionID uint64
}

// RelayOocPacket is the packet with Opcode (RELAYOOC)
type RelayOocPacket struct {
	Opcode    Opcode
	SessionID uint64
}

// RelayHelpPacket is the packet with Opcode (RELAYHELP)
type RelayHelpPacket struct {
	Opcode    Opcode
	SessionID uint64
}

// RelayPartyChatPacket is the packet with Opcode (RELAYPCHAT)
type RelayPartyChatPacket struct {
	Opcode    Opcode
	SessionID uint64
}

// RelayGuildChatPacket is the packet with Opcode (RELAYGCHAT)
type RelayGuildChatPacket struct {
	Opcode    Opcode
	SessionID uint64
}

// RelayWhisperPacket is the packet with Opcode (RELAYWHISPER)
type RelayWhisperPacket struct {
	Opcode    Opcode
	SessionID uint64
	Message   string
}

// ActivateCmdPacket is the packet with Opcode (ACTIVATECMD)
type ActivateCmdPacket struct {
	Opcode    Opcode
	SessionID uint64
}

// EnvironmentUpdatePacket is the packet with Opcode (ENVUPDATE)
type EnvironmentUpdatePacket struct {
	Opcode    Opcode
	SessionID uint64
}

// DialogueTextPacket is the packet with Opcode (DIALOGUETEXT)
type DialogueTextPacket struct {
	Opcode    Opcode
	SessionID uint64
}

// DialogueCmdPacket is the packet with Opcode (DIALOGUECMD)
type DialogueCmdPacket struct {
	Opcode    Opcode
	SessionID uint64
}

// SendItemPacket is the packet with Opcode (SENDITEM)
type SendItemPacket struct {
	Opcode    Opcode
	SessionID uint64
}

// InitShopPacket is the packet with Opcode (INITSHOP)
type InitShopPacket struct {
	Opcode    Opcode
	SessionID uint64
}

// ShopItemPacket is the packet with Opcode (SHOPITEM)
type ShopItemPacket struct {
	Opcode    Opcode
	SessionID uint64
}

// BuyItemPacket is the packet with Opcode (BUYITEM)
type BuyItemPacket struct {
	Opcode    Opcode
	SessionID uint64
}

// InitLootPacket is the packet with Opcode (INITLOOT)
type InitLootPacket struct {
	Opcode    Opcode
	SessionID uint64
}

// LootItemPacket is the packet with Opcode (LOOTITEM)
type LootItemPacket struct {
	Opcode    Opcode
	SessionID uint64
}

// TakeLootPacket is the packet with Opcode (TAKELOOT)
type TakeLootPacket struct {
	Opcode    Opcode
	SessionID uint64
}

// InitTradePacket is the packet with Opcode (INITTRADE)
type InitTradePacket struct {
	Opcode    Opcode
	SessionID uint64
}

// OfferItemPacket is the packet with Opcode (OFFERITEM)
type OfferItemPacket struct {
	Opcode    Opcode
	SessionID uint64
}

// PullItemPacket is the packet with Opcode (PULLITEM)
type PullItemPacket struct {
	Opcode    Opcode
	SessionID uint64
}

// TradeItemPacket is the packet with Opcode (TRADEITEM)
type TradeItemPacket struct {
	Opcode    Opcode
	SessionID uint64
}

// RemoveItemPacket is the packet with Opcode (REMITEM)
type RemoveItemPacket struct {
	Opcode    Opcode
	SessionID uint64
}

// AcceptTradePacket is the packet with Opcode (ACCTRADE)
type AcceptTradePacket struct {
	Opcode    Opcode
	SessionID uint64
}

// UnacceptTradePacket is the packet with Opcode (UNACCTRADE)
type UnacceptTradePacket struct {
	Opcode    Opcode
	SessionID uint64
}

// CommunicateTradePacket is the packet with Opcode (COMMTRADE)
type CommunicateTradePacket struct {
	Opcode    Opcode
	SessionID uint64
}

// FinalizeTradePacket is the packet with Opcode (FINTRADE)
type FinalizeTradePacket struct {
	Opcode    Opcode
	SessionID uint64
}

// InitCombatPacket is the packet with Opcode (INITCOMBAT)
type InitCombatPacket struct {
	Opcode    Opcode
	SessionID uint64
}

// AddNpcCombatPacket is the packet with Opcode (ADDNPCCOMBAT)
type AddNpcCombatPacket struct {
	Opcode    Opcode
	SessionID uint64
}

// AddPcCombatPacket is the packet with Opcode (ADDPCCOMBAT)
type AddPcCombatPacket struct {
	Opcode    Opcode
	SessionID uint64
}

// RemoveNpcCombatPacket is the packet with Opcode (REMNPCCOMBAT)
type RemoveNpcCombatPacket struct {
	Opcode    Opcode
	SessionID uint64
}

// RemovePcCombatPacket is the packet with Opcode (REMPCCOMBAT)
type RemovePcCombatPacket struct {
	Opcode    Opcode
	SessionID uint64
}
