package protocol

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
