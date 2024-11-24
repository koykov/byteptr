package byteptr

import (
	"unsafe"
)

// InitStr is a deprecated version of InitString.
// DEPRECATED: use InitString instead.
func InitStr(s string, offset, len int) Byteptr {
	p := Byteptr{}
	p.InitString(s, offset, len)
	return p
}

// TakeAddr is a deprecated version of TakeAddress.
// DEPRECATED: use TakeAddress instead.
func (p *Byteptr) TakeAddr(s []byte) *Byteptr {
	if s == nil {
		return p
	}
	h := (*bheader)(unsafe.Pointer(&s))
	p.addr = h.ptr
	p.max = h.c
	return p
}

// TakeStrAddr is a deprecated version of TakeStringAddress.
// DEPRECATED: use TakeStringAddress instead.
func (p *Byteptr) TakeStrAddr(s string) *Byteptr {
	if len(s) == 0 {
		return p
	}
	h := (*sheader)(unsafe.Pointer(&s))
	p.addr = h.ptr
	p.max = h.l
	return p
}

// InitStr is a deprecated version of InitString.
// DEPRECATED: use InitString instead.
func (p *Byteptr) InitStr(s string, offset, len int) *Byteptr {
	p.TakeStringAddress(s).SetOffset(offset).SetLen(len)
	return p
}
