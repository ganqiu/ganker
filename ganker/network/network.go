package network

import (
	"fmt"
	"net"
	"os"
	"path"
	"runtime"
	"strings"

	json "github.com/goccy/go-json"

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

// configEndpointIpAndRoute config ip address and route for the endpoint
func ConfigEndpointIpAndRoute(endpoint *NetPoint, pid string) error {
	// get t
	veth, err := netlink.LinkByName(endpoint.Device.PeerName)
	if err != nil {
		return fmt.Errorf("get peer link %s error: %v", endpoint.Device.PeerName, err)
	}

	// set endpoint ip address to container net ns, after functioned, the process will quit the container net ns
	defer enterContainerNetns(&veth, pid)()

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
func enterContainerNetns(link *netlink.Link, pid string) func() {

	// we can get the container netns file from /proc/[pid]/ns/net
	function, err := os.OpenFile(fmt.Sprintf("/proc/%s/ns/net", pid), os.O_RDONLY, 0)
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

func ConfigurePortMapping(endpoint *NetPoint) error {

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

func (n *Net) Dump(configPath string) error {

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

func (n *Net) Load(path string) error {

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
func (n *Net) Remove(Path string) error {
	if _, err := os.Stat(Path); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return os.Remove(path.Join(Path, n.Name))
}
