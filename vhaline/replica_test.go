// package vhaline provides AP high availability
// from a chain of processes.
package vhaline

import (
	//"bytes"
	"fmt"
	"net"
	"testing"
	"time"

	//"github.com/glycerine/cryrand"
	cv "github.com/glycerine/goconvey/convey"
	"github.com/glycerine/sshego"
	tf "github.com/glycerine/tmframe2"
)

func Test001FailureChecks(t *testing.T) {

	cv.Convey("We regularly check for parent (and child) failure.", t, func() {

		ttl := 600 * time.Millisecond

		// reserve all ports up front, so we don't
		// get test bleed one into the other across
		// lingering network connections.
		n := 1
		cfgs := make([]Cfg, n*3)
		for i := 0; i < n*3; i++ {
			lsn, port := getAvailPort()
			cfgs[i].TTL = ttl
			cfgs[i].Lsn = lsn
			cfgs[i].Addr = fmt.Sprintf("127.0.0.1:%v", port)
		}

		for k := 0; k < n; k++ {
			fmt.Printf("on run k= %v of %v.\n", k, n)
			root, err := NewReplica(&cfgs[k], "alister")
			panicOn(err)
			err = root.Start()
			panicOn(err)

			mid, err := NewReplica(&cfgs[n+k], "bella")
			panicOn(err)
			panicOn(mid.AddParent(root.Me.Addr))
			err = mid.Start()
			panicOn(err)

			tail, err := NewReplica(&cfgs[n+n+k], "chaplin")
			panicOn(err)
			panicOn(tail.AddParent(mid.Me.Addr))
			err = tail.Start()
			panicOn(err)

			<-mid.ParentFirstContactSuccessful.Chan
			<-tail.ParentFirstContactSuccessful.Chan

			root.ParentMustHaveFailed()
			root.ChildMustNotHaveFailed() // panic here, on run 7, 12

			mid.ParentMustNotHaveFailed()
			mid.ChildMustNotHaveFailed()

			tail.ParentMustNotHaveFailed()
			tail.ChildMustHaveFailed()

			root.FailStop()

			last := mid.hcc.lastNum()
			//p("7777 last = %v", last) // typically last is -1 here.

			// each health check is 1/3 of TTL, so
			// wait at least 1 full TTL to be sure we
			// should have detected; i.e. 7 heatbeats.
			// wait for 2 health checks to have been done.
			//t0 := time.Now()
			err = mid.hcc.waitUntil(last + 8)
			panicOn(err)
			//p("22222 after waituntil elap=%v, cur beat = %v", time.Since(t0), mid.hcc.lastNum())

			mid.ParentMustHaveFailed()
			mid.ChildMustNotHaveFailed()

			tail.FailStop()

			last = mid.hcc.lastNum()
			//p("27777 last = %v", last)
			// wait for 2 health checks to have been done.
			err = mid.hcc.waitUntil(last + 8)
			panicOn(err)

			mid.ParentMustHaveFailed()
			mid.ChildMustHaveFailed() // sometimes child will still be up??? before run 17

			// finish shutdown/cleanup
			mid.FailStop()
		}
		// should get here without panicing on any
		// the Must calls.
		cv.So(true, cv.ShouldBeTrue)
	})
}

func Test002MiddleRole(t *testing.T) {
	cv.Convey("if upstream is not nil, then we are a middle of last node. We listen for checkpoints from upstream, persistent them, rotate them, and copy them to our child (if we have one).", t, func() {

		origdir, tmpdir := MakeAndMoveToTempDir()
		_, _ = origdir, tmpdir
		defer TempDirCleanup(origdir, tmpdir)

		a, b, c := threeNodeTestSetup()

		a.Cfg.SetVerbosity(DEBUG)
		b.Cfg.SetVerbosity(DEBUG)
		c.Cfg.SetVerbosity(DEBUG)

		now := time.Now()

		data := []byte("hello-world")
		frm, err := tf.NewFrame(now, tf.Evtnum(Checkpoint), 0, 0, data, ZebraSchemaId64, 0)
		panicOn(err)

		apre := "alister-stuff"
		a.SetFilePrefix(apre)

		bpre := "bella-stuff"
		b.SetFilePrefix(bpre)

		cpre := "chaplin-stuff"
		c.SetFilePrefix(cpre)

		err = a.SaveCheckpoint(frm)
		panicOn(err)

		time.Sleep(100 * time.Millisecond)

		// check on alister's write:
		rot := tf.Rotator{
			Prefix: apre,
		}
		afrm2, err := rot.InitialRestoreState(time.Now())
		panicOn(err)
		p("cool. recovered frame frm2: '%s'", afrm2)
		cv.So(afrm2.Data, cv.ShouldResemble, frm.Data)

		// bella receives a frame from alister and
		// a) writes it to disk
		// c) passes it to chaplin
		// d) chaplin writes to disk too.

		// check on bella's write
		brot := tf.Rotator{
			Prefix: bpre,
		}
		bfrm2, err := brot.InitialRestoreState(time.Now())
		panicOn(err)
		cv.So(bfrm2, cv.ShouldNotBeNil)
		p("cool. recovered frame bfrm2: '%s'; =%#v", bfrm2, bfrm2)
		cv.So(bfrm2.Data,
			cv.ShouldResemble,
			frm.Data)

		time.Sleep(100 * time.Millisecond)

		// check on chaplin's write
		crot := tf.Rotator{
			Prefix: cpre,
		}
		cfrm2, err := crot.InitialRestoreState(time.Now())
		panicOn(err)
		cv.So(cfrm2, cv.ShouldNotBeNil)
		p("cool. recovered frame cfrm2: '%s'", cfrm2)
		cv.So(cfrm2.Data,
			cv.ShouldResemble,
			frm.Data)

	})
}

