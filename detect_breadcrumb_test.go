package awtree

import "testing"

func TestDetect_Breadcrumb_AngleBracket(t *testing.T) {
	g := NewGrid(5, 80)
	g.SetText(1, 2, "Home > Settings > Display", DefaultColor, DefaultColor, 0)

	m := Detect(g)

	found := false
	for _, el := range m.Elements {
		if el.Type == ElementBreadcrumb {
			found = true
			if el.Label != "Home > Settings > Display" {
				t.Errorf("breadcrumb label = %q, want %q", el.Label, "Home > Settings > Display")
			}
			if el.Bounds.Row != 1 || el.Bounds.Col != 2 {
				t.Errorf("breadcrumb position = (%d,%d), want (1,2)", el.Bounds.Row, el.Bounds.Col)
			}
			if el.Bounds.Width != 25 || el.Bounds.Height != 1 {
				t.Errorf("breadcrumb size = %dx%d, want 25x1", el.Bounds.Width, el.Bounds.Height)
			}
		}
	}
	if !found {
		t.Fatal("breadcrumb not detected")
	}
}

func TestDetect_Breadcrumb_Slash(t *testing.T) {
	g := NewGrid(5, 80)
	g.SetText(0, 0, "File / Edit / View / Help", DefaultColor, DefaultColor, 0)

	m := Detect(g)

	found := false
	for _, el := range m.Elements {
		if el.Type == ElementBreadcrumb {
			found = true
			if el.Label != "File / Edit / View / Help" {
				t.Errorf("breadcrumb label = %q, want %q", el.Label, "File / Edit / View / Help")
			}
		}
	}
	if !found {
		t.Fatal("slash-separated breadcrumb not detected")
	}
}

func TestDetect_Breadcrumb_TooFewSegments(t *testing.T) {
	g := NewGrid(5, 80)
	g.SetText(1, 2, "Home > Settings", DefaultColor, DefaultColor, 0)

	m := Detect(g)

	for _, el := range m.Elements {
		if el.Type == ElementBreadcrumb {
			t.Errorf("breadcrumb should not be detected with only 2 segments, got label %q", el.Label)
		}
	}
}
