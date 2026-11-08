package byteptr_test

import (
	"runtime"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/koykov/byteptr"
)

type stageKV struct {
	name                         string
	src                          []byte
	offsetK, lenK, offsetV, lenV int
	wantK, wantV                 string
	wantOffsetK, wantLenK        int
	wantOffsetV, wantLenV        int
}

func stagesKV() []stageKV {
	return []stageKV{
		{
			name: "full windows", src: []byte("hello"),
			offsetK: 0, lenK: 5, offsetV: 0, lenV: 5,
			wantK: "hello", wantV: "hello",
			wantOffsetK: 0, wantLenK: 5, wantOffsetV: 0, wantLenV: 5,
		},
		{
			name: "independent windows", src: []byte("hello"),
			offsetK: 0, lenK: 2, offsetV: 2, lenV: 3,
			wantK: "he", wantV: "llo",
			wantOffsetK: 0, wantLenK: 2, wantOffsetV: 2, wantLenV: 3,
		},
		{
			name: "overlapping windows", src: []byte("hello"),
			offsetK: 1, lenK: 3, offsetV: 0, lenV: 2,
			wantK: "ell", wantV: "he",
			wantOffsetK: 1, wantLenK: 3, wantOffsetV: 0, wantLenV: 2,
		},
		{
			name: "zero length windows", src: []byte("hello"),
			offsetK: 0, lenK: 0, offsetV: 0, lenV: 0,
			wantK: "", wantV: "",
			wantOffsetK: 0, wantLenK: 0, wantOffsetV: 0, wantLenV: 0,
		},
		{
			name: "K valid V invalid offset", src: []byte("hello"),
			offsetK: 0, lenK: 2, offsetV: 6, lenV: 1,
			wantK: "he", wantV: "",
			wantOffsetK: 0, wantLenK: 2, wantOffsetV: 0, wantLenV: 0,
		},
		{
			name: "K invalid V valid", src: []byte("hello"),
			offsetK: 6, lenK: 1, offsetV: 0, lenV: 2,
			wantK: "", wantV: "he",
			wantOffsetK: 0, wantLenK: 0, wantOffsetV: 0, wantLenV: 2,
		},
		{
			name: "K invalid length V valid", src: []byte("hello"),
			offsetK: 0, lenK: 6, offsetV: 0, lenV: 2,
			wantK: "", wantV: "he",
			wantOffsetK: 0, wantLenK: 0, wantOffsetV: 0, wantLenV: 2,
		},
		{
			name: "K valid V invalid length", src: []byte("hello"),
			offsetK: 0, lenK: 2, offsetV: 0, lenV: 6,
			wantK: "he", wantV: "",
			wantOffsetK: 0, wantLenK: 2, wantOffsetV: 0, wantLenV: 0,
		},
		{
			name: "negative K window", src: []byte("hello"),
			offsetK: -1, lenK: 2, offsetV: 0, lenV: 2,
			wantK: "", wantV: "he",
			wantOffsetK: 0, wantLenK: 0, wantOffsetV: 0, wantLenV: 2,
		},
		{
			name: "negative V window", src: []byte("hello"),
			offsetK: 0, lenK: 2, offsetV: -1, lenV: 2,
			wantK: "he", wantV: "",
			wantOffsetK: 0, wantLenK: 2, wantOffsetV: 0, wantLenV: 0,
		},
		{
			name: "tail windows", src: []byte("hello"),
			offsetK: 0, lenK: 5, offsetV: 1, lenV: 4,
			wantK: "hello", wantV: "ello",
			wantOffsetK: 0, wantLenK: 5, wantOffsetV: 1, wantLenV: 4,
		},
		{
			name: "empty windows at max", src: []byte("hello"),
			offsetK: 5, lenK: 0, offsetV: 5, lenV: 0,
			wantK: "", wantV: "",
			wantOffsetK: 5, wantLenK: 0, wantOffsetV: 5, wantLenV: 0,
		},
		{
			name: "K offset at max with non-zero length", src: []byte("hello"),
			offsetK: 5, lenK: 1, offsetV: 0, lenV: 2,
			wantK: "", wantV: "he",
			wantOffsetK: 0, wantLenK: 0, wantOffsetV: 0, wantLenV: 2,
		},
		{
			name: "V offset at max with non-zero length", src: []byte("hello"),
			offsetK: 0, lenK: 2, offsetV: 5, lenV: 1,
			wantK: "he", wantV: "",
			wantOffsetK: 0, wantLenK: 2, wantOffsetV: 0, wantLenV: 0,
		},
		{
			name: "nil slice", src: nil,
			offsetK: 0, lenK: 0, offsetV: 0, lenV: 0,
			wantK: "", wantV: "",
			wantOffsetK: 0, wantLenK: 0, wantOffsetV: 0, wantLenV: 0,
		},
		{
			name: "nil slice with windows", src: nil,
			offsetK: 0, lenK: 2, offsetV: 2, lenV: 3,
			wantK: "", wantV: "",
			wantOffsetK: 0, wantLenK: 0, wantOffsetV: 0, wantLenV: 0,
		},
		{
			name: "empty non-nil slice", src: []byte{},
			offsetK: 0, lenK: 0, offsetV: 0, lenV: 0,
			wantK: "", wantV: "",
			wantOffsetK: 0, wantLenK: 0, wantOffsetV: 0, wantLenV: 0,
		},
		{
			name: "empty non-nil slice with V length", src: []byte{},
			offsetK: 0, lenK: 0, offsetV: 0, lenV: 1,
			wantK: "", wantV: "",
			wantOffsetK: 0, wantLenK: 0, wantOffsetV: 0, wantLenV: 0,
		},
	}
}

