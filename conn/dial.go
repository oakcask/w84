package conn

import (
	"context"
	"net"
)

// Dial calls net.DialContext and immediatelly close connection,
// then returns nil on success,
// or an error if DialContext generated error.
func Dial(ctx context.Context, addr net.Addr) error {
	d := net.Dialer{}

	conn, e := d.DialContext(ctx, addr.Network(), addr.String())
	if conn != nil {
		defer closeConn(conn)
	}

	return e
}

// closeConn closes connection and ignore error.
func closeConn(conn net.Conn) {
	_ = conn.Close()
}
