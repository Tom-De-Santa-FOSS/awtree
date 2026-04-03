package awtree

// detectDialogs finds centered panels that contain at least one button,
// treating them as modal dialog boxes.
func detectDialogs(g *Grid, panels []Element, buttons []Element) []Element {
	gridCenterCol := g.Cols / 2
	var results []Element

	for _, p := range panels {
		panelCenterCol := p.Bounds.Col + p.Bounds.Width/2
		if abs(panelCenterCol-gridCenterCol) > 3 {
			continue
		}

		if !containsButton(p, buttons) {
			continue
		}

		results = append(results, Element{
			Type:   ElementDialog,
			Label:  p.Label,
			Bounds: p.Bounds,
		})
	}
	return results
}

func containsButton(panel Element, buttons []Element) bool {
	for _, b := range buttons {
		if boundsContain(panel.Bounds, b.Bounds) {
			return true
		}
	}
	return false
}

func boundsContain(outer, inner Rect) bool {
	return inner.Row >= outer.Row &&
		inner.Col >= outer.Col &&
		inner.Row+inner.Height <= outer.Row+outer.Height &&
		inner.Col+inner.Width <= outer.Col+outer.Width
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
