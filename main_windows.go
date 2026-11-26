package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"unsafe"

	"github.com/FSNH-MIMITU/url-shortcut-maker/core"
)

const (
	clientW = 640
	clientH = 462
)

const (
	idURL = 1001 + iota
	idIconPath
	idPick
	idBrowser
	idName
	idPreview
	idGenerate
	idStatus
	idDesc
	idHintBrowser
	idHintName
	idSubtitle
	idLang
)

const (
	statusNeutral = iota
	statusOK
	statusError
)

type app struct {
	hInst uintptr
	hwnd  uintptr

	hURL, hIconPath, hPick, hBrowser, hName   uintptr
	hPreview, hGenerate, hStatus              uintptr
	hSubtitle, hDesc, hHintBrowser, hHintName uintptr
	hLabelURL, hLabelIcon, hLabelOpenWith     uintptr
	hLabelName, hLabelPreview, hLabelLang     uintptr
	hLang                                     uintptr

	controls []uintptr

	font  uintptr
	brush uintptr
	dpi   int32
	lang  string

	colorLabel uint32
	colorHint  uint32
	colorOK    uint32
	colorErr   uint32

	selectedImage string
	hPreviewIcon  uintptr

	browsers   []core.Browser
	statusKind int
}

var ui *app

func main() {
	runtime.LockOSThread()
	core.InitCOM()

	ui = &app{
		colorLabel: rgb(0x37, 0x41, 0x51),
		colorHint:  rgb(0x6B, 0x72, 0x80),
		colorOK:    rgb(0x05, 0x96, 0x69),
		colorErr:   rgb(0xDC, 0x26, 0x26),
	}
	ui.lang = detectSystemLanguage()
	ui.hInst = syscallN(procGetModuleHandleW, 0)
	ui.init()
	ui.run()
}

func (a *app) init() {
	icc := initCommonControlsEx{size: uint32(unsafe.Sizeof(initCommonControlsEx{})), icc: 0x4000 /* ICC_STANDARD_CLASSES */}
	syscallN(procInitCommonControlsEx, uintptr(unsafe.Pointer(&icc)))

	a.dpi = getDPI()
	a.brush = syscallN(procCreateSolidBrush, uintptr(rgb(0xF7, 0xF8, 0xFA)))
	a.font = a.createFont(9, fwNormal, fontFor(a.lang))
	a.browsers = core.DetectBrowsers()
	a.registerClass()
	a.createWindow()
}

func getDPI() int32 {
	hdc := syscallN(procGetDC, 0)
	if hdc == 0 {
		return 96
	}
	defer syscallN(procReleaseDC, 0, hdc)
	dpi := int32(syscallN(procGetDeviceCaps, hdc, logPixelsY))
	if dpi <= 0 {
		dpi = 96
	}
	return dpi
}

func (a *app) createFont(pt int32, weight int32, family string) uintptr {
	height := -(pt * a.dpi / 72)
	return syscallN(procCreateFontW,
		uintptr(height), 0, 0, 0,
		uintptr(weight), 0, 0, 0,
		defaultCharset, 0, 0, clearTypeQuality, 0,
		uintptr(unsafe.Pointer(utf16Ptr(family))))
}

func (a *app) registerClass() {
	hIcon := syscallN(procLoadIconW, a.hInst, 1)
	if hIcon == 0 {
		hIcon = syscallN(procLoadIconW, 0, idiApplication)
	}
	wc := wndclassex{
		cbSize:        uint32(unsafe.Sizeof(wndclassex{})),
		style:         csHRedraw | csVRedraw,
		lpfnWndProc:   syscall.NewCallback(wndProc),
		hInstance:     a.hInst,
		hIcon:         hIcon,
		hCursor:       syscallN(procLoadCursorW, 0, idcArrow),
		hbrBackground: a.brush,
		lpszClassName: utf16Ptr("QuickLinkMakerWnd"),
	}
	if syscallN(procRegisterClassExW, uintptr(unsafe.Pointer(&wc))) == 0 {
		a.fatal(tr(a.lang).ErrWindowClass)
	}
}

