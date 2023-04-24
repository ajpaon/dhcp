package server4

import (
	"fmt"
	"net"
)

// NewIPv4UDPConn fails on Windows. Use WithConn() to pass the connection.
func NewIPv4UDPConn(iface string, addr *net.UDPAddr) (*net.UDPConn, error) {
	conn, err := net.ListenUDP("udp4", addr)
	if err != nil {
		return nil, fmt.Errorf("unable to create UDP connection: %v", err)
	}
	return conn, nil
}
