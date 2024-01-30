package container

import (
	"fmt"
	networks "go_docker_learning/ganker/network"
	"net"
	"os"
	"path/filepath"
	"text/tabwriter"
)

var (
	NetConfigRootPath = "./networks/netconfig/"
	netDriver         = map[string]networks.NetDriver{}
	network           = map[string]*networks.Net{}
)

const ipamAllocatorPath = "./networks/ipam.json"

var ipAllocator = &networks.IPAM{
	SubnetAllocator: ipamAllocatorPath,
}

// CreateNet create a net with the driver and subnet
func CreateNet(driver, subnet, name string) error {
	// parse subnet
	_, ipNet, err1 := net.ParseCIDR(subnet)
	if err1 != nil {
		return fmt.Errorf("parse subnet %s failed, err: %v", subnet, err1)
	}

	// Allocate ip for the subnet
	gatewayIp, err := ipAllocator.Allocate(ipNet)

	if err != nil {
		return fmt.Errorf("allocate ip for subnet %s failed, err: %v", subnet, err)
	}

	ipNet.IP = gatewayIp

	nw, err := netDriver[driver].Create(ipNet.String(), name)
	if err != nil {
		return fmt.Errorf("create network failed, err: %v", err)
	}

	return nw.Dump(NetConfigRootPath)

}

// Connect connect a container to the net
func Connect(containerId, netName string) error {
	info, err := getContainerInfo(containerId)
	if err != nil {
		return fmt.Errorf("get container info error: %v", err)
	}
	// get net from network map
	nw, ok := network[netName]
	if !ok {
		return fmt.Errorf("no such network: %s", netName)
	}

	// allocate ip for the container
	ip, err := ipAllocator.Allocate(nw.IpRange)
	if err != nil {
		return fmt.Errorf("allocate ip for subnet %s failed, err: %v", nw.IpRange.String(), err)
	}

	// construct netpoint
	netEndPoint := &networks.NetPoint{
		ID:          fmt.Sprintf("%s-%s", info.ContainerId, nw.Name),
		IP:          ip,
		PortMapping: info.PortMapping,
		Net:         nw,
	}

	if err := netDriver[nw.Driver].Connect(nw, netEndPoint); err != nil {
		return fmt.Errorf("connect network %s failed, err: %v", nw.Name, err)
	}

	if err := networks.ConfigEndpointIpAndRoute(netEndPoint, info.Pid); err != nil {
		return fmt.Errorf("config endpoint ip and route error: %v", err)
	}
	return networks.ConfigurePortMapping(netEndPoint)
}

// load all net config to network map
func InitNet() error {
	var bridgeDriver = networks.BridgeNetDriver{}
	netDriver[bridgeDriver.Name()] = &bridgeDriver

	if err := os.MkdirAll(NetConfigRootPath, 0644); err != nil {
		return fmt.Errorf("mkdir %s error: %v", NetConfigRootPath, err)
	}

	// check all net config file
	if err := filepath.Walk(NetConfigRootPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		// load fileName as net name
		_, fileName := filepath.Split(path)
		nw := &networks.Net{
			Name: fileName,
		}

		// call load function to load net config file
		if err := nw.Load(path); err != nil {
			return fmt.Errorf("load net config file %s error: %v", path, err)
		}

		// add net to network map
		network[nw.Name] = nw
		return nil
	}); err != nil {
		return err
	}
	return nil
}

// show all the nets
func ListNet() error {
	table := tabwriter.NewWriter(os.Stdout, 12, 1, 3, ' ', 0)
	fmt.Fprint(table, "NAME\tDRIVER\tSUBNET\tGATEWAY\n")

	for _, nw := range network {
		fmt.Fprintf(table, "%s\t%s\t%s\t%s\n", nw.Name, nw.Driver, nw.IpRange.String(), nw.IpRange.IP.String())
	}

	if err := table.Flush(); err != nil {
		return fmt.Errorf("flush table error: %v", err)
	}
	return nil
}
func DeleteNet(netName string) error {
	// check if the net exists
	nw, ok := network[netName]
	if !ok {
		return fmt.Errorf("no such network: %s", netName)
	}

	// Release ip for the subnet
	if err := ipAllocator.Release(nw.IpRange, &nw.IpRange.IP); err != nil {
		return fmt.Errorf("release ip for subnet %s failed, err: %v", nw.IpRange.String(), err)
	}

	// delete the net device and config file
	if err := netDriver[nw.Driver].Delete(nw); err != nil {
		return fmt.Errorf("delete network %s failed, err: %v", nw.Name, err)
	}

	return nw.Remove(NetConfigRootPath)
}
