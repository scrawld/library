package util

import (
	"fmt"
	"net"
)

// ParseIPPrefix 解析 IP 字符串，如果是 IPv4 返回原地址，如果是 IPv6 返回 /64 前缀
func ParseIPPrefix(ipStr string) (string, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return "", fmt.Errorf("invalid IP: %s", ipStr)
	}

	// IPv4 直接返回
	if ip4 := ip.To4(); ip4 != nil {
		return ip4.String(), nil
	}

	// IPv6 返回 /64 前缀
	mask := net.CIDRMask(64, 128)
	pref := ip.Mask(mask)
	return pref.String() + "/64", nil
}
