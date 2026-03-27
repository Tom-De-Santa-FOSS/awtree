package awtree

import "unicode"

// buttonPairs defines bracket pairs that indicate buttons.
var buttonPairs = [][2]rune{
	{'[', ']'},
	{'<', '>'},
	{'(', ')'},
}

// detectButtons finds bracketed text patterns like [Save], <OK>, (Yes).
func detectButtons(g *Grid) []Element {
	var buttons []Element

	for row := 0; row < g.Rows; row++ {
		for col := 0; col < g.Cols; col++ {
			ch := g.At(row, col).Char
			for _, pair := range buttonPairs {
				if ch != pair[0] {
					continue
				}
				if btn, ok := traceButton(g, row, col, pair[0], pair[1]); ok {
					buttons = append(buttons, btn)
				}
			}
		}
	}

	return buttons
}

// traceButton attempts to find a closing bracket and extract the label.
func traceButton(g *Grid, row, col int, open, close rune) (Element, bool) {
	var label []rune
	maxWidth := 30 // Buttons shouldn't be wider than this.

	for c := col + 1; c < g.Cols && c < col+maxWidth; c++ {
		ch := g.At(row, c).Char
		if ch == close {
			text := string(label)
			if !isButtonLabel(text) {
				return Element{}, false
			}

			focused := g.At(row, col).Attrs&AttrReverse != 0

			return Element{
				Type:    ElementButton,
				Label:   text,
				Focused: focused,
				Bounds: Rect{
					Row:    row,
					Col:    col,
					Width:  c - col + 1,
					Height: 1,
				},
			}, true
		}
		if ch == '\n' || ch == 0 {
			break
		}
		label = append(label, ch)
	}

	return Element{}, false
}

// isButtonLabel validates that text looks like a button label:
// short, starts with a letter, mostly printable.
func isButtonLabel(text string) bool {
	if len(text) == 0 || len(text) > 20 {
		return false
	}

	runes := []rune(text)
	if !unicode.IsLetter(runes[0]) && !unicode.IsUpper(runes[0]) {
		return false
	}

	for _, r := range runes {
		if !unicode.IsPrint(r) {
			return false
		}
	}

	return true
}
