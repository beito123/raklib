package protocol

/*
	Raklib

	Copyright (c) 2018 beito

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.
*/

import (
	"bytes"
	"github.com/beito123/raklib"
	"github.com/beito123/raklib/binary"
)

type BasePacket struct {
	binary.RaknetStream
}

func (base *BasePacket) Encode(pk raklib.Packet) error {
	if base.Buffer == nil {
		base.Buffer = &bytes.Buffer{}
	}

	err := base.PutByte(pk.ID())
	if err != nil {
		return err
	}

	return nil
}

func (base *BasePacket) Decode(pk raklib.Packet) error {
	if base.Buffer == nil {
		return raklib.NoSetBufferError{}
	}

	base.Skip(1) // id

	return nil
}

type PingDataPacket struct {
	BasePacket

	Time int64
}

func (PingDataPacket) ID() byte {
	return IDPingDataPacket
}

func (PingDataPacket) New() raklib.Packet {
	return new(PingDataPacket)
}

func (pk *PingDataPacket) Encode() error {
	err := pk.BasePacket.Encode(pk)

	err = pk.PutLong(pk.Time)
	if err != nil {
		return err
	}

	return nil
}

func (pk *PingDataPacket) Decode() error {
	err := pk.BasePacket.Decode(pk)
	if err != nil {
		return err
	}

	err = pk.Long(&pk.Time)
	if err != nil {
		return err
	}

	return nil
}

type UnconnectedPingPacket struct {
	BasePacket

	Time  int64
	Magic string
}

func (UnconnectedPingPacket) ID() byte {
	return IDUnconnectedPing
}

func (UnconnectedPingPacket) New() raklib.Packet {
	return new(UnconnectedPingPacket)
}

func (pk *UnconnectedPingPacket) Encode() error {
	err := pk.BasePacket.Encode(pk)
	if err != nil {
		return err
	}

	err = pk.Long(&pk.Time)
	if err != nil {
		return err
	}

	pk.HexString(16, &pk.Magic)

	return nil
}

func (pk *UnconnectedPingPacket) Decode() error {
	err := pk.BasePacket.Decode(pk)
	if err != nil {
		return err
	}

	err = pk.PutLong(pk.Time)
	if err != nil {
		return err
	}

	err = pk.PutHexString(raklib.Magic)
	if err != nil {
		return err
	}

	return nil
}

type UnconnectedPingOpenConnections struct {
	BasePacket
}

func (UnconnectedPingOpenConnections) ID() byte {
	return IDUnconnectedPingOpenConnections
}

func (UnconnectedPingOpenConnections) New() raklib.Packet {
	return new(UnconnectedPingOpenConnections)
}

func (pk *UnconnectedPingOpenConnections) Encode() error {
	err := pk.BasePacket.Encode(pk)
	if err != nil {
		return err
	}

	return nil
}

func (pk *UnconnectedPingOpenConnections) Decode() error {
	err := pk.BasePacket.Decode(pk)
	if err != nil {
		return err
	}

	return err
}

type PongDataPacket struct {
	BasePacket
}

func (PongDataPacket) ID() byte {
	return IDPongDataPacket
}

func (PongDataPacket) New() raklib.Packet {
	return new(PongDataPacket)
}

func (pk *PongDataPacket) Encode() error {
	err := pk.BasePacket.Encode(pk)
	if err != nil {
		return err
	}

	return nil
}

func (pk *PongDataPacket) Decode() error {
	err := pk.BasePacket.Decode(pk)
	if err != nil {
		return err
	}

	return err
}

type OpenConnectionRequest1Packet struct {
	BasePacket

	Magic    string
	Protocol byte
	MTU      []byte
}

func (OpenConnectionRequest1Packet) ID() byte {
	return IDOpenConnectionRequest1
}

func (OpenConnectionRequest1Packet) New() raklib.Packet {
	return new(OpenConnectionRequest1Packet)
}

func (pk *OpenConnectionRequest1Packet) Encode() error {
	err := pk.BasePacket.Encode(pk)
	if err != nil {
		return err
	}

	err = pk.PutHexString(raklib.Magic)
	if err != nil {
		return err
	}

	err = pk.PutByte(pk.Protocol)
	if err != nil {
		return err
	}

	m := make([]byte, 46) //

	err = pk.Put(m)
	if err != nil {
		return err
	}

	return nil
}

func (pk *OpenConnectionRequest1Packet) Decode() error {
	err := pk.BasePacket.Decode(pk)
	if err != nil {
		return err
	}

	pk.HexString(16, &pk.Magic)

	err = pk.Byte(&pk.Protocol)
	if err != nil {
		return err
	}

	pk.MTU = pk.Get(-1)

	return nil
}

type OpenConnectionReply1Packet struct {
	BasePacket

	Magic      string
	ServerUUID int64
	Security   bool
	MTU        uint16
}

func (OpenConnectionReply1Packet) ID() byte {
	return IDOpenConnectionReply1
}

