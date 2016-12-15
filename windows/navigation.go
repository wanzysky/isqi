package windows

type Pair struct {
	index int
	value Naviable
}

type Naviable interface {
	Display()
	Clear()
}

type Navigation struct {
	windows []Naviable
	current Pair
}

var Nav *Navigation

func NewNavigatoin(window Naviable) *Navigation {
	var nav Navigation
	nav.windows = make([]Naviable, 0)
	nav.windows = append(nav.windows, window)
	nav.current = Pair{0, window}
	window.Display()
	Nav = &nav
	return Nav
}

func (nav *Navigation) Push(window Naviable) *Navigation {
	current_index := nav.current.index
	nav.current.value.Clear()
	nav.windows = append(nav.windows[:current_index+1], window)
	nav.current.index += 1
	nav.current.value = window
	nav.current.value.Display()
	return nav
}

func (nav *Navigation) Back() *Navigation {
	if nav.current.index > 0 {
		nav.current.index -= 1
		nav.current.value.Clear()
		nav.current.value = nav.windows[nav.current.index]
		nav.current.value.Display()
	}
	return nav
}
