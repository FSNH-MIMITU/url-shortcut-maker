# 网址快捷方式生成器 / URL Shortcut Maker

把任意网址 + 自选图片图标，一键生成带自定义图标的 Windows 桌面快捷方式。

A small Windows tool that turns a URL plus an image of your choice into a desktop shortcut with that image as its icon. Optionally pick which browser it opens with.

## 功能

- 填写网址即可，自动补全 `https://` 并校验格式
- 选择本地图片作为图标（PNG / JPG / JPEG / ICO / GIF / BMP / WebP），自动转换为 16–256 共 7 种尺寸的 `.ico`
- 可选择用哪个浏览器打开：启动时自动扫描本机已安装的浏览器（Edge、Chrome、Firefox、Brave、Opera、Vivaldi、360、QQ、搜狗、傲游、IE）
- 所选浏览器不存在时自动回落到系统默认浏览器
- 快捷方式名称可自定义，留空则按域名自动生成；重名自动加序号
- 单文件 exe（约 3 MB），无需安装、无需运行库，仅支持 64 位 Windows 10 / 11

## 下载

到 [Releases](../../releases) 页面下载 `url-shortcut-maker.exe`，双击运行。

> 程序未做代码签名，首次运行时 Windows 可能提示「已保护你的电脑」，点击「更多信息 → 仍要运行」即可。

## 使用

1. 填写要打开的网址，例如 `www.example.com`
2. 点击「浏览…」选择一张图片，左侧会显示图标预览
3. 选择打开方式（默认浏览器或指定浏览器），可填名称
4. 点击「生成快捷方式」（或直接回车）

桌面会出现带所选图标的快捷方式，双击即用指定浏览器打开该网址。

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
- 转换后的 `.ico` 保存在 `%LOCALAPPDATA%\QuickLinkMaker\icons\`，由快捷方式通过 `IconLocation` 引用。**该目录不要删除**，否则快捷方式图标会丢失。
- 图标转换使用 `golang.org/x/image` 的 Catmull-Rom 缩放；小尺寸写 BMP（32 位 BGRA + AND 掩码），128 / 256 写 PNG，兼容性最好。

## 许可证

[MIT](LICENSE)。第三方依赖（Go 标准库、`golang.org/x/image`、`golang.org/x/sys`）均为 BSD 3-Clause，声明见 [THIRD_PARTY_NOTICES.md](THIRD_PARTY_NOTICES.md)。