func Test003ReconfigContactAttempted(t *testing.T) {
	cv.Convey("b1) if we know about the chain and we can't contact our parent, then we try to contact in turn each grand-parent, then great-grand-parent, then great-great-grand-parent, etc. each time we try every 500msec for up to 10 seconds.", t, func() {
	})
}

func Test007ChildFailAllowsReplacement(t *testing.T) {
	cv.Convey("* IIa) child failures should be detected. Once detected and cleared, we should allow another, different node to subscribe as our new child.", t, func() {

		// already implemented, but needs a test
	})
}

func Test009Dedup(t *testing.T) {
	cv.Convey("We should have a means of dedup-ing the checkpoints so we can recognize that we've already gotten a checkpoint and we don't propagate it downstream. If we only propagate things that are new to us, that is much more efficient/saves on bandwidth. Since the checkpoints come in a pair of payload frame + checkpoint frame, so a) we can simply drop any messages that are too old. and b) we can just use the checksum frame's signature hash to do deduplication.", t, func() {

		// can skip this for now, not critical
	})
}

func Test010SshSecuresConnections(t *testing.T) {

	if !testing.Short() {
		t.Skip("skipping test 010 that takes a while. use -short (ironic) to run me.")
	}

	cv.Convey("We connect as ssh client to our parent, who acts as sshd", t, func() {

		ttl := 3 * time.Second
		heartbeat := 1 * time.Second

		cfgs := make([]Cfg, 2)
		for i := 0; i < 2; i++ {
			lsn, port := getAvailPort()
			cfgs[i].TTL = ttl
			cfgs[i].Lsn = lsn
			cfgs[i].Addr = fmt.Sprintf("127.0.0.1:%v", port)
			cfgs[i].HeartbeatDur = heartbeat
		}

		alisterAddr := cfgs[0].Addr

		alister, err := NewReplica(&cfgs[0], "alister")
		panicOn(err)
		err = alister.Start()
		panicOn(err)
		defer func() {
			alister.Stop()
		}()

		srvCfg := sshego.NewSshegoConfig()
		srvCfg.SkipTOTP = true
		srvCfg.SkipPassphrase = true

		// take care of things that would be configured
		// upon installation.
		SetupSshdTestConfig(srvCfg)

		// these ports will be set by config/cmd line
		// options; for testing just grab some unused ports.
		sshdLsn, sshdLsnPort := GetAvailPort() // sshd local listen
		xportLsn, xport := GetAvailPort()      // xport
		sshdLsn.Close()
		xportLsn.Close()

		srvCfg.SshegoSystemMutexPort = xport
		srvCfg.EmbeddedSSHd.Title = "esshd"
		srvCfg.EmbeddedSSHd.Addr = fmt.Sprintf("127.0.0.1:%v", sshdLsnPort)
		srvCfg.EmbeddedSSHd.ParseAddr()
		srvCfg.EmbeddedSSHdHostDbPath = srvCfg.Tempdir + "/server_hostdb"

		srvCfg.NewEsshd()
		srvCfg.Esshd.Start()
		defer func() {
			srvCfg.Esshd.Stop()
			<-srvCfg.Esshd.Halt.Done.Chan
			TempDirCleanup(srvCfg.Origdir, srvCfg.Tempdir)
		}()

		// create a new acct
		mylogin, _, rsaPath, _, err := TestCreateNewAccount(srvCfg)
		panicOn(err)
		_, _ = mylogin, rsaPath

		cliCfg := sshego.NewSshegoConfig()
		SetupSshdTestConfig(cliCfg)

		// allow server to be discovered
		cliCfg.AddIfNotKnown = true
		cliCfg.TestAllowOneshotConnect = true

		// tell the client not to run an esshd
		cliCfg.EmbeddedSSHd.Addr = ""
		cliCfg.RemoteToLocal.Listen.Addr = ""

		dc := sshego.DialConfig{
			ClientKnownHostsPath: cliCfg.ClientKnownHostsPath,
			Mylogin:              mylogin,
			RsaPath:              rsaPath,
			Sshdhost:             srvCfg.EmbeddedSSHd.Host,
			Sshdport:             srvCfg.EmbeddedSSHd.Port,
			DownstreamHostPort:   alisterAddr,
			TofuAddIfNotKnown:    true,
		}

		// first time we add the server key
		_, _, err = dc.Dial()
		cv.So(err.Error(), cv.ShouldContainSubstring, "Re-run without -new")

		// second time we connect based on that server key
		dc.TofuAddIfNotKnown = false

		// tell client how to dial alister/root.
		cfgs[1].CliDial = func() (cliConn net.Conn, err error) {
			c, _, err := dc.Dial()
			return c, err
		}

		bella, err := NewReplica(&cfgs[1], "bella")
		panicOn(err)
		panicOn(bella.AddParent(alister.Me.Addr))
		err = bella.Start()
		panicOn(err)
		defer func() {
			bella.Stop()
		}()

		time.Sleep(10 * time.Second)

		cv.So(true, cv.ShouldEqual, true) // we should get here.
	})
}

