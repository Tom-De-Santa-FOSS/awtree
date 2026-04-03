package awtree

import "testing"

func TestDetect_NilGrid(t *testing.T) {
	m := Detect(nil)
	if len(m.Elements) != 0 {
		t.Fatalf("expected empty, got %d elements", len(m.Elements))
	}
}

func TestDetect_EmptyGrid(t *testing.T) {
	g := NewGrid(24, 80)
	m := Detect(g)
	if len(m.Elements) != 0 {
		t.Fatalf("expected empty, got %d elements", len(m.Elements))
	}
}

func TestDetect_Panel(t *testing.T) {
	g := NewGrid(10, 20)
	// Draw a simple box.
	g.Set(0, 0, Cell{Char: '┌'})
	for c := 1; c < 9; c++ {
		g.Set(0, c, Cell{Char: '─'})
	}
	g.Set(0, 9, Cell{Char: '┐'})

	for r := 1; r < 4; r++ {
		g.Set(r, 0, Cell{Char: '│'})
		g.Set(r, 9, Cell{Char: '│'})
	}

	g.Set(4, 0, Cell{Char: '└'})
	for c := 1; c < 9; c++ {
		g.Set(4, c, Cell{Char: '─'})
	}
	g.Set(4, 9, Cell{Char: '┘'})

	m := Detect(g)

	found := false
	for _, el := range m.Elements {
		if el.Type == ElementPanel {
			found = true
			if el.Bounds.Width != 10 || el.Bounds.Height != 5 {
				t.Errorf("panel bounds = %dx%d, want 10x5", el.Bounds.Width, el.Bounds.Height)
			}
		}
	}
	if !found {
		t.Fatal("no panel detected")
	}
}

func TestDetect_PanelWithTitle(t *testing.T) {
	g := NewGrid(10, 30)
	// ┌─ Title ─┐
	g.Set(0, 0, Cell{Char: '┌'})
	g.Set(0, 1, Cell{Char: '─'})
	g.Set(0, 2, Cell{Char: ' '})
	g.SetText(0, 3, "Title", DefaultColor, DefaultColor, 0)
	g.Set(0, 8, Cell{Char: ' '})
	g.Set(0, 9, Cell{Char: '─'})
	g.Set(0, 10, Cell{Char: '┐'})

	for r := 1; r < 4; r++ {
		g.Set(r, 0, Cell{Char: '│'})
		g.Set(r, 10, Cell{Char: '│'})
	}

	g.Set(4, 0, Cell{Char: '└'})
	for c := 1; c < 10; c++ {
		g.Set(4, c, Cell{Char: '─'})
	}
	g.Set(4, 10, Cell{Char: '┘'})

	m := Detect(g)

	for _, el := range m.Elements {
		if el.Type == ElementPanel {
			if el.Label != "Title" {
				t.Errorf("panel label = %q, want %q", el.Label, "Title")
			}
			return
		}
	}
	t.Fatal("no panel detected")
}

func TestDetect_Button(t *testing.T) {
	g := NewGrid(5, 40)
	g.SetText(2, 10, "[Save]", DefaultColor, DefaultColor, 0)
	g.SetText(2, 20, "<Cancel>", DefaultColor, DefaultColor, 0)

	m := Detect(g)

	labels := make(map[string]bool)
	for _, el := range m.Elements {
		if el.Type == ElementButton {
			labels[el.Label] = true
		}
	}

	if !labels["Save"] {
		t.Error("Save button not detected")
	}
	if !labels["Cancel"] {
		t.Error("Cancel button not detected")
	}
}

func TestDetect_FocusedButton(t *testing.T) {
	g := NewGrid(5, 40)
	g.SetText(2, 10, "[OK]", DefaultColor, DefaultColor, AttrReverse)

	m := Detect(g)

	for _, el := range m.Elements {
		if el.Type == ElementButton && el.Label == "OK" {
			if !el.Focused {
				t.Error("expected focused button")
			}
			return
		}
	}
	t.Fatal("focused button not detected")
}

func TestDetect_ReverseRegion(t *testing.T) {
	g := NewGrid(10, 40)
	g.SetText(3, 2, "  Selected Item  ", DefaultColor, 4, AttrReverse)

	m := Detect(g)

	found := false
	for _, el := range m.Elements {
		if el.Focused && el.Label == "Selected Item" {
			found = true
		}
	}
	if !found {
		t.Fatal("reverse-video region not detected")
	}
}

func TestDetect_StatusBar(t *testing.T) {
	g := NewGrid(24, 80)
	// Fill bottom row with colored background.
	for c := 0; c < 80; c++ {
		g.Set(23, c, Cell{Char: ' ', BG: 4})
	}
	g.SetText(23, 2, "status: ready", DefaultColor, 4, 0)

	m := Detect(g)

	found := false
	for _, el := range m.Elements {
		if el.Type == ElementStatusBar {
			found = true
		}
	}
	if !found {
		t.Fatal("status bar not detected")
	}
}

