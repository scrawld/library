package util

import (
	"fmt"
	"net/netip"
)

// ParseIPPrefix 解析 IP 字符串，如果是 IPv4 返回原地址，如果是 IPv6 返回前 64 位地址
func ParseIPPrefix(ipStr string) (string, error) {
	ip, err := netip.ParseAddr(ipStr)
	if err != nil {
		return "", fmt.Errorf("parse addr error: %s", err)
	}
	if ip.Is4() {
		return ip.String(), nil
	}
	prefix := netip.PrefixFrom(ip, 64)
	return prefix.Masked().Addr().String(), nil
}