func (a *app) createWindow() {
	style := uint32(wsCaption | wsSysMenu | wsMinimizeBox | wsClipChildren)
	rc := rect{right: a.sc(clientW), bottom: a.sc(clientH)}
	syscallN(procAdjustWindowRectEx, uintptr(unsafe.Pointer(&rc)), uintptr(style), 0, 0)
	w := rc.right - rc.left
	h := rc.bottom - rc.top
	x := (int32(syscallN(procGetSystemMetrics, smCxScreen)) - w) / 2
	y := (int32(syscallN(procGetSystemMetrics, smCyScreen)) - h) / 3
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}
	a.hwnd = syscallN(procCreateWindowExW,
		0,
		uintptr(unsafe.Pointer(utf16Ptr("QuickLinkMakerWnd"))),
		uintptr(unsafe.Pointer(utf16Ptr(tr(a.lang).Title))),
		uintptr(style),
		uintptr(x), uintptr(y), uintptr(w), uintptr(h),
		0, 0, a.hInst, 0)
	if a.hwnd == 0 {
		a.fatal(tr(a.lang).ErrWindowCreate)
	}
}

func (a *app) run() {
	syscallN(procShowWindow, a.hwnd, swShow)
	syscallN(procUpdateWindow, a.hwnd)
	var m msg
	for {
		r := int32(syscallN(procGetMessageW, uintptr(unsafe.Pointer(&m)), 0, 0, 0))
		if r <= 0 {
			break
		}
		// 让 Tab 切换焦点、回车触发「生成快捷方式」
		if syscallN(procIsDialogMessageW, a.hwnd, uintptr(unsafe.Pointer(&m))) != 0 {
			continue
		}
		syscallN(procTranslateMessage, uintptr(unsafe.Pointer(&m)))
		syscallN(procDispatchMessageW, uintptr(unsafe.Pointer(&m)))
	}
}

func wndProc(hwnd, message, wParam, lParam uintptr) uintptr {
	a := ui
	if a == nil {
		return syscallN(procDefWindowProcW, hwnd, message, wParam, lParam)
	}
	switch message {
	case wmCreate:
		a.hwnd = hwnd
		a.createControls()
		return 0
	case wmGetDefID:
		// 告诉消息循环默认按钮是「生成快捷方式」，回车即可触发生成
		return uintptr(dcHasDefID<<16 | idGenerate)
	case wmCommand:
		id := int(wParam & 0xFFFF)
		code := int((wParam >> 16) & 0xFFFF)
		switch id {
		case idLang:
			if code == cbnSelChange {
				sel := int(syscallN(procSendMessageW, a.hLang, cbGetCurSel, 0, 0))
				if sel >= 0 && sel < len(languages) {
					a.applyLanguage(languages[sel].code)
				}
			}
		case idPick:
			if code == 0 {
				a.pickImage()
			}
		case idGenerate:
			if code == 0 {
				a.generate()
			}
		}
		return 0
	case wmCtlColorStatic:
		hdc := wParam
		syscallN(procSetTextColor, hdc, uintptr(a.staticColor(lParam)))
		syscallN(procSetBkMode, hdc, 1 /* TRANSPARENT */)
		return a.brush
	case wmDestroy:
		syscallN(procPostQuitMessage, 0)
		return 0
	}
	return syscallN(procDefWindowProcW, hwnd, message, wParam, lParam)
}

func (a *app) staticColor(h uintptr) uint32 {
	switch h {
	case a.hStatus:
		switch a.statusKind {
		case statusOK:
			return a.colorOK
		case statusError:
			return a.colorErr
		}
		return a.colorHint
	case a.hSubtitle, a.hDesc, a.hHintBrowser, a.hHintName:
		return a.colorHint
	}
	return a.colorLabel
}

func (a *app) ctl(class, text string, style, exStyle uint32, x, y, w, h int32, id int) uintptr {
	hwnd := syscallN(procCreateWindowExW,
		uintptr(exStyle),
		uintptr(unsafe.Pointer(utf16Ptr(class))),
		uintptr(unsafe.Pointer(utf16Ptr(text))),
		uintptr(style|wsChild|wsVisible),
		uintptr(a.sc(x)), uintptr(a.sc(y)), uintptr(a.sc(w)), uintptr(a.sc(h)),
		a.hwnd, uintptr(id), a.hInst, 0)
	if hwnd != 0 {
		syscallN(procSendMessageW, hwnd, wmSetFont, a.font, 1)
		a.controls = append(a.controls, hwnd)
	}
	return hwnd
}

func (a *app) static(text string, x, y, w, h int32) uintptr {
	return a.ctl("STATIC", text, ssLeft, 0, x, y, w, h, 0)
}

