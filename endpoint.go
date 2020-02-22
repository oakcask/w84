package w84

import (
	"strings"
)

// EndPoint is a pair of two strings which represent
// Network and Address.
// And also implements net.Addr so it can be passed to
// net module functions.
type EndPoint struct {
	Family, Address string
}

// Network returns address family described by
// the EndPoint.
func (ep EndPoint) Network() string {
	return ep.Family
}

// String returns network address described by
// the EndPoint.
func (ep EndPoint) String() string {
	return ep.Address
}

var addrFamilyPrefixes = []string{
	"tcp:", "tcp4:", "tcp6:",
	"udp:", "udp4:", "udp6:",
	"ip:", "ip4:", "ip6:",
	"unix:", "unixgram:", "unixpacket:",
}

// ParseEndPoint parses a string and returns
// EndPoint. Note that this function does not validate
// address format or availability.
// This function never returns nil.
func ParseEndPoint(s string) *EndPoint {
	for _, prefix := range addrFamilyPrefixes {
		if strings.HasPrefix(s, prefix) {
			return &EndPoint{Family: prefix[:len(prefix)-1], Address: s[len(prefix):]}
		}
	}

	return &EndPoint{Family: "tcp", Address: s}
}
