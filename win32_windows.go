package main

import (
	"syscall"
	"unicode/utf16"
	"unsafe"
)

var (
	modUser32   = syscall.NewLazyDLL("user32.dll")
	modGdi32    = syscall.NewLazyDLL("gdi32.dll")
	modKernel32 = syscall.NewLazyDLL("kernel32.dll")
	modComdlg32 = syscall.NewLazyDLL("comdlg32.dll")
	modComctl32 = syscall.NewLazyDLL("comctl32.dll")

	procRegisterClassExW     = modUser32.NewProc("RegisterClassExW")
	procCreateWindowExW      = modUser32.NewProc("CreateWindowExW")
	procDefWindowProcW       = modUser32.NewProc("DefWindowProcW")
	procShowWindow           = modUser32.NewProc("ShowWindow")
	procUpdateWindow         = modUser32.NewProc("UpdateWindow")
	procGetMessageW          = modUser32.NewProc("GetMessageW")
	procTranslateMessage     = modUser32.NewProc("TranslateMessage")
	procDispatchMessageW     = modUser32.NewProc("DispatchMessageW")
	procPostQuitMessage      = modUser32.NewProc("PostQuitMessage")
	procSendMessageW         = modUser32.NewProc("SendMessageW")
	procSetWindowTextW       = modUser32.NewProc("SetWindowTextW")
	procGetWindowTextW       = modUser32.NewProc("GetWindowTextW")
	procGetWindowTextLengthW = modUser32.NewProc("GetWindowTextLengthW")
	procMessageBoxW          = modUser32.NewProc("MessageBoxW")
	procLoadImageW           = modUser32.NewProc("LoadImageW")
	procLoadIconW            = modUser32.NewProc("LoadIconW")
	procLoadCursorW          = modUser32.NewProc("LoadCursorW")
	procDestroyIcon          = modUser32.NewProc("DestroyIcon")
	procGetSystemMetrics     = modUser32.NewProc("GetSystemMetrics")
	procAdjustWindowRectEx   = modUser32.NewProc("AdjustWindowRectEx")
	procGetDC                = modUser32.NewProc("GetDC")
	procReleaseDC            = modUser32.NewProc("ReleaseDC")
	procSetFocus             = modUser32.NewProc("SetFocus")
	procInvalidateRect       = modUser32.NewProc("InvalidateRect")
	procIsDialogMessageW     = modUser32.NewProc("IsDialogMessageW")

	procGetDeviceCaps    = modGdi32.NewProc("GetDeviceCaps")
	procCreateFontW      = modGdi32.NewProc("CreateFontW")
	procCreateSolidBrush = modGdi32.NewProc("CreateSolidBrush")
	procSetBkMode        = modGdi32.NewProc("SetBkMode")
	procSetTextColor     = modGdi32.NewProc("SetTextColor")

	procGetModuleHandleW = modKernel32.NewProc("GetModuleHandleW")

	procGetOpenFileNameW     = modComdlg32.NewProc("GetOpenFileNameW")
	procInitCommonControlsEx = modComctl32.NewProc("InitCommonControlsEx")
)

const (
	wsChild        = 0x40000000
	wsVisible      = 0x10000000
	wsTabStop      = 0x00010000
	wsClipChildren = 0x02000000
	wsVScroll      = 0x00200000
	wsBorder       = 0x00800000
	wsCaption      = 0x00C00000
	wsSysMenu      = 0x00080000
	wsMinimizeBox  = 0x00020000
	wsExClientEdge = 0x00000200

	esAutoHScroll = 0x0080
	esReadOnly    = 0x0800

	bsDefPushButton = 0x00000001

	cbsDropDownList = 0x0003

	ssLeft        = 0x0000
	ssIcon        = 0x00000003
	ssCenterImage = 0x00000200

	csHRedraw = 0x0002
	csVRedraw = 0x0001

	swShow = 5

	idcArrow       = 32512
	idiApplication = 32512

	wmCreate         = 0x0001
	wmDestroy        = 0x0002
	wmSetFont        = 0x0030
	wmCommand        = 0x0111
	wmCtlColorStatic = 0x0138
	wmGetText        = 0x000D
	wmGetDefID       = 0x0400
	dcHasDefID       = 0x534B

	emSetSel = 0x00B1

	cbAddString = 0x0143
	cbSetCurSel = 0x014E
	cbGetCurSel = 0x0147

	stmSetImage = 0x0172
	imageIcon   = 1

	lrLoadFromFile = 0x0010

	smCxScreen = 0
	smCyScreen = 1

	logPixelsY = 88

	fwNormal         = 400
	fwBold           = 700
	defaultCharset   = 1
	clearTypeQuality = 5

	mbOK              = 0x00000000
	mbIconError       = 0x00000010
	mbIconWarning     = 0x00000030
	mbIconInformation = 0x00000040

	ofnFileMustExist = 0x00001000
	ofnPathMustExist = 0x00000800
	ofnNoChangeDir   = 0x00000008
	ofnExplorer      = 0x00080000
)

type wndclassex struct {
	cbSize        uint32
	style         uint32
	lpfnWndProc   uintptr
	cbClsExtra    int32
	cbWndExtra    int32
	hInstance     uintptr
	hIcon         uintptr
	hCursor       uintptr
	hbrBackground uintptr
	lpszMenuName  *uint16
	lpszClassName *uint16
	hIconSm       uintptr
}

type msg struct {
	hwnd     uintptr
	message  uint32
	wParam   uintptr
	lParam   uintptr
	time     uint32
	ptX      int32
	ptY      int32
	lPrivate uint32
}

type rect struct {
	left   int32
	top    int32
	right  int32
	bottom int32
}

type openfilename struct {
	lStructSize       uint32
	hwndOwner         uintptr
	hInstance         uintptr
	lpstrFilter       *uint16
	lpstrCustomFilter *uint16
	nMaxCustFilter    uint32
	nFilterIndex      uint32
	lpstrFile         *uint16
	nMaxFile          uint32
	lpstrFileTitle    *uint16
	nMaxFileTitle     uint32
	lpstrInitialDir   *uint16
	lpstrTitle        *uint16
	flags             uint32
	nFileOffset       uint16
	nFileExtension    uint16
	lpstrDefExt       *uint16
	lCustData         uintptr
	lpfnHook          uintptr
	lpTemplateName    *uint16
	pvReserved        uintptr
	dwReserved        uint32
	flagsEx           uint32
}

// 编译期校验结构体大小与 Win64 ABI 一致。
var _ = [1]struct{}{}[unsafe.Sizeof(openfilename{})-152]
var _ = [1]struct{}{}[unsafe.Sizeof(msg{})-48]
var _ = [1]struct{}{}[unsafe.Sizeof(wndclassex{})-80]

type initCommonControlsEx struct {
	size uint32
	icc  uint32
}

func syscallN(p *syscall.LazyProc, args ...uintptr) uintptr {
	r, _, _ := syscall.SyscallN(p.Addr(), args...)
	return r
}

func utf16Ptr(s string) *uint16 {
	u := utf16.Encode([]rune(s))
	u = append(u, 0)
	return &u[0]
}

func rgb(r, g, b uint8) uint32 {
	return uint32(r) | uint32(g)<<8 | uint32(b)<<16
}

func (a *app) sc(v int32) int32 {
	return v * a.dpi / 96
}