func TestDetect_MenuItems_VerticalListWithOneHighlighted(t *testing.T) {
	g := NewGrid(10, 30)
	// 3 items at same column, one reverse-video.
	g.SetText(2, 2, "  Open File  ", DefaultColor, DefaultColor, 0)
	g.SetText(3, 2, "  Save File  ", DefaultColor, DefaultColor, AttrReverse)
	g.SetText(4, 2, "  Close All  ", DefaultColor, DefaultColor, 0)

	m := Detect(g)

	var menuItems []Element
	for _, el := range m.Elements {
		if el.Type == ElementMenuItem {
			menuItems = append(menuItems, el)
		}
	}

	if len(menuItems) != 3 {
		t.Fatalf("expected exactly 3 menu items, got %d", len(menuItems))
	}

	focusCount := 0
	for _, el := range menuItems {
		if el.Focused {
			focusCount++
			if el.Label != "Save File" {
				t.Errorf("focused item label = %q, want %q", el.Label, "Save File")
			}
		}
	}
	if focusCount != 1 {
		t.Errorf("expected 1 focused item, got %d", focusCount)
	}
}

func TestDetect_TabBar(t *testing.T) {
	g := NewGrid(10, 60)
	// Horizontal tabs: one bold/reverse (active), others normal.
	g.SetText(0, 0, " Files ", DefaultColor, 4, AttrReverse|AttrBold)
	g.SetText(0, 7, " Edit ", DefaultColor, DefaultColor, 0)
	g.SetText(0, 13, " View ", DefaultColor, DefaultColor, 0)

	m := Detect(g)

	var tabs []Element
	for _, el := range m.Elements {
		if el.Type == ElementTab {
			tabs = append(tabs, el)
		}
	}

	if len(tabs) < 3 {
		t.Fatalf("expected 3 tabs, got %d", len(tabs))
	}

	focusCount := 0
	for _, el := range tabs {
		if el.Focused {
			focusCount++
			if el.Label != "Files" {
				t.Errorf("active tab = %q, want %q", el.Label, "Files")
			}
		}
	}
	if focusCount != 1 {
		t.Errorf("expected 1 focused tab, got %d", focusCount)
	}
}

func TestDetect_InputField_UnderscorePlaceholder(t *testing.T) {
	g := NewGrid(10, 40)
	// Label followed by underscores (common input pattern).
	g.SetText(3, 2, "Name: ", DefaultColor, DefaultColor, 0)
	g.SetText(3, 8, "________________", DefaultColor, DefaultColor, AttrUnderline)

	m := Detect(g)

	found := false
	for _, el := range m.Elements {
		if el.Type == ElementInput {
			found = true
			if el.Bounds.Row != 3 {
				t.Errorf("input row = %d, want 3", el.Bounds.Row)
			}
		}
	}
	if !found {
		t.Fatal("input field not detected")
	}
}

func TestDetect_InputField_DistinctBackground(t *testing.T) {
	g := NewGrid(10, 40)
	// Input field: region with distinct BG color.
	g.SetText(5, 10, "              ", DefaultColor, 7, 0) // white bg, empty
	g.SetText(5, 10, "hello", DefaultColor, 7, 0)

	m := Detect(g)

	found := false
	for _, el := range m.Elements {
		if el.Type == ElementInput {
			found = true
			if el.Label != "hello" {
				t.Errorf("input value = %q, want %q", el.Label, "hello")
			}
		}
	}
	if !found {
		t.Fatal("input field with distinct BG not detected")
	}
}

func TestDetect_BracketedNonButton_ArrayNotation(t *testing.T) {
	g := NewGrid(5, 40)
	// [0] and [123] should NOT be buttons (start with digit).
	g.SetText(1, 5, "[0]", DefaultColor, DefaultColor, 0)
	g.SetText(2, 5, "[123]", DefaultColor, DefaultColor, 0)

	m := Detect(g)

	for _, el := range m.Elements {
		if el.Type == ElementButton {
			t.Errorf("false positive button: %q", el.Label)
		}
	}
}

func TestDetect_NestedPanels(t *testing.T) {
	g := NewGrid(12, 30)
	// Outer panel.
	g.Set(0, 0, Cell{Char: '┌'})
	for c := 1; c < 19; c++ {
		g.Set(0, c, Cell{Char: '─'})
	}
	g.Set(0, 19, Cell{Char: '┐'})
	for r := 1; r < 9; r++ {
		g.Set(r, 0, Cell{Char: '│'})
		g.Set(r, 19, Cell{Char: '│'})
	}
	g.Set(9, 0, Cell{Char: '└'})
	for c := 1; c < 19; c++ {
		g.Set(9, c, Cell{Char: '─'})
	}
	g.Set(9, 19, Cell{Char: '┘'})

	// Inner panel.
	g.Set(2, 2, Cell{Char: '┌'})
	for c := 3; c < 9; c++ {
		g.Set(2, c, Cell{Char: '─'})
	}
	g.Set(2, 9, Cell{Char: '┐'})
	for r := 3; r < 6; r++ {
		g.Set(r, 2, Cell{Char: '│'})
		g.Set(r, 9, Cell{Char: '│'})
	}
	g.Set(6, 2, Cell{Char: '└'})
	for c := 3; c < 9; c++ {
		g.Set(6, c, Cell{Char: '─'})
	}
	g.Set(6, 9, Cell{Char: '┘'})

	m := Detect(g)

	panelCount := 0
	for _, el := range m.Elements {
		if el.Type == ElementPanel {
			panelCount++
		}
	}
	if panelCount != 2 {
		t.Errorf("expected 2 panels (nested), got %d", panelCount)
	}
}

func TestDetect_MenuBar_TopRowWithBG(t *testing.T) {
	g := NewGrid(24, 80)
	// Top row with colored BG → menu bar, not status bar.
	for c := 0; c < 80; c++ {
		g.Set(0, c, Cell{Char: ' ', BG: 2})
	}
	g.SetText(0, 2, "File  Edit  View  Help", DefaultColor, 2, 0)

	m := Detect(g)

	found := false
	for _, el := range m.Elements {
		if el.Type == ElementMenuBar {
			found = true
		}
		if el.Type == ElementStatusBar && el.Bounds.Row == 0 {
			t.Error("top row should be MenuBar, not StatusBar")
		}
	}
	if !found {
		t.Fatal("menu bar not detected on top row")
	}
}

