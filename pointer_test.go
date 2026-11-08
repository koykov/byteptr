package byteptr_test

import (
	"runtime"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/koykov/byteptr"
)

func TestInit(t *testing.T) {
	tests := []struct {
		name       string
		src        []byte
		offset     int
		length     int
		wantNil    bool
		want       string
		wantOffset int
		wantLen    int
	}{
		{
			name: "full window", src: []byte("hello"), offset: 0, length: 5,
			want: "hello", wantOffset: 0, wantLen: 5,
		},
		{
			name: "zero length at start", src: []byte("hello"), offset: 0, length: 0,
			want: "", wantOffset: 0, wantLen: 0,
		},
		{
			name: "prefix", src: []byte("hello"), offset: 0, length: 2,
			want: "he", wantOffset: 0, wantLen: 2,
		},
		{
			name: "middle window", src: []byte("hello"), offset: 1, length: 3,
			want: "ell", wantOffset: 1, wantLen: 3,
		},
		{
			name: "offset three length two", src: []byte("hello"), offset: 3, length: 2,
			want: "lo", wantOffset: 3, wantLen: 2,
		},
		{
			name: "last byte", src: []byte("hello"), offset: 4, length: 1,
			want: "o", wantOffset: 4, wantLen: 1,
		},
		{
			name: "offset at max empty window", src: []byte("hello"), offset: 5, length: 0,
			want: "", wantOffset: 5, wantLen: 0,
		},
		{
			name: "tail window", src: []byte("hello"), offset: 1, length: 4,
			want: "ello", wantOffset: 1, wantLen: 4,
		},
		{
			name: "offset plus length exceeds max", src: []byte("hello"), offset: 1, length: 5,
			want: "", wantOffset: 0, wantLen: 0,
		},
		{
			name: "offset beyond max", src: []byte("hello"), offset: 6, length: 5,
			want: "", wantOffset: 0, wantLen: 0,
		},
		{
			name: "offset beyond max with zero length", src: []byte("hello"), offset: 6, length: 0,
			want: "", wantOffset: 0, wantLen: 0,
		},
		{
			name: "negative offset", src: []byte("hello"), offset: -1, length: 2,
			want: "", wantOffset: 0, wantLen: 0,
		},
		{
			name: "negative length", src: []byte("hello"), offset: 0, length: -1,
			want: "", wantOffset: 0, wantLen: 0,
		},
		{
			name: "nil slice", src: nil, offset: 0, length: 0,
			wantNil: true, wantOffset: 0, wantLen: 0,
		},
		{
			name: "nil slice with non-zero window", src: nil, offset: 0, length: 5,
			wantNil: true, wantOffset: 0, wantLen: 0,
		},
		{
			name: "empty non-nil slice", src: []byte{}, offset: 0, length: 0,
			want: "", wantOffset: 0, wantLen: 0,
		},
		{
			name: "empty non-nil slice with non-zero length", src: []byte{}, offset: 0, length: 1,
			want: "", wantOffset: 0, wantLen: 0,
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			p := byteptr.Init(tc.src, tc.offset, tc.length)
			assert.Equal(t, tc.wantOffset, p.Offset())
			assert.Equal(t, tc.wantLen, p.Len())

			got := p.Bytes()
			if tc.wantNil {
				assert.Nil(t, got)
			} else {
				require.NotNil(t, got)
				assert.Equal(t, tc.want, string(got))
				assert.Equal(t, len(got), cap(got))
			}
			assert.Equal(t, tc.want, p.String())
		})
	}

	t.Run("reused pointer resets on nil source", func(t *testing.T) {
		p := byteptr.Init(mutableBytes("hello"), 0, 0)
		p.Init(nil, 0, 5)
		assert.Nil(t, p.Bytes())
		assert.Equal(t, 0, p.Offset())
		assert.Equal(t, 0, p.Len())
		assert.Equal(t, "", p.String())
	})

	t.Run("fresh pointer with nil source yields nil bytes", func(t *testing.T) {
		p := &byteptr.Pointer{}
		p.Init(nil, 0, 0)
		assert.Nil(t, p.Bytes())
		assert.Equal(t, 0, p.Offset())
		assert.Equal(t, 0, p.Len())
		assert.Equal(t, "", p.String())
	})

	t.Run("empty non-nil source is not treated like nil", func(t *testing.T) {
		p := byteptr.Init([]byte{}, 0, 0)
		got := p.Bytes()
		require.NotNil(t, got)
		assert.Equal(t, 0, len(got))
		assert.Equal(t, 0, cap(got))
		assert.Equal(t, "", p.String())
	})

	t.Run("length may extend to capacity", func(t *testing.T) {
		src := make([]byte, 3, 10)
		copy(src, "abc")

		p := byteptr.Init(src, 0, 10)
		assert.Equal(t, 0, p.Offset())
		assert.Equal(t, 10, p.Len())

		got := p.Bytes()
		require.NotNil(t, got)
		require.Len(t, got, 10)
		assert.Equal(t, 10, cap(got))
		assert.Equal(t, "abc", string(got[:3]))
		for _, b := range got[3:] {
			assert.Equal(t, byte(0), b)
		}
		runtime.KeepAlive(src)
	})
}

