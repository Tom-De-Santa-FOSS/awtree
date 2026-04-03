package awtree

import "testing"

func TestDetect_TabBar(t *testing.T) {
	g := NewGrid(10, 60)
	// Horizontal tabs: one bold/reverse (active), others normal.
	g.SetText(0, 0, " Files ", DefaultColor, 4, AttrReverse|AttrBold)
	g.SetText(0, 7, " Edit ", DefaultColor, DefaultColor, 0)
	g.SetText(0, 13, " View ", DefaultColor, DefaultColor, 0)

	m := Detect(g)

	var tabs []Element
	for _, el := range m.Elements {
		if el.Type == ElementTab {
			tabs = append(tabs, el)
		}
	}

	if len(tabs) < 3 {
		t.Fatalf("expected 3 tabs, got %d", len(tabs))
	}

	focusCount := 0
	for _, el := range tabs {
		if el.Focused {
			focusCount++
			if el.Label != "Files" {
				t.Errorf("active tab = %q, want %q", el.Label, "Files")
			}
		}
	}
	if focusCount != 1 {
		t.Errorf("expected 1 focused tab, got %d", focusCount)
	}
}

func TestDetect_RowWithTwoActiveSegments_NotDetectedAsTabs(t *testing.T) {
	g := NewGrid(5, 40)
	g.SetText(0, 0, " Files ", DefaultColor, 4, AttrReverse)
	g.SetText(0, 10, " Edit ", DefaultColor, 4, AttrReverse)
	g.SetText(0, 20, " View ", DefaultColor, DefaultColor, 0)

	m := Detect(g)
	for _, el := range m.Elements {
		if el.Type == ElementTab {
			t.Error("row with 2 active segments should not produce tabs")
		}
	}
}
