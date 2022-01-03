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
	offset              XY
	bufheight, bufwidth uint
}

// NewPadBox docs here.
func NewPadBox(box Box, height, width uint) (PadBox, error) {
	pad, err := gc.NewPad(int(height), int(width))
	return PadBox{box, pad, XY{0, 0}, height, width}, err
}

// Resize docs here.
func (pb *PadBox) Resize(box Box) { pb.Box = box }

// ResizeBuffer docs here.
func (pb *PadBox) ResizeBuffer(height, width uint) {
	// temporary fix to avoid text being shadowed near end of scroll.
	pb.Pad.Resize(int(height+pb.height), int(width+pb.width))
	pb.bufheight, pb.bufwidth = height, width
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

// BoundBuffer docs here.
func (pb *PadBox) BoundBuffer(x, y int) XY {
	// -1 so last line appears when pad is completely scrolled right.
	if x > int(pb.bufwidth)-1 {
		x = int(pb.bufwidth) - 1
	} else if x < 0 {
		x = 0
	}
	// -1 so last line appears when pad is completely scrolled down.
	if y > int(pb.bufheight)-1 {
		y = int(pb.bufheight) - 1
	} else if y < 0 {
		y = 0
	}
	return XY{x, y}
}

func (pb *PadBox) BoundBufferXY(xy XY) XY { return pb.BoundBuffer(xy.X, xy.Y) }

// MoveCursorXY docs here.
func (pb *PadBox) MoveCursorXY(xy XY) XY { return pb.MoveCursor(uint(xy.X), uint(xy.Y)) }

// Scroll docs here.
func (pb *PadBox) Scroll(x, y int) {
	pb.offset = pb.BoundBufferXY(pb.offset.Move(x, y))
}

// MoveCursor docs here.
func (pb *PadBox) MoveCursor(x, y uint) XY {
	// TODO: fix offset accordingly.
	xy := pb.BoundBuffer(int(x), int(y))
	pb.Pad.Move(int(xy.Y), int(xy.X))
	return xy
}
