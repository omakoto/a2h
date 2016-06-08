package engine

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/omakoto/mlib"
	"io"
	"math"
	"strconv"
	"strings"
	"text/template"
)

var (
	StandardColors = []int{
		// Standard VGA (http://en.wikipedia.org/wiki/ANSI_escape_code#Colors)
		rgb(0, 0, 0),
		rgb(170, 0, 0),
		rgb(0, 170, 0),
		rgb(170, 85, 0),
		rgb(0, 0, 170),
		rgb(170, 0, 170),
		rgb(0, 170, 170),
		rgb(170, 170, 170),
	}

	BrightColors = []int{
		rgb(85, 85, 85),
		rgb(255, 85, 85),
		rgb(85, 255, 85),
		rgb(255, 255, 85),
		rgb(85, 85, 255),
		rgb(255, 85, 255),
		rgb(85, 255, 255),
		rgb(255, 255, 255),
	}
)

var (
	Gamma             = flag.Float64("gamma", 1.0, "gamma value for RGB conversion")
	Title             = flag.String("title", "A2H", "HTML Title")
	BgColor           = flag.String("bg-color", "#000000", "Background color")
	TextColor         = flag.String("text-color", "#ffffff", "Background color")
	FontSize          = flag.String("font-size", "9pt", "Font size)")
	AutoFlash         = flag.Bool("auto-flush", false, "Auto flush")
	NoConvertControls = flag.Bool("no-convert-controls", false, "Don't convert control characters")
)

// Color manipulation

func rgb(r, g, b int) int {
	return r<<16 | g<<8 | b
}

func gamma(v float64) float64 {
	return math.Max(0, math.Min(1, math.Pow(v, *Gamma)))
}

func gammaRgb(rgbValue int) int {
	r := gamma(float64((rgbValue>>16)&255) / 255.0)
	g := gamma(float64((rgbValue>>8)&255) / 255.0)
	b := gamma(float64((rgbValue)&255) / 255.0)
	return rgb(int(r*255), int(g*255), int(b*255))
}

func xterm256ColortoRgb(value int) int {
	if value < 8 {
		return StandardColors[value]
	}
	if value < 16 {
		return BrightColors[value-8]
	}
	if 232 <= value && value <= 256 {
		// Gray
		level := (value-232)*10 + 8
		return rgb(level, level, level)
	}

	value -= 16

	b := value % 6
	g := (value / 6) % 6
	r := (value / 36) % 6
	return rgb(int(float64(r)*42.5), int(float64(g)*42.5), int(float64(b)*42.5))
}

func getIndexColor(index int, bold bool) int {
	if bold {
		return BrightColors[index]
	} else {
		return StandardColors[index]
	}
}

// Converter
const (
	defaultColor = -1000
)

type Converter struct {
	fg, bg int // positive: rgb, negative: index, -1000

	bold      bool
	faint     bool
	italic    bool
	underline bool
	blink     bool
	negative  bool
	conceal   bool
	crossout  bool

	inDiv  bool
	inSpan bool

	rows int

	buf *bufio.Writer
}

func NewConverter(w io.Writer) Converter {
	c := Converter{buf: bufio.NewWriter(w)}
	c.reset()
	return c
}

func (c *Converter) reset() {
	c.fg = defaultColor
	c.bg = defaultColor
	c.bold = false
	c.faint = false
	c.italic = false
	c.underline = false
	c.blink = false
	c.negative = false
	c.conceal = false
	c.crossout = false
}

func (c *Converter) hasAttr() bool {
	return (c.fg != defaultColor ||
		c.bg != defaultColor ||
		c.bold ||
		c.faint ||
		c.italic ||
		c.underline ||
		c.blink ||
		c.negative ||
		c.conceal ||
		c.crossout)
}

func (c *Converter) startDiv() {
	if !c.inDiv {
		c.inDiv = true
		c.buf.WriteString("<div>")
		c.rows++
	}
}

func (c *Converter) closeDiv() {
	c.closeSpan()
	if c.inDiv {
		c.inDiv = false
		c.buf.WriteString("</div>\n")
	}
}

func (c *Converter) closeSpan() {
	if c.inSpan {
		c.inSpan = false
		c.buf.WriteString("</span>")
	}
}

func parseInt(s string, default_ int) int {
	v, err := strconv.Atoi(s)
	if err == nil {
		return v
	} else {
		return default_
	}
}

func setColorForRgb(i int, vals []string) (int, int) {
	ret := 0
	next := vals[i]
	if next == "5" {
		if i+1 < len(vals) {
			// Xterm 256 colors
			i++
			ret = xterm256ColortoRgb(parseInt(vals[i], 0))
			i++
		}
	} else if next == "2" {
		// Kterm 24 bit color
		if i+3 < len(vals) {
			i++
			ret = rgb(parseInt(vals[i], 0), parseInt(vals[i+1], 0), parseInt(vals[i+2], 0))
			i += 3
		}
	}
	return ret, i
}

