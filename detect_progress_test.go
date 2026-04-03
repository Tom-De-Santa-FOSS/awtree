package awtree

import "testing"

func TestDetect_ProgressBar_BlockChars(t *testing.T) {
	g := NewGrid(5, 40)
	for c := 5; c < 9; c++ {
		g.Set(2, c, Cell{Char: '█'})
	}
	for c := 9; c < 13; c++ {
		g.Set(2, c, Cell{Char: '░'})
	}

	m := Detect(g)
	found := false
	for _, el := range m.Elements {
		if el.Type == ElementProgressBar {
			found = true
			if el.Bounds.Row != 2 || el.Bounds.Col != 5 || el.Bounds.Width != 8 || el.Bounds.Height != 1 {
				t.Errorf("bounds = %+v, want {Row:2 Col:5 Width:8 Height:1}", el.Bounds)
			}
		}
	}
	if !found {
		t.Fatal("block character progress bar not detected")
	}
}

func TestDetect_ProgressBar_ASCII(t *testing.T) {
	g := NewGrid(5, 40)
	g.SetText(1, 3, "[====    ]", DefaultColor, DefaultColor, 0)

	m := Detect(g)
	found := false
	for _, el := range m.Elements {
		if el.Type == ElementProgressBar {
			found = true
			if el.Bounds.Row != 1 || el.Bounds.Col != 3 || el.Bounds.Width != 10 || el.Bounds.Height != 1 {
				t.Errorf("bounds = %+v, want {Row:1 Col:3 Width:10 Height:1}", el.Bounds)
			}
		}
	}
	if !found {
		t.Fatal("ASCII progress bar not detected")
	}
}

func TestDetect_ProgressBar_WithPercentage(t *testing.T) {
	g := NewGrid(5, 40)
	for c := 5; c < 10; c++ {
		g.Set(2, c, Cell{Char: '█'})
	}
	for c := 10; c < 15; c++ {
		g.Set(2, c, Cell{Char: '░'})
	}
	g.SetText(2, 16, "45%", DefaultColor, DefaultColor, 0)

	m := Detect(g)
	found := false
	for _, el := range m.Elements {
		if el.Type == ElementProgressBar {
			found = true
			if el.Label != "45%" {
				t.Errorf("label = %q, want %q", el.Label, "45%")
			}
		}
	}
	if !found {
		t.Fatal("progress bar with percentage not detected")
	}
}
