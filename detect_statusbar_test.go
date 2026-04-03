package awtree

import "testing"

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
