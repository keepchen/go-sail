package utils

import "net"

type ipImpl struct {
}

type IIP interface {
	// GetLocal retrieves the local IP address (unicast address)
	GetLocal() (string, error)
}

var ipi IIP = &ipImpl{}

// IP instantiates IP utilities
func IP() IIP {
	return ipi
}

// GetLocal retrieves the local IP address (unicast address)
func (ipImpl) GetLocal() (string, error) {
	addrList, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrList {
		if ipNet, ok := addr.(*net.IPNet); ok && ipNet.IP.IsGlobalUnicast() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String(), nil
			}
		}
	}

	return "", nil
}
