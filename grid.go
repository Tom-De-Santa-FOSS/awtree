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

// Cell holds the character and style data for a single terminal cell.
type Cell struct {
	Char  rune
	FG    Color
	BG    Color
	Attrs Attr
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
			cells[r][c] = Cell{Char: ' ', FG: DefaultColor, BG: DefaultColor}
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
		g.Cells[row][col] = c
	}
}

// SetText writes a plain string at (row, col) with the given attributes.
func (g *Grid) SetText(row, col int, text string, fg, bg Color, attrs Attr) {
	i := 0
	for _, ch := range text {
		g.Set(row, col+i, Cell{Char: ch, FG: fg, BG: bg, Attrs: attrs})
		i++
	}
}
