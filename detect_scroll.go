package awtree

// detectScrollIndicators finds vertical scroll indicators in the grid.
// It detects two patterns:
//   - Arrow pairs: ▲ at top and ▼ at bottom of the same column, with 2+ rows between them
//   - Block scrollbar: vertical run of 3+ █, ▓, or ┃ chars in a single column
func detectScrollIndicators(g *Grid) []Element {
	var results []Element

	for col := 0; col < g.Cols; col++ {
		// Look for arrow pairs: ▲ then ▼ below in same column.
		for row := 0; row < g.Rows; row++ {
			if g.At(row, col).Char == '▲' {
				// Scan downward for ▼ with at least 2 rows between.
				for r2 := row + 3; r2 < g.Rows; r2++ {
					if g.At(r2, col).Char == '▼' {
						results = append(results, Element{
							Type:  ElementScrollIndicator,
							Bounds: Rect{
								Row:    row,
								Col:    col,
								Width:  1,
								Height: r2 - row + 1,
							},
						})
						break
					}
				}
			}
		}

		// Look for block scrollbar: vertical run of 3+ block chars.
		row := 0
		for row < g.Rows {
			ch := g.At(row, col).Char
			if ch == '█' || ch == '▓' || ch == '┃' {
				start := row
				for row < g.Rows {
					c := g.At(row, col).Char
					if c != '█' && c != '▓' && c != '┃' {
						break
					}
					row++
				}
				length := row - start
				if length >= 3 {
					results = append(results, Element{
						Type:  ElementScrollIndicator,
						Bounds: Rect{
							Row:    start,
							Col:    col,
							Width:  1,
							Height: length,
						},
					})
				}
			} else {
				row++
			}
		}
	}

	return results
}
