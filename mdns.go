package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/micro/mdns"
)

func publishmDNS(ifname, port string, logger *log.Logger) (*mdns.Server, error) {
	host, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	info := []string{"Stream publisher service"}
	localIP, err := getLocalIP(ifname)
	if err != nil {
		return nil, err
	}

	iPort, err := strconv.Atoi(port)
	if err != nil {
		return nil, err
	}

	service, err := mdns.NewMDNSService(host, "stream_publisher._tcp", "", host+".", iPort, []net.IP{localIP}, info)
	if err != nil {
		return nil, err
	}

	// Create the mDNS server
	server, err := mdns.NewServer(&mdns.Config{Zone: service})
	if err != nil {
		return nil, err
	}

	return server, nil
}

func getLocalIP(ifname string) (net.IP, error) {
	iface, err := net.InterfaceByName(ifname)
	if err != nil {
		return nil, err
	}
	addrs, err := iface.Addrs()
	if err != nil {
		return nil, err
	}
	for _, addr := range addrs {
		ip, ok := addr.(*net.IPNet)
		if !ok {
			return nil, fmt.Errorf("ip: %v is not a net.IPNet", ip)
		}
		if ip.IP.To4() == nil { // this is an IPv6 address, ignore
			continue
		}
		return ip.IP, nil
	}
	return nil, errors.New("No ipv4 addresses found")
}
