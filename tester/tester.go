package tester

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/oakcask/w84"
)

var errNotStarted = errors.New("timeout exceeded (perhaps you set timeout too short. hou about making it longer)")

func notStarted(addr net.Addr, time time.Time) w84.Report {
	return failed(addr, time, errNotStarted)
}

func succeeded(addr net.Addr, time time.Time) w84.Report {
	return &report{
		addr:    addr,
		err:     nil,
		updated: time,
	}
}

func failed(addr net.Addr, time time.Time, e error) w84.Report {
	return &report{
		addr:    addr,
		err:     e,
		updated: time,
	}
}

type singleResult struct {
	ticket int
	report w84.Report
}

func runSingle(ctx context.Context, ticket int, config w84.Config, addr net.Addr, ch chan<- singleResult) {
	e := config.DialFunc(ctx, addr)
	if e == nil {
		ch <- singleResult{
			ticket: ticket,
			report: succeeded(addr, config.Clock.Now()),
		}
		return
	}
	ch <- singleResult{
		ticket: ticket,
		report: failed(addr, config.Clock.Now(), e),
	}
}

// Run invokes connectivity tests.
func Run(ctx context.Context, config w84.Config, addrs []net.Addr) []w84.Report {
	reports := make([]w84.Report, len(addrs))
	now := config.Clock.Now()

	for idx, addr := range addrs {
		reports[idx] = notStarted(addr, now)
	}

	ch := make(chan singleResult)

	ctx1, cancel := context.WithTimeout(ctx, config.Timeout)
	defer cancel()

	for idx, addr := range addrs {
		go runSingle(ctx, idx, config, addr, ch)
	}

	done := 0

Wait:
	for done < len(addrs) {
		select {
		case r := <-ch:
			reports[r.ticket] = r.report
			if r.report.Err() == nil {
				done++
			}
		case <-ctx1.Done():
			break Wait
		}
	}

	return reports
}
