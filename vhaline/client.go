package vhaline

import (
	"fmt"
	"net"
	"time"

	"github.com/glycerine/idem"
	tf "github.com/glycerine/tmframe2"
)

type client struct {
	rw             *cpair
	me             NodeInfo
	parent         NodeInfo
	ArrivingNoteCh chan *Note
	OutboundNoteCh chan *Note
	started        bool

	ParentConnectionLost *idem.IdemCloseChan
	lasterr              lasterr

	replica *Replica

	// optional, maybe supplied by replica to
	// Dial through ssh tunnel to parent.
	// might be unset/nil.
	cliDial func() (cliConn net.Conn, err error)

	needsRestart bool
}

func newClient(me *NodeInfo, r *Replica) *client {
	return &client{
		me:                   *me,
		parent:               r.Parent,
		ArrivingNoteCh:       make(chan *Note),
		OutboundNoteCh:       make(chan *Note),
		ParentConnectionLost: idem.NewIdemCloseChan(),
		replica:              r,
		cliDial:              r.Cfg.CliDial,
		needsRestart:         true, // by default, cleared upon successful start().
	}
}

// client pair: a reader and a writer
type cpair struct {
	addr   string
	reader *cliReader
	writer *cliWriter
	conn   net.Conn
	halt   idem.Halter
	client *client
}

func newCliRwpair(addr string, conn net.Conn, client *client) *cpair {
	p := &cpair{
		reader: newCliReader(addr, conn, client),
		writer: newCliWriter(addr, conn, client),
		client: client,
	}
	p.reader.writeme = p.writer.writeme
	return p
}

func (rw *cpair) start() error {
	err := rw.reader.start()
	if err != nil {
		return err
	}
	return rw.writer.start()
}

func (rw *cpair) stop() {
	rw.reader.halt.ReqStop.Close()
	<-rw.reader.halt.Done.Chan
	rw.writer.halt.ReqStop.Close()
	<-rw.writer.halt.Done.Chan
}

func (c *client) start() error {
	c.started = true
	var conn net.Conn
	var err error
	if c.cliDial == nil {
		c.replica.ilog("client dialing parent at '%s'", c.parent.Addr)
		conn, err = net.Dial("tcp", c.parent.Addr)
		if err != nil {
			return err
		}
	} else {
		conn, err = c.cliDial()
		if err != nil {
			return err
		}
	}

	c.rw = newCliRwpair(c.parent.Addr, conn, c)
	err = c.rw.start()
	if err != nil {
		return fmt.Errorf("client.start() error from "+
			"newCliRwpair to parent '%s': '%s'", c.parent.Addr, err)
	}

	c.needsRestart = false
	return nil
}

func (c *client) stop() {
	if c.rw != nil {
		c.rw.stop()
	}
}

// client writer

type cliWriter struct {
	addr     string
	conn     net.Conn
	halt     idem.Halter
	writeme  chan []byte
	client   *client
	writebuf []byte
}

func newCliWriter(addr string, conn net.Conn, client *client) *cliWriter {
	w := &cliWriter{
		addr:     addr,
		conn:     conn,
		halt:     *idem.NewHalter(),
		writeme:  make(chan []byte),
		client:   client,
		writebuf: make([]byte, MaxMsgSize),
	}
	return w
}

func (w *cliWriter) start() error {
	go func() {
		for {
			select {
			case note := <-w.client.OutboundNoteCh:

				var bts []byte
				var err error
				if note.Num == Checkpoint {
					bts, err = note.cp.Marshal(w.writebuf[:0])
					panicOn(err)
				} else {
					w.client.replica.dlog("client (to parent) writer has '%s'", note.Num)
					bts, err = note.MarshalMsg(w.writebuf[:0])
					panicOn(err)

					now := note.SendTm
					bts, err = tf.NewMarshalledFrame(nil, now,
						tf.Evtnum(note.Num), 0, 0, bts, ZebraSchemaId64, 0)
					panicOn(err)
					// debug:
					/*
					if note.Num == ToParentPing {
						w.client.replica.alog("client (to parent) writer has '%s' marshalled as: bts(len=%v)='%#v'", note.Num, len(bts), bts)
					}*/
				}
				n, err := w.conn.Write(bts)
				panicOn(err)
				if n != len(bts) {
					panic("short write")
				}
				// important piece of logging. do not delete:
				w.client.replica.dlog("client wrote %s to %s/%s",
					note.Num, note.To.Nickname, note.To.Addr)
				if note.Num == ToParentPing {
					// debug, why no frame wrapper?
					var frm tf.Frame
					frm.Unmarshal(bts, true)
					w.client.replica.dlog("ToParentPing contents were: %s",
						frm)
				}
				
			case <-w.halt.ReqStop.Chan:
				w.client.replica.dlog("ReqStop received, shutting down")
				// shutdown requested
				w.halt.Done.Close()
				return
			}
		}

	}()
	return nil
}

