package ui

import (
	"github.com/gbin/goncurses"
	"github.com/gpontesss/bib/bib"
)

// NewVersionPad docs here.
func NewVersionPad(vsr *bib.Version, height, width, y, x int) (VersionPad, error) {
	pad, err := goncurses.NewPad(height, width)
	if err != nil {
		return VersionPad{}, err
	}

	pad.ScrollOk(true)

	return VersionPad{
		pad:     pad,
		version: vsr,
		height:  height,
		width:   width,
		offset:  0,
		x:       x,
		y:       y,
	}, nil
}

// VersionPad docs here.
type VersionPad struct {
	pad           *goncurses.Pad
	version       *bib.Version
	height, width int
	x, y          int
	offset        int
}

// Scroll docs here.
func (vsrp *VersionPad) Scroll(scaler int) {
	vsrp.offset += scaler
	if vsrp.offset < 0 {
		vsrp.offset = 0
	}
}

// NoutRefresh docs here.
func (vsrp *VersionPad) NoutRefresh() {
	vsrp.pad.NoutRefresh(
		// there won't be horizontal offsets for now.
		vsrp.offset, 0,
		vsrp.y, vsrp.x,
		vsrp.height, vsrp.width)
}

// Delete docs here.
func (vsrp *VersionPad) Delete() { vsrp.pad.Delete() }

// DisplayRef docs here.
func (vsrp *VersionPad) LoadRef(ref *bib.Ref) {
	verses := ref.Verses(vsrp.version)

	vsrp.pad.Clear()
	vsrp.pad.MovePrintln(0, 0, vsrp.version.Name, " ", ref)
	for _, verse := range verses {
		vsrp.pad.Print(verse)
	}
}
