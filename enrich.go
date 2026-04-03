package awtree

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

var explicitShortcutRE = regexp.MustCompile(`(?i)\b(?:ctrl|alt)\+[a-z0-9]+\b`)
var caretShortcutRE = regexp.MustCompile(`\^([A-Za-z])\b`)
var functionShortcutRE = regexp.MustCompile(`\bF([1-9]|1[0-2])\b`)

func buildElementMap(g *Grid, elements []Element) *ElementMap {
	viewport := gridViewport(g)
	tree := BuildTree(elements)
	enriched := enrichElements(g, tree, viewport)

	return &ElementMap{
		Elements: enriched,
		Viewport: viewport,
		Scrolled: hasElementType(enriched, ElementScrollIndicator),
	}
}

func gridViewport(g *Grid) Rect {
	if g == nil || g.Rows <= 0 || g.Cols <= 0 {
		return Rect{}
	}
	return Rect{Row: 0, Col: 0, Width: g.Cols, Height: g.Rows}
}

func enrichElements(g *Grid, elements []Element, viewport Rect) []Element {
	if len(elements) == 0 {
		return nil
	}

	enriched := make([]Element, len(elements))
	copy(enriched, elements)

	idxByID := make(map[int]int, len(enriched))
	parentByID := make(map[int]int, len(enriched))
	for i, el := range enriched {
		idxByID[el.ID] = i
	}
	for _, el := range enriched {
		for _, childID := range el.Children {
			parentByID[childID] = el.ID
		}
	}

	for i := range enriched {
		enriched[i].Enabled = isElementEnabled(g, enriched[i].Bounds)
		enriched[i].Checked = inferChecked(enriched[i])
		enriched[i].Selected = (enriched[i].Type == ElementMenuItem || enriched[i].Type == ElementTab) && enriched[i].Focused
		enriched[i].Shortcut = inferShortcut(g, enriched[i])
		enriched[i].Visible, enriched[i].Clipped, enriched[i].VisibleBounds = visibilityForBounds(enriched[i].Bounds, viewport)
	}

	assignStableRefs(enriched, idxByID, parentByID)
	inferRoles(enriched, idxByID, parentByID)

	return enriched
}

func hasElementType(elements []Element, want ElementType) bool {
	for _, el := range elements {
		if el.Type == want {
			return true
		}
	}
	return false
}

func isElementEnabled(g *Grid, bounds Rect) bool {
	if g == nil {
		return true
	}
	for row := bounds.Row; row < bounds.Row+bounds.Height; row++ {
		for col := bounds.Col; col < bounds.Col+bounds.Width; col++ {
			if g.At(row, col).Attrs&AttrFaint != 0 {
				return false
			}
		}
	}
	return true
}

func inferChecked(el Element) bool {
	if el.Type != ElementCheckbox {
		return false
	}

	label := strings.TrimSpace(el.Label)
	switch {
	case strings.HasPrefix(label, "☑"), strings.HasPrefix(label, "☒"), strings.HasPrefix(label, "✓"), strings.HasPrefix(label, "✗"):
		return true
	case strings.HasPrefix(label, "[x]"), strings.HasPrefix(label, "[X]"), strings.HasPrefix(label, "[*]"):
		return true
	case strings.HasPrefix(label, "(x)"), strings.HasPrefix(label, "(X)"), strings.HasPrefix(label, "(*)"):
		return true
	default:
		return false
	}
}

func inferShortcut(g *Grid, el Element) string {
	if supportsUnderlineShortcut(el.Type) {
		if shortcut := shortcutFromUnderline(g, el.Bounds); shortcut != "" {
			return shortcut
		}
	}
	return shortcutFromLabel(el.Label)
}

func supportsUnderlineShortcut(t ElementType) bool {
	switch t {
	case ElementButton, ElementMenuItem, ElementTab:
		return true
	default:
		return false
	}
}

func shortcutFromUnderline(g *Grid, bounds Rect) string {
	if g == nil {
		return ""
	}
	for row := bounds.Row; row < bounds.Row+bounds.Height; row++ {
		for col := bounds.Col; col < bounds.Col+bounds.Width; col++ {
			cell := g.At(row, col)
			if cell.Attrs&AttrUnderline == 0 {
				continue
			}
			if !unicode.IsLetter(cell.Char) && !unicode.IsDigit(cell.Char) {
				continue
			}
			return "Alt+" + strings.ToUpper(string(cell.Char))
		}
	}
	return ""
}

