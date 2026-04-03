package awtree

import (
	"encoding/json"
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

// jsonElement is the JSON-friendly representation of Element.
type jsonElement struct {
	ID       int          `json:"id"`
	Type     string       `json:"type"`
	Label    string       `json:"label"`
	Bounds   Rect         `json:"bounds"`
	Focused  bool         `json:"focused"`
	Children []int        `json:"children,omitempty"`
}

type jsonElementMap struct {
	Elements []jsonElement `json:"elements"`
}

// SerializeJSON produces a structured JSON representation of an ElementMap.
func SerializeJSON(m *ElementMap) string {
	out := jsonElementMap{Elements: make([]jsonElement, 0)}
	if m != nil {
		for _, el := range m.Elements {
			out.Elements = append(out.Elements, jsonElement{
				ID:       el.ID,
				Type:     el.Type.String(),
				Label:    el.Label,
				Bounds:   el.Bounds,
				Focused:  el.Focused,
				Children: el.Children,
			})
		}
	}
	b, _ := json.Marshal(out)
	return string(b)
}

func serializeElement(el Element) string {
	typeName := el.Type.ShortString()
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
