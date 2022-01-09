package wol

import (
	"errors"
	"net"

	"github.com/mdlayher/ethernet"
	"github.com/mdlayher/raw"
)

var ErrInvalidMAC = errors.New("invalid mac address")

type WOLType string

const (
	Ethernet WOLType = "ethernet"
	UDP      WOLType = "udp"
)

const WOLEtherType = 0x0842

type WOL struct {
	Type      WOLType
	MAC       net.HardwareAddr
	IP        net.IP
	Port      int
	Interface string
}

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

	c, err := raw.ListenPacket(iface, WOLEtherType, nil)
	if err != nil {
		return err
	}
	defer c.Close()

	f := &ethernet.Frame{
		Destination: w.MAC,
		Source:      iface.HardwareAddr,
		EtherType:   WOLEtherType,
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

func (w *WOL) Send() error {
	switch w.Type {
	case Ethernet:
		return w.sendEthernet()
	case UDP:
		return w.sendUDP()
	}

	return nil
}
