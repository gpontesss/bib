package ui

import (
	"unicode/utf8"

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
	width := rp.width - (2 * rp.padding)

	mp.MovePrint(vertpadding, rp.padding, rp.vsr.Name, " ", rp.ref) // header

	linei := vertpadding + 1
	for _, verse := range rp.ref.Verses(rp.vsr) {
		text := verse.String()

		linelen := 0
		wordsiter := NewWordIter(text)

		for wordsiter.Next() {
			// TODO: what if words are bigger than the max line length?
			word := wordsiter.Value()

			// RuneCountInString is used for compatibility with UTF-8 strings,
			// for, in some cases, len will return a number greater than
			// desired.
			if linelen+utf8.RuneCountInString(word) > width {
				linelen = 0
				linei++
			}
			mp.MovePrint(linei, rp.padding+linelen, word, " ")
			// +1 accounts for space
			linelen += utf8.RuneCountInString(word) + 1
		}
		linei++
	}

	return linei
}
