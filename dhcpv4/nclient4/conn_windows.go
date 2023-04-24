package nclient4

import (
	"errors"
	"net"
)

func NewRawUDPConn(iface string, port int) (net.PacketConn, error) {
	return nil, errors.New("raw UDP connection not supported on windows")
}
