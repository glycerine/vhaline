package vhaline

import (
	"fmt"
	"github.com/glycerine/idem"
	tf "github.com/glycerine/tmframe2"
	"log"
	mr "math/rand"
	"os"
	"time"
)

// Replica is the main node
type Replica struct {
	Cfg    Cfg
	Me     NodeInfo
	Parent NodeInfo
	Child  NodeInfo

	parentLiveness liveness
	childLiveness  liveness

	// chain is passed around periodically
	Chain map[string]*NodeInfo
	Mypos int // my position in the chain

	client *client
	server *server

	Halt                         idem.Halter
	ParentFirstContactSuccessful *idem.IdemCloseChan
	ChildFirstContactSuccessful  *idem.IdemCloseChan
	Pid                          int
	Logger                       *log.Logger
	CheckpointArrivedCh          chan *tf.Frame

	// TTL: after no-contact
	// from parent/child of this period,
	// we declare them to have failed.
	TTL time.Duration

	// HeartbeatDur determines how
	// frequently we ping our parent/child.
	// Deafult is TTL/3 if not set explicitly.
	HeartbeatDur time.Duration

	// how many heartbeats since we started?
	// see hcc.beat()

	rsrc *mr.Rand

	hcc        *healthCheckCounter
	lastHealth time.Time
	nextHealth time.Time

	Rot *tf.Rotator

	// we do an async send on this when this node becomes the master
	ParentFailedNotification     chan bool
	ChildFailedNotification      chan bool
	ParentRejectedUsNotification chan bool

	lastChainReport time.Time
}

func newMathRandSource() *mr.Rand {
	return mr.New(mr.NewSource(cryptoRandInt64()))
}

// NewReplica creates a Replica but does not
// start listening on addr. Call Start()
// to do that.
func NewReplica(cfg *Cfg, nickname string) (*Replica, error) {

	if cfg == nil {
		cfg = &Cfg{
			TTL:          time.Second * 3, //default
			HeartbeatDur: 1 * time.Second,
		}
	} else {
		// make our own copy, avoid data races in the tests.
		tmp := *cfg
		cfg = &tmp
	}

	if cfg.Addr == "" {
		lsn, port := getAvailPort()
		cfg.Addr = fmt.Sprintf("127.0.0.1:%v", port)
		cfg.Lsn = lsn
	}

	ahp := &AddrHostPort{Addr: cfg.Addr}
	ahp.Title = "addr"
	err := ahp.ParseAddr()
	if err != nil {
		return nil, err
	}
	me := newNodeInfo()
	me.Addr = cfg.Addr
	me.Nickname = nickname

	r := &Replica{
		Cfg:  *cfg,
		Me:   *me,
		Halt: *idem.NewHalter(),

		ParentFirstContactSuccessful: idem.NewIdemCloseChan(),
		ChildFirstContactSuccessful:  idem.NewIdemCloseChan(),
		Pid: os.Getpid(),

		//Logger: log.New(os.Stderr, "", log.LstdFlags|log.LUTC|log.Llongfile|log.Lmicroseconds),
		Logger: log.New(os.Stderr, "", log.LstdFlags|log.LUTC|log.Lshortfile|log.Lmicroseconds),

		TTL:          cfg.TTL,
		HeartbeatDur: cfg.HeartbeatDur,
		rsrc:         newMathRandSource(),

		parentLiveness: newLiveness(cfg.TTL),
		childLiveness:  newLiveness(cfg.TTL),
		Rot: &tf.Rotator{
			MaxFileSizeBytes: 10 * 1024 * 1024,
			NumFilesToKeep:   3,
		},
		CheckpointArrivedCh: make(chan *tf.Frame),

		ParentFailedNotification:     make(chan bool, 20),
		ChildFailedNotification:      make(chan bool, 20),
		ParentRejectedUsNotification: make(chan bool, 20),
	}
	r.hcc = newHealthCheckCounter(r.Halt.ReqStop.Chan, r.TTL, r)
	r.server = newServer(me, r)
	r.client = newClient(me, r)
	if r.HeartbeatDur == 0 {
		r.HeartbeatDur = r.TTL / 3
	}
	r.dlog("binding '%s'", r.Me.Addr)
	r.dlog("r.Me = '%#v'", r.Me)

	return r, nil
}

