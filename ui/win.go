package ui

import gc "github.com/gbin/goncurses"

type WinBox struct {
	Box
	*gc.Window
}

// MinBox docs here.
func MinBox() Box { return Box{XY{0, 0}, 1, 1} }

// NewBoxWin docs here.
func NewBoxWin(box Box) (WinBox, error) {
	win, err := gc.NewWindow(
		int(box.height), int(box.width),
		int(box.nw.Y), int(box.nw.X))
	return WinBox{box, win}, err
}

// Sub docs here.
func (wb *WinBox) Sub(box Box) WinBox {
	// TODO: should sub stuff be relative to parent? (seems to make sense)
	return WinBox{
		box,
		wb.Window.Sub(
			int(box.height), int(box.width),
			int(box.nw.Y), int(box.nw.X)),
	}
}

// Resize docs here.
func (wb *WinBox) Resize(height, width uint) {
	wb.ResizeBox(wb.Box.Resize(height, width))
}

// MoveXY docs here.
func (wb *WinBox) MoveXY(xy XY) {
	wb.Window.MoveWindow(int(xy.Y), int(xy.X))
}

// Resize docs here.
func (wb *WinBox) ResizeBox(box Box) {
	if wb.Box.nw != box.nw {
		wb.MoveXY(box.nw)
	}
	wb.Window.Resize(int(box.height), int(box.width))
	wb.Box = box
}

// PadBox docs here.
type PadBox struct {
	Box
	*gc.Pad
	offset XY
}

// NewPadBox docs here.
func NewPadBox(box Box, height, width int) (PadBox, error) {
	pad, err := gc.NewPad(height, width)
	return PadBox{box, pad, XY{0, 0}}, err
}

// Resize docs here.
func (pb *PadBox) Resize(box Box) { pb.Box = box }

// ResizeBuffer docs here.
func (pb *PadBox) ResizeBuffer(height, width uint) {
	pb.Pad.Resize(int(height), int(width))
}

// Refresh docs here.
func (pb *PadBox) Refresh() {
	pb.Pad.Refresh(
		int(pb.offset.Y), int(pb.offset.X),
		int(pb.Box.NW().Y), int(pb.Box.NW().X),
		int(pb.Box.SE().Y), int(pb.Box.SE().X),
	)
}

// NoutRefresh docs here.
func (pb *PadBox) NoutRefresh() {
	pb.Pad.NoutRefresh(
		int(pb.offset.Y), int(pb.offset.X),
		int(pb.Box.NW().Y), int(pb.Box.NW().X),
		int(pb.Box.SE().Y), int(pb.Box.SE().X),
	)
}

// Offset docs here.
func (pb *PadBox) Offset(offset XY) { pb.offset = offset }