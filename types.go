package w84

import (
	"context"
	"net"
	"time"
	"github.com/oakcask/stand"
)

// DialFunc is signature of function that try dialing
// to address given as net.Addr. A context will
// given for setting deadline (or timeout) for the operation.
type DialFunc func(context.Context, net.Addr) error

// Report represents the result of connectivity test for
// one site where Addr() points at.
type Report interface {
	// Addr is a site to test connectivity.
	Addr() net.Addr
	// Err will be nil if the connectivity is good,
	// or non-nil if any error occured during test processs.
	Err() error
	// Updated is the time when the report was constructed.
	Updated() time.Time
}

// Config is connectivity test configuration.
type Config struct {
	// Timeout specifies total test deadline in time.Duration
	Timeout time.Duration
	// Clock is a time source which will be utilized to record time for Report.
	Clock stand.Clock
	// DialFunc is connecting method which used within the test process.
	DialFunc DialFunc
}

