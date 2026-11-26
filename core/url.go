package core

import (
	"errors"
	"net/url"
	"regexp"
	"strings"
)

var schemeRe = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9+.\-]*://`)

// 网址校验的哨兵错误，便于界面层本地化。
var (
	ErrURLRequired          = errors.New("网址不能为空")
	ErrURLInvalid           = errors.New("网址格式不正确")
	ErrURLHostRequired      = errors.New("请填写完整的网址，例如 www.example.com")
	ErrURLSchemeUnsupported = errors.New("仅支持 http / https 网址")
)

// NormalizeURL 补全协议并校验网址。
func NormalizeURL(raw string) (string, error) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return "", ErrURLRequired
	}
	if !schemeRe.MatchString(s) {
		s = "https://" + s
	}
	u, err := url.Parse(s)
	if err != nil {
		return "", ErrURLInvalid
	}
	if u.Host == "" || u.Hostname() == "" {
		return "", ErrURLHostRequired
	}
	switch strings.ToLower(u.Scheme) {
	case "http", "https":
	default:
		return "", ErrURLSchemeUnsupported
	}
	return u.String(), nil
}

// DefaultName 根据网址推导一个合适的快捷方式名称。
func DefaultName(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil || u.Hostname() == "" {
		return "网站快捷方式"
	}
	host := strings.TrimPrefix(strings.ToLower(u.Hostname()), "www.")
	path := strings.Trim(u.Path, "/")
	if path == "" {
		return host
	}
	first := strings.SplitN(path, "/", 2)[0]
	if len([]rune(first)) > 24 {
		first = string([]rune(first)[:24])
	}
	if first == "" {
		return host
	}
	return host + " - " + first
}
