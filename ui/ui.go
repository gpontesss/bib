package ui

import (
	"github.com/gpontesss/bib/bib"
	gc "github.com/rthornton128/goncurses"
)

// UI docs here.
type UI struct {
	Versions []bib.Version
	stdscr   *gc.Window
	pads     []VersionPad
	curpadi  int
}

// Init docs here.
func (ui *UI) Init() error {
	var err error
	if ui.stdscr, err = gc.Init(); err != nil {
		return err
	}
	gc.Cursor(1)   // Shows cursor.
	gc.Echo(false) // Does not echo typing.

	rows, cols := ui.stdscr.MaxYX()
	padheight := rows
	padwidth := cols / len(ui.Versions)

	ui.pads = make([]VersionPad, 0, len(ui.Versions))
	for i := range ui.Versions {
		vsr := &ui.Versions[i]
		pad, err := NewVersionPad(
			vsr,
			padheight, padwidth,
			0, i*padwidth,
			1)
		if err != nil {
			return err
		}
		ui.pads = append(ui.pads, pad)
	}

	// by default, first pad is selected.
	ui.curpadi = 0
	return nil
}

func (ui *UI) nextpad() {
	ui.curpadi++
	if ui.curpadi >= len(ui.pads) {
		ui.curpadi = 0
	}
}

func (ui *UI) prevpad() {
	ui.curpadi--
	if ui.curpadi < 0 {
		ui.curpadi = len(ui.pads) - 1
	}
}

func (ui *UI) curpad() *VersionPad { return &ui.pads[ui.curpadi] }

// End docs here.
func (ui *UI) End() {
	for _, pad := range ui.pads {
		pad.Delete()
	}
	gc.End()
}

// Loop docs here.
func (ui *UI) Loop() {
	// Initially loads reference.
	ref := bib.Ref{BookName: "Genesis", ChapterNum: 1, VerseNum: 1, Offset: 100}
	for i := range ui.pads {
		pad := &ui.pads[i]
		pad.LoadRef(&ref)
	}

	for {
		ui.curpad().Refresh()
		gc.Update()
		switch ui.curpad().GetChar() {
		case 'q': // Quits.
			return
		case 'g': // Goes to top of text.
			ui.curpad().Goto(0)
		case 'G': // Goes to bottom of text.
			ui.curpad().Goto(ui.curpad().maxoffset)
		case 'k': // Move cursor up.
			ui.curpad().MoveCursor(-1, 0)
			// ui.curpad().Scroll(-1)
		case 'j': // Move cursor down.
			ui.curpad().MoveCursor(1, 0)
			// ui.curpad().Scroll(1)
		case 'h': // Move cursor left.
			ui.curpad().MoveCursor(0, -1)
		case 'l': // Move cursor right.
			ui.curpad().MoveCursor(0, 1)
		case 'u': // Move cursor a bunch up.
			ui.curpad().Scroll(-10)
			ui.curpad().MoveCursor(-10, 0)
		case 'd': // Move cursor a bunch down.
			ui.curpad().Scroll(10)
			ui.curpad().MoveCursor(10, 0)
		case 'n': // Advances chapter.
			ref.ChapterNum++
			if ref.ChapterNum > 50 {
				ref.ChapterNum = 50
			}
			ui.curpad().LoadRef(&ref)
		case 'p': // Retrocedes chapter.
			ref.ChapterNum--
			if ref.ChapterNum < 1 {
				ref.ChapterNum = 1
			}
			ui.curpad().LoadRef(&ref)
		case 'L': // Changes pad focus to the one on the right.
			ui.nextpad()
			ui.curpad().GotoCursor(ui.curpad().cursory, ui.curpad().cursorx)
		case 'H': // Changes pad focus to the one on the left.
			ui.prevpad()
			ui.curpad().GotoCursor(0, 0)
		}
	}
}
