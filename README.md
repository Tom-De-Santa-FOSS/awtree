# awtree

Structured element detection for terminal UIs. Takes a styled character grid (characters + ANSI attributes), produces a labeled element map — the terminal equivalent of a browser's accessibility tree.

## What it does

```
Raw terminal grid          →  awtree.Detect()  →  Structured elements
┌─ Files ─┐                                        [1:panel:"Files" 0,0 20x10]
│ doc.txt  │                                        [2:item:"doc.txt"* 1,2]
│ main.go  │                                        [3:btn:"Open" 11,5 w6]
└──────────┘                                        [4:status:"ready" 23,0 w80]
[Open] [Cancel]
```

## Usage

```go
// Build a styled grid (or convert from your terminal emulator).
g := awtree.NewGrid(24, 80)
g.SetText(5, 10, "[Save]", awtree.DefaultColor, awtree.DefaultColor, awtree.AttrReverse)

// Detect elements.
elements := awtree.Detect(g)

// Serialize for LLM consumption (~3 tokens per element).
fmt.Println(awtree.Serialize(elements))
// [1:btn*:"Save" 5,10 w6]
```

## Detection heuristics

| Priority | Signal | Detects |
|----------|--------|---------|
| 1 | Box-drawing chars (┌─┐│└┘) | Panels, windows, dialogs |
| 2 | Reverse-video regions | Focused/selected elements |
| 3 | Edge rows with distinct BG | Status bars, menu bars |
| 4 | Bracketed text ([Save], \<OK\>) | Buttons |

## Install

```
go get github.com/Tom-De-Santa-FOSS/awtree
```

## License

MIT
