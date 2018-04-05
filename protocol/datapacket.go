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
	"github.com/beito123/binary"
	"math"
	"bytes"
	"github.com/beito123/raklib"
)

// Ref: http://www.jenkinssoftware.com/raknet/manual/Doxygen/PacketPriority_8h.html#e41fa01235e99dced384d137fa874a7e

// Reliability decides reliable and ordered of packet
type Reliability int

const (
	// Unreliable is normal UDP packet.
	Unreliable                    Reliability = iota
	UnreliableSequenced           
	Reliable                      
	ReliableOrdered               
	ReliableSequenced             
	UnreliableWithACKReceipt      
	ReliableWithACKReceipt        
	ReliableOrderedWithACKReceipt 
)

func (r Reliability) IsReliable() bool {
	return r == Reliable || r == ReliableOrdered ||
		r == ReliableSequenced || r == ReliableWithACKReceipt ||
		r == ReliableOrderedWithACKReceipt
}

func (r Reliability) IsOrdered() bool {
	return r == UnreliableSequenced || r == ReliableOrdered ||
		r == ReliableSequenced || r == ReliableOrderedWithACKReceipt
}

func (r Reliability) IsSequenced() bool {
	return r == UnreliableSequenced || r == ReliableSequenced
}

func (r Reliability) IsNeededACK() bool {
	return r == UnreliableWithACKReceipt || r == ReliableWithACKReceipt ||
		r == ReliableOrderedWithACKReceipt
}

func (r Reliability) ToBinary() byte {
	var b byte

	if r.IsReliable() {
		b |= 1 << 2
	}

	if r.IsOrdered() {
		b |= 1 << 1
	}

	if r.IsSequenced() {
		b |= 1
	}

	return b
}

func ReliabilityFromBinary(b byte) Reliability {
	if (b & 0x04) > 0 { // Reliable
		if (b & 0x02) > 0 { // Ordered
			return ReliableOrdered
		} else if (b & 0x01) > 0 { // Sequenced
			return ReliableSequenced
		} else {
			return Reliable
		}
	} else { // Unreliable
		if (b & 0x01) > 0 { // Sequenced
			return UnreliableSequenced
		}
	}

	return Unreliable
}

