package init

import "fmt"

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
	Question    string
	cursorPlace int
	Items       []*Item
	MaxWidth    int
}

func (s *SelectBox) ToggleCursor() {
}

func (s *SelectBox) Answer() string {
	return fmt.Sprintf("%v", s.cursorPlace)
}

func (s *SelectBox) Message() string {
	return s.String()
}

func (s SelectBox) String() string {
	var ret = fmt.Sprintf("<-- %v -->\n", s.Question)

	for _, item := range s.Items {
		ret += item.String()
	}
	return ret
}

// カーソルを移動する
func (s *SelectBox) Up() {
	if s.cursorPlace == 0 {
		return
	}
	s.Items[s.cursorPlace].ToggleCursor()
	s.cursorPlace--
	s.Items[s.cursorPlace].ToggleCursor()
}

func (s *SelectBox) Down() {
	if s.cursorPlace == len(s.Items)-1 {
		return
	}
	s.Items[s.cursorPlace].ToggleCursor()
	s.cursorPlace++
	s.Items[s.cursorPlace].ToggleCursor()
}

func NewSelectBox(q string, elems []string, w ...int) *SelectBox {
	var items = []*Item{}
	var ret = &SelectBox{}
	if len(elems) == 0 {
		return nil
	}
	ret.Question = q

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

func genCursor(b bool) string {
	if b {
		return "_"
	}
	return ""
}

type TextBox struct {
	Question string
	input    string
	MaxWidth int
	cursor   bool
}

func (t *TextBox) ToggleCursor() {
	t.cursor = t.cursor != true
}

func (t *TextBox) Answer() string {
	return t.input
}

func (t *TextBox) Message() string {
	var ret string
	ret = fmt.Sprintf("<-- %v -->\n", t.Question)
	ret = ret + fmt.Sprintf("=> %v%v\n", t.input, genCursor(t.cursor))
	return ret
}

func (t *TextBox) Subst(s string) {
	t.input = s
}

func (t *TextBox) Add(r rune) {
	t.cursor = true
	if 0x20 <= uint8(r) && uint8(r) <= 0x7f {
		t.input += string([]rune{r})
	}
}

func (t *TextBox) BS() {
	t.cursor = true
	if len(t.input) == 0 {
		return
	}
	runes := []rune(t.input)
	runes = runes[:len(runes)-1]
	t.input = string(runes)
}

func NewTextBox(q string, w ...int) *TextBox {
	return &TextBox{
		Question: q,
	}
}

type PolarQuestionBox struct {
	Question string
	input    string
	defo     string
	cursor   bool
}

func (t *PolarQuestionBox) ToggleCursor() {
	t.cursor = t.cursor != true
}

func (t *PolarQuestionBox) Add(r rune) {
	t.cursor = true
	if 0x20 <= uint8(r) && uint8(r) <= 0x7f {
		t.input += string([]rune{r})
	}
}

func (t *PolarQuestionBox) BS() {
	t.cursor = true
	if len(t.input) == 0 {
		return
	}
	runes := []rune(t.input)
	runes = runes[:len(runes)-1]
	t.input = string(runes)
}

func (t *PolarQuestionBox) Answer() string {
	if t.input == "yes" ||
		t.input == "y" ||
		t.input == "Y" ||
		t.input == "Yes" ||
		t.input == "YES" {
		return "y"
	} else if t.input == "n" ||
		t.input == "N" ||
		t.input == "no" ||
		t.input == "No" ||
		t.input == "NO" {
		return "n"
	} else if t.input == "" {
		return ""
	}
	return t.input
}

func (t *PolarQuestionBox) Message() string {
	var ret string
	var def string
	if t.defo != "" {
		def = t.defo
	}
	ret = fmt.Sprintf("<-- %v -->\n", t.Question)
	ret = ret + fmt.Sprintf("(Y/n) %v%v%v\n", def, t.input, genCursor(t.cursor))
	return ret
}

func NewPolarQuestionBox(q string, d ...string) *PolarQuestionBox {
	var def string
	if len(d) != 0 {
		def = d[0]
	}
	return &PolarQuestionBox{
		Question: q,
		defo:     def,
	}
}