func shortcutFromLabel(label string) string {
	if match := explicitShortcutRE.FindString(label); match != "" {
		parts := strings.SplitN(strings.ToLower(match), "+", 2)
		if len(parts) == 2 {
			return strings.Title(parts[0]) + "+" + strings.ToUpper(parts[1])
		}
	}
	if match := caretShortcutRE.FindStringSubmatch(label); len(match) == 2 {
		return "Ctrl+" + strings.ToUpper(match[1])
	}
	if match := functionShortcutRE.FindString(label); match != "" {
		return strings.ToUpper(match)
	}
	return ""
}

func visibilityForBounds(bounds, viewport Rect) (bool, bool, *Rect) {
	visibleBounds, ok := intersectRect(bounds, viewport)
	if !ok {
		return false, false, nil
	}
	clipped := !sameRect(bounds, visibleBounds)
	return true, clipped, &visibleBounds
}

func intersectRect(a, b Rect) (Rect, bool) {
	row := maxInt(a.Row, b.Row)
	col := maxInt(a.Col, b.Col)
	bottom := minInt(a.Row+a.Height, b.Row+b.Height)
	right := minInt(a.Col+a.Width, b.Col+b.Width)
	if bottom <= row || right <= col {
		return Rect{}, false
	}
	return Rect{Row: row, Col: col, Width: right - col, Height: bottom - row}, true
}

func assignStableRefs(elements []Element, idxByID map[int]int, parentByID map[int]int) {
	childrenByParent := make(map[int][]int)
	var roots []int
	for _, el := range elements {
		parentID := parentByID[el.ID]
		if parentID == 0 {
			roots = append(roots, el.ID)
			continue
		}
		childrenByParent[parentID] = append(childrenByParent[parentID], el.ID)
	}

	var assign func(parentID int, prefix string, ids []int)
	assign = func(parentID int, prefix string, ids []int) {
		counts := make(map[ElementType]int)
		for _, id := range ids {
			idx := idxByID[id]
			counts[elements[idx].Type]++
			segment := fmt.Sprintf("%s[%d]", elements[idx].Type.String(), counts[elements[idx].Type])
			if prefix == "" {
				elements[idx].Ref = segment
			} else {
				elements[idx].Ref = prefix + "/" + segment
			}
			assign(id, elements[idx].Ref, childrenByParent[id])
		}
	}

	assign(0, "", roots)
}

func inferRoles(elements []Element, idxByID map[int]int, parentByID map[int]int) {
	for i := range elements {
		switch elements[i].Type {
		case ElementDialog:
			elements[i].Role = "dialog"
		case ElementMenuBar:
			elements[i].Role = "menubar"
		case ElementTab:
			elements[i].Role = "tab"
		case ElementStatusBar:
			elements[i].Role = "status"
		case ElementButton:
			elements[i].Role = "button"
		case ElementInput:
			elements[i].Role = "textbox"
		case ElementCheckbox:
			elements[i].Role = "checkbox"
		case ElementProgressBar:
			elements[i].Role = "progressbar"
		case ElementTable:
			elements[i].Role = "table"
		case ElementSeparator:
			elements[i].Role = "separator"
		case ElementBreadcrumb:
			elements[i].Role = "navigation"
		}
	}

	for i := range elements {
		if elements[i].Type != ElementPanel {
			continue
		}
		childTypes := make(map[ElementType]int)
		for _, childID := range elements[i].Children {
			childTypes[elements[idxByID[childID]].Type]++
		}
		switch {
		case childTypes[ElementInput] > 0 || childTypes[ElementCheckbox] > 0:
			elements[i].Role = "form"
		case childTypes[ElementMenuItem] > 0:
			elements[i].Role = "listbox"
		case childTypes[ElementTab] > 1:
			elements[i].Role = "tablist"
		case elements[i].Bounds.Row <= 1 && childTypes[ElementMenuItem] > 0:
			elements[i].Role = "navigation"
		}
	}

	for i := range elements {
		if elements[i].Type != ElementMenuItem {
			continue
		}
		parentID := parentByID[elements[i].ID]
		if parentID == 0 {
			continue
		}
		if elements[idxByID[parentID]].Role == "listbox" {
			elements[i].Role = "option"
		}
	}
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
