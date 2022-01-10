package wol

import (
	"errors"
	"net"

	"github.com/mdlayher/ethernet"
	"github.com/mdlayher/raw"
)

// ErrInvalidMAC represents invalid MAC address error
var ErrInvalidMAC = errors.New("invalid mac address")

// Type can be one of "ethernet" or "udp"
type Type string

const (
	// Ethernet type is used to encapsulate WoL magic packet in Ethernet frame
	Ethernet Type = "ethernet"
	// UDP type is used to encapsulate WoL magic packet in UDP datagram
	UDP Type = "udp"
)

// EtherType value used in Ethernet frame for WoL
const EtherType = 0x0842

// WOL holds configuration for the service
type WOL struct {
	Type      Type
	MAC       net.HardwareAddr
	IP        net.IP
	Port      int
	Interface string
}

// NewEther creates ethernet-encapsulated WoL instance
func NewEther(mac net.HardwareAddr, iface string) (*WOL, error) {
	if len(mac) != 6 {
		return nil, ErrInvalidMAC
	}

	return &WOL{
		Type:      Ethernet,
		MAC:       mac,
		Interface: iface,
	}, nil
}

// NewUDP creates udp-encapsulated WoL instance
func NewUDP(mac net.HardwareAddr, ip net.IP, port int) (*WOL, error) {
	if len(mac) != 6 {
		return nil, ErrInvalidMAC
	}

	return &WOL{
		Type: UDP,
		MAC:  mac,
		IP:   ip,
		Port: port,
	}, nil
}

func (w *WOL) createMagicPacket() []byte {
	packet := make([]byte, 102)
	copy(packet[0:], []byte{255, 255, 255, 255, 255, 255})

	for i := 1; i < 17; i++ {
		copy(packet[i*6:], w.MAC)
	}

	return packet
}

func (w *WOL) sendUDP() error {
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{IP: w.IP, Port: w.Port})
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Write(w.createMagicPacket())
	if err != nil {
		return err
	}

	return nil
}

func (w *WOL) sendEthernet() error {
	iface, err := net.InterfaceByName(w.Interface)
	if err != nil {
		return err
	}

	c, err := raw.ListenPacket(iface, EtherType, nil)
	if err != nil {
		return err
	}
	defer c.Close()

	f := &ethernet.Frame{
		Destination: w.MAC,
		Source:      iface.HardwareAddr,
		EtherType:   EtherType,
		Payload:     w.createMagicPacket(),
	}

	b, err := f.MarshalBinary()
	if err != nil {
		return err
	}

	addr := &raw.Addr{HardwareAddr: w.MAC}
	if _, err := c.WriteTo(b, addr); err != nil {
		return err
	}

	return nil
}

// Send sends WoL magic packet
func (w *WOL) Send() error {
	switch w.Type {
	case Ethernet:
		return w.sendEthernet()
	case UDP:
		return w.sendUDP()
	}

	return nil
}
