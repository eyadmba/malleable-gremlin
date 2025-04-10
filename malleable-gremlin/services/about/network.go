package about

import (
	"net"
	// "strings"
	// "github.com/shirou/gopsutil/v3/net"
)

type NetworkInfo struct {
	Interfaces []InterfaceInfo `json:"interfaces"`
}

type InterfaceInfo struct {
	Name         string   `json:"name"`
	HardwareAddr string   `json:"hardware_addr"`
	Addresses    []string `json:"addresses"`
	Flags        []string `json:"flags"`
	MTU          int      `json:"mtu"`
}

func GetNetworkInfo() (*NetworkInfo, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	netInfo := &NetworkInfo{
		Interfaces: make([]InterfaceInfo, 0, len(interfaces)),
	}

	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		addrStrings := make([]string, 0, len(addrs))
		for _, addr := range addrs {
			addrStrings = append(addrStrings, addr.String())
		}

		flags := make([]string, 0)
		if iface.Flags&net.FlagUp != 0 {
			flags = append(flags, "up")
		}
		if iface.Flags&net.FlagLoopback != 0 {
			flags = append(flags, "loopback")
		}
		if iface.Flags&net.FlagBroadcast != 0 {
			flags = append(flags, "broadcast")
		}
		if iface.Flags&net.FlagMulticast != 0 {
			flags = append(flags, "multicast")
		}

		netInfo.Interfaces = append(netInfo.Interfaces, InterfaceInfo{
			Name:         iface.Name,
			HardwareAddr: iface.HardwareAddr.String(),
			Addresses:    addrStrings,
			Flags:        flags,
			MTU:          iface.MTU,
		})
	}

	return netInfo, nil
}