func startBackgroundTestTcpServer(serverDone chan bool, payloadByteCount int, confirmationPayload string, confirmationReply string, tcpSrvLsn net.Listener) {
	go func() {
		p("startBackgroundTestTcpServer() about to call Accept().")
		tcpServerConn, err := tcpSrvLsn.Accept()
		panicOn(err)
		p("startBackgroundTestTcpServer() progress: got Accept() back: %v",
			tcpServerConn)

		b := make([]byte, payloadByteCount)
		n, err := tcpServerConn.Read(b)
		panicOn(err)
		if n != payloadByteCount {
			panic(fmt.Errorf("read too short! got %v but expected %v", n, payloadByteCount))
		}
		saw := string(b)

		if saw != confirmationPayload {
			panic(fmt.Errorf("expected '%s', but saw '%s'", confirmationPayload, saw))
		}

		p("success! server got expected confirmation payload of '%s'", saw)

		// reply back
		n, err = tcpServerConn.Write([]byte(confirmationReply))
		panicOn(err)
		if n != payloadByteCount {
			panic(fmt.Errorf("write too short! got %v but expected %v", n, payloadByteCount))
		}
		//tcpServerConn.Close()
		close(serverDone)
	}()
}

func Test012ErrorIfCannotContactParent(t *testing.T) {
	cv.Convey("if upstream is not nil, then on startup, we should error out and halt if cannot contact upstream", t, func() {

		// don't start root, just get
		// an unused port to simulate it not being available.
		lsn, rootport := getAvailPort()
		lsn.Close()
		rootAddr := fmt.Sprintf("127.0.0.1:%v", rootport)

		ttl := 300 * time.Millisecond
		cfg := &Cfg{TTL: ttl}

		mid, err := NewReplica(cfg, "mid")
		panicOn(err)
		err = mid.AddParent(rootAddr)
		panicOn(err)
		err = mid.Start()
		cv.So(err, cv.ShouldNotBeNil)
		p("err = %v", err)
	})
}

func Test013ParentReportsGoodsConnection(t *testing.T) {
	cv.Convey("if upstream is not nil, then on startup, the parent should report when child is connected", t, func() {

		ttl := 300 * time.Millisecond
		cfg := &Cfg{TTL: ttl}

		root, err := NewReplica(cfg, "root")
		panicOn(err)
		err = root.Start()
		panicOn(err)

		mid, err := NewReplica(cfg, "mid")
		panicOn(err)
		err = mid.AddParent(root.Me.Addr)
		cv.So(err, cv.ShouldBeNil)
		err = mid.Start()
		panicOn(err)

		select {
		case <-mid.ParentFirstContactSuccessful.Chan:
		case <-time.After(time.Second):
			panic("mid did not contact parent within a second")
		}
		select {
		case <-root.ChildFirstContactSuccessful.Chan:
		case <-time.After(time.Second):
			panic("root did not register child contact within a second")
		}

	})
}

