package raklib

/*
	Raklib

	Copyright (c) 2018 beito

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU Lesser General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.
*/

import (
	"bytes"
	"math"
	"net"

	"github.com/beito123/binary"
)

const (
	IDPingDataPacket                 = 0x00
	IDUnconnectedPing                = 0x01
	IDUnconnectedPingOpenConnections = 0x02
	IDPongDataPacket                 = 0x03
	IDOpenConnectionRequest1         = 0x05
	IDOpenConnectionReply1           = 0x06
	IDOpenConnectionRequest2         = 0x07
	IDOpenConnectionReply2           = 0x08
	IDClientConnectDataPacket        = 0x09
	IDServerHandshakeDataPacket      = 0x10
	IDClientHandshakeDataPacket      = 0x13
	IDClientDisconnectDataPacket     = 0x15
	IDDataPacket0                    = 0x80
	IDDataPacket1                    = 0x81
	IDDataPacket2                    = 0x82
	IDDataPacket3                    = 0x83
	IDDataPacket4                    = 0x84
	IDDataPacket5                    = 0x85
	IDDataPacket6                    = 0x86
	IDDataPacket7                    = 0x87
	IDDataPacket8                    = 0x88
	IDDataPacket9                    = 0x89
	IDDataPacketA                    = 0x8A
	IDDataPacketB                    = 0x8B
	IDDataPacketC                    = 0x8C
	IDDataPacketD                    = 0x8D
	IDDataPacketE                    = 0x8E
	IDDataPacketF                    = 0x8F
	IDUnconnectedPong                = 0x1c
	IDAdvertiseSystem                = 0x1d
	IDNACK                           = 0xa0
	IDACK                            = 0xc0
	IDUnknownPacket                  = 0xff
)

type Packet interface {
	ID() byte
	New() Packet
	Encode() error
	Decode() error
}

type BasePacket struct {
	binary.Stream
}

func (base *BasePacket) Encode(pk Packet) error {
	if base.Buffer == nil {
		base.Buffer = &bytes.Buffer{}
	}

	err := base.PutByte(pk.ID())
	if err != nil {
		return err
	}

	return nil
}

func (base *BasePacket) Decode(pk Packet) error {
	if base.Buffer == nil {
		return NoSetBufferError{}
	}

	base.Skip(1) // id

	return nil
}

// OpenConnectionRequest2 .
// Client To Server
type OpenConnectionRequest2 struct {
	BasePacket
	Magic         string
	ServerAddress string
	ServerPort    uint16
	MtuSize       uint16
	ClientID      int64
}

// ID .
func (OpenConnectionRequest2) ID() byte {
	return IDOpenConnectionRequest2
}

func (OpenConnectionRequest2) New() Packet {
	return new(OpenConnectionRequest2)
}

// Encode encodes a packet
func (pk *OpenConnectionRequest2) Encode() error {
	pk.BasePacket.Encode(pk)

	err := pk.PutHexString(Magic)
	if err != nil {
		return err
	}

	err = pk.PutAddress(pk.ServerAddress, pk.ServerPort, 4)
	if err != nil {
		return err
	}

	err = pk.PutShort(pk.MtuSize)
	if err != nil {
		return err
	}

	err = pk.PutLong(pk.ClientID)
	if err != nil {
		return err
	}

	return nil
}

// Decode .
func (pk *OpenConnectionRequest2) Decode() error {
	pk.BasePacket.Decode(pk)

	pk.HexString(16, &pk.Magic)

	err := pk.Address(&pk.ServerAddress, &pk.ServerPort)
	if err != nil {
		return err
	}

	err = pk.Short(&pk.MtuSize)
	if err != nil {
		return err
	}

	err = pk.Long(&pk.ClientID)
	if err != nil {
		return err
	}

	return nil
}

// OpenConnectionReply2 Server To Client
type OpenConnectionReply2 struct {
	BasePacket
	Magic         string
	ServerID      int64
	ClientAddress string
	ClientPort    uint16
	MtuSize       uint16
}

// ID .
func (OpenConnectionReply2) ID() byte {
	return IDOpenConnectionReply2
}

func (OpenConnectionReply2) New() Packet {
	return new(OpenConnectionReply2)
}

