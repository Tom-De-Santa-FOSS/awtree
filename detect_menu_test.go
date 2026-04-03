package awtree

import "testing"

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

	if len(menuItems) != 3 {
		t.Fatalf("expected exactly 3 menu items, got %d", len(menuItems))
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

func TestDetect_MenuItems_NoDuplicatesWithMultipleHighlights(t *testing.T) {
	g := NewGrid(10, 30)
	// Two highlighted items in same column range with shared siblings.
	g.SetText(2, 2, "  Open File  ", DefaultColor, DefaultColor, 0)
	g.SetText(3, 2, "  Save File  ", DefaultColor, DefaultColor, AttrReverse)
	g.SetText(4, 2, "  Close All  ", DefaultColor, DefaultColor, AttrReverse)
	g.SetText(5, 2, "  Quit       ", DefaultColor, DefaultColor, 0)

	m := Detect(g)

	var menuItems []Element
	for _, el := range m.Elements {
		if el.Type == ElementMenuItem {
			menuItems = append(menuItems, el)
		}
	}

	if len(menuItems) != 4 {
		t.Fatalf("expected exactly 4 menu items, got %d", len(menuItems))
	}

	// Check no duplicate rows.
	seen := make(map[int]bool)
	for _, el := range menuItems {
		if seen[el.Bounds.Row] {
			t.Errorf("duplicate menu item at row %d", el.Bounds.Row)
		}
		seen[el.Bounds.Row] = true
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
