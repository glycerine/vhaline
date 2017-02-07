package vhaline

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/glycerine/idem"
	tf "github.com/glycerine/tmframe2"
)

type server struct {
	me      NodeInfo
	rws     []*spair
	halt    idem.Halter
	started bool

	addr2pair      *AtomicAddrToPair
	ArrivingNoteCh chan *Note

	ChildConnectionLost *idem.IdemCloseChan
	lasterr             lasterr

	replica         *Replica
	pairHasShutdown chan *spair
}

func (s *server) GetPair(addr string) *spair {
	return s.addr2pair.Get(addr)
}

type lasterr struct {
	err error
	mux sync.Mutex
}

// set only takes the first err
// it is given, if there is a 2nd one,
// it is silently dropped. So no
// race can clobber the first error.
func (e *lasterr) set(err error) {
	e.mux.Lock()
	if e.err == nil {
		e.err = err
	}
	e.mux.Unlock()
}

func (e *lasterr) get() (err error) {
	e.mux.Lock()
	err = e.err
	e.mux.Unlock()
	return
}

func newServer(me *NodeInfo, r *Replica) *server {
	return &server{
		me:        *me,
		halt:      *idem.NewHalter(),
		addr2pair: NewAtomicAddrToPair(),

		// only one ArrivingNoteChan: how we talk to replica
		ArrivingNoteCh:  make(chan *Note),
		pairHasShutdown: make(chan *spair, 10),

		// if we loose connection to child, we Close() this
		ChildConnectionLost: idem.NewIdemCloseChan(),
		replica:             r,
	}
}

// server pair: a reader and a writer
type spair struct {
	addr           string
	reader         *srvReader
	writer         *srvWriter
	conn           net.Conn
	server         *server
	halt           *idem.Halter
	remoteAddr     string // differentiate so we only handle the correct Outbounds
	ArrivingNoteCh chan *Note
	OutboundNoteCh chan *Note
}

func newSpair(addr string, conn net.Conn, server *server) *spair {

	// halt ties all three structs and two goroutines together.
	halt := idem.NewHalter()
	p := &spair{
		server:         server,
		halt:           halt,
		remoteAddr:     conn.RemoteAddr().String(),
		ArrivingNoteCh: server.ArrivingNoteCh,
		OutboundNoteCh: make(chan *Note),
	}
	p.writer = newSrvWriter(addr, conn, server, halt, p)
	p.reader = newSrvReader(addr, conn, server, halt, p)
	p.reader.writeme = p.writer.writeme

	// lookup addr, get pair to write to.
	server.addr2pair.Set(p.remoteAddr, p)
	return p
}

func (rw *spair) start() error {
	err := rw.reader.start()
	if err != nil {
		return err
	}
	return rw.writer.start()
}

func (rw *spair) stop() {
	rw.halt.ReqStop.Close()
	<-rw.halt.Done.Chan
	select {
	case rw.server.pairHasShutdown <- rw:
	case <-rw.server.halt.ReqStop.Chan:
	}
	if rw != nil && rw.conn != nil {
		rw.conn.Close()
	}
}

func (s *server) stop() {
	for _, sp := range s.rws {
		sp.stop()
	}
}

// remove rw from s.rws
func (s *server) removeRw(rw *spair) {
	n := len(s.rws)
	for i := range s.rws {
		if s.rws[i] == rw {
			if i == n-1 {
				s.rws = s.rws[:i]
			} else {
				s.rws = append(s.rws[:i], s.rws[i+1:]...)
			}
			return
		}
	}
	s.addr2pair.Del(rw.remoteAddr)
}

