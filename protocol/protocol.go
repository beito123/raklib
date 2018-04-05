package protocol

/*
	Raklib

	Copyright (c) 2018 beito

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.
*/

import "github.com/beito123/raklib"

type Protocol struct {
	packets []raklib.Packet
}

func (pro *Protocol) registerPackets() {
	pro.packets = make([]raklib.Packet, 0xff)

	pro.packets[IDPingDataPacket] = &PingDataPacket{}
	pro.packets[IDUnconnectedPing] = &UnconnectedPingPacket{}
	pro.packets[IDUnconnectedPingOpenConnections] = &UnconnectedPingOpenConnections{}
	pro.packets[IDPongDataPacket] = &PongDataPacket{}
	pro.packets[IDOpenConnectionRequest1] = &OpenConnectionRequest1Packet{}
	pro.packets[IDOpenConnectionReply1] = &OpenConnectionReply1Packet{}
	pro.packets[IDOpenConnectionRequest2] = &OpenConnectionRequest2Packet{}
	pro.packets[IDOpenConnectionReply2] = &OpenConnectionReply2Packet{}
	pro.packets[IDClientConnectDataPacket] = &ClientConnectDataPacket{}
	pro.packets[IDServerHandshakeDataPacket] = &ServerHandshakeDataPacket{}
	pro.packets[IDClientHandshakeDataPacket] = &ClientHandshakeDataPacket{}
	pro.packets[IDClientDisconnectDataPacket] = &ClientDisconnectDataPacket{}
	pro.packets[IDDataPacket0] = &DataPacket0{}
	pro.packets[IDDataPacket1] = &DataPacket1{}
	pro.packets[IDDataPacket2] = &DataPacket2{}
	pro.packets[IDDataPacket3] = &DataPacket3{}
	pro.packets[IDDataPacket4] = &DataPacket4{}
	pro.packets[IDDataPacket5] = &DataPacket5{}
	pro.packets[IDDataPacket6] = &DataPacket6{}
	pro.packets[IDDataPacket7] = &DataPacket7{}
	pro.packets[IDDataPacket8] = &DataPacket8{}
	pro.packets[IDDataPacket9] = &DataPacket9{}
	pro.packets[IDDataPacketA] = &DataPacketA{}
	pro.packets[IDDataPacketB] = &DataPacketB{}
	pro.packets[IDDataPacketC] = &DataPacketC{}
	pro.packets[IDDataPacketD] = &DataPacketD{}
	pro.packets[IDDataPacketE] = &DataPacketE{}
	pro.packets[IDDataPacketF] = &DataPacketF{}
	pro.packets[IDUnconnectedPong] = &UnconnectedPongPacket{}
}

func (pro *Protocol) Packet(id byte) raklib.Packet {
	pk := pro.packets[id]
	if pk == nil {
		return nil
	}

	return pk.New()
}