type kvStringWindowCase struct {
	name                         string
	src                          string
	offsetK, lenK, offsetV, lenV int
	wantK, wantV                 string
	wantOffsetK, wantLenK        int
	wantOffsetV, wantLenV        int
}

func kvStringWindowCases() []kvStringWindowCase {
	return []kvStringWindowCase{
		{
			name: "full windows", src: "hello",
			offsetK: 0, lenK: 5, offsetV: 0, lenV: 5,
			wantK: "hello", wantV: "hello",
			wantOffsetK: 0, wantLenK: 5, wantOffsetV: 0, wantLenV: 5,
		},
		{
			name: "independent windows", src: "hello",
			offsetK: 0, lenK: 2, offsetV: 2, lenV: 3,
			wantK: "he", wantV: "llo",
			wantOffsetK: 0, wantLenK: 2, wantOffsetV: 2, wantLenV: 3,
		},
		{
			name: "overlapping windows", src: "hello",
			offsetK: 1, lenK: 3, offsetV: 0, lenV: 2,
			wantK: "ell", wantV: "he",
			wantOffsetK: 1, wantLenK: 3, wantOffsetV: 0, wantLenV: 2,
		},
		{
			name: "zero length windows", src: "hello",
			offsetK: 0, lenK: 0, offsetV: 0, lenV: 0,
			wantK: "", wantV: "",
			wantOffsetK: 0, wantLenK: 0, wantOffsetV: 0, wantLenV: 0,
		},
		{
			name: "K valid V invalid offset", src: "hello",
			offsetK: 0, lenK: 2, offsetV: 6, lenV: 1,
			wantK: "he", wantV: "",
			wantOffsetK: 0, wantLenK: 2, wantOffsetV: 0, wantLenV: 0,
		},
		{
			name: "K invalid V valid", src: "hello",
			offsetK: 6, lenK: 1, offsetV: 0, lenV: 2,
			wantK: "", wantV: "he",
			wantOffsetK: 0, wantLenK: 0, wantOffsetV: 0, wantLenV: 2,
		},
		{
			name: "negative K window", src: "hello",
			offsetK: -1, lenK: 2, offsetV: 0, lenV: 2,
			wantK: "", wantV: "he",
			wantOffsetK: 0, wantLenK: 0, wantOffsetV: 0, wantLenV: 2,
		},
		{
			name: "negative V window", src: "hello",
			offsetK: 0, lenK: 2, offsetV: -1, lenV: 2,
			wantK: "he", wantV: "",
			wantOffsetK: 0, wantLenK: 2, wantOffsetV: 0, wantLenV: 0,
		},
		{
			name: "tail windows", src: "hello",
			offsetK: 0, lenK: 5, offsetV: 1, lenV: 4,
			wantK: "hello", wantV: "ello",
			wantOffsetK: 0, wantLenK: 5, wantOffsetV: 1, wantLenV: 4,
		},
		{
			name: "empty windows at max", src: "hello",
			offsetK: 5, lenK: 0, offsetV: 5, lenV: 0,
			wantK: "", wantV: "",
			wantOffsetK: 5, wantLenK: 0, wantOffsetV: 5, wantLenV: 0,
		},
		{
			name: "K offset at max with non-zero length", src: "hello",
			offsetK: 5, lenK: 1, offsetV: 0, lenV: 2,
			wantK: "", wantV: "he",
			wantOffsetK: 0, wantLenK: 0, wantOffsetV: 0, wantLenV: 2,
		},
		{
			name: "V offset at max with non-zero length", src: "hello",
			offsetK: 0, lenK: 2, offsetV: 5, lenV: 1,
			wantK: "he", wantV: "",
			wantOffsetK: 0, wantLenK: 2, wantOffsetV: 0, wantLenV: 0,
		},
		{
			name: "empty string", src: "",
			offsetK: 0, lenK: 0, offsetV: 0, lenV: 0,
			wantK: "", wantV: "",
			wantOffsetK: 0, wantLenK: 0, wantOffsetV: 0, wantLenV: 0,
		},
		{
			name: "empty string with windows", src: "",
			offsetK: 0, lenK: 2, offsetV: 2, lenV: 3,
			wantK: "", wantV: "",
			wantOffsetK: 0, wantLenK: 0, wantOffsetV: 0, wantLenV: 0,
		},
		{
			name: "empty string with offset beyond max", src: "",
			offsetK: 5, lenK: 0, offsetV: 0, lenV: 0,
			wantK: "", wantV: "",
			wantOffsetK: 0, wantLenK: 0, wantOffsetV: 0, wantLenV: 0,
		},
	}
}

