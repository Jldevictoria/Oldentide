﻿using UnityEngine;
using Oldentide.Networking;
using System;
using System.IO;
using System.Net;
using System.Net.Sockets;
using System.Runtime.InteropServices;
using MessagePack;

public class NetworkInterface : MonoBehaviour {

	public string serverIp = "goblin.oldentide.com";
	public int serverPort = 1337;
	IPEndPoint serverEndPoint;
	IPEndPoint clientEndPoint;
	Socket clientSocket;


	int packetNumber = 1;
	int session = 0;

	// This is the number of bytes that the oldentide header is
	const int HEADER_SIZE = 3;
	const int PACKET_MAX_SIZE = 512;

	// Use this for initialization
	void Start() {
		// Set up Server End Point for sending packets.
		IPHostEntry serverHostEntry = Dns.GetHostEntry(serverIp);
		IPAddress serverIpAddress = serverHostEntry.AddressList[0];
		serverEndPoint = new IPEndPoint(serverIpAddress, serverPort);
		Debug.Log("Server IPEndPoint: " + serverEndPoint.ToString());
		// Set up Client End Point for receiving packets.
		IPHostEntry clientHostEntry = Dns.GetHostEntry(Dns.GetHostName());
		IPAddress clientIpAddress = IPAddress.Any;
		foreach (IPAddress ip in clientHostEntry.AddressList) {
			if (ip.AddressFamily == AddressFamily.InterNetwork) {
				clientIpAddress = ip;
			}
		}
		clientEndPoint = new IPEndPoint(clientIpAddress, serverPort);
		Debug.Log("Client IPEndPoint: " + clientEndPoint.ToString());
		// Create socket for client and bind to Client End Point (Ip/Port).
		clientSocket = new Socket(AddressFamily.InterNetwork, SocketType.Dgram, ProtocolType.Udp);
		try {
			clientSocket.Bind(clientEndPoint);
		}
		catch (Exception e) {
			Debug.Log("Winsock error: " + e.ToString());
		}
	}



	// For the MessagePack library, using:
	// https://github.com/neuecc/MessagePack-CSharp
	// instead of:
	// https://github.com/msgpack/msgpack-cli

	// To install, simply download the Unity zip from the releases page:
	// https://github.com/neuecc/MessagePack-CSharp/releases
	// Unzip it into a folder NOT under Oldentide/
	// (don't want unity automatically ingesting the files just yet)
	// In the Unity editor, click Assets -> Import Package
	// Select the MessagePack .unitypackage file and click import.
	// e.g. ...\MessagePack.Unity.1.2.0\MessagePack.Unity.1.2.0.unitypackage"
	// Now all the MessagePack files are loaded into the project!
	// To use it in code, put "using MessagePack" at the top of any C# files


	// Update is called once per frame
	void Update() {

		if(Input.GetKeyDown(KeyCode.M)){
			Debug.Log("Sending a single CONNECT packet via message pack!!");

			PACKET_CONNECT pp;
			pp.sessionId = session;
			pp.packetId = packetNumber;
			packetNumber++;

			byte [] msgpackData = MessagePackSerializer.Serialize(pp);
			SendDataTo(clientSocket, serverEndPoint, Oldentide.Networking.PTYPE.CONNECT, msgpackData);

			// Wait for the response
			byte[] receivedMsgpackData;
			Oldentide.Networking.PTYPE packetType = ReceiveDataFrom(out receivedMsgpackData);

			Debug.Log("Server responded with packet " + packetType);


			// TODO: What to do with this switch statement? Put it into a function of its own? Put it in a while loop? Use async stuff so it doesn't block the update function?

			switch(packetType){
				case Oldentide.Networking.PTYPE.CONNECT:
					ConnectHandler(receivedMsgpackData);
					break;
				case Oldentide.Networking.PTYPE.UNITY:
					UnityHandler(receivedMsgpackData);
					break;
				default:
					Debug.Log("Unknown packet received from the server!!! " + packetType);
					break;
			}

	   }

	}


	////
	/// Packet handlers
	//

