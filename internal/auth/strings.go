package auth

import "unsafe"

func ZeroString(s string) {
	ptr := unsafe.Pointer(unsafe.StringData(s))
	sz := len(s)

	for i := 0; i < sz; i++ {
		*(*byte)(ptr) = 0
		ptr = unsafe.Pointer(uintptr(ptr) + 1)
	}
}
