package tester

import (
	"context"
	"errors"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/oakcask/stand"
	"github.com/oakcask/w84"
)

func TestRunSuccessful(t *testing.T) {
	dialledAddrsChan := make(chan net.Addr)

	config := w84.Config{
		Timeout: time.Duration(1) * time.Second,
		Clock:   stand.Pause(stand.SystemClock),
		DialFunc: func(ctx context.Context, addr net.Addr) error {
			select {
			case <-ctx.Done():
				return fmt.Errorf("test was timed out unexpectedly")
			default:
				dialledAddrsChan <- addr
				return nil // which means success
			}
		},
	}

	addrs := []net.Addr{
		w84.ParseEndPoint("example.com:22"),
		w84.ParseEndPoint("api.example.com:80"),
	}

	// make sure test not be hangged up.
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(1)*time.Second)
	dialledAddrs := make([]net.Addr, 0)
	collectorDone := make(chan bool)
	go (func() {
	Collect:
		for {
			select {
			case <-ctx.Done():
				// this will cause panic if program is bugged and infinitely calling "DialFunc" above.
				// but it's ok. test will fail faster.
				close(dialledAddrsChan)
				break Collect
			case item, ok := <-dialledAddrsChan:
				if !ok {
					break Collect
				}
				dialledAddrs = append(dialledAddrs, item)
			}
		}
		collectorDone <- true
	})()

	actualReports := Run(ctx, config, addrs)
	cancel()
	<-collectorDone

	if len(actualReports) != len(addrs) {
		t.Errorf("got %v of length %d while expecting %d in length", actualReports, len(actualReports), len(addrs))
	}

	for i, r := range actualReports {
		e := r.Err()
		if e != nil {
			t.Errorf("got %v while expecting nil as Err() of element at %d of return value", e, i)
		}
		addr := r.Addr()
		if addr.Network() != addrs[i].Network() || addr.String() != addrs[i].String() {
			t.Errorf("got %v while expecting %v as Addr() of element at %d of return value", addr, addrs[i], i)
		}
		updated := r.Updated()
		if updated != config.Clock.Now() {
			t.Errorf("got %v while expecting %v as Updated() of element at %d of return value", updated, config.Clock.Now(), i)
		}
	}

	for _, inAddr := range addrs {
		found := false
		for _, dialledAddr := range dialledAddrs {
			if dialledAddr.Network() == inAddr.Network() && dialledAddr.String() == inAddr.String() {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("%v not found in actual passed addrs: %v", inAddr, dialledAddrs)
		}
	}
}

func TestRunFailure(t *testing.T) {
	dialledAddrsChan := make(chan net.Addr)
	simulatedError := fmt.Errorf("simulated connection error")

	config := w84.Config{
		Timeout: time.Duration(1) * time.Second,
		Clock:   stand.Pause(stand.SystemClock),
		DialFunc: func(ctx context.Context, addr net.Addr) error {
			dialledAddrsChan <- addr
			select {
			case <-ctx.Done():
				return simulatedError
			default:
				return simulatedError
			}
		},
	}

	addrs := []net.Addr{
		w84.ParseEndPoint("example.com:22"),
		w84.ParseEndPoint("api.example.com:80"),
	}

	// how do we test functions which may wait until deadline?
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(1)*time.Second)
	dialledAddrs := make([]net.Addr, 0)
	collectorDone := make(chan bool)
	go (func() {
	Collect:
		// this will prevent infinite loop of calling "DialFunc"
		for len(dialledAddrs) < 6 {
			select {
			case <-ctx.Done():
				// this will cause panic if program is bugged and infinitely calling "DialFunc" above.
				// but it's ok. test will fail faster.
				close(dialledAddrsChan)
				break Collect
			case item, ok := <-dialledAddrsChan:
				if !ok {
					break Collect
				}
				dialledAddrs = append(dialledAddrs, item)
			}
		}
		collectorDone <- true
	})()

	actualReportsChan := make(chan []w84.Report)
	go (func() {
		actualReportsChan <- Run(ctx, config, addrs)
	})()

	<-collectorDone
	cancel()

	actualReports := <-actualReportsChan

	if len(actualReports) != len(addrs) {
		t.Errorf("got %v of length %d while expecting %d in length", actualReports, len(actualReports), len(addrs))
	}

	for i, r := range actualReports {
		e := r.Err()
		if !errors.Is(e, simulatedError) {
			t.Errorf("got %v while expecting %v as Err() of element at %d of return value", e, simulatedError, i)
		}
		addr := r.Addr()
		if addr.Network() != addrs[i].Network() || addr.String() != addrs[i].String() {
			t.Errorf("got %v while expecting %v as Addr() of element at %d of return value", addr, addrs[i], i)
		}
		updated := r.Updated()
		if updated != config.Clock.Now() {
			t.Errorf("got %v while expecting %v as Updated() of element at %d of return value", updated, config.Clock.Now(), i)
		}
	}

	for _, inAddr := range addrs {
		found := false
		for _, dialledAddr := range dialledAddrs {
			if dialledAddr.Network() == inAddr.Network() && dialledAddr.String() == inAddr.String() {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("%v not found in actual passed addrs: %v", inAddr, dialledAddrs)
		}
	}
}