func TestDetect_RealisticTUI(t *testing.T) {
	g := NewGrid(24, 80)

	// Menu bar (row 0, colored BG).
	for c := 0; c < 80; c++ {
		g.Set(0, c, Cell{Char: ' ', BG: 4})
	}
	g.SetText(0, 1, "File  Edit  View", DefaultColor, 4, 0)

	// Panel (rows 1-20).
	g.Set(1, 0, Cell{Char: '┌'})
	for c := 1; c < 59; c++ {
		g.Set(1, c, Cell{Char: '─'})
	}
	g.Set(1, 59, Cell{Char: '┐'})
	for r := 2; r < 20; r++ {
		g.Set(r, 0, Cell{Char: '│'})
		g.Set(r, 59, Cell{Char: '│'})
	}
	g.Set(20, 0, Cell{Char: '└'})
	for c := 1; c < 59; c++ {
		g.Set(20, c, Cell{Char: '─'})
	}
	g.Set(20, 59, Cell{Char: '┘'})

	// Buttons on row 21.
	g.SetText(21, 30, "[Save]", DefaultColor, DefaultColor, 0)
	g.SetText(21, 40, "[Cancel]", DefaultColor, DefaultColor, 0)

	// Status bar (bottom row).
	for c := 0; c < 80; c++ {
		g.Set(23, c, Cell{Char: ' ', BG: 2})
	}
	g.SetText(23, 2, "Ready", DefaultColor, 2, 0)

	m := Detect(g)

	types := make(map[ElementType]int)
	for _, el := range m.Elements {
		types[el.Type]++
	}

	if types[ElementPanel] < 1 {
		t.Error("missing panel")
	}
	if types[ElementButton] < 2 {
		t.Errorf("expected 2+ buttons, got %d", types[ElementButton])
	}
	if types[ElementMenuBar] < 1 {
		t.Error("missing menu bar")
	}
	if types[ElementStatusBar] < 1 {
		t.Error("missing status bar")
	}

	// All IDs unique and > 0.
	ids := make(map[int]bool)
	for _, el := range m.Elements {
		if el.ID == 0 {
			t.Error("ID 0 found")
		}
		if ids[el.ID] {
			t.Errorf("duplicate ID %d", el.ID)
		}
		ids[el.ID] = true
	}
}

func TestDetect_MenuItems_NoDuplicatesWithMultipleHighlights(t *testing.T) {
	g := NewGrid(10, 30)
	// Two highlighted items in same column range with shared siblings.
	g.SetText(2, 2, "  Open File  ", DefaultColor, DefaultColor, 0)
	g.SetText(3, 2, "  Save File  ", DefaultColor, DefaultColor, AttrReverse)
	g.SetText(4, 2, "  Close All  ", DefaultColor, DefaultColor, AttrReverse)
	g.SetText(5, 2, "  Quit       ", DefaultColor, DefaultColor, 0)

	m := Detect(g)

	var menuItems []Element
	for _, el := range m.Elements {
		if el.Type == ElementMenuItem {
			menuItems = append(menuItems, el)
		}
	}

	if len(menuItems) != 4 {
		t.Fatalf("expected exactly 4 menu items, got %d", len(menuItems))
	}

	// Check no duplicate rows.
	seen := make(map[int]bool)
	for _, el := range menuItems {
		if seen[el.Bounds.Row] {
			t.Errorf("duplicate menu item at row %d", el.Bounds.Row)
		}
		seen[el.Bounds.Row] = true
	}
}

func TestOverlapsAny_BoundingBoxIntersection(t *testing.T) {
	// A reverse region starting inside a menu item's bounds should overlap.
	menuItem := Element{
		Type:   ElementMenuItem,
		Bounds: Rect{Row: 3, Col: 2, Width: 13, Height: 1},
	}
	// Reverse region starts at col 5 (inside menu item col 2..14).
	reverseRegion := Element{
		Bounds: Rect{Row: 3, Col: 5, Width: 8, Height: 1},
	}

	if !overlapsAny(reverseRegion, []Element{menuItem}) {
		t.Error("expected overlapsAny to return true for bounding-box intersection")
	}
}

func TestOverlapsAny_NoOverlap(t *testing.T) {
	menuItem := Element{
		Bounds: Rect{Row: 3, Col: 2, Width: 13, Height: 1},
	}
	other := Element{
		Bounds: Rect{Row: 5, Col: 2, Width: 13, Height: 1},
	}

	if overlapsAny(other, []Element{menuItem}) {
		t.Error("expected no overlap for different rows")
	}
}

