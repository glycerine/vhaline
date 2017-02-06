package vhaline

import (
	cr "crypto/rand"
	"encoding/binary"
	mr "math/rand"
)

func mathRandHexString(n int, rsrc *mr.Rand) string {
	by := mathRandBytes(rsrc, n/2+1)

	m := len(by)
	p := m - 1
	t := p * 2
	res := make([]rune, m*2)
	for i := 0; i < m; i++ {
		r := byte(by[p-i])
		res[t-i*2] = encode16(r >> 4)
		res[t-i*2+1] = encode16(r & 0x0F)
	}

	// we could be 1 byte larger than need, so
	// truncate here
	s := string(res)[:n]
	//fmt.Printf("source bytes: %x\n", by)
	//fmt.Printf("hexstring conversion: %s\n", s)
	return s
}

func mathRandBytes(rsrc *mr.Rand, n int) []byte {
	b := make([]byte, n)
	_, err := rsrc.Read(b)
	if err != nil {
		panic(err)
	}
	return b
}

// Use crypto/rand to get an random int64.
func cryptoRandInt64() int64 {
	b := make([]byte, 8)
	_, err := cr.Read(b)
	if err != nil {
		panic(err)
	}
	r := int64(binary.LittleEndian.Uint64(b))
	return r
}

var enc16 string = "0123456789abcdef"
var e16 []rune = []rune(enc16)

// nibble must be between 0 and 15 inclusive.
func encode16(nibble byte) rune {
	return e16[nibble]
}
