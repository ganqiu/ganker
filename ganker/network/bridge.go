package network

import (
	"fmt"
	"net"
	"strings"

	"github.com/coreos/go-iptables/iptables"
	"github.com/vishvananda/netlink"
)

// BridgeNetDriver is a driver for bridge network
type BridgeNetDriver struct{}

func (d *BridgeNetDriver) Name() string {
	return "bridge"
}

// create Bridge network
func (d *BridgeNetDriver) Create(subnet, name string) (*Net, error) {
	// get subnet and gateway ip address
	ip, ipRange, err := net.ParseCIDR(subnet)
	if err != nil {
		return nil, fmt.Errorf("fail to parse subnet : %v", err)
	}

	ipRange.IP = ip
	// init network struct
	n := &Net{
		Name:    name,
		IpRange: ipRange,
		Driver:  d.Name(),
	}

	// init bridge network
	if err := d.initDriver(n); err != nil {
		return nil, fmt.Errorf("fail to init bridge network : %v", err)
	}

	return n, nil
}

// delete Bridge network
func (d *BridgeNetDriver) Delete(net *Net) error {
	bridgeName := net.Name

	// find bridge interface device
	interf, err := netlink.LinkByName(bridgeName)
	if err != nil {
		return fmt.Errorf("get bridge interface %v error : %v", bridgeName, err)
	}

	if err := netlink.LinkDel(interf); err != nil {
		return fmt.Errorf("delete bridge interface %v error : %v", bridgeName, err)
	}

	return err
}

// Connect create a veth pair and connect it to bridge interface
func (d *BridgeNetDriver) Connect(net *Net, endpoint *NetPoint) error {
	bridgeName := net.Name

	// find bridge interface device
	interf, err := netlink.LinkByName(bridgeName)
	if err != nil {
		return fmt.Errorf("get bridge interface %v error : %v", bridgeName, err)
	}

	// configure veth pair
	la := netlink.NewLinkAttrs()
	la.Name = endpoint.IDAddress[:5]
	// mount one end of veth pair to the bridge in network (as Master)
	la.MasterIndex = interf.Attrs().Index

	// create veth pair , one end of veth pair is mounted to the bridge in network(as Peer)
	endpoint.Device = netlink.Veth{
		LinkAttrs: la,
		PeerName:  endpoint.IDAddress[:5],
	}

	// create veth pair interface
	if err := netlink.LinkAdd(&endpoint.Device); err != nil {
		return fmt.Errorf("create veth pair interface error : %v", err)
	}

	// set veth pair interface up
	if err := netlink.LinkSetUp(&endpoint.Device); err != nil {
		return fmt.Errorf("set veth pair interface up error : %v", err)
	}

	return nil
}

func (d *BridgeNetDriver) Disconnect(net *Net, endpoint *NetPoint) error {

	if err := netlink.LinkDel(&endpoint.Device); err != nil {
		return fmt.Errorf("delete veth pair interface error : %v", err)
	}

	return nil
}

// init a bridge network
func (d *BridgeNetDriver) initDriver(n *Net) error {
	// create bridge interface
	brName := n.Name
	if err := d.createBrigdeInterface(brName); err != nil {
		return fmt.Errorf("fail to create bridge interface  :' %v'", err)
	}

	// set bridge interface address and router
	gateway := *n.IpRange
	gateway.IP = n.IpRange.IP

	// set bridge interface ip and router
	if err := setIP(brName, gateway.String()); err != nil {
		return fmt.Errorf("fail to set bridge interface ip : %v", err)
	}

	// set bridge interface up
	if err := setUpInterface(brName); err != nil {
		return fmt.Errorf("fail to set bridge interface up : %v", err)
	}

	// set iptables rule for bridge interface
	if err := setIPTableUp(brName, n.IpRange); err != nil {
		return fmt.Errorf("fail to set iptables : %v", err)
	}

	return nil
}

// createBridgeInterface creates a origin bridge interface
func (d *BridgeNetDriver) createBrigdeInterface(bridgeName string) error {

	// check if bridge interface exists
	interf, err := net.InterfaceByName(bridgeName)
	if interf != nil || err == nil {
		return fmt.Errorf("bridge interface already exists")
	}
	if !strings.Contains(err.Error(), "no such network interface") {
		return err
	}

	// create bridge interface
	link := netlink.NewLinkAttrs()
	link.Name = bridgeName
	bridge := &netlink.Bridge{LinkAttrs: link}

	// add bridge interface
	if err := netlink.LinkAdd(bridge); err != nil {
		return fmt.Errorf(" add bridge interface %v error : %v", bridgeName, err)
	}
	return nil
}

// set bridge interface ip
func setIP(name string, IP string) error {
	// get bridge interface that need to set ip
	interf, err := netlink.LinkByName(name)
	if err != nil {
		return fmt.Errorf("get bridge interface %v error : %v", name, err)
	}

	// ipNet include network segment and mask
	ipNet, err := netlink.ParseIPNet(IP)
	if err != nil {
		return fmt.Errorf("parse ip %v error : %v", IP, err)
	}

	addr := &netlink.Addr{
		IPNet: ipNet,
		Label: "",
		Flags: 0,
		Scope: 0,
	}

	//  if network segment where ip is located is
	if err := netlink.AddrAdd(interf, addr); err != nil {
		return fmt.Errorf("add ip %v to bridge interface %v error : %v", IP, name, err)
	}
	return nil

}

// setUp set bridge interface up
func setUpInterface(name string) error {
	interf, err := netlink.LinkByName(name)
	if err != nil {
		return fmt.Errorf("get bridge interface %v error : %v", name, err)
	}

	// set bridge interface up
	if err := netlink.LinkSetUp(interf); err != nil {
		return fmt.Errorf("set bridge interface %v up error : %v", name, err)
	}

	return nil
}

// set MASQUERADE rule for iptables in bridge mode
// MASQUERADE will replace the source ip of the packet with the ip of the NIC that the packet is sent from
func setIPTableUp(name string, subnet *net.IPNet) error {

	iptable, err := iptables.New()
	if err != nil {
		return fmt.Errorf("create iptables error : %v", err)
	}

	command := []string{"-s", subnet.String(), "! -o", name, "-j", "MASQUERADE"}
	iptable.Delete("nat", "POSTROUTING", command...)

	if err := iptable.Append("nat", "POSTROUTING", command...); err != nil {
		return fmt.Errorf("iptables append error : %v", err)
	}

	return nil
}
