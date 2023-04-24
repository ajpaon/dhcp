package interfaces

import "errors"

// BindToInterface fails on Windows.
func BindToInterface(fd int, ifname string) error {
	return errors.New("SO_BINDTODEVICE is not supported on windows")
}
