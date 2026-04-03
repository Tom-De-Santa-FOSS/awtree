package awtree

import "testing"

func TestDetect_Separator_BasicDash(t *testing.T) {
	g := NewGrid(5, 20)
	for c := 3; c < 13; c++ {
		g.Set(2, c, Cell{Char: '─'})
	}

	m := Detect(g)

	found := false
	for _, el := range m.Elements {
		if el.Type == ElementSeparator {
			found = true
			if el.Bounds.Row != 2 {
				t.Errorf("separator row = %d, want 2", el.Bounds.Row)
			}
			if el.Bounds.Col != 3 {
				t.Errorf("separator col = %d, want 3", el.Bounds.Col)
			}
			if el.Bounds.Width != 10 {
				t.Errorf("separator width = %d, want 10", el.Bounds.Width)
			}
			if el.Bounds.Height != 1 {
				t.Errorf("separator height = %d, want 1", el.Bounds.Height)
			}
			if el.Label != "" {
				t.Errorf("separator label = %q, want empty", el.Label)
			}
		}
	}
	if !found {
		t.Fatal("separator not detected")
	}
}

func TestDetect_Separator_ShortDash(t *testing.T) {
	g := NewGrid(5, 20)
	g.Set(2, 5, Cell{Char: '-'})
	g.Set(2, 6, Cell{Char: '-'})

	m := Detect(g)

	for _, el := range m.Elements {
		if el.Type == ElementSeparator {
			t.Error("2-char dash run should not be detected as separator")
		}
	}
}

func TestDetect_Separator_PanelBorderNotSeparator(t *testing.T) {
	g := NewGrid(5, 20)
	g.Set(2, 2, Cell{Char: '┌'})
	for c := 3; c < 13; c++ {
		g.Set(2, c, Cell{Char: '─'})
	}
	g.Set(2, 13, Cell{Char: '┐'})

	m := Detect(g)

	for _, el := range m.Elements {
		if el.Type == ElementSeparator {
			t.Error("panel border row should not be detected as separator")
		}
	}
}
