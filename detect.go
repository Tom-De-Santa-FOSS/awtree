package awtree

// Detect analyzes a styled terminal grid and returns detected UI elements.
// Elements are assigned sequential IDs starting from 1.
func Detect(g *Grid) *ElementMap {
	if g == nil || g.Rows == 0 || g.Cols == 0 {
		return &ElementMap{}
	}

	var elements []Element
	nextID := 1

	// Detection order matches confidence ranking:
	// 1. Panels (box-drawing boundaries)
	// 2. Reverse-video regions (focused/selected elements)
	// 3. Status/menu bars (edge rows with distinct styling)
	// 4. Buttons (bracketed text patterns)
	// 5. Menu items (repeated vertical structure with highlight)
	// 6. Input fields (cursor-adjacent editable areas)

	panels := detectPanels(g)
	for i := range panels {
		panels[i].ID = nextID
		nextID++
	}
	elements = append(elements, panels...)

	buttons := detectButtons(g)
	for i := range buttons {
		buttons[i].ID = nextID
		nextID++
	}
	elements = append(elements, buttons...)

	menuItems := detectMenuItems(g)
	for i := range menuItems {
		menuItems[i].ID = nextID
		nextID++
	}
	elements = append(elements, menuItems...)

	// Standalone reverse regions (not part of a menu).
	reverseRegions := detectReverseRegions(g)
	for i := range reverseRegions {
		// Skip if already covered by a menu item at same position.
		if overlapsAny(reverseRegions[i], menuItems) {
			continue
		}
		reverseRegions[i].ID = nextID
		nextID++
		elements = append(elements, reverseRegions[i])
	}

	statusBars := detectStatusBars(g)
	for i := range statusBars {
		statusBars[i].ID = nextID
		nextID++
	}
	elements = append(elements, statusBars...)

	return &ElementMap{Elements: elements}
}

func overlapsAny(el Element, others []Element) bool {
	for _, o := range others {
		if el.Bounds.Row == o.Bounds.Row && el.Bounds.Col == o.Bounds.Col {
			return true
		}
	}
	return false
}
