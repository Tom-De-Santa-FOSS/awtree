package awtree

// detectStatusBars finds rows at the top or bottom of the screen where a
// majority of cells share a distinct background color, indicating a
// status bar or menu bar.
func detectStatusBars(g *Grid) []Element {
	var elements []Element

	// Check first and last rows.
	edges := []int{0, g.Rows - 1}
	if g.Rows > 2 {
		edges = append(edges, 1, g.Rows-2) // Also check second/penultimate rows.
	}

	seen := make(map[int]bool)
	for _, row := range edges {
		if row < 0 || row >= g.Rows || seen[row] {
			continue
		}
		seen[row] = true

		if bar, ok := detectBarRow(g, row); ok {
			elements = append(elements, bar)
		}
	}

	return elements
}

// detectBarRow checks if a row looks like a status/menu bar:
// majority of cells share a non-default background color.
func detectBarRow(g *Grid, row int) (Element, bool) {
	bgCounts := make(map[Color]int)
	nonEmpty := 0

	for col := 0; col < g.Cols; col++ {
		cell := g.At(row, col)
		if cell.Char == 0 || cell.Char == ' ' {
			bgCounts[cell.BG]++
			continue
		}
		nonEmpty++
		bgCounts[cell.BG]++
	}

	// Find the dominant background color.
	var dominantBG Color
	maxCount := 0
	for bg, count := range bgCounts {
		if count > maxCount {
			maxCount = count
			dominantBG = bg
		}
	}

	// Must cover >60% of cells and not be the default color.
	threshold := g.Cols * 60 / 100
	if maxCount < threshold || dominantBG == DefaultColor {
		return Element{}, false
	}

	// Extract visible text.
	var label []rune
	for col := 0; col < g.Cols; col++ {
		cell := g.At(row, col)
		if cell.Char != 0 {
			label = append(label, cell.Char)
		}
	}

	typ := ElementStatusBar
	if row == 0 || row == 1 {
		typ = ElementMenuBar
	}

	return Element{
		Type:  typ,
		Label: trimSpaces(label),
		Bounds: Rect{
			Row:    row,
			Col:    0,
			Width:  g.Cols,
			Height: 1,
		},
	}, true
}
