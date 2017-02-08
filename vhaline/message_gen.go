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
	const maxFields0zrez = 7

	// -- templateDecodeMsg starts here--
	var totalEncodedFields0zrez uint32
	totalEncodedFields0zrez, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	encodedFieldsLeft0zrez := totalEncodedFields0zrez
	missingFieldsLeft0zrez := maxFields0zrez - totalEncodedFields0zrez

	var nextMiss0zrez int32 = -1
	var found0zrez [maxFields0zrez]bool
	var curField0zrez string

doneWithStruct0zrez:
	// First fill all the encoded fields, then
	// treat the remaining, missing fields, as Nil.
	for encodedFieldsLeft0zrez > 0 || missingFieldsLeft0zrez > 0 {
		//fmt.Printf("encodedFieldsLeft: %v, missingFieldsLeft: %v, found: '%v', fields: '%#v'\n", encodedFieldsLeft0zrez, missingFieldsLeft0zrez, msgp.ShowFound(found0zrez[:]), decodeMsgFieldOrder0zrez)
		if encodedFieldsLeft0zrez > 0 {
			encodedFieldsLeft0zrez--
			field, err = dc.ReadMapKeyPtr()
			if err != nil {
				return
			}
			curField0zrez = msgp.UnsafeString(field)
		} else {
			//missing fields need handling
			if nextMiss0zrez < 0 {
				// tell the reader to only give us Nils
				// until further notice.
				dc.PushAlwaysNil()
				nextMiss0zrez = 0
			}
			for nextMiss0zrez < maxFields0zrez && (found0zrez[nextMiss0zrez] || decodeMsgFieldSkip0zrez[nextMiss0zrez]) {
				nextMiss0zrez++
			}
			if nextMiss0zrez == maxFields0zrez {
				// filled all the empty fields!
				break doneWithStruct0zrez
			}
			missingFieldsLeft0zrez--
			curField0zrez = decodeMsgFieldOrder0zrez[nextMiss0zrez]
		}
		//fmt.Printf("switching on curField: '%v'\n", curField0zrez)
		switch curField0zrez {
		// -- templateDecodeMsg ends here --

		case "Id":
			found0zrez[0] = true
			z.Id, err = dc.ReadString()
			if err != nil {
				return
			}
		case "Addr":
			found0zrez[1] = true
			z.Addr, err = dc.ReadString()
			if err != nil {
				return
			}
		case "Host":
			found0zrez[2] = true
			z.Host, err = dc.ReadString()
			if err != nil {
				return
			}
		case "Port":
			found0zrez[3] = true
			z.Port, err = dc.ReadInt()
			if err != nil {
				return
			}
		case "Parent":
			found0zrez[4] = true
			z.Parent, err = dc.ReadString()
			if err != nil {
				return
			}
		case "Child":
			found0zrez[5] = true
			z.Child, err = dc.ReadString()
			if err != nil {
				return
			}
		case "Nickname":
			found0zrez[6] = true
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
	if nextMiss0zrez != -1 {
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
var decodeMsgFieldOrder0zrez = []string{"Id", "Addr", "Host", "Port", "Parent", "Child", "Nickname"}

var decodeMsgFieldSkip0zrez = []bool{false, false, false, false, false, false, false}

// fieldsNotEmpty supports omitempty tags
func (z *NodeInfo) fieldsNotEmpty(isempty []bool) uint32 {
	return 7
}

// EncodeMsg implements msgp.Encodable
func (z *NodeInfo) EncodeMsg(en *msgp.Writer) (err error) {
	if p, ok := interface{}(z).(msgp.PreSave); ok {
		p.PreSaveHook()
	}

	// map header, size 7
	// write "Id"
	err = en.Append(0x87, 0xa2, 0x49, 0x64)
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
	// write "Host"
	err = en.Append(0xa4, 0x48, 0x6f, 0x73, 0x74)
	if err != nil {
		return err
	}
	err = en.WriteString(z.Host)
	if err != nil {
		return
	}
	// write "Port"
	err = en.Append(0xa4, 0x50, 0x6f, 0x72, 0x74)
	if err != nil {
		return err
	}
	err = en.WriteInt(z.Port)
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
	// map header, size 7
	// string "Id"
	o = append(o, 0x87, 0xa2, 0x49, 0x64)
	o = msgp.AppendString(o, z.Id)
	// string "Addr"
	o = append(o, 0xa4, 0x41, 0x64, 0x64, 0x72)
	o = msgp.AppendString(o, z.Addr)
	// string "Host"
	o = append(o, 0xa4, 0x48, 0x6f, 0x73, 0x74)
	o = msgp.AppendString(o, z.Host)
	// string "Port"
	o = append(o, 0xa4, 0x50, 0x6f, 0x72, 0x74)
	o = msgp.AppendInt(o, z.Port)
	// string "Parent"
	o = append(o, 0xa6, 0x50, 0x61, 0x72, 0x65, 0x6e, 0x74)
	o = msgp.AppendString(o, z.Parent)
	// string "Child"
	o = append(o, 0xa5, 0x43, 0x68, 0x69, 0x6c, 0x64)
	o = msgp.AppendString(o, z.Child)
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
	const maxFields1zitd = 7

	// -- templateUnmarshalMsg starts here--
	var totalEncodedFields1zitd uint32
	if !nbs.AlwaysNil {
		totalEncodedFields1zitd, bts, err = nbs.ReadMapHeaderBytes(bts)
		if err != nil {
			return
		}
	}
	encodedFieldsLeft1zitd := totalEncodedFields1zitd
	missingFieldsLeft1zitd := maxFields1zitd - totalEncodedFields1zitd

	var nextMiss1zitd int32 = -1
	var found1zitd [maxFields1zitd]bool
	var curField1zitd string

doneWithStruct1zitd:
	// First fill all the encoded fields, then
	// treat the remaining, missing fields, as Nil.
	for encodedFieldsLeft1zitd > 0 || missingFieldsLeft1zitd > 0 {
		//fmt.Printf("encodedFieldsLeft: %v, missingFieldsLeft: %v, found: '%v', fields: '%#v'\n", encodedFieldsLeft1zitd, missingFieldsLeft1zitd, msgp.ShowFound(found1zitd[:]), unmarshalMsgFieldOrder1zitd)
		if encodedFieldsLeft1zitd > 0 {
			encodedFieldsLeft1zitd--
			field, bts, err = nbs.ReadMapKeyZC(bts)
			if err != nil {
				return
			}
			curField1zitd = msgp.UnsafeString(field)
		} else {
			//missing fields need handling
			if nextMiss1zitd < 0 {
				// set bts to contain just mnil (0xc0)
				bts = nbs.PushAlwaysNil(bts)
				nextMiss1zitd = 0
			}
			for nextMiss1zitd < maxFields1zitd && (found1zitd[nextMiss1zitd] || unmarshalMsgFieldSkip1zitd[nextMiss1zitd]) {
				nextMiss1zitd++
			}
			if nextMiss1zitd == maxFields1zitd {
				// filled all the empty fields!
				break doneWithStruct1zitd
			}
			missingFieldsLeft1zitd--
			curField1zitd = unmarshalMsgFieldOrder1zitd[nextMiss1zitd]
		}
		//fmt.Printf("switching on curField: '%v'\n", curField1zitd)
		switch curField1zitd {
		// -- templateUnmarshalMsg ends here --

		case "Id":
			found1zitd[0] = true
			z.Id, bts, err = nbs.ReadStringBytes(bts)

			if err != nil {
				return
			}
		case "Addr":
			found1zitd[1] = true
			z.Addr, bts, err = nbs.ReadStringBytes(bts)

			if err != nil {
				return
			}
		case "Host":
			found1zitd[2] = true
			z.Host, bts, err = nbs.ReadStringBytes(bts)

			if err != nil {
				return
			}
		case "Port":
			found1zitd[3] = true
			z.Port, bts, err = nbs.ReadIntBytes(bts)

			if err != nil {
				return
			}
		case "Parent":
			found1zitd[4] = true
			z.Parent, bts, err = nbs.ReadStringBytes(bts)

			if err != nil {
				return
			}
		case "Child":
			found1zitd[5] = true
			z.Child, bts, err = nbs.ReadStringBytes(bts)

			if err != nil {
				return
			}
		case "Nickname":
			found1zitd[6] = true
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
	if nextMiss1zitd != -1 {
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
var unmarshalMsgFieldOrder1zitd = []string{"Id", "Addr", "Host", "Port", "Parent", "Child", "Nickname"}

var unmarshalMsgFieldSkip1zitd = []bool{false, false, false, false, false, false, false}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *NodeInfo) Msgsize() (s int) {
	s = 1 + 3 + msgp.StringPrefixSize + len(z.Id) + 5 + msgp.StringPrefixSize + len(z.Addr) + 5 + msgp.StringPrefixSize + len(z.Host) + 5 + msgp.IntSize + 7 + msgp.StringPrefixSize + len(z.Parent) + 6 + msgp.StringPrefixSize + len(z.Child) + 9 + msgp.StringPrefixSize + len(z.Nickname)
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
	const maxFields2zrfx = 6

	// -- templateDecodeMsg starts here--
	var totalEncodedFields2zrfx uint32
	totalEncodedFields2zrfx, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	encodedFieldsLeft2zrfx := totalEncodedFields2zrfx
	missingFieldsLeft2zrfx := maxFields2zrfx - totalEncodedFields2zrfx

	var nextMiss2zrfx int32 = -1
	var found2zrfx [maxFields2zrfx]bool
	var curField2zrfx string

doneWithStruct2zrfx:
	// First fill all the encoded fields, then
	// treat the remaining, missing fields, as Nil.
	for encodedFieldsLeft2zrfx > 0 || missingFieldsLeft2zrfx > 0 {
		//fmt.Printf("encodedFieldsLeft: %v, missingFieldsLeft: %v, found: '%v', fields: '%#v'\n", encodedFieldsLeft2zrfx, missingFieldsLeft2zrfx, msgp.ShowFound(found2zrfx[:]), decodeMsgFieldOrder2zrfx)
		if encodedFieldsLeft2zrfx > 0 {
			encodedFieldsLeft2zrfx--
			field, err = dc.ReadMapKeyPtr()
			if err != nil {
				return
			}
			curField2zrfx = msgp.UnsafeString(field)
		} else {
			//missing fields need handling
			if nextMiss2zrfx < 0 {
				// tell the reader to only give us Nils
				// until further notice.
				dc.PushAlwaysNil()
				nextMiss2zrfx = 0
			}
			for nextMiss2zrfx < maxFields2zrfx && (found2zrfx[nextMiss2zrfx] || decodeMsgFieldSkip2zrfx[nextMiss2zrfx]) {
				nextMiss2zrfx++
			}
			if nextMiss2zrfx == maxFields2zrfx {
				// filled all the empty fields!
				break doneWithStruct2zrfx
			}
			missingFieldsLeft2zrfx--
			curField2zrfx = decodeMsgFieldOrder2zrfx[nextMiss2zrfx]
		}
		//fmt.Printf("switching on curField: '%v'\n", curField2zrfx)
		switch curField2zrfx {
		// -- templateDecodeMsg ends here --

		case "Num":
			found2zrfx[0] = true
			{
				var zcmn int
				zcmn, err = dc.ReadInt()
				z.Num = NoteEvt(zcmn)
			}
			if err != nil {
				return
			}
		case "From":
			found2zrfx[1] = true
			err = z.From.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "To":
			found2zrfx[2] = true
			err = z.To.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "ChainInfo":
			found2zrfx[3] = true
			var zjql uint32
			zjql, err = dc.ReadMapHeader()
			if err != nil {
				return
			}
			if z.ChainInfo == nil && zjql > 0 {
				z.ChainInfo = make(map[string]NodeInfo, zjql)
			} else if len(z.ChainInfo) > 0 {
				for key, _ := range z.ChainInfo {
					delete(z.ChainInfo, key)
				}
			}
			for zjql > 0 {
				zjql--
				var zqid string
				var zzce NodeInfo
				zqid, err = dc.ReadString()
				if err != nil {
					return
				}
				err = zzce.DecodeMsg(dc)
				if err != nil {
					return
				}
				z.ChainInfo[zqid] = zzce
			}
		case "SendTm":
			found2zrfx[4] = true
			z.SendTm, err = dc.ReadTime()
			if err != nil {
				return
			}
		case "Nonce":
			found2zrfx[5] = true
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
	if nextMiss2zrfx != -1 {
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
var decodeMsgFieldOrder2zrfx = []string{"Num", "From", "To", "ChainInfo", "SendTm", "Nonce"}

var decodeMsgFieldSkip2zrfx = []bool{false, false, false, false, false, false}

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
	for zqid, zzce := range z.ChainInfo {
		err = en.WriteString(zqid)
		if err != nil {
			return
		}
		err = zzce.EncodeMsg(en)
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
	for zqid, zzce := range z.ChainInfo {
		o = msgp.AppendString(o, zqid)
		o, err = zzce.MarshalMsg(o)
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
	const maxFields3zsjk = 6

	// -- templateUnmarshalMsg starts here--
	var totalEncodedFields3zsjk uint32
	if !nbs.AlwaysNil {
		totalEncodedFields3zsjk, bts, err = nbs.ReadMapHeaderBytes(bts)
		if err != nil {
			return
		}
	}
	encodedFieldsLeft3zsjk := totalEncodedFields3zsjk
	missingFieldsLeft3zsjk := maxFields3zsjk - totalEncodedFields3zsjk

	var nextMiss3zsjk int32 = -1
	var found3zsjk [maxFields3zsjk]bool
	var curField3zsjk string

doneWithStruct3zsjk:
	// First fill all the encoded fields, then
	// treat the remaining, missing fields, as Nil.
	for encodedFieldsLeft3zsjk > 0 || missingFieldsLeft3zsjk > 0 {
		//fmt.Printf("encodedFieldsLeft: %v, missingFieldsLeft: %v, found: '%v', fields: '%#v'\n", encodedFieldsLeft3zsjk, missingFieldsLeft3zsjk, msgp.ShowFound(found3zsjk[:]), unmarshalMsgFieldOrder3zsjk)
		if encodedFieldsLeft3zsjk > 0 {
			encodedFieldsLeft3zsjk--
			field, bts, err = nbs.ReadMapKeyZC(bts)
			if err != nil {
				return
			}
			curField3zsjk = msgp.UnsafeString(field)
		} else {
			//missing fields need handling
			if nextMiss3zsjk < 0 {
				// set bts to contain just mnil (0xc0)
				bts = nbs.PushAlwaysNil(bts)
				nextMiss3zsjk = 0
			}
			for nextMiss3zsjk < maxFields3zsjk && (found3zsjk[nextMiss3zsjk] || unmarshalMsgFieldSkip3zsjk[nextMiss3zsjk]) {
				nextMiss3zsjk++
			}
			if nextMiss3zsjk == maxFields3zsjk {
				// filled all the empty fields!
				break doneWithStruct3zsjk
			}
			missingFieldsLeft3zsjk--
			curField3zsjk = unmarshalMsgFieldOrder3zsjk[nextMiss3zsjk]
		}
		//fmt.Printf("switching on curField: '%v'\n", curField3zsjk)
		switch curField3zsjk {
		// -- templateUnmarshalMsg ends here --

		case "Num":
			found3zsjk[0] = true
			{
				var zmlp int
				zmlp, bts, err = nbs.ReadIntBytes(bts)

				if err != nil {
					return
				}
				z.Num = NoteEvt(zmlp)
			}
		case "From":
			found3zsjk[1] = true
			bts, err = z.From.UnmarshalMsg(bts)
			if err != nil {
				return
			}
			if err != nil {
				return
			}
		case "To":
			found3zsjk[2] = true
			bts, err = z.To.UnmarshalMsg(bts)
			if err != nil {
				return
			}
			if err != nil {
				return
			}
		case "ChainInfo":
			found3zsjk[3] = true
			if nbs.AlwaysNil {
				if len(z.ChainInfo) > 0 {
					for key, _ := range z.ChainInfo {
						delete(z.ChainInfo, key)
					}
				}

			} else {

				var zgti uint32
				zgti, bts, err = nbs.ReadMapHeaderBytes(bts)
				if err != nil {
					return
				}
				if z.ChainInfo == nil && zgti > 0 {
					z.ChainInfo = make(map[string]NodeInfo, zgti)
				} else if len(z.ChainInfo) > 0 {
					for key, _ := range z.ChainInfo {
						delete(z.ChainInfo, key)
					}
				}
				for zgti > 0 {
					var zqid string
					var zzce NodeInfo
					zgti--
					zqid, bts, err = nbs.ReadStringBytes(bts)
					if err != nil {
						return
					}
					bts, err = zzce.UnmarshalMsg(bts)
					if err != nil {
						return
					}
					if err != nil {
						return
					}
					z.ChainInfo[zqid] = zzce
				}
			}
		case "SendTm":
			found3zsjk[4] = true
			z.SendTm, bts, err = nbs.ReadTimeBytes(bts)

			if err != nil {
				return
			}
		case "Nonce":
			found3zsjk[5] = true
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
	if nextMiss3zsjk != -1 {
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
var unmarshalMsgFieldOrder3zsjk = []string{"Num", "From", "To", "ChainInfo", "SendTm", "Nonce"}

var unmarshalMsgFieldSkip3zsjk = []bool{false, false, false, false, false, false}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *Note) Msgsize() (s int) {
	s = 1 + 4 + msgp.IntSize + 5 + z.From.Msgsize() + 3 + z.To.Msgsize() + 10 + msgp.MapHeaderSize
	if z.ChainInfo != nil {
		for zqid, zzce := range z.ChainInfo {
			_ = zzce
			_ = zqid
			s += msgp.StringPrefixSize + len(zqid) + zzce.Msgsize()
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
		var zjyo int
		zjyo, err = dc.ReadInt()
		(*z) = NoteEvt(zjyo)
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
		var zuwk int
		zuwk, bts, err = nbs.ReadIntBytes(bts)

		if err != nil {
			return
		}
		(*z) = NoteEvt(zuwk)
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
