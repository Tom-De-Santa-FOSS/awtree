package awtree

import (
	"encoding/json"
	"testing"
)

func TestSerializeJSON_Phase4TreeAndDescriptions(t *testing.T) {
	m := &ElementMap{
		Viewport: Rect{Row: 0, Col: 0, Width: 40, Height: 10},
		Debug:    &DebugInfo{Config: DefaultDetectConfig(), Events: []DebugEvent{{Detector: "buttons", Accepted: true, Reason: "matched"}}},
		Elements: BuildTree([]Element{
			{ID: 1, Type: ElementPanel, Label: "Main", Description: "form \"Main\"", Bounds: Rect{Row: 0, Col: 0, Width: 20, Height: 8}},
			{ID: 2, Type: ElementButton, Label: "Save", Description: "button \"Save\" focused", Focused: true, Bounds: Rect{Row: 1, Col: 2, Width: 6, Height: 1}},
		}),
	}

	var parsed map[string]any
	if err := json.Unmarshal([]byte(SerializeJSON(m)), &parsed); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	tree := parsed["tree"].([]any)
	if len(tree) != 1 {
		t.Fatalf("tree len = %d, want 1", len(tree))
	}
	root := tree[0].(map[string]any)
	if root["description"] != `form "Main"` {
		t.Fatalf("root description = %v, want form description", root["description"])
	}
	children := root["children"].([]any)
	if len(children) != 1 || children[0].(map[string]any)["label"] != "Save" {
		t.Fatalf("tree children = %+v, want nested Save button", children)
	}
	debug := parsed["debug"].(map[string]any)
	if debug["config"] == nil {
		t.Fatal("expected debug config in JSON")
	}
}

func TestGrid_SetText_DoubleWidthCharacters(t *testing.T) {
	g := NewGrid(2, 10)
	g.SetText(0, 0, "界🙂A", DefaultColor, DefaultColor, 0)

	if got := g.At(0, 0); got.Char != '界' || got.Width != 2 {
		t.Fatalf("first cell = %+v, want wide CJK rune", got)
	}
	if got := g.At(0, 1); !got.Continuation {
		t.Fatalf("second cell = %+v, want continuation", got)
	}
	if got := g.At(0, 2); got.Char != '🙂' || got.Width != 2 {
		t.Fatalf("third cell = %+v, want wide emoji rune", got)
	}
	if got := g.At(0, 4).Char; got != 'A' {
		t.Fatalf("cell[0,4] = %q, want A", got)
	}
}

func TestColor_RGBColorRoundTrip(t *testing.T) {
	color := RGBColor(0x12, 0x34, 0x56)
	if !color.IsRGB() {
		t.Fatal("expected RGB color")
	}
	r, g, b, ok := color.RGB()
	if !ok || r != 0x12 || g != 0x34 || b != 0x56 {
		t.Fatalf("RGB() = %02x %02x %02x %v, want 12 34 56 true", r, g, b, ok)
	}
	if PaletteColor(4).IsRGB() {
		t.Fatal("palette color should not report RGB")
	}
}
