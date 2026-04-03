package awtree

import "unicode"

// ASCII fill characters used inside bracketed progress bars.
func isASCIIFill(ch rune) bool {
	return ch == '=' || ch == '#'
}

// ASCII empty characters used inside bracketed progress bars.
func isASCIIEmpty(ch rune) bool {
	return ch == ' ' || ch == '-'
}

// findPercentLabel scans the row for a "N%" pattern near the bar.
// It looks within 5 columns before and after the bar bounds.
func findPercentLabel(g *Grid, row, barCol, barWidth int) string {
	searchStart := barCol - 5
	if searchStart < 0 {
		searchStart = 0
	}
	searchEnd := barCol + barWidth + 5
	if searchEnd > g.Cols {
		searchEnd = g.Cols
	}

	for c := searchStart; c < searchEnd; c++ {
		ch := g.At(row, c).Char
		if unicode.IsDigit(ch) {
			// Collect digits.
			numStart := c
			for c < searchEnd && unicode.IsDigit(g.At(row, c).Char) {
				c++
			}
			// Check for '%' immediately after.
			if c < searchEnd && g.At(row, c).Char == '%' {
				c++
				// Build the label string.
				var label []rune
				for i := numStart; i < c; i++ {
					label = append(label, g.At(row, i).Char)
				}
				return string(label)
			}
		}
	}
	return ""
}

// detectProgressBars scans the grid for progress bar patterns.
func detectProgressBars(g *Grid) []Element {
	var results []Element

	for r := 0; r < g.Rows; r++ {
		c := 0
		for c < g.Cols {
			ch := g.At(r, c).Char
			if ch == '█' {
				// Block char pattern: ████░░░░
				fillStart := c
				for c < g.Cols && g.At(r, c).Char == '█' {
					c++
				}
				fillCount := c - fillStart

				for c < g.Cols && g.At(r, c).Char == '░' {
					c++
				}
				emptyCount := c - fillStart - fillCount

				if fillCount >= 3 && emptyCount >= 3 {
					barWidth := fillCount + emptyCount
					label := findPercentLabel(g, r, fillStart, barWidth)
					results = append(results, Element{
						Type:   ElementProgressBar,
						Label:  label,
						Bounds: Rect{Row: r, Col: fillStart, Width: barWidth, Height: 1},
					})
				}
			} else if ch == '[' {
				// ASCII pattern: [====    ] or [###---]
				barStart := c
				c++ // skip '['

				// Count fill chars.
				fillCount := 0
				for c < g.Cols && isASCIIFill(g.At(r, c).Char) {
					fillCount++
					c++
				}

				// Count empty chars.
				emptyCount := 0
				for c < g.Cols && isASCIIEmpty(g.At(r, c).Char) {
					emptyCount++
					c++
				}

				// Expect closing bracket.
				if c < g.Cols && g.At(r, c).Char == ']' && fillCount >= 1 && emptyCount >= 1 {
					c++ // skip ']'
					barWidth := c - barStart
					label := findPercentLabel(g, r, barStart, barWidth)
					results = append(results, Element{
						Type:   ElementProgressBar,
						Label:  label,
						Bounds: Rect{Row: r, Col: barStart, Width: barWidth, Height: 1},
					})
				}
			} else {
				c++
			}
		}
	}

	return results
}
