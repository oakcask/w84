package tester

import (
	"net"
	"time"
)

type report struct {
	addr net.Addr
	err error
	updated time.Time
}

func (r *report) Addr() net.Addr {
	return r.addr
}

func (r *report) Err() error {
	return r.err
}

func (r *report) Updated() time.Time {
	return r.updated
}
