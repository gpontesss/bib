package ui

import (
	gc "github.com/gbin/goncurses"
	"github.com/gpontesss/bib/bib"
)

// UI docs here.
type UI struct {
	Versions []*bib.Version
	vsrmenu  VersionMenu
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
	for i, vsr := range ui.Versions {
		pad, err := NewVersionPad(
			vsr,
			// for some reason, if height is used integrally, nothing is
			// rendered.
			padheight-1, padwidth,
			0, i*padwidth,
			1)
		if err != nil {
			return err
		}
		ui.pads = append(ui.pads, pad)
	}

	// by default, first pad is selected.
	ui.curpadi = 0

	ui.vsrmenu, err = NewVersionMenu(
		rows/4, cols/4,
		rows/2, cols/2,
		ui.Versions...,
	)
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

// Refresh docs here.
func (ui *UI) Refresh(all bool) {
	if all {
		for i := range ui.pads {
			pad := &ui.pads[i]
			pad.NoutRefresh()
		}
	} else {
		ui.curpad().NoutRefresh()
	}
	gc.Update()
}

// Loop docs here.
func (ui *UI) Loop() {
	// Initially loads reference.
	ref := bib.Ref{BookName: "Genesis", ChapterNum: 1}
	for i := range ui.pads {
		pad := &ui.pads[i]
		pad.LoadRef(&ref)
	}

	// Initializes cursor in right position
	curpad := ui.curpad()
	curpad.MoveCursor(curpad.miny(), curpad.minx())

	for {
		ui.Refresh(false)
		curpad := ui.curpad()
		switch curpad.GetChar() {
		case 'q': // Quits.
			return
		case 'g': // Goes to top of text.
			curpad.GotoCursor(curpad.miny(), curpad.minx())
		case 'G': // Goes to bottom of text.
			curpad.GotoCursor(curpad.maxy(), curpad.minx())
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
			if next := curpad.RefLoaded().Chapter(curpad.vsr).Next(); next != nil {
				ref := next.Ref()
				curpad.LoadRef(&ref)
			}
		case 'p': // Retrocedes chapter.
			if prev := curpad.RefLoaded().Chapter(curpad.vsr).Previous(); prev != nil {
				ref := prev.Ref()
				curpad.LoadRef(&ref)
			}
		case 'L': // Changes pad focus to the one on the right.
			ui.nextpad()
			ui.curpad().GotoCursor(ui.curpad().cursory, ui.curpad().cursorx)
		case 'H': // Changes pad focus to the one on the left.
			ui.prevpad()
			ui.curpad().GotoCursor(ui.curpad().cursory, ui.curpad().cursorx)
		case gc.KEY_TAB:
			if vsr, err := ui.vsrmenu.Select(); err == nil && vsr != nil {
				curpad.SetVersion(vsr)
				// Refreshes reference to show updated version's text.
				curpad.LoadRef(&curpad.refloaded)
			}
			// always refreshes all for removing menu window "shadow".
			ui.Refresh(true)
			// and moves cursor to where it should be, in the active pad.
			ui.curpad().GotoCursor(ui.curpad().cursory, ui.curpad().cursorx)
		}
	}
}
