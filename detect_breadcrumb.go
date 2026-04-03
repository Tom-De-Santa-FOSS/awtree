package awtree

import (
	"strings"
	"unicode/utf8"
)

// breadcrumbSeparators are the separator patterns that delimit breadcrumb
// segments. Each must be space-padded (e.g. " > " not ">").
var breadcrumbSeparators = []string{" > ", " / ", " » ", " › ", " → "}

// detectBreadcrumbs finds path-like breadcrumb trails such as
// "Home > Settings > Display" or "File / Edit / View".
func detectBreadcrumbs(g *Grid) []Element {
	var elements []Element
	maxWidth := g.Cols * 80 / 100

	for row := 0; row < g.Rows; row++ {
		// Extract the full row text and find the non-space bounds.
		text, startCol := extractRowText(g, row)
		if utf8.RuneCountInString(text) == 0 {
			continue
		}

		// Try each separator type.
		for _, sep := range breadcrumbSeparators {
			parts := strings.Split(text, sep)
			if len(parts) < 3 {
				continue
			}

			// Verify all segments are non-empty after trimming.
			valid := true
			for _, p := range parts {
				if strings.TrimSpace(p) == "" {
					valid = false
					break
				}
			}
			if !valid {
				continue
			}

			if utf8.RuneCountInString(text) > maxWidth {
				continue
			}

			elements = append(elements, Element{
				Type:  ElementBreadcrumb,
				Label: text,
				Bounds: Rect{
					Row:    row,
					Col:    startCol,
					Width:  utf8.RuneCountInString(text),
					Height: 1,
				},
			})
			break // One breadcrumb per row, first matching separator wins.
		}
	}

	return elements
}

// extractRowText returns the trimmed text content of a row and the column
// where the first non-space character appears.
func extractRowText(g *Grid, row int) (string, int) {
	var runes []rune
	for col := 0; col < g.Cols; col++ {
		ch := g.At(row, col).Char
		if ch == 0 {
			ch = ' '
		}
		runes = append(runes, ch)
	}

	// Find first and last non-space positions.
	first := -1
	last := -1
	for i, r := range runes {
		if r != ' ' {
			if first == -1 {
				first = i
			}
			last = i
		}
	}

	if first == -1 {
		return "", 0
	}

	return string(runes[first : last+1]), first
}