func Test014RecoveryBeforeTTLAllowed(t *testing.T) {

	cv.Convey("For a TTL (ttl; time-to-live) of 3 seconds and a heartbeat every 1 second, if the parent fails the first two heartbeats but recovers before the 3rd heartbeat, then the connection between the parent and the child issuing the pings should recover and remain intact. This allows some temporary network fluxuation or packet loss before the child failsover to be a write-master. Hence the parent/child, at each hearbeat prior to TTL expiry, must attempt to re-establish their connections on each heartbeat if that connection has been broken before the to ttl expiry.", t, func() {

		a, b := twoNodeTestSetup()
		_, _ = a, b
		// okay, initial setup is done, now to
		// have the network go away for
		// two heartbeats...

		// should get here without panicing on any
		// the Must calls.
		cv.So(true, cv.ShouldBeTrue)
	})
}

func oneNodeTestSetup(ttl, heartbeat time.Duration) (a *Replica) {

	cfgs := make([]Cfg, 1)
	for i := 0; i < 1; i++ {
		lsn, port := getAvailPort()
		cfgs[i].TTL = ttl
		cfgs[i].Lsn = lsn
		cfgs[i].Addr = fmt.Sprintf("127.0.0.1:%v", port)
		cfgs[i].HeartbeatDur = heartbeat
	}

	root, err := NewReplica(&cfgs[0], "alister")
	panicOn(err)
	err = root.Start()
	panicOn(err)

	root.ParentMustHaveFailed()
	root.ChildMustHaveFailed()

	return root
}

func twoNodeTestSetup() (a, b *Replica) {

	ttl := 3 * time.Second
	heartbeat := 1 * time.Second

	cfgs := make([]Cfg, 2)
	for i := 0; i < 2; i++ {
		lsn, port := getAvailPort()
		cfgs[i].TTL = ttl
		cfgs[i].Lsn = lsn
		cfgs[i].Addr = fmt.Sprintf("127.0.0.1:%v", port)
		cfgs[i].HeartbeatDur = heartbeat
	}

	root, err := NewReplica(&cfgs[0], "alister")
	panicOn(err)
	err = root.Start()
	panicOn(err)

	mid, err := NewReplica(&cfgs[1], "bella")
	panicOn(err)
	panicOn(mid.AddParent(root.Me.Addr))
	err = mid.Start()
	panicOn(err)

	<-root.ChildFirstContactSuccessful.Chan
	<-mid.ParentFirstContactSuccessful.Chan

	root.ParentMustHaveFailed()
	root.ChildMustNotHaveFailed()

	mid.ParentMustNotHaveFailed()
	mid.ChildMustHaveFailed()

	return root, mid
}

func threeNodeTestSetup() (a, b, c *Replica) {
	ttl := 3 * time.Second
	heartbeat := 1 * time.Second

	n := 1 // 10
	cfgs := make([]Cfg, n*3)
	for i := 0; i < n*3; i++ {
		lsn, port := getAvailPort()
		cfgs[i].TTL = ttl
		cfgs[i].HeartbeatDur = heartbeat
		cfgs[i].Lsn = lsn
		cfgs[i].Addr = fmt.Sprintf("127.0.0.1:%v", port)
	}

	root, err := NewReplica(&cfgs[0], "alister")
	panicOn(err)
	err = root.Start()
	panicOn(err)

	mid, err := NewReplica(&cfgs[1], "bella")
	panicOn(err)
	panicOn(mid.AddParent(root.Me.Addr))
	err = mid.Start()
	panicOn(err)

	tail, err := NewReplica(&cfgs[2], "chaplin")
	panicOn(err)
	panicOn(tail.AddParent(mid.Me.Addr))
	err = tail.Start()
	panicOn(err)

	<-root.ChildFirstContactSuccessful.Chan
	<-mid.ParentFirstContactSuccessful.Chan
	<-mid.ChildFirstContactSuccessful.Chan
	<-tail.ParentFirstContactSuccessful.Chan

	root.ParentMustHaveFailed()
	root.ChildMustNotHaveFailed() // panic here, on run 7, 12

	mid.ParentMustNotHaveFailed()
	mid.ChildMustNotHaveFailed()

	tail.ParentMustNotHaveFailed()
	tail.ChildMustHaveFailed()

	return root, mid, tail
}
