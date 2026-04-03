package awtree

// separatorChars are horizontal line characters that form standalone separators.
var separatorChars = map[rune]bool{
	'─': true, '━': true, '═': true, '-': true,
}

// borderCornerChars are characters that indicate a panel border row.
var borderCornerChars = map[rune]bool{
	'┌': true, '┐': true, '└': true, '┘': true,
	'┬': true, '┴': true, '├': true, '┤': true,
	'╔': true, '╗': true, '╚': true, '╝': true,
	'╭': true, '╮': true, '╰': true, '╯': true,
	'┏': true, '┓': true, '┗': true, '┛': true,
	'+': true,
}

// detectSeparators finds horizontal separator/divider lines on the grid.
func detectSeparators(g *Grid) []Element {
	var result []Element

	for r := 0; r < g.Rows; r++ {
		// Skip rows that look like panel borders.
		if isPanelBorderRow(g, r) {
			continue
		}

		c := 0
		for c < g.Cols {
			ch := g.At(r, c).Char
			if !separatorChars[ch] {
				c++
				continue
			}
			// Found the start of a potential separator run.
			start := c
			for c < g.Cols && separatorChars[g.At(r, c).Char] {
				c++
			}
			width := c - start
			if width >= 3 {
				result = append(result, Element{
					Type:   ElementSeparator,
					Label:  "",
					Bounds: Rect{Row: r, Col: start, Width: width, Height: 1},
				})
			}
		}
	}
	return result
}

// isPanelBorderRow checks if a row starts or ends with border corner characters.
func isPanelBorderRow(g *Grid, r int) bool {
	// Find first non-space char on the row.
	first := rune(0)
	last := rune(0)
	for c := 0; c < g.Cols; c++ {
		ch := g.At(r, c).Char
		if ch != ' ' && ch != 0 {
			if first == 0 {
				first = ch
			}
			last = ch
		}
	}
	return borderCornerChars[first] || borderCornerChars[last]
}
