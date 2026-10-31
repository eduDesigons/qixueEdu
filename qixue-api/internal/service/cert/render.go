package cert

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"io"
	"log/slog"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode"

	xdraw "golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
	_ "golang.org/x/image/webp"
)

const (
	canvasW = 800
	canvasH = 500
	scale   = 2
)

type certBox struct {
	Type     string   `json:"type"`
	Key      string   `json:"key"`
	Content  string   `json:"content"`
	Src      string   `json:"src"`
	X        float64  `json:"x"`
	Y        float64  `json:"y"`
	W        float64  `json:"w"`
	H        float64  `json:"h"`
	Color    string   `json:"color"`
	Align    string   `json:"align"`
	FontSize *float64 `json:"fontSize"`
}

type renderValues struct {
	Name   string
	Course string
	Time   string
	No     string
}

func (v renderValues) asMap() map[string]string {
	return map[string]string{
		"name": v.Name, "course": v.Course, "time": v.Time, "no": v.No,
	}
}

func (s *Issuer) scheduleRender(userCertID int) {
	if userCertID <= 0 {
		return
	}
	go func() {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("证书异步渲染 panic", "user_cert_id", userCertID, "recover", r)
			}
		}()
		if err := s.RenderAndStoreImage(context.Background(), userCertID); err != nil {
			slog.Error("证书异步渲染失败", "user_cert_id", userCertID, "err", err)
		}
	}()
}

// RenderAndStoreImage renders a certificate PNG and stores the image URL.
func (s *Issuer) RenderAndStoreImage(ctx context.Context, userCertID int) error {
	if userCertID <= 0 {
		return nil
	}
	record, err := s.certs.GetUserCertByID(ctx, userCertID)
	if err != nil {
		return err
	}
	tpl, err := s.certs.LoadTemplateRender(ctx, record.CertID)
	if err != nil {
		return err
	}
	if strings.TrimSpace(tpl.TemplateImage) == "" {
		return fmt.Errorf("证书模板缺少背景图")
	}

	nick, account, _ := s.certs.GetUserDisplay(ctx, record.UserID)
	name := strings.TrimSpace(nick)
	if name == "" {
		name = strings.TrimSpace(account)
	}

	subject := ""
	if record.CourseID > 0 {
		subject, _ = s.certs.GetCourseTitle(ctx, record.CourseID)
	} else if record.PathID > 0 {
		subject, _ = s.certs.GetPathTitle(ctx, record.PathID)
	}

	issued := time.Now()
	if record.IssuedAt != nil {
		issued = *record.IssuedAt
	}
	values := renderValues{
		Name:   name,
		Course: subject,
		Time:   formatCertDate(issued),
		No:     record.CertNo,
	}

	key, err := s.render(tpl.TemplateImage, tpl.Params, values, "", scale)
	if err != nil {
		return err
	}
	return s.certs.UpdateUserCertImage(ctx, userCertID, key, time.Now())
}

func (s *Issuer) render(templateImage string, paramsRaw json.RawMessage, values renderValues, objectKey string, outScale int) (string, error) {
	if outScale <= 0 {
		outScale = scale
	}
	src, err := s.loadImage(templateImage)
	if err != nil {
		return "", err
	}
	outW := canvasW * outScale
	outH := canvasH * outScale
	canvas := image.NewNRGBA(image.Rect(0, 0, outW, outH))
	xdraw.Draw(canvas, canvas.Bounds(), image.White, image.Point{}, xdraw.Src)
	drawCover(canvas, src, outW, outH)

	fontPath := findFont()
	var otFont *opentype.Font
	if fontPath != "" {
		otFont, err = parseFontFile(fontPath)
		if err != nil {
			slog.Error("证书字体加载失败", "font", fontPath, "err", err)
			otFont = nil
		}
	} else {
		slog.Warn("证书字体未找到，将跳过文字绘制（印章与背景仍会合成）; 候选: storage/fonts/NotoSansSC-Medium.ttf, C:/Windows/Fonts/msyh.ttc|msyh.ttf|simhei.ttf")
	}

	boxes := parseBoxes(paramsRaw)
	valMap := values.asMap()
	for _, box := range boxes {
		typ := resolveBoxType(box)
		x := int(math.Round(box.X * float64(outScale)))
		y := int(math.Round(box.Y * float64(outScale)))
		w := int(math.Round(defaultF(box.W, 200) * float64(outScale)))
		h := int(math.Round(defaultF(box.H, 40) * float64(outScale)))

		if typ == "stamp" {
			stampSrc := strings.TrimSpace(box.Src)
			if stampSrc == "" {
				continue
			}
			stamp, loadErr := s.loadImage(stampSrc)
			if loadErr != nil {
				slog.Error("证书图片绘制失败", "src", stampSrc, "err", loadErr)
				continue
			}
			drawStampInBox(canvas, stamp, x, y, w, h)
			continue
		}

		var text string
		if typ == "text" {
			text = interpolate(box.Content, valMap)
		} else {
			text = valMap[box.Key]
		}
		if text == "" || otFont == nil {
			continue
		}
		col := parseHexColor(box.Color)
		align := box.Align
		if align == "" {
			align = "center"
		}
		var drawSize *int
		if box.FontSize != nil && *box.FontSize > 0 {
			sz := int(math.Round(*box.FontSize * float64(outScale)))
			drawSize = &sz
		}
		drawTextInBox(canvas, text, x, y, w, h, otFont, col, align, drawSize, outScale)
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, canvas); err != nil {
		return "", err
	}
	key := objectKey
	if key == "" {
		key = fmt.Sprintf("cert/%s/%s.png", time.Now().Format("200601"), uniqCertName())
	}
	if s.storage == nil {
		return "", fmt.Errorf("Bunny storage 未初始化")
	}
	uploaded, err := s.storage.UploadBytes(buf.Bytes(), key, "image/png")
	if err != nil {
		return "", err
	}
	return uploaded, nil
}

