package views

import (
	ui "github.com/gizak/termui"
	"github.com/mattn/go-runewidth"
	"image"
	"strings"
)

type CollectionView struct {
	BaseView
	headers  []string
	contents []string
	h_views  []*ui.Par
	c_views  []*ui.Par
	offset   int
}

func NewCollectionView(rect image.Rectangle, headers, contents []string) *CollectionView {
	col := &CollectionView{}
	col.rect = rect
	col.headers = headers
	col.contents = contents
	return col
}

func (col *CollectionView) OffsetContents() ([]string, []string) {
	return col.headers[col.offset:], col.contents[col.offset:]
}

func (col *CollectionView) Up() {
	if col.offset == 0 {
		return
	}
	col.offset -= 1
	col.Clear()
	col.Display()
}

func (col *CollectionView) Down() {
	if col.offset == len(col.headers)-1 {
		return
	}
	col.offset += 1
	col.Clear()
	col.Display()
}

func (col *CollectionView) Draw() {
}

func (col *CollectionView) Sync() {
	headers, contents := col.OffsetContents()
	line_width := col.rect.Dx()
	width_of_content := line_width - 2
	lines := col.rect.Dy()
	for index, header := range headers {
		if lines <= 0 {
			break
		}
		content := contents[index]
		content_view := ui.NewPar(content)
		content_view.BorderLabel = header
		lines_need := 2
		for _, line := range strings.Split(content, "\n") {
			content_length := runewidth.StringWidth(line)
			lines_need += content_length/width_of_content + 1
		}

		content_view.Width = line_width
		content_view.Height = lines_need
		content_view.X = col.rect.Min.X
		content_view.Y = col.rect.Max.Y - lines

		lines -= lines_need
		ui.Render(content_view)
	}
}

func (col *CollectionView) Clear() {
	ui.ClearArea(col.rect, DEFAULT_BG_COLOR)
}

func (col *CollectionView) Display() {
	col.Sync()
	col.Draw()
}
