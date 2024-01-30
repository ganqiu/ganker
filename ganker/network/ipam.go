package network

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	json "github.com/goccy/go-json"
)

type IPAM struct {
	SubnetAllocator string            `json:"subnet_allocator"`
	Subnets         map[string]string `json:"subnets"` //map of subnet and its allocated ip ,key is subnet, value is gateway ip(Bitmap string)
}

// load ipam from config file
func (ipam *IPAM) load() error {

	// check if ipamAllocatorPath exists,if not exists
	if _, err := os.Stat(ipam.SubnetAllocator); err != nil {
		return err
	}
	configFile, err := os.Open(ipam.SubnetAllocator)
	if err != nil {
		return err
	}
	defer configFile.Close()

	if err := json.NewDecoder(configFile).Decode(&ipam.Subnets); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// dump store ipam to config file
func (ipam *IPAM) dump() error {
	if err := os.MkdirAll(ipam.SubnetAllocator, 0644); err != nil {
		return err
	}
	configFile, err := os.OpenFile(ipam.SubnetAllocator, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer configFile.Close()

	if err := json.NewEncoder(configFile).Encode(ipam.Subnets); err != nil {
		return err
	}
	return nil

}

// Allocate allocate a ip from subnet
func (ipam *IPAM) Allocate(subnet *net.IPNet) (ip net.IP, err error) {
	ipam.Subnets = make(map[string]string)

	// load allocated ip from config file
	if err := ipam.load(); err != nil {
		return nil, fmt.Errorf("load ipam error: %v", err)
	}

	// regenerate subnet instance
	_, subnet, _ = net.ParseCIDR(subnet.String())

	// Size returns the number of leading ones and total bits in the mask
	// leading one means the number of 1 in the mask(as subnet identification length )
	prefix, size := subnet.Mask.Size()
	ipAddr := subnet.String()
	// if subnet is not allocated, init it
	if _, exist := ipam.Subnets[ipAddr]; !exist {
		// 1 << uint8(size-prefix) equals to 2 ^ (size-prefix)
		// it means the number of ip can be allocated
		ipam.Subnets[ipAddr] = strings.Repeat("0", 1<<uint8(size-prefix))
	}
	var AllocatedIP net.IP
	for c := range ipam.Subnets[ipAddr] {
		if ipam.Subnets[ipAddr][c] == '0' {
			// set ip allocated
			ipalloc := []byte(ipam.Subnets[ipAddr]) //as string can not be modified, so we need to convert it to byte array
			ipalloc[c] = '1'
			ipam.Subnets[ipAddr] = string(ipalloc)
			// get subnet IP address,like 192.168.1.1
			subnetIP := subnet.IP
			// as net.IP is an array of uint(4), like [192,168,1,0]
			// so we can get it through:
			for t := uint(4); t > 0; t-- {
				[]byte(subnetIP)[4-t] += uint8(c >> ((t - 1) * 8))
			}

			// as we have allocated ip from 1, so we need to add 1
			subnetIP[3] += 1
			AllocatedIP = subnetIP
			break
		}
	}

	// store new subnet to config file
	if err := ipam.dump(); err != nil {
		return nil, fmt.Errorf("dump ipam error: %v", err)
	}

	return AllocatedIP, nil
}

// Release release a ip from subnet
func (ipam *IPAM) Release(subnet *net.IPNet, ipaddr *net.IP) error {
	ipam.Subnets = make(map[string]string)
	_, subnet, _ = net.ParseCIDR(subnet.String())

	// load allocated ip from config file
	if err := ipam.load(); err != nil {
		return fmt.Errorf("load ipam error: %v", err)
	}

	// get subnet IP address by index
	c := 0
	releaseIP := ipaddr.To4()

	// as IP is allocated from 1, so we need to minus 1
	releaseIP[3] -= 1
	for t := uint(4); t > 0; t-- {
		// To reverse the process of obtaining an index and assigning IP addresses: subtract each corresponding digit of the IP addresses,
		// then left-shift the resulting values and add them to the index
		c += int(releaseIP[t-1]-subnet.IP[t-1]) << ((4 - t) * 8)
	}

	// set index to 0
	ipalloc := []byte(ipam.Subnets[subnet.String()])
	ipalloc[c] = '0'
	ipam.Subnets[subnet.String()] = string(ipalloc)
	// save new subnet to config file
	if err := ipam.dump(); err != nil {
		return fmt.Errorf("dump ipam error: %v", err)
	}
	return nil
}