func (s *Issuer) loadImage(pathOrURL string) (image.Image, error) {
	pathOrURL = strings.TrimSpace(pathOrURL)
	if pathOrURL == "" {
		return nil, fmt.Errorf("证书图片路径为空")
	}
	var bin []byte
	var err error
	if strings.HasPrefix(strings.ToLower(pathOrURL), "http://") || strings.HasPrefix(strings.ToLower(pathOrURL), "https://") {
		bin, err = downloadHTTP(pathOrURL)
	} else {
		if s.storage == nil {
			return nil, fmt.Errorf("Bunny storage 未初始化")
		}
		bin, err = s.storage.GetContents(pathOrURL)
		if err != nil {
			return nil, fmt.Errorf("无法从存储读取证书图片: %s (%v)", pathOrURL, err)
		}
	}
	if len(bin) == 0 {
		return nil, fmt.Errorf("证书图片内容为空: %s", pathOrURL)
	}
	img, _, err := image.Decode(bytes.NewReader(bin))
	if err != nil {
		return nil, fmt.Errorf("证书图片无法识别: %s", pathOrURL)
	}
	return img, nil
}

func downloadHTTP(url string) ([]byte, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "QixueCertIssuer/1.0")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("无法下载证书图片: %s", url)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("无法下载证书图片: %s (%s)", url, resp.Status)
	}
	return io.ReadAll(io.LimitReader(resp.Body, 32<<20))
}

func drawCover(dst *image.NRGBA, src image.Image, dw, dh int) {
	sb := src.Bounds()
	sw, sh := sb.Dx(), sb.Dy()
	if sw <= 0 || sh <= 0 {
		return
	}
	sc := math.Max(float64(dw)/float64(sw), float64(dh)/float64(sh))
	nw := int(math.Round(float64(sw) * sc))
	nh := int(math.Round(float64(sh) * sc))
	ox := (dw - nw) / 2
	oy := (dh - nh) / 2
	scaled := image.NewNRGBA(image.Rect(0, 0, nw, nh))
	xdraw.CatmullRom.Scale(scaled, scaled.Bounds(), src, sb, xdraw.Over, nil)
	xdraw.Draw(dst, image.Rect(ox, oy, ox+nw, oy+nh), scaled, image.Point{}, xdraw.Over)
}

func drawStampInBox(dst *image.NRGBA, src image.Image, x, y, w, h int) {
	sb := src.Bounds()
	sw, sh := sb.Dx(), sb.Dy()
	if sw <= 0 || sh <= 0 || w <= 0 || h <= 0 {
		return
	}
	sc := math.Min(float64(w)/float64(sw), float64(h)/float64(sh))
	nw := max(1, int(math.Round(float64(sw)*sc)))
	nh := max(1, int(math.Round(float64(sh)*sc)))
	ox := x + (w-nw)/2
	oy := y + (h-nh)/2
	scaled := image.NewNRGBA(image.Rect(0, 0, nw, nh))
	xdraw.CatmullRom.Scale(scaled, scaled.Bounds(), src, sb, xdraw.Over, nil)
	xdraw.Draw(dst, image.Rect(ox, oy, ox+nw, oy+nh), scaled, image.Point{}, xdraw.Over)
}

