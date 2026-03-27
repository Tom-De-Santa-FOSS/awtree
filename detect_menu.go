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

	// Group highlighted regions by (col, width) to avoid scanning siblings
	// multiple times when several highlights share the same column range.
	type colKey struct{ col, width int }
	groups := make(map[colKey][]region)
	var groupOrder []colKey
	for _, h := range highlighted {
		k := colKey{h.col, h.width}
		if _, exists := groups[k]; !exists {
			groupOrder = append(groupOrder, k)
		}
		groups[k] = append(groups[k], h)
	}

	var elements []Element

	for _, k := range groupOrder {
		group := groups[k]

		// Find the extent of all highlights in this column range.
		minRow, maxRow := group[0].row, group[0].row
		for _, h := range group[1:] {
			if h.row < minRow {
				minRow = h.row
			}
			if h.row > maxRow {
				maxRow = h.row
			}
		}

		// Build a set of highlighted rows for focused marking.
		focusedRows := make([]bool, g.Rows)
		focusedLabels := make([]string, g.Rows)
		for _, h := range group {
			focusedRows[h.row] = true
			focusedLabels[h.row] = h.label
		}

		// Scan upward from the topmost highlight.
		topRow := minRow
		for r := minRow - 1; r >= 0; r-- {
			label := extractLineText(g, r, k.col, k.width)
			if label == "" {
				break
			}
			topRow = r
		}

		// Scan downward from the bottommost highlight.
		bottomRow := maxRow
		for r := maxRow + 1; r < g.Rows; r++ {
			label := extractLineText(g, r, k.col, k.width)
			if label == "" {
				break
			}
			bottomRow = r
		}

		// Only emit as menu items if there are multiple items.
		count := bottomRow - topRow + 1
		if count < 2 {
			continue
		}

		for r := topRow; r <= bottomRow; r++ {
			label := ""
			if focusedRows[r] {
				label = focusedLabels[r]
			} else {
				label = extractLineText(g, r, k.col, k.width)
			}
			elements = append(elements, Element{
				Type:    ElementMenuItem,
				Label:   label,
				Focused: focusedRows[r],
				Bounds:  Rect{Row: r, Col: k.col, Width: k.width, Height: 1},
			})
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
