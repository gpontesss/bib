package ui

// XY docs here.
type XY struct {
	X, Y int
}

// Move docs here.
func (xy *XY) Move(x, y int) XY { return XY{X: xy.X + x, Y: xy.Y + y} }

// Bound docs here.
func (xy XY) Bound(box *Box) XY {
	nw, se := box.NW(), box.SE()
	if xy.X < nw.X {
		xy.X = nw.X
	} else if xy.X > se.X {
		xy.X = se.X
	}

	if xy.Y < nw.Y {
		xy.Y = nw.Y
	} else if xy.Y > se.Y {
		xy.Y = se.Y
	}

	return xy
}

// Box docs here.
type Box struct {
	nw            XY
	height, width int
}

// NW docs here.
func (box *Box) NW() XY { return box.nw }

// NE docs here.
func (box *Box) NE() XY { return box.nw.Move(box.width, 0) }

// SW docs here.
func (box *Box) SW() XY { return box.nw.Move(0, box.height) }

// SE docs here.
func (box *Box) SE() XY { return box.nw.Move(box.width, box.height) }

// Pad docs here.
func (box Box) Pad(pad int) Box { return box.VertPad(pad, pad).HorPad(pad, pad) }

// VertPad docs here.
func (box Box) VertPad(top, bottom int) Box {
	return Box{
		nw:     box.nw.Move(0, top),
		height: box.height - (top + bottom),
		width:  box.width,
	}
}

// HorPad docs here.
func (box Box) HorPad(left, right int) Box {
	return Box{
		nw:     box.nw.Move(left, 0),
		height: box.height,
		width:  box.width - (left + right),
	}
}
