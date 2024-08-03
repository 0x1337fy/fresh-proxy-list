package config

import (
	"net"
)

var PrivateIPs = []net.IPNet{
	*ParseCIDR("10.0.0.0/8"),
	*ParseCIDR("172.16.0.0/12"),
	*ParseCIDR("192.168.0.0/16"),
	*ParseCIDR("169.254.0.0/16"), // link-local
	*ParseCIDR("240.0.0.0/4"),    // reserved for special use
	*ParseCIDR("224.0.0.0/4"),    // multicast
}

func ParseCIDR(cidr string) *net.IPNet {
	_, netIP, err := net.ParseCIDR(cidr)
	if err != nil {
		panic(err)
	}
	return netIP
}
