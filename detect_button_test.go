package awtree

import "testing"

func TestDetect_Button(t *testing.T) {
	g := NewGrid(5, 40)
	g.SetText(2, 10, "[Save]", DefaultColor, DefaultColor, 0)
	g.SetText(2, 20, "<Cancel>", DefaultColor, DefaultColor, 0)

	m := Detect(g)

	labels := make(map[string]bool)
	for _, el := range m.Elements {
		if el.Type == ElementButton {
			labels[el.Label] = true
		}
	}

	if !labels["Save"] {
		t.Error("Save button not detected")
	}
	if !labels["Cancel"] {
		t.Error("Cancel button not detected")
	}
}

func TestDetect_FocusedButton(t *testing.T) {
	g := NewGrid(5, 40)
	g.SetText(2, 10, "[OK]", DefaultColor, DefaultColor, AttrReverse)

	m := Detect(g)

	for _, el := range m.Elements {
		if el.Type == ElementButton && el.Label == "OK" {
			if !el.Focused {
				t.Error("expected focused button")
			}
			return
		}
	}
	t.Fatal("focused button not detected")
}

func TestDetect_BracketedNonButton_ArrayNotation(t *testing.T) {
	g := NewGrid(5, 40)
	// [0] and [123] should NOT be buttons (start with digit).
	g.SetText(1, 5, "[0]", DefaultColor, DefaultColor, 0)
	g.SetText(2, 5, "[123]", DefaultColor, DefaultColor, 0)

	m := Detect(g)

	for _, el := range m.Elements {
		if el.Type == ElementButton {
			t.Errorf("false positive button: %q", el.Label)
		}
	}
}
