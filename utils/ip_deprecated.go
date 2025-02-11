package utils

import "net"

// GetLocalIP 获取本地ip地址（单播地址）
//
// Deprecated: GetLocalIP is deprecated,it will be removed in the future.
//
// Please use IP().GetLocal() instead.
func GetLocalIP() (string, error) {
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
