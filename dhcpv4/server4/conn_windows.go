package server4

import (
	"errors"
	"fmt"
	"net"
	"os"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"golang.org/x/sys/windows"
)

// NewIPv4UDPConn fails on Windows. Use WithConn() to pass the connection.
func NewIPv4UDPConn(iface string, addr *net.UDPAddr) (*net.UDPConn, error) {
	fd, err := windows.Socket(windows.AF_INET, windows.SOCK_DGRAM, windows.IPPROTO_UDP)
	if err != nil {
		return nil, fmt.Errorf("cannot get a UDP socket: %v", err)
	}
	f := os.NewFile(uintptr(fd), "")
	// net.FilePacketConn dups the FD, so we have to close this in any case.
	defer f.Close()

	// Allow broadcasting.
	if err := windows.SetsockoptInt(fd, windows.SOL_SOCKET, windows.SO_BROADCAST, 1); err != nil {
		return nil, fmt.Errorf("cannot set broadcasting on socket: %v", err)
	}
	// Allow reusing the addr to aid debugging.
	if err := windows.SetsockoptInt(fd, windows.SOL_SOCKET, windows.SO_REUSEADDR, 1); err != nil {
		return nil, fmt.Errorf("cannot set reuseaddr on socket: %v", err)
	}

	if addr == nil {
		addr = &net.UDPAddr{Port: dhcpv4.ServerPort}
	}
	// Bind to the port.
	saddr := windows.SockaddrInet4{Port: addr.Port}
	if addr.IP != nil && addr.IP.To4() == nil {
		return nil, fmt.Errorf("wrong address family (expected v4) for %s", addr.IP)
	}
	copy(saddr.Addr[:], addr.IP.To4())
	if err := windows.Bind(fd, &saddr); err != nil {
		return nil, fmt.Errorf("cannot bind to port %d: %v", addr.Port, err)
	}

	conn, err := net.FilePacketConn(f)
	if err != nil {
		return nil, err
	}
	udpconn, ok := conn.(*net.UDPConn)
	if !ok {
		return nil, errors.New("BUG(dhcp4): incorrect socket type, expected UDP")
	}
	return udpconn, nil
}
