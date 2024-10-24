package network

import (
	"errors"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vishvananda/netlink"
)

type NetlinkMock struct {
	mock.Mock
}

func (m *NetlinkMock) LinkAdd(link netlink.Link) error {
	args := m.Called(link)
	return args.Error(0)
}

func (m *NetlinkMock) ParseAddr(s string) (*netlink.Addr, error) {
	args := m.Called(s)
	return args.Get(0).(*netlink.Addr), args.Error(1)
}

func (m *NetlinkMock) AddrAdd(link netlink.Link, addr *netlink.Addr) error {
	args := m.Called(link, addr)
	return args.Error(0)
}

func (m *NetlinkMock) LinkSetUp(link netlink.Link) error {
	args := m.Called(link)
	return args.Error(0)
}

func (m *NetlinkMock) InterfaceByName(name string) (*net.Interface, error) {
	args := m.Called(name)
	iface := &net.Interface{Name: name} // Return a dummy interface
	return iface, args.Error(0)
}

func SetupMocks() *NetlinkMock {
	// Create and configure the mock object
	mockLink := new(NetlinkMock)

	// Mock what happens when the bridge is successfully created
	mockLink.On("LinkAdd", mock.Anything).Return(nil)
	mockLink.On("ParseAddr", bridgeAddr).Return(&netlink.Addr{}, nil)
	mockLink.On("AddrAdd", mock.Anything, mock.Anything).Return(nil)
	mockLink.On("LinkSetUp", mock.Anything).Return(nil)
	mockLink.On("InterfaceByName", bridgeName).Return(nil, errors.New("interface not found"))

	return mockLink
}

func TestCreateBridge(t *testing.T) {
	netlinkMock := SetupMocks() // Get the mocked instance

	bridge := NewBridge() // Create a new Bridge object

	// Test bridge creation without issues when the bridge does *not* already exist.
	err := bridge.Create()

	// Check the bridge creation succeeded with no errors
	assert.NoError(t, err)

	// Now, verify that the mock methods were called as expected
	netlinkMock.AssertCalled(t, "InterfaceByName", bridgeName)
	netlinkMock.AssertCalled(t, "LinkAdd", mock.Anything)
	netlinkMock.AssertCalled(t, "ParseAddr", bridgeAddr)
	netlinkMock.AssertCalled(t, "AddrAdd", mock.Anything, mock.Anything)
	netlinkMock.AssertCalled(t, "LinkSetUp", mock.Anything)
}

func TestCreateBridge_AlreadyExists(t *testing.T) {
	// We simulate that the interface already exists by not returning an error from `InterfaceByName`.
	mockLink := new(NetlinkMock)
	mockLink.On("InterfaceByName", bridgeName).Return(&net.Interface{Name: bridgeName}, nil) // no error means it exists

	bridge := NewBridge()

	// Call the function and expect it to detect the bridge exists and return early with no error.
	err := bridge.Create()

	// Expect no error and that the bridge creation exited early.
	assert.NoError(t, err)

	// Ensure the LinkAdd, AddrAdd, and LinkSetUp methods were NOT called since the bridge exists.
	mockLink.AssertNotCalled(t, "LinkAdd", mock.Anything)
	mockLink.AssertNotCalled(t, "AddrAdd", mock.Anything, mock.Anything)
	mockLink.AssertNotCalled(t, "LinkSetUp", mock.Anything)
}
