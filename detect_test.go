package awtree

import "testing"

func TestDetect_NilGrid(t *testing.T) {
	m := Detect(nil)
	if len(m.Elements) != 0 {
		t.Fatalf("expected empty, got %d elements", len(m.Elements))
	}
}

func TestDetect_EmptyGrid(t *testing.T) {
	g := NewGrid(24, 80)
	m := Detect(g)
	if len(m.Elements) != 0 {
		t.Fatalf("expected empty, got %d elements", len(m.Elements))
	}
}

func TestDetect_RealisticTUI(t *testing.T) {
	g := NewGrid(24, 80)

	// Menu bar (row 0, colored BG).
	for c := 0; c < 80; c++ {
		g.Set(0, c, Cell{Char: ' ', BG: 4})
	}
	g.SetText(0, 1, "File  Edit  View", DefaultColor, 4, 0)

	// Panel (rows 1-20).
	g.Set(1, 0, Cell{Char: '┌'})
	for c := 1; c < 59; c++ {
		g.Set(1, c, Cell{Char: '─'})
	}
	g.Set(1, 59, Cell{Char: '┐'})
	for r := 2; r < 20; r++ {
		g.Set(r, 0, Cell{Char: '│'})
		g.Set(r, 59, Cell{Char: '│'})
	}
	g.Set(20, 0, Cell{Char: '└'})
	for c := 1; c < 59; c++ {
		g.Set(20, c, Cell{Char: '─'})
	}
	g.Set(20, 59, Cell{Char: '┘'})

	// Buttons on row 21.
	g.SetText(21, 30, "[Save]", DefaultColor, DefaultColor, 0)
	g.SetText(21, 40, "[Cancel]", DefaultColor, DefaultColor, 0)

	// Status bar (bottom row).
	for c := 0; c < 80; c++ {
		g.Set(23, c, Cell{Char: ' ', BG: 2})
	}
	g.SetText(23, 2, "Ready", DefaultColor, 2, 0)

	m := Detect(g)

	types := make(map[ElementType]int)
	for _, el := range m.Elements {
		types[el.Type]++
	}

	if types[ElementPanel] < 1 {
		t.Error("missing panel")
	}
	if types[ElementButton] < 2 {
		t.Errorf("expected 2+ buttons, got %d", types[ElementButton])
	}
	if types[ElementMenuBar] < 1 {
		t.Error("missing menu bar")
	}
	if types[ElementStatusBar] < 1 {
		t.Error("missing status bar")
	}

	// All IDs unique and > 0.
	ids := make(map[int]bool)
	for _, el := range m.Elements {
		if el.ID == 0 {
			t.Error("ID 0 found")
		}
		if ids[el.ID] {
			t.Errorf("duplicate ID %d", el.ID)
		}
		ids[el.ID] = true
	}
}

func TestDetect_FirstElementIDIsOne(t *testing.T) {
	g := NewGrid(5, 20)
	g.SetText(2, 5, "[OK]", DefaultColor, DefaultColor, 0)

	m := Detect(g)

	if len(m.Elements) == 0 {
		t.Fatal("expected at least one element")
	}
	if m.Elements[0].ID != 1 {
		t.Errorf("first element ID = %d, want 1", m.Elements[0].ID)
	}
}

func TestDetect_SequentialIDs(t *testing.T) {
	g := NewGrid(10, 40)
	// Panel.
	g.Set(0, 0, Cell{Char: '┌'})
	g.Set(0, 1, Cell{Char: '─'})
	g.Set(0, 2, Cell{Char: '┐'})
	g.Set(1, 0, Cell{Char: '│'})
	g.Set(1, 2, Cell{Char: '│'})
	g.Set(2, 0, Cell{Char: '└'})
	g.Set(2, 1, Cell{Char: '─'})
	g.Set(2, 2, Cell{Char: '┘'})

	// Button.
	g.SetText(5, 10, "[OK]", DefaultColor, DefaultColor, 0)

	m := Detect(g)

	ids := make(map[int]bool)
	for _, el := range m.Elements {
		if el.ID == 0 {
			t.Error("element has ID 0, should start at 1")
		}
		if ids[el.ID] {
			t.Errorf("duplicate ID %d", el.ID)
		}
		ids[el.ID] = true
	}
}

func TestOverlapsAny_BoundingBoxIntersection(t *testing.T) {
	// A reverse region starting inside a menu item's bounds should overlap.
	menuItem := Element{
		Type:   ElementMenuItem,
		Bounds: Rect{Row: 3, Col: 2, Width: 13, Height: 1},
	}
	// Reverse region starts at col 5 (inside menu item col 2..14).
	reverseRegion := Element{
		Bounds: Rect{Row: 3, Col: 5, Width: 8, Height: 1},
	}

	if !overlapsAny(reverseRegion, []Element{menuItem}) {
		t.Error("expected overlapsAny to return true for bounding-box intersection")
	}
}

func TestOverlapsAny_NoOverlap(t *testing.T) {
	menuItem := Element{
		Bounds: Rect{Row: 3, Col: 2, Width: 13, Height: 1},
	}
	other := Element{
		Bounds: Rect{Row: 5, Col: 2, Width: 13, Height: 1},
	}

	if overlapsAny(other, []Element{menuItem}) {
		t.Error("expected no overlap for different rows")
	}
}

func TestIsButtonLabel(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"", false},
		{"Save", true},
		{"OK", true},
		{"Cancel", true},
		{"0", false},   // starts with digit
		{"123", false}, // starts with digit
		{"a very long button label that exceeds twenty chars", false},
		{string([]rune{0x01}), false}, // non-printable
	}
	for _, tt := range tests {
		got := isButtonLabel(tt.input, DefaultDetectConfig().MaxButtonLabelLen)
		if got != tt.want {
			t.Errorf("isButtonLabel(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestTrimSpaces(t *testing.T) {
	tests := []struct {
		input []rune
		want  string
	}{
		{[]rune("  hello  "), "hello"},
		{[]rune{}, ""},
		{[]rune("   "), ""},
		{[]rune("hello"), "hello"},
		{[]rune("  hello"), "hello"},
		{[]rune("hello  "), "hello"},
	}
	for _, tt := range tests {
		got := trimSpaces(tt.input)
		if got != tt.want {
			t.Errorf("trimSpaces(%q) = %q, want %q", string(tt.input), got, tt.want)
		}
	}
}

func TestExtractLineText(t *testing.T) {
	g := NewGrid(5, 20)
	g.SetText(2, 5, "Hello", DefaultColor, DefaultColor, 0)

	got := extractLineText(g, 2, 5, 5)
	if got != "Hello" {
		t.Errorf("extractLineText = %q, want %q", got, "Hello")
	}

	// Blank region returns empty string.
	got = extractLineText(g, 3, 5, 5)
	if got != "" {
		t.Errorf("extractLineText blank = %q, want empty", got)
	}
}

func TestDetect_ReverseRegion(t *testing.T) {
	g := NewGrid(10, 40)
	g.SetText(3, 2, "  Selected Item  ", DefaultColor, 4, AttrReverse)

	m := Detect(g)

	found := false
	for _, el := range m.Elements {
		if el.Focused && el.Label == "Selected Item" {
			found = true
		}
	}
	if !found {
		t.Fatal("reverse-video region not detected")
	}
}
