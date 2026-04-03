package awtree

import "testing"

func TestDetect_InputField_UnderscorePlaceholder(t *testing.T) {
	g := NewGrid(10, 40)
	// Label followed by underscores (common input pattern).
	g.SetText(3, 2, "Name: ", DefaultColor, DefaultColor, 0)
	g.SetText(3, 8, "________________", DefaultColor, DefaultColor, AttrUnderline)

	m := Detect(g)

	found := false
	for _, el := range m.Elements {
		if el.Type == ElementInput {
			found = true
			if el.Bounds.Row != 3 {
				t.Errorf("input row = %d, want 3", el.Bounds.Row)
			}
		}
	}
	if !found {
		t.Fatal("input field not detected")
	}
}

func TestDetect_InputField_DistinctBackground(t *testing.T) {
	g := NewGrid(10, 40)
	// Input field: region with distinct BG color.
	g.SetText(5, 10, "              ", DefaultColor, 7, 0) // white bg, empty
	g.SetText(5, 10, "hello", DefaultColor, 7, 0)

	m := Detect(g)

	found := false
	for _, el := range m.Elements {
		if el.Type == ElementInput {
			found = true
			if el.Label != "hello" {
				t.Errorf("input value = %q, want %q", el.Label, "hello")
			}
		}
	}
	if !found {
		t.Fatal("input field with distinct BG not detected")
	}
}

func TestDetect_InputField_TooShortUnderlinedSpan_NotDetected(t *testing.T) {
	g := NewGrid(5, 20)
	// Width=2, below minimum threshold of 3.
	g.Set(2, 5, Cell{Char: '_', Attrs: AttrUnderline})
	g.Set(2, 6, Cell{Char: '_', Attrs: AttrUnderline})

	m := Detect(g)
	for _, el := range m.Elements {
		if el.Type == ElementInput {
			t.Error("short underlined span should not be detected as input")
		}
	}
}

func TestDetect_InputField_WideDistinctBGSpan_NotDetectedAsInput(t *testing.T) {
	g := NewGrid(10, 20) // 20 cols, threshold is 12
	for c := 0; c < 16; c++ {
		g.Set(5, c, Cell{Char: ' ', BG: 7})
	}

	m := Detect(g)
	for _, el := range m.Elements {
		if el.Type == ElementInput {
			t.Error("wide BG span should not be classified as input")
		}
	}
}
