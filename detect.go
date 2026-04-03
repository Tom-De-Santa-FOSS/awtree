package awtree

// Detect analyzes a styled terminal grid and returns detected UI elements.
// Elements are assigned sequential IDs starting from 1.
func Detect(g *Grid, opts ...Option) *ElementMap {
	if g == nil || g.Rows == 0 || g.Cols == 0 {
		return &ElementMap{}
	}

	cfg := applyDetectOptions(opts)
	dbg := newDebugCollector(cfg)
	b := elementBuilder{nextID: 1, dbg: dbg}

	// Detection order matches confidence ranking.
	panels := detectPanels(g)
	b.addAll("panels", panels, "box-drawing border matched")

	b.addAll("tables", detectTables(g), "repeated column separators matched")
	buttons := detectButtons(g, cfg, dbg)
	b.addAll("buttons", buttons, "balanced button brackets matched")
	b.addAll("dialogs", detectDialogs(g, panels, buttons), "centered panel with actions matched")
	b.addAll("checkboxes", detectCheckboxes(g), "checkbox glyph matched")
	menuItems := detectMenuItems(g)
	b.addAll("menu_items", menuItems, "menu row pattern matched")
	b.addAll("inputs", detectInputs(g, cfg, dbg), "input field styling matched")
	b.addAll("progress_bars", detectProgressBars(g), "progress bar fill pattern matched")
	b.addAll("separators", detectSeparators(g), "separator line matched")
	b.addAll("scroll_indicators", detectScrollIndicators(g), "scroll affordance matched")

	tabs := detectTabs(g)
	b.addAll("tabs", tabs, "tab strip pattern matched")

	// Standalone reverse regions — skip those already claimed by menus/tabs.
	for _, r := range detectReverseRegions(g) {
		if !overlapsAny(r, menuItems) && !overlapsAny(r, tabs) {
			b.add("reverse_regions", r, "reverse-video span matched")
		}
	}

	b.addAll("breadcrumbs", detectBreadcrumbs(g), "breadcrumb separators matched")
	b.addAll("status_bars", detectStatusBars(g, cfg, dbg), "edge row background threshold matched")

	return buildElementMap(g, b.elements, dbg)
}

// elementBuilder assigns sequential IDs and collects elements.
type elementBuilder struct {
	elements []Element
	nextID   int
	dbg      *debugCollector
}

func (b *elementBuilder) add(detector string, el Element, reason string) {
	el.ID = b.nextID
	b.nextID++
	b.elements = append(b.elements, el)
	b.dbg.accept(detector, el, reason)
}

func (b *elementBuilder) addAll(detector string, els []Element, reason string) {
	for _, el := range els {
		b.add(detector, el, reason)
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
