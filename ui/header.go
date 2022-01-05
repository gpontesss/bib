package ui

import gc "github.com/gbin/goncurses"

// Header docs here.
type Header struct {
	WinBox
	text string
}

// CastHeader docs here.
func CastHeader(win WinBox, text string) Header {
	return Header{win, text}
}

// NewHeader docs here.
func NewHeader(box Box, text string) (Header, error) {
	win, err := NewBoxWin(box)
	return Header{win, text}, err
}

// Refresh docs here.
func (hdr *Header) Refresh() {
	hdr.Draw()
	hdr.WinBox.Refresh()
}

// NoutRefresh docs here.
func (hdr *Header) NoutRefresh() {
	hdr.Draw()
	hdr.WinBox.NoutRefresh()
}

// SetText docs here.
func (hdr *Header) SetText(text string) { hdr.text = text }

// Draw docs here.
func (hdr *Header) Draw() {
	hdr.Erase()
	hdr.SetBackground(gc.ColorPair(2))
	hdr.AttrOn(gc.ColorPair(2) | gc.A_BOLD)
	// TODO: wrap text as necessary.
	hdr.MovePrint(0, 0, hdr.text)
	hdr.AttrOff(gc.ColorPair(2) | gc.A_BOLD)
}