func drawTextInBox(dst *image.NRGBA, text string, x, y, w, h int, otFont *opentype.Font, col color.Color, align string, preferredSize *int, outScale int) {
	if otFont == nil || text == "" || w <= 0 || h <= 0 {
		return
	}
	outScale = max(1, outScale)
	maxW := max(8, w)
	maxLines := 6
	text = normalizeCertSpaces(text)

	var size int
	if preferredSize != nil && *preferredSize > 0 {
		size = max(10, min(200, *preferredSize))
	} else {
		size = max(12*outScale, int(math.Round(float64(h)*0.42)))
	}
	face, err := opentype.NewFace(otFont, &opentype.FaceOptions{Size: float64(size), DPI: 72})
	if err != nil {
		return
	}
	defer face.Close()

	indentPx := 0
	body := text
	if align == "left" {
		body, indentPx = pullLeadingIndent(face, size, text)
	}
	lines := wrapText(face, body, maxW, maxLines, indentPx)
	lineH := float64(size) * 1.35
	totalH := lineH * float64(len(lines))
	firstMidY := float64(y) + (float64(h)-totalH)/2 + lineH/2

	src := image.NewUniform(col)
	for i, line := range lines {
		if line == "" {
			continue
		}
		bounds, _ := font.BoundString(face, line)
		midY := firstMidY + float64(i)*lineH
		minY := bounds.Min.Y.Round()
		maxY := bounds.Max.Y.Round()
		ty := int(math.Round(midY - float64(minY+maxY)/2.0))
		lineIndent := 0
		if align == "left" && i == 0 {
			lineIndent = indentPx
		}
		var tx int
		switch align {
		case "left":
			tx = x + lineIndent - bounds.Min.X.Round()
		case "right":
			tx = x + w - bounds.Max.X.Round()
		default:
			tx = x + w/2 - (bounds.Min.X.Round()+bounds.Max.X.Round())/2
		}
		d := &font.Drawer{
			Dst:  dst,
			Src:  src,
			Face: face,
			Dot:  fixed.P(tx, ty),
		}
		d.DrawString(line)
	}
}

func pullLeadingIndent(face font.Face, size int, text string) (string, int) {
	runes := []rune(text)
	i := 0
	for i < len(runes) && unicode.IsSpace(runes[i]) {
		i++
	}
	if i == 0 {
		return text, 0
	}
	ws := string(runes[:i])
	body := string(runes[i:])
	if body == "" {
		return "", 0
	}
	marker := "国"
	with, _ := font.BoundString(face, ws+marker)
	only, _ := font.BoundString(face, marker)
	indent := (with.Max.X - with.Min.X).Round() - (only.Max.X - only.Min.X).Round()
	if indent <= 0 {
		indent = int(math.Round(float64(i) * float64(size) * 0.33))
	}
	if indent < 0 {
		indent = 0
	}
	return body, indent
}

func wrapText(face font.Face, text string, maxW, maxLines, firstLineIndent int) []string {
	maxLines = max(1, maxLines)
	paragraphs := regexp.MustCompile(`\r\n|\n|\r`).Split(text, -1)
	lines := make([]string, 0, 4)
	isFirstLine := true
	for _, para := range paragraphs {
		if para == "" {
			if len(lines) < maxLines {
				lines = append(lines, "")
				isFirstLine = false
			}
			continue
		}
		chars := []rune(para)
		cur := ""
		for _, ch := range chars {
			try := cur + string(ch)
			limit := maxW
			if isFirstLine && firstLineIndent > 0 {
				limit = max(8, maxW-firstLineIndent)
			}
			bounds, _ := font.BoundString(face, try)
			tw := (bounds.Max.X - bounds.Min.X).Round()
			if cur != "" && tw > limit {
				lines = append(lines, cur)
				cur = string(ch)
				isFirstLine = false
				if len(lines) >= maxLines {
					return lines
				}
			} else {
				cur = try
			}
		}
		if cur != "" && len(lines) < maxLines {
			lines = append(lines, cur)
			isFirstLine = false
		}
		if len(lines) >= maxLines {
			break
		}
	}
	if len(lines) == 0 {
		return []string{text}
	}
	return lines
}

