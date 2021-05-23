package byteptr

import (
	"reflect"
	"unsafe"
)

type Byteptr struct {
	addr uint64
	offset, limit int
}

func (p *Byteptr) TakeAddr(s []byte) {
	if s == nil {
		return
	}
	h := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	p.addr = uint64(h.Data)
}

func (p *Byteptr) TakeStrAddr(s string) {
	if len(s) == 0 {
		return
	}
	h := (*reflect.StringHeader)(unsafe.Pointer(&s))
	p.addr = uint64(h.Data)
}

func (p *Byteptr) Init(offset, limit int) {
	p.offset, p.limit = offset, limit
}

func (p *Byteptr) SetOffset(offset int) {
	p.offset = offset
}

func (p *Byteptr) Offset() int {
	return p.offset
}

func (p *Byteptr) SetLimit(limit int) {
	p.limit = limit
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

func (p *Byteptr) Reset() {
	p.addr = 0
	p.Init(0, 0)
}
