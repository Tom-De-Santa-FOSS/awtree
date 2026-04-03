package awtree

import "testing"

func TestDetect_ScrollIndicator_ArrowPair(t *testing.T) {
	g := NewGrid(10, 20)
	g.Set(0, 10, Cell{Char: '▲'})
	g.Set(5, 10, Cell{Char: '▼'})

	m := Detect(g)

	found := false
	for _, el := range m.Elements {
		if el.Type == ElementScrollIndicator {
			found = true
			if el.Bounds.Col != 10 {
				t.Errorf("scroll indicator col = %d, want 10", el.Bounds.Col)
			}
			if el.Bounds.Row != 0 {
				t.Errorf("scroll indicator row = %d, want 0", el.Bounds.Row)
			}
			if el.Bounds.Height != 6 {
				t.Errorf("scroll indicator height = %d, want 6", el.Bounds.Height)
			}
			if el.Bounds.Width != 1 {
				t.Errorf("scroll indicator width = %d, want 1", el.Bounds.Width)
			}
		}
	}
	if !found {
		t.Fatal("scroll indicator (arrow pair) not detected")
	}
}

func TestDetect_ScrollIndicator_BlockThumb(t *testing.T) {
	g := NewGrid(10, 20)
	g.Set(2, 15, Cell{Char: '█'})
	g.Set(3, 15, Cell{Char: '█'})
	g.Set(4, 15, Cell{Char: '█'})
	g.Set(5, 15, Cell{Char: '█'})

	m := Detect(g)

	found := false
	for _, el := range m.Elements {
		if el.Type == ElementScrollIndicator {
			found = true
			if el.Bounds.Col != 15 {
				t.Errorf("scroll indicator col = %d, want 15", el.Bounds.Col)
			}
			if el.Bounds.Row != 2 {
				t.Errorf("scroll indicator row = %d, want 2", el.Bounds.Row)
			}
			if el.Bounds.Height != 4 {
				t.Errorf("scroll indicator height = %d, want 4", el.Bounds.Height)
			}
		}
	}
	if !found {
		t.Fatal("scroll indicator (block thumb) not detected")
	}
}

func TestDetect_ScrollIndicator_TooShort(t *testing.T) {
	g := NewGrid(10, 20)
	g.Set(3, 12, Cell{Char: '█'})
	g.Set(4, 12, Cell{Char: '█'})

	m := Detect(g)

	for _, el := range m.Elements {
		if el.Type == ElementScrollIndicator {
			t.Errorf("short block run (2 chars) should not be detected as scroll indicator")
		}
	}
}