	void ConnectHandler(byte [] msgpackData) {
		var data = MessagePackSerializer.Deserialize<PACKET_CONNECT>(msgpackData);
		Debug.Log("Connect handler! sessionId: " + data.sessionId + " ; packetId: " + data.packetId);

		if(data.sessionId != session){
			Debug.Log("Setting new session to " + data.sessionId);
			session = data.sessionId;
		}
		else {
			Debug.Log("Session is already set to " + session);
		}

		// Send a Unity packet, just for kicks
		PACKET_UNITY pp;
		// pp.packetType = Oldentide.Networking.PTYPE.UNITY;
		pp.sessionId = data.sessionId;
		pp.packetId = packetNumber;
		packetNumber++;
		pp.data1 = 1;
		pp.data2 = 8;
		pp.data3 = 16;
		pp.data4 = 64;
		pp.data5 = 255;

		byte [] msgpackDataToSendUnity = MessagePackSerializer.Serialize(pp);
		SendDataTo(clientSocket, serverEndPoint, Oldentide.Networking.PTYPE.UNITY, msgpackDataToSendUnity);

		// // Wait for a response
		// byte[] receivedMsgpackDataUnity;
		// Oldentide.Networking.PTYPE packetType = ReceiveDataFrom(out receivedMsgpackDataUnity);

		// if(packetType == Oldentide.Networking.PTYPE.UNITY){
		// 	UnityHandler(receivedMsgpackDataUnity);
		// }
		// else {
		// 	Debug.Log("Unknown packet received instead of unity packet...: " + packetType);
		// }


	}



	void UnityHandler(byte [] msgpackData) {
		var data = MessagePackSerializer.Deserialize<PACKET_UNITY>(msgpackData);
		Debug.Log("Unity handler! Data1: " + data.data1 + "; Data2: " + data.data2 + "; data3: " + data.data3 + ";");
	}


	////
	/// Util functions
	//

	void SendDataTo(Socket clientSocket, IPEndPoint serverEndPoint, Oldentide.Networking.PTYPE packetType, byte [] msgpackData){
		byte[] packetToSend = new byte[msgpackData.Length + HEADER_SIZE];

		// Prepend header data
		packetToSend[0] = (byte) packetType;
		// Convert msgpack length to a byte array
		byte[] msgpackLength = BitConverter.GetBytes(msgpackData.Length);
		// Copy the 2 bytes of the msgpack length into the packet at location 1
		Array.Copy(msgpackLength, 0, packetToSend, 1, 2);
		// Copy the msgpack data into the packet
		Array.Copy(msgpackData, 0, packetToSend, HEADER_SIZE, msgpackData.Length);

		Debug.Log("Sending packet data: ");
		PrintHexString(packetToSend);
		clientSocket.SendTo(packetToSend, serverEndPoint);
	}


	// TODO: Is there any way to receive from only the server, and not just anybody?
	Oldentide.Networking.PTYPE ReceiveDataFrom(out byte[] data){
		byte[] packetToReceive = new byte[PACKET_MAX_SIZE];
		IPEndPoint sender = new IPEndPoint(IPAddress.Any, 0);
		EndPoint senderRemote = (EndPoint)sender;
		clientSocket.ReceiveFrom(packetToReceive, ref senderRemote);


		// Get the packet type
		Oldentide.Networking.PTYPE packetType = (Oldentide.Networking.PTYPE) packetToReceive[0];
		// Get the length of the msgpack data
		ushort msgpackLength = BitConverter.ToUInt16(packetToReceive, 1);
		// Copy the msgpack data into the packet
		data = new byte[msgpackLength];
		Array.Copy(packetToReceive, HEADER_SIZE, data, 0, msgpackLength);

		Debug.Log("Receiving msgpack data: ");
		PrintHexString(data);

		return packetType;
	}

	public void PrintHexString(byte [] bytearray){
		string hexstring = "0x";
		int len = bytearray.Length;
		for (int i = 0; i < len; i++) {
			hexstring += String.Format("{0:X2}", bytearray[i]);
		}
		Debug.Log(hexstring);
	}

	// public byte [] StructureToByteArray(object obj){
	// 	int len = Marshal.SizeOf(obj);
	// 	byte [] arr = new byte[len];
	// 	IntPtr ptr = Marshal.AllocHGlobal(len);
	// 	Marshal.StructureToPtr(obj, ptr, true);
	// 	Marshal.Copy(ptr, arr, 0, len);
	// 	Marshal.FreeHGlobal(ptr);
	// 	return arr;
	// }


	// public void ByteArrayToStructure(byte [] bytearray, ref object obj){
	// 	int len = Marshal.SizeOf(obj);
	// 	IntPtr i = Marshal.AllocHGlobal(len);
	// 	Marshal.Copy(bytearray,0, i,len);
	// 	obj = Marshal.PtrToStructure(i, obj.GetType());
	// 	Marshal.FreeHGlobal(i);
	// }



}