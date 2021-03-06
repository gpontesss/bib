package ui

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// XY docs here.
type XY struct {
	X, Y int
}

// Neg docs here.
func (xy XY) Neg() XY { return XY{-xy.X, -xy.Y} }

// Eq docs here.
func (xy XY) Eq(xy2 XY) bool { return xy.X == xy2.X && xy.Y == xy2.Y }

// RelTo docs here.
func (xy *XY) RelTo(xy2 XY) XY { return XY{xy2.X - xy.X, xy2.Y - xy.Y} }

// Move docs here.
func (xy XY) Move(x, y int) XY {
	return XY{
		X: xy.X + x,
		Y: xy.Y + y,
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

// CenteredBox docs here.
func CenteredBox(center XY, height, width uint) Box {
	return Box{center.Move(int(width/2), int(height/2)), height, width}
}

// Center docs here.
func (box Box) Center() XY { return box.nw.Move(int(box.width/2), int(box.height/2)) }

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
func (box Box) Move(hor, vert int) Box { return Box{box.nw.Move(hor, vert), box.height, box.width} }

// MoveXY docs here.
func (box Box) MoveXY(nw XY) Box {
	box.nw = nw
	return box
}

// Resize docs here.
func (box Box) Resize(height, width uint) Box { return Box{box.nw, height, width} }

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
				XY{box.nw.X + int(divwidth)*i, box.nw.Y},
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
				XY{box.nw.X, box.nw.Y + int(divheight)*i},
				divheight,
				box.width,
			}
		},
	}
}
