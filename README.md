# 网址快捷方式生成器 / URL Shortcut Maker

把任意网址 + 自选图片图标，一键生成带自定义图标的 Windows 桌面快捷方式。

A small Windows tool that turns a URL plus an image of your choice into a desktop shortcut with that image as its icon. Optionally pick which browser it opens with.

## 功能

- 填写网址即可，自动补全 `https://` 并校验格式
- 选择本地图片作为图标（PNG / JPG / JPEG / ICO / GIF / BMP / WebP），自动转换为 16–256 共 7 种尺寸的 `.ico`
- 可选择用哪个浏览器打开：启动时自动扫描本机已安装的浏览器（Edge、Chrome、Firefox、Brave、Opera、Vivaldi、360、QQ、搜狗、傲游、IE）
- 所选浏览器不存在时自动回落到系统默认浏览器
- 快捷方式名称可自定义，留空则按域名自动生成；重名自动加序号
- 界面支持 10 种语言（简体中文、繁體中文、English、日本語、한국어、Deutsch、Français、Español、Português、Русский），启动时自动跟随系统语言，也可在窗口左下角随时切换，切换即时生效
- 单文件 exe（约 3.2 MB），无需安装、无需运行库，仅支持 64 位 Windows 10 / 11

## 下载

到 [Releases](../../releases) 页面下载 `url-shortcut-maker.exe`，双击运行。

> 程序未做代码签名，首次运行时 Windows 可能提示「已保护你的电脑」，点击「更多信息 → 仍要运行」即可。

## 使用

1. 填写要打开的网址，例如 `www.example.com`
2. 点击「浏览…」选择一张图片，左侧会显示图标预览
3. 选择打开方式（默认浏览器或指定浏览器），可填名称
4. 点击「生成快捷方式」（或直接回车）

桌面会出现带所选图标的快捷方式，双击即用指定浏览器打开该网址。

界面左下角可切换界面语言；首次启动会按 Windows 显示语言自动选择。

## 从源码构建

需要 Go 1.26 或更高版本：

```bash
go install github.com/akavel/rsrc@latest

# 生成内嵌的资源（程序图标 + 清单文件，含 DPI 感知与 Common Controls v6）
rsrc -arch amd64 -manifest app.manifest -ico app.ico -o rsrc.syso

# 编译出无控制台窗口的单文件 exe
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath \
  -ldflags "-H windowsgui -s -w" -o url-shortcut-maker.exe .
```

## 实现说明

- 纯 Win32 API（`syscall` 直接调用 user32 / gdi32 / shell32 / ole32），没有任何 GUI 框架依赖。
- 快捷方式通过 `IShellLink` + `IPersistFile` 生成：
  - 指定浏览器时，目标为该浏览器的 exe，参数为带引号的网址；
  - 使用默认浏览器时，目标直接是网址本身，由系统按默认程序打开。
- 浏览器检测读取注册表 `App Paths`（HKLM / HKCU / WOW6432Node），再回退到常见安装目录，并校验 exe 是否存在。
- 界面文案集中在 `i18n_windows.go`：启动时用 `GetUserDefaultUILanguage` 匹配系统显示语言，未收录的语言回落到英文；切换语言时同步更换字体（雅黑 / Segoe UI / Yu Gothic UI / Malgun Gothic）以保证各语种字形正常显示。
- 转换后的 `.ico` 保存在 `%LOCALAPPDATA%\QuickLinkMaker\icons\`，由快捷方式通过 `IconLocation` 引用。**该目录不要删除**，否则快捷方式图标会丢失。
- 图标转换使用 `golang.org/x/image` 的 Catmull-Rom 缩放；小尺寸写 BMP（32 位 BGRA + AND 掩码），128 / 256 写 PNG，兼容性最好。

## English

**URL Shortcut Maker** — turn any URL plus an image into a Windows desktop shortcut that uses that image as its icon.

Features:

- Enter a URL, `https://` is added automatically and the format is validated
- Pick a local image (PNG / JPG / JPEG / ICO / GIF / BMP / WebP); it is converted into an `.ico` with 7 sizes from 16 to 256 px
- Choose which browser opens the shortcut — installed browsers (Edge, Chrome, Firefox, Brave, Opera, Vivaldi, 360, QQ, Sogou, Maxthon, IE) are detected automatically on startup
- If the chosen browser is not installed, the system default browser is used as a fallback
- Custom shortcut name; defaults to the domain name, duplicates get a numeric suffix
- UI in 10 languages (Chinese Simplified/Traditional, English, Japanese, Korean, German, French, Spanish, Portuguese, Russian) — follows the system UI language on startup and can be switched at the bottom-left of the window
- Single ~3.2 MB exe, no installer and no runtime needed, 64-bit Windows 10 / 11 only

Download `url-shortcut-maker.exe` from [Releases](../../releases) and double-click it. The binary is not code-signed, so SmartScreen may warn on first run — choose *More info → Run anyway*.

Usage: fill in the URL, click *Browse…* to select an icon image, pick the browser (default or specific), then click *Create shortcut*. The shortcut appears on your desktop. The UI language follows the Windows display language on first run and can be changed at the bottom-left of the window.

Build from source with Go 1.26+:

```bash
go install github.com/akavel/rsrc@latest
rsrc -arch amd64 -manifest app.manifest -ico app.ico -o rsrc.syso
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath \
  -ldflags "-H windowsgui -s -w" -o url-shortcut-maker.exe .
```

License: [MIT](LICENSE). Third-party dependencies are BSD 3-Clause (see [THIRD_PARTY_NOTICES.md](THIRD_PARTY_NOTICES.md)).

## 许可证

[MIT](LICENSE)。第三方依赖（Go 标准库、`golang.org/x/image`、`golang.org/x/sys`）均为 BSD 3-Clause，声明见 [THIRD_PARTY_NOTICES.md](THIRD_PARTY_NOTICES.md)。
