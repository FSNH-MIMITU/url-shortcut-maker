package core

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"image"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"strings"

	xdraw "golang.org/x/image/draw"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/webp"
	_ "image/gif"
	_ "image/jpeg"
)

// 生成的 ico 中包含的尺寸（从大到小）。
var icoSizes = []int{256, 128, 64, 48, 32, 24, 16}

// 图片处理的哨兵错误，便于界面层本地化。
var (
	ErrImageUnsupported = errors.New("不支持的图片格式")
	ErrImageEmpty       = errors.New("图片为空")
	ErrImageSizeInvalid = errors.New("图片尺寸无效")
	ErrImageTooLarge    = errors.New("图片像素过大，请使用较小的图片")
	ErrImageRead        = errors.New("无法读取图片文件")
)

type icoEntry struct {
	size  int
	data  []byte
	isPNG bool
}

// DecodeImage 解码常见格式的图片文件。
func DecodeImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, ErrImageUnsupported
	}
	return img, nil
}

// BuildICO 将任意图片转换为多尺寸 ico 字节。
func BuildICO(src image.Image) ([]byte, error) {
	if src == nil {
		return nil, ErrImageEmpty
	}
	b := src.Bounds()
	if b.Dx() <= 0 || b.Dy() <= 0 {
		return nil, ErrImageSizeInvalid
	}
	if b.Dx()*b.Dy() > 80*1000*1000 {
		return nil, ErrImageTooLarge
	}

	entries := make([]icoEntry, 0, len(icoSizes))
	for _, size := range icoSizes {
		canvas := fitSquare(src, size)
		if size >= 128 {
			var buf bytes.Buffer
			if err := png.Encode(&buf, canvas); err != nil {
				return nil, err
			}
			entries = append(entries, icoEntry{size: size, data: buf.Bytes(), isPNG: true})
		} else {
			entries = append(entries, icoEntry{size: size, data: encodeDIB(canvas)})
		}
	}

	var out bytes.Buffer
	// ICONDIR
	binary.Write(&out, binary.LittleEndian, uint16(0)) // reserved
	binary.Write(&out, binary.LittleEndian, uint16(1)) // type = icon
	binary.Write(&out, binary.LittleEndian, uint16(len(entries)))

	offset := 6 + 16*len(entries)
	for _, e := range entries {
		var dim byte
		if e.size >= 256 {
			dim = 0
		} else {
			dim = byte(e.size)
		}
		out.WriteByte(dim)
		out.WriteByte(dim)
		out.WriteByte(0)                                    // color count
		out.WriteByte(0)                                    // reserved
		binary.Write(&out, binary.LittleEndian, uint16(1))  // planes
		binary.Write(&out, binary.LittleEndian, uint16(32)) // bit count
		binary.Write(&out, binary.LittleEndian, uint32(len(e.data)))
		binary.Write(&out, binary.LittleEndian, uint32(offset))
		offset += len(e.data)
	}
	for _, e := range entries {
		out.Write(e.data)
	}
	return out.Bytes(), nil
}

// fitSquare 等比缩放并居中放置到 size×size 的透明画布上。
func fitSquare(src image.Image, size int) *image.NRGBA {
	dst := image.NewNRGBA(image.Rect(0, 0, size, size))
	b := src.Bounds()
	sw, sh := b.Dx(), b.Dy()
	if sw <= 0 || sh <= 0 {
		return dst
	}
	scale := math.Min(float64(size)/float64(sw), float64(size)/float64(sh))
	dw := int(math.Round(float64(sw) * scale))
	dh := int(math.Round(float64(sh) * scale))
	if dw < 1 {
		dw = 1
	}
	if dh < 1 {
		dh = 1
	}
	ox, oy := (size-dw)/2, (size-dh)/2
	xdraw.CatmullRom.Scale(dst, image.Rect(ox, oy, ox+dw, oy+dh), src, b, xdraw.Over, nil)
	return dst
}

