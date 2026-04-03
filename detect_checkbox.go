package awtree

// unicodeCheckboxChars maps Unicode checkbox/radio characters to true.
var unicodeCheckboxChars = map[rune]bool{
	'☐': true,
	'☑': true,
	'☒': true,
	'✓': true,
	'✗': true,
}

// asciiBracketPairs maps opening brackets to their expected closing bracket.
var asciiBracketPairs = map[rune]rune{
	'[': ']',
	'(': ')',
}

// asciiCheckboxInners lists valid characters inside ASCII checkbox brackets.
var asciiCheckboxInners = map[rune]bool{
	'x': true,
	'X': true,
	'*': true,
	' ': true,
}

// detectCheckboxes finds checkbox and radio button patterns on the grid.
func detectCheckboxes(g *Grid) []Element {
	var results []Element

	for row := 0; row < g.Rows; row++ {
		col := 0
		for col < g.Cols {
			ch := g.At(row, col).Char

			if unicodeCheckboxChars[ch] {
				el := buildCheckboxFromCol(g, row, col, 1)
				results = append(results, el)
				col += el.Bounds.Width
				continue
			}

			// ASCII checkbox [x]/[ ] or radio (x)/( ).
			if closer, ok := asciiBracketPairs[ch]; ok && col+2 < g.Cols {
				inner := g.At(row, col+1).Char
				end := g.At(row, col+2).Char
				if asciiCheckboxInners[inner] && end == closer {
					el := buildCheckboxFromCol(g, row, col, 3)
					results = append(results, el)
					col += el.Bounds.Width
					continue
				}
			}

			col++
		}
	}

	return results
}

// buildCheckboxFromCol constructs a checkbox Element starting at (row, col)
// where indicatorWidth is the width of the indicator itself (1 for unicode,
// 3 for ASCII like [x]).
func buildCheckboxFromCol(g *Grid, row, col, indicatorWidth int) Element {
	// Collect the indicator characters.
	var label []rune
	for c := col; c < col+indicatorWidth && c < g.Cols; c++ {
		label = append(label, g.At(row, c).Char)
	}

	// Collect trailing label text on the same row.
	for c := col + indicatorWidth; c < g.Cols; c++ {
		ch := g.At(row, c).Char
		if ch == 0 || ch == '\n' {
			break
		}
		if ch == ' ' {
			// Check if rest of line is blank.
			allBlank := true
			for cc := c; cc < g.Cols; cc++ {
				if g.At(row, cc).Char != ' ' && g.At(row, cc).Char != 0 {
					allBlank = false
					break
				}
			}
			if allBlank {
				break
			}
		}
		label = append(label, ch)
	}

	// Trim trailing spaces.
	for len(label) > 0 && label[len(label)-1] == ' ' {
		label = label[:len(label)-1]
	}

	text := string(label)
	focused := g.At(row, col).Attrs&AttrReverse != 0

	return Element{
		Type:    ElementCheckbox,
		Label:   text,
		Focused: focused,
		Bounds: Rect{
			Row:    row,
			Col:    col,
			Width:  len(label),
			Height: 1,
		},
	}
}
