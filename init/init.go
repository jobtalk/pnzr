package init

func isMultiByteChar(r rune) bool {
	c := string([]rune{r})
	return 1 != len(c)
}

func string2Lines(s string, w int) []*Line {
	var ret = []*Line{}
	runes := []rune(s)
	wCounter := 0
	var stringBuffer string
	if w == 0 {
		return []*Line{NewLine(s)}
	}

	for _, r := range runes {
		if isMultiByteChar(r) {
			if wCounter+2 <= w {
				wCounter += 2
				stringBuffer += string([]rune{r})
			} else {
				l := NewLine(stringBuffer)
				ret = append(ret, l)
				wCounter = 2
				stringBuffer = string([]rune{r})
			}
		} else {
			if wCounter+1 <= w {
				wCounter += 1
				stringBuffer += string([]rune{r})
			} else {
				l := NewLine(stringBuffer)
				ret = append(ret, l)
				wCounter = 1
				stringBuffer = string([]rune{r})
			}
		}
	}
	ret = append(ret, NewLine(stringBuffer))
	return ret
}

type Line struct {
	hasCursor bool
	S         string
}

func NewLine(s string) *Line {
	return &Line{S: s}
}

func (l Line) String() string {
	if l.hasCursor {
		return "> " + l.S
	}
	return "  " + l.S
}

func (l *Line) SetCursor() {
	l.hasCursor = true
}

func (l *Line) DeleteCursor() {
	l.hasCursor = false
}

func (l *Line) ToggleCursor() {
	l.hasCursor = true != l.hasCursor
}

type Item struct {
	Lines []*Line
}

func (i Item) String() string {
	var ret string
	for _, line := range i.Lines {
		ret = ret + line.String() + "\n"
	}
	return ret
}

func (i *Item) SetCursor() {
	i.Lines[0].SetCursor()
}

func (i *Item) DeleteCursor() {
	i.Lines[0].DeleteCursor()
}

func (i *Item) ToggleCursor() {
	i.Lines[0].ToggleCursor()
}

func NewItem(s string, w ...int) *Item {
	if len(w) == 0 {
		return &Item{
			[]*Line{NewLine(s)},
		}
	}

	return &Item{string2Lines(s, w[0])}
}

type SelectBox struct {
	cursorPlace int
	Items       []*Item
	MaxWidth    int
}

func NewSelectBox(elems []string, w ...int) *SelectBox {
	var items = []*Item{}
	var ret = &SelectBox{}
	if len(elems) == 0 {
		return nil
	}

	if len(w) == 0 {
		ret.MaxWidth = 0
	} else {
		ret.MaxWidth = w[0]
	}
	ret.cursorPlace = 0

	for _, elem := range elems {
		items = append(items, NewItem(elem, w...))
	}
	ret.Items = items
	ret.Items[0].SetCursor()
	return ret
}
