package conn

import (
	"context"
	"net"
	"testing"
	"time"

	"golang.org/x/net/nettest"
)

func TestDialSuccesful(t *testing.T) {
	networks := []string{"tcp", "tcp4", "tcp6", "unix", "unixpacket"}
	for _, network := range networks {
		listener, err := nettest.NewLocalListener(network)
		if err != nil {
			t.Logf("Skipping %s test", network)
			continue
		}
		defer closeListener(listener)

		ch := make(chan error)
		defer close(ch)

		addr := listener.Addr()
		// Timeout is for just makeing test die smoothly.
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(1)*time.Second)
		defer cancel()

		go (func() {
			ch <- Dial(ctx, addr)
		})()
		go (func() {
			if conn, _ := listener.Accept(); conn != nil {
				closeConn(conn)
			}
		})()

		select {
		case <-ctx.Done():
			t.Errorf("Dial does not returned within timeout.")
		case actual := <-ch:
			if actual != nil {
				t.Errorf("got %v while expecting nil.", actual)
			}
		}
	}
}

func TestDialTimeout(t *testing.T) {
	networks := []string{"tcp", "tcp4", "tcp6", "unix", "unixpacket"}
	for _, network := range networks {
		listener, err := nettest.NewLocalListener(network)
		if err != nil {
			t.Logf("Skipping %s test", network)
			continue
		}
		defer closeListener(listener)

		ch := make(chan error)
		defer close(ch)

		addr := listener.Addr()
		// Preparing "dead" context
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		go (func() {
			ch <- Dial(ctx, addr)
		})()
		go (func() {
			if conn, _ := listener.Accept(); conn != nil {
				closeConn(conn)
			}
		})()

		select {
		case <-time.After(time.Duration(1) * time.Second):
			t.Errorf("Dial does not returned within timeout.")
		case actual := <-ch:
			if actual == nil {
				t.Errorf("got nil while expecting non-nil.")
			}
		}
	}
}

// closeListener closes listener and ignore error.
func closeListener(listener net.Listener) {
	_ = listener.Close()
}
