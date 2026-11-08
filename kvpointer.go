package byteptr

import "unsafe"

type KVPointer struct {
	addr uintptr
	max,
	offsetK, lenK int
	offsetV, lenV int
}

func InitKV(s []byte, offsetK, lenK, offsetV, lenV int) (p KVPointer) {
	p.InitKV(s, offsetK, lenK, offsetV, lenV)
	return
}

func InitStringKV(s string, offsetK, lenK, offsetV, lenV int) (p KVPointer) {
	p.InitStringKV(s, offsetK, lenK, offsetV, lenV)
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

func (p *KVPointer) InitKV(s []byte, offsetK, lenK, offsetV, lenV int) *KVPointer {
	if s == nil {
		return p.Reset()
	}
	p.TakeAddress(s)
	p.offsetK, p.lenK = 0, 0
	if offsetK >= 0 && lenK >= 0 && offsetK <= p.max && lenK <= p.max-offsetK {
		p.offsetK, p.lenK = offsetK, lenK
	}
	p.offsetV, p.lenV = 0, 0
	if offsetV >= 0 && lenV >= 0 && offsetV <= p.max && lenV <= p.max-offsetV {
		p.offsetV, p.lenV = offsetV, lenV
	}
	return p
}

func (p *KVPointer) InitStringKV(s string, offsetK, lenK, offsetV, lenV int) *KVPointer {
	if len(s) == 0 {
		return p.Reset()
	}
	p.TakeStringAddress(s)
	p.offsetK, p.lenK = 0, 0
	if offsetK >= 0 && lenK >= 0 && offsetK <= p.max && lenK <= p.max-offsetK {
		p.offsetK, p.lenK = offsetK, lenK
	}
	p.offsetV, p.lenV = 0, 0
	if offsetV >= 0 && lenV >= 0 && offsetV <= p.max && lenV <= p.max-offsetV {
		p.offsetV, p.lenV = offsetV, lenV
	}
	return p
}

func (p *KVPointer) SetOffsetK(offset int) *KVPointer {
	if offset >= 0 && offset <= p.max {
		p.offsetK = offset
	}
	return p
}

func (p *KVPointer) OffsetK() int {
	return p.offsetK
}

func (p *KVPointer) SetLenK(len int) *KVPointer {
	if len >= 0 && len <= p.max-p.offsetK {
		p.lenK = len
	}
	return p
}

func (p *KVPointer) LenK() int {
	return p.lenK
}

func (p *KVPointer) SetOffsetV(offset int) *KVPointer {
	if offset >= 0 && offset <= p.max {
		p.offsetV = offset
	}
	return p
}

func (p *KVPointer) OffsetV() int {
	return p.offsetV
}

func (p *KVPointer) SetLenV(len int) *KVPointer {
	if len >= 0 && len <= p.max-p.offsetV {
		p.lenV = len
	}
	return p
}

func (p *KVPointer) LenV() int {
	return p.lenV
}

func (p *KVPointer) BytesK() []byte {
	if p.addr == 0 || p.offsetK < 0 || p.lenK < 0 || p.lenK > p.max-p.offsetK {
		return nil
	}
	h := bheader{
		ptr: p.addr + uintptr(p.offsetK),
		l:   p.lenK,
		c:   p.lenK,
	}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func (p *KVPointer) StringK() string {
	if p.addr == 0 || p.offsetK < 0 || p.lenK < 0 || p.lenK > p.max-p.offsetK {
		return ""
	}
	h := sheader{
		ptr: p.addr + uintptr(p.offsetK),
		l:   p.lenK,
	}
	return *(*string)(unsafe.Pointer(&h))
}

func (p *KVPointer) BytesV() []byte {
	if p.addr == 0 || p.offsetV < 0 || p.lenV < 0 || p.lenV > p.max-p.offsetV {
		return nil
	}
	h := bheader{
		ptr: p.addr + uintptr(p.offsetV),
		l:   p.lenV,
		c:   p.lenV,
	}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func (p *KVPointer) StringV() string {
	if p.addr == 0 || p.offsetV < 0 || p.lenV < 0 || p.lenV > p.max-p.offsetV {
		return ""
	}
	h := sheader{
		ptr: p.addr + uintptr(p.offsetV),
		l:   p.lenV,
	}
	return *(*string)(unsafe.Pointer(&h))
}

func (p *KVPointer) Reset() *KVPointer {
	p.addr, p.max, p.offsetK, p.lenK, p.offsetV, p.lenV = 0, 0, 0, 0, 0, 0
	return p
}

var _, _ = InitKV, InitStringKV
