package network

import (
	"fmt"
	"net"
	"strings"
	"time"
)

func IsAny(ip string) bool {
	return isAnyV4(ip) || isAnyV6(ip)
}

func isAnyV4(ip string) bool { return ip == "0.0.0.0" }

func isAnyV6(ip string) bool { return ip == "::" || ip == "[::]" }

func DialUntilReachable(addr string, maxWait time.Duration) error {
	done := time.Now().Add(maxWait)
	for time.Now().Before(done) {
		c, err := net.Dial("tcp", addr)
		if err == nil {
			c.Close()
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("%v unreachable for %v", addr, maxWait)
}

func ContainsPortNumber(url string) bool {
	return strings.LastIndex(url, ":") > strings.LastIndex(url, "]")
}
