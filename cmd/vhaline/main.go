package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	// profiling:
	_ "net/http/pprof"
	"runtime/pprof"

	"github.com/glycerine/vhaline/vhaline"
)

const ProgramName = "vhaline"

func main() {

	myflags := flag.NewFlagSet("myflags", flag.ExitOnError)
	cfg := &VhalineConfig{}
	cfg.DefineFlags(myflags)

	err := myflags.Parse(os.Args[1:])
	err = cfg.ValidateConfig()
	if err != nil {
		log.Fatalf("%s command line flag error: '%s'", ProgramName, err)
	}

	c := &vhaline.Cfg{
		Addr:         cfg.Addr,
		TTL:          time.Millisecond * time.Duration(cfg.TtlMillisec),
		HeartbeatDur: time.Millisecond * time.Duration(cfg.BeatMillisec),
	}
	if cfg.Info {
		c.Verbosity = vhaline.INFO
	}
	if cfg.Debug {
		c.Verbosity = vhaline.DEBUG
	}
	me, err := vhaline.NewReplica(c, cfg.Nickname)
	panicOn(err)
	if cfg.Upstream != "" {
		err = me.AddParent(cfg.Upstream)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		p("added upstream parent '%s'", cfg.Upstream)
	}
	comment := ""
	root := ""
	if me.IsRoot() {
		root = "root "
	} else {
		root = "mid "
		comment = fmt.Sprintf("my parent is '%s'", me.Parent.Addr)
	}
	if c.Verbosity >= vhaline.INFO {
		log.Printf("%snode '%v' started on %v. %v", root, me.Me.Id[:8], me.Me.Addr, comment)
	}
	log.Printf("%snode '%s' using ttl=%v and beat=%v.", root, me.Me.Id[:8], c.TTL, c.HeartbeatDur)
	err = me.Start()
	panicOn(err)

	if cfg.CpuProfile {
		f, err := os.Create(fmt.Sprintf("./cpu-profile.%v", os.Getpid()))
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if cfg.MemProfile {
		f, err := os.Create(fmt.Sprintf("./memory-profile.%v", os.Getpid()))
		if err != nil {
			log.Fatal(err)
		}
		pprof.WriteHeapProfile(f)
		f.Close()
	}

	if cfg.CpuProfile || cfg.MemProfile {
		go func() {
			log.Println(http.ListenAndServe("localhost:6060", nil))
		}()
	}

	// demonstrate how to listen for events:

	par0 := me.ParentFirstContactSuccessful.Chan
	chd0 := me.ChildFirstContactSuccessful.Chan
	pfail := me.ParentFailedNotification
	cfail := me.ChildFailedNotification
	for {
		select {
		case <-par0:
			log.Printf("contacted parent event.")
			par0 = nil
		case <-chd0:
			log.Printf("contacted child event.")
			chd0 = nil
		case <-pfail:
			log.Printf("parent failed event.")
		case <-cfail:
			log.Printf("child failed event.")
		case <-me.ParentRejectedUsNotification:
			log.Printf("serious problem: parent rejected us for another.")
			os.Exit(1)
		}
	}
}
