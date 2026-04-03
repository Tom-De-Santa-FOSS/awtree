package awtree

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestSerialize_Empty(t *testing.T) {
	m := &ElementMap{}
	got := Serialize(m)
	if got != "(no elements detected)" {
		t.Errorf("got %q, want '(no elements detected)'", got)
	}
}

func TestSerialize_Nil(t *testing.T) {
	got := Serialize(nil)
	if got != "(no elements detected)" {
		t.Errorf("got %q", got)
	}
}

func TestSerialize_Button(t *testing.T) {
	m := &ElementMap{
		Elements: []Element{
			{
				ID: 1, Type: ElementButton, Label: "Save", Focused: true,
				Bounds: Rect{Row: 12, Col: 35, Width: 6, Height: 1},
			},
		},
	}
	got := Serialize(m)
	if !strings.Contains(got, `btn*`) {
		t.Errorf("expected focused btn marker, got %q", got)
	}
	if !strings.Contains(got, `"Save"`) {
		t.Errorf("expected label, got %q", got)
	}
}

func TestSerialize_Panel(t *testing.T) {
	m := &ElementMap{
		Elements: []Element{
			{
				ID: 1, Type: ElementPanel, Label: "Files",
				Bounds: Rect{Row: 0, Col: 0, Width: 40, Height: 20},
			},
		},
	}
	got := Serialize(m)
	if !strings.Contains(got, "panel") {
		t.Errorf("expected panel type, got %q", got)
	}
	if !strings.Contains(got, "40x20") {
		t.Errorf("expected WxH for multi-row, got %q", got)
	}
}

func TestSerializeJSON_returns_valid_json(t *testing.T) {
	m := &ElementMap{
		Elements: []Element{
			{ID: 1, Type: ElementPanel, Label: "Files", Bounds: Rect{0, 0, 40, 20}},
			{ID: 2, Type: ElementButton, Label: "Save", Focused: true, Bounds: Rect{12, 35, 6, 1}},
		},
	}
	got := SerializeJSON(m)

	var parsed map[string]any
	if err := json.Unmarshal([]byte(got), &parsed); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, got)
	}

	elements, ok := parsed["elements"].([]any)
	if !ok {
		t.Fatal("expected elements array")
	}
	if len(elements) != 2 {
		t.Fatalf("expected 2 elements, got %d", len(elements))
	}

	first := elements[0].(map[string]any)
	if first["type"] != "panel" {
		t.Errorf("first type = %v, want panel", first["type"])
	}
	if first["label"] != "Files" {
		t.Errorf("first label = %v, want Files", first["label"])
	}

	second := elements[1].(map[string]any)
	if second["focused"] != true {
		t.Errorf("second focused = %v, want true", second["focused"])
	}
}

func TestSerializeJSON_empty(t *testing.T) {
	got := SerializeJSON(&ElementMap{})
	var parsed map[string]any
	if err := json.Unmarshal([]byte(got), &parsed); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	elements := parsed["elements"].([]any)
	if len(elements) != 0 {
		t.Errorf("expected empty array, got %d elements", len(elements))
	}
}

func TestSerializeJSON_nil(t *testing.T) {
	got := SerializeJSON(nil)
	var parsed map[string]any
	if err := json.Unmarshal([]byte(got), &parsed); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
}

func TestSerialize_MultipleElements(t *testing.T) {
	m := &ElementMap{
		Elements: []Element{
			{ID: 1, Type: ElementPanel, Label: "Main", Bounds: Rect{0, 0, 80, 24}},
			{ID: 2, Type: ElementButton, Label: "OK", Bounds: Rect{22, 35, 4, 1}},
		},
	}
	got := Serialize(m)
	parts := strings.Split(got, " [")
	if len(parts) < 2 {
		t.Errorf("expected multiple elements separated by space, got %q", got)
	}
}