func (r *Replica) SetFilePrefix(prefix string) {
	r.Rot.Prefix = prefix
}

func (r *Replica) IsRoot() bool {
	r.dlog("IsRoot() call sees r.Parent: %#v", r.Parent)
	return r.Parent.Addr == ""
}

func (r *Replica) IsTail() bool {
	return r.Child.Addr == "" || r.Child.Id == ""
}

func (r *Replica) IsMiddle() bool {
	return !r.IsRoot() && !r.IsTail()
}

func (r *Replica) AddChild(child *NodeInfo) error {
	if r.Child.Addr != "" {
		return fmt.Errorf("already have child")
	}
	r.Child = *child
	return nil
}

func (r *Replica) FailStop() {
	r.client.stop()
	r.server.stop()
	r.Halt.ReqStop.Close()
	<-r.Halt.Done.Chan
}

func (r *Replica) AddParent(parentAddr string) error {
	if parentAddr == "" {
		return fmt.Errorf("empty parentAddr in call to AddParent.")
	}
	if r.Parent.Addr != "" {
		return fmt.Errorf("already have parent")
	}
	if parentAddr == r.Me.Addr {
		return fmt.Errorf("parent address '%s' conflicts with my own '%s'", parentAddr, r.Me.Addr)
	}
	if r.client.started {
		return fmt.Errorf("Relica.client already started")
	}

	r.Parent = NodeInfo{}
	r.Parent.Addr = parentAddr

	return nil
}

func (m *Replica) Stop() {
	m.server.stop()
	m.Halt.ReqStop.Close()
	<-m.Halt.Done.Chan
}

