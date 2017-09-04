using UnityEngine;
using Oldentide.Networking;
using System;
using System.IO;
using System.Collections;
using System.Collections.Generic;
using System.Net;
using System.Net.Sockets;
using System.Runtime.InteropServices;
using System.Threading;
using MessagePack;
using UnityEngine.UI;


public class NetworkInterface : MonoBehaviour {

    public InputField messageInput;
    public Text messages;

    ServerConfig serverConfig;

    IPEndPoint serverEndPoint;
    IPEndPoint clientEndPoint;
    Socket clientSocket;

    Queue<Packet_t> packetQueue = new Queue<Packet_t>();

    int packetNumber = 1;

    // Set this to false to stop the thread, or else it will keep running
    bool listenForPackets = true;
    bool waitingForPacket = false;

    // This is the number of bytes that the oldentide header is
    const int HEADER_SIZE = 3;
    const int PACKET_MAX_SIZE = 512;

    private void Awake() {
        serverConfig = gameObject.GetComponent<ServerConfig>();
    }

    // Use this for initialization
    void Start() {
        // Set up Server End Point for sending packets.
        IPHostEntry serverHostEntry = Dns.GetHostEntry(serverConfig.serverIp);
        IPAddress serverIpAddress = serverHostEntry.AddressList[0];
        serverEndPoint = new IPEndPoint(serverIpAddress, serverConfig.serverPort);
        Debug.Log("Server IPEndPoint: " + serverEndPoint.ToString());
        // Set up Client End Point for receiving packets.
        IPHostEntry clientHostEntry = Dns.GetHostEntry(Dns.GetHostName());
        IPAddress clientIpAddress = IPAddress.Any;
        foreach (IPAddress ip in clientHostEntry.AddressList) {
            if (ip.AddressFamily == AddressFamily.InterNetwork) {
                clientIpAddress = ip;
            }
        }
        clientEndPoint = new IPEndPoint(clientIpAddress, serverConfig.serverPort);
        Debug.Log("Client IPEndPoint: " + clientEndPoint.ToString());
        // Create socket for client and bind to Client End Point (Ip/Port).
        clientSocket = new Socket(AddressFamily.InterNetwork, SocketType.Dgram, ProtocolType.Udp);
        try {
            clientSocket.Bind(clientEndPoint);
        }
        catch (Exception e) {
            Debug.Log("Winsock error: " + e.ToString());
        }



        // TODO: Get this to work!!!
        // // Submit the user command
        // messageInput.OnSubmit.AddListener(delegate(BaseEventData eventData) {
        //     Debug.Log("Sending command...");
        //     BroadcastAction(messageInput.text);
        //     messageInput.text = "";
        // });
    }


