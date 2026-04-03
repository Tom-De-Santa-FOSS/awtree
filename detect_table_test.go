package awtree

import "testing"

func TestDetect_Table_Simple(t *testing.T) {
	g := NewGrid(5, 20)
	g.SetText(0, 0, " Alice ", DefaultColor, DefaultColor, 0)
	g.Set(0, 7, Cell{Char: '│'})
	g.SetText(0, 8, " 30", DefaultColor, DefaultColor, 0)
	g.SetText(1, 0, " Bob   ", DefaultColor, DefaultColor, 0)
	g.Set(1, 7, Cell{Char: '│'})
	g.SetText(1, 8, " 25", DefaultColor, DefaultColor, 0)
	g.SetText(2, 0, " Carol ", DefaultColor, DefaultColor, 0)
	g.Set(2, 7, Cell{Char: '│'})
	g.SetText(2, 8, " 40", DefaultColor, DefaultColor, 0)

	m := Detect(g)
	found := false
	for _, el := range m.Elements {
		if el.Type == ElementTable {
			found = true
			if el.Label != "3x2" {
				t.Errorf("table label = %q, want %q", el.Label, "3x2")
			}
			if el.Bounds.Height != 3 {
				t.Errorf("table height = %d, want 3", el.Bounds.Height)
			}
		}
	}
	if !found {
		t.Fatal("no table detected")
	}
}

func TestDetect_Table_WithHeader(t *testing.T) {
	g := NewGrid(6, 20)
	g.SetText(0, 0, " Name  ", DefaultColor, DefaultColor, AttrBold)
	g.Set(0, 7, Cell{Char: '│', Attrs: AttrBold})
	g.SetText(0, 8, " Age", DefaultColor, DefaultColor, AttrBold)
	for c := 0; c < 14; c++ {
		g.Set(1, c, Cell{Char: '─'})
	}
	g.Set(1, 7, Cell{Char: '┼'})
	g.SetText(2, 0, " Alice ", DefaultColor, DefaultColor, 0)
	g.Set(2, 7, Cell{Char: '│'})
	g.SetText(2, 8, " 30", DefaultColor, DefaultColor, 0)
	g.SetText(3, 0, " Bob   ", DefaultColor, DefaultColor, 0)
	g.Set(3, 7, Cell{Char: '│'})
	g.SetText(3, 8, " 25", DefaultColor, DefaultColor, 0)

	m := Detect(g)
	found := false
	for _, el := range m.Elements {
		if el.Type == ElementTable {
			found = true
			if el.Label != "3x2" {
				t.Errorf("table label = %q, want %q", el.Label, "3x2")
			}
			if el.Bounds.Height != 4 {
				t.Errorf("table height = %d, want 4", el.Bounds.Height)
			}
		}
	}
	if !found {
		t.Fatal("no table with header detected")
	}
}

func TestDetect_Table_ASCII(t *testing.T) {
	g := NewGrid(7, 25)
	g.SetText(0, 0, " Name  ", DefaultColor, DefaultColor, 0)
	g.Set(0, 7, Cell{Char: '|'})
	g.SetText(0, 8, " Age ", DefaultColor, DefaultColor, 0)
	g.Set(0, 13, Cell{Char: '|'})
	g.SetText(0, 14, " City", DefaultColor, DefaultColor, 0)
	for c := 0; c < 20; c++ {
		g.Set(1, c, Cell{Char: '-'})
	}
	g.Set(1, 7, Cell{Char: '+'})
	g.Set(1, 13, Cell{Char: '+'})
	for r := 2; r <= 4; r++ {
		g.SetText(r, 0, " Name  ", DefaultColor, DefaultColor, 0)
		g.Set(r, 7, Cell{Char: '|'})
		g.SetText(r, 8, " val ", DefaultColor, DefaultColor, 0)
		g.Set(r, 13, Cell{Char: '|'})
		g.SetText(r, 14, " City", DefaultColor, DefaultColor, 0)
	}

	m := Detect(g)
	found := false
	for _, el := range m.Elements {
		if el.Type == ElementTable {
			found = true
			if el.Label != "4x3" {
				t.Errorf("table label = %q, want %q", el.Label, "4x3")
			}
			if el.Bounds.Height != 5 {
				t.Errorf("table height = %d, want 5", el.Bounds.Height)
			}
		}
	}
	if !found {
		t.Fatal("no ASCII table detected")
	}
}

func TestDetect_Table_PanelNotFalsePositive(t *testing.T) {
	g := NewGrid(6, 12)
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
	for _, el := range m.Elements {
		if el.Type == ElementTable {
			t.Errorf("panel falsely detected as table: %+v", el)
		}
	}
}
