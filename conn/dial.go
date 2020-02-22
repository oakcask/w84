package conn

import (
	"context"
	"net"
)

func Dial(ctx context.Context, addr net.Addr) error {
	d := net.Dialer{}

	conn, e := d.DialContext(ctx, addr.Network(), addr.String())
	if conn != nil {
		defer conn.Close()
	}

	return e
}