func TestInitKV(t *testing.T) {
	for _, tc := range stagesKV() {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			p := byteptr.InitKV(tc.src, tc.offsetK, tc.lenK, tc.offsetV, tc.lenV)

			assert.Equal(t, tc.wantOffsetK, p.OffsetK())
			assert.Equal(t, tc.wantLenK, p.LenK())
			assert.Equal(t, tc.wantOffsetV, p.OffsetV())
			assert.Equal(t, tc.wantLenV, p.LenV())

			gotK := p.BytesK()
			gotV := p.BytesV()
			if tc.src == nil {
				assert.Nil(t, gotK)
				assert.Nil(t, gotV)
				return
			}

			require.NotNil(t, gotK)
			assert.Equal(t, tc.wantK, string(gotK))
			assert.Equal(t, len(gotK), cap(gotK))

			require.NotNil(t, gotV)
			assert.Equal(t, tc.wantV, string(gotV))
			assert.Equal(t, len(gotV), cap(gotV))

			assert.Equal(t, tc.wantK, p.StringK())
			assert.Equal(t, tc.wantV, p.StringV())
		})
	}

	t.Run("reused pointer resets on nil source", func(t *testing.T) {
		p := byteptr.InitKV(mutableBytes("hello"), 0, 2, 0, 2)
		require.Equal(t, "he", p.StringK())
		require.Equal(t, "he", p.StringV())

		p.InitKV(nil, 0, 5, 0, 5)

		assert.Equal(t, "", p.StringK())
		assert.Equal(t, "", p.StringV())
		assert.Nil(t, p.BytesK())
		assert.Nil(t, p.BytesV())
		assert.Equal(t, 0, p.OffsetK())
		assert.Equal(t, 0, p.LenK())
		assert.Equal(t, 0, p.OffsetV())
		assert.Equal(t, 0, p.LenV())
	})
}

func TestInitStringKV(t *testing.T) {
	for _, tc := range kvStringWindowCases() {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			p := byteptr.InitStringKV(tc.src, tc.offsetK, tc.lenK, tc.offsetV, tc.lenV)

			assert.Equal(t, tc.wantOffsetK, p.OffsetK())
			assert.Equal(t, tc.wantLenK, p.LenK())
			assert.Equal(t, tc.wantOffsetV, p.OffsetV())
			assert.Equal(t, tc.wantLenV, p.LenV())

			gotK := p.BytesK()
			gotV := p.BytesV()
			if tc.src == "" {
				assert.Nil(t, gotK)
				assert.Nil(t, gotV)
				return
			}

			require.NotNil(t, gotK)
			assert.Equal(t, tc.wantK, string(gotK))
			assert.Equal(t, len(gotK), cap(gotK))

			require.NotNil(t, gotV)
			assert.Equal(t, tc.wantV, string(gotV))
			assert.Equal(t, len(gotV), cap(gotV))

			assert.Equal(t, tc.wantK, p.StringK())
			assert.Equal(t, tc.wantV, p.StringV())
		})
	}

	t.Run("reused pointer resets on empty source", func(t *testing.T) {
		p := byteptr.InitStringKV("hello", 0, 2, 0, 2)
		require.Equal(t, "he", p.StringK())
		require.Equal(t, "he", p.StringV())

		p.InitStringKV("", 0, 5, 0, 5)

		assert.Equal(t, "", p.StringK())
		assert.Equal(t, "", p.StringV())
		assert.Nil(t, p.BytesK())
		assert.Nil(t, p.BytesV())
		assert.Equal(t, 0, p.OffsetK())
		assert.Equal(t, 0, p.LenK())
		assert.Equal(t, 0, p.OffsetV())
		assert.Equal(t, 0, p.LenV())
	})
}