// client reader

type cliReader struct {
	addr    string
	conn    net.Conn
	halt    idem.Halter
	writeme chan []byte
	client  *client
	frr     *tf.FrameReader

	readbuf []byte
}

const MaxMsgSize = 16 * 1024 * 1026 // 16MB

func newCliReader(addr string, c net.Conn, client *client) *cliReader {
	r := &cliReader{
		addr:    addr,
		conn:    c,
		halt:    *idem.NewHalter(),
		client:  client,
		readbuf: make([]byte, MaxMsgSize),
	}
	r.frr = tf.NewFrameReaderUsingBuffer(r.conn, r.readbuf)
	return r
}

func (r *cliReader) start() error {

	go func() {
		defer func() {
			r.halt.Done.Close()
		}()
		var payload *Note
		var cliCh chan *Note
		var frm *tf.Frame
		var err error
		for {
			frm = nil
			cliCh = nil
			r.conn.SetReadDeadline(time.Now().Add(time.Millisecond * 20))
			// ssh doesn't support read deadlines, ignore this error.

			frm, _, err, _ = r.frr.NextFrame(nil)

			//			n, err := r.conn.Read(r.readbuf)
			isTimeout := false
			if err != nil {
				if neterr, isNetErr := err.(net.Error); isNetErr {
					if neterr.Timeout() {
						isTimeout = true
					}
				}
				cliCh = nil
			}
			if isTimeout {
				cliCh = nil
			} else {
				if err != nil {
					// "EOF"
					// parent process probably died.
					r.client.replica.dlog("client reader sees error '%s': concluding parent connection lost.", err)
					r.client.ParentConnectionLost.Close()
					return
					//panicOn(err) // EOF on shutdown.
				}

				evt := frm.GetEvtnum()
				if evt == tf.EvChecksum || evt == tf.EvChecksumXL2 {
					// skip the checksum, it has
					// a filepos that won't be correct
					// for the new write to disk, and
					// the locator will mess with recovery.
					r.client.replica.dlog("ignoring checksum frame.")

					// TODO: check the checksum against the
					// payload frame we just received.
					continue
				}

				payload = &Note{}
				switch NoteEvt(evt) {
				case Checkpoint:
					rot := r.client.replica.Rot
					msg, err := rot.WriteCheckpoint(frm)
					if err != nil {
						r.client.replica.dlog("error writing checkpoint to disk: '%s'", msg)
						continue
					}
					r.client.replica.dlog("client reader saved checkpoint to '%s'.",
						rot.CurFile)
					select {
					case r.client.replica.CheckpointArrivedCh <- frm:
					case <-r.halt.ReqStop.Chan:
						// shutdown requested
						r.halt.Done.Close()
						return
					}
					continue
				default:
					_, err := payload.UnmarshalMsg(frm.Data)
					panicOn(err)
				}

				r.client.replica.dlog("client received '%s'", payload.Num)
				cliCh = r.client.ArrivingNoteCh
			}
			if payload == nil {
				select {
				case <-r.halt.ReqStop.Chan:
					// shutdown requested
					r.halt.Done.Close()
					return
				default:
					// no shutdown request as of yet.
				}
			} else {
				select {
				case cliCh <- payload:
					payload = nil
					cliCh = nil

				// case for shutdown:
				case <-r.halt.ReqStop.Chan:
					// shutdown requested
					r.halt.Done.Close()
					return
				}
			}
		}

	}()
	return nil
}
