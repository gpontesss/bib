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

	// pad.ScrollOk(true)

	return VersionPad{
		pad:       pad,
		version:   vsr,
		height:    height,
		width:     width,
		padding:   padding,
		offset:    0,
		maxoffset: 0,
		x:         x,
		y:         y,
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
}

// NoutRefresh docs here.
func (vsrp *VersionPad) NoutRefresh() {
	vsrp.pad.NoutRefresh(
		// there won't be horizontal offsets for now.
		vsrp.offset, 0,
		vsrp.y, vsrp.x,
		// subtracts one, since coordinate is 0 based.
		vsrp.height-1, vsrp.width-1)
}

// Refresh docs here.
func (vsrp *VersionPad) Refresh() {
	vsrp.pad.Refresh(
		// there won't be horizontal offsets for now.
		vsrp.offset, 0,
		vsrp.y, vsrp.x,
		// subtracts one, since coordinate is 0 based.
		vsrp.height-1, vsrp.width-1)
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
	vsrp.maxoffset = refp.Print(vsrp.pad) - (vsrp.width / 2)
	if vsrp.maxoffset < 0 {
		vsrp.maxoffset = 0
	}
	vsrp.NoutRefresh()
}
