package awtree

// boxDrawingCornerTL maps top-left corner characters.
var boxDrawingCornerTL = map[rune]bool{
	'┌': true, '╔': true, '╒': true, '╓': true,
	'┏': true, '┍': true, '┎': true,
	'+': true, // ASCII fallback
}

// boxDrawingCornerTR maps top-right corner characters.
var boxDrawingCornerTR = map[rune]bool{
	'┐': true, '╗': true, '╕': true, '╖': true,
	'┓': true, '┑': true, '┒': true,
	'+': true,
}

// boxDrawingCornerBL maps bottom-left corner characters.
var boxDrawingCornerBL = map[rune]bool{
	'└': true, '╚': true, '╘': true, '╙': true,
	'┗': true, '┕': true, '┖': true,
	'+': true,
}

// boxDrawingCornerBR maps bottom-right corner characters.
var boxDrawingCornerBR = map[rune]bool{
	'┘': true, '╝': true, '╛': true, '╜': true,
	'┛': true, '┙': true, '┚': true,
	'+': true,
}

// boxDrawingHorizontal maps horizontal line characters.
var boxDrawingHorizontal = map[rune]bool{
	'─': true, '═': true, '━': true, '-': true,
}

// boxDrawingVertical maps vertical line characters.
var boxDrawingVertical = map[rune]bool{
	'│': true, '║': true, '┃': true, '|': true,
}

// detectPanels finds rectangular regions bounded by box-drawing characters.
func detectPanels(g *Grid) []Element {
	var panels []Element

	for row := 0; row < g.Rows-1; row++ {
		for col := 0; col < g.Cols-1; col++ {
			ch := g.At(row, col).Char
			if !boxDrawingCornerTL[ch] {
				continue
			}

			if panel, ok := tracePanel(g, row, col); ok {
				// Extract title from top border if present.
				panel.Label = extractPanelTitle(g, row, col, panel.Bounds.Width)
				panels = append(panels, panel)
			}
		}
	}

	return panels
}

// tracePanel attempts to trace a complete box starting from a top-left corner.
func tracePanel(g *Grid, startRow, startCol int) (Element, bool) {
	// Find top-right corner by scanning horizontally.
	endCol := -1
	for c := startCol + 1; c < g.Cols; c++ {
		ch := g.At(startRow, c).Char
		if boxDrawingCornerTR[ch] {
			endCol = c
			break
		}
		if !boxDrawingHorizontal[ch] && ch != ' ' && !isBoxTitle(ch) {
			break
		}
	}
	if endCol == -1 {
		return Element{}, false
	}

	// Find bottom-left corner by scanning vertically.
	endRow := -1
	for r := startRow + 1; r < g.Rows; r++ {
		ch := g.At(r, startCol).Char
		if boxDrawingCornerBL[ch] {
			endRow = r
			break
		}
		if !boxDrawingVertical[ch] {
			break
		}
	}
	if endRow == -1 {
		return Element{}, false
	}

	// Verify bottom-right corner exists.
	if !boxDrawingCornerBR[g.At(endRow, endCol).Char] {
		return Element{}, false
	}

	return Element{
		Type: ElementPanel,
		Bounds: Rect{
			Row:    startRow,
			Col:    startCol,
			Width:  endCol - startCol + 1,
			Height: endRow - startRow + 1,
		},
	}, true
}

// isBoxTitle returns true if the character could be part of a panel title in the border.
func isBoxTitle(ch rune) bool {
	return ch != 0 && ch != '\n' && ch != '\r'
}

// extractPanelTitle extracts text embedded in the top border of a panel.
func extractPanelTitle(g *Grid, row, col, width int) string {
	var title []rune
	inTitle := false

	for c := col + 1; c < col+width-1; c++ {
		ch := g.At(row, c).Char
		if boxDrawingHorizontal[ch] {
			if inTitle {
				break
			}
			continue
		}
		if ch == ' ' && !inTitle {
			inTitle = true
			continue
		}
		if ch == ' ' && inTitle {
			// Check if there's more title text ahead.
			next := g.At(row, c+1).Char
			if boxDrawingHorizontal[next] || boxDrawingCornerTR[next] {
				break
			}
		}
		inTitle = true
		title = append(title, ch)
	}

	return string(title)
}
