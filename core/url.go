package core

import (
	"errors"
	"net/url"
	"regexp"
	"strings"
)

var schemeRe = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9+.\-]*://`)

// NormalizeURL 补全协议并校验网址。
func NormalizeURL(raw string) (string, error) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return "", errors.New("网址不能为空")
	}
	if !schemeRe.MatchString(s) {
		s = "https://" + s
	}
	u, err := url.Parse(s)
	if err != nil {
		return "", errors.New("网址格式不正确")
	}
	if u.Host == "" || u.Hostname() == "" {
		return "", errors.New("请填写完整的网址，例如 www.example.com")
	}
	switch strings.ToLower(u.Scheme) {
	case "http", "https":
	default:
		return "", errors.New("仅支持 http / https 网址")
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