    // Update is called once per frame
    void Update() {
        if(Input.GetKeyDown(KeyCode.Alpha7)){
            ConnectToServer();
        }
        if(Input.GetKeyDown(KeyCode.Alpha8)){
            ListCharactersAction();
        }
        if(Input.GetKeyDown(KeyCode.Return)){
            Debug.Log("Sending command...");
            // Get the text input and send it
            BroadcastAction(messageInput.text);
            // Clear the input
            messageInput.text = "";
        }

        // if(Input.GetKeyDown(KeyCode.Alpha6)){
            // messageInput.Select();
            // messageInput.ActivateInputField();
            // messageInput.text = "sdf";
            // Debug.Log("Selecting text input");
            // Debug.Log(messageInput.enabled);
        // }


        ////
        ///     Listen for any packets
        //

        // Continually wait for a packet to come in from the server
        // If we are already waiting for a packet, give control back to unity
        if(listenForPackets && !waitingForPacket){
            Debug.Log("Start waiting for a packet!");
            waitingForPacket = true;
            // Wait for the response asynchronously with a callback (this callback is executed in a different thread!)
            ReceiveDataFrom(delegate(Packet_t p) {
                // Now that we have a packet, copy it into the packet queue
                lock(packetQueue){
                    packetQueue.Enqueue(p);
                }
                waitingForPacket = false;
            }); // End ReceiveDataFrom
        } // End packetListener while loop


        ////
        ///     Service any packets in the packet queue
        //

        // Wait on packets in the packetQueue and service them
        // This thread shouldn't have to wait too long on this lock
        // It only needs to wait on the listener thread to enqueue an item
        // Don't take too long to service any packet though
        Packet_t returnPacket;
        lock(packetQueue){
            if(packetQueue.Count > 0){
                returnPacket = packetQueue.Dequeue();
            }
            else {
                // No packet, so leave Update() and check next time
                return;
            }
        }

        Debug.Log("Found a packet in the queue!");

        switch(returnPacket.type) {
            case Oldentide.Networking.PTYPE.ERROR: {
                var data = MessagePackSerializer.Deserialize<PacketError>(returnPacket.data);
                Debug.Log("ERROR: " + data.errorMsg);
                messages.text += "ERROR: " + data.errorMsg;
                break;
            }
            case Oldentide.Networking.PTYPE.SENDSERVERCOMMAND: {
                var data = MessagePackSerializer.Deserialize<PacketSendservercommand>(returnPacket.data);
                if(data.sessionId != OldentidePlayerInformation.sessionId){
                    Debug.Log("Invalid session! Ignoring SENDSERVERCOMMAND packet...");
                    return;
                }
                Debug.Log(data.command);
                messages.text += "\n" + data.command;
                break;
            }

            case Oldentide.Networking.PTYPE.CONNECT: {
                var data = MessagePackSerializer.Deserialize<PacketConnect>(returnPacket.data);
                // Debug.Log("Connect packet response! sessionId: " + data.sessionId + " ; packetId: " + data.packetId);
                if(data.sessionId != OldentidePlayerInformation.sessionId) {
                    Debug.Log("Setting new session to " + data.sessionId);
                    OldentidePlayerInformation.sessionId = data.sessionId;
                }
                else {
                    Debug.Log("Session is already set to " + OldentidePlayerInformation.sessionId);
                }
                break;
            }

            case Oldentide.Networking.PTYPE.SENDPLAYERCOMMAND: {
                var data = MessagePackSerializer.Deserialize<PacketSendplayercommand>(returnPacket.data);
                if(data.sessionId != OldentidePlayerInformation.sessionId) {
                    Debug.Log("Invalid session! Ignoring packet response");
                    return;
                }
                Debug.Log("SENDPLAYERCOMMAND response: " + data.command);
                break;
            }

            case Oldentide.Networking.PTYPE.LISTCHARACTERS: {
                var data = MessagePackSerializer.Deserialize<PacketListcharacters>(returnPacket.data);
                if(data.sessionId != OldentidePlayerInformation.sessionId) {
                    Debug.Log("Invalid session! Ignoring packet response");
                    return;
                }

                Debug.Log("ListCharacter response! sessionId: " + data.sessionId + " ; packetId: " + data.packetId);
                // Print out the character list
                foreach (string element in data.characterArray) {
                    Debug.Log(element);
                }
                break;
            }

            default: {
                Debug.Log("Unhandled packet type: " + returnPacket.type);
                break;
            }

        } // End packet switch


    } // End update

    void OnApplicationQuit() {
        Debug.Log("Application ending after " + Time.time + " seconds");
        listenForPackets = false;
        Debug.Log("Closing socket...");
        // Close the socket the is most likely blocking right now
        // This should stop it
        clientSocket.Close();
    }


    public void ConnectToServer(){
        Debug.Log("sending a single CONNECT packet via message pack!!");

        PacketConnect pp = new PacketConnect();
        pp.sessionId = OldentidePlayerInformation.sessionId;
        pp.packetId = packetNumber;
        packetNumber++;

        byte [] msgpackData = MessagePackSerializer.Serialize(pp);
        // Synchronously send, since it's UDP
        SendDataTo(clientSocket, serverEndPoint, Oldentide.Networking.PTYPE.CONNECT, msgpackData);
    }


    void ListCharactersAction(){
        PacketListcharacters pp = new PacketListcharacters();
        pp.sessionId = OldentidePlayerInformation.sessionId;
        pp.packetId = packetNumber;
        // Set to empty array, NOT null. This makes a difference in messagepack
        pp.characterArray = new string [0];
        packetNumber++;

        byte [] sendMsgpackData = MessagePackSerializer.Serialize(pp);
        SendDataTo(clientSocket, serverEndPoint, Oldentide.Networking.PTYPE.LISTCHARACTERS, sendMsgpackData);

    }

