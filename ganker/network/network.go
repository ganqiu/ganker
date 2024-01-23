package network

import (
	"encoding/json"
	"fmt"
	"go_docker_learning/ganker/container"
	"net"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"text/tabwriter"

	"github.com/coreos/go-iptables/iptables"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
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
	ID          string           `json:"id"`
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

	return nw.dump(NetConfigRootPath)

}

// Connect connect a container to the net
func Connect(netName string, info *container.Info) error {
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
	netEndPoint := &NetPoint{
		ID:          fmt.Sprintf("%s-%s", info.ContainerId, nw.Name),
		IP:          ip,
		PortMapping: info.PortMapping,
		Net:         nw,
	}

	if err := netDriver[nw.Driver].Connect(nw, netEndPoint); err != nil {
		return fmt.Errorf("connect network %s failed, err: %v", nw.Name, err)
	}

	if err := configEndpointIpAndRoute(netEndPoint, info); err != nil {
		return fmt.Errorf("config endpoint ip and route error: %v", err)
	}
	return configurePortMapping(netEndPoint)
}

// load all net config to network map
func InitNet() error {
	var bridgeDriver = BridgeNetDriver{}
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
		nw := &Net{
			Name: fileName,
		}

		// call load function to load net config file
		if err := nw.load(path); err != nil {
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

	return nw.remove(NetConfigRootPath)
}

// configEndpointIpAndRoute config ip address and route for the endpoint
func configEndpointIpAndRoute(endpoint *NetPoint, info *container.Info) error {
	// get t
	veth, err := netlink.LinkByName(endpoint.Device.PeerName)
	if err != nil {
		return fmt.Errorf("get peer link %s error: %v", endpoint.Device.PeerName, err)
	}

	// set endpoint ip address to container net ns, after functioned, the process will quit the container net ns
	defer enterContainerNetns(&veth, info)()

	// get the container ip address and subnet
	interfaceIp := *endpoint.Net.IpRange
	interfaceIp.IP = endpoint.IP

	// set ip address to the container
	if err := setIP(endpoint.Device.PeerName, interfaceIp.String()); err != nil {
		return fmt.Errorf("set interface %s ip error: %v", endpoint.Device.PeerName, err)
	}

	// set up veth endpoint
	if err := setUpInterface(endpoint.Device.PeerName); err != nil {
		return fmt.Errorf("set interface %s up error: %v", endpoint.Device.PeerName, err)
	}

	// set default route to the container
	if err := setUpInterface("lo"); err != nil {
		return fmt.Errorf("set interface lo up error: %v", err)
	}

	// all out request should be sent to the gateway(veth endpoint)
	_, cidr, _ := net.ParseCIDR("0.0.0.0/0") // 0.0.0.0/0 means all ip address

	// construct net route
	route := netlink.Route{
		LinkIndex: veth.Attrs().Index,
		Gw:        endpoint.Net.IpRange.IP,
		Dst:       cidr,
	}

	if err := netlink.RouteAdd(&route); err != nil {
		return fmt.Errorf("set route error: %v", err)
	}

	return nil

}

// enterContainerNetns return a function pointer, when the function is called, the process will quit the container netns
func enterContainerNetns(link *netlink.Link, info *container.Info) func() {

	// we can get the container netns file from /proc/[pid]/ns/net
	function, err := os.OpenFile(fmt.Sprintf("/proc/%s/ns/net", info.Pid), os.O_RDONLY, 0)
	if err != nil {
		fmt.Printf("open container netns error: %v", err)
		return nil
	}

	// get the container netns file descriptor
	nsFD := function.Fd()

	// lock the container netns file
	runtime.LockOSThread()

	// move one end of the veth pair to the container netns
	if err := netlink.LinkSetNsFd(*link, int(nsFD)); err != nil {
		fmt.Printf("set link %s netns error: %v", (*link).Attrs().Name, err)
		return nil
	}

	// get the current namespace
	preNs, err := netns.Get()
	if err != nil {
		fmt.Printf("get current netns error: %v", err)
		return nil
	}

	// set the container netns as the current namespace(set current process to container net namespace)
	if err = netns.Set(netns.NsHandle(nsFD)); err != nil {
		fmt.Printf("set netns error: %v", err)
		return nil
	}

	// return current namespace function
	// in container net namespace,
	return func() {
		// set current namespace to original namespace
		netns.Set(preNs)
		// close the container netns file
		preNs.Close()
		// unlock the thread
		runtime.UnlockOSThread()
		// close namespace file
		function.Close()
	}
}

func configurePortMapping(endpoint *NetPoint) error {

	for _, pm := range endpoint.PortMapping {
		mapArray := strings.Split(pm, ":")
		if len(mapArray) != 2 {
			return fmt.Errorf("invalid port mapping: %s", pm)
		}

		iptable, err := iptables.New()
		if err != nil {
			return fmt.Errorf("create iptables error: %v", err)
		}

		// add DNAT rule to nat table
		command := []string{"-p", "tcp", "-m", "tcp", "--dport", mapArray[0], "-j", "DNAT", "--to-destination", fmt.Sprintf("%s:%s", endpoint.IP.String(), mapArray[1])}
		if err := iptable.Append("nat", "POSTROUTING", command...); err != nil {
			fmt.Printf("iptables append %v:%v port mapping error: %v", mapArray[0], mapArray[1], err)
			continue
		}
	}
	return nil
}

func (n *Net) dump(configPath string) error {

	if err := os.MkdirAll(configPath, 0644); err != nil {
		return err
	}

	file, err := os.OpenFile(path.Join(configPath, n.Name), os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("create config file %s error: %v", n.Name, err)
	}

	if err := json.NewEncoder(file).Encode(&n); err != nil {
		return fmt.Errorf("encode config file %s error: %v", n.Name, err)
	}
	return nil
}

func (n *Net) load(path string) error {

	configFile, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open config file %s error: %v", path, err)
	}
	defer configFile.Close()

	if err := json.NewDecoder(configFile).Decode(&n); err != nil {
		return fmt.Errorf("decode config file %s error: %v", path, err)
	}

	return nil
}

// remove remove the net config file
func (n *Net) remove(Path string) error {
	if _, err := os.Stat(Path); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return os.Remove(path.Join(Path, n.Name))
}
