package awtree

import "testing"

func BenchmarkDetect_EmptyGrid(b *testing.B) {
	g := NewGrid(24, 80)
	b.ResetTimer()
	for b.Loop() {
		Detect(g)
	}
}

func BenchmarkDetect_RealisticTUI(b *testing.B) {
	g := buildRealisticGrid()
	b.ResetTimer()
	for b.Loop() {
		Detect(g)
	}
}

func BenchmarkSerialize(b *testing.B) {
	g := buildRealisticGrid()
	m := Detect(g)
	b.ResetTimer()
	for b.Loop() {
		Serialize(m)
	}
}

func buildRealisticGrid() *Grid {
	g := NewGrid(24, 80)

	// Menu bar.
	for c := 0; c < 80; c++ {
		g.Set(0, c, Cell{Char: ' ', BG: 4})
	}
	g.SetText(0, 1, "File  Edit  View", DefaultColor, 4, 0)

	// Panel.
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

	// Menu items inside panel.
	g.SetText(3, 2, "  Open File  ", DefaultColor, DefaultColor, 0)
	g.SetText(4, 2, "  Save File  ", DefaultColor, DefaultColor, AttrReverse)
	g.SetText(5, 2, "  Close All  ", DefaultColor, DefaultColor, 0)

	// Buttons.
	g.SetText(21, 30, "[Save]", DefaultColor, DefaultColor, 0)
	g.SetText(21, 40, "[Cancel]", DefaultColor, DefaultColor, 0)

	// Input field.
	g.SetText(22, 5, "Search: ", DefaultColor, DefaultColor, 0)
	g.SetText(22, 13, "              ", DefaultColor, 7, 0)

	// Status bar.
	for c := 0; c < 80; c++ {
		g.Set(23, c, Cell{Char: ' ', BG: 2})
	}
	g.SetText(23, 2, "Ready", DefaultColor, 2, 0)

	return g
}
