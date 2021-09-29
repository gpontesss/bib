package ui

import (
	"unicode/utf8"

	gc "github.com/gbin/goncurses"
	"github.com/gpontesss/bib/bib"
)

// RefPrinter docs here.
type RefPrinter struct {
	ref   *bib.Ref
	vsr   *bib.Version
	width int
	vrss  []*bib.Verse
}

// NewRefPrinter docs here.
func NewRefPrinter(ref *bib.Ref, vsr *bib.Version, width int) RefPrinter {
	return RefPrinter{
		ref:   ref,
		vsr:   vsr,
		width: width,
	}
}

// caches verses.
func (rp *RefPrinter) verses() []*bib.Verse {
	if rp.vrss == nil {
		rp.vrss = rp.ref.Verses(rp.vsr)
	}
	return rp.vrss
}

// LinesRequired docs here.
func (rp *RefPrinter) LinesRequired() int {
	linei := 0
	for _, verse := range rp.verses() {
		ref := verse.Ref()

		// RuneCountInString is used for compatibility with UTF-8 strings, for,
		// in some cases, len will return a number greater than desired.
		// +1 accounts for space.
		linelen := utf8.RuneCountInString(ref.String()) + 1

		wordsiter := NewWordIter(verse.Text)
		for wordsiter.Next() {
			// TODO: what if words are bigger than the max line length?
			word := wordsiter.Value()
			wordlen := utf8.RuneCountInString(word)

			if linelen+wordlen > rp.width {
				linelen = 0
				linei++
			}
			// +1 accounts for space.
			linelen += wordlen + 1
		}
		linei++
	}
	return linei
}

// Print docs here.
func (rp *RefPrinter) Print(pad *gc.Pad) int {
	linei := 0
	for _, verse := range rp.ref.Verses(rp.vsr) {
		ref := verse.Ref()
		refstr := ref.String()

		pad.AttrOn(gc.ColorPair(1) | gc.A_BOLD)
		pad.MovePrint(linei, 0, refstr)
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
			wordlen := utf8.RuneCountInString(word)

			if linelen+wordlen > rp.width {
				linelen = 0
				linei++
			}
			pad.MovePrint(linei, linelen, word)
			// +1 accounts for space
			linelen += wordlen + 1
		}
		linei++
	}

	return linei
}
