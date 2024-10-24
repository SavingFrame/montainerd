package network

import "net"

func InterfaceExists(name string) bool {
	_, err := net.InterfaceByName(name)
	return err == nil
}