func TestInitString(t *testing.T) {
	tests := []struct {
		name       string
		src        string
		offset     int
		length     int
		wantNil    bool
		want       string
		wantOffset int
		wantLen    int
	}{
		{
			name: "full window", src: "hello", offset: 0, length: 5,
			want: "hello", wantOffset: 0, wantLen: 5,
		},
		{
			name: "zero length at start", src: "hello", offset: 0, length: 0,
			want: "", wantOffset: 0, wantLen: 0,
		},
		{
			name: "middle window", src: "hello", offset: 1, length: 3,
			want: "ell", wantOffset: 1, wantLen: 3,
		},
		{
			name: "last byte", src: "hello", offset: 4, length: 1,
			want: "o", wantOffset: 4, wantLen: 1,
		},
		{
			name: "offset at max empty window", src: "hello", offset: 5, length: 0,
			want: "", wantOffset: 5, wantLen: 0,
		},
		{
			name: "offset at max with non-zero length", src: "hello", offset: 5, length: 1,
			want: "", wantOffset: 0, wantLen: 0,
		},
		{
			name: "tail window", src: "hello", offset: 1, length: 4,
			want: "ello", wantOffset: 1, wantLen: 4,
		},
		{
			name: "offset plus length exceeds max", src: "hello", offset: 1, length: 5,
			want: "", wantOffset: 0, wantLen: 0,
		},
		{
			name: "offset beyond max", src: "hello", offset: 6, length: 2,
			want: "", wantOffset: 0, wantLen: 0,
		},
		{
			name: "negative offset", src: "hello", offset: -1, length: 2,
			want: "", wantOffset: 0, wantLen: 0,
		},
		{
			name: "negative length", src: "hello", offset: 0, length: -1,
			want: "", wantOffset: 0, wantLen: 0,
		},
		{
			name: "empty string", src: "", offset: 0, length: 0,
			wantNil: true, wantOffset: 0, wantLen: 0,
		},
		{
			name: "empty string with non-zero window", src: "", offset: 0, length: 3,
			wantNil: true, wantOffset: 0, wantLen: 0,
		},
		{
			name: "empty string with offset beyond max", src: "", offset: 5, length: 0,
			wantNil: true, wantOffset: 0, wantLen: 0,
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			p := byteptr.InitString(tc.src, tc.offset, tc.length)
			assert.Equal(t, tc.wantOffset, p.Offset())
			assert.Equal(t, tc.wantLen, p.Len())

			got := p.Bytes()
			if tc.wantNil {
				assert.Nil(t, got)
			} else {
				require.NotNil(t, got)
				assert.Equal(t, tc.want, string(got))
				assert.Equal(t, len(got), cap(got))
			}
			assert.Equal(t, tc.want, p.String())
		})
	}

	t.Run("reused pointer resets on empty source", func(t *testing.T) {
		p := byteptr.InitString("hello", 0, 5)
		require.Equal(t, "hello", p.String())

		p.InitString("", 0, 5)

		assert.Equal(t, "", p.String())
		assert.Nil(t, p.Bytes())
		assert.Equal(t, 0, p.Offset())
		assert.Equal(t, 0, p.Len())
	})
}