// EncapsulatedPacket .
type EncapsulatedPacket struct {
	binary.Stream
	Flags  byte
	Length uint16

	ReliableIndex binary.Triad // only if reliable

	SequenceIndex binary.Triad // sequenced

	OrderIndex   binary.Triad // order
	OrderChannel byte

	SplitCount int32 // Split
	SplitID    uint16
	SplitIndex int32

	Body []byte // Buffer

	// When Encode, use under the vars

	Reliability Reliability
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

// Encode encodes packet to binary
func (epk *EncapsulatedPacket) Encode() error {
	epk.EncodeFlags()

	err := epk.PutByte(epk.Flags)
	if err != nil {
		return err
	}

	err = epk.PutShort(uint16(epk.Buffer.Len() << 3))
	if err != nil {
		return err
	}

	if epk.Reliability.IsReliable() { // Reliable
		err = epk.PutTriad(epk.ReliableIndex)
		if err != nil {
			return err
		}
	}

	if epk.Reliability.IsOrdered() { // Ordered
		err = epk.PutTriad(epk.OrderIndex)
		if err != nil {
			return err
		}

		err = epk.PutByte(epk.OrderChannel)
		if err != nil {
			return err
		}
	}

	if epk.Reliability.IsSequenced() { // Sequenced
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

// EncodeFlags encodes the internal variables to binary
func (epk *EncapsulatedPacket) EncodeFlags() {
	var splitFlag byte
	if epk.HasSplit {
		splitFlag = 1
	}

	var flags byte

	flags |= epk.Reliability.ToBinary() << 5
	flags |= splitFlag << 4

	epk.Flags = flags
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
	// xxx: Reliable, Ordered, Sequenced. yyy: Has Split. zzzz: empty
	epk.Reliability = ReliabilityFromBinary(epk.Flags >> 5)
	epk.HasSplit = (epk.Flags & 0x10) > 0

	if epk.Reliability.IsReliable() { // Reliable
		err = epk.Triad(&epk.ReliableIndex)
		if err != nil {
			return err
		}
	}

	if epk.Reliability.IsOrdered() { // Ordered
		err = epk.Triad(&epk.OrderIndex)
		if err != nil {
			return err
		}

		err = epk.Byte(&epk.OrderChannel)
		if err != nil {
			return err
		}
	}

	if epk.Reliability.IsSequenced() { // Sequenced
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

	bodyLen := int(math.Ceil(float64(epk.Length / 8)))
	epk.Body = epk.Get(bodyLen)

	return nil
}

func (epk *EncapsulatedPacket) GetBuffer() []byte {
	return epk.Bytes()
}

func (epk *EncapsulatedPacket) Len() int {
	ln := 3 // flags(1byte), short(2byte)
	ln += len(epk.Body)

	if epk.Reliability.IsReliable() {
		ln += 3
	}

	if epk.Reliability.IsSequenced() {
		ln += 3
	}

	if epk.Reliability.IsOrdered() {
		ln += 4
	}

	if epk.HasSplit {
		ln += 10
	}

	return ln

}

// DataPacket

type DataPacket struct {
	BasePacket
	Index   binary.Triad
	Packets []*EncapsulatedPacket
}

func (DataPacket) ID() byte {
	return 0xff
}

func (DataPacket) New() raklib.Packet {
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

func (DataPacket0) New() raklib.Packet {
	return new(DataPacket0)
}

type DataPacket1 struct {
	DataPacket
}

func (DataPacket1) ID() byte {
	return IDDataPacket1
}

func (DataPacket1) New() raklib.Packet {
	return new(DataPacket1)
}

type DataPacket2 struct {
	DataPacket
}

func (DataPacket2) ID() byte {
	return IDDataPacket2
}

func (DataPacket2) New() raklib.Packet {
	return new(DataPacket2)
}

type DataPacket3 struct {
	DataPacket
}

func (DataPacket3) ID() byte {
	return IDDataPacket3
}

func (DataPacket3) New() raklib.Packet {
	return new(DataPacket3)
}

type DataPacket4 struct {
	DataPacket
}

func (DataPacket4) ID() byte {
	return IDDataPacket4
}

func (DataPacket4) New() raklib.Packet {
	return new(DataPacket4)
}

type DataPacket5 struct {
	DataPacket
}

func (DataPacket5) ID() byte {
	return IDDataPacket5
}

func (DataPacket5) New() raklib.Packet {
	return new(DataPacket5)
}

type DataPacket6 struct {
	DataPacket
}

func (DataPacket6) ID() byte {
	return IDDataPacket6
}

func (DataPacket6) New() raklib.Packet {
	return new(DataPacket6)
}

type DataPacket7 struct {
	DataPacket
}

func (DataPacket7) ID() byte {
	return IDDataPacket7
}

func (DataPacket7) New() raklib.Packet {
	return new(DataPacket7)
}

type DataPacket8 struct {
	DataPacket
}

func (DataPacket8) ID() byte {
	return IDDataPacket8
}

func (DataPacket8) New() raklib.Packet {
	return new(DataPacket8)
}

type DataPacket9 struct {
	DataPacket
}

func (DataPacket9) ID() byte {
	return IDDataPacket9
}

func (DataPacket9) New() raklib.Packet {
	return new(DataPacket9)
}

type DataPacketA struct {
	DataPacket
}

func (DataPacketA) ID() byte {
	return IDDataPacketA
}

func (DataPacketA) New() raklib.Packet {
	return new(DataPacketA)
}

type DataPacketB struct {
	DataPacket
}

func (DataPacketB) ID() byte {
	return IDDataPacketB
}

func (DataPacketB) New() raklib.Packet {
	return new(DataPacketB)
}

type DataPacketC struct {
	DataPacket
}

func (DataPacketC) ID() byte {
	return IDDataPacketC
}

func (DataPacketC) New() raklib.Packet {
	return new(DataPacketC)
}

type DataPacketD struct {
	DataPacket
}

func (DataPacketD) ID() byte {
	return IDDataPacketD
}

func (DataPacketD) New() raklib.Packet {
	return new(DataPacketD)
}

type DataPacketE struct {
	DataPacket
}

func (DataPacketE) ID() byte {
	return IDDataPacketE
}

func (DataPacketE) New() raklib.Packet {
	return new(DataPacketE)
}

type DataPacketF struct {
	DataPacket
}

func (DataPacketF) ID() byte {
	return IDDataPacketF
}

func (DataPacketF) New() raklib.Packet {
	return new(DataPacketF)
}