func TestIsButtonLabel(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"", false},
		{"Save", true},
		{"OK", true},
		{"Cancel", true},
		{"0", false},   // starts with digit
		{"123", false}, // starts with digit
		{"a very long button label that exceeds twenty chars", false},
		{string([]rune{0x01}), false}, // non-printable
	}
	for _, tt := range tests {
		got := isButtonLabel(tt.input, DefaultDetectConfig().MaxButtonLabelLen)
		if got != tt.want {
			t.Errorf("isButtonLabel(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestTrimSpaces(t *testing.T) {
	tests := []struct {
		input []rune
		want  string
	}{
		{[]rune("  hello  "), "hello"},
		{[]rune{}, ""},
		{[]rune("   "), ""},
		{[]rune("hello"), "hello"},
		{[]rune("  hello"), "hello"},
		{[]rune("hello  "), "hello"},
	}
	for _, tt := range tests {
		got := trimSpaces(tt.input)
		if got != tt.want {
			t.Errorf("trimSpaces(%q) = %q, want %q", string(tt.input), got, tt.want)
		}
	}
}

func TestExtractLineText(t *testing.T) {
	g := NewGrid(5, 20)
	g.SetText(2, 5, "Hello", DefaultColor, DefaultColor, 0)

	got := extractLineText(g, 2, 5, 5)
	if got != "Hello" {
		t.Errorf("extractLineText = %q, want %q", got, "Hello")
	}

	// Blank region returns empty string.
	got = extractLineText(g, 3, 5, 5)
	if got != "" {
		t.Errorf("extractLineText blank = %q, want empty", got)
	}
}

func TestDetect_FirstElementIDIsOne(t *testing.T) {
	g := NewGrid(5, 20)
	g.SetText(2, 5, "[OK]", DefaultColor, DefaultColor, 0)

	m := Detect(g)

	if len(m.Elements) == 0 {
		t.Fatal("expected at least one element")
	}
	if m.Elements[0].ID != 1 {
		t.Errorf("first element ID = %d, want 1", m.Elements[0].ID)
	}
}

func TestDetect_InputField_TooShortUnderlinedSpan_NotDetected(t *testing.T) {
	g := NewGrid(5, 20)
	// Width=2, below minimum threshold of 3.
	g.Set(2, 5, Cell{Char: '_', Attrs: AttrUnderline})
	g.Set(2, 6, Cell{Char: '_', Attrs: AttrUnderline})

	m := Detect(g)
	for _, el := range m.Elements {
		if el.Type == ElementInput {
			t.Error("short underlined span should not be detected as input")
		}
	}
}

func TestDetect_InputField_WideDistinctBGSpan_NotDetectedAsInput(t *testing.T) {
	g := NewGrid(10, 20) // 20 cols, threshold is 12
	for c := 0; c < 16; c++ {
		g.Set(5, c, Cell{Char: ' ', BG: 7})
	}

	m := Detect(g)
	for _, el := range m.Elements {
		if el.Type == ElementInput {
			t.Error("wide BG span should not be classified as input")
		}
	}
}

func TestDetect_RowWithTwoActiveSegments_NotDetectedAsTabs(t *testing.T) {
	g := NewGrid(5, 40)
	g.SetText(0, 0, " Files ", DefaultColor, 4, AttrReverse)
	g.SetText(0, 10, " Edit ", DefaultColor, 4, AttrReverse)
	g.SetText(0, 20, " View ", DefaultColor, DefaultColor, 0)

	m := Detect(g)
	for _, el := range m.Elements {
		if el.Type == ElementTab {
			t.Error("row with 2 active segments should not produce tabs")
		}
	}
}

func TestDetect_RoundedCornerPanel(t *testing.T) {
	g := NewGrid(5, 12)
	// Lazygit-style rounded corners: ╭╮╰╯
	g.Set(0, 0, Cell{Char: '╭'})
	for c := 1; c < 9; c++ {
		g.Set(0, c, Cell{Char: '─'})
	}
	g.Set(0, 9, Cell{Char: '╮'})
	for r := 1; r < 4; r++ {
		g.Set(r, 0, Cell{Char: '│'})
		g.Set(r, 9, Cell{Char: '│'})
	}
	g.Set(4, 0, Cell{Char: '╰'})
	for c := 1; c < 9; c++ {
		g.Set(4, c, Cell{Char: '─'})
	}
	g.Set(4, 9, Cell{Char: '╯'})

	m := Detect(g)

	found := false
	for _, el := range m.Elements {
		if el.Type == ElementPanel {
			found = true
			if el.Bounds.Width != 10 || el.Bounds.Height != 5 {
				t.Errorf("panel bounds = %dx%d, want 10x5", el.Bounds.Width, el.Bounds.Height)
			}
		}
	}
	if !found {
		t.Error("rounded corner panel not detected")
	}
}

func TestDetect_SequentialIDs(t *testing.T) {
	g := NewGrid(10, 40)
	// Panel.
	g.Set(0, 0, Cell{Char: '┌'})
	g.Set(0, 1, Cell{Char: '─'})
	g.Set(0, 2, Cell{Char: '┐'})
	g.Set(1, 0, Cell{Char: '│'})
	g.Set(1, 2, Cell{Char: '│'})
	g.Set(2, 0, Cell{Char: '└'})
	g.Set(2, 1, Cell{Char: '─'})
	g.Set(2, 2, Cell{Char: '┘'})

	// Button.
	g.SetText(5, 10, "[OK]", DefaultColor, DefaultColor, 0)

	m := Detect(g)

	ids := make(map[int]bool)
	for _, el := range m.Elements {
		if el.ID == 0 {
			t.Error("element has ID 0, should start at 1")
		}
		if ids[el.ID] {
			t.Errorf("duplicate ID %d", el.ID)
		}
		ids[el.ID] = true
	}
}

// --- Checkbox/Radio tests ---

func TestDetect_Checkbox_Unicode(t *testing.T) {
	g := NewGrid(5, 40)
	g.Set(1, 2, Cell{Char: '☐', FG: DefaultColor, BG: DefaultColor})
	g.SetText(1, 3, " Accept terms", DefaultColor, DefaultColor, 0)
	g.Set(2, 2, Cell{Char: '☑', FG: DefaultColor, BG: DefaultColor})
	g.SetText(2, 3, " Send newsletter", DefaultColor, DefaultColor, 0)

	m := Detect(g)
	var checkboxes []Element
	for _, el := range m.Elements {
		if el.Type == ElementCheckbox {
			checkboxes = append(checkboxes, el)
		}
	}
	if len(checkboxes) != 2 {
		t.Fatalf("expected 2 checkboxes, got %d", len(checkboxes))
	}
	if checkboxes[0].Label != "☐ Accept terms" {
		t.Errorf("checkbox[0] label = %q, want %q", checkboxes[0].Label, "☐ Accept terms")
	}
	if checkboxes[1].Label != "☑ Send newsletter" {
		t.Errorf("checkbox[1] label = %q, want %q", checkboxes[1].Label, "☑ Send newsletter")
	}
}

func TestDetect_Checkbox_ASCII(t *testing.T) {
	g := NewGrid(5, 40)
	g.SetText(0, 0, "[x] Enable logging", DefaultColor, DefaultColor, 0)
	g.SetText(1, 0, "[ ] Verbose mode", DefaultColor, DefaultColor, 0)
	g.SetText(2, 0, "[X] Dark theme", DefaultColor, DefaultColor, 0)
	g.SetText(3, 0, "[*] Auto-save", DefaultColor, DefaultColor, 0)

	m := Detect(g)
	var checkboxes []Element
	for _, el := range m.Elements {
		if el.Type == ElementCheckbox {
			checkboxes = append(checkboxes, el)
		}
	}
	if len(checkboxes) != 4 {
		t.Fatalf("expected 4 checkboxes, got %d", len(checkboxes))
	}
	wantLabels := []string{"[x] Enable logging", "[ ] Verbose mode", "[X] Dark theme", "[*] Auto-save"}
	for i, want := range wantLabels {
		if checkboxes[i].Label != want {
			t.Errorf("checkbox[%d] label = %q, want %q", i, checkboxes[i].Label, want)
		}
	}
}

func TestDetect_Radio_ASCII(t *testing.T) {
	g := NewGrid(5, 40)
	g.SetText(0, 0, "( ) Option A", DefaultColor, DefaultColor, 0)
	g.SetText(1, 0, "(x) Option B", DefaultColor, DefaultColor, 0)
	g.SetText(2, 0, "(*) Option C", DefaultColor, DefaultColor, 0)

	m := Detect(g)
	var checkboxes []Element
	for _, el := range m.Elements {
		if el.Type == ElementCheckbox {
			checkboxes = append(checkboxes, el)
		}
	}
	if len(checkboxes) != 3 {
		t.Fatalf("expected 3 radio buttons, got %d", len(checkboxes))
	}
	wantLabels := []string{"( ) Option A", "(x) Option B", "(*) Option C"}
	for i, want := range wantLabels {
		if checkboxes[i].Label != want {
			t.Errorf("radio[%d] label = %q, want %q", i, checkboxes[i].Label, want)
		}
	}
}

func TestDetect_Checkbox_Focused(t *testing.T) {
	g := NewGrid(5, 40)
	g.Set(1, 2, Cell{Char: '☑', FG: DefaultColor, BG: DefaultColor, Attrs: AttrReverse})
	g.SetText(1, 3, " Agree", DefaultColor, DefaultColor, AttrReverse)
	g.SetText(2, 2, "[x]", DefaultColor, DefaultColor, AttrReverse)
	g.SetText(2, 5, " Confirm", DefaultColor, DefaultColor, AttrReverse)

	m := Detect(g)
	var checkboxes []Element
	for _, el := range m.Elements {
		if el.Type == ElementCheckbox {
			checkboxes = append(checkboxes, el)
		}
	}
	if len(checkboxes) != 2 {
		t.Fatalf("expected 2 checkboxes, got %d", len(checkboxes))
	}
	for i, cb := range checkboxes {
		if !cb.Focused {
			t.Errorf("checkbox[%d] should be focused", i)
		}
	}
}

// --- ProgressBar tests ---

func TestDetect_ProgressBar_BlockChars(t *testing.T) {
	g := NewGrid(5, 40)
	for c := 5; c < 9; c++ {
		g.Set(2, c, Cell{Char: '█'})
	}
	for c := 9; c < 13; c++ {
		g.Set(2, c, Cell{Char: '░'})
	}

	m := Detect(g)
	found := false
	for _, el := range m.Elements {
		if el.Type == ElementProgressBar {
			found = true
			if el.Bounds.Row != 2 || el.Bounds.Col != 5 || el.Bounds.Width != 8 || el.Bounds.Height != 1 {
				t.Errorf("bounds = %+v, want {Row:2 Col:5 Width:8 Height:1}", el.Bounds)
			}
		}
	}
	if !found {
		t.Fatal("block character progress bar not detected")
	}
}

func TestDetect_ProgressBar_ASCII(t *testing.T) {
	g := NewGrid(5, 40)
	g.SetText(1, 3, "[====    ]", DefaultColor, DefaultColor, 0)

	m := Detect(g)
	found := false
	for _, el := range m.Elements {
		if el.Type == ElementProgressBar {
			found = true
			if el.Bounds.Row != 1 || el.Bounds.Col != 3 || el.Bounds.Width != 10 || el.Bounds.Height != 1 {
				t.Errorf("bounds = %+v, want {Row:1 Col:3 Width:10 Height:1}", el.Bounds)
			}
		}
	}
	if !found {
		t.Fatal("ASCII progress bar not detected")
	}
}

func TestDetect_ProgressBar_WithPercentage(t *testing.T) {
	g := NewGrid(5, 40)
	for c := 5; c < 10; c++ {
		g.Set(2, c, Cell{Char: '█'})
	}
	for c := 10; c < 15; c++ {
		g.Set(2, c, Cell{Char: '░'})
	}
	g.SetText(2, 16, "45%", DefaultColor, DefaultColor, 0)

	m := Detect(g)
	found := false
	for _, el := range m.Elements {
		if el.Type == ElementProgressBar {
			found = true
			if el.Label != "45%" {
				t.Errorf("label = %q, want %q", el.Label, "45%")
			}
		}
	}
	if !found {
		t.Fatal("progress bar with percentage not detected")
	}
}

// --- Table tests ---

func TestDetect_Table_Simple(t *testing.T) {
	g := NewGrid(5, 20)
	g.SetText(0, 0, " Alice ", DefaultColor, DefaultColor, 0)
	g.Set(0, 7, Cell{Char: '│'})
	g.SetText(0, 8, " 30", DefaultColor, DefaultColor, 0)
	g.SetText(1, 0, " Bob   ", DefaultColor, DefaultColor, 0)
	g.Set(1, 7, Cell{Char: '│'})
	g.SetText(1, 8, " 25", DefaultColor, DefaultColor, 0)
	g.SetText(2, 0, " Carol ", DefaultColor, DefaultColor, 0)
	g.Set(2, 7, Cell{Char: '│'})
	g.SetText(2, 8, " 40", DefaultColor, DefaultColor, 0)

	m := Detect(g)
	found := false
	for _, el := range m.Elements {
		if el.Type == ElementTable {
			found = true
			if el.Label != "3x2" {
				t.Errorf("table label = %q, want %q", el.Label, "3x2")
			}
			if el.Bounds.Height != 3 {
				t.Errorf("table height = %d, want 3", el.Bounds.Height)
			}
		}
	}
	if !found {
		t.Fatal("no table detected")
	}
}

func TestDetect_Table_WithHeader(t *testing.T) {
	g := NewGrid(6, 20)
	g.SetText(0, 0, " Name  ", DefaultColor, DefaultColor, AttrBold)
	g.Set(0, 7, Cell{Char: '│', Attrs: AttrBold})
	g.SetText(0, 8, " Age", DefaultColor, DefaultColor, AttrBold)
	for c := 0; c < 14; c++ {
		g.Set(1, c, Cell{Char: '─'})
	}
	g.Set(1, 7, Cell{Char: '┼'})
	g.SetText(2, 0, " Alice ", DefaultColor, DefaultColor, 0)
	g.Set(2, 7, Cell{Char: '│'})
	g.SetText(2, 8, " 30", DefaultColor, DefaultColor, 0)
	g.SetText(3, 0, " Bob   ", DefaultColor, DefaultColor, 0)
	g.Set(3, 7, Cell{Char: '│'})
	g.SetText(3, 8, " 25", DefaultColor, DefaultColor, 0)

	m := Detect(g)
	found := false
	for _, el := range m.Elements {
		if el.Type == ElementTable {
			found = true
			if el.Label != "3x2" {
				t.Errorf("table label = %q, want %q", el.Label, "3x2")
			}
			if el.Bounds.Height != 4 {
				t.Errorf("table height = %d, want 4", el.Bounds.Height)
			}
		}
	}
	if !found {
		t.Fatal("no table with header detected")
	}
}

func TestDetect_Table_ASCII(t *testing.T) {
	g := NewGrid(7, 25)
	g.SetText(0, 0, " Name  ", DefaultColor, DefaultColor, 0)
	g.Set(0, 7, Cell{Char: '|'})
	g.SetText(0, 8, " Age ", DefaultColor, DefaultColor, 0)
	g.Set(0, 13, Cell{Char: '|'})
	g.SetText(0, 14, " City", DefaultColor, DefaultColor, 0)
	for c := 0; c < 20; c++ {
		g.Set(1, c, Cell{Char: '-'})
	}
	g.Set(1, 7, Cell{Char: '+'})
	g.Set(1, 13, Cell{Char: '+'})
	for r := 2; r <= 4; r++ {
		g.SetText(r, 0, " Name  ", DefaultColor, DefaultColor, 0)
		g.Set(r, 7, Cell{Char: '|'})
		g.SetText(r, 8, " val ", DefaultColor, DefaultColor, 0)
		g.Set(r, 13, Cell{Char: '|'})
		g.SetText(r, 14, " City", DefaultColor, DefaultColor, 0)
	}

	m := Detect(g)
	found := false
	for _, el := range m.Elements {
		if el.Type == ElementTable {
			found = true
			if el.Label != "4x3" {
				t.Errorf("table label = %q, want %q", el.Label, "4x3")
			}
			if el.Bounds.Height != 5 {
				t.Errorf("table height = %d, want 5", el.Bounds.Height)
			}
		}
	}
	if !found {
		t.Fatal("no ASCII table detected")
	}
}

func TestDetect_Table_PanelNotFalsePositive(t *testing.T) {
	g := NewGrid(6, 12)
	g.Set(0, 0, Cell{Char: '┌'})
	for c := 1; c < 9; c++ {
		g.Set(0, c, Cell{Char: '─'})
	}
	g.Set(0, 9, Cell{Char: '┐'})
	for r := 1; r < 4; r++ {
		g.Set(r, 0, Cell{Char: '│'})
		g.Set(r, 9, Cell{Char: '│'})
	}
	g.Set(4, 0, Cell{Char: '└'})
	for c := 1; c < 9; c++ {
		g.Set(4, c, Cell{Char: '─'})
	}
	g.Set(4, 9, Cell{Char: '┘'})

	m := Detect(g)
	for _, el := range m.Elements {
		if el.Type == ElementTable {
			t.Errorf("panel falsely detected as table: %+v", el)
		}
	}
}

// --- Separator tests ---

func TestDetect_Separator_BasicDash(t *testing.T) {
	g := NewGrid(5, 20)
	for c := 3; c < 13; c++ {
		g.Set(2, c, Cell{Char: '─'})
	}

	m := Detect(g)

	found := false
	for _, el := range m.Elements {
		if el.Type == ElementSeparator {
			found = true
			if el.Bounds.Row != 2 {
				t.Errorf("separator row = %d, want 2", el.Bounds.Row)
			}
			if el.Bounds.Col != 3 {
				t.Errorf("separator col = %d, want 3", el.Bounds.Col)
			}
			if el.Bounds.Width != 10 {
				t.Errorf("separator width = %d, want 10", el.Bounds.Width)
			}
			if el.Bounds.Height != 1 {
				t.Errorf("separator height = %d, want 1", el.Bounds.Height)
			}
			if el.Label != "" {
				t.Errorf("separator label = %q, want empty", el.Label)
			}
		}
	}
	if !found {
		t.Fatal("separator not detected")
	}
}

func TestDetect_Separator_ShortDash(t *testing.T) {
	g := NewGrid(5, 20)
	g.Set(2, 5, Cell{Char: '-'})
	g.Set(2, 6, Cell{Char: '-'})

	m := Detect(g)

	for _, el := range m.Elements {
		if el.Type == ElementSeparator {
			t.Error("2-char dash run should not be detected as separator")
		}
	}
}

func TestDetect_Separator_PanelBorderNotSeparator(t *testing.T) {
	g := NewGrid(5, 20)
	g.Set(2, 2, Cell{Char: '┌'})
	for c := 3; c < 13; c++ {
		g.Set(2, c, Cell{Char: '─'})
	}
	g.Set(2, 13, Cell{Char: '┐'})

	m := Detect(g)

	for _, el := range m.Elements {
		if el.Type == ElementSeparator {
			t.Error("panel border row should not be detected as separator")
		}
	}
}

// --- Dialog tests ---

func TestDetect_Dialog_CenteredPanelWithButton(t *testing.T) {
	g := NewGrid(24, 80)
	topRow, botRow := 8, 14
	leftCol, rightCol := 30, 50
	g.Set(topRow, leftCol, Cell{Char: '┌'})
	g.Set(topRow, rightCol, Cell{Char: '┐'})
	g.Set(botRow, leftCol, Cell{Char: '└'})
	g.Set(botRow, rightCol, Cell{Char: '┘'})
	for c := leftCol + 1; c < rightCol; c++ {
		g.Set(topRow, c, Cell{Char: '─'})
		g.Set(botRow, c, Cell{Char: '─'})
	}
	for r := topRow + 1; r < botRow; r++ {
		g.Set(r, leftCol, Cell{Char: '│'})
		g.Set(r, rightCol, Cell{Char: '│'})
	}
	g.SetText(topRow, leftCol+2, "─ Confirm ─", DefaultColor, DefaultColor, 0)
	g.SetText(12, 38, "[OK]", DefaultColor, DefaultColor, 0)

	m := Detect(g)

	found := false
	for _, el := range m.Elements {
		if el.Type == ElementDialog {
			found = true
			if el.Bounds.Row != topRow || el.Bounds.Col != leftCol {
				t.Errorf("dialog bounds start = (%d,%d), want (%d,%d)", el.Bounds.Row, el.Bounds.Col, topRow, leftCol)
			}
		}
	}
	if !found {
		t.Fatal("dialog not detected for centered panel with button")
	}
}

func TestDetect_Dialog_CenteredPanelNoButtons(t *testing.T) {
	g := NewGrid(24, 80)
	topRow, botRow := 8, 14
	leftCol, rightCol := 30, 50
	g.Set(topRow, leftCol, Cell{Char: '┌'})
	g.Set(topRow, rightCol, Cell{Char: '┐'})
	g.Set(botRow, leftCol, Cell{Char: '└'})
	g.Set(botRow, rightCol, Cell{Char: '┘'})
	for c := leftCol + 1; c < rightCol; c++ {
		g.Set(topRow, c, Cell{Char: '─'})
		g.Set(botRow, c, Cell{Char: '─'})
	}
	for r := topRow + 1; r < botRow; r++ {
		g.Set(r, leftCol, Cell{Char: '│'})
		g.Set(r, rightCol, Cell{Char: '│'})
	}
	g.SetText(10, 33, "Loading...", DefaultColor, DefaultColor, 0)

	m := Detect(g)

	for _, el := range m.Elements {
		if el.Type == ElementDialog {
			t.Fatal("centered panel without buttons should not be detected as dialog")
		}
	}
}

func TestDetect_Dialog_OffCenterPanelNotDialog(t *testing.T) {
	g := NewGrid(24, 80)
	topRow, botRow := 8, 14
	leftCol, rightCol := 0, 20
	g.Set(topRow, leftCol, Cell{Char: '┌'})
	g.Set(topRow, rightCol, Cell{Char: '┐'})
	g.Set(botRow, leftCol, Cell{Char: '└'})
	g.Set(botRow, rightCol, Cell{Char: '┘'})
	for c := leftCol + 1; c < rightCol; c++ {
		g.Set(topRow, c, Cell{Char: '─'})
		g.Set(botRow, c, Cell{Char: '─'})
	}
	for r := topRow + 1; r < botRow; r++ {
		g.Set(r, leftCol, Cell{Char: '│'})
		g.Set(r, rightCol, Cell{Char: '│'})
	}
	g.SetText(12, 5, "[OK]", DefaultColor, DefaultColor, 0)

	m := Detect(g)

	for _, el := range m.Elements {
		if el.Type == ElementDialog {
			t.Fatal("off-center panel should not be detected as dialog")
		}
	}
}

// --- ScrollIndicator tests ---

func TestDetect_ScrollIndicator_ArrowPair(t *testing.T) {
	g := NewGrid(10, 20)
	g.Set(0, 10, Cell{Char: '▲'})
	g.Set(5, 10, Cell{Char: '▼'})

	m := Detect(g)

	found := false
	for _, el := range m.Elements {
		if el.Type == ElementScrollIndicator {
			found = true
			if el.Bounds.Col != 10 {
				t.Errorf("scroll indicator col = %d, want 10", el.Bounds.Col)
			}
			if el.Bounds.Row != 0 {
				t.Errorf("scroll indicator row = %d, want 0", el.Bounds.Row)
			}
			if el.Bounds.Height != 6 {
				t.Errorf("scroll indicator height = %d, want 6", el.Bounds.Height)
			}
			if el.Bounds.Width != 1 {
				t.Errorf("scroll indicator width = %d, want 1", el.Bounds.Width)
			}
		}
	}
	if !found {
		t.Fatal("scroll indicator (arrow pair) not detected")
	}
}

func TestDetect_ScrollIndicator_BlockThumb(t *testing.T) {
	g := NewGrid(10, 20)
	g.Set(2, 15, Cell{Char: '█'})
	g.Set(3, 15, Cell{Char: '█'})
	g.Set(4, 15, Cell{Char: '█'})
	g.Set(5, 15, Cell{Char: '█'})

	m := Detect(g)

	found := false
	for _, el := range m.Elements {
		if el.Type == ElementScrollIndicator {
			found = true
			if el.Bounds.Col != 15 {
				t.Errorf("scroll indicator col = %d, want 15", el.Bounds.Col)
			}
			if el.Bounds.Row != 2 {
				t.Errorf("scroll indicator row = %d, want 2", el.Bounds.Row)
			}
			if el.Bounds.Height != 4 {
				t.Errorf("scroll indicator height = %d, want 4", el.Bounds.Height)
			}
		}
	}
	if !found {
		t.Fatal("scroll indicator (block thumb) not detected")
	}
}

func TestDetect_ScrollIndicator_TooShort(t *testing.T) {
	g := NewGrid(10, 20)
	g.Set(3, 12, Cell{Char: '█'})
	g.Set(4, 12, Cell{Char: '█'})

	m := Detect(g)

	for _, el := range m.Elements {
		if el.Type == ElementScrollIndicator {
			t.Errorf("short block run (2 chars) should not be detected as scroll indicator")
		}
	}
}

// --- Breadcrumb tests ---

func TestDetect_Breadcrumb_AngleBracket(t *testing.T) {
	g := NewGrid(5, 80)
	g.SetText(1, 2, "Home > Settings > Display", DefaultColor, DefaultColor, 0)

	m := Detect(g)

	found := false
	for _, el := range m.Elements {
		if el.Type == ElementBreadcrumb {
			found = true
			if el.Label != "Home > Settings > Display" {
				t.Errorf("breadcrumb label = %q, want %q", el.Label, "Home > Settings > Display")
			}
			if el.Bounds.Row != 1 || el.Bounds.Col != 2 {
				t.Errorf("breadcrumb position = (%d,%d), want (1,2)", el.Bounds.Row, el.Bounds.Col)
			}
			if el.Bounds.Width != 25 || el.Bounds.Height != 1 {
				t.Errorf("breadcrumb size = %dx%d, want 25x1", el.Bounds.Width, el.Bounds.Height)
			}
		}
	}
	if !found {
		t.Fatal("breadcrumb not detected")
	}
}

func TestDetect_Breadcrumb_Slash(t *testing.T) {
	g := NewGrid(5, 80)
	g.SetText(0, 0, "File / Edit / View / Help", DefaultColor, DefaultColor, 0)

	m := Detect(g)

	found := false
	for _, el := range m.Elements {
		if el.Type == ElementBreadcrumb {
			found = true
			if el.Label != "File / Edit / View / Help" {
				t.Errorf("breadcrumb label = %q, want %q", el.Label, "File / Edit / View / Help")
			}
		}
	}
	if !found {
		t.Fatal("slash-separated breadcrumb not detected")
	}
}

func TestDetect_Breadcrumb_TooFewSegments(t *testing.T) {
	g := NewGrid(5, 80)
	g.SetText(1, 2, "Home > Settings", DefaultColor, DefaultColor, 0)

	m := Detect(g)

	for _, el := range m.Elements {
		if el.Type == ElementBreadcrumb {
			t.Errorf("breadcrumb should not be detected with only 2 segments, got label %q", el.Label)
		}
	}
}
