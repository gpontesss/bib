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
	ref            *bib.Ref
	vsr            *bib.Version
	width, padding int
}

// NewRefPrinter docs here.
func NewRefPrinter(ref *bib.Ref, vsr *bib.Version, width, padding int) RefPrinter {
	return RefPrinter{
		ref:     ref,
		vsr:     vsr,
		width:   width,
		padding: padding,
	}
}

// SetVersion docs here.
func (rp *RefPrinter) SetVersion(vsr *bib.Version) { rp.vsr = vsr }

// Print docs here.
func (rp *RefPrinter) Print(mp MovePrinter) int {
	// height/width ration ~ 2
	vertpadding := (rp.padding / 2)
	mp.MovePrint(vertpadding, rp.padding, rp.vsr.Name, " ", rp.ref) // header

	linei := vertpadding + 1
	for _, verse := range rp.ref.Verses(rp.vsr) {
		text := verse.String()
		textlen := len(text)

		width := rp.width - (2 * rp.padding)
		// A more secure way of doing ceiling division with integers.
		linecount := textlen / width
		if (textlen % width) != 0 {
			linecount++
		}

		for i := 0; i < linecount; i++ {
			linei++
			// TODO: padding?
			lowi := i * width
			highi := (i + 1) * width
			if highi > textlen {
				highi = textlen
			}
			mp.MovePrint(linei, rp.padding, text[lowi:highi])
		}
	}

	return linei
}