func TestTakeAddress(t *testing.T) {
	tests := []struct {
		name  string
		check func(t *testing.T)
	}{
		{
			name: "nil is a no-op and preserves prior window",
			check: func(t *testing.T) {
				p := byteptr.Init([]byte("hello"), 1, 3)
				p.TakeAddress(nil)
				assert.Equal(t, 1, p.Offset())
				assert.Equal(t, 3, p.Len())
				assert.Equal(t, "ell", p.String())
			},
		},
		{
			name: "non-nil sets address and capacity",
			check: func(t *testing.T) {
				p := &byteptr.Pointer{}
				p.TakeAddress([]byte("world"))
				assert.Equal(t, 0, p.Offset())
				assert.Equal(t, 0, p.Len())
				p.SetOffset(4).SetLen(1)
				assert.Equal(t, "d", p.String())
			},
		},
		{
			name: "empty non-nil sets address with zero capacity",
			check: func(t *testing.T) {
				p := &byteptr.Pointer{}
				p.TakeAddress([]byte{})
				p.SetOffset(0).SetLen(0)
				got := p.Bytes()
				require.NotNil(t, got)
				assert.Equal(t, 0, len(got))
				assert.Equal(t, 0, cap(got))
				assert.Equal(t, "", p.String())
			},
		},
		{
			name: "preserves offset and length when re-taking a new buffer",
			check: func(t *testing.T) {
				p := byteptr.Init([]byte("hello"), 1, 3)
				p.TakeAddress([]byte("abcdefg"))
				assert.Equal(t, 1, p.Offset())
				assert.Equal(t, 3, p.Len())
				assert.Equal(t, "bcd", p.String())
			},
		},
		{
			name: "smaller buffer makes getters return nil",
			check: func(t *testing.T) {
				p := byteptr.Init([]byte("hello"), 3, 2)
				require.Equal(t, "lo", p.String())
				p.TakeAddress([]byte("ab"))
				assert.Equal(t, 3, p.Offset())
				assert.Equal(t, 2, p.Len())
				assert.Nil(t, p.Bytes())
				assert.Equal(t, "", p.String())
			},
		},
		{
			name: "empty non-nil buffer after a valid window makes getters return nil",
			check: func(t *testing.T) {
				p := byteptr.Init([]byte("hello"), 1, 3)
				require.Equal(t, "ell", p.String())
				p.TakeAddress([]byte{})
				assert.Nil(t, p.Bytes())
				assert.Equal(t, "", p.String())
			},
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			tc.check(t)
		})
	}
}

func TestTakeStringAddress(t *testing.T) {
	tests := []struct {
		name  string
		check func(t *testing.T)
	}{
		{
			name: "empty string on a fresh pointer is a no-op",
			check: func(t *testing.T) {
				p := &byteptr.Pointer{}
				p.TakeStringAddress("")
				assert.Equal(t, 0, p.Offset())
				assert.Equal(t, 0, p.Len())
				assert.Nil(t, p.Bytes())
				assert.Equal(t, "", p.String())
			},
		},
		{
			name: "empty string is a no-op and preserves prior window",
			check: func(t *testing.T) {
				p := byteptr.InitString("hello", 1, 3)
				p.TakeStringAddress("")
				assert.Equal(t, 1, p.Offset())
				assert.Equal(t, 3, p.Len())
				assert.Equal(t, "ell", p.String())
			},
		},
		{
			name: "non-empty sets address and capacity",
			check: func(t *testing.T) {
				p := &byteptr.Pointer{}
				p.TakeStringAddress("abcdef")
				p.SetOffset(2).SetLen(3)
				assert.Equal(t, "cde", p.String())
			},
		},
		{
			name: "preserves offset and length when re-taking a new string",
			check: func(t *testing.T) {
				p := byteptr.InitString("hello", 1, 3)
				p.TakeStringAddress("abcdefg")
				assert.Equal(t, 1, p.Offset())
				assert.Equal(t, 3, p.Len())
				assert.Equal(t, "bcd", p.String())
			},
		},
		{
			name: "smaller string makes getters return nil",
			check: func(t *testing.T) {
				p := byteptr.InitString("hello", 3, 2)
				require.Equal(t, "lo", p.String())
				p.TakeStringAddress("ab")
				assert.Equal(t, 3, p.Offset())
				assert.Equal(t, 2, p.Len())
				assert.Nil(t, p.Bytes())
				assert.Equal(t, "", p.String())
			},
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			tc.check(t)
		})
	}
}

