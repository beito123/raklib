package binary

import (
	"github.com/beito123/binary"
	"encoding/hex"
	"strconv"
	"strings"
	"github.com/beito123/raklib"
)

// RaknetStream is binary stream for Raknet
type RaknetStream struct {
	binary.Stream
}

// String sets string(len short, str string) got from buffer to value
func (bs *RaknetStream) String(value *string) error {
	var n uint16
	err := bs.Short(&n)
	if err != nil {
		return err
	}

	*value = string(bs.Get(int(n)))
	return nil
}

// PutString puts string(len short, str string) to Buffer
func (bs *RaknetStream) PutString(value string) error {
	n := uint16(len(value))
	err := bs.PutShort(n)
	if err != nil {
		return err
	}
	return bs.Put([]byte(value))
}

// HexString gets hex string from Buffer (for Magic)
func (bs *RaknetStream) HexString(n int, value *string) {
	*value = hex.EncodeToString(bs.Buffer.Bytes())
}

// PutHexString puts hex string to Buffer
func (bs *RaknetStream) PutHexString(value string) error {
	bytes, err := hex.DecodeString(value)
	if err != nil {
		return err
	}
	return bs.Put(bytes)
}

// Address sets address got from Buffer to addr and port
// address(version byte, address byte x4, port ushort)
func (bs *RaknetStream) Address(addr *string, port *uint16) error {
	var version byte
	err := bs.Byte(&version)
	if err != nil {
		return err
	}

	var address string

	if version == 4 {
		var bytes byte
		for i := 0; i < 4; i++ {
			err = bs.Byte(&bytes)
			if err != nil {
				return err
			}

			address = address + strconv.Itoa(int(^bytes&0xff))
			if i < 3 {
				address = address + "."
			}
		}
		addr = &address

		err = bs.Short(port)
		if err != nil {
			return err
		}
	} else {
		// IPv6
	}

	return nil
}

// PutAddress puts address to Buffer
// address(version byte, address byte x4, port ushort)
func (bs *RaknetStream) PutAddress(addr string, port uint16, version byte) error {
	err := bs.PutByte(version)
	if err != nil {
		return err
	}

	if version == 4 {
		for _, str := range strings.Split(addr, ".") {
			i, _ := strconv.Atoi(str)
			err = bs.PutByte(^byte(i) & 0xff)
			if err != nil {
				return err
			}
		}
		err = bs.PutShort(port)
		if err != nil {
			return err
		}
	} else {
		// ipv6
	}

	return nil
}

// AddressSystemAddress sets address got from Buffer to SystemAddress
func (bs *RaknetStream) AddressSystemAddress(addr *raklib.SystemAddress) error {
	var add string
	var port uint16

	err := bs.Address(&add, &port)
	if err != nil {
		return err
	}

	naddr := raklib.NewSystemAddress(add, port)

	*addr = *naddr

	return nil
}

// PutAddressSystemAddress puts address from UDPAddr to Buffer
func (bs *RaknetStream) PutAddressSystemAddress(addr raklib.SystemAddress) error {
	return bs.PutAddress(addr.IP.String(), addr.Port, byte(addr.Version()))
}
