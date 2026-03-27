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

func TestDetect_Panel(t *testing.T) {
	g := NewGrid(10, 20)
	// Draw a simple box.
	g.Set(0, 0, Cell{Char: '┌'})
	for c := 1; c < 9; c++ {
		g.Set(0, c, Cell{Char: '─'})
	}
	g.Set(0, 9, Cell{Char: '┐'})

	for r := 1; r < 4; r++ {
		g.Set(r, 0, Cell{Char: '│'})
		g.Set(r, 9, Cell{Char: '│'})
	}

	g.Set(4, 0, Cell{Char: '└'})
	for c := 1; c < 9; c++ {
		g.Set(4, c, Cell{Char: '─'})
	}
	g.Set(4, 9, Cell{Char: '┘'})

	m := Detect(g)

	found := false
	for _, el := range m.Elements {
		if el.Type == ElementPanel {
			found = true
			if el.Bounds.Width != 10 || el.Bounds.Height != 5 {
				t.Errorf("panel bounds = %dx%d, want 10x5", el.Bounds.Width, el.Bounds.Height)
			}
		}
	}
	if !found {
		t.Fatal("no panel detected")
	}
}

func TestDetect_PanelWithTitle(t *testing.T) {
	g := NewGrid(10, 30)
	// ┌─ Title ─┐
	g.Set(0, 0, Cell{Char: '┌'})
	g.Set(0, 1, Cell{Char: '─'})
	g.Set(0, 2, Cell{Char: ' '})
	g.SetText(0, 3, "Title", DefaultColor, DefaultColor, 0)
	g.Set(0, 8, Cell{Char: ' '})
	g.Set(0, 9, Cell{Char: '─'})
	g.Set(0, 10, Cell{Char: '┐'})

	for r := 1; r < 4; r++ {
		g.Set(r, 0, Cell{Char: '│'})
		g.Set(r, 10, Cell{Char: '│'})
	}

	g.Set(4, 0, Cell{Char: '└'})
	for c := 1; c < 10; c++ {
		g.Set(4, c, Cell{Char: '─'})
	}
	g.Set(4, 10, Cell{Char: '┘'})

	m := Detect(g)

	for _, el := range m.Elements {
		if el.Type == ElementPanel {
			if el.Label != "Title" {
				t.Errorf("panel label = %q, want %q", el.Label, "Title")
			}
			return
		}
	}
	t.Fatal("no panel detected")
}

func TestDetect_Button(t *testing.T) {
	g := NewGrid(5, 40)
	g.SetText(2, 10, "[Save]", DefaultColor, DefaultColor, 0)
	g.SetText(2, 20, "<Cancel>", DefaultColor, DefaultColor, 0)

	m := Detect(g)

	labels := make(map[string]bool)
	for _, el := range m.Elements {
		if el.Type == ElementButton {
			labels[el.Label] = true
		}
	}

	if !labels["Save"] {
		t.Error("Save button not detected")
	}
	if !labels["Cancel"] {
		t.Error("Cancel button not detected")
	}
}

