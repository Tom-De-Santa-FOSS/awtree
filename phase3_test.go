package awtree

import (
	"encoding/json"
	"testing"
)

func TestDetect_Phase3InteractiveStateAndShortcut(t *testing.T) {
	g := NewGrid(8, 40)
	g.Set(1, 1, Cell{Char: '☑'})
	g.SetText(1, 3, "Enable feature", DefaultColor, DefaultColor, 0)
	g.SetText(2, 2, " File ", DefaultColor, DefaultColor, AttrReverse)
	g.Set(2, 3, Cell{Char: 'F', FG: DefaultColor, BG: DefaultColor, Attrs: AttrReverse | AttrUnderline})
	g.SetText(3, 2, "[Save (Ctrl+S)]", DefaultColor, DefaultColor, 0)
	g.SetText(4, 2, "[Disabled]", DefaultColor, DefaultColor, AttrFaint)

	m := Detect(g)

	var checked, selected, explicitShortcut, disabled *Element
	for i := range m.Elements {
		el := &m.Elements[i]
		switch {
		case el.Type == ElementCheckbox:
			checked = el
		case el.Type == ElementMenuItem && el.Label == "File":
			selected = el
		case el.Type == ElementButton && el.Label == "Save (Ctrl+S)":
			explicitShortcut = el
		case el.Type == ElementButton && el.Label == "Disabled":
			disabled = el
		}
	}

	if checked == nil || !checked.Checked {
		t.Fatal("expected checked checkbox")
	}
	if selected == nil || !selected.Selected || selected.Shortcut != "Alt+F" {
		t.Fatalf("menu item = %+v, want selected with Alt+F shortcut", selected)
	}
	if explicitShortcut == nil || explicitShortcut.Shortcut != "Ctrl+S" {
		t.Fatalf("button shortcut = %+v, want Ctrl+S", explicitShortcut)
	}
	if disabled == nil || disabled.Enabled {
		t.Fatalf("disabled button = %+v, want Enabled=false", disabled)
	}
}

func TestDetect_Phase3RolesRefsAndViewport(t *testing.T) {
	g := NewGrid(10, 40)
	// Panel with input and checkbox should become a form.
	g.Set(0, 0, Cell{Char: '┌'})
	for c := 1; c < 19; c++ {
		g.Set(0, c, Cell{Char: '─'})
		g.Set(5, c, Cell{Char: '─'})
	}
	g.Set(0, 19, Cell{Char: '┐'})
	g.Set(5, 0, Cell{Char: '└'})
	g.Set(5, 19, Cell{Char: '┘'})
	for r := 1; r < 5; r++ {
		g.Set(r, 0, Cell{Char: '│'})
		g.Set(r, 19, Cell{Char: '│'})
	}
	g.SetText(2, 2, "[x] Email me", DefaultColor, DefaultColor, 0)
	g.SetText(3, 2, "_______", DefaultColor, DefaultColor, AttrUnderline)
	// Scroll indicator should mark map as scrolled.
	g.Set(1, 30, Cell{Char: '▲'})
	g.Set(6, 30, Cell{Char: '▼'})

	m := Detect(g)
	if m.Viewport.Width != 40 || m.Viewport.Height != 10 {
		t.Fatalf("viewport = %+v, want 40x10", m.Viewport)
	}
	if !m.Scrolled {
		t.Fatal("expected scrolled=true when scroll indicator exists")
	}

	var panel, checkbox, input *Element
	for i := range m.Elements {
		el := &m.Elements[i]
		switch el.Type {
		case ElementPanel:
			panel = el
		case ElementCheckbox:
			checkbox = el
		case ElementInput:
			input = el
		}
	}

	if panel == nil || panel.Role != "form" || panel.Ref == "" {
		t.Fatalf("panel = %+v, want form role and stable ref", panel)
	}
	if checkbox == nil || checkbox.Role != "checkbox" || !checkbox.Visible || checkbox.VisibleBounds == nil {
		t.Fatalf("checkbox = %+v, want checkbox role and visible bounds", checkbox)
	}
	if input == nil || input.Role != "textbox" {
		t.Fatalf("input = %+v, want textbox role", input)
	}
	if input.Shortcut != "" {
		t.Fatalf("input shortcut = %q, want empty", input.Shortcut)
	}
}

