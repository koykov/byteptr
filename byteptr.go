package byteptr

import (
	"reflect"
	"unsafe"
)

type Byteptr struct {
	addr uint64
	max,
	offset, len int
}

func (p *Byteptr) TakeAddr(s []byte) *Byteptr {
	if s == nil {
		return p
	}
	h := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	p.addr = uint64(h.Data)
	p.max = h.Cap
	return p
}

func (p *Byteptr) TakeStrAddr(s string) *Byteptr {
	if len(s) == 0 {
		return p
	}
	h := (*reflect.StringHeader)(unsafe.Pointer(&s))
	p.addr = uint64(h.Data)
	p.max = h.Len
	return p
}

func (p *Byteptr) Init(s []byte, offset, len int) *Byteptr {
	p.TakeAddr(s).SetOffset(offset).SetLen(len)
	return p
}

func (p *Byteptr) InitStr(s string, offset, len int) *Byteptr {
	p.TakeStrAddr(s).SetOffset(offset).SetLen(len)
	return p
}

func (p *Byteptr) SetOffset(offset int) *Byteptr {
	if offset <= p.max {
		p.offset = offset
	}
	return p
}

func (p *Byteptr) Offset() int {
	return p.offset
}

func (p *Byteptr) SetLen(len int) *Byteptr {
	if len <= p.max {
		p.len = len
	}
	return p
}

func (p *Byteptr) Len() int {
	return p.len
}

func (p *Byteptr) Bytes() []byte {
	if p.addr == 0 || p.offset < 0 || p.len < 0 {
		return nil
	}
	h := reflect.SliceHeader{
		Data: uintptr(p.addr + uint64(p.offset)),
		Len:  p.len,
		Cap:  p.len,
	}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func (p *Byteptr) String() string {
	if p.addr == 0 || p.offset < 0 || p.len < 0 {
		return ""
	}
	h := reflect.StringHeader{
		Data: uintptr(p.addr + uint64(p.offset)),
		Len:  p.len,
	}
	return *(*string)(unsafe.Pointer(&h))
}

func (p *Byteptr) Reset() *Byteptr {
	p.addr, p.offset, p.len = 0, 0, 0
	return p
}
