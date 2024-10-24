package network

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/vishvananda/netlink"
)

type Veth struct {
	namePrefix        string
	hostVeth          *net.Interface
	containerVeth     *net.Interface
	hostVethName      string
	containerVethName string
}

type Veather interface {
	Create() (*net.Interface, *net.Interface, error)
	GetInterfaceByName(hostVethName, containerVethName string) (*net.Interface, *net.Interface, error)
	MoveToNamespace(pid int) error
}

func (v *Veth) Create() (*net.Interface, *net.Interface, error) {
	hostVethName := fmt.Sprintf("%s0", v.namePrefix)
	containerVethName := fmt.Sprintf("%s1", v.namePrefix)
	if InterfaceExists(hostVethName) {
		slog.Info("Veth with name ", slog.Any("name", hostVethName), " already exists")
		return v.GetInterfaceByName(hostVethName, containerVethName)
	}
	vethLinkAttrs := netlink.LinkAttrs{Name: hostVethName}
	veth := &netlink.Veth{
		LinkAttrs: vethLinkAttrs,
		PeerName:  containerVethName,
	}
	if err := netlink.LinkAdd(veth); err != nil {
		slog.Error("Error creating veth: ", slog.Any("error", err))
		return nil, nil, err
	}
	if err := netlink.LinkSetUp(veth); err != nil {
		slog.Error("Error setting up veth: ", slog.Any("error", err))
		return nil, nil, err
	}

	return v.GetInterfaceByName(hostVethName, containerVethName)
}

func (v *Veth) GetInterfaceByName(hostVethName, containerVethName string) (*net.Interface, *net.Interface, error) {
	hostVethInterface, err := net.InterfaceByName(hostVethName)
	if err != nil {
		return nil, nil, err
	}
	v.hostVeth = hostVethInterface
	containerVethInterface, err := net.InterfaceByName(containerVethName)
	if err != nil {
		return nil, nil, err
	}
	v.containerVeth = containerVethInterface
	return hostVethInterface, containerVethInterface, nil
}
