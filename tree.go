package awtree

// BuildTree populates Children by assigning each element to the tightest
// container that fully contains it.
func BuildTree(elements []Element) []Element {
	if len(elements) == 0 {
		return nil
	}

	tree := make([]Element, len(elements))
	copy(tree, elements)

	for i := range tree {
		tree[i].Children = nil
	}

	parents := make([]int, len(tree))
	for i := range parents {
		parents[i] = -1
	}

	for childIdx := range tree {
		bestParent := -1
		bestArea := 0
		bestPriority := 0

		for parentIdx := range tree {
			if parentIdx == childIdx || !canContainChildren(tree[parentIdx].Type) {
				continue
			}
			if !strictlyContains(tree[parentIdx].Bounds, tree[childIdx].Bounds) {
				continue
			}

			area := rectArea(tree[parentIdx].Bounds)
			priority := containerPriority(tree[parentIdx].Type)
			if bestParent == -1 || area < bestArea ||
				(area == bestArea && priority > bestPriority) ||
				(area == bestArea && priority == bestPriority && tree[parentIdx].ID < tree[bestParent].ID) {
				bestParent = parentIdx
				bestArea = area
				bestPriority = priority
			}
		}

		parents[childIdx] = bestParent
	}

	for childIdx, parentIdx := range parents {
		if parentIdx >= 0 {
			tree[parentIdx].Children = append(tree[parentIdx].Children, tree[childIdx].ID)
		}
	}

	return tree
}

func canContainChildren(t ElementType) bool {
	switch t {
	case ElementPanel, ElementDialog, ElementMenuBar:
		return true
	default:
		return false
	}
}

func containerPriority(t ElementType) int {
	switch t {
	case ElementDialog:
		return 2
	case ElementPanel, ElementMenuBar:
		return 1
	default:
		return 0
	}
}

func strictlyContains(outer, inner Rect) bool {
	return boundsContain(outer, inner) && !sameRect(outer, inner)
}

func sameRect(a, b Rect) bool {
	return a.Row == b.Row && a.Col == b.Col && a.Width == b.Width && a.Height == b.Height
}

func rectArea(r Rect) int {
	return r.Width * r.Height
}
