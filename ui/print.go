package ui

import (
	"unicode/utf8"

	gc "github.com/gbin/goncurses"
	"github.com/gpontesss/bib/bib"
)

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
func (rp *RefPrinter) Print(pad *gc.Pad) int {
	// height/width ration ~ 2
	vertpadding := (rp.padding / 2)
	width := rp.width - (2 * rp.padding)

	linei := vertpadding
	for _, verse := range rp.ref.Verses(rp.vsr) {
		ref := verse.Ref()
		refstr := ref.String()

		pad.AttrOn(gc.ColorPair(1) | gc.A_BOLD)
		pad.MovePrint(linei, rp.padding, refstr)
		pad.AttrOff(gc.ColorPair(1) | gc.A_BOLD)

		// RuneCountInString is used for compatibility with UTF-8 strings,
		// for, in some cases, len will return a number greater than
		// desired.
		// +1 accounts for space
		linelen := utf8.RuneCountInString(refstr) + 1
		wordsiter := NewWordIter(verse.Text)

		for wordsiter.Next() {
			// TODO: what if words are bigger than the max line length?
			word := wordsiter.Value()

			if linelen+utf8.RuneCountInString(word) > width {
				linelen = 0
				linei++
			}
			pad.MovePrint(linei, rp.padding+linelen, word)
			// +1 accounts for space
			linelen += utf8.RuneCountInString(word) + 1
		}
		linei++
	}

	return linei
}
