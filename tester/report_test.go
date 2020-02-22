package tester

import (
	"time"
	"testing"
	"errors"
	"github.com/oakcask/w84"
)

var testReportImplExamples = []struct {
	in report
} {
	{
		report{
			addr: &w84.EndPoint{"tcp", "example.com:80"},
			err: errors.New("Connection failed."),
			updated: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	},
}

func TestReportImpl(t *testing.T) {
	for _, tt := range testReportImplExamples {
		if tt.in.Addr() != tt.in.addr {
			t.Errorf("%v: got %v while expect %v for Addr()", tt.in, tt.in.Addr(), tt.in.addr)
		}
		if tt.in.Err() != tt.in.err {
			t.Errorf("%v: got %v while expect %v for Err()", tt.in, tt.in.Err(), tt.in.err)
		}
		if tt.in.Updated() != tt.in.updated {
			t.Errorf("%v: got %v while expect %v for Updated()", tt.in, tt.in.Updated(), tt.in.updated)
		}
	}
}
