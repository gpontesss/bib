package ui

import (
	"github.com/gpontesss/bib/bib"
)

// MovePrinter docs here.
type MovePrinter interface {
	MovePrint(y, x int, args ...interface{})
	MovePrintln(y, x int, args ...interface{})
	MovePrintf(y, x int, fmt string, args ...interface{})
}

// RefPrinter docs here.
type RefPrinter struct {
	ref   *bib.Ref
	vsr   *bib.Version
	width int
}

// NewRefPrinter docs here.
func NewRefPrinter(ref *bib.Ref, vsr *bib.Version, width int) RefPrinter {
	return RefPrinter{
		ref:   ref,
		vsr:   vsr,
		width: width,
	}
}

// SetVersion docs here.
func (rp *RefPrinter) SetVersion(vsr *bib.Version) { rp.vsr = vsr }

// Print docs here.
func (rp *RefPrinter) Print(mp MovePrinter) int {
	mp.MovePrint(0, 0, rp.vsr.Name, " ", rp.ref) // header

	linei := 1
	for _, verse := range rp.ref.Verses(rp.vsr) {
		text := verse.String()
		textlen := len(text)

		// A more secure way of doing ceiling division with integers.
		linecount := textlen / rp.width
		if (textlen % rp.width) != 0 {
			linecount++
		}

		for i := 0; i < linecount; i++ {
			linei++
			// TODO: padding?
			lowi := i * rp.width
			highi := (i + 1) * rp.width
			if highi > textlen {
				highi = textlen
			}
			mp.MovePrint(linei, 0, text[lowi:highi])
		}
	}

	return linei
}