func resolveBoxType(box certBox) string {
	typ := strings.TrimSpace(box.Type)
	switch typ {
	case "field", "text", "stamp":
		return typ
	case "image":
		return "stamp"
	}
	key := box.Key
	if strings.HasPrefix(key, "text_") || strings.HasPrefix(key, "text-") {
		return "text"
	}
	if strings.HasPrefix(key, "stamp_") || strings.HasPrefix(key, "stamp-") ||
		strings.HasPrefix(key, "image_") || strings.HasPrefix(key, "image-") {
		return "stamp"
	}
	return "field"
}

var interpolateRe = regexp.MustCompile(`(?i)[\{\x{FF5B}]\s*(name|course|time|no)\s*[\}\x{FF5D}]`)

func interpolate(content string, values map[string]string) string {
	content = normalizeCertSpaces(content)
	return interpolateRe.ReplaceAllStringFunc(content, func(m string) string {
		sub := interpolateRe.FindStringSubmatch(m)
		if len(sub) < 2 {
			return ""
		}
		return values[strings.ToLower(sub[1])]
	})
}

func normalizeCertSpaces(text string) string {
	text = strings.ReplaceAll(text, "\u3000", "   ")
	var b strings.Builder
	b.Grow(len(text))
	for _, r := range text {
		switch r {
		case '\u00A0', '\u2000', '\u2001', '\u2002', '\u2003', '\u2004',
			'\u2005', '\u2006', '\u2007', '\u2008', '\u2009', '\u200A',
			'\u202F', '\u205F':
			b.WriteByte(' ')
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

func formatCertDate(t time.Time) string {
	loc := time.FixedZone("CST", 8*3600)
	if l, err := time.LoadLocation("Asia/Shanghai"); err == nil {
		loc = l
	}
	t = t.In(loc)
	return fmt.Sprintf("%d年%d月%d日", t.Year(), int(t.Month()), t.Day())
}

func parseBoxes(raw json.RawMessage) []certBox {
	if len(raw) == 0 || string(raw) == "null" {
		return nil
	}
	var boxes []certBox
	if err := json.Unmarshal(raw, &boxes); err == nil {
		return boxes
	}
	var asString string
	if json.Unmarshal(raw, &asString) == nil {
		asString = strings.TrimSpace(asString)
		if asString == "" {
			return nil
		}
		_ = json.Unmarshal([]byte(asString), &boxes)
	}
	return boxes
}

func parseHexColor(hexStr string) color.NRGBA {
	hexStr = strings.TrimPrefix(strings.TrimSpace(hexStr), "#")
	if len(hexStr) != 6 {
		hexStr = "333333"
	}
	var r, g, b uint8
	if _, err := fmt.Sscanf(hexStr, "%02x%02x%02x", &r, &g, &b); err != nil {
		return color.NRGBA{R: 0x33, G: 0x33, B: 0x33, A: 255}
	}
	return color.NRGBA{R: r, G: g, B: b, A: 255}
}

func parseFontFile(path string) (*opentype.Font, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	lower := strings.ToLower(path)
	if strings.HasSuffix(lower, ".ttc") {
		coll, err := opentype.ParseCollection(data)
		if err != nil {
			return nil, err
		}
		if coll.NumFonts() < 1 {
			return nil, fmt.Errorf("empty font collection")
		}
		return coll.Font(0)
	}
	return opentype.Parse(data)
}

// findFont tries Noto Sans SC under qixue-api/storage/fonts, then Windows 微软雅黑 / 黑体.
// Walks up from cwd so it works when the process starts in qixue-api or a subdir.
func findFont() string {
	relNames := []string{
		filepath.Join("storage", "fonts", "NotoSansSC-Medium.ttf"),
		filepath.Join("storage", "fonts", "simhei.ttf"),
	}
	if dir, err := os.Getwd(); err == nil {
		for i := 0; i < 8; i++ {
			for _, name := range relNames {
				p := filepath.Join(dir, name)
				if fileExists(p) {
					return p
				}
			}
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
			dir = parent
		}
	}
	for _, f := range []string{
		`C:\Windows\Fonts\msyh.ttc`,
		`C:\Windows\Fonts\msyh.ttf`,
		`C:\Windows\Fonts\simhei.ttf`,
	} {
		if fileExists(f) {
			return f
		}
	}
	return ""
}

func fileExists(path string) bool {
	st, err := os.Stat(path)
	return err == nil && !st.IsDir()
}

func uniqCertName() string {
	var b [8]byte
	_, _ = rand.Read(b[:])
	return "uc_" + hex.EncodeToString(b[:]) + fmt.Sprintf("%x", time.Now().UnixNano()&0xffff)
}

func defaultF(v, def float64) float64 {
	if v == 0 {
		return def
	}
	return v
}