// Encode .
func (pk *OpenConnectionReply2) Encode() error {
	pk.BasePacket.Encode(pk)

	err := pk.PutHexString(Magic)
	if err != nil {
		return err
	}

	err = pk.PutLong(pk.ServerID)
	if err != nil {
		return err
	}

	err = pk.PutAddress(pk.ClientAddress, pk.ClientPort, 4)
	if err != nil {
		return err
	}

	err = pk.PutShort(pk.MtuSize)
	if err != nil {
		return err
	}

	err = pk.PutByte(0) // security
	if err != nil {
		return err
	}

	return nil
}

// Decode .
func (pk *OpenConnectionReply2) Decode() error {
	pk.BasePacket.Decode(pk)

	pk.HexString(16, &pk.Magic)

	err := pk.Long(&pk.ServerID)
	if err != nil {
		return err
	}

	err = pk.Address(&pk.ClientAddress, &pk.ClientPort)
	if err != nil {
		return err
	}

	err = pk.Short(&pk.MtuSize)
	if err != nil {
		return err
	}

	return nil
}

/*
	EncapsulatedPacket
*/

// EncapsulatedPacket .
type EncapsulatedPacket struct {
	binary.Stream
	Flags  byte // Base
	Length uint16

	ReliableIndex binary.Triad // only if reliable

	SequenceIndex binary.Triad // sequenced

	OrderIndex   binary.Triad // order
	OrderChannel byte

	SplitCount int32 // Split
	SplitID    uint16
	SplitIndex int32

	Body []byte // Buffer

	Reliability byte // var flags
	HasSplit    bool

	// NeedACK       bool //ACK
	// IdentifierACK int
}

// NewEncapsulatedPacket .
func NewEncapsulatedPacket(buffer *bytes.Buffer) *EncapsulatedPacket {
	return &EncapsulatedPacket{
		Stream: binary.Stream{
			Buffer: buffer,
		},
	}
}

// Encode . ToBinary
func (epk *EncapsulatedPacket) Encode() error {

	err := epk.PutByte(epk.Flags)
	if err != nil {
		return err
	}

	err = epk.PutShort(uint16(epk.Buffer.Len() << 3))
	if err != nil {
		return err
	}

	epk.Reliability = epk.Flags >> 5
	epk.HasSplit = (epk.Flags & 0x10) > 0

	if (epk.Reliability & 0x04) > 0 { // Reliable
		err = epk.PutTriad(epk.ReliableIndex)
		if err != nil {
			return err
		}
	}

	if (epk.Reliability & 0x02) > 0 { // Ordered
		err = epk.PutTriad(epk.OrderIndex)
		if err != nil {
			return err
		}

		err = epk.PutByte(epk.OrderChannel)
		if err != nil {
			return err
		}
	}

	if (epk.Reliability & 0x01) > 0 { // Sequenced
		err = epk.PutTriad(epk.SequenceIndex)
		if err != nil {
			return err
		}
	}

	if epk.HasSplit {
		err = epk.PutInt(epk.SplitCount)
		if err != nil {
			return err
		}

		err = epk.PutShort(epk.SplitID)
		if err != nil {
			return err
		}

		err = epk.PutInt(epk.SplitIndex)
		if err != nil {
			return err
		}
	}

	epk.Put(epk.Buffer.Bytes())

	return nil
}

func (epk *EncapsulatedPacket) EncodeFlags() {
	splitFlag := 0
	if epk.HasSplit {
		splitFlag = 1
	}

	epk.Flags = epk.Reliability | byte(splitFlag)
}

// Decode . FromBinary
func (epk *EncapsulatedPacket) Decode() error {

	err := epk.Byte(&epk.Flags)
	if err != nil {
		return err
	}

	err = epk.Short(&epk.Length)
	if err != nil {
		return err
	}

	// ref: http://www.jenkinssoftware.com/raknet/manual/reliabilitytypes.html

	// xxx y zzzz
	epk.Reliability = epk.Flags >> 5
	epk.HasSplit = (epk.Flags & 0x10) > 0

	if (epk.Reliability & 0x04) > 0 { // Reliable
		err = epk.Triad(&epk.ReliableIndex)
		if err != nil {
			return err
		}
	}

	if (epk.Reliability & 0x02) > 0 { // Ordered
		err = epk.Triad(&epk.OrderIndex)
		if err != nil {
			return err
		}

		err = epk.Byte(&epk.OrderChannel)
		if err != nil {
			return err
		}
	}

	if (epk.Reliability & 0x01) > 0 { // Sequenced
		err = epk.Triad(&epk.SequenceIndex)
		if err != nil {
			return err
		}
	}

	if epk.HasSplit {
		err = epk.Int(&epk.SplitCount)
		if err != nil {
			return err
		}

		err = epk.Short(&epk.SplitID)
		if err != nil {
			return err
		}

		err = epk.Int(&epk.SplitIndex)
		if err != nil {
			return err
		}
	}

	len := int(math.Ceil(float64(epk.Length / 8)))
	epk.Body = epk.Get(len)

	return nil
}

