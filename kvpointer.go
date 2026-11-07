package byteptr

import "unsafe"

type KVPointer struct {
	addr uintptr
	max,
	offsetK, lenK int
	offsetV, lenV int
}

func InitKV(s []byte, offset, len int) (p KVPointer) {
	p.InitKV(s, offset, len)
	return
}

func InitStringKV(s string, offset, len int) (p KVPointer) {
	p.InitStringKV(s, offset, len)
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

func (p *KVPointer) InitKV(s []byte, offset, len int) *KVPointer {
	p.TakeAddress(s).SetOffsetK(offset).SetLenK(len)
	return p
}

func (p *KVPointer) InitStringKV(s string, offset, len int) *KVPointer {
	p.TakeStringAddress(s).SetOffsetK(offset).SetLenK(len)
	return p
}

func (p *KVPointer) SetOffsetK(offset int) *KVPointer {
	if offset <= p.max {
		p.offsetK = offset
	}
	return p
}

func (p *KVPointer) OffsetK() int {
	return p.offsetK
}

func (p *KVPointer) SetLenK(len int) *KVPointer {
	if len <= p.max {
		p.lenK = len
	}
	return p
}

func (p *KVPointer) LenK() int {
	return p.lenK
}

func (p *KVPointer) BytesK() []byte {
	if p.addr == 0 || p.offsetK < 0 || p.lenK < 0 {
		return nil
	}
	h := bheader{
		ptr: p.addr + uintptr(p.offsetK),
		l:   p.lenK,
		c:   p.lenK,
	}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func (p *KVPointer) kString() string {
	if p.addr == 0 || p.offsetK < 0 || p.lenK < 0 {
		return ""
	}
	h := sheader{
		ptr: p.addr + uintptr(p.offsetK),
		l:   p.lenK,
	}
	return *(*string)(unsafe.Pointer(&h))
}

func (p *KVPointer) Reset() *KVPointer {
	p.addr, p.max, p.offsetK, p.lenK, p.offsetK, p.lenK = 0, 0, 0, 0, 0, 0
	return p
}

var _, _ = InitKV, InitStringKV