func TestSetOffset(t *testing.T) {
	tests := []struct {
		name       string
		initOffset int
		initLen    int
		set        int
		want       int
	}{
		{name: "zero accepted", initOffset: 2, initLen: 2, set: 0, want: 0},
		{name: "in range accepted", initOffset: 2, initLen: 2, set: 4, want: 4},
		{name: "max accepted", initOffset: 0, initLen: 0, set: 5, want: 5},
		{name: "max plus one ignored", initOffset: 2, initLen: 2, set: 6, want: 2},
		{name: "negative ignored", initOffset: 2, initLen: 2, set: -1, want: 2},
		{name: "unchanged value accepted", initOffset: 3, initLen: 1, set: 3, want: 3},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			p := byteptr.Init([]byte("hello"), tc.initOffset, tc.initLen)
			p.SetOffset(tc.set)
			assert.Equal(t, tc.want, p.Offset())
		})
	}
}

func TestSetLen(t *testing.T) {
	tests := []struct {
		name       string
		initOffset int
		initLen    int
		apply      func(p *byteptr.Pointer)
		want       int
	}{
		{
			name: "zero accepted", initOffset: 2, initLen: 3,
			apply: func(p *byteptr.Pointer) { p.SetLen(0) }, want: 0,
		},
		{
			name: "max minus offset accepted", initOffset: 2, initLen: 0,
			apply: func(p *byteptr.Pointer) { p.SetLen(3) }, want: 3,
		},
		{
			name: "max minus offset plus one ignored", initOffset: 2, initLen: 1,
			apply: func(p *byteptr.Pointer) { p.SetLen(4) }, want: 1,
		},
		{
			name: "negative ignored", initOffset: 2, initLen: 1,
			apply: func(p *byteptr.Pointer) { p.SetLen(-1) }, want: 1,
		},
		{
			name: "zero at max offset accepted", initOffset: 5, initLen: 0,
			apply: func(p *byteptr.Pointer) { p.SetLen(0) }, want: 0,
		},
		{
			name: "one at max offset ignored", initOffset: 5, initLen: 0,
			apply: func(p *byteptr.Pointer) { p.SetLen(1) }, want: 0,
		},
		{
			name: "respects newly set offset", initOffset: 0, initLen: 0,
			apply: func(p *byteptr.Pointer) { p.SetOffset(3).SetLen(2) }, want: 2,
		},
		{
			name: "rejects length beyond newly set offset", initOffset: 0, initLen: 0,
			apply: func(p *byteptr.Pointer) { p.SetOffset(3).SetLen(3) }, want: 0,
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			p := &byteptr.Pointer{}
			p.Init([]byte("hello"), tc.initOffset, tc.initLen)
			tc.apply(p)
			assert.Equal(t, tc.want, p.Len())
		})
	}
}

