package core

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

var (
	folderIDDesktop      = guid{0xB4BFCC3A, 0xDB2C, 0x424C, [8]byte{0xB0, 0x29, 0x7F, 0xE9, 0x9A, 0x87, 0xC6, 0x41}}
	folderIDLocalAppData = guid{0xF1B32785, 0x6FBA, 0x4FCF, [8]byte{0x9D, 0x55, 0x7B, 0x8E, 0x7F, 0x15, 0x70, 0x91}}

	clsidShellLink  = guid{0x00021401, 0x0000, 0x0000, [8]byte{0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x46}}
	iidIShellLinkW  = guid{0x000214F9, 0x0000, 0x0000, [8]byte{0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x46}}
	iidIPersistFile = guid{0x0000010B, 0x0000, 0x0000, [8]byte{0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x46}}
)

// 快捷方式 COM 接口的方法序号（vtable 下标）。
const (
	vtQueryInterface   = 0
	vtRelease          = 2
	vtSetDescription   = 7
	vtSetWorkingDir    = 9
	vtSetArguments     = 11
	vtSetIconLocation  = 17
	vtSetPath          = 20
	vtPersistFileSave  = 6
	clsctxInProcServer = 1
)

func knownFolder(id *guid) (string, error) {
	var p *uint16
	hr := syscallN(procSHGetKnownFolderPath, id.ptr(), 0, 0, uintptr(unsafe.Pointer(&p)))
	if failedHR(hr) || p == nil {
		return "", fmt.Errorf("无法获取系统文件夹路径 (0x%08X)", uint32(hr))
	}
	defer syscallN(procCoTaskMemFree, uintptr(unsafe.Pointer(p)))
	return readUTF16(p, 32768), nil
}

func readUTF16(p *uint16, max int) string {
	var units []uint16
	for i := 0; i < max; i++ {
		c := *(*uint16)(unsafe.Add(unsafe.Pointer(p), uintptr(i)*2))
		if c == 0 {
			break
		}
		units = append(units, c)
	}
	return string(utf16.Decode(units))
}

// DesktopDir 返回当前用户桌面目录（自动处理 OneDrive 重定向）。
func DesktopDir() (string, error) {
	dir, err := knownFolder(&folderIDDesktop)
	if err != nil || dir == "" {
		if home, e := os.UserHomeDir(); e == nil {
			return filepath.Join(home, "Desktop"), nil
		}
		return "", fmt.Errorf("无法定位桌面文件夹")
	}
	return dir, nil
}

// IconsDir 返回存放生成后 ico 文件的目录，并确保其存在。
func IconsDir() (string, error) {
	base, err := knownFolder(&folderIDLocalAppData)
	if err != nil || base == "" {
		base, err = os.UserCacheDir()
		if err != nil {
			return "", fmt.Errorf("无法定位应用数据目录")
		}
	}
	dir := filepath.Join(base, "QuickLinkMaker", "icons")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return dir, nil
}

func comCall(obj unsafe.Pointer, index int, args ...uintptr) uintptr {
	vtbl := *(*unsafe.Pointer)(obj)
	fn := *(*uintptr)(unsafe.Add(vtbl, uintptr(index)*unsafe.Sizeof(uintptr(0))))
	all := make([]uintptr, 0, len(args)+1)
	all = append(all, uintptr(obj))
	all = append(all, args...)
	r, _, _ := syscall.SyscallN(fn, all...)
	return r
}

func comRelease(obj unsafe.Pointer) {
	if obj != nil {
		comCall(obj, vtRelease)
	}
}

// CreateShortcut 生成一个带自定义图标的快捷方式。
// target 可以是 exe 路径，也可以是网址（此时点击后由默认浏览器打开）。
func CreateShortcut(lnkPath, target, arguments, workDir, description, iconPath string) error {
	var link unsafe.Pointer
	hr := syscallN(procCoCreateInstance,
		clsidShellLink.ptr(), 0, clsctxInProcServer,
		iidIShellLinkW.ptr(),
		uintptr(unsafe.Pointer(&link)))
	if failedHR(hr) || link == nil {
		return fmt.Errorf("无法创建快捷方式对象 (0x%08X)", uint32(hr))
	}
	defer comRelease(link)

	tPtr := utf16Ptr(target)
	if hr := comCall(link, vtSetPath, uintptr(unsafe.Pointer(tPtr))); failedHR(hr) {
		return fmt.Errorf("设置快捷方式目标失败 (0x%08X)", uint32(hr))
	}
	var aPtr, wPtr, dPtr, iPtr *uint16
	if arguments != "" {
		aPtr = utf16Ptr(arguments)
		comCall(link, vtSetArguments, uintptr(unsafe.Pointer(aPtr)))
	}
	if workDir != "" {
		wPtr = utf16Ptr(workDir)
		comCall(link, vtSetWorkingDir, uintptr(unsafe.Pointer(wPtr)))
	}
	if description != "" {
		dPtr = utf16Ptr(description)
		comCall(link, vtSetDescription, uintptr(unsafe.Pointer(dPtr)))
	}
	if iconPath != "" {
		iPtr = utf16Ptr(iconPath)
		comCall(link, vtSetIconLocation, uintptr(unsafe.Pointer(iPtr)), 0)
	}

	var pf unsafe.Pointer
	hr = comCall(link, vtQueryInterface, iidIPersistFile.ptr(), uintptr(unsafe.Pointer(&pf)))
	if failedHR(hr) || pf == nil {
		return fmt.Errorf("无法保存快捷方式 (0x%08X)", uint32(hr))
	}
	defer comRelease(pf)

	lPtr := utf16Ptr(lnkPath)
	hr = comCall(pf, vtPersistFileSave, uintptr(unsafe.Pointer(lPtr)), 1 /* TRUE */)
	runtime.KeepAlive(tPtr)
	runtime.KeepAlive(aPtr)
	runtime.KeepAlive(wPtr)
	runtime.KeepAlive(dPtr)
	runtime.KeepAlive(iPtr)
	runtime.KeepAlive(lPtr)
	if failedHR(hr) {
		return fmt.Errorf("写入快捷方式文件失败 (0x%08X)", uint32(hr))
	}
	NotifyShell(lnkPath)
	return nil
}

// NotifyShell 通知资源管理器刷新，使桌面立即显示新图标。
func NotifyShell(path string) {
	p := utf16Ptr(path)
	// SHCNE_CREATE | SHCNF_PATHW
	syscallN(procSHChangeNotify, 2, 0x0005, uintptr(unsafe.Pointer(p)), 0)
	runtime.KeepAlive(p)
}