    void BroadcastAction(object command){
        Debug.Log("Sending a broadcast!");
        PacketSendplayercommand pp = new PacketSendplayercommand();
        pp.sessionId = OldentidePlayerInformation.sessionId;
        pp.packetId = packetNumber;
        pp.command = (string)command;
        packetNumber++;

        byte [] sendMsgpackData = MessagePackSerializer.Serialize(pp);
        SendDataTo(clientSocket, serverEndPoint, Oldentide.Networking.PTYPE.SENDPLAYERCOMMAND, sendMsgpackData);
    }


    ////
    /// Util functions
    //

    // Everything below here is analogous to the Utils C++ class in Server - wrappers to
    // the (UDP) socket functions that will be used multiple times

    // Gameplan: Use a delegate function, so the user can pass in their own callback function
    // That way, we can hide the complexities of the other stuff

    // How to use delegates and async callbacks:
    // http://stackoverflow.com/questions/2082615/pass-method-as-parameter-using-c-sharp
    // http://stackoverflow.com/questions/6866347/lambda-anonymous-function-as-a-parameter
    // https://docs.microsoft.com/en-us/dotnet/articles/csharp/programming-guide/statements-expressions-operators/anonymous-functions

    // Create a class to pass data to ReceiveDataFromCallback
    public class CallbackParams{
        public byte[] buffer;
        public Action<Packet_t> callback;
        // Note: Action<...> is a no-return delegate; Func<..., retType> has a return type
        // public Func<string, int> callback;
    }

    // Wrapper function for BeginReceiveFrom() that deals with
    // The callback params include the msgpack data and the packetType
    void ReceiveDataFrom(Action<Packet_t> callback) {
        // byte[] packetToReceive = new byte[PACKET_MAX_SIZE];
        IPEndPoint sender = new IPEndPoint(IPAddress.Any, 0);
        EndPoint senderRemote = (EndPoint)sender;

        // Use BeginReceiveFrom and EndReceiveFrom:
        // https://msdn.microsoft.com/en-us/library/system.net.sockets.socket.beginreceivefrom(v=vs.110).aspx
        // http://stackoverflow.com/questions/7860398/c-sharp-stateobject-class-location

        CallbackParams so2 = new CallbackParams();
        so2.buffer = new byte[PACKET_MAX_SIZE];
        so2.callback = callback;
        clientSocket.BeginReceiveFrom(so2.buffer, 0, PACKET_MAX_SIZE, 0, ref senderRemote, new AsyncCallback(ReceiveDataFromCallback), so2);
    }

    void ReceiveDataFromCallback(IAsyncResult result){
        // Get the params passed to this callback
        CallbackParams so2 = (CallbackParams) result.AsyncState;

        // // ???? I'm not quite sure why EndReceiveFrom() is even necessary
        // // It seems to work without it just fine
        // // see https://msdn.microsoft.com/en-us/library/system.net.sockets.socket.endreceivefrom(v=vs.110).aspx
        // IPEndPoint sender = new IPEndPoint(IPAddress.Any, 0);
        // EndPoint senderRemote = (EndPoint)sender;
        // // Block and wait for a packet to be put in so2.buffer
        // int read = clientSocket.EndReceiveFrom(result, ref senderRemote);
        // Debug.Log("Got packet, read " + read + " bytes");

        // Get the packet type
        Oldentide.Networking.PTYPE packetType = (Oldentide.Networking.PTYPE) so2.buffer[0];
        // Get the length of the msgpack data
        ushort msgpackLength = BitConverter.ToUInt16(so2.buffer, 1);
        // Copy the msgpack data into the packet
        byte[] data = new byte[msgpackLength];
        Array.Copy(so2.buffer, HEADER_SIZE, data, 0, msgpackLength);

        Debug.Log("Receiving msgpack data: ");
        PrintHexString(data);

        Packet_t returnPacket = new Packet_t(){
            data = data,
            type = packetType
        };
        // Return execution back to the original caller, passing back the data
        so2.callback(returnPacket);
    }

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

    public void PrintHexString(byte [] bytearray){
        string hexstring = "0x";
        int len = bytearray.Length;
        for (int i = 0; i < len; i++) {
            hexstring += String.Format("{0:X2}", bytearray[i]);
        }
        Debug.Log(hexstring);
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

}