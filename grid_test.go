package termtree

import "testing"

func TestNewGrid(t *testing.T) {
	g := NewGrid(24, 80)
	if g.Rows != 24 || g.Cols != 80 {
		t.Fatalf("dims = %dx%d, want 24x80", g.Rows, g.Cols)
	}
	// All cells should be spaces with default colors.
	c := g.At(0, 0)
	if c.Char != ' ' {
		t.Errorf("default char = %q, want space", c.Char)
	}
	if c.FG != DefaultColor || c.BG != DefaultColor {
		t.Errorf("default colors = %d/%d, want %d/%d", c.FG, c.BG, DefaultColor, DefaultColor)
	}
}

func TestGrid_OutOfBounds(t *testing.T) {
	g := NewGrid(5, 5)
	c := g.At(-1, 0)
	if c.Char != 0 {
		t.Error("out of bounds should return zero Cell")
	}
	c = g.At(10, 10)
	if c.Char != 0 {
		t.Error("out of bounds should return zero Cell")
	}
	// Set out of bounds should not panic.
	g.Set(-1, 0, Cell{Char: 'x'})
	g.Set(10, 10, Cell{Char: 'x'})
}

func TestGrid_SetText(t *testing.T) {
	g := NewGrid(5, 20)
	g.SetText(2, 3, "Hello", 1, 2, AttrBold)

	for i, ch := range "Hello" {
		c := g.At(2, 3+i)
		if c.Char != ch {
			t.Errorf("cell[2,%d] = %q, want %q", 3+i, c.Char, ch)
		}
		if c.Attrs&AttrBold == 0 {
			t.Errorf("cell[2,%d] missing bold attr", 3+i)
		}
	}
}