func TestBytes(t *testing.T) {
	tests := []struct {
		name    string
		src     []byte
		offset  int
		length  int
		wantNil bool
		want    string
	}{
		{name: "nil slice", src: nil, offset: 0, length: 0, wantNil: true, want: ""},
		{name: "nil slice with window", src: nil, offset: 0, length: 5, wantNil: true, want: ""},
		{name: "empty non-nil slice", src: []byte{}, offset: 0, length: 0, wantNil: false, want: ""},
		{name: "full window", src: []byte("hello"), offset: 0, length: 5, wantNil: false, want: "hello"},
		{name: "middle window", src: []byte("hello"), offset: 1, length: 3, wantNil: false, want: "ell"},
		{name: "last byte", src: []byte("hello"), offset: 4, length: 1, wantNil: false, want: "o"},
		{name: "empty window at max", src: []byte("hello"), offset: 5, length: 0, wantNil: false, want: ""},
		{name: "invalid window becomes empty non-nil", src: []byte("hello"), offset: 1, length: 5, wantNil: false, want: ""},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			p := byteptr.Init(tc.src, tc.offset, tc.length)
			got := p.Bytes()
			if tc.wantNil {
				assert.Nil(t, got)
				return
			}
			require.NotNil(t, got)
			assert.Equal(t, tc.want, string(got))
			assert.Equal(t, len(got), cap(got))
		})
	}

	t.Run("zero value pointer returns nil", func(t *testing.T) {
		p := &byteptr.Pointer{}
		assert.Nil(t, p.Bytes())
	})

	t.Run("view mutation is visible in source", func(t *testing.T) {
		buf := mutableBytes("hello")
		p := byteptr.Init(buf, 1, 3)
		view := p.Bytes()
		require.NotNil(t, view)
		require.Len(t, view, 3)
		view[0] = 'X'
		assert.Equal(t, byte('X'), buf[1])
		assert.Equal(t, "hXllo", string(buf))
		runtime.KeepAlive(buf)
	})

	t.Run("source mutation is visible in view", func(t *testing.T) {
		buf := mutableBytes("hello")
		p := byteptr.Init(buf, 1, 3)
		view := p.Bytes()
		require.NotNil(t, view)
		require.Len(t, view, 3)
		buf[2] = 'Y'
		assert.Equal(t, "eYl", string(view))
		runtime.KeepAlive(buf)
	})

	t.Run("source mutation is visible through string", func(t *testing.T) {
		buf := mutableBytes("hello")
		p := byteptr.Init(buf, 1, 3)
		s := p.String()
		buf[1] = 'Z'
		assert.Equal(t, "Zll", s)
		runtime.KeepAlive(buf)
	})
}

func TestString(t *testing.T) {
	tests := []struct {
		name    string
		src     []byte
		offset  int
		length  int
		wantNil bool
		want    string
	}{
		{name: "nil slice", src: nil, offset: 0, length: 0, wantNil: true, want: ""},
		{name: "missing address with window", src: nil, offset: 0, length: 3, wantNil: true, want: ""},
		{name: "empty non-nil slice", src: []byte{}, offset: 0, length: 0, wantNil: false, want: ""},
		{name: "full window", src: []byte("hello"), offset: 0, length: 5, wantNil: false, want: "hello"},
		{name: "middle window", src: []byte("hello"), offset: 1, length: 3, wantNil: false, want: "ell"},
		{name: "empty window at max", src: []byte("hello"), offset: 5, length: 0, wantNil: false, want: ""},
		{name: "invalid window becomes empty", src: []byte("hello"), offset: 1, length: 5, wantNil: false, want: ""},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			p := byteptr.Init(tc.src, tc.offset, tc.length)
			assert.Equal(t, tc.want, p.String())

			got := p.Bytes()
			if tc.wantNil {
				assert.Nil(t, got)
				return
			}
			require.NotNil(t, got)
			assert.Equal(t, tc.want, string(got))
		})
	}

	t.Run("zero value pointer returns empty string", func(t *testing.T) {
		p := &byteptr.Pointer{}
		assert.Equal(t, "", p.String())
	})

	t.Run("source mutation is visible through previously obtained string", func(t *testing.T) {
		buf := mutableBytes("hello")
		p := byteptr.Init(buf, 1, 3)
		s := p.String()
		buf[1] = 'Q'
		assert.Equal(t, "Qll", s)
		runtime.KeepAlive(buf)
	})

	t.Run("bytes and string views share memory", func(t *testing.T) {
		buf := mutableBytes("hello")
		p := byteptr.Init(buf, 1, 3)
		s := p.String()

		view := p.Bytes()
		require.NotNil(t, view)
		require.Len(t, view, 3)
		view[0] = 'R'

		assert.Equal(t, "Rll", s)
		runtime.KeepAlive(buf)
	})
}