func TestKVTakeAddress(t *testing.T) {
	tests := []struct {
		name  string
		check func(t *testing.T)
	}{
		{
			name: "nil is a no-op and preserves both windows",
			check: func(t *testing.T) {
				p := byteptr.InitKV([]byte("hello"), 0, 2, 2, 3)
				p.TakeAddress(nil)
				assert.Equal(t, 0, p.OffsetK())
				assert.Equal(t, 2, p.LenK())
				assert.Equal(t, 2, p.OffsetV())
				assert.Equal(t, 3, p.LenV())
				assert.Equal(t, "he", p.StringK())
				assert.Equal(t, "llo", p.StringV())
			},
		},
		{
			name: "non-nil sets address and capacity, preserving windows",
			check: func(t *testing.T) {
				p := byteptr.InitStringKV("hello", 0, 2, 2, 3)
				p.TakeAddress([]byte("abcdefgh"))
				assert.Equal(t, 0, p.OffsetK())
				assert.Equal(t, 2, p.LenK())
				assert.Equal(t, 2, p.OffsetV())
				assert.Equal(t, 3, p.LenV())
				assert.Equal(t, "ab", p.StringK())
				assert.Equal(t, "cde", p.StringV())
			},
		},
		{
			name: "empty non-nil sets address with zero capacity",
			check: func(t *testing.T) {
				p := &byteptr.KVPointer{}
				p.TakeAddress([]byte{})
				gotK := p.BytesK()
				gotV := p.BytesV()
				require.NotNil(t, gotK)
				require.NotNil(t, gotV)
				assert.Equal(t, 0, len(gotK))
				assert.Equal(t, 0, cap(gotK))
				assert.Equal(t, 0, len(gotV))
				assert.Equal(t, 0, cap(gotV))
				assert.Equal(t, "", p.StringK())
				assert.Equal(t, "", p.StringV())
			},
		},
		{
			name: "preserves offsets and lengths when re-taking a new buffer",
			check: func(t *testing.T) {
				p := byteptr.InitKV([]byte("hello"), 1, 3, 2, 2)
				require.Equal(t, "ell", p.StringK())
				require.Equal(t, "ll", p.StringV())
				p.TakeAddress([]byte("abcdefghij"))
				assert.Equal(t, 1, p.OffsetK())
				assert.Equal(t, 3, p.LenK())
				assert.Equal(t, 2, p.OffsetV())
				assert.Equal(t, 2, p.LenV())
				assert.Equal(t, "bcd", p.StringK())
				assert.Equal(t, "cd", p.StringV())
			},
		},
		{
			name: "smaller buffer makes both getters return nil",
			check: func(t *testing.T) {
				p := byteptr.InitKV([]byte("hello"), 3, 2, 3, 2)
				require.Equal(t, "lo", p.StringK())
				require.Equal(t, "lo", p.StringV())
				p.TakeAddress([]byte("ab"))
				assert.Equal(t, 3, p.OffsetK())
				assert.Equal(t, 2, p.LenK())
				assert.Equal(t, 3, p.OffsetV())
				assert.Equal(t, 2, p.LenV())
				assert.Nil(t, p.BytesK())
				assert.Nil(t, p.BytesV())
				assert.Equal(t, "", p.StringK())
				assert.Equal(t, "", p.StringV())
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

func TestKVTakeStringAddress(t *testing.T) {
	tests := []struct {
		name  string
		check func(t *testing.T)
	}{
		{
			name: "empty string on a fresh pointer is a no-op",
			check: func(t *testing.T) {
				p := &byteptr.KVPointer{}
				p.TakeStringAddress("")
				assert.Equal(t, 0, p.OffsetK())
				assert.Equal(t, 0, p.LenK())
				assert.Equal(t, 0, p.OffsetV())
				assert.Equal(t, 0, p.LenV())
				assert.Nil(t, p.BytesK())
				assert.Nil(t, p.BytesV())
				assert.Equal(t, "", p.StringK())
				assert.Equal(t, "", p.StringV())
			},
		},
		{
			name: "empty string is a no-op and preserves both windows",
			check: func(t *testing.T) {
				p := byteptr.InitStringKV("hello", 0, 2, 2, 3)
				p.TakeStringAddress("")
				assert.Equal(t, 0, p.OffsetK())
				assert.Equal(t, 2, p.LenK())
				assert.Equal(t, 2, p.OffsetV())
				assert.Equal(t, 3, p.LenV())
				assert.Equal(t, "he", p.StringK())
				assert.Equal(t, "llo", p.StringV())
			},
		},
		{
			name: "non-empty sets address and capacity",
			check: func(t *testing.T) {
				p := &byteptr.KVPointer{}
				p.TakeStringAddress("abcdef")
				p.SetOffsetK(1).SetLenK(2)
				p.SetOffsetV(4).SetLenV(2)
				assert.Equal(t, "bc", p.StringK())
				assert.Equal(t, "ef", p.StringV())
			},
		},
		{
			name: "preserves offsets and lengths when re-taking a new string",
			check: func(t *testing.T) {
				p := byteptr.InitStringKV("hello", 0, 2, 2, 3)
				p.TakeStringAddress("abcdefgh")
				assert.Equal(t, 0, p.OffsetK())
				assert.Equal(t, 2, p.LenK())
				assert.Equal(t, 2, p.OffsetV())
				assert.Equal(t, 3, p.LenV())
				assert.Equal(t, "ab", p.StringK())
				assert.Equal(t, "cde", p.StringV())
			},
		},
		{
			name: "smaller string makes both getters return nil",
			check: func(t *testing.T) {
				p := byteptr.InitStringKV("hello", 3, 2, 3, 2)
				require.Equal(t, "lo", p.StringK())
				require.Equal(t, "lo", p.StringV())
				p.TakeStringAddress("ab")
				assert.Equal(t, 3, p.OffsetK())
				assert.Equal(t, 2, p.LenK())
				assert.Equal(t, 3, p.OffsetV())
				assert.Equal(t, 2, p.LenV())
				assert.Nil(t, p.BytesK())
				assert.Nil(t, p.BytesV())
				assert.Equal(t, "", p.StringK())
				assert.Equal(t, "", p.StringV())
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

func TestKVSetters(t *testing.T) {
	tests := []struct {
		name    string
		prepare func(p *byteptr.KVPointer)
		apply   func(p *byteptr.KVPointer)
		get     func(p *byteptr.KVPointer) int
		want    int
	}{
		{
			name:    "SetOffsetK zero accepted",
			prepare: func(p *byteptr.KVPointer) { p.InitKV([]byte("hello"), 2, 2, 0, 0) },
			apply:   func(p *byteptr.KVPointer) { p.SetOffsetK(0) },
			get:     func(p *byteptr.KVPointer) int { return p.OffsetK() },
			want:    0,
		},
		{
			name:    "SetOffsetK in range accepted",
			prepare: func(p *byteptr.KVPointer) { p.InitKV([]byte("hello"), 2, 2, 0, 0) },
			apply:   func(p *byteptr.KVPointer) { p.SetOffsetK(4) },
			get:     func(p *byteptr.KVPointer) int { return p.OffsetK() },
			want:    4,
		},
		{
			name:    "SetOffsetK max accepted",
			prepare: func(p *byteptr.KVPointer) { p.InitKV([]byte("hello"), 0, 0, 0, 0) },
			apply:   func(p *byteptr.KVPointer) { p.SetOffsetK(5) },
			get:     func(p *byteptr.KVPointer) int { return p.OffsetK() },
			want:    5,
		},
		{
			name:    "SetOffsetK max plus one ignored",
			prepare: func(p *byteptr.KVPointer) { p.InitKV([]byte("hello"), 2, 2, 0, 0) },
			apply:   func(p *byteptr.KVPointer) { p.SetOffsetK(6) },
			get:     func(p *byteptr.KVPointer) int { return p.OffsetK() },
			want:    2,
		},
		{
			name:    "SetOffsetK negative ignored",
			prepare: func(p *byteptr.KVPointer) { p.InitKV([]byte("hello"), 2, 2, 0, 0) },
			apply:   func(p *byteptr.KVPointer) { p.SetOffsetK(-1) },
			get:     func(p *byteptr.KVPointer) int { return p.OffsetK() },
			want:    2,
		},
		{
			name:    "SetLenK zero accepted",
			prepare: func(p *byteptr.KVPointer) { p.InitKV([]byte("hello"), 2, 3, 0, 0) },
			apply:   func(p *byteptr.KVPointer) { p.SetLenK(0) },
			get:     func(p *byteptr.KVPointer) int { return p.LenK() },
			want:    0,
		},
		{
			name:    "SetLenK max minus offset accepted",
			prepare: func(p *byteptr.KVPointer) { p.InitKV([]byte("hello"), 2, 0, 0, 0) },
			apply:   func(p *byteptr.KVPointer) { p.SetLenK(3) },
			get:     func(p *byteptr.KVPointer) int { return p.LenK() },
			want:    3,
		},
		{
			name:    "SetLenK max minus offset plus one ignored",
			prepare: func(p *byteptr.KVPointer) { p.InitKV([]byte("hello"), 2, 1, 0, 0) },
			apply:   func(p *byteptr.KVPointer) { p.SetLenK(4) },
			get:     func(p *byteptr.KVPointer) int { return p.LenK() },
			want:    1,
		},
		{
			name:    "SetLenK negative ignored",
			prepare: func(p *byteptr.KVPointer) { p.InitKV([]byte("hello"), 2, 1, 0, 0) },
			apply:   func(p *byteptr.KVPointer) { p.SetLenK(-1) },
			get:     func(p *byteptr.KVPointer) int { return p.LenK() },
			want:    1,
		},
		{
			name:    "SetLenK zero at max offset accepted",
			prepare: func(p *byteptr.KVPointer) { p.InitKV([]byte("hello"), 5, 0, 0, 0) },
			apply:   func(p *byteptr.KVPointer) { p.SetLenK(0) },
			get:     func(p *byteptr.KVPointer) int { return p.LenK() },
			want:    0,
		},
		{
			name:    "SetLenK one at max offset ignored",
			prepare: func(p *byteptr.KVPointer) { p.InitKV([]byte("hello"), 5, 0, 0, 0) },
			apply:   func(p *byteptr.KVPointer) { p.SetLenK(1) },
			get:     func(p *byteptr.KVPointer) int { return p.LenK() },
			want:    0,
		},
		{
			name:    "SetLenK respects newly set offset",
			prepare: func(p *byteptr.KVPointer) { p.InitKV([]byte("hello"), 0, 0, 0, 0) },
			apply:   func(p *byteptr.KVPointer) { p.SetOffsetK(3).SetLenK(2) },
			get:     func(p *byteptr.KVPointer) int { return p.LenK() },
			want:    2,
		},
		{
			name:    "SetLenK rejects length beyond newly set offset",
			prepare: func(p *byteptr.KVPointer) { p.InitKV([]byte("hello"), 0, 0, 0, 0) },
			apply:   func(p *byteptr.KVPointer) { p.SetOffsetK(3).SetLenK(3) },
			get:     func(p *byteptr.KVPointer) int { return p.LenK() },
			want:    0,
		},
		{
			name:    "SetOffsetV zero accepted",
			prepare: func(p *byteptr.KVPointer) { p.InitKV([]byte("hello"), 0, 0, 2, 2) },
			apply:   func(p *byteptr.KVPointer) { p.SetOffsetV(0) },
			get:     func(p *byteptr.KVPointer) int { return p.OffsetV() },
			want:    0,
		},
		{
			name:    "SetOffsetV in range accepted",
			prepare: func(p *byteptr.KVPointer) { p.InitKV([]byte("hello"), 0, 0, 2, 2) },
			apply:   func(p *byteptr.KVPointer) { p.SetOffsetV(4) },
			get:     func(p *byteptr.KVPointer) int { return p.OffsetV() },
			want:    4,
		},
		{
			name:    "SetOffsetV max accepted",
			prepare: func(p *byteptr.KVPointer) { p.InitKV([]byte("hello"), 0, 0, 0, 0) },
			apply:   func(p *byteptr.KVPointer) { p.SetOffsetV(5) },
			get:     func(p *byteptr.KVPointer) int { return p.OffsetV() },
			want:    5,
		},
		{
			name:    "SetOffsetV max plus one ignored",
			prepare: func(p *byteptr.KVPointer) { p.InitKV([]byte("hello"), 0, 0, 2, 2) },
			apply:   func(p *byteptr.KVPointer) { p.SetOffsetV(6) },
			get:     func(p *byteptr.KVPointer) int { return p.OffsetV() },
			want:    2,
		},
		{
			name:    "SetOffsetV negative ignored",
			prepare: func(p *byteptr.KVPointer) { p.InitKV([]byte("hello"), 0, 0, 2, 2) },
			apply:   func(p *byteptr.KVPointer) { p.SetOffsetV(-1) },
			get:     func(p *byteptr.KVPointer) int { return p.OffsetV() },
			want:    2,
		},
		{
			name:    "SetLenV zero accepted",
			prepare: func(p *byteptr.KVPointer) { p.InitKV([]byte("hello"), 0, 0, 2, 3) },
			apply:   func(p *byteptr.KVPointer) { p.SetLenV(0) },
			get:     func(p *byteptr.KVPointer) int { return p.LenV() },
			want:    0,
		},
		{
			name:    "SetLenV max minus offset accepted",
			prepare: func(p *byteptr.KVPointer) { p.InitKV([]byte("hello"), 0, 0, 2, 0) },
			apply:   func(p *byteptr.KVPointer) { p.SetLenV(3) },
			get:     func(p *byteptr.KVPointer) int { return p.LenV() },
			want:    3,
		},
		{
			name:    "SetLenV max minus offset plus one ignored",
			prepare: func(p *byteptr.KVPointer) { p.InitKV([]byte("hello"), 0, 0, 2, 1) },
			apply:   func(p *byteptr.KVPointer) { p.SetLenV(4) },
			get:     func(p *byteptr.KVPointer) int { return p.LenV() },
			want:    1,
		},
		{
			name:    "SetLenV negative ignored",
			prepare: func(p *byteptr.KVPointer) { p.InitKV([]byte("hello"), 0, 0, 2, 1) },
			apply:   func(p *byteptr.KVPointer) { p.SetLenV(-1) },
			get:     func(p *byteptr.KVPointer) int { return p.LenV() },
			want:    1,
		},
		{
			name:    "SetLenV zero at max offset accepted",
			prepare: func(p *byteptr.KVPointer) { p.InitKV([]byte("hello"), 0, 0, 5, 0) },
			apply:   func(p *byteptr.KVPointer) { p.SetLenV(0) },
			get:     func(p *byteptr.KVPointer) int { return p.LenV() },
			want:    0,
		},
		{
			name:    "SetLenV one at max offset ignored",
			prepare: func(p *byteptr.KVPointer) { p.InitKV([]byte("hello"), 0, 0, 5, 0) },
			apply:   func(p *byteptr.KVPointer) { p.SetLenV(1) },
			get:     func(p *byteptr.KVPointer) int { return p.LenV() },
			want:    0,
		},
		{
			name:    "SetLenV respects newly set offset",
			prepare: func(p *byteptr.KVPointer) { p.InitKV([]byte("hello"), 0, 0, 0, 0) },
			apply:   func(p *byteptr.KVPointer) { p.SetOffsetV(3).SetLenV(2) },
			get:     func(p *byteptr.KVPointer) int { return p.LenV() },
			want:    2,
		},
		{
			name:    "SetLenV rejects length beyond newly set offset",
			prepare: func(p *byteptr.KVPointer) { p.InitKV([]byte("hello"), 0, 0, 0, 0) },
			apply:   func(p *byteptr.KVPointer) { p.SetOffsetV(3).SetLenV(3) },
			get:     func(p *byteptr.KVPointer) int { return p.LenV() },
			want:    0,
		},
		{
			name:    "SetOffsetK does not change V offset",
			prepare: func(p *byteptr.KVPointer) { p.InitKV([]byte("hello"), 0, 2, 2, 3) },
			apply:   func(p *byteptr.KVPointer) { p.SetOffsetK(1) },
			get:     func(p *byteptr.KVPointer) int { return p.OffsetV() },
			want:    2,
		},
		{
			name:    "SetOffsetV does not change K offset",
			prepare: func(p *byteptr.KVPointer) { p.InitKV([]byte("hello"), 0, 2, 2, 3) },
			apply:   func(p *byteptr.KVPointer) { p.SetOffsetV(1) },
			get:     func(p *byteptr.KVPointer) int { return p.OffsetK() },
			want:    0,
		},
		{
			name:    "SetLenK does not change V length",
			prepare: func(p *byteptr.KVPointer) { p.InitKV([]byte("hello"), 0, 2, 2, 3) },
			apply:   func(p *byteptr.KVPointer) { p.SetLenK(0) },
			get:     func(p *byteptr.KVPointer) int { return p.LenV() },
			want:    3,
		},
		{
			name:    "SetLenV does not change K length",
			prepare: func(p *byteptr.KVPointer) { p.InitKV([]byte("hello"), 0, 2, 2, 3) },
			apply:   func(p *byteptr.KVPointer) { p.SetLenV(0) },
			get:     func(p *byteptr.KVPointer) int { return p.LenK() },
			want:    2,
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			p := &byteptr.KVPointer{}
			tc.prepare(p)
			tc.apply(p)
			assert.Equal(t, tc.want, tc.get(p))
		})
	}
}

func TestKVBytesK(t *testing.T) {
	for _, tc := range stagesKV() {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			p := byteptr.InitKV(tc.src, tc.offsetK, tc.lenK, tc.offsetV, tc.lenV)
			got := p.BytesK()
			if tc.src == nil {
				assert.Nil(t, got)
				return
			}
			require.NotNil(t, got)
			assert.Equal(t, tc.wantK, string(got))
			assert.Equal(t, len(got), cap(got))
		})
	}

	t.Run("view mutation is visible in source and V window", func(t *testing.T) {
		buf := mutableBytes("hello")
		p := byteptr.InitKV(buf, 1, 3, 0, 5)
		view := p.BytesK()
		require.NotNil(t, view)
		require.Len(t, view, 3)
		view[0] = 'X'
		assert.Equal(t, "hXllo", string(buf))
		assert.Equal(t, "hXllo", p.StringV())
		runtime.KeepAlive(buf)
	})

	t.Run("source mutation is visible in K view", func(t *testing.T) {
		buf := mutableBytes("hello")
		p := byteptr.InitKV(buf, 1, 3, 0, 5)
		view := p.BytesK()
		require.NotNil(t, view)
		require.Len(t, view, 3)
		buf[2] = 'Y'
		assert.Equal(t, "eYl", string(view))
		runtime.KeepAlive(buf)
	})
}

func TestKVStringK(t *testing.T) {
	for _, tc := range stagesKV() {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			p := byteptr.InitKV(tc.src, tc.offsetK, tc.lenK, tc.offsetV, tc.lenV)
			assert.Equal(t, tc.wantK, p.StringK())
		})
	}

	t.Run("obtained K string observes source mutation", func(t *testing.T) {
		buf := mutableBytes("hello")
		p := byteptr.InitKV(buf, 1, 3, 0, 5)
		s := p.StringK()
		buf[1] = 'Z'
		assert.Equal(t, "Zll", s)
		runtime.KeepAlive(buf)
	})

	t.Run("K string shares memory with bytes", func(t *testing.T) {
		buf := mutableBytes("hello")
		p := byteptr.InitKV(buf, 1, 3, 0, 5)
		s := p.StringK()

		view := p.BytesK()
		require.NotNil(t, view)
		require.Len(t, view, 3)
		view[0] = 'R'

		assert.Equal(t, "Rll", s)
		runtime.KeepAlive(buf)
	})
}

func TestKVBytesV(t *testing.T) {
	for _, tc := range stagesKV() {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			p := byteptr.InitKV(tc.src, tc.offsetK, tc.lenK, tc.offsetV, tc.lenV)
			got := p.BytesV()
			if tc.src == nil {
				assert.Nil(t, got)
				return
			}
			require.NotNil(t, got)
			assert.Equal(t, tc.wantV, string(got))
			assert.Equal(t, len(got), cap(got))
		})
	}

	t.Run("view mutation is visible in source and K window", func(t *testing.T) {
		buf := mutableBytes("hello")
		p := byteptr.InitKV(buf, 0, 2, 1, 3)
		view := p.BytesV()
		require.NotNil(t, view)
		require.Len(t, view, 3)
		view[0] = 'X'
		assert.Equal(t, "hXllo", string(buf))
		assert.Equal(t, "hX", p.StringK())
		runtime.KeepAlive(buf)
	})

	t.Run("source mutation is visible in V view", func(t *testing.T) {
		buf := mutableBytes("hello")
		p := byteptr.InitKV(buf, 0, 2, 1, 3)
		view := p.BytesV()
		require.NotNil(t, view)
		require.Len(t, view, 3)
		buf[2] = 'Y'
		assert.Equal(t, "eYl", string(view))
		runtime.KeepAlive(buf)
	})
}

func TestKVStringV(t *testing.T) {
	for _, tc := range stagesKV() {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			p := byteptr.InitKV(tc.src, tc.offsetK, tc.lenK, tc.offsetV, tc.lenV)
			assert.Equal(t, tc.wantV, p.StringV())
		})
	}

	t.Run("obtained V string observes source mutation", func(t *testing.T) {
		buf := mutableBytes("hello")
		p := byteptr.InitKV(buf, 0, 5, 1, 3)
		s := p.StringV()
		buf[1] = 'Z'
		assert.Equal(t, "Zll", s)
		runtime.KeepAlive(buf)
	})

	t.Run("V string shares memory with bytes", func(t *testing.T) {
		buf := mutableBytes("hello")
		p := byteptr.InitKV(buf, 0, 5, 1, 3)
		s := p.StringV()

		view := p.BytesV()
		require.NotNil(t, view)
		require.Len(t, view, 3)
		view[0] = 'R'

		assert.Equal(t, "Rll", s)
		runtime.KeepAlive(buf)
	})
}

