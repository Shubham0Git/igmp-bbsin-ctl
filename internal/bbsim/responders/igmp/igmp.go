package igmp

import (
	//	"context"
	//	"crypto/md5"
	//	"errors"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	bbsim "github.com/opencord/bbsim/internal/bbsim/types"
	"net"
	"time"
)

func SendIGMPMembershipReportV2(ponPortId uint32, onuId uint32, macAddress net.HardwareAddr, stream bbsim.Stream) {
	igmp := createIgmpV2Packet()
	pkt, err := serializeIgmpPacket(ponPortId, onuId, macAddress, igmp)

	if err != nil {
		//TODO : Error Handling
	}
	if pkt != nil {
	}
	//TODO : create bbsim message and send to stream
}

func createIgmpV2Packet() *layers.IGMP {

	igmpDefault := layers.IGMP{
		Type:            layers.IGMPMembershipReportV2,
		MaxResponseTime: time.Duration(1),
		Checksum:        0,
		GroupAddress:    net.IPv4zero,
		Version:         2,
	}

	calculateChecksum(igmpDefault)

	//returning igmp packet
	return &igmpDefault
}

func calculateChecksum(igmpDefault layers.IGMP) {
	//TODO : calculate checksum as per rfc 2236 and set into igmpDefault

}

func serializeIgmpPacket(intfId uint32, onuId uint32, srcMac net.HardwareAddr, igmp *layers.IGMP) ([]byte, error) {
	buffer := gopacket.NewSerializeBuffer()
	options := gopacket.SerializeOptions{
		ComputeChecksums: true,
		FixLengths:       true,
	}

	ethernetLayer := &layers.Ethernet{
		SrcMAC:       srcMac,
		DstMAC:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		EthernetType: layers.EthernetTypeIPv4,
	}

	ipLayer := &layers.IPv4{
		Version:  4,
		TOS:      0x10,
		TTL:      128,
		SrcIP:    []byte{0, 0, 0, 0},
		DstIP:    []byte{255, 255, 255, 255},
		Protocol: layers.IPProtocolIGMP,
	}

	if err := gopacket.SerializeLayers(buffer, options, ethernetLayer, ipLayer); err != nil {
		return nil, err
	}

	bytes := buffer.Bytes()
	return bytes, nil
}

func SendIGMPLeaveGroup(stream bbsim.Stream) {
	//TODO : Implement IGMPLeave
}
