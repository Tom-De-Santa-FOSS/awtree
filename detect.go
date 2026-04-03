package awtree

const (
	// majorityThresholdPct is the percentage of row width above which a span is
	// considered a bar (status/menu), not an input field.
	majorityThresholdPct = 60

	// maxButtonWidth is the maximum scan width for a button bracket pair.
	maxButtonWidth = 30

	// maxButtonLabelLen is the maximum allowed button label length.
	maxButtonLabelLen = 20
)

// Detect analyzes a styled terminal grid and returns detected UI elements.
// Elements are assigned sequential IDs starting from 1.
func Detect(g *Grid) *ElementMap {
	if g == nil || g.Rows == 0 || g.Cols == 0 {
		return &ElementMap{}
	}

	b := elementBuilder{nextID: 1}

	// Detection order matches confidence ranking.
	panels := detectPanels(g)
	b.addAll(panels)

	b.addAll(detectTables(g))
	b.addAll(detectButtons(g))
	b.addAll(detectCheckboxes(g))
	menuItems := detectMenuItems(g)
	b.addAll(menuItems)
	b.addAll(detectInputs(g))
	b.addAll(detectProgressBars(g))

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