func (a *app) createControls() {
	s := tr(a.lang)

	a.hSubtitle = a.ctl("STATIC", s.Subtitle, ssLeft, 0, 24, 18, 592, 22, idSubtitle)

	a.hLabelURL = a.static(s.LabelURL, 24, 62, 74, 24)
	a.hURL = a.ctl("EDIT", "", wsBorder|wsTabStop|esAutoHScroll, 0, 104, 60, 512, 30, idURL)

	a.hLabelIcon = a.static(s.LabelIcon, 24, 106, 74, 24)
	a.hIconPath = a.ctl("EDIT", "", wsBorder|wsTabStop|esAutoHScroll|esReadOnly, 0, 104, 104, 392, 30, idIconPath)
	a.hPick = a.ctl("BUTTON", s.Browse, wsTabStop, 0, 506, 103, 110, 32, idPick)

	a.hLabelOpenWith = a.static(s.LabelOpenWith, 24, 152, 74, 24)
	a.hBrowser = a.ctl("COMBOBOX", "", cbsDropDownList|wsTabStop|wsVScroll, 0, 104, 148, 260, 220, idBrowser)
	a.hHintBrowser = a.ctl("STATIC", s.HintBrowser, ssLeft, 0, 376, 150, 240, 40, idHintBrowser)

	a.hLabelName = a.static(s.LabelName, 24, 204, 74, 24)
	a.hName = a.ctl("EDIT", "", wsBorder|wsTabStop|esAutoHScroll, 0, 104, 200, 260, 30, idName)
	a.hHintName = a.ctl("STATIC", s.HintName, ssLeft, 0, 376, 206, 240, 22, idHintName)

	a.hLabelPreview = a.ctl("STATIC", s.LabelPreview, ssLeft, 0, 24, 230, 92, 20, 0)
	a.hPreview = a.ctl("STATIC", "", ssIcon|ssCenterImage, wsExClientEdge, 24, 252, 96, 96, idPreview)
	a.hDesc = a.ctl("STATIC", s.Desc, ssLeft, 0, 140, 258, 476, 84, idDesc)

	a.hLabelLang = a.static(s.LabelLang, 24, 374, 64, 22)
	a.hLang = a.ctl("COMBOBOX", "", cbsDropDownList|wsTabStop|wsVScroll, 0, 92, 368, 240, 300, idLang)

	a.hGenerate = a.ctl("BUTTON", s.Generate, bsDefPushButton|wsTabStop, 0, 472, 368, 144, 40, idGenerate)
	a.hStatus = a.ctl("STATIC", s.StatusInit, ssLeft, 0, 24, 424, 592, 22, idStatus)

	for _, l := range languages {
		p := utf16Ptr(l.name)
		syscallN(procSendMessageW, a.hLang, cbAddString, 0, uintptr(unsafe.Pointer(p)))
		runtime.KeepAlive(p)
	}
	if idx := languageIndex(a.lang); idx >= 0 {
		syscallN(procSendMessageW, a.hLang, cbSetCurSel, uintptr(idx), 0)
	}

	a.fillBrowserList(s.DefaultBrowser, 0)

	syscallN(procSetFocus, a.hURL)
}

// fillBrowserList 重建浏览器下拉框的选项，并尽量保留原来的选中项。
func (a *app) fillBrowserList(defaultText string, keep int) {
	p := utf16Ptr(defaultText)
	syscallN(procSendMessageW, a.hBrowser, cbAddString, 0, uintptr(unsafe.Pointer(p)))
	runtime.KeepAlive(p)
	for _, b := range a.browsers {
		bp := utf16Ptr(b.Name)
		syscallN(procSendMessageW, a.hBrowser, cbAddString, 0, uintptr(unsafe.Pointer(bp)))
		runtime.KeepAlive(bp)
	}
	if keep < 0 || keep > len(a.browsers) {
		keep = 0
	}
	syscallN(procSendMessageW, a.hBrowser, cbSetCurSel, uintptr(keep), 0)
}

