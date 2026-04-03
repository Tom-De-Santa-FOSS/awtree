package awtree

import "testing"

func TestDetect_Checkbox_Unicode(t *testing.T) {
	g := NewGrid(5, 40)
	g.Set(1, 2, Cell{Char: '☐', FG: DefaultColor, BG: DefaultColor})
	g.SetText(1, 3, " Accept terms", DefaultColor, DefaultColor, 0)
	g.Set(2, 2, Cell{Char: '☑', FG: DefaultColor, BG: DefaultColor})
	g.SetText(2, 3, " Send newsletter", DefaultColor, DefaultColor, 0)

	m := Detect(g)
	var checkboxes []Element
	for _, el := range m.Elements {
		if el.Type == ElementCheckbox {
			checkboxes = append(checkboxes, el)
		}
	}
	if len(checkboxes) != 2 {
		t.Fatalf("expected 2 checkboxes, got %d", len(checkboxes))
	}
	if checkboxes[0].Label != "☐ Accept terms" {
		t.Errorf("checkbox[0] label = %q, want %q", checkboxes[0].Label, "☐ Accept terms")
	}
	if checkboxes[1].Label != "☑ Send newsletter" {
		t.Errorf("checkbox[1] label = %q, want %q", checkboxes[1].Label, "☑ Send newsletter")
	}
}

func TestDetect_Checkbox_ASCII(t *testing.T) {
	g := NewGrid(5, 40)
	g.SetText(0, 0, "[x] Enable logging", DefaultColor, DefaultColor, 0)
	g.SetText(1, 0, "[ ] Verbose mode", DefaultColor, DefaultColor, 0)
	g.SetText(2, 0, "[X] Dark theme", DefaultColor, DefaultColor, 0)
	g.SetText(3, 0, "[*] Auto-save", DefaultColor, DefaultColor, 0)

	m := Detect(g)
	var checkboxes []Element
	for _, el := range m.Elements {
		if el.Type == ElementCheckbox {
			checkboxes = append(checkboxes, el)
		}
	}
	if len(checkboxes) != 4 {
		t.Fatalf("expected 4 checkboxes, got %d", len(checkboxes))
	}
	wantLabels := []string{"[x] Enable logging", "[ ] Verbose mode", "[X] Dark theme", "[*] Auto-save"}
	for i, want := range wantLabels {
		if checkboxes[i].Label != want {
			t.Errorf("checkbox[%d] label = %q, want %q", i, checkboxes[i].Label, want)
		}
	}
}

func TestDetect_Radio_ASCII(t *testing.T) {
	g := NewGrid(5, 40)
	g.SetText(0, 0, "( ) Option A", DefaultColor, DefaultColor, 0)
	g.SetText(1, 0, "(x) Option B", DefaultColor, DefaultColor, 0)
	g.SetText(2, 0, "(*) Option C", DefaultColor, DefaultColor, 0)

	m := Detect(g)
	var checkboxes []Element
	for _, el := range m.Elements {
		if el.Type == ElementCheckbox {
			checkboxes = append(checkboxes, el)
		}
	}
	if len(checkboxes) != 3 {
		t.Fatalf("expected 3 radio buttons, got %d", len(checkboxes))
	}
	wantLabels := []string{"( ) Option A", "(x) Option B", "(*) Option C"}
	for i, want := range wantLabels {
		if checkboxes[i].Label != want {
			t.Errorf("radio[%d] label = %q, want %q", i, checkboxes[i].Label, want)
		}
	}
}

func TestDetect_Checkbox_Focused(t *testing.T) {
	g := NewGrid(5, 40)
	g.Set(1, 2, Cell{Char: '☑', FG: DefaultColor, BG: DefaultColor, Attrs: AttrReverse})
	g.SetText(1, 3, " Agree", DefaultColor, DefaultColor, AttrReverse)
	g.SetText(2, 2, "[x]", DefaultColor, DefaultColor, AttrReverse)
	g.SetText(2, 5, " Confirm", DefaultColor, DefaultColor, AttrReverse)

	m := Detect(g)
	var checkboxes []Element
	for _, el := range m.Elements {
		if el.Type == ElementCheckbox {
			checkboxes = append(checkboxes, el)
		}
	}
	if len(checkboxes) != 2 {
		t.Fatalf("expected 2 checkboxes, got %d", len(checkboxes))
	}
	for i, cb := range checkboxes {
		if !cb.Focused {
			t.Errorf("checkbox[%d] should be focused", i)
		}
	}
}
