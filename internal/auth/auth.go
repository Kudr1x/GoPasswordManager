package auth

import (
	"unsafe"
)

var masterPassword []byte

func StorePassword(password string) {
	b := []byte(password)

	masterPassword = make([]byte, len(b))
	copy(masterPassword, b)

	zeroBytes(b)
}

func GetPassword() []byte {
	defer zeroBytes(masterPassword)

	tmp := make([]byte, len(masterPassword))
	copy(tmp, masterPassword)
	return tmp
}

func zeroBytes(b []byte) {
	for i := range b {
		b[i] = 0
	}

	ptr := unsafe.Pointer(&b[0])
	sz := len(b)
	for i := 0; i < sz; i++ {
		*(*byte)(ptr) = 0
		ptr = unsafe.Pointer(uintptr(ptr) + 1)
	}
}
