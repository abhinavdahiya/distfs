package main

import (
	"errors"
	"net"
	"os"

	"github.com/distfs/disk"
)

var (
	ErrIPNotFound = "IP address for the node not found"
)

//Defines he heatbeat payload
//Host info:
// Space <- Storage struct
type Heartbeat struct {
	IPAddr string
	Host   string
	Space  Storage
}

func Create() (*Heartbeat, error) {
	addr, err := GetLocalAddr()
	if err != nil {
		return nil, errors.Wrap(err, ErrIPNotFound)
	}

	space, err := GetSpace()
	if err != nil {
		return nil, err
	}

	host, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	return &Heartbeat{
		IPAddr: addr.String(),
		Host:   host,
		Space:  space,
	}, nil
}

//Defines disk space of host
//Payload in heartbeat
type Storage struct {
	Remaining int64
	Total     int64
}

//Get disk space
func GetSpace() (Storage, error) {
	total, free, err := disk.Space()
	if err != nil {
		return Storage{}, err
	}
	return Storage{free, total}, nil
}

//Get local IPaddress for the node
func GetLocalAddr() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, i := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := i.Addrs()
		if err != nil {
			return nil, err
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip, nil
		}
	}
	return nil, errors.New("Not connected to a network!")
}
