package utils

import "net"

type ipImpl struct {
}

type IIP interface {
	// GetLocal 获取本地ip地址（单播地址）
	GetLocal() (string, error)
}

var ipi IIP = &ipImpl{}

// IP 实例化ip工具类
func IP() IIP {
	return ipi
}

// GetLocal 获取本地ip地址（单播地址）
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
