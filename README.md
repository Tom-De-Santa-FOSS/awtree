# awtree

Structured element detection for terminal UIs. Takes a styled character grid and produces a labeled element map — the terminal equivalent of a browser's accessibility tree.

## Highlights

- Detects 15 element types: panels, buttons, inputs, menus, tabs, tables, dialogs, and more
- CSS-like query engine: `panel > button[label="Save"]:focused`
- Compact token-efficient serialization for LLM consumption
- ARIA-like roles, shortcuts, focus/checked state, and tree hierarchy
- Used by [awn](https://github.com/Tom-De-Santa-FOSS/awn) for AI-driven TUI automation

## Install

```
go get github.com/Tom-De-Santa-FOSS/awtree
```

## Usage

```go
g := awtree.NewGrid(24, 80)
g.SetText(5, 10, "[Save]", awtree.DefaultColor, awtree.DefaultColor, awtree.AttrReverse)

elements := awtree.Detect(g)
fmt.Println(awtree.Serialize(elements))
// [1:btn*:"Save" 5,10 w6]
```

### Querying

```go
results := elements.Query(`panel > button[label="Save"]:focused`)
results  = elements.Query(`checkbox:checked`)
results  = elements.Query(`menu_item:nth(2)`)
```

Supports type, ID (`#5`), attribute (`[label="Save"]`), pseudo-class (`:focused`, `:checked`, `:disabled`, `:selected`), descendant/child combinators, and `:nth()`.

### Output Formats

**Compact text** — token-efficient for LLMs:

```
[1:panel:"File Browser" 0,0 40x20] [2:btn*:"Save" 12,35 6x1]
```

**JSON** — structured with flattened elements and hierarchical tree:

```go
fmt.Println(awtree.SerializeJSON(elements))
// {"elements":[...],"tree":[...],"viewport":{...}}
```

## Detection

| Signal | Detects |
|--------|---------|
| Box-drawing chars | Panels, windows, dialogs |
| Reverse-video regions | Focused/selected elements |
| Edge rows with distinct BG | Status bars, menu bars |
| Bracketed text `[Save]` | Buttons |
| Cursor-adjacent fields | Inputs |
| Checkbox/radio glyphs | Checkboxes |
| Vertical item lists | Menu items |
| Tab-bar labels | Tabs |
| Column separators | Tables |
| Fill patterns | Progress bars |
| Centered panel + buttons | Dialogs |
| Horizontal rules | Separators |
| Arrow/block affordances | Scroll indicators |
| Path-like trails | Breadcrumbs |

## License

MIT
