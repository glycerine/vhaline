package vhaline

import (
	"net"
	"time"
)

type Cfg struct {

	// Addr is the host:port address to bind as server and listen on.
	Addr string

	// if not nil, the Lsn is aready bound to Addr.
	Lsn net.Listener

	// if we hear nothing after TTL,
	// despite pings, we declare that
	// parent (child) failed.
	TTL time.Duration

	// call to make or remake a client connection
	CliDial func() (cliConn net.Conn, err error)

	// how often to ping neighbors (parent and child)
	// to see if they are alive. Defaults to TTL/3.
	HeartbeatDur time.Duration

	Verbosity Verbosity
}

type Verbosity int

const (
	QUIET Verbosity = 0 // default
	INFO  Verbosity = 1 // -info
	DEBUG Verbosity = 2 // -debug
)
