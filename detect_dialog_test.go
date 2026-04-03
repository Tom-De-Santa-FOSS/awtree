package awtree

import "testing"

func TestDetect_Dialog_CenteredPanelWithButton(t *testing.T) {
	g := NewGrid(24, 80)
	topRow, botRow := 8, 14
	leftCol, rightCol := 30, 50
	g.Set(topRow, leftCol, Cell{Char: '┌'})
	g.Set(topRow, rightCol, Cell{Char: '┐'})
	g.Set(botRow, leftCol, Cell{Char: '└'})
	g.Set(botRow, rightCol, Cell{Char: '┘'})
	for c := leftCol + 1; c < rightCol; c++ {
		g.Set(topRow, c, Cell{Char: '─'})
		g.Set(botRow, c, Cell{Char: '─'})
	}
	for r := topRow + 1; r < botRow; r++ {
		g.Set(r, leftCol, Cell{Char: '│'})
		g.Set(r, rightCol, Cell{Char: '│'})
	}
	g.SetText(topRow, leftCol+2, "─ Confirm ─", DefaultColor, DefaultColor, 0)
	g.SetText(12, 38, "[OK]", DefaultColor, DefaultColor, 0)

	m := Detect(g)

	found := false
	for _, el := range m.Elements {
		if el.Type == ElementDialog {
			found = true
			if el.Bounds.Row != topRow || el.Bounds.Col != leftCol {
				t.Errorf("dialog bounds start = (%d,%d), want (%d,%d)", el.Bounds.Row, el.Bounds.Col, topRow, leftCol)
			}
		}
	}
	if !found {
		t.Fatal("dialog not detected for centered panel with button")
	}
}

func TestDetect_Dialog_CenteredPanelNoButtons(t *testing.T) {
	g := NewGrid(24, 80)
	topRow, botRow := 8, 14
	leftCol, rightCol := 30, 50
	g.Set(topRow, leftCol, Cell{Char: '┌'})
	g.Set(topRow, rightCol, Cell{Char: '┐'})
	g.Set(botRow, leftCol, Cell{Char: '└'})
	g.Set(botRow, rightCol, Cell{Char: '┘'})
	for c := leftCol + 1; c < rightCol; c++ {
		g.Set(topRow, c, Cell{Char: '─'})
		g.Set(botRow, c, Cell{Char: '─'})
	}
	for r := topRow + 1; r < botRow; r++ {
		g.Set(r, leftCol, Cell{Char: '│'})
		g.Set(r, rightCol, Cell{Char: '│'})
	}
	g.SetText(10, 33, "Loading...", DefaultColor, DefaultColor, 0)

	m := Detect(g)

	for _, el := range m.Elements {
		if el.Type == ElementDialog {
			t.Fatal("centered panel without buttons should not be detected as dialog")
		}
	}
}

func TestDetect_Dialog_OffCenterPanelNotDialog(t *testing.T) {
	g := NewGrid(24, 80)
	topRow, botRow := 8, 14
	leftCol, rightCol := 0, 20
	g.Set(topRow, leftCol, Cell{Char: '┌'})
	g.Set(topRow, rightCol, Cell{Char: '┐'})
	g.Set(botRow, leftCol, Cell{Char: '└'})
	g.Set(botRow, rightCol, Cell{Char: '┘'})
	for c := leftCol + 1; c < rightCol; c++ {
		g.Set(topRow, c, Cell{Char: '─'})
		g.Set(botRow, c, Cell{Char: '─'})
	}
	for r := topRow + 1; r < botRow; r++ {
		g.Set(r, leftCol, Cell{Char: '│'})
		g.Set(r, rightCol, Cell{Char: '│'})
	}
	g.SetText(12, 5, "[OK]", DefaultColor, DefaultColor, 0)

	m := Detect(g)

	for _, el := range m.Elements {
		if el.Type == ElementDialog {
			t.Fatal("off-center panel should not be detected as dialog")
		}
	}
}