func TestReset(t *testing.T) {
	t.Run("clears an initialized pointer", func(t *testing.T) {
		p := byteptr.Init([]byte("hello"), 1, 3)
		require.Equal(t, "ell", p.String())

		p.Reset()

		assert.Equal(t, 0, p.Offset())
		assert.Equal(t, 0, p.Len())
		assert.Nil(t, p.Bytes())
		assert.Equal(t, "", p.String())
	})

	t.Run("reset on a zero value is safe", func(t *testing.T) {
		p := &byteptr.Pointer{}
		p.Reset()
		assert.Nil(t, p.Bytes())
		assert.Equal(t, "", p.String())
	})

	t.Run("no stale window after reset and re-take", func(t *testing.T) {
		p := byteptr.Init([]byte("hello"), 1, 3)
		p.Reset()
		p.TakeAddress([]byte("world"))
		assert.Equal(t, 0, p.Offset())
		assert.Equal(t, 0, p.Len())
		assert.Equal(t, "", p.String())
	})
}

func TestChaining(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(p *byteptr.Pointer) *byteptr.Pointer
		verify func(t *testing.T, p *byteptr.Pointer)
	}{
		{
			name: "TakeAddress",
			mutate: func(p *byteptr.Pointer) *byteptr.Pointer {
				p.Init([]byte("hello"), 1, 3)
				return p.TakeAddress([]byte("abcdefg"))
			},
			verify: func(t *testing.T, p *byteptr.Pointer) {
				assert.Equal(t, 1, p.Offset())
				assert.Equal(t, "bcd", p.String())
			},
		},
		{
			name: "TakeStringAddress",
			mutate: func(p *byteptr.Pointer) *byteptr.Pointer {
				p.InitString("hello", 1, 3)
				return p.TakeStringAddress("abcdefg")
			},
			verify: func(t *testing.T, p *byteptr.Pointer) {
				assert.Equal(t, 1, p.Offset())
				assert.Equal(t, "bcd", p.String())
			},
		},
		{
			name: "Init",
			mutate: func(p *byteptr.Pointer) *byteptr.Pointer {
				return p.Init([]byte("hello"), 1, 3)
			},
			verify: func(t *testing.T, p *byteptr.Pointer) {
				assert.Equal(t, 1, p.Offset())
				assert.Equal(t, "ell", p.String())
			},
		},
		{
			name: "InitString",
			mutate: func(p *byteptr.Pointer) *byteptr.Pointer {
				return p.InitString("hello", 1, 3)
			},
			verify: func(t *testing.T, p *byteptr.Pointer) {
				assert.Equal(t, 1, p.Offset())
				assert.Equal(t, "ell", p.String())
			},
		},
		{
			name: "SetOffset",
			mutate: func(p *byteptr.Pointer) *byteptr.Pointer {
				p.Init([]byte("hello"), 0, 0)
				return p.SetOffset(2)
			},
			verify: func(t *testing.T, p *byteptr.Pointer) {
				assert.Equal(t, 2, p.Offset())
			},
		},
		{
			name: "SetLen",
			mutate: func(p *byteptr.Pointer) *byteptr.Pointer {
				p.Init([]byte("hello"), 0, 0)
				return p.SetLen(3)
			},
			verify: func(t *testing.T, p *byteptr.Pointer) {
				assert.Equal(t, 3, p.Len())
			},
		},
		{
			name: "Reset",
			mutate: func(p *byteptr.Pointer) *byteptr.Pointer {
				p.Init([]byte("hello"), 1, 3)
				return p.Reset()
			},
			verify: func(t *testing.T, p *byteptr.Pointer) {
				assert.Nil(t, p.Bytes())
				assert.Equal(t, "", p.String())
			},
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			p := &byteptr.Pointer{}
			got := tc.mutate(p)
			require.Same(t, p, got)
			tc.verify(t, p)
		})
	}
}

func TestPointerSize(t *testing.T) {
	var p byteptr.Pointer
	assert.Equal(t, 4*unsafe.Sizeof(uintptr(0)), unsafe.Sizeof(p))
}

// mutableBytes returns a writable copy of s. A plain []byte("...") conversion
// may be placed in read-only memory by the compiler, which breaks aliasing
// tests that mutate the buffer through an unsafe-derived view.
func mutableBytes(s string) []byte {
	b := make([]byte, len(s))
	copy(b, s)
	return b
}
