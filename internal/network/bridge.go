package network

import (
	"log/slog"
	"net"

	"github.com/vishvananda/netlink"
)

const (
	bridgeName = "montainerd0"
	bridgeAddr = "172.25.0.0/16"
)

type Bridge struct {
	name  string
	addr  *netlink.Addr
	iface *net.Interface
}
type Bridger interface {
	Create() error
	Exists() bool
	Attach(hostVeth *net.Interface) error
}

func (i *Bridge) Create() error {
	if i.Exists() {
		slog.Info("Bridge with name montainerd0 already exists")
		return nil
	}
	linkAttrs := netlink.LinkAttrs{Name: i.name}
	link := &netlink.Bridge{LinkAttrs: linkAttrs}
	if err := netlink.LinkAdd(link); err != nil {
		slog.Error("Error creating bridge: ", slog.Any("error", err))
		return err
	}
	address, err := netlink.ParseAddr(bridgeAddr)
	if err != nil {
		slog.Error("Error parsing address: ", slog.Any("error", err))
		return err
	}
	if err := netlink.AddrAdd(link, address); err != nil {
		slog.Error("Error adding address to bridge: ", slog.Any("error", err))
		return err
	}
	if err := netlink.LinkSetUp(link); err != nil {
		slog.Error("Error setting up bridge: ", slog.Any("error", err))
		return err
	}

	i.iface, _ = net.InterfaceByName(i.name)
	i.addr = address
	return nil
}

func (i *Bridge) Exists() bool {
	return InterfaceExists(i.name)
}

func (i *Bridge) Attach(hostVeth *net.Interface) error {
	bridgeLink, err := netlink.LinkByName(i.name)
	if err != nil {
		slog.Error("Error getting bridge link: ", slog.Any("error", err))
		return err
	}
	hostVethLink, err := netlink.LinkByName(hostVeth.Name)
	if err != nil {
		slog.Error("Error getting host veth link: ", slog.Any("error", err))
		return err
	}
	return netlink.LinkSetMaster(hostVethLink, bridgeLink)
}

func NewBridge() *Bridge {
	return &Bridge{name: bridgeName}
}