func (epk *EncapsulatedPacket) GetBuffer() []byte {
	return epk.Bytes()
}

/*
	DataPacket
*/

type DataPacket struct {
	BasePacket
	Index   binary.Triad
	Packets []*EncapsulatedPacket
}

func (DataPacket) ID() byte {
	return 0xff
}

func (DataPacket) New() Packet {
	return new(DataPacket)
}

func (bp *DataPacket) Encode() error {
	bp.BasePacket.Encode(bp)

	bp.PutTriad(bp.Index)
	for _, pk := range bp.Packets {
		pk.Encode()
		bp.Put(pk.GetBuffer())
	}
	return nil
}

func (bp *DataPacket) Decode() error {
	bp.BasePacket.Decode(bp)

	bp.Triad(&bp.Index)

	for bp.Buffer.Len() > 0 {
		epk := NewEncapsulatedPacket(bp.Buffer)
		bp.Packets = append(bp.Packets, epk)
	}

	return nil
}

type DataPacket0 struct {
	DataPacket
}

func (DataPacket0) ID() byte {
	return IDDataPacket0
}

func (DataPacket0) New() Packet {
	return new(DataPacket0)
}

type DataPacket1 struct {
	DataPacket
}

func (DataPacket1) ID() byte {
	return IDDataPacket1
}

func (DataPacket1) New() Packet {
	return new(DataPacket1)
}

type DataPacket2 struct {
	DataPacket
}

func (DataPacket2) ID() byte {
	return IDDataPacket2
}

func (DataPacket2) New() Packet {
	return new(DataPacket2)
}

type DataPacket3 struct {
	DataPacket
}

func (DataPacket3) ID() byte {
	return IDDataPacket3
}

func (DataPacket3) New() Packet {
	return new(DataPacket3)
}

type DataPacket4 struct {
	DataPacket
}

func (DataPacket4) ID() byte {
	return IDDataPacket4
}

func (DataPacket4) New() Packet {
	return new(DataPacket4)
}

type DataPacket5 struct {
	DataPacket
}

func (DataPacket5) ID() byte {
	return IDDataPacket5
}

func (DataPacket5) New() Packet {
	return new(DataPacket5)
}

type DataPacket6 struct {
	DataPacket
}

func (DataPacket6) ID() byte {
	return IDDataPacket6
}

func (DataPacket6) New() Packet {
	return new(DataPacket6)
}

type DataPacket7 struct {
	DataPacket
}

func (DataPacket7) ID() byte {
	return IDDataPacket7
}

func (DataPacket7) New() Packet {
	return new(DataPacket7)
}

type DataPacket8 struct {
	DataPacket
}

func (DataPacket8) ID() byte {
	return IDDataPacket8
}

func (DataPacket8) New() Packet {
	return new(DataPacket8)
}

type DataPacket9 struct {
	DataPacket
}

func (DataPacket9) ID() byte {
	return IDDataPacket9
}

func (DataPacket9) New() Packet {
	return new(DataPacket9)
}

type DataPacketA struct {
	DataPacket
}

func (DataPacketA) ID() byte {
	return IDDataPacketA
}

func (DataPacketA) New() Packet {
	return new(DataPacketA)
}

type DataPacketB struct {
	DataPacket
}

func (DataPacketB) ID() byte {
	return IDDataPacketB
}

func (DataPacketB) New() Packet {
	return new(DataPacketB)
}

type DataPacketC struct {
	DataPacket
}

func (DataPacketC) ID() byte {
	return IDDataPacketC
}

func (DataPacketC) New() Packet {
	return new(DataPacketC)
}

type DataPacketD struct {
	DataPacket
}

func (DataPacketD) ID() byte {
	return IDDataPacketD
}

func (DataPacketD) New() Packet {
	return new(DataPacketD)
}

type DataPacketE struct {
	DataPacket
}

