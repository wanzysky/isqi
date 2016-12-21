package views

import (
	ui "github.com/gizak/termui"
	"github.com/mattn/go-runewidth"
	termbox "github.com/nsf/termbox-go"
	"go-webterm"
	"image"
	"sync"
	"unicode/utf8"
)

type ConsoleView struct {
	BaseView
	typing         bool
	value          []byte
	cursor         int
	virtual_cursor int
	view           *ui.Par
	sync.Mutex
}

type Executable interface {
	Exec(string)
}

var input_chan chan rune

func NewConsoleView(rect image.Rectangle) *ConsoleView {
	console := ConsoleView{}
	console.typing = true
	console.rect = rect
	view := ui.NewPar("")
	console.view = view
	console.LocateCursor()
	return &console
}

func (con *ConsoleView) Sync() {
	if con.view == nil {
		return
	}

	con.view.Width = con.rect.Dx()
	con.view.Height = con.rect.Dy()
	con.view.X = con.rect.Min.X
	con.view.Y = con.rect.Min.Y
	con.view.Text = string(con.value)
}

func (con *ConsoleView) Draw() {
	if con.typing {
		cursor := con.LocateCursor()
		termbox.SetCursor(cursor.X, cursor.Y)
	} else {
		termbox.HideCursor()
	}

	ui.Render(con.view)
}

func (con *ConsoleView) Clear() {
}

func (con *ConsoleView) Display() {
	con.Sync()
	con.Draw()
}
func (con *ConsoleView) Continue() {
	con.typing = true
	con.Draw()
}

func (con *ConsoleView) Stop() {
	con.typing = false
	con.Draw()
}

func (con *ConsoleView) Val() string {
	return string(con.value)
}

func (con *ConsoleView) Key(key_str string) {
	if !con.typing {
		return
	}

	switch key_str {
	case "<left>":
		con.Left()
	case "<right>":
		con.Right()
	case "<up>", "<down>":
	case "<enter>":
		con.Type(' ')
	case "<space>":
		con.Type(' ')
	case "C-8":
		con.Delete()
	default:
		debuger.Logf("str got : %s", key_str)
		char, _ := utf8.DecodeRuneInString(key_str)
		con.Type(char)
	}
}

func (con *ConsoleView) Type(char rune) {
	con.Lock()
	defer con.Unlock()

	var buf [utf8.UTFMax]byte
	n := utf8.EncodeRune(buf[:], char)
	offset := runewidth.RuneWidth(char)
	con.value = append(con.value, buf[:n]...)
	con.CursorRight(offset)
	con.VirtualCursorRight(n)
	con.Sync()
	con.Draw()
}

func (con *ConsoleView) Left() {
	if con.virtual_cursor == 0 {
		return
	}
	word, virtual_width := con.RuneBeforeCursor()
	debuger.Logf("%c\n", word)
	con.CursorRight(-runewidth.RuneWidth(word))
	con.VirtualCursorRight(-virtual_width)
	con.Sync()
	con.Draw()
}

func (con *ConsoleView) Right() {
	if con.virtual_cursor == len(con.value) {
		return
	}

	word, virtual_width := con.RuneOnCursor()
	debuger.Logf("%c\n", word)
	con.CursorRight(runewidth.RuneWidth(word))
	con.VirtualCursorRight(virtual_width)
	con.Sync()
	con.Draw()
}

func (con *ConsoleView) Delete() {
	if con.virtual_cursor == 0 {
		return
	}
	char, virtual_width := con.RuneBeforeCursor()
	if con.virtual_cursor < len(con.value) {
		con.value = append(con.value[:con.virtual_cursor-virtual_width], con.value[con.virtual_cursor:]...)
	} else {
		con.value = con.value[:con.virtual_cursor-virtual_width]
	}

	con.CursorRight(-runewidth.RuneWidth(char))
	con.VirtualCursorRight(-virtual_width)
	con.Display()
}

func (con *ConsoleView) CursorRight(n int) {
	con.cursor += n
}

func (con *ConsoleView) VirtualCursorRight(n int) {
	con.virtual_cursor += n
}

func (con *ConsoleView) LocateCursor() image.Point {
	width := con.rect.Dx() - 2
	x := con.cursor%width + 1
	y := 1
	if con.cursor > width {
		y = con.cursor/width + 1
	}
	return image.Pt(x, y)
}

func (con *ConsoleView) RuneBeforeCursor() (rune, int) {
	return utf8.DecodeLastRune(con.value[:con.virtual_cursor])
}
func (con *ConsoleView) RuneOnCursor() (rune, int) {
	return utf8.DecodeRune(con.value[con.virtual_cursor:])
}
