package sys

import (
	"reflect"
	"time"
	"unsafe"
)

// ToTimeDuration convert int to time.Duration seconds.
func ToTimeDuration(v int) time.Duration {
	return time.Second * time.Duration(v)
}

// ToString convert without memory allocations.
func ToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// ToBytes convert without memory allocations.
func ToBytes(s string) (b []byte) {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh.Data = sh.Data
	bh.Cap = sh.Len
	bh.Len = sh.Len
	return b
}