func TestElementMapQuery_Phase3Selectors(t *testing.T) {
	m := &ElementMap{Elements: BuildTree([]Element{
		{ID: 1, Type: ElementPanel, Label: "Main", Enabled: true, Visible: true, Bounds: Rect{Row: 0, Col: 0, Width: 20, Height: 8}},
		{ID: 2, Type: ElementButton, Label: "Save", Focused: true, Enabled: true, Bounds: Rect{Row: 1, Col: 1, Width: 6, Height: 1}},
		{ID: 3, Type: ElementPanel, Label: "Nested", Enabled: true, Visible: true, Bounds: Rect{Row: 2, Col: 1, Width: 14, Height: 4}},
		{ID: 4, Type: ElementButton, Label: "Cancel", Enabled: false, Visible: true, Bounds: Rect{Row: 6, Col: 22, Width: 8, Height: 1}},
		{ID: 5, Type: ElementCheckbox, Label: "[x] Auto-save", Enabled: true, Checked: true, Visible: true, Bounds: Rect{Row: 3, Col: 2, Width: 13, Height: 1}},
		{ID: 6, Type: ElementButton, Label: "Apply", Enabled: true, Visible: true, Bounds: Rect{Row: 4, Col: 2, Width: 7, Height: 1}},
	})}
	m.Elements[0].Ref = "panel[1]"
	m.Elements[1].Ref = "panel[1]/button[1]"
	m.Elements[2].Ref = "panel[1]/panel[1]"
	m.Elements[3].Ref = "button[1]"
	m.Elements[4].Ref = "panel[1]/panel[1]/checkbox[1]"
	m.Elements[5].Ref = "panel[1]/panel[1]/button[1]"

	if got := m.Query("button"); len(got) != 3 {
		t.Fatalf("button query len = %d, want 3", len(got))
	}
	if got := m.Query(`[ref="panel[1]"] > button`); len(got) != 1 || got[0].Label != "Save" {
		t.Fatalf("outer panel > button = %+v, want Save", got)
	}
	if got := m.Query("button:focused"); len(got) != 1 || got[0].Label != "Save" {
		t.Fatalf("button:focused = %+v, want Save", got)
	}
	if got := m.Query(`[label="Save"]`); len(got) != 1 || got[0].ID != 2 {
		t.Fatalf("label selector = %+v, want ID 2", got)
	}
	if got := m.Query("checkbox:checked"); len(got) != 1 || got[0].ID != 5 {
		t.Fatalf("checkbox:checked = %+v, want ID 5", got)
	}
	if got := m.Query("button:disabled"); len(got) != 1 || got[0].ID != 4 {
		t.Fatalf("button:disabled = %+v, want ID 4", got)
	}
	if got := m.QueryOne("#5"); got == nil || got.Label != "[x] Auto-save" {
		t.Fatalf("QueryOne(#5) = %+v, want checkbox", got)
	}
	if got := m.Query("button:nth(2)"); len(got) != 1 || got[0].ID != 4 {
		t.Fatalf("button:nth(2) = %+v, want ID 4", got)
	}
	if got := m.Query(`[ref="panel[1]/button[1]"]`); len(got) != 1 || got[0].ID != 2 {
		t.Fatalf("ref selector = %+v, want ID 2", got)
	}
	if got := m.Query("panel button"); len(got) != 2 {
		t.Fatalf("panel button = %+v, want Save and Apply", got)
	}
	if got := m.Query("dialog button"); len(got) != 0 {
		t.Fatalf("no-match query should be empty, got %+v", got)
	}
}

func TestSerializeJSON_Phase3Fields(t *testing.T) {
	m := &ElementMap{
		Viewport: Rect{Row: 0, Col: 0, Width: 80, Height: 24},
		Scrolled: true,
		Elements: []Element{{
			ID:            1,
			Type:          ElementCheckbox,
			Label:         "[x] Save",
			Focused:       true,
			Enabled:       true,
			Checked:       true,
			Selected:      false,
			Visible:       true,
			Role:          "checkbox",
			Shortcut:      "Ctrl+S",
			Ref:           "checkbox[1]",
			Bounds:        Rect{Row: 2, Col: 2, Width: 9, Height: 1},
			VisibleBounds: &Rect{Row: 2, Col: 2, Width: 9, Height: 1},
		}},
	}

	var parsed map[string]any
	if err := json.Unmarshal([]byte(SerializeJSON(m)), &parsed); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if parsed["scrolled"] != true {
		t.Fatalf("scrolled = %v, want true", parsed["scrolled"])
	}
	viewport := parsed["viewport"].(map[string]any)
	if viewport["width"] != float64(80) {
		t.Fatalf("viewport width = %v, want 80", viewport["width"])
	}
	first := parsed["elements"].([]any)[0].(map[string]any)
	if first["checked"] != true || first["role"] != "checkbox" || first["ref"] != "checkbox[1]" {
		t.Fatalf("unexpected phase3 json element fields: %+v", first)
	}
	if first["shortcut"] != "Ctrl+S" || first["visible"] != true {
		t.Fatalf("unexpected shortcut/visible fields: %+v", first)
	}
}
