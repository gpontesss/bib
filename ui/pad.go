package ui

import (
	gc "github.com/gbin/goncurses"
	"github.com/gpontesss/bib/bib"
)

// NewVersionPad docs here.
func NewVersionPad(
	vsr *bib.Version, box Box, padding uint) (vsrp VersionPad, err error) {

	if vsrp.WinBox, err = NewBoxWin(MinBox()); err != nil {
		return VersionPad{}, err
	}
	vsrp.header = vsrp.Sub(Box{})

	if vsrp.pad, err = NewPadBox(Box{}, 1, 1); err != nil {
		return VersionPad{}, err
	}

	vsrp.vsr = vsr

	vsrp.Resize(box, padding)
	return vsrp, nil
}

// VersionPad docs here.
type VersionPad struct {
	WinBox

	header WinBox
	pad    PadBox

	display Box
	padding uint

	cursor            XY
	offset, maxoffset uint

	vsr       *bib.Version
	refloaded bib.Ref
}

func (vsrp *VersionPad) Resize(box Box, padding uint) {
	vsrp.display = box.Pad(padding)

	vsrp.ResizeBox(box)
	vsrp.header.ResizeBox(Box{vsrp.display.NW(), 1, vsrp.display.width})
	vsrp.pad.Resize(vsrp.display.VertPad(vsrp.header.height, 1))

	// Forces pad refresh, while reloading text with new dimensions.
	if ref := vsrp.RefLoaded(); ref != nil {
		vsrp.LoadRef(vsrp.RefLoaded())
	}
}

func (vsrp *VersionPad) minx() uint { return 0 }
func (vsrp *VersionPad) maxx() uint { return vsrp.display.width }
func (vsrp *VersionPad) miny() uint { return 0 }
func (vsrp *VersionPad) maxy() uint { return uint(vsrp.maxoffset - 1) }

// SetVersion docs here.
func (vsrp *VersionPad) SetVersion(vsr *bib.Version) { vsrp.vsr = vsr }

// MoveCursor docs here.
func (vsrp *VersionPad) MoveCursor(yoffset, xoffset int) {
	vsrp.cursor = vsrp.pad.MoveCursorXY(vsrp.cursor.Move(xoffset, yoffset))

	// fixes offset if needed.
	if ydiff := vsrp.cursor.Y - (int(vsrp.pad.height) + vsrp.pad.offset.Y); ydiff > 0 {
		vsrp.pad.Scroll(0, ydiff)
	} else if ydiff := vsrp.cursor.Y - vsrp.pad.offset.Y; ydiff < 0 {
		vsrp.pad.Scroll(0, ydiff)
	}
}

// GotoCursor docs here.
func (vsrp *VersionPad) GotoCursor(y, x uint) {
	offset := vsrp.cursor.RelTo(XY{int(x), int(y)})
	vsrp.MoveCursor(offset.Y, offset.X)
}

// Scroll docs here.
func (vsrp *VersionPad) Scroll(offset int) {
	vsrp.pad.Scroll(0, offset)
	vsrp.cursor = vsrp.pad.MoveCursorXY(vsrp.cursor.Move(0, offset))
}

// NoutRefresh docs here.
func (vsrp *VersionPad) NoutRefresh() {
	vsrp.WinBox.NoutRefresh()
	vsrp.header.NoutRefresh()
	vsrp.pad.NoutRefresh()
}

// Refresh docs here.
func (vsrp *VersionPad) Refresh() {
	vsrp.WinBox.Refresh()
	vsrp.header.Refresh()
	vsrp.pad.Refresh()
}

// GetChar docs here.
func (vsrp *VersionPad) GetChar() gc.Key { return vsrp.pad.GetChar() }

// Delete docs here.
func (vsrp *VersionPad) Delete() {
	// TODO: include other components too.
	vsrp.pad.Delete()
}

// RefLoaded docs here.
func (vsrp *VersionPad) RefLoaded() *bib.Ref { return &vsrp.refloaded }

// LoadRef docs here.
func (vsrp *VersionPad) LoadRef(ref *bib.Ref) {
	vsrp.refloaded = *ref
	vsrp.pad.Erase()

	refp := NewRefPrinter(&vsrp.refloaded, vsrp.vsr, vsrp.display.width)
	vsrp.maxoffset = refp.LinesRequired()
	vsrp.pad.ResizeBuffer(
		// +height avoids text shadows at end when scrolling near end of text.
		// vsrp.maxoffset+vsrp.display.height,
		vsrp.maxoffset,
		vsrp.display.width)

	vsrp.WinBox.Erase()
	vsrp.WinBox.NoutRefresh()

	vsrp.header.SetBackground(gc.ColorPair(2))
	vsrp.header.AttrOn(gc.ColorPair(2) | gc.A_BOLD)
	vsrp.header.MovePrint(0, 0, vsrp.vsr.Name, " ", &vsrp.refloaded) // header
	vsrp.header.AttrOff(gc.ColorPair(2) | gc.A_BOLD)

	refp.Print(&vsrp.pad)

	vsrp.offset = 0
	vsrp.GotoCursor(vsrp.miny(), vsrp.minx())

	vsrp.NoutRefresh()
}
