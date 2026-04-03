package awtree

// Attr is a bitmask of cell style attributes.
type Attr uint16

const (
	AttrBold          Attr = 1 << iota
	AttrFaint              //nolint:revive // consistent naming
	AttrItalic             //nolint:revive
	AttrUnderline          //nolint:revive
	AttrBlink              //nolint:revive
	AttrReverse            //nolint:revive
	AttrConceal            //nolint:revive
	AttrStrikethrough      //nolint:revive
)

// Color represents a terminal color. Negative values mean "default".
type Color int32

const DefaultColor Color = -1

const trueColorFlag Color = 1 << 24

// PaletteColor returns a palette-backed terminal color.
func PaletteColor(index int) Color {
	return Color(index)
}

// RGBColor packs a 24-bit true color value into a Color.
func RGBColor(r, g, b uint8) Color {
	return trueColorFlag | Color(r)<<16 | Color(g)<<8 | Color(b)
}

// IsRGB reports whether c stores a 24-bit true color value.
func (c Color) IsRGB() bool {
	return c >= trueColorFlag
}

// RGB returns the unpacked RGB components for a true color.
func (c Color) RGB() (uint8, uint8, uint8, bool) {
	if !c.IsRGB() {
		return 0, 0, 0, false
	}
	value := c - trueColorFlag
	return uint8(value >> 16), uint8(value >> 8), uint8(value), true
}

// Cell holds the character and style data for a single terminal cell.
type Cell struct {
	Char         rune
	FG           Color
	BG           Color
	Attrs        Attr
	Width        int
	Continuation bool
}

// Grid is a styled 2D character grid — the input to element detection.
type Grid struct {
	Rows  int
	Cols  int
	Cells [][]Cell
}

// NewGrid creates an empty grid with the given dimensions.
func NewGrid(rows, cols int) *Grid {
	cells := make([][]Cell, rows)
	for r := range cells {
		cells[r] = make([]Cell, cols)
		for c := range cells[r] {
			cells[r][c] = Cell{Char: ' ', FG: DefaultColor, BG: DefaultColor, Width: 1}
		}
	}
	return &Grid{Rows: rows, Cols: cols, Cells: cells}
}

// At returns the cell at (row, col). Returns a zero Cell if out of bounds.
func (g *Grid) At(row, col int) Cell {
	if row < 0 || row >= g.Rows || col < 0 || col >= g.Cols {
		return Cell{}
	}
	return g.Cells[row][col]
}

// Set writes a cell at (row, col). No-op if out of bounds.
func (g *Grid) Set(row, col int, c Cell) {
	if row >= 0 && row < g.Rows && col >= 0 && col < g.Cols {
		if c.Width == 0 && !c.Continuation {
			c.Width = RuneWidth(c.Char)
		}
		g.Cells[row][col] = c
	}
}

// SetText writes a plain string at (row, col) with the given attributes.
func (g *Grid) SetText(row, col int, text string, fg, bg Color, attrs Attr) {
	i := 0
	for _, ch := range text {
		width := RuneWidth(ch)
		g.Set(row, col+i, Cell{Char: ch, FG: fg, BG: bg, Attrs: attrs, Width: width})
		for offset := 1; offset < width; offset++ {
			g.Set(row, col+i+offset, Cell{FG: fg, BG: bg, Attrs: attrs, Continuation: true})
		}
		i += width
	}
}

// RuneWidth reports the number of terminal cells needed to render r.
func RuneWidth(r rune) int {
	if r == 0 {
		return 1
	}
	if r < 0x20 || (r >= 0x7f && r < 0xa0) {
		return 0
	}
	if isWideRune(r) {
		return 2
	}
	return 1
}

func isWideRune(r rune) bool {
	switch {
	case r >= 0x1100 && r <= 0x115F:
		return true
	case r >= 0x2329 && r <= 0x232A:
		return true
	case r >= 0x2E80 && r <= 0xA4CF:
		return true
	case r >= 0xAC00 && r <= 0xD7A3:
		return true
	case r >= 0xF900 && r <= 0xFAFF:
		return true
	case r >= 0xFE10 && r <= 0xFE19:
		return true
	case r >= 0xFE30 && r <= 0xFE6F:
		return true
	case r >= 0xFF00 && r <= 0xFF60:
		return true
	case r >= 0xFFE0 && r <= 0xFFE6:
		return true
	case r >= 0x1F300 && r <= 0x1FAFF:
		return true
	default:
		return false
	}
}
