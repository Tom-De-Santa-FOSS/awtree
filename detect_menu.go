package awtree

// detectMenuItems finds vertical lists of text lines where at least one is
// reverse-video (highlighted). Non-highlighted siblings at the same column
// range are also included as menu items.
func detectMenuItems(g *Grid) []Element {
	// Find all reverse-video regions first.
	type region struct {
		row, col, width int
		label           string
	}
	var highlighted []region

	for row := 0; row < g.Rows; row++ {
		col := 0
		for col < g.Cols {
			cell := g.At(row, col)
			if cell.Attrs&AttrReverse == 0 {
				col++
				continue
			}
			startCol := col
			var label []rune
			for col < g.Cols && g.At(row, col).Attrs&AttrReverse != 0 {
				ch := g.At(row, col).Char
				if ch != 0 {
					label = append(label, ch)
				}
				col++
			}
			text := trimSpaces(label)
			if len(text) > 0 {
				highlighted = append(highlighted, region{
					row: row, col: startCol, width: col - startCol, label: text,
				})
			}
		}
	}

	var elements []Element

	for _, h := range highlighted {
		// Look for sibling lines above and below with text at the same column range.
		siblings := []Element{{
			Type:    ElementMenuItem,
			Label:   h.label,
			Focused: true,
			Bounds:  Rect{Row: h.row, Col: h.col, Width: h.width, Height: 1},
		}}

		// Scan upward.
		for r := h.row - 1; r >= 0; r-- {
			label := extractLineText(g, r, h.col, h.width)
			if label == "" {
				break
			}
			siblings = append(siblings, Element{
				Type:   ElementMenuItem,
				Label:  label,
				Bounds: Rect{Row: r, Col: h.col, Width: h.width, Height: 1},
			})
		}

		// Scan downward.
		for r := h.row + 1; r < g.Rows; r++ {
			label := extractLineText(g, r, h.col, h.width)
			if label == "" {
				break
			}
			siblings = append(siblings, Element{
				Type:   ElementMenuItem,
				Label:  label,
				Bounds: Rect{Row: r, Col: h.col, Width: h.width, Height: 1},
			})
		}

		// Only emit as menu items if there are siblings (>1 item).
		if len(siblings) > 1 {
			elements = append(elements, siblings...)
		}
	}

	return elements
}

// extractLineText gets trimmed text from a row at the given column range.
// Returns empty string if the region is blank.
func extractLineText(g *Grid, row, col, width int) string {
	var label []rune
	for c := col; c < col+width && c < g.Cols; c++ {
		ch := g.At(row, c).Char
		if ch != 0 {
			label = append(label, ch)
		}
	}
	return trimSpaces(label)
}
