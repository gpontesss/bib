package ui

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// XY docs here.
type XY struct {
	X, Y uint
}

// Move docs here.
func (xy *XY) Move(x, y int) XY {
	return XY{
		X: uint(max(0, int(xy.X)+x)),
		Y: uint(max(0, int(xy.Y)+y)),
	}
}

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
	height, width uint
}

// NW docs here.
func (box *Box) NW() XY { return box.nw }

// NE docs here.
func (box *Box) NE() XY { return box.nw.Move(int(box.width), 0) }

// SW docs here.
func (box *Box) SW() XY { return box.nw.Move(0, int(box.height)) }

// SE docs here.
func (box *Box) SE() XY { return box.nw.Move(int(box.width), int(box.height)) }

// Pad docs here.
func (box Box) Pad(pad uint) Box { return box.VertPad(pad, pad).HorPad(pad, pad) }

// Move docs here.
func (box Box) Move(hor, vert int) Box {
	return Box{box.nw.Move(hor, vert), box.height, box.width}
}

func (box Box) Resize(height, width uint) Box {
	return Box{box.nw, height, width}
}

// VertPad docs here.
func (box Box) VertPad(top, bottom uint) Box {
	return Box{
		nw:     box.nw.Move(0, int(top)),
		height: box.height - (top + bottom),
		width:  box.width,
	}
}

// HorPad docs here.
func (box Box) HorPad(left, right uint) Box {
	return Box{
		nw:     box.nw.Move(int(left), 0),
		height: box.height,
		width:  box.width - (left + right),
	}
}

// BoxIter docs here.
type BoxIter struct {
	Next  func() bool
	Value func() Box
	Index func() int
}

// VertDiv docs here.
func (box Box) VertDiv(n uint) BoxIter {
	divwidth := box.width / n
	i := -1

	return BoxIter{
		Index: func() int { return i },
		Next: func() bool {
			i++
			return i < int(n)
		},
		Value: func() Box {
			return Box{
				// uint cast is safe for i is never negative after a `Next` call.
				XY{box.nw.X + divwidth*uint(i), box.nw.Y},
				box.height,
				divwidth,
			}
		},
	}
}

// HorDiv docs here.
func (box Box) HorDiv(n uint) BoxIter {
	divheight := box.height / n
	i := -1

	return BoxIter{
		Index: func() int { return i },
		Next: func() bool {
			i++
			return i < int(n)
		},
		Value: func() Box {
			return Box{
				// uint cast is safe for i is never negative after a `Next` call.
				XY{box.nw.X, box.nw.Y + divheight*uint(i)},
				divheight,
				box.width,
			}
		},
	}
}
