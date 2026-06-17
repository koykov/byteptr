package byteptr

import (
	"unsafe"
)

// Byteptr is a legacy alias of Pointer.
// Deprecated: use Pointer instead.
type Byteptr = Pointer

// InitStr is a deprecated version of InitString.
// DEPRECATED: use InitString instead.
func InitStr(s string, offset, len int) Pointer {
	p := Pointer{}
	p.InitString(s, offset, len)
	return p
}

// TakeAddr is a deprecated version of TakeAddress.
// DEPRECATED: use TakeAddress instead.
func (p *Pointer) TakeAddr(s []byte) *Pointer {
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
func (p *Pointer) TakeStrAddr(s string) *Pointer {
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
func (p *Pointer) InitStr(s string, offset, len int) *Pointer {
	p.TakeStringAddress(s).SetOffset(offset).SetLen(len)
	return p
}