// applyLanguage 切换界面语言：更新全部文案、字体，并重建依赖语言的列表项。
func (a *app) applyLanguage(code string) {
	a.lang = code
	s := tr(code)
	a.applyFont(code)

	set := func(h uintptr, text string) {
		if h != 0 {
			p := utf16Ptr(text)
			syscallN(procSetWindowTextW, h, uintptr(unsafe.Pointer(p)))
			runtime.KeepAlive(p)
		}
	}
	set(a.hwnd, s.Title)
	set(a.hSubtitle, s.Subtitle)
	set(a.hLabelURL, s.LabelURL)
	set(a.hLabelIcon, s.LabelIcon)
	set(a.hLabelOpenWith, s.LabelOpenWith)
	set(a.hLabelName, s.LabelName)
	set(a.hLabelPreview, s.LabelPreview)
	set(a.hLabelLang, s.LabelLang)
	set(a.hHintBrowser, s.HintBrowser)
	set(a.hHintName, s.HintName)
	set(a.hDesc, s.Desc)
	set(a.hPick, s.Browse)
	set(a.hGenerate, s.Generate)

	if idx := languageIndex(code); idx >= 0 {
		syscallN(procSendMessageW, a.hLang, cbSetCurSel, uintptr(idx), 0)
	}

	keep := int(syscallN(procSendMessageW, a.hBrowser, cbGetCurSel, 0, 0))
	syscallN(procSendMessageW, a.hBrowser, cbResetContent, 0, 0)
	a.fillBrowserList(s.DefaultBrowser, keep)

	if a.statusKind == statusNeutral {
		a.setStatus(s.StatusInit, statusNeutral)
	}
	syscallN(procInvalidateRect, a.hwnd, 0, 1)
}

// applyFont 按界面语言更换字体（确保 CJK / 韩文 / 西里尔字形可正常显示）。
func (a *app) applyFont(code string) {
	if a.font != 0 {
		syscallN(procDeleteObject, a.font)
	}
	a.font = a.createFont(9, fwNormal, fontFor(code))
	for _, h := range a.controls {
		syscallN(procSendMessageW, h, wmSetFont, a.font, 1)
	}
}

