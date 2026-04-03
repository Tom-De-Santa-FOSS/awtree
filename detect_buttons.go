package awtree

import (
	"unicode"
	"unicode/utf8"
)

// buttonPairs defines bracket pairs that indicate buttons.
var buttonPairs = [][2]rune{
	{'[', ']'},
	{'<', '>'},
	{'(', ')'},
}

// detectButtons finds bracketed text patterns like [Save], <OK>, (Yes).
func detectButtons(g *Grid, cfg DetectConfig, dbg *debugCollector) []Element {
	var buttons []Element

	for row := 0; row < g.Rows; row++ {
		for col := 0; col < g.Cols; col++ {
			ch := g.At(row, col).Char
			for _, pair := range buttonPairs {
				if ch != pair[0] {
					continue
				}
				if btn, ok := traceButton(g, row, col, pair[0], pair[1], cfg, dbg); ok {
					buttons = append(buttons, btn)
				}
			}
		}
	}

	return buttons
}

// traceButton attempts to find a closing bracket and extract the label.
func traceButton(g *Grid, row, col int, open, close rune, cfg DetectConfig, dbg *debugCollector) (Element, bool) {
	var label []rune
	maxWidth := cfg.MaxButtonWidth

	for c := col + 1; c < g.Cols && c < col+maxWidth; c++ {
		ch := g.At(row, c).Char
		if ch == close {
			text := string(label)
			if !isButtonLabel(text, cfg.MaxButtonLabelLen) {
				dbg.reject("buttons", "button label failed validation", Rect{Row: row, Col: col, Width: c - col + 1, Height: 1}, text)
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
			dbg.reject("buttons", "button scan hit an empty cell before closing bracket", Rect{Row: row, Col: col, Width: c - col + 1, Height: 1}, string(label))
			return Element{}, false
		}
		label = append(label, ch)
	}
	dbg.reject("buttons", "button exceeded configured scan width", Rect{Row: row, Col: col, Width: maxWidth, Height: 1}, string(label))

	return Element{}, false
}

// isButtonLabel validates that text looks like a button label:
// short, starts with a letter, mostly printable.
func isButtonLabel(text string, maxButtonLabelLen int) bool {
	if runeCount := utf8.RuneCountInString(text); runeCount == 0 || runeCount > maxButtonLabelLen {
		return false
	}

	r, _ := utf8.DecodeRuneInString(text)
	if !unicode.IsLetter(r) {
		return false
	}

	for _, r := range text {
		if !unicode.IsPrint(r) {
			return false
		}
	}

	return true
}
