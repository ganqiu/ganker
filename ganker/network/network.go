package network

import (
	"go_docker_learning/ganker/container"
	"net"

	"github.com/vishvananda/netlink"
)

var (
	NetConfigRootPath = "./networks/netconfig/"
	netDriver         = map[string]NetDriver{}
	network           = map[string]*Net{}
)

// Net is a collection of conatiners, which can communicate with each other, like mounted on the linux bridge device
type Net struct {
	Name    string     // net name
	Driver  string     // driver name
	IpRange *net.IPNet // ip range
}

// Netpoint is a network endpoint in the net, which is used to connect container to the net
type NetPoint struct {
	IDAddress   string           `json:"id"`
	IP          net.IP           `json:"ip"`
	MACAddress  net.HardwareAddr `json:"mac"`
	Device      netlink.Veth     `json:"device"`
	PortMapping []string         `json:"portmapping"`
	Net         *Net             `json:"net"`
}

type NetDriver interface {
	Name() string                                    // return driver name
	Create(subnet string, name string) (*Net, error) // create a net
	Delete(net *Net) error                           // delete a net
	Connect(net *Net, endpoint *NetPoint) error      // connect a container to the net
	Disconnect(net *Net, endpoint *NetPoint) error   // disconnect a container from the net
}

func CreateNet(driver, subnet, name string) error {
	return nil
}
func Connect(netName string, info *container.Info) error {
	return nil
}
func InitNet() error {
	return nil
}
func ListNet() error {
	return nil
}
func DeleteNet(netName string) error {
	return nil
}

func configEndpointIpAndRoute(endpoint *NetPoint, info *container.Info) error {
	return nil
}

func enterContainerNetns(link *netlink.Link, info *container.Info) error {
	return nil
}

func configurePortMapping(endpoint *NetPoint) error {
	return nil
}

func (n *Net) dump() error {
	return nil
}

func (n *Net) load() error {
	return nil
}