func (m *Replica) Start() error {
	m.dlog("Start called.")
	if m.Parent.Addr == "" {
		// I am root
		m.Me.Role = "root"
	} else {
		// TODO: change this if parent fails and we take over as root.
		m.Me.Role = "non-root"
	}

	// log rotation
	if m.Rot.Prefix == "" {
		m.Rot.Prefix = "replica." + m.Me.Id + ".log"
	}

	err := m.server.start()
	if err != nil {
		return err
	}

	// Initialize next point in time to do health check.
	// Generally we want to try to get three before declaring
	// the node dead.
	m.nextHealth = time.Now().Add(m.TTL / 3)
	if m.TTL/3 == 0 {
		return fmt.Errorf("m.TTL=%v is too small; m.TTL/3 == 0", m.TTL)
	}

	if m.Parent.Addr != "" {
		m.dlog("Replica.Start() client -> contact parent at '%s'", m.Parent.Addr)
		m.client.me = m.Me
		m.client.parent = m.Parent

		err = m.client.start()
		if err != nil {
			return fmt.Errorf("error contacting parent: %v", err)
		}
		note := newNote(FromChildConnect, &m.Me, &m.Parent, m.rsrc)

		select {
		case m.client.OutboundNoteCh <- note:
			m.dlog("issued FromChildConnect to '%s'", m.Parent.Addr)
		case <-m.Halt.ReqStop.Chan:
			m.Halt.Done.Close()
			return fmt.Errorf("shutdown requested before initial FromChildConnect could be started")
		}

	}

	// typical m usage pattern
	go func() {
		defer func() {
			m.Halt.ReqStop.Close()
			m.Halt.Done.Close()
			m.dlog("Start() is exiting. This replica is shutting down.")
		}()
		loop := -1

		for {
			loop++
			_ = loop
			//m.dlog("top of Start select loop, %v loop.", loop)

			select {
			case cp := <-m.CheckpointArrivedCh:
				m.ilog("a checkpoint frame arrived.")
				childAlive, _, _ := m.childLiveness.isAlive(time.Now())
				if childAlive {
					note := newNote(Checkpoint, &m.Me, &m.Child, m.rsrc)
					note.cp = cp
					err := m.sendToChild(note)
					if err != nil {
						m.dlog("exiting Start on client checkpoint conveyance problem: '%s'", err)
						return
					}
				}

			case <-m.server.ChildConnectionLost.Chan:
				m.ilog("child connection lost")
				//m.childLiveness.reset()
				//m.Child = NodeInfo{} // get rid of address

				// we need to get a new ChildConnectonLost
				// since the old one has been used, and it
				// can only be Closed() once. Since this
				// select case we are in is the only place we receive on
				// ChildConnectionLost.Chan, there should
				// be no race in replacing the channel vie Reinit().
				//
				// WARNING: if other locations start to
				// receive from ChildConnectionLost.Chan
				// they should also have a timeout, or
				// else this Reinit() could cause them
				// to deadlock forever. It will re-make
				// the Chan.
				//
				// This avoids endly busy looping on
				// this case.
				m.server.ChildConnectionLost.Reinit()
				// have to wait for client to try and
				// reconnect; we can't do this ourselves.

			case <-m.client.ParentConnectionLost.Chan:
				// parent failed
				m.ilog("parent connection lost")
				// and avoid endless busy looping here too.
				m.client.ParentConnectionLost.Reinit()
				// don't declare gone right away, wait until
				// liveness checks expire.
				m.client.stop()
				m.client = newClient(&m.Me, m)
				err := m.client.start()
				if err != nil {
					m.ilog("could not restart client: '%s'", err)
				}
				//m.parentLiveness.reset()

			case <-time.After(m.HeartbeatDur):
				// heartbeat
				// fires only if no other select activity
				err := m.doHealthCheck()
				if err != nil {
					m.dlog("exiting on err from health check comm: '%s'", err)
					return
				}

			case note := <-m.client.ArrivingNoteCh:
				//m.dlog("new note from parent: '%s'", note.Num)
				m.ParentFirstContactSuccessful.Close()
				err := m.handleFromParent(note)
				if err != nil {
					m.dlog("exiting on err from parent comm: '%s'", err)
					return
				}

			case note := <-m.server.ArrivingNoteCh:
				//m.dlog("new note from child: '%s'", note.Num)
				m.ChildFirstContactSuccessful.Close()
				err := m.handleFromChild(note)
				if err != nil {
					m.dlog("exiting on err from child comm: '%s'", err)
					return
				}

			case <-m.Halt.ReqStop.Chan:
				m.dlog("shutdown requested.")
				return
			}
		}
	}()
	return nil
}

func nodeInfoFromAddr(addr string) *NodeInfo {
	return &NodeInfo{Addr: addr}
}

