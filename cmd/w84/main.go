package main

import (
	"context"
	"flag"
	"fmt"
	"time"
	"os"
	"net"

	"github.com/oakcask/w84"
	"github.com/oakcask/w84/conn"
	"github.com/oakcask/stand"
	"github.com/oakcask/w84/tester"
)

const (
	exitCodeOk = 0,
	exitCodeConnectivityTestFailure = 1,
	exitCodeNoArguments = 2,
)

var timeoutDefault = time.Duration(15)*time.Second
var timeout = flag.Duration("timeout", timeoutDefault, "Timeout duration to wait for the connection.")
var verboseDefault = false
var verbose = flag.Bool("verbose", verboseDefault, "Verbose: report more about test result.")

func init() {
	flag.DurationVar(timeout, "t", timeoutDefault, "alias to timeout")
	flag.BoolVar(verbose, "v", verboseDefault, "alias to verbose")
	flag.Parse()
}


func main() {
	os.Exit(programMain())
}

func programMain() int {
	if flag.NArg() < 1 {
		flag.PrintDefaults()
		return exitCodeNoArguments
	}

	endPoints := parseEndPoints(flag.Args())

	ctx := context.Background()
	config := w84.Config{
		Timeout: *timeout,
		Clock: stand.SystemClock,
		DialFunc: conn.Dial,
	}

	reports := tester.Run(ctx, config, endPoints)

	hasError := false
	for _, r := range reports {
		printReport(r)
		if r.Err() != nil {
			hasError = true
		}
	}

	if hasError {
		return exitCodeConnectivityTestFailure
	}
	return exitCodeOk
}

func parseEndPoints(args []string) []net.Addr {
	addrs := make([]net.Addr, len(args))
	for i, e := range args {
		addrs[i] = w84.ParseEndPoint(e)
	}

	return addrs
}

func printReport(r w84.Report) {
	if !*verbose {
		return
	}
	e := r.Err()
	if e == nil {
		fmt.Printf("%v %v -- %v", r.Updated().Format(time.RFC3339), r.Addr(), "OK")
	} else {
		fmt.Printf("%v %v -- %v", r.Updated().Format(time.RFC3339), r.Addr(), r.Err())
	}
}