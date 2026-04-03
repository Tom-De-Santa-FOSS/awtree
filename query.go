package awtree

import (
	"strconv"
	"strings"
	"unicode"
)

type selectorStep struct {
	combinator string
	typeName   string
	id         int
	attrs      []selectorAttr
	pseudos    []string
	nth        int
}

type selectorAttr struct {
	name     string
	contains bool
	value    string
}

func (m *ElementMap) Query(selector string) []Element {
	indices := m.queryIndices(selector)
	if len(indices) == 0 {
		return []Element{}
	}
	results := make([]Element, 0, len(indices))
	for _, idx := range indices {
		results = append(results, m.Elements[idx])
	}
	return results
}

func (m *ElementMap) QueryOne(selector string) *Element {
	indices := m.queryIndices(selector)
	if len(indices) == 0 {
		return nil
	}
	return &m.Elements[indices[0]]
}

func (m *ElementMap) queryIndices(selector string) []int {
	if m == nil || strings.TrimSpace(selector) == "" {
		return nil
	}
	steps := parseSelector(selector)
	if len(steps) == 0 {
		return nil
	}
	parentByID := make(map[int]int)
	for _, el := range m.Elements {
		for _, childID := range el.Children {
			parentByID[childID] = el.ID
		}
	}

	var current []int
	for stepIndex, step := range steps {
		var next []int
		prevSet := make(map[int]bool, len(current))
		for _, idx := range current {
			prevSet[m.Elements[idx].ID] = true
		}

		for i, el := range m.Elements {
			if !matchesStep(el, step) {
				continue
			}
			if stepIndex == 0 || satisfiesCombinator(el.ID, step.combinator, prevSet, parentByID) {
				next = append(next, i)
			}
		}

		if step.nth > 0 {
			if step.nth <= len(next) {
				next = []int{next[step.nth-1]}
			} else {
				next = nil
			}
		}
		current = next
	}

	if current == nil {
		return []int{}
	}
	return current
}

func parseSelector(selector string) []selectorStep {
	var steps []selectorStep
	var buf strings.Builder
	bracketDepth := 0
	parenDepth := 0
	pendingCombinator := ""

	flush := func() {
		text := strings.TrimSpace(buf.String())
		buf.Reset()
		if text == "" {
			return
		}
		step := parseSelectorStep(text)
		step.combinator = pendingCombinator
		steps = append(steps, step)
		pendingCombinator = ""
	}

	for i := 0; i < len(selector); i++ {
		ch := selector[i]
		switch ch {
		case '[':
			bracketDepth++
		case ']':
			if bracketDepth > 0 {
				bracketDepth--
			}
		case '(':
			parenDepth++
		case ')':
			if parenDepth > 0 {
				parenDepth--
			}
		case '>':
			if bracketDepth == 0 && parenDepth == 0 {
				flush()
				pendingCombinator = "child"
				continue
			}
		}

		if unicode.IsSpace(rune(ch)) && bracketDepth == 0 && parenDepth == 0 {
			flush()
			if len(steps) > 0 {
				pendingCombinator = "descendant"
			}
			continue
		}

		buf.WriteByte(ch)
	}
	flush()
	return steps
}

func parseSelectorStep(text string) selectorStep {
	step := selectorStep{id: -1}
	i := 0
	for i < len(text) && (unicode.IsLetter(rune(text[i])) || text[i] == '_') {
		i++
	}
	step.typeName = text[:i]

	for i < len(text) {
		switch text[i] {
		case '#':
			j := i + 1
			for j < len(text) && unicode.IsDigit(rune(text[j])) {
				j++
			}
			if id, err := strconv.Atoi(text[i+1 : j]); err == nil {
				step.id = id
			}
			i = j
		case '[':
			j := i + 1
			inQuotes := false
			for j < len(text) {
				if text[j] == '"' {
					inQuotes = !inQuotes
				} else if text[j] == ']' && !inQuotes {
					break
				}
				j++
			}
			if j <= len(text) {
				step.attrs = append(step.attrs, parseSelectorAttr(text[i+1:j]))
			}
			i = j + 1
		case ':':
			j := i + 1
			for j < len(text) && (unicode.IsLetter(rune(text[j])) || text[j] == '_') {
				j++
			}
			name := text[i+1 : j]
			if name == "nth" && j < len(text) && text[j] == '(' {
				k := j + 1
				for k < len(text) && text[k] != ')' {
					k++
				}
				if n, err := strconv.Atoi(text[j+1 : k]); err == nil {
					step.nth = n
				}
				i = k + 1
				continue
			}
			step.pseudos = append(step.pseudos, name)
			i = j
		default:
			i++
		}
	}

	return step
}

func parseSelectorAttr(text string) selectorAttr {
	if strings.Contains(text, "~=") {
		parts := strings.SplitN(text, "~=", 2)
		return selectorAttr{name: strings.TrimSpace(parts[0]), contains: true, value: trimSelectorValue(parts[1])}
	}
	parts := strings.SplitN(text, "=", 2)
	if len(parts) != 2 {
		return selectorAttr{}
	}
	return selectorAttr{name: strings.TrimSpace(parts[0]), value: trimSelectorValue(parts[1])}
}

func trimSelectorValue(value string) string {
	value = strings.TrimSpace(value)
	value = strings.TrimPrefix(value, `"`)
	value = strings.TrimSuffix(value, `"`)
	return value
}

func matchesStep(el Element, step selectorStep) bool {
	if step.typeName != "" && el.Type.String() != step.typeName {
		return false
	}
	if step.id >= 0 && el.ID != step.id {
		return false
	}
	for _, attr := range step.attrs {
		var actual string
		switch attr.name {
		case "label":
			actual = el.Label
		case "ref":
			actual = el.Ref
		default:
			return false
		}
		if attr.contains {
			if !strings.Contains(actual, attr.value) {
				return false
			}
		} else if actual != attr.value {
			return false
		}
	}
	for _, pseudo := range step.pseudos {
		switch pseudo {
		case "checked":
			if !el.Checked {
				return false
			}
		case "focused":
			if !el.Focused {
				return false
			}
		case "disabled":
			if el.Enabled {
				return false
			}
		case "selected":
			if !el.Selected {
				return false
			}
		default:
			return false
		}
	}
	return true
}

func satisfiesCombinator(id int, combinator string, prevSet map[int]bool, parentByID map[int]int) bool {
	switch combinator {
	case "child":
		return prevSet[parentByID[id]]
	case "descendant":
		for parentID := parentByID[id]; parentID != 0; parentID = parentByID[parentID] {
			if prevSet[parentID] {
				return true
			}
		}
		return false
	default:
		return true
	}
}