// child comm goes via m.server
func (m *Replica) handleFromChild(note *Note) error {
	//m.dlog("handleFromChild called with '%s'", note.Num)

	now := time.Now()

	if m.Child.Id == "" {
		// remember the first one.
		m.Child = note.From
	}

	if note.From.Addr != m.Child.Addr {
		m.ilog("rejecting new child '%s' b/c already have '%s'",
			note.From.Str(),
			m.Child.Str())
		m.sendToChild(newNote(AlreadyHaveChild, &m.Me, &note.From, m.rsrc))
		// give the message a little time to be sent before
		// killing the client connection
		pair := m.server.GetPair(note.From.Addr)
		go func(pair *spair) {
			time.Sleep(2 * time.Second)
			// we shutdown the client connection
			pair.halt.ReqStop.Close()
		}(pair)
		return nil
	}
	// display the lastest process nonce.
	m.Child = note.From

	m.heardFromChild(now)

	switch note.Num {

	case Error:

	case FromChildConnect:
		err := m.recordNewChildConnect(note)
		if err != nil {
			return err
		}
		return m.sendToChild(newNote(ToChildConnectAck, &m.Me, &m.Child, m.rsrc))

	case ToChildConnectAck:
		m.ilog("got FromChildConnect from child %s' at '%s'",
			m.Child.Id, m.Child.Addr)
		panic("should never happen that child sends FromChildConnectAck")

	case ToParentPing:
		m.dlog("sees from child(%s): ToParentPing.", m.Child.Str())
		return m.sendToChild(newNote(FromParentPingAck, &m.Me, &note.From, m.rsrc))

	case FromParentPingAck:
		m.dlog("got FromParentPingAck from child(%s)", m.Child.Str())
		panic("should never happen that child sends FromParentPingAck")

	case ToChildPing:
		m.dlog("got ToChildPing from child(%s)", m.Child.Str())
		panic("should never happen that child sends ToChildPing")
	case FromChildPingAck:
		m.dlog("sees from child: FromChildPingAck.")
		// nothing more

	case ChainInfo:
		m.dlog("sees from child(%s): ChainInfo.", m.Child.Str())
		return m.recvdChainInfo(note)

	case ChainInfoAck:
		m.dlog("sees from child(%s): ChainInfoAck.", m.Child.Str())
		// nothing more

	case Checkpoint:
		m.dlog("got Checkpoint from child(%s)", m.Child.Str())
		panic("should never happen that child sends Checkpoint")

	case CheckpointAck:
		m.dlog("sees from child(%s) CheckpointAck.", m.Child.Str())
		m.recvdAckCheckpoint(note)
		// nothing further

	default:
		panic(fmt.Sprintf("unhandled note.Num '%v' from child comm", note.Num))
	}
	return nil
}

// parent comm goes via m.client
func (m *Replica) handleFromParent(note *Note) error {
	//m.dlog("handleFromParent called with '%s'", note.Num)
	m.heardFromParent(time.Now())

	if m.Parent.Id == "" {
		// remember the first one, for the Id.
		if m.Parent.Addr != note.From.Addr {
			panic("next line needs refinement, because we got from a different address than expected; or else there is something else funky going on. either way, figure it out.")
		}
		m.Parent = note.From
	}

	if m.Parent.Id == "" {
		m.Parent.Id = note.From.Id
	}

	if note.From.Id != m.Parent.Id {
		newid := note.From.Id
		oldid := m.Parent.Id
		if len(newid) > 8 {
			newid = newid[:8]
		}
		if len(oldid) > 8 {
			oldid = oldid[:8]
		}
		m.ilog("parent Id switched; must be new instance. saw '%s' from '%s', but expected '%s' from '%s'",
			newid, note.From.Addr,
			oldid, m.Parent.Addr)
		// prevents the parent from recovering...
		// panic("should never happen that m.Parent.Id changes.")
		m.Parent.Id = note.From.Id
	}

	switch note.Num {

	case AlreadyHaveChild:
		m.ilog("got AlreadyHaveChild from parent(%s). Stopping client.", m.Parent.Str())
		m.client.stop()
		if len(m.ParentRejectedUsNotification) < cap(m.ParentRejectedUsNotification) {
			m.ParentRejectedUsNotification <- true
		}

		return nil

	case RestartLink:
		m.ilog("got RestartLink request from parent(%s). Restarting client.", m.Parent.Str())

		m.client.stop()
		m.client = newClient(&m.Me, m)
		err := m.client.start()
		if err != nil {
			m.alog("could not restart client: '%s'", err)
			m.client.needsRestart = true
			return fmt.Errorf("serious problem: could not restart client: '%s'", err)
		}
		return nil

	case Error:

	case FromChildConnect:
		m.ilog("got FromChildConnect from parent %s' at '%s'", m.Parent.Id, m.Parent.Addr)
		panic("should never happen that parent sends FromChildConnect.")

	case ToChildConnectAck:
		m.ilog("sees from parent: ToChildConnectAck.")
		// nothing further

	case ToParentPing:
		m.dlog("got ToParentPing from parent %s' at '%s'", m.Parent.Id, m.Parent.Addr)
		panic("should never happen that parent sends ToParentPing to its child.")

	case FromParentPingAck:
		m.dlog("sees from parent: FromParentPingAck.")
		// nothing further

	case ToChildPing:
		m.dlog("sees from parent: ToChildPing.")
		return m.sendToParent(newNote(FromChildPingAck, &m.Me, &note.From, m.rsrc))

	case FromChildPingAck:
		m.dlog("got FromChildPingAck from parent %s' at '%s'", m.Parent.Id, m.Parent.Addr)
		panic("should never happen that parent sends FromChildPingAck to its child.")

	case ChainInfo:
		m.ilog("sees from parent: ChainInfo.")
		m.mergeChainInfo(note)
	case ChainInfoAck:
		m.ilog("sees from parent: ChainInfoAck.")
		// nothing further
	case Checkpoint:
		m.ilog("sees from parent: Checkpoint.")
		return m.recvdCheckpoint(note)

	case CheckpointAck:
		m.ilog("got CheckpointAck from parent %s' at '%s'", m.Parent.Id, m.Parent.Addr)
		panic("should never happen that parent sends CheckpointAck to its child.")

	default:
		panic(fmt.Sprintf("unhandled note.Num '%v' from child comm", note.Num))
	}
	return nil
}

