package core

import (
	"syscall"
	"unicode/utf16"
	"unsafe"
)

var (
	modShell32 = syscall.NewLazyDLL("shell32.dll")
	modOle32   = syscall.NewLazyDLL("ole32.dll")

	procSHGetKnownFolderPath = modShell32.NewProc("SHGetKnownFolderPath")
	procSHChangeNotify       = modShell32.NewProc("SHChangeNotify")
	procCoInitializeEx       = modOle32.NewProc("CoInitializeEx")
	procCoCreateInstance     = modOle32.NewProc("CoCreateInstance")
	procCoTaskMemFree        = modOle32.NewProc("CoTaskMemFree")
)

type guid struct {
	Data1 uint32
	Data2 uint16
	Data3 uint16
	Data4 [8]byte
}

func (g *guid) ptr() uintptr { return uintptr(unsafe.Pointer(g)) }

func syscallN(p *syscall.LazyProc, args ...uintptr) uintptr {
	r, _, _ := syscall.SyscallN(p.Addr(), args...)
	return r
}

func failedHR(hr uintptr) bool { return int32(hr) < 0 }

func utf16Ptr(s string) *uint16 {
	u := utf16.Encode([]rune(s))
	u = append(u, 0)
	return &u[0]
}

// InitCOM 在当前线程以 STA 方式初始化 COM。
func InitCOM() {
	syscallN(procCoInitializeEx, 0, 2 /* COINIT_APARTMENTTHREADED */)
}
