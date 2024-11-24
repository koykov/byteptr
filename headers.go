package byteptr

type bheader struct {
	ptr  uintptr
	l, c int
}

type sheader struct {
	ptr uintptr
	l   int
}
