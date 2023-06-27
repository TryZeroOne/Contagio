package methods

import (
	"bytes"
	"encoding/binary"
	"math/rand"
	"net"
	"syscall"
	"time"
)

type TCPHeader struct {
	SourcePort      uint16
	DestinationPort uint16
	SequenceNumber  uint32
	AckNumber       uint32
	DataOffset      uint8
	Flags           uint8
	WindowSize      uint16
	Checksum        uint16
	UrgentPointer   uint16
	Options         []byte
}

type IpHeader struct {
	Version  int
	Len      int
	TOS      int
	TotalLen int
	ID       int
	Flags    int
	FragOff  int
	TTL      int
	Protocol int
	Checksum int
	Src      net.IP
	Dst      net.IP
	Options  []byte
}

type PsdHeader struct {
	SrcAddr   [4]uint8
	DstAddr   [4]uint8
	Zero      uint8
	ProtoType uint8
	TcpLength uint16
}

func CreateIpHeader(srcIp, dstIp net.IP) []byte {

	defer Catch()

	h := &IpHeader{
		ID:       1,
		TTL:      1,
		Protocol: syscall.IPPROTO_TCP,
		Checksum: 0,
		Src:      srcIp,
		Dst:      dstIp,
	}
	h.TTL = rand.Intn(255-60) + 60
	return h.Marshal()
}

func (h *IpHeader) Marshal() []byte {
	defer Catch()

	if h == nil {
		return nil
	}

	hdrlen := 20 + len(h.Options)
	b := make([]byte, hdrlen)

	b[0] = byte(4<<4 | (hdrlen >> 2 & 0x0f))
	b[1] = byte(h.TOS)

	binary.BigEndian.PutUint16(b[2:4], uint16(h.TotalLen))
	binary.BigEndian.PutUint16(b[4:6], uint16(h.ID))

	flagsAndFragOff := (h.FragOff & 0x1fff) | int(h.Flags<<13)
	binary.BigEndian.PutUint16(b[6:8], uint16(flagsAndFragOff))

	b[8] = byte(h.TTL)
	b[9] = byte(h.Protocol)

	binary.BigEndian.PutUint16(b[10:12], uint16(h.Checksum))

	if ip := h.Src.To4(); ip != nil {
		copy(b[12:16], ip[:net.IPv4len])
	}

	if ip := h.Dst.To4(); ip != nil {
		copy(b[16:20], ip[:net.IPv4len])
	} else {
		return nil
	}

	if len(h.Options) > 0 {
		copy(b[20:], h.Options)
	}

	return b
}

func CreateTcpHeader(srcIp, dstIp net.IP, dstPort int, flags uint8) []byte {
	rand.NewSource(time.Now().UnixNano())

	h := &TCPHeader{
		SourcePort:      9765,
		DestinationPort: uint16(dstPort),
		SequenceNumber:  690,
		AckNumber:       0,
		Flags:           flags,
		WindowSize:      2048,
		UrgentPointer:   0,
	}

	h.SourcePort = uint16(rand.Intn(1<<16-1)%16383 + 49152)
	h.SequenceNumber = uint32(rand.Intn(1<<10 - 1))
	h.WindowSize = uint16(rand.Intn(10<<10-1000) + 1000)
	h.Marshal()

	var psdheader PsdHeader

	copy(psdheader.SrcAddr[:4], srcIp)
	copy(psdheader.DstAddr[:4], dstIp)

	psdheader.Zero = 0
	psdheader.ProtoType = syscall.IPPROTO_TCP
	psdheader.TcpLength = uint16(20)

	var buffer bytes.Buffer

	binary.Write(&buffer, binary.BigEndian, psdheader)
	buffs := h.Marshal()
	buffer.Write(buffs)
	h.Checksum = uint16(checksum(buffer.Bytes()))

	return h.Marshal()
}

func checksum(data []byte) uint16 {
	var (
		sum    uint32
		length int = len(data)
		index  int
	)
	for length > 1 {
		sum += uint32(data[index])<<8 + uint32(data[index+1])
		index += 2
		length -= 2
	}
	if length > 0 {
		sum += uint32(data[index])
	}
	sum += (sum >> 16)

	return uint16(^sum)
}

func (h *TCPHeader) Marshal() []byte {
	defer Catch()

	if h == nil {
		return nil
	}

	hdrlen := 20 + len(h.Options)
	b := make([]byte, hdrlen)

	binary.BigEndian.PutUint16(b[0:2], h.SourcePort)
	binary.BigEndian.PutUint16(b[2:4], h.DestinationPort)

	binary.BigEndian.PutUint32(b[4:8], h.SequenceNumber)
	binary.BigEndian.PutUint32(b[8:12], h.AckNumber)

	b[12] = uint8(hdrlen / 4 << 4)
	b[13] = uint8(h.Flags)

	binary.BigEndian.PutUint16(b[14:16], uint16(h.WindowSize))
	binary.BigEndian.PutUint16(b[16:18], uint16(h.Checksum))
	binary.BigEndian.PutUint16(b[18:20], uint16(h.UrgentPointer))

	if len(h.Options) > 0 {
		copy(b[20:], h.Options)
	}

	return b
}