func (c *Converter) convertCsi(csi string) {
	vals := strings.Split(csi, ";")

	for i := 0; i < len(vals); {
		code := parseInt(vals[i], 0) // first code
		i += 1
		if code == 0 {
			c.reset()
		} else if code == 1 {
			c.bold = true
		} else if code == 2 {
			c.faint = true
		} else if code == 3 {
			c.italic = true
		} else if code == 4 {
			c.underline = true
		} else if code == 5 {
			c.blink = true
		} else if code == 7 {
			c.negative = true
		} else if code == 8 {
			c.conceal = true
		} else if code == 9 {
			c.crossout = true
		} else if code == 21 {
			c.bold = false
		} else if code == 22 {
			c.bold = false
			c.faint = false
		} else if code == 23 {
			c.italic = false
		} else if code == 24 {
			c.underline = false
		} else if code == 25 {
			c.blink = false
		} else if code == 27 {
			c.negative = false
		} else if code == 28 {
			c.conceal = false
		} else if code == 29 {
			c.crossout = false
		} else if 30 <= code && code <= 37 {
			c.fg = -(code - 30 + 1) // FG color, index
		} else if 40 <= code && code <= 47 {
			c.bg = -(code - 40 + 1) // BG color, index
		} else if code == 38 {
			c.fg, i = setColorForRgb(i, vals)
		} else if code == 48 {
			c.bg, i = setColorForRgb(i, vals)
		} else {
			// Unknown
		}
	}

	if !c.hasAttr() {
		c.closeSpan()
		return
	}

	fg := c.fg
	bg := c.bg
	// Convert index color to RGB
	if fg < 0 && fg != defaultColor {
		fg = getIndexColor(-fg-1, c.bold)
	}
	if bg < 0 && bg != defaultColor {
		bg = getIndexColor(-bg-1, false)
	}

	c.closeSpan()
	c.buf.WriteString("<span ")
	if c.blink {
		c.buf.WriteString("class=\"blink\" ")
	}
	c.buf.WriteString("style=\"")
	c.inSpan = true

	if c.bold {
		c.buf.WriteString("font-weight:bold;")
	}
	if c.faint {
		c.buf.WriteString("opacity:0.5;")
	}
	if c.italic {
		c.buf.WriteString("font-style:italic;")
	}
	if c.underline {
		c.buf.WriteString("text-decoration:underline;")
	}
	if c.crossout {
		c.buf.WriteString("text-decoration:line-through;")
	}
	var b, f string
	if bg == defaultColor {
		b = *BgColor
	} else {
		b = fmt.Sprintf("#%06x", gammaRgb(bg))
	}

	if fg == defaultColor {
		f = *TextColor
	} else {
		f = fmt.Sprintf("#%06x", gammaRgb(fg))
	}
	if c.negative {
		f, b = b, f
	}
	if c.conceal {
		f = b
	}

	if f != *TextColor {
		c.buf.WriteString(fmt.Sprintf("color:%s;", f))
	}
	if b != *BgColor {
		c.buf.WriteString(fmt.Sprintf("background-color:%s;", b))
	}

	c.buf.WriteString("\">")

}

func isCsiEnd(b byte) bool {
	return 64 <= b && b <= 126
}

func peek(line string, index int) int {
	if index < len(line) {
		return int(line[index])
	} else {
		return -1
	}
}

func (c *Converter) convert(line string) {
	c.startDiv()

	size := len(line)
outer:
	for i := 0; i < size; i++ {
		b := line[i]
		switch b {
		case '&':
			c.buf.WriteString("&amp;")
			continue
		case '<':
			c.buf.WriteString("&lt;")
			continue
		case '>':
			c.buf.WriteString("&gt;")
			continue
		case '\a': // bell, ignore.
			continue
		case 0x0d:
			if peek(line, i+1) == 0x0a {
				// CR followed by LF, ignore.
				continue
			}
			fallthrough
		case 0x0a:
			c.closeDiv()
			if peek(line, i+1) != -1 {
				c.startDiv()
			}
			continue
		case '\x1b':
			i++
			switch peek(line, i) {
			case -1:
				continue
			case '[': // CSI
				i++
				csiStart := i
				for i < size && !isCsiEnd(line[i]) {
					i++
				}
				if i >= size {
					continue
				}
				if line[i] == 'm' {
					c.convertCsi(line[csiStart:i])
				}
				continue
			case ']':
				i++
				for {
					n := peek(line, i)
					if n == -1 || n == '\a' {
						continue outer
					}
					if n == '\x1b' && peek(line, i+1) == '\\' {
						i++
						continue outer
					}
					i++
				}
				continue
			case '(':
				i++
				continue
			case 'c':
				c.reset()
				c.closeSpan()
				continue
			default: // just eat the next byte
				continue
			}
			continue
		}
		if !*NoConvertControls && 0 <= b && b <= 31 && b != '\t' {
			c.buf.WriteByte('^')
			c.buf.WriteByte(b + '@')
			continue
		}
		c.buf.WriteByte(b)
	}
	c.closeDiv()
}

type TemplateParams struct {
	Title           string
	BackgroundColor string
	TextColor       string
	FontSize        string
	RowCount        int
}

// TODO Get the input from argument too.
func (c *Converter) Convert() {
	defer c.buf.Flush()

	// Header
	params := TemplateParams{
		Title:           *Title,
		BackgroundColor: *BgColor,
		TextColor:       *TextColor,
		FontSize:        *FontSize,
	}

	tmpl, err := template.New("h").Parse(HtmlHeader)
	mlib.Check(err)
	err = tmpl.Execute(c.buf, params)
	mlib.Check(err)

	// Body
	for line := range mlib.ReadFilesFromArgs() {
		c.convert(line)
		if *AutoFlash {
			c.buf.Flush()
		}
	}

	// Footer
	params.RowCount = c.rows
	tmpl, err = template.New("f").Parse(HtmlFooter)
	mlib.Check(err)
	err = tmpl.Execute(c.buf, params)
	mlib.Check(err)
}
