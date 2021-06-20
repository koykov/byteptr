package byteptr

import (
	"reflect"
	"unsafe"
)

// Byteptr is a pointer-free representation of bytes/string types.
//
// Similar to reflect.SliceHeader/reflect.StringHeader structs.
type Byteptr struct {
	addr uintptr
	max,
	offset, len int
}

// Take address of underlying byte array of bytes s.
func (p *Byteptr) TakeAddr(s []byte) *Byteptr {
	if s == nil {
		return p
	}
	h := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	p.addr = h.Data
	p.max = h.Cap
	return p
}

// Take address of underlying byte array of string s.
func (p *Byteptr) TakeStrAddr(s string) *Byteptr {
	if len(s) == 0 {
		return p
	}
	h := (*reflect.StringHeader)(unsafe.Pointer(&s))
	p.addr = h.Data
	p.max = h.Len
	return p
}

// Init p with address of the byte array s and offset/length.
func (p *Byteptr) Init(s []byte, offset, len int) *Byteptr {
	p.TakeAddr(s).SetOffset(offset).SetLen(len)
	return p
}

// Init p with address of the string s and offset/length.
func (p *Byteptr) InitStr(s string, offset, len int) *Byteptr {
	p.TakeStrAddr(s).SetOffset(offset).SetLen(len)
	return p
}

// Set offset from previously taken address.
func (p *Byteptr) SetOffset(offset int) *Byteptr {
	if offset <= p.max {
		p.offset = offset
	}
	return p
}

// Get offset.
func (p *Byteptr) Offset() int {
	return p.offset
}

// Set length of sub-slice.
func (p *Byteptr) SetLen(len int) *Byteptr {
	if len <= p.max {
		p.len = len
	}
	return p
}

// Get length.
func (p *Byteptr) Len() int {
	return p.len
}

// Get byte sub-slice using offset from previously take address with length len.
func (p *Byteptr) Bytes() []byte {
	if p.addr == 0 || p.offset < 0 || p.len < 0 {
		return nil
	}
	h := reflect.SliceHeader{
		Data: p.addr + uintptr(p.offset),
		Len:  p.len,
		Cap:  p.len,
	}
	return *(*[]byte)(unsafe.Pointer(&h))
}

// Get substring.
//
// See Bytes().
func (p *Byteptr) String() string {
	if p.addr == 0 || p.offset < 0 || p.len < 0 {
		return ""
	}
	h := reflect.StringHeader{
		Data: p.addr + uintptr(p.offset),
		Len:  p.len,
	}
	return *(*string)(unsafe.Pointer(&h))
}

// Reset all fields.
//
// Made to use with pools.
func (p *Byteptr) Reset() *Byteptr {
	p.addr, p.max, p.offset, p.len = 0, 0, 0, 0
	return p
}