func (s *server) start() error {
	s.started = true
	var err error
	var lsn net.Listener
	if s.replica.Cfg.Lsn == nil {
		lsn, err = net.Listen("tcp", s.me.Addr)
		if err != nil {
			return err
		}
	} else {
		// client already holding this open for us.
		lsn = s.replica.Cfg.Lsn
	}
	tcpLsn, ok := lsn.(*net.TCPListener)
	if !ok {
		panic("didn't get *TCPListener???")
	}
	s.replica.dlog("server is listening on '%s'", s.me.Addr)
	go func() {
		for {
			tcpLsn.SetDeadline(time.Now().Add(20 * time.Millisecond))
			conn, err := tcpLsn.Accept()

			isTimeout := false
			if err != nil {
				neterr, isNeterr := err.(net.Error)
				if isNeterr {
					if neterr.Timeout() {
						isTimeout = true
					}
				}
			}
			if isTimeout {
				select {
				case <-s.halt.ReqStop.Chan:
					for _, rw := range s.rws {
						rw.stop()
					}
					s.halt.Done.Close()
					return
				case rw := <-s.pairHasShutdown:
					//p("pair %p has shutdown", rw)
					s.removeRw(rw) // allow spair to be garbage collected
				case <-time.After(30 * time.Millisecond):

				}
				continue
			}
			panicOn(err)
			p := newSpair(s.me.Addr, conn, s)
			s.rws = append(s.rws, p)
			err = p.start()
			panicOn(err)

			select {
			case <-time.After(20 * time.Millisecond):
				// service other clients too.

			case <-p.halt.ReqStop.Chan:

			case <-s.halt.ReqStop.Chan:
				for _, rw := range s.rws {
					rw.stop()
				}
				s.halt.Done.Close()
				return
			case rw := <-s.pairHasShutdown:
				//p("pair %p has shutdown", rw)
				s.removeRw(rw) // allow spair to be garbage collected
			}
		}
	}()

	return nil
}

type srvReader struct {
	addr    string
	conn    net.Conn
	halt    *idem.Halter
	writeme chan []byte
	server  *server
	pair    *spair
	readbuf []byte
	frr     *tf.FrameReader
}

func newSrvReader(addr string, c net.Conn, server *server, halt *idem.Halter, pair *spair) *srvReader {
	r := &srvReader{
		addr: addr,
		conn: c,
		halt: halt,
		// writeme is made by the srvWriter
		server:  server,
		readbuf: make([]byte, MaxMsgSize),
		pair:    pair,
	}
	r.frr = tf.NewFrameReaderUsingBuffer(r.conn, r.readbuf)
	return r
}

func (r *srvReader) start() error {
	go func() {
		defer func() {
			// halt our pair as well
			r.halt.ReqStop.Close()
			r.halt.Done.Close()
		}()
		var payload *Note
		var srvCh chan *Note
		var frm *tf.Frame
		var err error
		var nbytes int64
		for {
			frm = nil
			srvCh = nil
			r.conn.SetReadDeadline(time.Now().Add(time.Millisecond * 20))
			// ignore errors, ssh doesn't do deadlines

			frm, nbytes, err, _ = r.frr.NextFrame(nil)
			//n, err := r.conn.Read(r.readbuf)
			isTimeout := false
			if err != nil {
				if neterr, isNetErr := err.(net.Error); isNetErr {
					if neterr.Timeout() {
						isTimeout = true
					}
				}
				srvCh = nil
			}
			if isTimeout {
				srvCh = nil
			} else {
				if err != nil && err.Error() == "EOF" {
					continue
				}
				if err != nil && strings.Contains(err.Error(), "connection reset by peer") {
					r.server.replica.dlog("server reader sees failure or shutdown of client. closing ChildConnectionLost.")
					r.server.lasterr.set(err)
					//p("srvReader is closing ChildConnectionLost")
					r.server.ChildConnectionLost.Close()
					return
				}
				if err != nil && err == tf.FrameTooLargeErr {
					// was getting some mal-formed message: ToParentPing
					// a note not wrapped in a frame.
					r.server.replica.alog("server reader saw unexpected FrameTooLargeErr: request was for %v bytes. Doing buffer reset.", nbytes)
					// try resetting our buffers before doing something drastic
					r.frr = tf.NewFrameReaderUsingBuffer(r.conn, r.readbuf)

					//r.server.ChildConnectionLost.Close()
					continue
				}
				if err != nil {
					r.server.replica.ilog("server reader saw unexpected error '%s'", err)
					r.server.ChildConnectionLost.Close()
					return
				}

				evt := frm.GetEvtnum()
				r.server.replica.dlog("server sees evt %v", evt)
				if evt == tf.EvChecksum || evt == tf.EvChecksumXL2 {
					// skip the checksum, it has
					// a filepos that won't be correct
					// for the new write to disk, and
					// the locator will mess with recovery.
					r.server.replica.dlog("ignoring checksum frame.")

					// TODO: check the checksum against the
					// payload frame we just received. possibly
					// use it to dedup...
					continue
				}

				payload = &Note{}
				payload.From.Addr = r.pair.remoteAddr

				switch NoteEvt(evt) {
				case Checkpoint:
					//if !r.server.replica.isDup(frm) {
					rot := r.server.replica.Rot
					msg, err := rot.WriteCheckpoint(frm)
					if err != nil {
						r.server.replica.dlog("error writing checkpoint to disk: '%s'", msg)
						continue
					}
					r.server.replica.dlog("server reader saved checkpoint to '%s'.",
						rot.CurFile)
					continue

				default:
					_, err := payload.UnmarshalMsg(frm.Data)
					if err != nil {
						payload.From.Addr = r.pair.remoteAddr
						r.server.replica.ilog("server could not unmarshal as Checkpoint or Note, cutting connection to restart the link.")
						r.server.replica.sendToChild(r.server.replica.newNote(RestartLink, &r.server.replica.Me, &payload.From))
						return
					}

					// make sure the remoteAddr is what we think it is.
					// NOTE: this overwrites what the client provided us for From.Addr.
					// But this is important to be consistent, because it
					// is how we key the addr2pair map.
					payload.From.Addr = r.pair.remoteAddr
				}
				r.server.replica.dlog("server received a '%s' message of len %v from '%#v'", payload.Num, len(frm.Data), payload.From)
				srvCh = r.server.ArrivingNoteCh
			}
			if payload == nil {
				select {
				case <-r.halt.ReqStop.Chan:
					// shutdown requested
					return

				default:
					// no shutdown request as of yet.
				}
			} else {
				//r.server.replica.dlog("srvReader sending payload on r.server.ArrivingNoteCh")
				select {
				case srvCh <- payload:
					//r.server.replica.dlog("payload went on srvCh")
					payload = nil
					srvCh = nil

				// case for shutdown:
				case <-r.halt.ReqStop.Chan:
					// shutdown requested
					return
				}
			}
		}

	}()
	return nil
}