func (DataPacketE) ID() byte {
	return IDDataPacketE
}

func (DataPacketE) New() Packet {
	return new(DataPacketE)
}

type DataPacketF struct {
	DataPacket
}

func (DataPacketF) ID() byte {
	return IDDataPacketF
}

func (DataPacketF) New() Packet {
	return new(DataPacketF)
}

type PingDataPacket struct {
	BasePacket

	Time int64
}

func (PingDataPacket) ID() byte {
	return IDPingDataPacket
}

func (PingDataPacket) New() Packet {
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

func (UnconnectedPingPacket) New() Packet {
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

	err = pk.PutHexString(Magic)
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

func (UnconnectedPingOpenConnections) New() Packet {
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

func (PongDataPacket) New() Packet {
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

func (OpenConnectionRequest1Packet) New() Packet {
	return new(OpenConnectionRequest1Packet)
}

func (pk *OpenConnectionRequest1Packet) Encode() error {
	err := pk.BasePacket.Encode(pk)
	if err != nil {
		return err
	}

	err = pk.PutHexString(Magic)
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

func (OpenConnectionReply1Packet) New() Packet {
	return new(OpenConnectionReply1Packet)
}

func (pk *OpenConnectionReply1Packet) Encode() error {
	err := pk.BasePacket.Encode(pk)
	if err != nil {
		return err
	}

	err = pk.PutHexString(Magic)
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
	ServerAddress net.UDPAddr
	MTU           uint16
	ClientUUID    int64
}

func (OpenConnectionRequest2Packet) ID() byte {
	return IDOpenConnectionRequest2
}

func (OpenConnectionRequest2Packet) New() Packet {
	return new(OpenConnectionRequest2Packet)
}

func (pk *OpenConnectionRequest2Packet) Encode() error {
	err := pk.BasePacket.Encode(pk)
	if err != nil {
		return err
	}

	err = pk.PutHexString(Magic)
	if err != nil {
		return err
	}

	err = pk.PutAddressUDPAddr(pk.ServerAddress)
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

	err = pk.AddressUDPAddr(&pk.ServerAddress)
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
	ClientAddress net.UDPAddr
	MTU           uint16
	Encryption    byte
}

func (OpenConnectionReply2Packet) ID() byte {
	return IDOpenConnectionReply2
}

func (OpenConnectionReply2Packet) New() Packet {
	return new(OpenConnectionReply2Packet)
}

func (pk *OpenConnectionReply2Packet) Encode() error {
	err := pk.BasePacket.Encode(pk)
	if err != nil {
		return err
	}

	err = pk.PutHexString(Magic)
	if err != nil {
		return err
	}

	err = pk.PutLong(pk.ServerUUID)
	if err != nil {
		return err
	}

	err = pk.PutAddressUDPAddr(pk.ClientAddress)
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

	err = pk.AddressUDPAddr(&pk.ClientAddress)
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

func (ClientConnectDataPacket) New() Packet {
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

	ClientAddr      net.UDPAddr
	SystemIndex     uint16
	SystemAddresses [10]net.UDPAddr
	RequestTime     int64
	Time            int64
}

func (ServerHandshakeDataPacket) ID() byte {
	return IDServerHandshakeDataPacket
}

func (ServerHandshakeDataPacket) New() Packet {
	return new(ServerHandshakeDataPacket)
}

func (pk *ServerHandshakeDataPacket) Encode() error {
	err := pk.BasePacket.Encode(pk)
	if err != nil {
		return err
	}

	err = pk.PutAddressUDPAddr(pk.ClientAddr)
	if err != nil {
		return err
	}

	err = pk.PutShort(pk.SystemIndex)
	if err != nil {
		return err
	}

	for i := 0; i < 10; i++ {
		err = pk.PutAddressUDPAddr(pk.SystemAddresses[i])
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

	err = pk.AddressUDPAddr(&pk.ClientAddr)
	if err != nil {
		return err
	}

	err = pk.Short(&pk.SystemIndex)
	if err != nil {
		return err
	}

	for i := 0; i < 10; i++ {
		err = pk.AddressUDPAddr(&pk.SystemAddresses[i])
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

func (ClientHandshakeDataPacket) New() Packet {
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

func (ClientDisconnectDataPacket) New() Packet {
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

func (UnconnectedPongPacket) New() Packet {
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

	err = pk.PutHexString(Magic)
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