func TestKVReset(t *testing.T) {
	t.Run("clears all six fields", func(t *testing.T) {
		p := byteptr.InitKV([]byte("hello"), 0, 2, 2, 3)
		require.Equal(t, "he", p.StringK())
		require.Equal(t, "llo", p.StringV())

		p.Reset()

		assert.Equal(t, 0, p.OffsetK())
		assert.Equal(t, 0, p.LenK())
		assert.Equal(t, 0, p.OffsetV())
		assert.Equal(t, 0, p.LenV())
		assert.Nil(t, p.BytesK())
		assert.Nil(t, p.BytesV())
		assert.Equal(t, "", p.StringK())
		assert.Equal(t, "", p.StringV())
	})

	t.Run("reset on a zero value is safe", func(t *testing.T) {
		p := &byteptr.KVPointer{}
		p.Reset()
		assert.Nil(t, p.BytesK())
		assert.Nil(t, p.BytesV())
		assert.Equal(t, "", p.StringK())
		assert.Equal(t, "", p.StringV())
	})

	t.Run("no stale window after reset and re-take", func(t *testing.T) {
		p := byteptr.InitKV([]byte("hello"), 0, 1, 2, 3)
		require.Equal(t, "h", p.StringK())
		require.Equal(t, "llo", p.StringV())

		p.Reset()
		p.TakeAddress([]byte("abcdef"))

		assert.Equal(t, 0, p.OffsetK())
		assert.Equal(t, 0, p.LenK())
		assert.Equal(t, 0, p.OffsetV())
		assert.Equal(t, 0, p.LenV())
		assert.Equal(t, "", p.StringK())
		assert.Equal(t, "", p.StringV())

		p.SetLenK(2)
		assert.Equal(t, "ab", p.StringK())
		assert.Equal(t, "", p.StringV())
	})
}

