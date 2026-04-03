package awtree

// ElementType classifies a detected TUI element.
type ElementType int

const (
	ElementPanel     ElementType = iota // Box-drawing bordered region
	ElementButton                       // Bracketed text like [Save], <OK>
	ElementInput                        // Cursor-adjacent editable field
	ElementMenuItem                     // Item in a vertical list/menu
	ElementStatusBar                    // Color-contiguous bar at screen edge
	ElementMenuBar                      // Horizontal menu at top
	ElementTab                          // Tab-bar label
	ElementText                         // Generic styled text region
	ElementCheckbox                     // Checkbox or radio button
	ElementProgressBar                  // Progress bar (block or ASCII)
	ElementTable                        // Tabular data with column separators
)

var elementTypeNames = [...]string{
	"panel",
	"button",
	"input",
	"menu_item",
	"status_bar",
	"menu_bar",
	"tab",
	"text",
	"checkbox",
	"progress_bar",
	"table",
}

func (t ElementType) String() string {
	if int(t) < len(elementTypeNames) {
		return elementTypeNames[t]
	}
	return "unknown"
}

var elementTypeShortNames = [...]string{
	"panel",
	"btn",
	"input",
	"item",
	"status",
	"menu",
	"tab",
	"text",
	"chk",
	"prog",
	"tbl",
}

func (t ElementType) ShortString() string {
	if int(t) < len(elementTypeShortNames) {
		return elementTypeShortNames[t]
	}
	return "?"
}

// Rect defines a rectangular region on the terminal grid.
type Rect struct {
	Row    int `json:"row"`
	Col    int `json:"col"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

// Element represents a detected interactive or structural TUI element.
type Element struct {
	ID       int         `json:"id"`
	Type     ElementType `json:"type"`
	Label    string      `json:"label"`
	Bounds   Rect        `json:"bounds"`
	Focused  bool        `json:"focused"`
	Children []int       `json:"children,omitempty"` // IDs of contained elements; populated by BuildTree
}

// ElementMap is the result of detecting elements on a styled grid.
type ElementMap struct {
	Elements []Element `json:"elements"`
}
