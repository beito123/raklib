package raklib

/*
	Raklib

	Copyright (c) 2018 beito

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.
*/

import (
	"net"
	"strconv"
)

const (
	// Version is version of the Raknet library
	Version = "v1.0.0"

	// ProtocolVersion is supported version of raknet protocol
	ProtocolVersion = 8
)

const (
	Magic = "00ffff00fefefefefdfdfdfd12345678"
)

// Packet is basic Raknet packet interface
type Packet interface {
	ID() byte
	New() Packet
	Encode() error
	Decode() error
}

// SystemAddress is internal address for Raknet
type SystemAddress struct {
	IP   net.IP
	Port uint16
}

// SetLoopback sets loopback address
func (addr *SystemAddress) SetLoopback() {
	if len(addr.IP) == net.IPv4len {
		addr.IP = net.ParseIP("127.0.0.1")
	} else {
		addr.IP = net.IPv6loopback // "::1"
	}
}

// IsLoopback returns whether this is loopback address
func (addr *SystemAddress) IsLoopback() bool {
	return addr.IP.IsLoopback()
}

// Version returns the ip address version (4 or 6)
func (addr *SystemAddress) Version() int {
	if len(addr.IP) == net.IPv6len {
		return 6
	}

	return 4
}

// Equal returns whether sub is the same address
func (addr *SystemAddress) Equal(sub *SystemAddress) bool {
	return addr.IP.Equal(sub.IP) && addr.Port == sub.Port
}

// String returns as string
// Format: 192.168.11.1:8080, [fc00::]:8080
func (addr *SystemAddress) String() string {
	if len(addr.IP) == net.IPv6len {
		return "[" + addr.IP.String() + "]:" + strconv.Itoa(int(addr.Port))
	}

	return addr.IP.String() + ":" + strconv.Itoa(int(addr.Port))
}

// NewSystemAddress returns a new SystemAddress from string
func NewSystemAddress(addr string, port uint16) *SystemAddress {
	return &SystemAddress{
		IP:   net.ParseIP(addr),
		Port: port,
	}
}

// NewSystemAddress returns a new SystemAddress from bytes
func NewSystemAddressBytes(addr []byte, port uint16) *SystemAddress {
	return &SystemAddress{
		IP:   net.IP(addr).To16(),
		Port: port,
	}
}