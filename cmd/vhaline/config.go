package main

import (
	"flag"
	"fmt"

	hc "github.com/glycerine/vhaline/vhaline"
)

// VhalineConfig is the top level, main config
type VhalineConfig struct {
	Upstream string
	upstream hc.AddrHostPort

	Addr string // my address
	addr hc.AddrHostPort

	Nickname string

	Info  bool
	Debug bool

	TtlMillisec  int
	BeatMillisec int

	CpuProfile, MemProfile bool
}

func NewVhalineConfig() *VhalineConfig {
	cfg := &VhalineConfig{}
	return cfg
}

// DefineFlags should be called before myflags.Parse().
func (c *VhalineConfig) DefineFlags(fs *flag.FlagSet) {

	fs.StringVar(&c.Upstream, "parent", "", "host:port address of our parent node, if any.")
	fs.StringVar(&c.Addr, "addr", ":9449", "bind this host:port pair")
	fs.StringVar(&c.Nickname, "name", "", "name for this process")
	fs.BoolVar(&c.Info, "info", false, "be somewhat more verbose about what is happening")
	fs.BoolVar(&c.Debug, "debug", false, "be very verbose, trace all internal ops")

	fs.IntVar(&c.TtlMillisec, "ttl", 0, "milliseconds after which we declare neighbor to have failed. Should be at least 3x the -beat setting, to require 3 failed contact attempts before declaring a node to have failed. Defaults to 3*beat if not set explicitly.")
	fs.IntVar(&c.BeatMillisec, "beat", 1000, "issue a health check every this many milliseconds. Typically should be ttl/3, or ttl/4 or more (so we make 3 or 4 attempts before failing the neighbor). Example: if -ttl 4000, then -beat 1000; this would heartbeat every 1sec, and declare the other node failed after 4sec of no-reply.")

	fs.BoolVar(&c.CpuProfile, "cpu", false, "activate cpu profiling")
	fs.BoolVar(&c.MemProfile, "mem", false, "activate memory profiling")
}

// ValidateConfig should be called after myflags.Parse().
func (c *VhalineConfig) ValidateConfig() error {

	c.upstream.Addr = c.Upstream
	c.upstream.Title = "parent"
	var err error
	err = c.upstream.ParseAddr()
	if err != nil {
		return err
	}

	c.addr.Addr = c.Addr
	c.addr.Title = "addr"
	c.addr.Required = true
	err = c.addr.ParseAddr()
	if err != nil {
		return err
	}

	if c.TtlMillisec == 0 {
		c.TtlMillisec = 3 * c.BeatMillisec
	}
	if c.TtlMillisec < 60 {
		return fmt.Errorf("-ttl must be at least 60 msec.")
	}
	if c.BeatMillisec < 20 {
		return fmt.Errorf("-beat must be at least 20 msec.")
	}
	if c.BeatMillisec*2 >= c.TtlMillisec {
		return fmt.Errorf("-beat must be less than ttl*2;"+
			" we have beat=%v msec, and ttl=%v msec",
			c.BeatMillisec, c.TtlMillisec)
	}

	return nil
}