func TestKVChaining(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(p *byteptr.KVPointer) *byteptr.KVPointer
		verify func(t *testing.T, p *byteptr.KVPointer)
	}{
		{
			name: "TakeAddress",
			mutate: func(p *byteptr.KVPointer) *byteptr.KVPointer {
				p.InitKV([]byte("hello"), 0, 2, 2, 3)
				return p.TakeAddress([]byte("abcdefgh"))
			},
			verify: func(t *testing.T, p *byteptr.KVPointer) {
				assert.Equal(t, "ab", p.StringK())
				assert.Equal(t, "cde", p.StringV())
			},
		},
		{
			name: "TakeStringAddress",
			mutate: func(p *byteptr.KVPointer) *byteptr.KVPointer {
				p.InitStringKV("hello", 0, 2, 2, 3)
				return p.TakeStringAddress("abcdefgh")
			},
			verify: func(t *testing.T, p *byteptr.KVPointer) {
				assert.Equal(t, "ab", p.StringK())
				assert.Equal(t, "cde", p.StringV())
			},
		},
		{
			name: "InitKV",
			mutate: func(p *byteptr.KVPointer) *byteptr.KVPointer {
				return p.InitKV([]byte("hello"), 0, 2, 2, 3)
			},
			verify: func(t *testing.T, p *byteptr.KVPointer) {
				assert.Equal(t, "he", p.StringK())
				assert.Equal(t, "llo", p.StringV())
			},
		},
		{
			name: "InitStringKV",
			mutate: func(p *byteptr.KVPointer) *byteptr.KVPointer {
				return p.InitStringKV("hello", 1, 2, 2, 3)
			},
			verify: func(t *testing.T, p *byteptr.KVPointer) {
				assert.Equal(t, "el", p.StringK())
				assert.Equal(t, "llo", p.StringV())
			},
		},
		{
			name: "SetOffsetK",
			mutate: func(p *byteptr.KVPointer) *byteptr.KVPointer {
				p.InitKV([]byte("hello"), 0, 0, 0, 0)
				return p.SetOffsetK(2)
			},
			verify: func(t *testing.T, p *byteptr.KVPointer) {
				assert.Equal(t, 2, p.OffsetK())
			},
		},
		{
			name: "SetLenK",
			mutate: func(p *byteptr.KVPointer) *byteptr.KVPointer {
				p.InitKV([]byte("hello"), 0, 0, 0, 0)
				return p.SetLenK(3)
			},
			verify: func(t *testing.T, p *byteptr.KVPointer) {
				assert.Equal(t, 3, p.LenK())
			},
		},
		{
			name: "SetOffsetV",
			mutate: func(p *byteptr.KVPointer) *byteptr.KVPointer {
				p.InitKV([]byte("hello"), 0, 0, 0, 0)
				return p.SetOffsetV(2)
			},
			verify: func(t *testing.T, p *byteptr.KVPointer) {
				assert.Equal(t, 2, p.OffsetV())
			},
		},
		{
			name: "SetLenV",
			mutate: func(p *byteptr.KVPointer) *byteptr.KVPointer {
				p.InitKV([]byte("hello"), 0, 0, 0, 0)
				return p.SetLenV(3)
			},
			verify: func(t *testing.T, p *byteptr.KVPointer) {
				assert.Equal(t, 3, p.LenV())
			},
		},
		{
			name: "Reset",
			mutate: func(p *byteptr.KVPointer) *byteptr.KVPointer {
				p.InitKV([]byte("hello"), 0, 2, 2, 3)
				return p.Reset()
			},
			verify: func(t *testing.T, p *byteptr.KVPointer) {
				assert.Nil(t, p.BytesK())
				assert.Nil(t, p.BytesV())
				assert.Equal(t, "", p.StringK())
				assert.Equal(t, "", p.StringV())
			},
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			p := &byteptr.KVPointer{}
			got := tc.mutate(p)
			require.Same(t, p, got)
			tc.verify(t, p)
		})
	}
}

func TestKVPointerSize(t *testing.T) {
	var p byteptr.KVPointer
	assert.Equal(t, 6*unsafe.Sizeof(uintptr(0)), unsafe.Sizeof(p))
}
