# awtree

Structured element detection for terminal UIs. Takes a styled character grid and produces a labeled element map — the terminal equivalent of a browser's accessibility tree.

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

## Detection

| Signal | Detects |
|--------|---------|
| Box-drawing chars | Panels, windows, dialogs |
| Reverse-video regions | Focused/selected elements |
| Edge rows with distinct BG | Status bars, menu bars |
| Bracketed text `[Save]` | Buttons |

## License

MIT
