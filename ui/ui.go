package ui

import (
	"github.com/gpontesss/bib/bib"
	"github.com/rthornton128/goncurses"
)

// UI docs here.
type UI struct {
	Versions []bib.Version
	stdscr   *goncurses.Window
	pads     []VersionPad
	curpadi  int
}

// Init docs here.
func (ui *UI) Init() error {
	var err error
	if ui.stdscr, err = goncurses.Init(); err != nil {
		return err
	}
	goncurses.Cursor(1)   // Shows cursor.
	goncurses.Echo(false) // Does not echo typing.

	rows, cols := ui.stdscr.MaxYX()
	padheight := rows
	padwidth := cols / len(ui.Versions)

	ui.pads = make([]VersionPad, 0, len(ui.Versions))
	for i := range ui.Versions {
		vsr := &ui.Versions[i]
		pad, err := NewVersionPad(
			vsr,
			padheight, padwidth,
			0, i*padwidth)
		if err != nil {
			return err
		}
		ui.pads = append(ui.pads, pad)
	}

	// by default, first pad is selected.
	ui.curpadi = 0
	return nil
}

func (ui *UI) nextpad()            { ui.curpadi++ }
func (ui *UI) curpad() *VersionPad { return &ui.pads[ui.curpadi] }

// End docs here.
func (ui *UI) End() {
	for _, pad := range ui.pads {
		pad.Delete()
	}
	goncurses.End()
}

// Loop docs here.
func (ui *UI) Loop() {
	// Initially loads reference.
	ref := bib.Ref{BookName: "Genesis", ChapterNum: 1, VerseNum: 1, Offset: 20}
	ui.curpad().LoadRef(&ref)

	for {
		ui.curpad().Refresh()
		goncurses.Update()
		switch ui.curpad().GetChar() {
		case 'q': // Quits.
			return
		case 'g': // Goes to top of text.
			ui.curpad().Goto(0)
			ui.curpad().NoutRefresh()
		case 'G': // Goes to bottom of text.
			ui.curpad().Goto(ui.curpad().maxoffset)
			ui.curpad().NoutRefresh()
		case 'k': // Scrolls up.
			ui.curpad().Scroll(-1)
			ui.curpad().NoutRefresh()
		case 'j': // Scrolls down.
			ui.curpad().Scroll(1)
			ui.curpad().NoutRefresh()
		case 'u': // Scrolls a bunch up.
			ui.curpad().Scroll(-10)
			ui.curpad().NoutRefresh()
		case 'd': // Scrolls a bunch down.
			ui.curpad().Scroll(10)
			ui.curpad().NoutRefresh()
		case 'l': // Advances chapter.
			ref.ChapterNum++
			if ref.ChapterNum > 50 {
				ref.ChapterNum = 50
			}
			ui.curpad().LoadRef(&ref)
			ui.curpad().NoutRefresh()
		case 'h': // Retrocedes chapter.
			ref.ChapterNum--
			if ref.ChapterNum < 1 {
				ref.ChapterNum = 1
			}
			ui.curpad().LoadRef(&ref)
			ui.curpad().NoutRefresh()
		}
	}
}
