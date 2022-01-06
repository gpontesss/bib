package ui

import (
	gc "github.com/gbin/goncurses"
	"github.com/gpontesss/bib/bib"
)

// VersionMenu docs here.
type VersionMenu struct {
	WinBox
	sub   WinBox
	vsrs  []*bib.Version
	menu  *gc.Menu
	items []*gc.MenuItem
}

// NewVersionMenu docs here.
func NewVersionMenu(box Box, vsrs ...*bib.Version) (vmenu VersionMenu, err error) {
	vmenu.vsrs = vsrs
	vmenu.items = make([]*gc.MenuItem, len(vmenu.vsrs))
	for i := range vmenu.vsrs {
		if vmenu.items[i], err = gc.NewItem(vmenu.vsrs[i].Name, ""); err != nil {
			return
		}
	}
	if vmenu.WinBox, err = NewBoxWin(MinBox()); err != nil {
		return
	}
	if vmenu.menu, err = gc.NewMenu(vmenu.items); err != nil {
		return
	}
	if err = vmenu.menu.SetWindow(vmenu.WinBox.Window); err != nil {
		return
	}

	vmenu.menu.Mark(" * ")
	vmenu.menu.SetBackground(0)

	vmenu.sub = vmenu.WinBox.Sub(MinBox())
	vmenu.menu.SubWindow(vmenu.sub.Window)

	vmenu.Resize(box)

	return vmenu, nil
}

// Resize docs here.
func (vmenu *VersionMenu) Resize(box Box) {
	vmenu.WinBox.ResizeBox(box)
	// -1 to consider the border; -3 to consider the top header
	vmenu.sub.ResizeBox(vmenu.HorPad(1, 1).VertPad(3, 1).MoveXY(XY{1, 3}))
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
	win := vmenu.Window
	win.Box(0, 0)

	// Menu header, considering the border offset.
	// TODO: migrate to Header.
	win.MovePrint(1, 1, "Select the version")
	win.HLine(2, 1, 0, int(vmenu.width)-2)
	win.NoutRefresh()
}

// Current docs here.
func (vmenu *VersionMenu) Current() *bib.Version {
	return vmenu.vsrs[vmenu.menu.Current(nil).Index()]
}

// Delete docs here.
func (vmenu *VersionMenu) Delete() {
	vmenu.sub.Delete()
	vmenu.WinBox.Delete()
	for _, item := range vmenu.items {
		if item != nil {
			item.Free()
		}
	}
	if vmenu.menu != nil {
		vmenu.menu.Free()
	}
}
