package core

import (
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sys/windows/registry"
)

// Browser 表示一个已安装的浏览器。
type Browser struct {
	Name string
	Path string
}

type browserCandidate struct {
	name  string
	exe   string
	paths []string
}

func candidates() []browserCandidate {
	pf := os.Getenv("ProgramFiles")
	pf86 := os.Getenv("ProgramFiles(x86)")
	lad := os.Getenv("LOCALAPPDATA")
	roam := os.Getenv("APPDATA")
	join := filepath.Join

	var list []browserCandidate
	add := func(name, exe string, paths ...string) {
		list = append(list, browserCandidate{name: name, exe: exe, paths: paths})
	}

	add("Microsoft Edge", "msedge.exe",
		join(pf86, `Microsoft\Edge\Application\msedge.exe`),
		join(pf, `Microsoft\Edge\Application\msedge.exe`))
	add("Google Chrome", "chrome.exe",
		join(pf, `Google\Chrome\Application\chrome.exe`),
		join(pf86, `Google\Chrome\Application\chrome.exe`),
		join(lad, `Google\Chrome\Application\chrome.exe`))
	add("Mozilla Firefox", "firefox.exe",
		join(pf, `Mozilla Firefox\firefox.exe`),
		join(pf86, `Mozilla Firefox\firefox.exe`))
	add("Brave", "brave.exe",
		join(pf, `BraveSoftware\Brave-Browser\Application\brave.exe`),
		join(pf86, `BraveSoftware\Brave-Browser\Application\brave.exe`))
	add("Opera", "opera.exe",
		join(lad, `Programs\Opera\opera.exe`),
		join(roam, `Opera\opera.exe`))
	add("Vivaldi", "vivaldi.exe",
		join(lad, `Vivaldi\Application\vivaldi.exe`),
		join(pf, `Vivaldi\Application\vivaldi.exe`))
	add("360 安全浏览器", "360se.exe",
		join(roam, `360se6\Application\360se.exe`),
		join(pf86, `360\360se6\Application\360se.exe`))
	add("360 极速浏览器", "360chrome.exe",
		join(roam, `360Chrome\Chrome\Application\360chrome.exe`),
		join(pf86, `360Chrome\Chrome\Application\360chrome.exe`))
	add("QQ 浏览器", "QQBrowser.exe",
		join(pf86, `Tencent\QQBrowser\QQBrowser.exe`),
		join(pf, `Tencent\QQBrowser\QQBrowser.exe`))
	add("搜狗浏览器", "SogouExplorer.exe",
		join(pf86, `SogouExplorer\SogouExplorer.exe`),
		join(pf, `SogouExplorer\SogouExplorer.exe`))
	add("傲游浏览器", "Maxthon.exe",
		join(pf86, `Maxthon\Application\Maxthon.exe`),
		join(pf, `Maxthon\Application\Maxthon.exe`))
	add("Internet Explorer", "iexplore.exe",
		join(pf, `Internet Explorer\iexplore.exe`),
		join(pf86, `Internet Explorer\iexplore.exe`))

	return list
}

// DetectBrowsers 扫描注册表与常见安装目录，返回已安装的浏览器列表。
func DetectBrowsers() []Browser {
	var out []Browser
	seen := map[string]bool{}
	for _, c := range candidates() {
		path := resolveBrowser(c)
		if path == "" {
			continue
		}
		key := strings.ToLower(filepath.Clean(path))
		if seen[key] {
			continue
		}
		seen[key] = true
		out = append(out, Browser{Name: c.name, Path: path})
	}
	return out
}

func resolveBrowser(c browserCandidate) string {
	if p := registryAppPath(c.exe); isExe(p) {
		return p
	}
	for _, p := range c.paths {
		if isExe(p) {
			return p
		}
	}
	return ""
}

func registryAppPath(exe string) string {
	sub := `Microsoft\Windows\CurrentVersion\App Paths\` + exe
	lookups := []struct {
		root registry.Key
		path string
	}{
		{registry.LOCAL_MACHINE, `SOFTWARE\` + sub},
		{registry.LOCAL_MACHINE, `SOFTWARE\WOW6432Node\` + sub},
		{registry.CURRENT_USER, `SOFTWARE\` + sub},
	}
	for _, l := range lookups {
		k, err := registry.OpenKey(l.root, l.path, registry.QUERY_VALUE)
		if err != nil {
			continue
		}
		v, _, err := k.GetStringValue("")
		k.Close()
		if err == nil {
			if p := strings.Trim(strings.TrimSpace(v), `"`); isExe(p) {
				return p
			}
		}
	}
	return ""
}

func isExe(path string) bool {
	if path == "" || !strings.EqualFold(filepath.Ext(path), ".exe") {
		return false
	}
	st, err := os.Stat(path)
	return err == nil && !st.IsDir()
}
