package ui

import (
	"github.com/gpontesss/bib/bib"
	gc "github.com/rthornton128/goncurses"
)

// NewVersionPad docs here.
func NewVersionPad(vsr *bib.Version, height, width, y, x, padding int) (VersionPad, error) {
	// TODO: how to determine appropriate height?
	pad, err := gc.NewPad(10*height, width)
	if err != nil {
		return VersionPad{}, err
	}

	// height/width ration ~ 2
	horpadding, vertpadding := padding, (padding / 2)

	vsrp := VersionPad{
		pad:    pad,
		vsr:    vsr,
		height: height, width: width,
		horpadding: horpadding, vertpadding: vertpadding,
		offset: 0, maxoffset: 0,
		x: x, y: y,
		cursorx: 0, cursory: 0,
		refloaded: bib.Ref{},
	}
	return vsrp, nil
}

// VersionPad docs here.
type VersionPad struct {
	pad                     *gc.Pad
	vsr                     *bib.Version
	height, width           int
	x, y                    int
	horpadding, vertpadding int
	offset, maxoffset       int
	cursorx, cursory        int
	refloaded               bib.Ref
}

func (vsrp *VersionPad) minx() int { return vsrp.horpadding }
func (vsrp *VersionPad) maxx() int { return vsrp.width - vsrp.horpadding }
func (vsrp *VersionPad) miny() int { return vsrp.vertpadding }
func (vsrp *VersionPad) maxy() int { return vsrp.maxoffset }

// SetVersion docs here.
func (vsrp *VersionPad) SetVersion(vsr *bib.Version) { vsrp.vsr = vsr }

// MoveCursor docs here.
func (vsrp *VersionPad) MoveCursor(yoffset, xoffset int) {
	vsrp.GotoCursor(vsrp.cursory+yoffset, vsrp.cursorx+xoffset)
}

// GotoCursor docs here.
func (vsrp *VersionPad) GotoCursor(y, x int) {
	if miny := vsrp.miny(); y < miny {
		y = miny
	} else if maxy := vsrp.maxy() - 1; y > maxy {
		y = maxy
	}
	if minx := vsrp.minx(); x < minx {
		x = minx
	} else if maxx := vsrp.maxx() - 1; x > maxx {
		x = maxx
	}

	if yoffset := y - vsrp.offset; yoffset < 0 {
		vsrp.Scroll(yoffset)
	} else if yoffset := y - vsrp.offset - vsrp.height + 1; yoffset > 0 {
		vsrp.Scroll(yoffset)
	}

	vsrp.cursory, vsrp.cursorx = y, x
	vsrp.pad.Move(vsrp.cursory, vsrp.cursorx)
}

// Scroll docs here.
func (vsrp *VersionPad) Scroll(offset int) {
	vsrp.offset = vsrp.offset + offset
	if vsrp.offset < 0 {
		vsrp.offset = 0
	} else if vsrp.offset > vsrp.maxoffset-1 {
		vsrp.offset = vsrp.maxoffset - 1
	}
	vsrp.GotoCursor(offset+vsrp.cursory, vsrp.cursorx)
}

// NoutRefresh docs here.
func (vsrp *VersionPad) NoutRefresh() {
	vsrp.pad.NoutRefresh(
		// there won't be horizontal offsets for now.
		vsrp.offset, 0,
		vsrp.y, vsrp.x,
		vsrp.y+vsrp.height, vsrp.x+vsrp.width)
}

// Refresh docs here.
func (vsrp *VersionPad) Refresh() {
	vsrp.pad.Refresh(
		// there won't be horizontal offsets for now.
		vsrp.offset, 0,
		vsrp.y, vsrp.x,
		vsrp.y+vsrp.height, vsrp.x+vsrp.width)
}

// GetChar docs here.
func (vsrp *VersionPad) GetChar() gc.Key { return vsrp.pad.GetChar() }

// Delete docs here.
func (vsrp *VersionPad) Delete() { vsrp.pad.Delete() }

// RefLoaded docs here.
func (vsrp *VersionPad) RefLoaded() *bib.Ref { return &vsrp.refloaded }

// LoadRef docs here.
func (vsrp *VersionPad) LoadRef(ref *bib.Ref) {
	vsrp.refloaded = *ref
	vsrp.pad.Erase()
	refp := NewRefPrinter(&vsrp.refloaded, vsrp.vsr, vsrp.width, vsrp.horpadding)
	vsrp.maxoffset = refp.Print(vsrp.pad) + 1
	if vsrp.maxoffset < 0 {
		vsrp.maxoffset = 0
	}
	vsrp.offset = 0
	vsrp.GotoCursor(vsrp.miny(), vsrp.minx())
	vsrp.NoutRefresh()
}
