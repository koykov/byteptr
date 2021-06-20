# Byteptr

Byteptr is a solution to manipulate with bytes/strings (and their parts) without
using pointers. It's a part of the reducing pointer policy.

Byteptr is similar of [refrect.SliceHeader](https://golang.org/pkg/reflect/#SliceHeader)
and [refrecl.StringHeader](https://golang.org/pkg/reflect/#StringHeader) structs, but
provides extra methods to manipulate with data in raw memory.