// encodeDIB 输出 32 位 BGRA（自下而上）像素 + 1bpp AND 掩码。
func encodeDIB(img *image.NRGBA) []byte {
	size := img.Bounds().Dx()
	rowBytes := size * 4
	maskRow := ((size + 31) / 32) * 4

	hdr := make([]byte, 40)
	binary.LittleEndian.PutUint32(hdr[0:], 40)
	binary.LittleEndian.PutUint32(hdr[4:], uint32(size))
	binary.LittleEndian.PutUint32(hdr[8:], uint32(size*2)) // 高度包含掩码
	binary.LittleEndian.PutUint16(hdr[12:], 1)
	binary.LittleEndian.PutUint16(hdr[14:], 32)
	binary.LittleEndian.PutUint32(hdr[16:], 0) // BI_RGB
	binary.LittleEndian.PutUint32(hdr[20:], uint32(rowBytes*size))

	var buf bytes.Buffer
	buf.Write(hdr)
	for y := size - 1; y >= 0; y-- {
		for x := 0; x < size; x++ {
			i := img.PixOffset(x, y)
			buf.Write([]byte{img.Pix[i+2], img.Pix[i+1], img.Pix[i], img.Pix[i+3]})
		}
	}
	for y := size - 1; y >= 0; y-- {
		row := make([]byte, maskRow)
		for x := 0; x < size; x++ {
			if img.Pix[img.PixOffset(x, y)+3] < 128 {
				row[x/8] |= 0x80 >> uint(x%8)
			}
		}
		buf.Write(row)
	}
	return buf.Bytes()
}

// IsICOData 判断数据是否为 ico 文件头。
func IsICOData(data []byte) bool {
	return len(data) >= 6 && data[0] == 0 && data[1] == 0 && data[2] == 1 && data[3] == 0
}

// PrepareIcon 把用户选择的图片转换成 ico 文件，返回 ico 路径。
// 若源文件本身就是 ico，则直接复用。
func PrepareIcon(srcPath, baseName, destDir string) (string, error) {
	data, err := os.ReadFile(srcPath)
	if err != nil {
		return "", fmt.Errorf("%w：%v", ErrImageRead, err)
	}
	sum := sha1.Sum(data)
	dst := filepath.Join(destDir, SanitizeName(baseName)+"-"+hex.EncodeToString(sum[:5])+".ico")

	if IsICOData(data) {
		if _, statErr := os.Stat(dst); statErr == nil {
			return dst, nil
		}
		if err := os.MkdirAll(destDir, 0o755); err != nil {
			return "", err
		}
		if err := os.WriteFile(dst, data, 0o644); err != nil {
			return "", err
		}
		return dst, nil
	}

	img, err := DecodeImage(srcPath)
	if err != nil {
		return "", err
	}
	ico, err := BuildICO(img)
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return "", err
	}
	if err := os.WriteFile(dst, ico, 0o644); err != nil {
		return "", err
	}
	return dst, nil
}

// SanitizeName 清理文件名中的非法字符。
func SanitizeName(s string) string {
	s = strings.Map(func(r rune) rune {
		switch r {
		case '\\', '/', ':', '*', '?', '"', '<', '>', '|':
			return '_'
		}
		if r < 0x20 {
			return -1
		}
		return r
	}, s)
	s = strings.TrimSpace(s)
	s = strings.Trim(s, " .")
	if rs := []rune(s); len(rs) > 60 {
		s = string(rs[:60])
		s = strings.Trim(strings.TrimSpace(s), " .")
	}
	// 避免出现 xxx.lnk.lnk 这类重复后缀
	lower := strings.ToLower(s)
	if strings.HasSuffix(lower, ".lnk") || strings.HasSuffix(lower, ".url") {
		s = strings.TrimSpace(s[:len(s)-4])
	}
	if s == "" {
		s = "快捷方式"
	}
	return s
}

// UniquePath 在目录下生成不重名的文件路径。
func UniquePath(dir, base, ext string) string {
	path := filepath.Join(dir, base+ext)
	if !fileExists(path) {
		return path
	}
	for i := 2; i < 1000; i++ {
		path = filepath.Join(dir, fmt.Sprintf("%s (%d)%s", base, i, ext))
		if !fileExists(path) {
			return path
		}
	}
	return filepath.Join(dir, fmt.Sprintf("%s-%d%s", base, os.Getpid(), ext))
}

func fileExists(path string) bool {
	st, err := os.Stat(path)
	return err == nil && !st.IsDir()
}
