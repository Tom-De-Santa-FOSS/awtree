package awtree

// detectInputs finds input-like regions on the grid:
// 1. Underlined spans (common text input pattern)
// 2. Horizontal spans with a distinct non-default BG color surrounded by
//    default-BG cells (form fields with highlighted background)
func detectInputs(g *Grid) []Element {
	var elements []Element

	for row := 0; row < g.Rows; row++ {
		// Strategy 1: underlined spans.
		elements = append(elements, findUnderlinedSpans(g, row)...)

		// Strategy 2: distinct-BG spans not at edge rows (those are status bars).
		if row > 1 && row < g.Rows-2 {
			elements = append(elements, findDistinctBGSpans(g, row)...)
		}
	}

	return elements
}

func findUnderlinedSpans(g *Grid, row int) []Element {
	var elements []Element
	col := 0

	for col < g.Cols {
		cell := g.At(row, col)
		if cell.Attrs&AttrUnderline == 0 {
			col++
			continue
		}

		startCol := col
		var label []rune
		for col < g.Cols && g.At(row, col).Attrs&AttrUnderline != 0 {
			ch := g.At(row, col).Char
			if ch != 0 && ch != '_' {
				label = append(label, ch)
			}
			col++
		}

		width := col - startCol
		if width < 3 {
			continue
		}

		elements = append(elements, Element{
			Type:  ElementInput,
			Label: trimSpaces(label),
			Bounds: Rect{
				Row:    row,
				Col:    startCol,
				Width:  width,
				Height: 1,
			},
		})
	}

	return elements
}

func findDistinctBGSpans(g *Grid, row int) []Element {
	var elements []Element
	col := 0

	for col < g.Cols {
		cell := g.At(row, col)
		if cell.BG == DefaultColor {
			col++
			continue
		}

		// Check that this isn't reverse-video (handled by other detectors).
		if cell.Attrs&AttrReverse != 0 {
			col++
			continue
		}

		bg := cell.BG
		startCol := col
		var label []rune
		for col < g.Cols && g.At(row, col).BG == bg && g.At(row, col).Attrs&AttrReverse == 0 {
			ch := g.At(row, col).Char
			if ch != 0 {
				label = append(label, ch)
			}
			col++
		}

		width := col - startCol
		if width < 3 {
			continue
		}

		// Must have default-BG on at least one side (not a full-row bar).
		leftDefault := startCol == 0 || g.At(row, startCol-1).BG == DefaultColor
		rightDefault := col >= g.Cols || g.At(row, col).BG == DefaultColor
		if !leftDefault && !rightDefault {
			continue
		}

		// Skip if it spans most of the row (likely a bar, not an input).
		if width > g.Cols*majorityThresholdPct/100 {
			continue
		}

		elements = append(elements, Element{
			Type:  ElementInput,
			Label: trimSpaces(label),
			Bounds: Rect{
				Row:    row,
				Col:    startCol,
				Width:  width,
				Height: 1,
			},
		})
	}

	return elements
}