func (OpenConnectionReply1Packet) New() raklib.Packet {
	return new(OpenConnectionReply1Packet)
}

func (pk *OpenConnectionReply1Packet) Encode() error {
	err := pk.BasePacket.Encode(pk)
	if err != nil {
		return err
	}

	err = pk.PutHexString(raklib.Magic)
	if err != nil {
		return err
	}

	err = pk.PutLong(pk.ServerUUID)
	if err != nil {
		return err
	}

	err = pk.PutBool(pk.Security)
	if err != nil {
		return err
	}

	err = pk.PutShort(pk.MTU)
	if err != nil {
		return err
	}

	return nil
}

func (pk *OpenConnectionReply1Packet) Decode() error {
	err := pk.BasePacket.Decode(pk)
	if err != nil {
		return err
	}

	pk.HexString(16, &pk.Magic)

	err = pk.Long(&pk.ServerUUID)
	if err != nil {
		return err
	}

	err = pk.Bool(&pk.Security)
	if err != nil {
		return err
	}

	err = pk.Short(&pk.MTU)
	if err != nil {
		return err
	}

	return err
}

type OpenConnectionRequest2Packet struct {
	BasePacket

	Magic         string
	ServerAddress raklib.SystemAddress
	MTU           uint16
	ClientUUID    int64
}

func (OpenConnectionRequest2Packet) ID() byte {
	return IDOpenConnectionRequest2
}

func (OpenConnectionRequest2Packet) New() raklib.Packet {
	return new(OpenConnectionRequest2Packet)
}

func (pk *OpenConnectionRequest2Packet) Encode() error {
	err := pk.BasePacket.Encode(pk)
	if err != nil {
		return err
	}

	err = pk.PutHexString(raklib.Magic)
	if err != nil {
		return err
	}

	err = pk.PutAddressSystemAddress(pk.ServerAddress)
	if err != nil {
		return err
	}

	err = pk.PutShort(pk.MTU)
	if err != nil {
		return err
	}

	err = pk.PutLong(pk.ClientUUID)
	if err != nil {
		return err
	}

	return nil
}

func (pk *OpenConnectionRequest2Packet) Decode() error {
	err := pk.BasePacket.Decode(pk)
	if err != nil {
		return err
	}

	pk.HexString(16, &pk.Magic)

	err = pk.AddressSystemAddress(&pk.ServerAddress)
	if err != nil {
		return err
	}

	err = pk.Short(&pk.MTU)
	if err != nil {
		return err
	}

	err = pk.Long(&pk.ClientUUID)
	if err != nil {
		return err
	}

	return err
}

type OpenConnectionReply2Packet struct {
	BasePacket

	Magic         string
	ServerUUID    int64
	ClientAddress raklib.SystemAddress
	MTU           uint16
	Encryption    byte
}

func (OpenConnectionReply2Packet) ID() byte {
	return IDOpenConnectionReply2
}

func (OpenConnectionReply2Packet) New() raklib.Packet {
	return new(OpenConnectionReply2Packet)
}

func (pk *OpenConnectionReply2Packet) Encode() error {
	err := pk.BasePacket.Encode(pk)
	if err != nil {
		return err
	}

	err = pk.PutHexString(raklib.Magic)
	if err != nil {
		return err
	}

	err = pk.PutLong(pk.ServerUUID)
	if err != nil {
		return err
	}

	err = pk.PutAddressSystemAddress(pk.ClientAddress)
	if err != nil {
		return err
	}

	err = pk.PutByte(pk.Encryption)
	if err != nil {
		return err
	}

	return nil
}

func (pk *OpenConnectionReply2Packet) Decode() error {
	err := pk.BasePacket.Decode(pk)
	if err != nil {
		return err
	}

	pk.HexString(16, &pk.Magic)

	err = pk.Long(&pk.ServerUUID)
	if err != nil {
		return err
	}

	err = pk.AddressSystemAddress(&pk.ClientAddress)
	if err != nil {
		return err
	}

	err = pk.Short(&pk.MTU)
	if err != nil {
		return err
	}

	return err
}

type ClientConnectDataPacket struct {
	BasePacket

	UUID int64
	Time int64
}

func (ClientConnectDataPacket) ID() byte {
	return IDClientConnectDataPacket
}

func (ClientConnectDataPacket) New() raklib.Packet {
	return new(ClientConnectDataPacket)
}

func (pk *ClientConnectDataPacket) Encode() error {
	err := pk.BasePacket.Encode(pk)
	if err != nil {
		return err
	}

	err = pk.PutLong(pk.UUID)
	if err != nil {
		return err
	}

	err = pk.PutLong(pk.Time)
	if err != nil {
		return err
	}

	return nil
}

func (pk *ClientConnectDataPacket) Decode() error {
	err := pk.BasePacket.Decode(pk)
	if err != nil {
		return err
	}

	err = pk.Long(&pk.UUID)
	if err != nil {
		return err
	}

	err = pk.Long(&pk.Time)
	if err != nil {
		return err
	}

	return err
}

