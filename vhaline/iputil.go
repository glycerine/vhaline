package vhaline

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

// getAvailPort asks the OS for an unused port,
// returning a bound net.Listener and the port number
// to which it is bound. The caller should
// Close() the listener when it is done with
// the port.
func getAvailPort() (net.Listener, int) {
	lsn, err := net.Listen("tcp", ":0")
	panicOn(err)
	r := lsn.Addr()
	return lsn, r.(*net.TCPAddr).Port
}

// waitUntilAddrAvailable returns -1 if the addr was
// alays unavailable after tries sleeps of dur time.
// Otherwise it returns the number of tries it took.
// Between attempts we wait 'dur' time before trying
// again.
func waitUntilAddrAvailable(addr string, dur time.Duration, tries int) int {
	for i := 0; i < tries; i++ {
		var isbound bool
		isbound = IsAlreadyBound(addr)
		if isbound {
			time.Sleep(dur)
		} else {
			fmt.Printf("\n took %v %v sleeps for address '%v' to become available.\n", i, dur, addr)
			return i
		}
	}
	return -1
}

func IsAlreadyBound(addr string) bool {

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return true
	}
	ln.Close()
	return false
}

// AddrHostPort is used to specify tunnel endpoints.
type AddrHostPort struct {
	Title    string
	Addr     string
	Host     string
	Port     uint64
	Required bool
}

// ParseAddr fills Host and Port from Addr, breaking Addr apart at the ':'
// using net.SplitHostPort()
func (a *AddrHostPort) ParseAddr() error {

	if a.Addr == "" {
		if a.Required {
			return fmt.Errorf("provide -%s ip:port", a.Title)
		}
		return nil
	}

	host, port, err := net.SplitHostPort(a.Addr)
	if err != nil {
		return fmt.Errorf("bad -%s ip:port given; net.SplitHostPort() gave: %s", a.Title, err)
	}
	a.Host = host
	if host == "" {
		//p("defaulting empty host to 127.0.0.1")
		a.Host = "127.0.0.1"
	} else {
		//p("in ParseAddr(%s), host is '%v'", a.Title, host)
	}
	if len(port) == 0 {
		return fmt.Errorf("empty -%s port; no port found in '%s'", a.Title, a.Addr)
	}
	a.Port, err = strconv.ParseUint(port, 10, 16)
	if err != nil {
		return fmt.Errorf("bad -%s port given; could not convert "+
			"to integer: %s", a.Title, err)
	}
	return nil
}
