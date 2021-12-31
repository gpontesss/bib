package ui

import (
	gc "github.com/gbin/goncurses"
	"github.com/gpontesss/bib/bib"
)

// NewVersionPad docs here.
func NewVersionPad(vsr *bib.Version, box Box, padding int) (VersionPad, error) {
	mainwin, err := gc.NewWindow(1, 1, 0, 0)
	if err != nil {
		return VersionPad{}, err
	}
	pad, err := gc.NewPad(1, 1)
	if err != nil {
		return VersionPad{}, err
	}
	header := mainwin.Sub(0, 0, 0, 0)

	vsrp := VersionPad{
		mainwin: mainwin,
		header:  header,
		pad:     pad,
		vsr:     vsr,
		box:     box,
	}

	vsrp.Resize(box, padding)
	return vsrp, nil
}

// VersionPad docs here.
type VersionPad struct {
	mainwin                 *gc.Window
	header                  *gc.Window
	pad                     *gc.Pad
	vsr                     *bib.Version
	box                     Box
	cursor                  XY
	horpadding, vertpadding int
	offset, maxoffset       int
	refloaded               bib.Ref
}

func (vsrp *VersionPad) minx() int { return 0 }
func (vsrp *VersionPad) maxx() int { return vsrp.box.width - (2 * vsrp.horpadding) }
func (vsrp *VersionPad) miny() int { return 0 }
func (vsrp *VersionPad) maxy() int { return vsrp.maxoffset - 1 }

// SetVersion docs here.
func (vsrp *VersionPad) SetVersion(vsr *bib.Version) { vsrp.vsr = vsr }

// MoveCursor docs here.
func (vsrp *VersionPad) MoveCursor(yoffset, xoffset int) {
	vsrp.GotoCursor(vsrp.cursor.Y+yoffset, vsrp.cursor.X+xoffset)
}

// GotoCursor docs here.
func (vsrp *VersionPad) GotoCursor(y, x int) {
	if miny := vsrp.miny(); y < miny {
		y = miny
	} else if maxy := vsrp.maxy(); y > maxy {
		y = maxy
	}
	if minx := vsrp.minx(); x < minx {
		x = minx
	} else if maxx := vsrp.maxx(); x > maxx {
		x = maxx
	}

	if yoffset := y - vsrp.offset; yoffset < 0 {
		vsrp.Scroll(yoffset)
	} else if yoffset := y - vsrp.offset - vsrp.box.height + (2 * vsrp.vertpadding) + 2; yoffset > 0 {
		vsrp.Scroll(yoffset)
	}

	vsrp.cursor = XY{x, y}
	vsrp.pad.Move(vsrp.cursor.Y, vsrp.cursor.X)
}

// Scroll docs here.
func (vsrp *VersionPad) Scroll(offset int) {
	vsrp.offset = vsrp.offset + offset
	if vsrp.offset < 0 {
		vsrp.offset = 0
	} else if vsrp.offset > vsrp.maxoffset {
		vsrp.offset = vsrp.maxoffset
	}
	vsrp.GotoCursor(vsrp.cursor.Y+offset, vsrp.cursor.X)
}

// NoutRefresh docs here.
func (vsrp *VersionPad) NoutRefresh() {
	vsrp.mainwin.NoutRefresh()
	vsrp.header.NoutRefresh()
	vsrp.pad.NoutRefresh(
		// there won't be horizontal offsets for now.
		vsrp.offset, 0,
		// +1 accounts the header.
		vsrp.box.nw.Y+vsrp.vertpadding+1,
		vsrp.box.nw.X+vsrp.horpadding,
		// -1 accounts the header.
		vsrp.box.nw.Y+vsrp.box.height-vsrp.vertpadding-1,
		vsrp.box.nw.X+vsrp.box.width-vsrp.horpadding)
}

// Refresh docs here.
func (vsrp *VersionPad) Refresh() {
	vsrp.mainwin.Refresh()
	vsrp.header.Refresh()
	vsrp.pad.Refresh(
		// there won't be horizontal offsets for now.
		vsrp.offset, 0,
		// +1 accounts the header.
		vsrp.box.nw.Y+vsrp.vertpadding+1,
		vsrp.box.nw.X+vsrp.horpadding,
		// -1 accounts the header.
		vsrp.box.nw.Y+vsrp.box.height-vsrp.vertpadding-1,
		vsrp.box.nw.X+vsrp.box.width-vsrp.horpadding)
}

func (vsrp *VersionPad) Resize(box Box, padding int) {
	// height/width ration ~ 2
	vsrp.horpadding, vsrp.vertpadding = padding, (padding / 2)

	vsrp.box = box

	vsrp.mainwin.MoveWindow(vsrp.box.nw.Y, vsrp.box.nw.X)
	vsrp.mainwin.Resize(vsrp.box.height, vsrp.box.width)

	vsrp.header.Resize(1, vsrp.box.width-(vsrp.horpadding*2))
	vsrp.header.MoveWindow(
		vsrp.box.nw.Y+vsrp.vertpadding,
		vsrp.box.nw.X+vsrp.horpadding)

	// Forces pad refresh, while reloading text with new dimensions.
	if ref := vsrp.RefLoaded(); ref != nil {
		vsrp.LoadRef(vsrp.RefLoaded())
	}
}

// GetChar docs here.
func (vsrp *VersionPad) GetChar() gc.Key { return vsrp.pad.GetChar() }

// Delete docs here.
func (vsrp *VersionPad) Delete() {
	// TODO: include other components too.
	if vsrp.pad != nil {
		vsrp.pad.Delete()
	}
}

// RefLoaded docs here.
func (vsrp *VersionPad) RefLoaded() *bib.Ref { return &vsrp.refloaded }

// LoadRef docs here.
func (vsrp *VersionPad) LoadRef(ref *bib.Ref) {
	vsrp.refloaded = *ref
	vsrp.pad.Erase()

	refp := NewRefPrinter(
		&vsrp.refloaded,
		vsrp.vsr,
		vsrp.box.width-(2*vsrp.horpadding))
	vsrp.maxoffset = refp.LinesRequired()
	// +height avoids text shadows at end when scrolling near end of text.
	vsrp.pad.Resize(
		vsrp.maxoffset+vsrp.box.height,
		vsrp.box.width-(2*vsrp.horpadding))

	vsrp.mainwin.Erase()
	vsrp.mainwin.NoutRefresh()

	vsrp.header.SetBackground(gc.ColorPair(2))
	vsrp.header.AttrOn(gc.ColorPair(2) | gc.A_BOLD)
	vsrp.header.MovePrint(0, 0, vsrp.vsr.Name, " ", &vsrp.refloaded) // header
	vsrp.header.AttrOff(gc.ColorPair(2) | gc.A_BOLD)

	refp.Print(vsrp.pad)

	vsrp.offset = 0
	vsrp.GotoCursor(vsrp.miny(), vsrp.minx())

	vsrp.NoutRefresh()
}
