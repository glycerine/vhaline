package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	// profiling:
	_ "net/http/pprof"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"syscall"

	"github.com/glycerine/vhaline/vhaline"
)

const ProgramName = "vhaline"

func main() {

	// handle SIGQUIT without stopping, to
	// get a stacktrace on the fly.
	sigChan := make(chan os.Signal)
	go func() {
		stacktrace := make([]byte, 8192)
		for _ = range sigChan {
			length := runtime.Stack(stacktrace, true)
			fmt.Println(string(stacktrace[:length]))
		}
	}()
	signal.Notify(sigChan, syscall.SIGQUIT)

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
		//p("added upstream parent '%s'", cfg.Upstream)
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
	log.Printf("%s", vhaline.GoVersion())
	err = me.Start()
	panicOn(err)

	if cfg.CpuProfile {
		fn := fmt.Sprintf("./cpu-profile.%v", os.Getpid())
		f, err := os.Create(fn)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
		go func() {
			for {
				select {
				case <-time.After(time.Minute):
					f.Sync()
				}
			}
		}()
	}

	if cfg.MemProfile {
		f, err := os.Create(fmt.Sprintf("./memory-profile.%v", os.Getpid()))
		if err != nil {
			log.Fatal(err)
		}
		pprof.WriteHeapProfile(f)
		f.Close()
	}

	if cfg.CpuProfile || cfg.MemProfile || cfg.WebProfile {
		// find an unused port, startin at 6060
		port := 6060
		ppaddr := ""
		i := 0
		for i = 0; i < 100; i++ {
			ppaddr = fmt.Sprintf("localhost:%v", port+i)
			if !addrAlreadyBound(ppaddr) {
				break
			}
		}
		if i == 100 {
			panic("could not find port for pprof web server in range 6060-6159")
		}
		log.Printf("pprof profiler providing web server on '%s'", ppaddr)
		go func() {
			log.Println(http.ListenAndServe(ppaddr, nil))
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
			log.Printf("cmd/vhaline/main.go: contacted parent event.")
			par0 = nil // prevent endless loop on the closed channel.

		case <-chd0:
			log.Printf("cmd/vhaline/main.go: contacted child event.")
			chd0 = nil

		case <-pfail:
			log.Printf("cmd/vhaline/main.go: parent failed event.")

		case <-cfail:
			log.Printf("cmd/vhaline/main.go:child failed event.")

		case <-me.ParentRejectedUsNotification:
			log.Printf("cmd/vhaline/main.go: serious problem: parent rejected us for another.")
			os.Exit(1)
		}
	}
}

func addrAlreadyBound(addr string) bool {

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return true
	}
	ln.Close()
	return false
}
