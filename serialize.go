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
	ID            int    `json:"id"`
	Type          string `json:"type"`
	Label         string `json:"label"`
	Description   string `json:"description,omitempty"`
	Bounds        Rect   `json:"bounds"`
	Focused       bool   `json:"focused"`
	Enabled       bool   `json:"enabled"`
	Checked       bool   `json:"checked"`
	Selected      bool   `json:"selected"`
	Visible       bool   `json:"visible"`
	Clipped       bool   `json:"clipped"`
	Role          string `json:"role,omitempty"`
	Shortcut      string `json:"shortcut,omitempty"`
	Ref           string `json:"ref,omitempty"`
	VisibleBounds *Rect  `json:"visible_bounds,omitempty"`
	Children      []int  `json:"children,omitempty"`
}

type jsonTreeElement struct {
	ID            int               `json:"id"`
	Type          string            `json:"type"`
	Label         string            `json:"label"`
	Description   string            `json:"description,omitempty"`
	Bounds        Rect              `json:"bounds"`
	Focused       bool              `json:"focused"`
	Enabled       bool              `json:"enabled"`
	Checked       bool              `json:"checked"`
	Selected      bool              `json:"selected"`
	Visible       bool              `json:"visible"`
	Clipped       bool              `json:"clipped"`
	Role          string            `json:"role,omitempty"`
	Shortcut      string            `json:"shortcut,omitempty"`
	Ref           string            `json:"ref,omitempty"`
	VisibleBounds *Rect             `json:"visible_bounds,omitempty"`
	Children      []jsonTreeElement `json:"children,omitempty"`
}

type jsonElementMap struct {
	Elements []jsonElement     `json:"elements"`
	Tree     []jsonTreeElement `json:"tree"`
	Viewport Rect              `json:"viewport"`
	Scrolled bool              `json:"scrolled"`
	Debug    *DebugInfo        `json:"debug,omitempty"`
}

// SerializeJSON produces a structured JSON representation of an ElementMap.
func SerializeJSON(m *ElementMap) string {
	out := jsonElementMap{Elements: make([]jsonElement, 0), Tree: make([]jsonTreeElement, 0)}
	if m != nil {
		out.Viewport = m.Viewport
		out.Scrolled = m.Scrolled
		out.Debug = m.Debug
		for _, el := range m.Elements {
			out.Elements = append(out.Elements, jsonElement{
				ID:            el.ID,
				Type:          el.Type.String(),
				Label:         el.Label,
				Description:   el.Description,
				Bounds:        el.Bounds,
				Focused:       el.Focused,
				Enabled:       el.Enabled,
				Checked:       el.Checked,
				Selected:      el.Selected,
				Visible:       el.Visible,
				Clipped:       el.Clipped,
				Role:          el.Role,
				Shortcut:      el.Shortcut,
				Ref:           el.Ref,
				VisibleBounds: el.VisibleBounds,
				Children:      el.Children,
			})
		}
		out.Tree = buildJSONTree(m.Elements)
	}
	b, _ := json.Marshal(out)
	return string(b)
}

func buildJSONTree(elements []Element) []jsonTreeElement {
	if len(elements) == 0 {
		return []jsonTreeElement{}
	}
	elementsByID := make(map[int]Element, len(elements))
	children := make(map[int][]Element, len(elements))
	childIDs := make(map[int]bool, len(elements))
	for _, el := range elements {
		elementsByID[el.ID] = el
	}
	for _, el := range elements {
		for _, childID := range el.Children {
			if child, ok := elementsByID[childID]; ok {
				children[el.ID] = append(children[el.ID], child)
				childIDs[childID] = true
			}
		}
	}
	var tree []jsonTreeElement
	for _, el := range elements {
		if childIDs[el.ID] {
			continue
		}
		tree = append(tree, jsonTreeElementFromElement(el, children))
	}
	return tree
}

func jsonTreeElementFromElement(el Element, children map[int][]Element) jsonTreeElement {
	item := jsonTreeElement{
		ID:            el.ID,
		Type:          el.Type.String(),
		Label:         el.Label,
		Description:   el.Description,
		Bounds:        el.Bounds,
		Focused:       el.Focused,
		Enabled:       el.Enabled,
		Checked:       el.Checked,
		Selected:      el.Selected,
		Visible:       el.Visible,
		Clipped:       el.Clipped,
		Role:          el.Role,
		Shortcut:      el.Shortcut,
		Ref:           el.Ref,
		VisibleBounds: el.VisibleBounds,
	}
	for _, child := range children[el.ID] {
		item.Children = append(item.Children, jsonTreeElementFromElement(child, children))
	}
	return item
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
