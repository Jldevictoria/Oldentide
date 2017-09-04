using System;
using System.Collections.Generic;
using MessagePack;

namespace Oldentide.Networking {

	public enum PTYPE {
	    // Reserve 0 for null/uninitialized
		GENERIC = 1,
		ACK = 2,
		CONNECT = 3,
		DISCONNECT = 4,
	    ERROR = 5,
	    // CREATEACCOUNT = 6,
	    // LOGIN = 7,
	    LISTCHARACTERS = 8,
	    SELECTCHARACTER = 9,
	    DELETECHARACTER = 10,
	    CREATECHARACTER = 11,
	    INITIALIZEGAME = 12,
	    UPDATEPC = 13,
	    UPDATENPC = 14,
	    SENDPLAYERCOMMAND = 15,
	    SENDSERVERCOMMAND = 16,
	    SENDPLAYERACTION = 17,
	    SENDSERVERACTION = 18,
	    UNITY = 19,
	};

	// Create a single class for a packet
	public class Packet_t {
		public byte [] data;
		public PTYPE type;
	}




	[MessagePackObject]
	public class PacketAck {
	    [Key(0)]
	    public int packetId;
	    [Key(1)]
	    public string sessionId;
	};

	[MessagePackObject]
	public class PacketConnect {
	    [Key(0)]
	    public int packetId;
	    [Key(1)]
	    public string sessionId;
	};

	[MessagePackObject]
	public class PacketDisconnect {
	    [Key(0)]
	    public int packetId;
	    [Key(1)]
	    public string sessionId;
	};

	[MessagePackObject]
	public class PacketListcharacters {
	    [Key(0)]
	    public int packetId;
	    [Key(1)]
	    public string sessionId;
	    [Key(2)]
	    public string [] characterArray;
	};

	[MessagePackObject]
	public class PacketError {
	    [Key(0)]
	    public int packetId;
	    [Key(1)]
	    public string sessionId;
	    [Key(2)]
	    public string errorMsg;
	    // TODO: Implement error number system
	    // [Key(3)]
	    // public int errorNumber;
	};

    [MessagePackObject]
	public class PacketSelectcharacter {
	    [Key(0)]
	    public int packetId;
	    [Key(1)]
	    public string sessionId;
	    [Key(2)]
	    public string character;
	};

    [MessagePackObject]
	public class PacketDeletecharacter {
	    [Key(0)]
	    public int packetId;
	    [Key(1)]
	    public string sessionId;
	    [Key(2)]
	    public string character;
	};

    [MessagePackObject]
	public class PacketCreatecharacter {
	    [Key(0)]
	    public int packetId;
	    [Key(1)]
	    public string sessionId;
	    [Key(2)]
	    public string firstName;
	    [Key(3)]
	    public string lastName;
	    [Key(4)]
	    public string race;
	    [Key(5)]
	    public string gender;
	    [Key(6)]
	    public string profession;
	    // c++ map -> c # SortedDictionary
	    // c++ unordered_map -> c # Dictionary
	    // msgpack shouldn't care about the order though for maps
	    [Key(7)]
	    public Dictionary<string, int> attributes;
	    // TODO: Create an attributes map -> map<string, int> attributes
	};

    [MessagePackObject]
	public class PacketInitializegame {
        [Key(0)]
        public int packetId;
        [Key(1)]
        public string sessionId;
	};

    [MessagePackObject]
	public class PacketUpdatepc {
        [Key(0)]
        public int packetId;
        [Key(1)]
        public string sessionId;
        [Key(2)]
        public string firstName;
        [Key(3)]
        public string lastName;
        [Key(4)]
        public string race;
        [Key(5)]
        public string gender;
        [Key(6)]
        public string profession;
        [Key(7)]
        public int level;
        [Key(8)]
        public int hp;
        [Key(9)]
        public int bp;
        [Key(10)]
        public int mp;
        [Key(11)]
        public int ep;
        [Key(12)]
        public int x;
        [Key(13)]
        public int y;
        [Key(14)]
        public int z;
        [Key(15)]
	    public float direction;
	};

    [MessagePackObject]
	public class PacketUpdatenpc {
        [Key(0)]
        public int packetId;
        [Key(1)]
        public string sessionId;
	};

    [MessagePackObject]
	public class PacketSendplayercommand {
        [Key(0)]
        public int packetId;
        [Key(1)]
        public string sessionId;
        [Key(2)]
        public string command;
    };

    [MessagePackObject]
    public class PacketSendservercommand {
        [Key(0)]
        public int packetId;
        [Key(1)]
        public string sessionId;
        [Key(2)]
	    public string command;
	};

    [MessagePackObject]
	public class PacketSendplayeraction {
        [Key(0)]
        public int packetId;
        [Key(1)]
        public string sessionId;
	};

    [MessagePackObject]
	public class PacketSendserveraction {
        [Key(0)]
        public int packetId;
        [Key(1)]
        public string sessionId;
	};

	// mark MessagePackObjectAttribute
	[MessagePackObject]
	public class PacketUnity {
	    [Key(0)]
	    public int packetId;
	    [Key(1)]
	    public string sessionId;
	    [Key(2)]
		public int data1;
	    [Key(3)]
		public int data2;
	    [Key(4)]
		public int data3;
	    [Key(5)]
		public int data4;
	    [Key(6)]
		public int data5;
	};
}