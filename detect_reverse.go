package awtree

// detectReverseRegions finds contiguous horizontal spans of reverse-video cells.
// These typically indicate focused/selected elements in TUIs.
func detectReverseRegions(g *Grid) []Element {
	var elements []Element

	for row := 0; row < g.Rows; row++ {
		col := 0
		for col < g.Cols {
			cell := g.At(row, col)
			if cell.Attrs&AttrReverse == 0 {
				col++
				continue
			}

			// Found start of reverse region — scan to end.
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
			if len(text) == 0 {
				continue
			}

			elements = append(elements, Element{
				Type:    ElementText,
				Label:   text,
				Focused: true,
				Bounds: Rect{
					Row:    row,
					Col:    startCol,
					Width:  col - startCol,
					Height: 1,
				},
			})
		}
	}

	return elements
}

func trimSpaces(runes []rune) string {
	// Trim leading spaces.
	start := 0
	for start < len(runes) && runes[start] == ' ' {
		start++
	}
	// Trim trailing spaces.
	end := len(runes)
	for end > start && runes[end-1] == ' ' {
		end--
	}
	return string(runes[start:end])
}