func getText(h uintptr) string {
	n := int32(syscallN(procGetWindowTextLengthW, h))
	if n <= 0 {
		return ""
	}
	buf := make([]uint16, n+2)
	syscallN(procGetWindowTextW, h, uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	return syscall.UTF16ToString(buf)
}

func openImageDialog(owner uintptr, s uiStrings) string {
	var buf [1024]uint16
	filter := s.FilterImages + " (*.png;*.jpg;*.jpeg;*.ico;*.gif;*.bmp;*.webp)\x00" +
		"*.png;*.jpg;*.jpeg;*.ico;*.gif;*.bmp;*.webp\x00" +
		s.FilterAll + " (*.*)\x00*.*\x00\x00"
	ofn := openfilename{
		lStructSize:  uint32(unsafe.Sizeof(openfilename{})),
		hwndOwner:    owner,
		lpstrFilter:  utf16Ptr(filter),
		nFilterIndex: 1,
		lpstrFile:    &buf[0],
		nMaxFile:     uint32(len(buf)),
		lpstrTitle:   utf16Ptr(s.DlgPickTitle),
		flags:        ofnFileMustExist | ofnPathMustExist | ofnNoChangeDir | ofnExplorer,
	}
	if syscallN(procGetOpenFileNameW, uintptr(unsafe.Pointer(&ofn))) == 0 {
		return ""
	}
	return syscall.UTF16ToString(buf[:])
}

func (a *app) loadIconFile(path string, size int32) uintptr {
	p := utf16Ptr(path)
	h := syscallN(procLoadImageW, 0, uintptr(unsafe.Pointer(p)), imageIcon,
		uintptr(size), uintptr(size), lrLoadFromFile)
	if h == 0 {
		h = syscallN(procLoadImageW, 0, uintptr(unsafe.Pointer(p)), imageIcon, 0, 0, lrLoadFromFile)
	}
	runtime.KeepAlive(p)
	return h
}

func (a *app) pickImage() {
	s := tr(a.lang)
	path := openImageDialog(a.hwnd, s)
	if path == "" {
		return
	}
	tmpDir := filepath.Join(os.TempDir(), "QuickLinkMakerPreview")
	icoPath, err := core.PrepareIcon(path, "preview", tmpDir)
	if err != nil {
		if isImageError(err) {
			a.alert(errText(err, s))
		} else {
			a.alert(s.ErrImageProcess + err.Error())
		}
		return
	}
	h := a.loadIconFile(icoPath, a.sc(80))
	if h == 0 {
		a.alert(s.ErrImageInvalid)
		return
	}
	if a.hPreviewIcon != 0 {
		syscallN(procDestroyIcon, a.hPreviewIcon)
	}
	a.hPreviewIcon = h
	syscallN(procSendMessageW, a.hPreview, stmSetImage, imageIcon, h)
	a.selectedImage = path
	p := utf16Ptr(path)
	syscallN(procSetWindowTextW, a.hIconPath, uintptr(unsafe.Pointer(p)))
	runtime.KeepAlive(p)
	a.setStatus(s.StatusIconPicked+filepath.Base(path), statusOK)
}

func (a *app) generate() {
	s := tr(a.lang)
	rawURL := strings.TrimSpace(getText(a.hURL))
	if rawURL == "" {
		a.alert(s.ErrNeedURL)
		syscallN(procSetFocus, a.hURL)
		return
	}
	if a.selectedImage == "" {
		a.alert(s.ErrNeedIcon)
		return
	}
	urlStr, err := core.NormalizeURL(rawURL)
	if err != nil {
		a.alert(errText(err, s))
		syscallN(procSetFocus, a.hURL)
		return
	}

	name := strings.TrimSpace(getText(a.hName))
	if name == "" {
		name = core.DefaultName(urlStr)
	}
	name = core.SanitizeName(name)

	iconsDir, err := core.IconsDir()
	if err != nil {
		a.alert(s.ErrIconsDir + err.Error())
		return
	}
	icoPath, err := core.PrepareIcon(a.selectedImage, name, iconsDir)
	if err != nil {
		if isImageError(err) {
			a.alert(errText(err, s))
		} else {
			a.alert(s.ErrIconProcess + err.Error())
		}
		return
	}

	desktop, err := core.DesktopDir()
	if err != nil {
		a.alert(err.Error())
		return
	}

	sel := int(syscallN(procSendMessageW, a.hBrowser, cbGetCurSel, 0, 0))
	var chosen *core.Browser
	note := ""
	if sel > 0 && sel-1 < len(a.browsers) {
		b := a.browsers[sel-1]
		if _, statErr := os.Stat(b.Path); statErr == nil {
			chosen = &b
		} else {
			note = fmt.Sprintf(s.BrowserFallbackFmt, b.Name)
		}
	}

	lnkPath := core.UniquePath(desktop, name, ".lnk")
	desc := fmt.Sprintf(s.ShortcutDescFmt, urlStr)
	openWith := s.SystemDefault

	var createErr error
	if chosen != nil {
		openWith = chosen.Name
		createErr = core.CreateShortcut(lnkPath, chosen.Path, `"`+urlStr+`"`, filepath.Dir(chosen.Path), desc, icoPath)
	} else {
		createErr = core.CreateShortcut(lnkPath, urlStr, "", "", desc, icoPath)
	}
	if createErr != nil {
		a.setStatus(s.ErrCreateFailed+createErr.Error(), statusError)
		a.alert(s.ErrCreateShortcut + createErr.Error())
		return
	}

	a.setStatus(s.StatusCreated+filepath.Base(lnkPath), statusOK)
	a.info(note + fmt.Sprintf(s.InfoCreatedFmt, filepath.Base(lnkPath)) +
		fmt.Sprintf(s.InfoOpenWithFmt, openWith) + s.InfoTail)
}

func (a *app) setStatus(text string, kind int) {
	a.statusKind = kind
	syscallN(procSetWindowTextW, a.hStatus, uintptr(unsafe.Pointer(utf16Ptr(text))))
	syscallN(procInvalidateRect, a.hStatus, 0, 1)
}

func (a *app) alert(msg string) {
	a.messageBox(msg, tr(a.lang).Title, mbOK|mbIconWarning)
}

func (a *app) info(msg string) {
	a.messageBox(msg, tr(a.lang).InfoTitle, mbOK|mbIconInformation)
}

func (a *app) fatal(msg string) {
	a.messageBox(msg, tr(a.lang).Title, mbOK|mbIconError)
	os.Exit(1)
}

func (a *app) messageBox(msg, title string, flags uint32) {
	mp := utf16Ptr(msg)
	tp := utf16Ptr(title)
	syscallN(procMessageBoxW, a.hwnd, uintptr(unsafe.Pointer(mp)), uintptr(unsafe.Pointer(tp)), uintptr(flags))
	runtime.KeepAlive(mp)
	runtime.KeepAlive(tp)
}
