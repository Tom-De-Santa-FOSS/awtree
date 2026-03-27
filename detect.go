package awtree

// Detect analyzes a styled terminal grid and returns detected UI elements.
// Elements are assigned sequential IDs starting from 1.
func Detect(g *Grid) *ElementMap {
	if g == nil || g.Rows == 0 || g.Cols == 0 {
		return &ElementMap{}
	}

	var b elementBuilder

	// Detection order matches confidence ranking.
	panels := detectPanels(g)
	b.addAll(panels)

	b.addAll(detectButtons(g))
	menuItems := detectMenuItems(g)
	b.addAll(menuItems)
	b.addAll(detectInputs(g))

	tabs := detectTabs(g)
	b.addAll(tabs)

	// Standalone reverse regions — skip those already claimed by menus/tabs.
	for _, r := range detectReverseRegions(g) {
		if !overlapsAny(r, menuItems) && !overlapsAny(r, tabs) {
			b.add(r)
		}
	}

	b.addAll(detectStatusBars(g))

	return &ElementMap{Elements: b.elements}
}

// elementBuilder assigns sequential IDs and collects elements.
type elementBuilder struct {
	elements []Element
	nextID   int
}

func (b *elementBuilder) add(el Element) {
	if b.nextID == 0 {
		b.nextID = 1
	}
	el.ID = b.nextID
	b.nextID++
	b.elements = append(b.elements, el)
}

func (b *elementBuilder) addAll(els []Element) {
	for _, el := range els {
		b.add(el)
	}
}

func overlapsAny(el Element, others []Element) bool {
	for _, o := range others {
		if rectsOverlap(el.Bounds, o.Bounds) {
			return true
		}
	}
	return false
}

func rectsOverlap(a, b Rect) bool {
	return a.Row < b.Row+b.Height && a.Row+a.Height > b.Row &&
		a.Col < b.Col+b.Width && a.Col+a.Width > b.Col
}
