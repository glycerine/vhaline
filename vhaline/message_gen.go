package vhaline

// NOTE: THIS FILE WAS PRODUCED BY THE
// ZEBRAPACK CODE GENERATION TOOL (github.com/glycerine/zebrapack)
// DO NOT EDIT

import (
	"github.com/glycerine/zebrapack/msgp"
)

// DecodeMsg implements msgp.Decodable
// We treat empty fields as if we read a Nil from the wire.
func (z *NodeInfo) DecodeMsg(dc *msgp.Reader) (err error) {
	var sawTopNil bool
	if dc.IsNil() {
		sawTopNil = true
		err = dc.ReadNil()
		if err != nil {
			return
		}
		dc.PushAlwaysNil()
	}

	var field []byte
	_ = field
	const maxFields0zuiw = 6

	// -- templateDecodeMsg starts here--
	var totalEncodedFields0zuiw uint32
	totalEncodedFields0zuiw, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	encodedFieldsLeft0zuiw := totalEncodedFields0zuiw
	missingFieldsLeft0zuiw := maxFields0zuiw - totalEncodedFields0zuiw

	var nextMiss0zuiw int32 = -1
	var found0zuiw [maxFields0zuiw]bool
	var curField0zuiw string

doneWithStruct0zuiw:
	// First fill all the encoded fields, then
	// treat the remaining, missing fields, as Nil.
	for encodedFieldsLeft0zuiw > 0 || missingFieldsLeft0zuiw > 0 {
		//fmt.Printf("encodedFieldsLeft: %v, missingFieldsLeft: %v, found: '%v', fields: '%#v'\n", encodedFieldsLeft0zuiw, missingFieldsLeft0zuiw, msgp.ShowFound(found0zuiw[:]), decodeMsgFieldOrder0zuiw)
		if encodedFieldsLeft0zuiw > 0 {
			encodedFieldsLeft0zuiw--
			field, err = dc.ReadMapKeyPtr()
			if err != nil {
				return
			}
			curField0zuiw = msgp.UnsafeString(field)
		} else {
			//missing fields need handling
			if nextMiss0zuiw < 0 {
				// tell the reader to only give us Nils
				// until further notice.
				dc.PushAlwaysNil()
				nextMiss0zuiw = 0
			}
			for nextMiss0zuiw < maxFields0zuiw && (found0zuiw[nextMiss0zuiw] || decodeMsgFieldSkip0zuiw[nextMiss0zuiw]) {
				nextMiss0zuiw++
			}
			if nextMiss0zuiw == maxFields0zuiw {
				// filled all the empty fields!
				break doneWithStruct0zuiw
			}
			missingFieldsLeft0zuiw--
			curField0zuiw = decodeMsgFieldOrder0zuiw[nextMiss0zuiw]
		}
		//fmt.Printf("switching on curField: '%v'\n", curField0zuiw)
		switch curField0zuiw {
		// -- templateDecodeMsg ends here --

		case "Id":
			found0zuiw[0] = true
			z.Id, err = dc.ReadString()
			if err != nil {
				return
			}
		case "Addr":
			found0zuiw[1] = true
			z.Addr, err = dc.ReadString()
			if err != nil {
				return
			}
		case "Parent":
			found0zuiw[2] = true
			z.Parent, err = dc.ReadString()
			if err != nil {
				return
			}
		case "Child":
			found0zuiw[3] = true
			z.Child, err = dc.ReadString()
			if err != nil {
				return
			}
		case "Role":
			found0zuiw[4] = true
			z.Role, err = dc.ReadString()
			if err != nil {
				return
			}
		case "Nickname":
			found0zuiw[5] = true
			z.Nickname, err = dc.ReadString()
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	if nextMiss0zuiw != -1 {
		dc.PopAlwaysNil()
	}

	if sawTopNil {
		dc.PopAlwaysNil()
	}

	if p, ok := interface{}(z).(msgp.PostLoad); ok {
		p.PostLoadHook()
	}

	return
}

// fields of NodeInfo
var decodeMsgFieldOrder0zuiw = []string{"Id", "Addr", "Parent", "Child", "Role", "Nickname"}

var decodeMsgFieldSkip0zuiw = []bool{false, false, false, false, false, false}

// fieldsNotEmpty supports omitempty tags
func (z *NodeInfo) fieldsNotEmpty(isempty []bool) uint32 {
	return 6
}

// EncodeMsg implements msgp.Encodable
func (z *NodeInfo) EncodeMsg(en *msgp.Writer) (err error) {
	if p, ok := interface{}(z).(msgp.PreSave); ok {
		p.PreSaveHook()
	}

	// map header, size 6
	// write "Id"
	err = en.Append(0x86, 0xa2, 0x49, 0x64)
	if err != nil {
		return err
	}
	err = en.WriteString(z.Id)
	if err != nil {
		return
	}
	// write "Addr"
	err = en.Append(0xa4, 0x41, 0x64, 0x64, 0x72)
	if err != nil {
		return err
	}
	err = en.WriteString(z.Addr)
	if err != nil {
		return
	}
	// write "Parent"
	err = en.Append(0xa6, 0x50, 0x61, 0x72, 0x65, 0x6e, 0x74)
	if err != nil {
		return err
	}
	err = en.WriteString(z.Parent)
	if err != nil {
		return
	}
	// write "Child"
	err = en.Append(0xa5, 0x43, 0x68, 0x69, 0x6c, 0x64)
	if err != nil {
		return err
	}
	err = en.WriteString(z.Child)
	if err != nil {
		return
	}
	// write "Role"
	err = en.Append(0xa4, 0x52, 0x6f, 0x6c, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteString(z.Role)
	if err != nil {
		return
	}
	// write "Nickname"
	err = en.Append(0xa8, 0x4e, 0x69, 0x63, 0x6b, 0x6e, 0x61, 0x6d, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteString(z.Nickname)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *NodeInfo) MarshalMsg(b []byte) (o []byte, err error) {
	if p, ok := interface{}(z).(msgp.PreSave); ok {
		p.PreSaveHook()
	}

	o = msgp.Require(b, z.Msgsize())
	// map header, size 6
	// string "Id"
	o = append(o, 0x86, 0xa2, 0x49, 0x64)
	o = msgp.AppendString(o, z.Id)
	// string "Addr"
	o = append(o, 0xa4, 0x41, 0x64, 0x64, 0x72)
	o = msgp.AppendString(o, z.Addr)
	// string "Parent"
	o = append(o, 0xa6, 0x50, 0x61, 0x72, 0x65, 0x6e, 0x74)
	o = msgp.AppendString(o, z.Parent)
	// string "Child"
	o = append(o, 0xa5, 0x43, 0x68, 0x69, 0x6c, 0x64)
	o = msgp.AppendString(o, z.Child)
	// string "Role"
	o = append(o, 0xa4, 0x52, 0x6f, 0x6c, 0x65)
	o = msgp.AppendString(o, z.Role)
	// string "Nickname"
	o = append(o, 0xa8, 0x4e, 0x69, 0x63, 0x6b, 0x6e, 0x61, 0x6d, 0x65)
	o = msgp.AppendString(o, z.Nickname)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *NodeInfo) UnmarshalMsg(bts []byte) (o []byte, err error) {
	return z.UnmarshalMsgWithCfg(bts, nil)
}
func (z *NodeInfo) UnmarshalMsgWithCfg(bts []byte, cfg *msgp.RuntimeConfig) (o []byte, err error) {
	var nbs msgp.NilBitsStack
	nbs.Init(cfg)
	var sawTopNil bool
	if msgp.IsNil(bts) {
		sawTopNil = true
		bts = nbs.PushAlwaysNil(bts[1:])
	}

	var field []byte
	_ = field
	const maxFields1zmlx = 6

	// -- templateUnmarshalMsg starts here--
	var totalEncodedFields1zmlx uint32
	if !nbs.AlwaysNil {
		totalEncodedFields1zmlx, bts, err = nbs.ReadMapHeaderBytes(bts)
		if err != nil {
			return
		}
	}
	encodedFieldsLeft1zmlx := totalEncodedFields1zmlx
	missingFieldsLeft1zmlx := maxFields1zmlx - totalEncodedFields1zmlx

	var nextMiss1zmlx int32 = -1
	var found1zmlx [maxFields1zmlx]bool
	var curField1zmlx string

doneWithStruct1zmlx:
	// First fill all the encoded fields, then
	// treat the remaining, missing fields, as Nil.
	for encodedFieldsLeft1zmlx > 0 || missingFieldsLeft1zmlx > 0 {
		//fmt.Printf("encodedFieldsLeft: %v, missingFieldsLeft: %v, found: '%v', fields: '%#v'\n", encodedFieldsLeft1zmlx, missingFieldsLeft1zmlx, msgp.ShowFound(found1zmlx[:]), unmarshalMsgFieldOrder1zmlx)
		if encodedFieldsLeft1zmlx > 0 {
			encodedFieldsLeft1zmlx--
			field, bts, err = nbs.ReadMapKeyZC(bts)
			if err != nil {
				return
			}
			curField1zmlx = msgp.UnsafeString(field)
		} else {
			//missing fields need handling
			if nextMiss1zmlx < 0 {
				// set bts to contain just mnil (0xc0)
				bts = nbs.PushAlwaysNil(bts)
				nextMiss1zmlx = 0
			}
			for nextMiss1zmlx < maxFields1zmlx && (found1zmlx[nextMiss1zmlx] || unmarshalMsgFieldSkip1zmlx[nextMiss1zmlx]) {
				nextMiss1zmlx++
			}
			if nextMiss1zmlx == maxFields1zmlx {
				// filled all the empty fields!
				break doneWithStruct1zmlx
			}
			missingFieldsLeft1zmlx--
			curField1zmlx = unmarshalMsgFieldOrder1zmlx[nextMiss1zmlx]
		}
		//fmt.Printf("switching on curField: '%v'\n", curField1zmlx)
		switch curField1zmlx {
		// -- templateUnmarshalMsg ends here --

		case "Id":
			found1zmlx[0] = true
			z.Id, bts, err = nbs.ReadStringBytes(bts)

			if err != nil {
				return
			}
		case "Addr":
			found1zmlx[1] = true
			z.Addr, bts, err = nbs.ReadStringBytes(bts)

			if err != nil {
				return
			}
		case "Parent":
			found1zmlx[2] = true
			z.Parent, bts, err = nbs.ReadStringBytes(bts)

			if err != nil {
				return
			}
		case "Child":
			found1zmlx[3] = true
			z.Child, bts, err = nbs.ReadStringBytes(bts)

			if err != nil {
				return
			}
		case "Role":
			found1zmlx[4] = true
			z.Role, bts, err = nbs.ReadStringBytes(bts)

			if err != nil {
				return
			}
		case "Nickname":
			found1zmlx[5] = true
			z.Nickname, bts, err = nbs.ReadStringBytes(bts)

			if err != nil {
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	if nextMiss1zmlx != -1 {
		bts = nbs.PopAlwaysNil()
	}

	if sawTopNil {
		bts = nbs.PopAlwaysNil()
	}
	o = bts
	if p, ok := interface{}(z).(msgp.PostLoad); ok {
		p.PostLoadHook()
	}

	return
}

// fields of NodeInfo
var unmarshalMsgFieldOrder1zmlx = []string{"Id", "Addr", "Parent", "Child", "Role", "Nickname"}

var unmarshalMsgFieldSkip1zmlx = []bool{false, false, false, false, false, false}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *NodeInfo) Msgsize() (s int) {
	s = 1 + 3 + msgp.StringPrefixSize + len(z.Id) + 5 + msgp.StringPrefixSize + len(z.Addr) + 7 + msgp.StringPrefixSize + len(z.Parent) + 6 + msgp.StringPrefixSize + len(z.Child) + 5 + msgp.StringPrefixSize + len(z.Role) + 9 + msgp.StringPrefixSize + len(z.Nickname)
	return
}

// DecodeMsg implements msgp.Decodable
// We treat empty fields as if we read a Nil from the wire.
func (z *Note) DecodeMsg(dc *msgp.Reader) (err error) {
	var sawTopNil bool
	if dc.IsNil() {
		sawTopNil = true
		err = dc.ReadNil()
		if err != nil {
			return
		}
		dc.PushAlwaysNil()
	}

	var field []byte
	_ = field
	const maxFields2zihb = 6

	// -- templateDecodeMsg starts here--
	var totalEncodedFields2zihb uint32
	totalEncodedFields2zihb, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	encodedFieldsLeft2zihb := totalEncodedFields2zihb
	missingFieldsLeft2zihb := maxFields2zihb - totalEncodedFields2zihb

	var nextMiss2zihb int32 = -1
	var found2zihb [maxFields2zihb]bool
	var curField2zihb string

doneWithStruct2zihb:
	// First fill all the encoded fields, then
	// treat the remaining, missing fields, as Nil.
	for encodedFieldsLeft2zihb > 0 || missingFieldsLeft2zihb > 0 {
		//fmt.Printf("encodedFieldsLeft: %v, missingFieldsLeft: %v, found: '%v', fields: '%#v'\n", encodedFieldsLeft2zihb, missingFieldsLeft2zihb, msgp.ShowFound(found2zihb[:]), decodeMsgFieldOrder2zihb)
		if encodedFieldsLeft2zihb > 0 {
			encodedFieldsLeft2zihb--
			field, err = dc.ReadMapKeyPtr()
			if err != nil {
				return
			}
			curField2zihb = msgp.UnsafeString(field)
		} else {
			//missing fields need handling
			if nextMiss2zihb < 0 {
				// tell the reader to only give us Nils
				// until further notice.
				dc.PushAlwaysNil()
				nextMiss2zihb = 0
			}
			for nextMiss2zihb < maxFields2zihb && (found2zihb[nextMiss2zihb] || decodeMsgFieldSkip2zihb[nextMiss2zihb]) {
				nextMiss2zihb++
			}
			if nextMiss2zihb == maxFields2zihb {
				// filled all the empty fields!
				break doneWithStruct2zihb
			}
			missingFieldsLeft2zihb--
			curField2zihb = decodeMsgFieldOrder2zihb[nextMiss2zihb]
		}
		//fmt.Printf("switching on curField: '%v'\n", curField2zihb)
		switch curField2zihb {
		// -- templateDecodeMsg ends here --

		case "Num":
			found2zihb[0] = true
			{
				var zkoe int
				zkoe, err = dc.ReadInt()
				z.Num = NoteEvt(zkoe)
			}
			if err != nil {
				return
			}
		case "From":
			found2zihb[1] = true
			err = z.From.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "To":
			found2zihb[2] = true
			err = z.To.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "ChainInfo":
			found2zihb[3] = true
			var zwfl uint32
			zwfl, err = dc.ReadMapHeader()
			if err != nil {
				return
			}
			if z.ChainInfo == nil && zwfl > 0 {
				z.ChainInfo = make(map[string]NodeInfo, zwfl)
			} else if len(z.ChainInfo) > 0 {
				for key, _ := range z.ChainInfo {
					delete(z.ChainInfo, key)
				}
			}
			for zwfl > 0 {
				zwfl--
				var zsgv string
				var zswi NodeInfo
				zsgv, err = dc.ReadString()
				if err != nil {
					return
				}
				err = zswi.DecodeMsg(dc)
				if err != nil {
					return
				}
				z.ChainInfo[zsgv] = zswi
			}
		case "SendTm":
			found2zihb[4] = true
			z.SendTm, err = dc.ReadTime()
			if err != nil {
				return
			}
		case "Nonce":
			found2zihb[5] = true
			z.Nonce, err = dc.ReadString()
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	if nextMiss2zihb != -1 {
		dc.PopAlwaysNil()
	}

	if sawTopNil {
		dc.PopAlwaysNil()
	}

	if p, ok := interface{}(z).(msgp.PostLoad); ok {
		p.PostLoadHook()
	}

	return
}

// fields of Note
var decodeMsgFieldOrder2zihb = []string{"Num", "From", "To", "ChainInfo", "SendTm", "Nonce"}

var decodeMsgFieldSkip2zihb = []bool{false, false, false, false, false, false}

// fieldsNotEmpty supports omitempty tags
func (z *Note) fieldsNotEmpty(isempty []bool) uint32 {
	return 6
}

// EncodeMsg implements msgp.Encodable
func (z *Note) EncodeMsg(en *msgp.Writer) (err error) {
	if p, ok := interface{}(z).(msgp.PreSave); ok {
		p.PreSaveHook()
	}

	// map header, size 6
	// write "Num"
	err = en.Append(0x86, 0xa3, 0x4e, 0x75, 0x6d)
	if err != nil {
		return err
	}
	err = en.WriteInt(int(z.Num))
	if err != nil {
		return
	}
	// write "From"
	err = en.Append(0xa4, 0x46, 0x72, 0x6f, 0x6d)
	if err != nil {
		return err
	}
	err = z.From.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "To"
	err = en.Append(0xa2, 0x54, 0x6f)
	if err != nil {
		return err
	}
	err = z.To.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "ChainInfo"
	err = en.Append(0xa9, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x49, 0x6e, 0x66, 0x6f)
	if err != nil {
		return err
	}
	err = en.WriteMapHeader(uint32(len(z.ChainInfo)))
	if err != nil {
		return
	}
	for zsgv, zswi := range z.ChainInfo {
		err = en.WriteString(zsgv)
		if err != nil {
			return
		}
		err = zswi.EncodeMsg(en)
		if err != nil {
			return
		}
	}
	// write "SendTm"
	err = en.Append(0xa6, 0x53, 0x65, 0x6e, 0x64, 0x54, 0x6d)
	if err != nil {
		return err
	}
	err = en.WriteTime(z.SendTm)
	if err != nil {
		return
	}
	// write "Nonce"
	err = en.Append(0xa5, 0x4e, 0x6f, 0x6e, 0x63, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteString(z.Nonce)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Note) MarshalMsg(b []byte) (o []byte, err error) {
	if p, ok := interface{}(z).(msgp.PreSave); ok {
		p.PreSaveHook()
	}

	o = msgp.Require(b, z.Msgsize())
	// map header, size 6
	// string "Num"
	o = append(o, 0x86, 0xa3, 0x4e, 0x75, 0x6d)
	o = msgp.AppendInt(o, int(z.Num))
	// string "From"
	o = append(o, 0xa4, 0x46, 0x72, 0x6f, 0x6d)
	o, err = z.From.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "To"
	o = append(o, 0xa2, 0x54, 0x6f)
	o, err = z.To.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "ChainInfo"
	o = append(o, 0xa9, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x49, 0x6e, 0x66, 0x6f)
	o = msgp.AppendMapHeader(o, uint32(len(z.ChainInfo)))
	for zsgv, zswi := range z.ChainInfo {
		o = msgp.AppendString(o, zsgv)
		o, err = zswi.MarshalMsg(o)
		if err != nil {
			return
		}
	}
	// string "SendTm"
	o = append(o, 0xa6, 0x53, 0x65, 0x6e, 0x64, 0x54, 0x6d)
	o = msgp.AppendTime(o, z.SendTm)
	// string "Nonce"
	o = append(o, 0xa5, 0x4e, 0x6f, 0x6e, 0x63, 0x65)
	o = msgp.AppendString(o, z.Nonce)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Note) UnmarshalMsg(bts []byte) (o []byte, err error) {
	return z.UnmarshalMsgWithCfg(bts, nil)
}
func (z *Note) UnmarshalMsgWithCfg(bts []byte, cfg *msgp.RuntimeConfig) (o []byte, err error) {
	var nbs msgp.NilBitsStack
	nbs.Init(cfg)
	var sawTopNil bool
	if msgp.IsNil(bts) {
		sawTopNil = true
		bts = nbs.PushAlwaysNil(bts[1:])
	}

	var field []byte
	_ = field
	const maxFields3zmry = 6

	// -- templateUnmarshalMsg starts here--
	var totalEncodedFields3zmry uint32
	if !nbs.AlwaysNil {
		totalEncodedFields3zmry, bts, err = nbs.ReadMapHeaderBytes(bts)
		if err != nil {
			return
		}
	}
	encodedFieldsLeft3zmry := totalEncodedFields3zmry
	missingFieldsLeft3zmry := maxFields3zmry - totalEncodedFields3zmry

	var nextMiss3zmry int32 = -1
	var found3zmry [maxFields3zmry]bool
	var curField3zmry string

doneWithStruct3zmry:
	// First fill all the encoded fields, then
	// treat the remaining, missing fields, as Nil.
	for encodedFieldsLeft3zmry > 0 || missingFieldsLeft3zmry > 0 {
		//fmt.Printf("encodedFieldsLeft: %v, missingFieldsLeft: %v, found: '%v', fields: '%#v'\n", encodedFieldsLeft3zmry, missingFieldsLeft3zmry, msgp.ShowFound(found3zmry[:]), unmarshalMsgFieldOrder3zmry)
		if encodedFieldsLeft3zmry > 0 {
			encodedFieldsLeft3zmry--
			field, bts, err = nbs.ReadMapKeyZC(bts)
			if err != nil {
				return
			}
			curField3zmry = msgp.UnsafeString(field)
		} else {
			//missing fields need handling
			if nextMiss3zmry < 0 {
				// set bts to contain just mnil (0xc0)
				bts = nbs.PushAlwaysNil(bts)
				nextMiss3zmry = 0
			}
			for nextMiss3zmry < maxFields3zmry && (found3zmry[nextMiss3zmry] || unmarshalMsgFieldSkip3zmry[nextMiss3zmry]) {
				nextMiss3zmry++
			}
			if nextMiss3zmry == maxFields3zmry {
				// filled all the empty fields!
				break doneWithStruct3zmry
			}
			missingFieldsLeft3zmry--
			curField3zmry = unmarshalMsgFieldOrder3zmry[nextMiss3zmry]
		}
		//fmt.Printf("switching on curField: '%v'\n", curField3zmry)
		switch curField3zmry {
		// -- templateUnmarshalMsg ends here --

		case "Num":
			found3zmry[0] = true
			{
				var zfbe int
				zfbe, bts, err = nbs.ReadIntBytes(bts)

				if err != nil {
					return
				}
				z.Num = NoteEvt(zfbe)
			}
		case "From":
			found3zmry[1] = true
			bts, err = z.From.UnmarshalMsg(bts)
			if err != nil {
				return
			}
			if err != nil {
				return
			}
		case "To":
			found3zmry[2] = true
			bts, err = z.To.UnmarshalMsg(bts)
			if err != nil {
				return
			}
			if err != nil {
				return
			}
		case "ChainInfo":
			found3zmry[3] = true
			if nbs.AlwaysNil {
				if len(z.ChainInfo) > 0 {
					for key, _ := range z.ChainInfo {
						delete(z.ChainInfo, key)
					}
				}

			} else {

				var ztwj uint32
				ztwj, bts, err = nbs.ReadMapHeaderBytes(bts)
				if err != nil {
					return
				}
				if z.ChainInfo == nil && ztwj > 0 {
					z.ChainInfo = make(map[string]NodeInfo, ztwj)
				} else if len(z.ChainInfo) > 0 {
					for key, _ := range z.ChainInfo {
						delete(z.ChainInfo, key)
					}
				}
				for ztwj > 0 {
					var zsgv string
					var zswi NodeInfo
					ztwj--
					zsgv, bts, err = nbs.ReadStringBytes(bts)
					if err != nil {
						return
					}
					bts, err = zswi.UnmarshalMsg(bts)
					if err != nil {
						return
					}
					if err != nil {
						return
					}
					z.ChainInfo[zsgv] = zswi
				}
			}
		case "SendTm":
			found3zmry[4] = true
			z.SendTm, bts, err = nbs.ReadTimeBytes(bts)

			if err != nil {
				return
			}
		case "Nonce":
			found3zmry[5] = true
			z.Nonce, bts, err = nbs.ReadStringBytes(bts)

			if err != nil {
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	if nextMiss3zmry != -1 {
		bts = nbs.PopAlwaysNil()
	}

	if sawTopNil {
		bts = nbs.PopAlwaysNil()
	}
	o = bts
	if p, ok := interface{}(z).(msgp.PostLoad); ok {
		p.PostLoadHook()
	}

	return
}

// fields of Note
var unmarshalMsgFieldOrder3zmry = []string{"Num", "From", "To", "ChainInfo", "SendTm", "Nonce"}

var unmarshalMsgFieldSkip3zmry = []bool{false, false, false, false, false, false}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *Note) Msgsize() (s int) {
	s = 1 + 4 + msgp.IntSize + 5 + z.From.Msgsize() + 3 + z.To.Msgsize() + 10 + msgp.MapHeaderSize
	if z.ChainInfo != nil {
		for zsgv, zswi := range z.ChainInfo {
			_ = zswi
			_ = zsgv
			s += msgp.StringPrefixSize + len(zsgv) + zswi.Msgsize()
		}
	}
	s += 7 + msgp.TimeSize + 6 + msgp.StringPrefixSize + len(z.Nonce)
	return
}

// DecodeMsg implements msgp.Decodable
// We treat empty fields as if we read a Nil from the wire.
func (z *NoteEvt) DecodeMsg(dc *msgp.Reader) (err error) {
	var sawTopNil bool
	if dc.IsNil() {
		sawTopNil = true
		err = dc.ReadNil()
		if err != nil {
			return
		}
		dc.PushAlwaysNil()
	}

	{
		var zlbr int
		zlbr, err = dc.ReadInt()
		(*z) = NoteEvt(zlbr)
	}
	if err != nil {
		return
	}
	if sawTopNil {
		dc.PopAlwaysNil()
	}

	if p, ok := interface{}(z).(msgp.PostLoad); ok {
		p.PostLoadHook()
	}

	return
}

// EncodeMsg implements msgp.Encodable
func (z NoteEvt) EncodeMsg(en *msgp.Writer) (err error) {
	if p, ok := interface{}(z).(msgp.PreSave); ok {
		p.PreSaveHook()
	}

	err = en.WriteInt(int(z))
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z NoteEvt) MarshalMsg(b []byte) (o []byte, err error) {
	if p, ok := interface{}(z).(msgp.PreSave); ok {
		p.PreSaveHook()
	}

	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendInt(o, int(z))
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *NoteEvt) UnmarshalMsg(bts []byte) (o []byte, err error) {
	return z.UnmarshalMsgWithCfg(bts, nil)
}
func (z *NoteEvt) UnmarshalMsgWithCfg(bts []byte, cfg *msgp.RuntimeConfig) (o []byte, err error) {
	var nbs msgp.NilBitsStack
	nbs.Init(cfg)
	var sawTopNil bool
	if msgp.IsNil(bts) {
		sawTopNil = true
		bts = nbs.PushAlwaysNil(bts[1:])
	}

	{
		var ztzt int
		ztzt, bts, err = nbs.ReadIntBytes(bts)

		if err != nil {
			return
		}
		(*z) = NoteEvt(ztzt)
	}
	if sawTopNil {
		bts = nbs.PopAlwaysNil()
	}
	o = bts
	if p, ok := interface{}(z).(msgp.PostLoad); ok {
		p.PostLoadHook()
	}

	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z NoteEvt) Msgsize() (s int) {
	s = msgp.IntSize
	return
}
