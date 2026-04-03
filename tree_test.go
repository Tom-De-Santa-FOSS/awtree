package awtree

import (
	"reflect"
	"testing"
)

func TestBuildTree_AssignsChildrenToTightestParent(t *testing.T) {
	elements := []Element{
		{ID: 1, Type: ElementPanel, Bounds: Rect{Row: 0, Col: 0, Width: 30, Height: 10}},
		{ID: 2, Type: ElementButton, Bounds: Rect{Row: 2, Col: 2, Width: 6, Height: 1}},
		{ID: 3, Type: ElementButton, Bounds: Rect{Row: 4, Col: 10, Width: 8, Height: 1}},
	}

	tree := BuildTree(elements)

	if got := tree[0].Children; !reflect.DeepEqual(got, []int{2, 3}) {
		t.Fatalf("children = %v, want [2 3]", got)
	}
	if len(tree[1].Children) != 0 || len(tree[2].Children) != 0 {
		t.Fatal("leaf elements should not have children")
	}
}

func TestBuildTree_NestedPanelsChooseDirectParent(t *testing.T) {
	elements := []Element{
		{ID: 1, Type: ElementPanel, Bounds: Rect{Row: 0, Col: 0, Width: 40, Height: 15}},
		{ID: 2, Type: ElementPanel, Bounds: Rect{Row: 2, Col: 2, Width: 20, Height: 8}},
		{ID: 3, Type: ElementButton, Bounds: Rect{Row: 4, Col: 4, Width: 6, Height: 1}},
	}

	tree := BuildTree(elements)

	if got := tree[0].Children; !reflect.DeepEqual(got, []int{2}) {
		t.Fatalf("outer children = %v, want [2]", got)
	}
	if got := tree[1].Children; !reflect.DeepEqual(got, []int{3}) {
		t.Fatalf("inner children = %v, want [3]", got)
	}
}

func TestBuildTree_PrefersDialogForEqualBounds(t *testing.T) {
	elements := []Element{
		{ID: 1, Type: ElementPanel, Bounds: Rect{Row: 8, Col: 30, Width: 21, Height: 7}},
		{ID: 2, Type: ElementDialog, Bounds: Rect{Row: 8, Col: 30, Width: 21, Height: 7}},
		{ID: 3, Type: ElementButton, Bounds: Rect{Row: 12, Col: 38, Width: 4, Height: 1}},
	}

	tree := BuildTree(elements)

	if len(tree[0].Children) != 0 {
		t.Fatalf("panel children = %v, want none", tree[0].Children)
	}
	if got := tree[1].Children; !reflect.DeepEqual(got, []int{3}) {
		t.Fatalf("dialog children = %v, want [3]", got)
	}
}
