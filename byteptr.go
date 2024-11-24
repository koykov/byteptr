package byteptr

import "unsafe"

// Byteptr is a pointer-free representation of bytes/string types.
//
// Similar to reflect.SliceHeader/reflect.StringHeader structs.
type Byteptr struct {
	addr uintptr
	max,
	offset, len int
}

// Init makes new ptr and set up with given params.
func Init(s []byte, offset, len int) Byteptr {
	p := Byteptr{}
	p.Init(s, offset, len)
	return p
}

// InitString makes new ptr and set up with given params.
func InitString(s string, offset, len int) Byteptr {
	p := Byteptr{}
	p.InitString(s, offset, len)
	return p
}

// TakeAddress takes address of underlying byte array of bytes s.
func (p *Byteptr) TakeAddress(s []byte) *Byteptr {
	if s == nil {
		return p
	}
	h := (*bheader)(unsafe.Pointer(&s))
	p.addr = h.ptr
	p.max = h.c
	return p
}

// TakeStringAddress takes address of underlying byte array of string s.
func (p *Byteptr) TakeStringAddress(s string) *Byteptr {
	if len(s) == 0 {
		return p
	}
	h := (*sheader)(unsafe.Pointer(&s))
	p.addr = h.ptr
	p.max = h.l
	return p
}

// Init ptr with address of the byte array s and offset/length.
func (p *Byteptr) Init(s []byte, offset, len int) *Byteptr {
	p.TakeAddress(s).SetOffset(offset).SetLen(len)
	return p
}

// InitString ptr with address of the string s and offset/length.
func (p *Byteptr) InitString(s string, offset, len int) *Byteptr {
	p.TakeStringAddress(s).SetOffset(offset).SetLen(len)
	return p
}

// SetOffset sets offset from previously taken address.
func (p *Byteptr) SetOffset(offset int) *Byteptr {
	if offset <= p.max {
		p.offset = offset
	}
	return p
}

// Offset returns offset.
func (p *Byteptr) Offset() int {
	return p.offset
}

// SetLen sets length of sub-slice.
func (p *Byteptr) SetLen(len int) *Byteptr {
	if len <= p.max {
		p.len = len
	}
	return p
}

func (p *Byteptr) Len() int {
	return p.len
}

// Bytes returns byte sub-slice using offset from previously take address with length len.
func (p *Byteptr) Bytes() []byte {
	if p.addr == 0 || p.offset < 0 || p.len < 0 {
		return nil
	}
	h := bheader{
		ptr: p.addr + uintptr(p.offset),
		l:   p.len,
		c:   p.len,
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
	h := sheader{
		ptr: p.addr + uintptr(p.offset),
		l:   p.len,
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

var _, _ = Init, InitString
