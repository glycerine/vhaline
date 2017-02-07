package vhaline

import (
	"bytes"
	crypto "crypto/rand"
	"fmt"
	tf "github.com/glycerine/tmframe2"
	"github.com/glycerine/zebrapack/msgp"
	mr "math/rand"
	"time"
)

//go:generate zebrapack -msgp -no-load

const ZebraSchemaId64 int64 = 0x133a156fb4705 // 338243018770181

type NoteEvt int

const (
	Error NoteEvt = -2

	FromChildConnect  NoteEvt = -3
	ToChildConnectAck NoteEvt = -4

	ToParentPing      NoteEvt = -5
	FromParentPingAck NoteEvt = -6

	ToChildPing      NoteEvt = -7
	FromChildPingAck NoteEvt = -8

	ChainInfo    NoteEvt = -9
	ChainInfoAck NoteEvt = -10

	// do not attach checksum to a checkpoint note/frame.
	Checkpoint    NoteEvt = -11
	CheckpointAck NoteEvt = -12

	// request from parent to child to
	// sever the tcp connection and
	// re-establish it.
	RestartLink      NoteEvt = -13
	AlreadyHaveChild NoteEvt = -14
)

func (e NoteEvt) String() string {
	switch e {
	case Error:
		return "Error"
	case FromChildConnect:
		return "FromChildConnect"
	case ToChildConnectAck:
		return "ToChildConnectAck"
	case ToParentPing:
		return "ToParentPing"
	case FromParentPingAck:
		return "FromParentPingAck"
	case ToChildPing:
		return "ToChildPing"
	case FromChildPingAck:
		return "FromChildPingAck"
	case ChainInfo:
		return "ChainInfo"
	case ChainInfoAck:
		return "ChainInfoAck"
	case Checkpoint:
		return "Checkpoint"
	case CheckpointAck:
		return "CheckpointAck"
	case RestartLink:
		return "RestartLink"
	case AlreadyHaveChild:
		return "AlreadyHaveChild"
	}
	return "*-unknown-NoteEvt-*"
}

// Note has no pointers to avoid data races as it
// gets passed around.
type Note struct {
	Num       NoteEvt
	From      NodeInfo
	To        NodeInfo
	ChainInfo map[string]NodeInfo
	SendTm    time.Time
	Nonce     string

	cp *tf.Frame // for Num == Checkpoint; local only
}

func newNote(evt NoteEvt, from, to *NodeInfo, rsrc *mr.Rand) *Note {
	n := &Note{
		Num:    evt,
		From:   *from,
		To:     *to,
		SendTm: time.Now().UTC(),
		Nonce:  mathRandHexString(40, rsrc),
	}
	return n
}

// NodeInfo models the node graph
// so the Replica can keep track of it.
type NodeInfo struct {
	Id   string // long random string in hex
	Addr string // ssh listening address of this node

	Parent string // nodeid
	Child  string // nodeid

	Role     string // root, middle, tail.
	Nickname string
}

func (n *NodeInfo) ShortId() string {
	if len(n.Id) > 8 {
		return n.Id[:8]
	}
	return ""
}

func (n *NodeInfo) Str() string {
	return fmt.Sprintf("%v/%v", n.ShortId(), n.Addr)
}

func newNodeId() string {
	b := make([]byte, 20)
	_, err := crypto.Read(b)
	panicOn(err)
	return fmt.Sprintf("%x", b)
}

func newNodeInfo() *NodeInfo {
	return &NodeInfo{
		Id: newNodeId(),
	}
}

func newNodeInfoWithoutNodeId() *NodeInfo {
	return &NodeInfo{}
}

func (n *Note) String() string {
	var js bytes.Buffer
	o, err := n.MarshalMsg(nil)
	panicOn(err)
	_, err = msgp.CopyToJSON(&js, bytes.NewBuffer(o))
	panicOn(err)
	pretty := prettyPrintJson(js.Bytes())
	return string(pretty)
}
