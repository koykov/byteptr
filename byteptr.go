package byteptr

import (
	"reflect"
	"unsafe"
)

type Byteptr struct {
	addr          uint64
	max           int
	offset, limit int
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

func (p *Byteptr) Init(s []byte, offset, limit int) *Byteptr {
	p.TakeAddr(s).SetOffset(offset).SetLimit(limit)
	return p
}

func (p *Byteptr) InitStr(s string, offset, limit int) *Byteptr {
	p.TakeStrAddr(s).SetOffset(offset).SetLimit(limit)
	return p
}

func (p *Byteptr) SetOffset(offset int) *Byteptr {
	if offset < p.max {
		p.offset = offset
	}
	return p
}

func (p *Byteptr) Offset() int {
	return p.offset
}

func (p *Byteptr) SetLimit(limit int) *Byteptr {
	if limit < p.max {
		p.limit = limit
	}
	return p
}

func (p *Byteptr) Limit() int {
	return p.limit
}

func (p *Byteptr) Bytes() []byte {
	if p.addr == 0 || p.offset < 0 || p.limit < 0 {
		return nil
	}
	h := reflect.SliceHeader{
		Data: uintptr(p.offset),
		Len:  p.limit,
		Cap:  p.limit,
	}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func (p *Byteptr) String() string {
	if p.addr == 0 || p.offset < 0 || p.limit < 0 {
		return ""
	}
	h := reflect.StringHeader{
		Data: uintptr(p.offset),
		Len:  p.limit,
	}
	return *(*string)(unsafe.Pointer(&h))
}

func (p *Byteptr) Reset() *Byteptr {
	p.addr, p.offset, p.limit = 0, 0, 0
	return p
}