func TestDetect_FocusedButton(t *testing.T) {
	g := NewGrid(5, 40)
	g.SetText(2, 10, "[OK]", DefaultColor, DefaultColor, AttrReverse)

	m := Detect(g)

	for _, el := range m.Elements {
		if el.Type == ElementButton && el.Label == "OK" {
			if !el.Focused {
				t.Error("expected focused button")
			}
			return
		}
	}
	t.Fatal("focused button not detected")
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

func TestDetect_StatusBar(t *testing.T) {
	g := NewGrid(24, 80)
	// Fill bottom row with colored background.
	for c := 0; c < 80; c++ {
		g.Set(23, c, Cell{Char: ' ', BG: 4})
	}
	g.SetText(23, 2, "status: ready", DefaultColor, 4, 0)

	m := Detect(g)

	found := false
	for _, el := range m.Elements {
		if el.Type == ElementStatusBar {
			found = true
		}
	}
	if !found {
		t.Fatal("status bar not detected")
	}
}

func TestDetect_MenuItems_VerticalListWithOneHighlighted(t *testing.T) {
	g := NewGrid(10, 30)
	// 3 items at same column, one reverse-video.
	g.SetText(2, 2, "  Open File  ", DefaultColor, DefaultColor, 0)
	g.SetText(3, 2, "  Save File  ", DefaultColor, DefaultColor, AttrReverse)
	g.SetText(4, 2, "  Close All  ", DefaultColor, DefaultColor, 0)

	m := Detect(g)

	var menuItems []Element
	for _, el := range m.Elements {
		if el.Type == ElementMenuItem {
			menuItems = append(menuItems, el)
		}
	}

	if len(menuItems) < 3 {
		t.Fatalf("expected 3 menu items, got %d", len(menuItems))
	}

	focusCount := 0
	for _, el := range menuItems {
		if el.Focused {
			focusCount++
			if el.Label != "Save File" {
				t.Errorf("focused item label = %q, want %q", el.Label, "Save File")
			}
		}
	}
	if focusCount != 1 {
		t.Errorf("expected 1 focused item, got %d", focusCount)
	}
}

func TestDetect_TabBar(t *testing.T) {
	g := NewGrid(10, 60)
	// Horizontal tabs: one bold/reverse (active), others normal.
	g.SetText(0, 0, " Files ", DefaultColor, 4, AttrReverse|AttrBold)
	g.SetText(0, 7, " Edit ", DefaultColor, DefaultColor, 0)
	g.SetText(0, 13, " View ", DefaultColor, DefaultColor, 0)

	m := Detect(g)

	var tabs []Element
	for _, el := range m.Elements {
		if el.Type == ElementTab {
			tabs = append(tabs, el)
		}
	}

	if len(tabs) < 3 {
		t.Fatalf("expected 3 tabs, got %d", len(tabs))
	}

	focusCount := 0
	for _, el := range tabs {
		if el.Focused {
			focusCount++
			if el.Label != "Files" {
				t.Errorf("active tab = %q, want %q", el.Label, "Files")
			}
		}
	}
	if focusCount != 1 {
		t.Errorf("expected 1 focused tab, got %d", focusCount)
	}
}

func TestDetect_InputField_UnderscorePlaceholder(t *testing.T) {
	g := NewGrid(10, 40)
	// Label followed by underscores (common input pattern).
	g.SetText(3, 2, "Name: ", DefaultColor, DefaultColor, 0)
	g.SetText(3, 8, "________________", DefaultColor, DefaultColor, AttrUnderline)

	m := Detect(g)

	found := false
	for _, el := range m.Elements {
		if el.Type == ElementInput {
			found = true
			if el.Bounds.Row != 3 {
				t.Errorf("input row = %d, want 3", el.Bounds.Row)
			}
		}
	}
	if !found {
		t.Fatal("input field not detected")
	}
}

func TestDetect_InputField_DistinctBackground(t *testing.T) {
	g := NewGrid(10, 40)
	// Input field: region with distinct BG color.
	g.SetText(5, 10, "              ", DefaultColor, 7, 0) // white bg, empty
	g.SetText(5, 10, "hello", DefaultColor, 7, 0)

	m := Detect(g)

	found := false
	for _, el := range m.Elements {
		if el.Type == ElementInput {
			found = true
			if el.Label != "hello" {
				t.Errorf("input value = %q, want %q", el.Label, "hello")
			}
		}
	}
	if !found {
		t.Fatal("input field with distinct BG not detected")
	}
}

func TestDetect_BracketedNonButton_ArrayNotation(t *testing.T) {
	g := NewGrid(5, 40)
	// [0] and [123] should NOT be buttons (start with digit).
	g.SetText(1, 5, "[0]", DefaultColor, DefaultColor, 0)
	g.SetText(2, 5, "[123]", DefaultColor, DefaultColor, 0)

	m := Detect(g)

	for _, el := range m.Elements {
		if el.Type == ElementButton {
			t.Errorf("false positive button: %q", el.Label)
		}
	}
}

func TestDetect_NestedPanels(t *testing.T) {
	g := NewGrid(12, 30)
	// Outer panel.
	g.Set(0, 0, Cell{Char: '┌'})
	for c := 1; c < 19; c++ {
		g.Set(0, c, Cell{Char: '─'})
	}
	g.Set(0, 19, Cell{Char: '┐'})
	for r := 1; r < 9; r++ {
		g.Set(r, 0, Cell{Char: '│'})
		g.Set(r, 19, Cell{Char: '│'})
	}
	g.Set(9, 0, Cell{Char: '└'})
	for c := 1; c < 19; c++ {
		g.Set(9, c, Cell{Char: '─'})
	}
	g.Set(9, 19, Cell{Char: '┘'})

	// Inner panel.
	g.Set(2, 2, Cell{Char: '┌'})
	for c := 3; c < 9; c++ {
		g.Set(2, c, Cell{Char: '─'})
	}
	g.Set(2, 9, Cell{Char: '┐'})
	for r := 3; r < 6; r++ {
		g.Set(r, 2, Cell{Char: '│'})
		g.Set(r, 9, Cell{Char: '│'})
	}
	g.Set(6, 2, Cell{Char: '└'})
	for c := 3; c < 9; c++ {
		g.Set(6, c, Cell{Char: '─'})
	}
	g.Set(6, 9, Cell{Char: '┘'})

	m := Detect(g)

	panelCount := 0
	for _, el := range m.Elements {
		if el.Type == ElementPanel {
			panelCount++
		}
	}
	if panelCount != 2 {
		t.Errorf("expected 2 panels (nested), got %d", panelCount)
	}
}

func TestDetect_MenuBar_TopRowWithBG(t *testing.T) {
	g := NewGrid(24, 80)
	// Top row with colored BG → menu bar, not status bar.
	for c := 0; c < 80; c++ {
		g.Set(0, c, Cell{Char: ' ', BG: 2})
	}
	g.SetText(0, 2, "File  Edit  View  Help", DefaultColor, 2, 0)

	m := Detect(g)

	found := false
	for _, el := range m.Elements {
		if el.Type == ElementMenuBar {
			found = true
		}
		if el.Type == ElementStatusBar && el.Bounds.Row == 0 {
			t.Error("top row should be MenuBar, not StatusBar")
		}
	}
	if !found {
		t.Fatal("menu bar not detected on top row")
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
