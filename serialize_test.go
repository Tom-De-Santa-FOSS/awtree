package awtree

import (
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
			{ID: 1, Type: ElementButton, Label: "Save", Focused: true,
				Bounds: Rect{Row: 12, Col: 35, Width: 6, Height: 1}},
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
			{ID: 1, Type: ElementPanel, Label: "Files",
				Bounds: Rect{Row: 0, Col: 0, Width: 40, Height: 20}},
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
