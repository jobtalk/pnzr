package init

import "testing"

func TestIsMultiByteChar(t *testing.T) {
	if isMultiByteChar('a') {
		t.Errorf("I wanted %v but it was %v.", false, isMultiByteChar('a'))
	}

	if !isMultiByteChar('あ') {
		t.Errorf("I wanted %v but it was %v.", true, isMultiByteChar('あ'))
	}
}

func TestString2Lines(t *testing.T) {
	sampleStringAscii := "0123456789"
	sampleStringMultiByte := "あいうえお"

	if lines := string2Lines(sampleStringAscii, 0); len(lines) != 1 {
		t.Error(string2Lines(sampleStringAscii, 0))
	} else if lines[0].S != "0123456789" {
		t.Errorf("I wanted %v but it was %v.", "0123456789", lines[0].String())
	}

	if lines := string2Lines(sampleStringAscii, 5); len(lines) != 2 {
		t.Error(string2Lines(sampleStringAscii, 5))
	} else if lines[0].S != "01234" {
		t.Errorf("I wanted %v but it was %v.", "01234", lines[0].String())
	} else if lines[1].S != "56789" {
		t.Errorf("I wanted %v but it was %v.", "56789", lines[1].String())
	}

	if lines := string2Lines(sampleStringAscii, 4); len(lines) != 3 {
		t.Error(string2Lines(sampleStringAscii, 5))
	} else if lines[0].S != "0123" {
		t.Errorf("I wanted %v but it was %v.", "0123", lines[0].String())
	} else if lines[1].S != "4567" {
		t.Errorf("I wanted %v but it was %v.", "4567", lines[1].String())
	} else if lines[2].S != "89" {
		t.Errorf("I wanted %v but it was %v.", "89", lines[2].String())
	}

	if lines := string2Lines(sampleStringMultiByte, 0); len(lines) != 1 {
		t.Error(string2Lines(sampleStringMultiByte, 0))
	} else if lines[0].S != "あいうえお" {
		t.Errorf("I wanted %v but it was %v.", "あいうえお", lines[0].String())
	}

	if lines := string2Lines(sampleStringMultiByte, 5); len(lines) != 3 {
		t.Error(string2Lines(sampleStringMultiByte, 5))
	} else if lines[0].S != "あい" {
		t.Errorf("I wanted %v but it was %v.", "あい", lines[0].String())
	} else if lines[1].S != "うえ" {
		t.Errorf("I wanted %v but it was %v.", "うえ", lines[0].String())
	} else if lines[2].S != "お" {
		t.Errorf("I wanted %v but it was %v.", "お", lines[0].String())
	}
}

func TestLineString(t *testing.T) {
	line := NewLine("hoge")
	if line.String() != "  hoge" {
		t.Errorf(line.String())
	}

	line.SetCursor()
	if line.String() != "> hoge" {
		t.Errorf(line.String())
	}

	line.DeleteCursor()
	if line.String() != "  hoge" {
		t.Errorf(line.String())
	}

	line.ToggleCursor()
	if line.String() != "> hoge" {
		t.Errorf(line.String())
	}

	line.ToggleCursor()
	if line.String() != "  hoge" {
		t.Errorf(line.String())
	}
}

func TestItemString(t *testing.T) {
	if item := NewItem("test"); len(item.Lines) != 1 {
		t.Errorf("can not create item: %v", item)
	} else if item.SetCursor(); item.Lines[0].hasCursor == false {
		t.Errorf("can not set cursor: %v", item)
	} else if item.DeleteCursor(); item.Lines[0].hasCursor == true {
		t.Errorf("can not delete cursor: %v", item)
	} else if item.ToggleCursor(); item.Lines[0].hasCursor == false {
		t.Errorf("can not delete cursor: %v", item)
	} else if item.ToggleCursor(); item.Lines[0].hasCursor == true {
		t.Errorf("can not delete cursor: %v", item)
	}

	if item := NewItem("0123456789", 5); len(item.Lines) != 2 {
		t.Errorf("can not create item: %v", item)
	} else if item.SetCursor(); item.Lines[0].hasCursor == false {
		t.Errorf("can not set cursor: %v", item)
	} else if item.Lines[1].hasCursor == true {
		t.Errorf("can not set cursor: %v", item)
	} else if item.DeleteCursor(); item.Lines[0].hasCursor == true {
		t.Errorf("can not set cursor: %v", item)
	} else if item.Lines[1].hasCursor == true {
		t.Errorf("can not set cursor: %v", item)
	} else if item.ToggleCursor(); item.Lines[0].hasCursor == false {
		t.Errorf("can not set cursor: %v", item)
	} else if item.Lines[1].hasCursor == true {
		t.Errorf("can not set cursor: %v", item)
	} else if item.ToggleCursor(); item.Lines[0].hasCursor == true {
		t.Errorf("can not set cursor: %v", item)
	} else if item.Lines[1].hasCursor == true {
		t.Errorf("can not set cursor: %v", item)
	}
}
