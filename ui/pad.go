package ui

import (
	"github.com/gpontesss/bib/bib"
	gc "github.com/rthornton128/goncurses"
)

// NewVersionPad docs here.
func NewVersionPad(vsr *bib.Version, height, width, y, x, padding int) (VersionPad, error) {
	// TODO: how to determine appropriate height?
	pad, err := gc.NewPad(10*height, width-1)
	if err != nil {
		return VersionPad{}, err
	}

	return VersionPad{
		pad:     pad,
		version: vsr,
		height:  height, width: width,
		padding: padding,
		offset:  0, maxoffset: 0,
		x: x, y: y,
		// TODO: deal with vert padding.
		cursorx: 0, cursory: 0,
	}, nil
}

// VersionPad docs here.
type VersionPad struct {
	pad               *gc.Pad
	version           *bib.Version
	height, width     int
	x, y              int
	padding           int
	offset, maxoffset int
	cursorx, cursory  int
}

// MoveCursor docs here.
func (vsrp *VersionPad) MoveCursor(yoffset, xoffset int) {
	y := vsrp.cursory + yoffset
	x := vsrp.cursorx + xoffset
	vsrp.GotoCursor(y, x)
}

// GotoCursor docs here.
func (vsrp *VersionPad) GotoCursor(y, x int) {
	if miny := 0; y < miny {
		y = miny
	} else if maxy := vsrp.maxoffset; y > maxy {
		y = maxy
	}
	if minx := 0; x < minx {
		x = minx
	} else if maxx := vsrp.width; x > maxx {
		x = maxx
	}

	if y-vsrp.offset < 0 {
		vsrp.Scroll(-1)
	} else if y-vsrp.offset > vsrp.height {
		vsrp.Scroll(1)
	}

	vsrp.cursory, vsrp.cursorx = y, x
	vsrp.pad.Move(vsrp.cursory, vsrp.cursorx)
}

// Scroll docs here.
func (vsrp *VersionPad) Scroll(offset int) {
	y := vsrp.offset + offset
	vsrp.Goto(y)
}

// Goto docs here.
func (vsrp *VersionPad) Goto(y int) {
	vsrp.offset = y
	if vsrp.offset < 0 {
		vsrp.offset = 0
	} else if vsrp.offset > vsrp.maxoffset {
		vsrp.offset = vsrp.maxoffset
	}
	vsrp.GotoCursor(y, 0)
}

// NoutRefresh docs here.
func (vsrp *VersionPad) NoutRefresh() {
	vsrp.pad.NoutRefresh(
		// there won't be horizontal offsets for now.
		vsrp.offset, 0,
		vsrp.y, vsrp.x,
		// subtracts one, since coordinate is 0 based.
		vsrp.y+vsrp.height-1, vsrp.x+vsrp.width-1)
}

// Refresh docs here.
func (vsrp *VersionPad) Refresh() {
	vsrp.pad.Refresh(
		// there won't be horizontal offsets for now.
		vsrp.offset, 0,
		vsrp.y, vsrp.x,
		// subtracts one, since coordinate is 0 based.
		vsrp.y+vsrp.height-1, vsrp.x+vsrp.width-1)
}

// GetChar docs here.
func (vsrp *VersionPad) GetChar() gc.Key { return vsrp.pad.GetChar() }

// Delete docs here.
func (vsrp *VersionPad) Delete() { vsrp.pad.Delete() }

// LoadRef docs here.
func (vsrp *VersionPad) LoadRef(ref *bib.Ref) {
	vsrp.pad.Erase()
	refp := NewRefPrinter(ref, vsrp.version, vsrp.width, vsrp.padding)
	vsrp.offset = 0
	vsrp.maxoffset = refp.Print(vsrp.pad) + 1
	if vsrp.maxoffset < 0 {
		vsrp.maxoffset = 0
	}
	vsrp.NoutRefresh()
}
