package ui

import (
	"github.com/gpontesss/bib/bib"
	gc "github.com/rthornton128/goncurses"
)

// NewVersionMenu docs here.
func NewVersionMenu(y, x, h, w int, vsrs ...*bib.Version) (VersionMenu, error) {
	var err error
	items := make([]*gc.MenuItem, len(vsrs))
	for i := range items {
		items[i], err = gc.NewItem(vsrs[i].Name, "")
		if err != nil {
			return VersionMenu{}, err
		}
	}

	vmenu := VersionMenu{
		vsrs:  vsrs,
		width: w, height: h,
	}
	vmenu.menu, err = gc.NewMenu(items)
	if err != nil {
		return VersionMenu{}, err
	}

	menuwin, err := gc.NewWindow(h, w, y, x)
	if err != nil {
		return VersionMenu{}, err
	}
	if err = vmenu.menu.SetWindow(menuwin); err != nil {
		return VersionMenu{}, nil
	}

	// -1 to consider the border; -3 to consider the top header
	derwin := menuwin.Derived(h-1-3, w-1, 3, 1)
	vmenu.menu.SubWindow(derwin)

	vmenu.menu.Mark(" * ")
	vmenu.menu.SetBackground(0)

	return vmenu, nil
}

// VersionMenu docs here.
type VersionMenu struct {
	menu          *gc.Menu
	vsrs          []*bib.Version
	width, height int
}

// Select docs here.
func (vmenu *VersionMenu) Select() (*bib.Version, error) {
	err := vmenu.menu.Post()
	if err != nil {
		return nil, err
	}
	defer vmenu.menu.UnPost()

	// Disables cursor while selecting.
	gc.Cursor(0)
	defer gc.Cursor(1)

	vmenu.Draw()

	for {
		gc.Update()
		switch vmenu.menu.Window().GetChar() {
		case 'q': // Quits without selecting anything.
			return nil, nil
		case gc.KEY_RETURN, gc.KEY_TAB: // Selects current version.
			return vmenu.Current(), nil
		case 'j': // Selects item bellow of current.
			vmenu.menu.Driver(gc.REQ_DOWN)
		case 'k': // Selects item on top of current.
			vmenu.menu.Driver(gc.REQ_UP)
		}
	}
}

// Draw docs here.
func (vmenu *VersionMenu) Draw() {
	menuwin := vmenu.menu.Window()
	menuwin.Box(0, 0)

	// Menu header, considering the border offset.
	menuwin.MovePrint(1, 1, "Select the version")
	menuwin.HLine(2, 1, 0, vmenu.width-2)
	menuwin.NoutRefresh()
}

// Current docs here.
func (vmenu *VersionMenu) Current() *bib.Version {
	return vmenu.vsrs[vmenu.menu.Current(nil).Index()]
}

// Delete docs here.
func (vmenu *VersionMenu) Delete() {
	for _, item := range vmenu.menu.Items() {
		item.Free()
	}
	vmenu.menu.Window().Delete()
	vmenu.menu.Free()
}