// goroutine safe
func (m *Replica) parentAvail() bool {
	now := time.Now().UTC()

	alive, _, _ := m.parentLiveness.isAlive(now)
	//alive, lastContact, ttl := m.parentLiveness.isAlive(now)
	//m.dlog("parentAvail has alive=%v, last-contact: %v, ttl: %v, now=%v", alive, lastContact, ttl, now)

	return alive
}

// goroutine safe
func (m *Replica) childAvail() bool {
	now := time.Now().UTC()

	alive, _, _ := m.childLiveness.isAlive(now)
	//alive, lastContact, ttl := m.childLiveness.isAlive(now)
	//	m.dlog("childAvail has alive=%v, last-contact: %v, ttl: %v, now=%v", alive, lastContact, ttl, now)

	return alive
}

func (m *Replica) doHealthCheck() (err error) {
	beatNum := m.hcc.beat()
	action := false

	defer func() {
		if action {
			m.dlog("doHealthCheck finished, returning err '%v'", err)
		}
	}()
	now := time.Now()

	///	if beatNum%10 == 0 && now.After(m.lastChainReport.Add(time.Minute)) {
	if beatNum%10 == 0 && now.After(m.lastChainReport.Add(time.Second)) {
		m.ilog("line status: parent(%s) -> me(%s) -> child(%s)",
			m.Parent.Str(), m.Me.Str(), m.Child.Str())
		m.lastChainReport = now
	}

	if m.Child.Addr != "" {

		// over TTL?
		calive, clastContact, cttl := m.childLiveness.isAlive(now)
		if !calive {
			id := m.Child.Id
			if len(id) > 8 {
				id = id[:8]
			}
			m.ilog("it's been %s (> ttl == %s) since last child contact,"+
				" declaring child '%s' at '%s' to have failed.",
				now.Sub(clastContact), cttl, id, m.Child.Addr)

			m.childLiveness.reset()
			m.Child = NodeInfo{} // get rid of address
			// async notify our library client, but don't block
			if len(m.ChildFailedNotification) < cap(m.ChildFailedNotification) {
				m.ChildFailedNotification <- true
			}
		} else {

			// done with TTL check

			action = true
			err = m.sendToChild(
				newNote(ToChildPing, &m.Me, &m.Child, m.rsrc))

			if err != nil {
				return err
			}
		}

	}

	if m.Parent.Addr != "" && beatNum > 0 {

		// over TTL?
		palive, plastContact, pttl := m.parentLiveness.isAlive(now)
		if !palive {
			// is this the first time we've tried to contact parent, and
			// we've never heard from them ever?
			if plastContact.IsZero() {

			}

			id := m.Parent.Id
			if len(id) > 8 {
				id = id[:8]
			}
			m.ilog("it's been %s (> ttl == %s; last-contact: '%v') since last parent contact, declaring parent '%s' at '%s' to have failed.", now.Sub(plastContact), pttl, plastContact, id, m.Parent.Addr)

			m.client.stop()
			m.client.needsRestart = true

			m.Parent = NodeInfo{}
			m.parentLiveness.reset()
			// async notify our library client; don't block
			if len(m.ParentFailedNotification) < cap(m.ParentFailedNotification) {
				m.ParentFailedNotification <- true
			}
		} else {
			// parent alive
			// done with TTL check

			// do we need to try and restart client?
			if m.client.needsRestart {
				err := m.client.start()
				if err != nil {
					m.ilog("in healthCheck, could not restart client: "+
						"'%s'. Making newClient", err)

					m.client = newClient(&m.Me, m)
					err := m.client.start()
					if err != nil {
						m.ilog("in healthCheck, after newClient(), still could not restart "+
							"client: '%s'", err)
					}
				}
			} else {

				action = true
				err = m.sendToParent(
					newNote(ToParentPing, &m.Me, &m.Parent, m.rsrc))

				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// sendToParent aka clientReply
func (m *Replica) sendToParent(reply *Note) error {
	m.dlog("sendToParent called with '%v'", reply.Num)

	select {

	case m.client.OutboundNoteCh <- reply:
		return nil

	case <-m.Halt.ReqStop.Chan:
		return fmt.Errorf("shutting down")
	}
}

// sendToChild, aka serverReply
func (m *Replica) sendToChild(reply *Note) error {
	m.dlog("sendToChild called with '%v'", reply.Num)

	pair := m.server.GetPair(reply.To.Addr)
	if pair == nil {
		panic(fmt.Sprintf("bad child address: not found in server.addr2pair : '%s'", reply.To.Addr))
	}

	select {

	case pair.OutboundNoteCh <- reply:
		return nil

	case <-m.Halt.ReqStop.Chan:
		return fmt.Errorf("shutting down")

	case <-time.After(time.Second):
		m.dlog("server may have shutdown. sendToChild could not send '%s' on m.server.OutboundNoteCh after 1 sec", reply.Num)
		return nil // TODO: should this be an error? probably/queue and retry the send later.
	}
}

func (m *Replica) mergeChainInfo(note *Note) error {
	panic("TODO")
	return nil
}

func (m *Replica) recordNewChildConnect(note *Note) error {
	id := note.From.Id
	if len(id) > 8 {
		id = id[:8]
	}
	m.ilog("NewChildConnect from child %s' at '%s'", id, note.From.Addr)

	return nil
}

func (m *Replica) recvdAckCheckpoint(note *Note) error {
	panic("TODO")
	return nil
}

func (m *Replica) recvdCheckpoint(note *Note) error {
	// save this, and rotate
	panic("TODO")
	return nil
}

func (m *Replica) recvdChainInfo(note *Note) error {
	panic("TODO")
	return nil
}

func (m *Replica) heardFromParent(now time.Time) {
	m.parentLiveness.heardFrom(now)
}

func (m *Replica) heardFromChild(now time.Time) {
	m.childLiveness.heardFrom(now)
}

// SaveCheckpoint is the main service API that the
// vhaline library provides to clients. The other
// services are the ParentFailedNotification chan bool
// and 	ChildFailedNotification  chan bool;
// which receives a true event when the parent
// (child) fail. When the parent fails then we
// become the master, and should start writing.
//
func (r *Replica) SaveCheckpoint(frm *tf.Frame) error {
	msg, err := r.Rot.WriteCheckpoint(frm)
	if err != nil {
		r.dlog("error writing checkpoint to disk: '%s'", msg)
		return err
	}
	r.dlog("server reader saved checkpoint to '%s'.", r.Rot.CurFile)

	select {
	case r.CheckpointArrivedCh <- frm:
	case <-r.Halt.ReqStop.Chan:
	}
	return nil
}
