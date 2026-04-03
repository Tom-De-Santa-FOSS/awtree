package awtree

import "testing"

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

func TestDetect_RoundedCornerPanel(t *testing.T) {
	g := NewGrid(5, 12)
	// Lazygit-style rounded corners: ╭╮╰╯
	g.Set(0, 0, Cell{Char: '╭'})
	for c := 1; c < 9; c++ {
		g.Set(0, c, Cell{Char: '─'})
	}
	g.Set(0, 9, Cell{Char: '╮'})
	for r := 1; r < 4; r++ {
		g.Set(r, 0, Cell{Char: '│'})
		g.Set(r, 9, Cell{Char: '│'})
	}
	g.Set(4, 0, Cell{Char: '╰'})
	for c := 1; c < 9; c++ {
		g.Set(4, c, Cell{Char: '─'})
	}
	g.Set(4, 9, Cell{Char: '╯'})

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
		t.Error("rounded corner panel not detected")
	}
}
