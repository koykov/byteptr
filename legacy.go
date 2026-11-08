package byteptr

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
	return p.TakeAddress(s)
}

// TakeStrAddr is a deprecated version of TakeStringAddress.
// DEPRECATED: use TakeStringAddress instead.
func (p *Pointer) TakeStrAddr(s string) *Pointer {
	return p.TakeStringAddress(s)
}

// InitStr is a deprecated version of InitString.
// DEPRECATED: use InitString instead.
func (p *Pointer) InitStr(s string, offset, len int) *Pointer {
	return p.InitString(s, offset, len)
}