type srvWriter struct {
	addr     string
	conn     net.Conn
	halt     *idem.Halter
	writeme  chan []byte
	server   *server
	pair     *spair
	writebuf []byte
}

func newSrvWriter(addr string, conn net.Conn, server *server, halt *idem.Halter, pair *spair) *srvWriter {
	w := &srvWriter{
		addr:     addr,
		conn:     conn,
		halt:     halt,
		writeme:  make(chan []byte),
		server:   server,
		writebuf: make([]byte, MaxMsgSize),
		pair:     pair,
	}
	return w
}

func (w *srvWriter) start() error {
	if w.conn == nil {
		return fmt.Errorf("srvWriter error: w.conn was nil")
	}
	go func() {
		defer func() {
			// halt our pair as well
			w.halt.ReqStop.Close()
			w.halt.Done.Close()
			w.pair.stop() // asymmetric, only here.
		}()
		for {
			select {
			case note := <-w.pair.OutboundNoteCh:

				// we can have two or more children temporarily
				// connected, but sure to route our reply
				// to the correct child (to tell the new one
				// to buzz off!).
				if note.To.Addr != w.pair.remoteAddr {
					panic("wrong server writer!")
					continue
				}

				var bts []byte
				var err error
				if note.Num == Checkpoint {
					bts, err = note.cp.Marshal(w.writebuf[:0])
				} else {
					bts, err = note.MarshalMsg(w.writebuf[:0])
					panicOn(err)

					now := note.SendTm
					bts, err = tf.NewMarshalledFrame(nil, now,
						tf.Evtnum(note.Num), 0, 0, bts, ZebraSchemaId64, 0)
					panicOn(err)
				}

				//w.server.replica.dlog("server sending note '%s' of len %v.",
				//	note.Num, len(bts))

			littleloop:
				for {
					n, err := w.conn.Write(bts)
					if n == len(bts) || err == nil {
						// if we wrote all our bits, then
						// we don't care about various errors, and we'll
						// find out about them on the next attempted write.
						break littleloop
					}
					bts = bts[n:]
					// err could be:
					// write tcp 127.0.0.1:49639->127.0.0.1:49642: write: broken pipe
					serr := err.Error()
					if strings.Contains(serr, "short write") {
						continue littleloop // try again.
					}
					if strings.Contains(serr, "broken pipe") {
						// failure or shutdown of client
						w.server.lasterr.set(err)
						//p("srvWriter is closing ChildConnectionLost")
						w.server.ChildConnectionLost.Close()
						return
					}
					// what else? might be something else we can/should recover from?
					panic(err) // TODO: remove this and just return like broken pipe above.
				}

			// case for shutdown:
			case <-w.halt.ReqStop.Chan:
				// shutdown requested
				return
			}
		}

	}()
	return nil
}