type ServerHandshakeDataPacket struct {
	BasePacket

	ClientAddr      raklib.SystemAddress
	SystemIndex     uint16
	SystemAddresses [10]raklib.SystemAddress
	RequestTime     int64
	Time            int64
}

func (ServerHandshakeDataPacket) ID() byte {
	return IDServerHandshakeDataPacket
}

func (ServerHandshakeDataPacket) New() raklib.Packet {
	return new(ServerHandshakeDataPacket)
}

func (pk *ServerHandshakeDataPacket) Encode() error {
	err := pk.BasePacket.Encode(pk)
	if err != nil {
		return err
	}

	err = pk.PutAddressSystemAddress(pk.ClientAddr)
	if err != nil {
		return err
	}

	err = pk.PutShort(pk.SystemIndex)
	if err != nil {
		return err
	}

	for i := 0; i < 10; i++ {
		err = pk.PutAddressSystemAddress(pk.SystemAddresses[i])
		if err != nil {
			return err
		}
	}

	err = pk.PutLong(pk.RequestTime)
	if err != nil {
		return err
	}

	err = pk.PutLong(pk.Time)
	if err != nil {
		return err
	}

	return nil
}

func (pk *ServerHandshakeDataPacket) Decode() error {
	err := pk.BasePacket.Decode(pk)
	if err != nil {
		return err
	}

	err = pk.AddressSystemAddress(&pk.ClientAddr)
	if err != nil {
		return err
	}

	err = pk.Short(&pk.SystemIndex)
	if err != nil {
		return err
	}

	for i := 0; i < 10; i++ {
		err = pk.AddressSystemAddress(&pk.SystemAddresses[i])
		if err != nil {
			return err
		}
	}

	err = pk.Long(&pk.RequestTime)
	if err != nil {
		return err
	}

	err = pk.Long(&pk.Time)
	if err != nil {
		return err
	}

	return nil
}

type ClientHandshakeDataPacket struct {
	BasePacket
}

func (ClientHandshakeDataPacket) ID() byte {
	return IDClientHandshakeDataPacket
}

func (ClientHandshakeDataPacket) New() raklib.Packet {
	return new(ClientHandshakeDataPacket)
}

func (pk *ClientHandshakeDataPacket) Encode() error {
	err := pk.BasePacket.Encode(pk)
	if err != nil {
		return err
	}

	return nil
}

func (pk *ClientHandshakeDataPacket) Decode() error {
	err := pk.BasePacket.Decode(pk)
	if err != nil {
		return err
	}

	return err
}

type ClientDisconnectDataPacket struct {
	BasePacket

	Time int64
	UUID int64
}

func (ClientDisconnectDataPacket) ID() byte {
	return IDClientDisconnectDataPacket
}

func (ClientDisconnectDataPacket) New() raklib.Packet {
	return new(ClientDisconnectDataPacket)
}

func (pk *ClientDisconnectDataPacket) Encode() error {
	err := pk.BasePacket.Encode(pk)
	if err != nil {
		return err
	}

	err = pk.PutLong(pk.Time)
	if err != nil {
		return err
	}

	err = pk.PutLong(pk.UUID)
	if err != nil {
		return err
	}

	return nil
}

func (pk *ClientDisconnectDataPacket) Decode() error {
	err := pk.BasePacket.Decode(pk)
	if err != nil {
		return err
	}

	err = pk.Long(&pk.Time)
	if err != nil {
		return err
	}

	err = pk.Long(&pk.UUID)
	if err != nil {
		return err
	}

	return err
}

type UnconnectedPongPacket struct {
	BasePacket

	PingID     int64
	ServerID   int64
	Magic      string
	ServerName string
}

func (UnconnectedPongPacket) ID() byte {
	return IDUnconnectedPong
}

func (UnconnectedPongPacket) New() raklib.Packet {
	return new(UnconnectedPongPacket)
}

func (pk *UnconnectedPongPacket) Encode() error {
	err := pk.BasePacket.Encode(pk)
	if err != nil {
		return err
	}

	err = pk.PutLong(pk.PingID)
	if err != nil {
		return err
	}

	err = pk.PutLong(pk.ServerID)
	if err != nil {
		return err
	}

	err = pk.PutHexString(raklib.Magic)
	if err != nil {
		return err
	}

	err = pk.PutString(pk.ServerName)
	if err != nil {
		return err
	}

	return nil
}

func (pk *UnconnectedPongPacket) Decode() error {
	err := pk.BasePacket.Decode(pk)
	if err != nil {
		return err
	}

	err = pk.Long(&pk.PingID)
	if err != nil {
		return err
	}

	err = pk.Long(&pk.ServerID)
	if err != nil {
		return err
	}

	err = pk.String(&pk.Magic)
	if err != nil {
		return err
	}

	err = pk.String(&pk.ServerName)
	if err != nil {
		return err
	}

	return err
}
