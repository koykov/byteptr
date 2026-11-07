package byteptr

import "unsafe"

type KVPointer struct {
	addr uintptr
	max,
	koffset, klen int
	voffset, vlen int
}

func InitKV(s []byte, offset, len int) (p KVPointer) {
	p.KVInit(s, offset, len)
	return
}

func InitStringKV(s string, offset, len int) (p KVPointer) {
	p.KVInitString(s, offset, len)
	return
}

func (p *KVPointer) TakeAddress(s []byte) *KVPointer {
	if s == nil {
		return p
	}
	h := (*bheader)(unsafe.Pointer(&s))
	p.addr = h.ptr
	p.max = h.c
	return p
}

func (p *KVPointer) TakeStringAddress(s string) *KVPointer {
	if len(s) == 0 {
		return p
	}
	h := (*sheader)(unsafe.Pointer(&s))
	p.addr = h.ptr
	p.max = h.l
	return p
}

func (p *KVPointer) KVInit(s []byte, offset, len int) *KVPointer {
	p.TakeAddress(s).SetOffsetK(offset).SetLenK(len)
	return p
}

func (p *KVPointer) KVInitString(s string, offset, len int) *KVPointer {
	p.TakeStringAddress(s).SetOffsetK(offset).SetLenK(len)
	return p
}

func (p *KVPointer) SetOffsetK(offset int) *KVPointer {
	if offset <= p.max {
		p.koffset = offset
	}
	return p
}

func (p *KVPointer) OffsetK() int {
	return p.koffset
}

func (p *KVPointer) SetLenK(len int) *KVPointer {
	if len <= p.max {
		p.klen = len
	}
	return p
}

func (p *KVPointer) LenK() int {
	return p.klen
}

func (p *KVPointer) BytesK() []byte {
	if p.addr == 0 || p.koffset < 0 || p.klen < 0 {
		return nil
	}
	h := bheader{
		ptr: p.addr + uintptr(p.koffset),
		l:   p.klen,
		c:   p.klen,
	}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func (p *KVPointer) kString() string {
	if p.addr == 0 || p.koffset < 0 || p.klen < 0 {
		return ""
	}
	h := sheader{
		ptr: p.addr + uintptr(p.koffset),
		l:   p.klen,
	}
	return *(*string)(unsafe.Pointer(&h))
}

func (p *KVPointer) Reset() *KVPointer {
	p.addr, p.max, p.koffset, p.klen, p.koffset, p.klen = 0, 0, 0, 0, 0, 0
	return p
}

var _, _ = InitKV, InitStringKV
