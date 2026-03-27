package awtree

// detectTabs finds horizontal sequences of text segments on the same row
// where at least one segment is styled differently (reverse/bold), indicating
// a tab bar with an active tab.
func detectTabs(g *Grid) []Element {
	var elements []Element

	for row := 0; row < g.Rows; row++ {
		segments := findTextSegments(g, row)
		if len(segments) < 2 {
			continue
		}

		// Count how many segments are "active" (reverse or bold+colored).
		activeCount := 0
		for _, s := range segments {
			if s.active {
				activeCount++
			}
		}

		// Tab bar pattern: exactly 1 active among 2+ segments.
		if activeCount != 1 {
			continue
		}

		for _, s := range segments {
			elements = append(elements, Element{
				Type:    ElementTab,
				Label:   s.label,
				Focused: s.active,
				Bounds: Rect{
					Row:    row,
					Col:    s.col,
					Width:  s.width,
					Height: 1,
				},
			})
		}
	}

	return elements
}

type textSegment struct {
	col    int
	width  int
	label  string
	active bool
}

// findTextSegments splits a row into contiguous non-empty text segments,
// tracking whether each segment has "active" styling (reverse video).
func findTextSegments(g *Grid, row int) []textSegment {
	var segments []textSegment

	col := 0
	for col < g.Cols {
		cell := g.At(row, col)

		// Skip empty cells between segments.
		if cell.Char == ' ' && cell.Attrs == 0 && cell.BG == DefaultColor {
			col++
			continue
		}

		// Start of a segment — scan to find its extent.
		startCol := col
		hasReverse := false
		var label []rune

		for col < g.Cols {
			c := g.At(row, col)

			// Segment ends at unstyled space (gap between tabs).
			if c.Char == ' ' && c.Attrs == 0 && c.BG == DefaultColor {
				break
			}

			if c.Attrs&AttrReverse != 0 {
				hasReverse = true
			}
			if c.Char != 0 {
				label = append(label, c.Char)
			}
			col++
		}

		text := trimSpaces(label)
		if len(text) > 0 {
			segments = append(segments, textSegment{
				col:    startCol,
				width:  col - startCol,
				label:  text,
				active: hasReverse,
			})
		}
	}

	return segments
}
