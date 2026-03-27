package termtree

import (
	"fmt"
	"strings"
)

// Serialize produces a compact, token-efficient text representation of an
// ElementMap for LLM consumption.
//
// Format: [id:type:"label" row,col wxh] with * suffix on type for focused.
//
// Example:
//
//	[1:panel:"File Browser" 0,0 40x20] [2:btn:"Save"* 12,35 6x1]
func Serialize(m *ElementMap) string {
	if m == nil || len(m.Elements) == 0 {
		return "(no elements detected)"
	}

	var b strings.Builder
	for i, el := range m.Elements {
		if i > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(serializeElement(el))
	}
	return b.String()
}

func serializeElement(el Element) string {
	typeName := shortType(el.Type)
	focus := ""
	if el.Focused {
		focus = "*"
	}

	size := ""
	if el.Bounds.Width > 0 && el.Bounds.Height > 1 {
		size = fmt.Sprintf(" %dx%d", el.Bounds.Width, el.Bounds.Height)
	} else if el.Bounds.Width > 0 {
		size = fmt.Sprintf(" w%d", el.Bounds.Width)
	}

	return fmt.Sprintf("[%d:%s%s:\"%s\" %d,%d%s]",
		el.ID, typeName, focus, el.Label,
		el.Bounds.Row, el.Bounds.Col, size)
}

func shortType(t ElementType) string {
	switch t {
	case ElementPanel:
		return "panel"
	case ElementButton:
		return "btn"
	case ElementInput:
		return "input"
	case ElementMenuItem:
		return "item"
	case ElementStatusBar:
		return "status"
	case ElementMenuBar:
		return "menu"
	case ElementTab:
		return "tab"
	case ElementText:
		return "text"
	default:
		return "?"
	}
}
