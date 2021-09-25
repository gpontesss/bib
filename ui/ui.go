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

	// Initializes cursor in right position
	curpad := ui.curpad()
	curpad.MoveCursor(curpad.miny(), curpad.minx())

	for {
		curpad := ui.curpad()
		curpad.Refresh()
		gc.Update()
		switch curpad.GetChar() {
		case 'q': // Quits.
			return
		case 'g': // Goes to top of text.
			curpad.GotoCursor(curpad.miny(), 0)
		case 'G': // Goes to bottom of text.
			curpad.GotoCursor(curpad.maxy(), 0)
		case '_':
			curpad.GotoCursor(curpad.cursory, curpad.minx())
		case '$':
			curpad.GotoCursor(curpad.cursory, curpad.maxx())
		case 'k': // Moves cursor up.
			curpad.MoveCursor(-1, 0)
		case 'j': // Moves cursor down.
			curpad.MoveCursor(1, 0)
		case 'h': // Moves cursor left.
			curpad.MoveCursor(0, -1)
		case 'l': // Moves cursor right.
			curpad.MoveCursor(0, 1)
		case 'u': // Moves cursor half-page up.
			curpad.Scroll(-curpad.height / 2)
		case 'd': // Moves cursor half-page down.
			curpad.Scroll(curpad.height / 2)
		case 'n': // Advances chapter.
			ref.ChapterNum++
			if ref.ChapterNum > 50 {
				ref.ChapterNum = 50
			}
			curpad.LoadRef(&ref)
		case 'p': // Retrocedes chapter.
			ref.ChapterNum--
			if ref.ChapterNum < 1 {
				ref.ChapterNum = 1
			}
			curpad.LoadRef(&ref)
		case 'L': // Changes pad focus to the one on the right.
			ui.nextpad()
			ui.curpad().GotoCursor(ui.curpad().cursory, ui.curpad().cursorx)
		case 'H': // Changes pad focus to the one on the left.
			ui.prevpad()
			ui.curpad().GotoCursor(ui.curpad().cursory, ui.curpad().cursorx)
		}
	}
}
