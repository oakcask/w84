package w84

import (
	"fmt"
)

func ExampleParseEndPoint() {
	endPoints := []string {
		"example.com",
		"tcp:192.168.192.168:80",
		"udp6:[::1]:51234",
		"unix:/var/run/docker.sock",
	}

	for _, s := range endPoints {
		ep := ParseEndPoint(s)
		fmt.Printf("%v => Network: %s, Address: %s\n", s, ep.Network(), ep.String())
	}

	// Output:
	// example.com => Network: tcp, Address: example.com
	// tcp:192.168.192.168:80 => Network: tcp, Address: 192.168.192.168:80
	// udp6:[::1]:51234 => Network: udp6, Address: [::1]:51234
	// unix:/var/run/docker.sock => Network: unix, Address: /var/run/docker.sock
}
