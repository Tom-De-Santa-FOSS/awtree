package awtree

import "testing"

func TestDetect_Phase4ConfigurableButtonThresholdsAndDebug(t *testing.T) {
	g := NewGrid(3, 48)
	g.SetText(1, 2, "[Save changes before exit]", DefaultColor, DefaultColor, 0)

	defaultMap := Detect(g, WithDebug(true))
	if len(defaultMap.Elements) != 0 {
		t.Fatalf("default detect found %+v, want no button", defaultMap.Elements)
	}
	if defaultMap.Debug == nil || len(defaultMap.Debug.Events) == 0 {
		t.Fatal("expected debug events for rejected button")
	}

	configured := Detect(g, WithDebug(true), WithMaxButtonLabelLen(32), WithMaxButtonWidth(32))
	if len(configured.Elements) != 1 || configured.Elements[0].Type != ElementButton {
		t.Fatalf("configured detect = %+v, want one button", configured.Elements)
	}
	if configured.Debug == nil || !configured.Debug.Config.Debug {
		t.Fatal("expected debug info with config echo")
	}
}

func TestDetect_Phase4ConfigurableBackgroundThreshold(t *testing.T) {
	g := NewGrid(6, 10)
	g.SetText(3, 1, "abcdefgh", DefaultColor, PaletteColor(7), 0)

	defaultMap := Detect(g)
	for _, el := range defaultMap.Elements {
		if el.Type == ElementInput {
			t.Fatalf("default config unexpectedly detected input: %+v", el)
		}
	}

	configured := Detect(g, WithMajorityThresholdPct(90))
	var input *Element
	for i := range configured.Elements {
		if configured.Elements[i].Type == ElementInput {
			input = &configured.Elements[i]
			break
		}
	}
	if input == nil || input.Bounds.Width != 8 {
		t.Fatalf("configured detect input = %+v, want width 8 input", input)
	}
}
