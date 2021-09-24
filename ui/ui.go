package ui

import (
	"github.com/gbin/goncurses"
	"github.com/gpontesss/bib/bib"
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
	goncurses.SlkAttributeOn(goncurses.FO_WRAP)

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
	for {
		switch ui.stdscr.GetChar() {
		case 'q':
			return
		case 'k':
			ui.curpad().Scroll(-1)
			ui.curpad().NoutRefresh()
		case 'j':
			ui.curpad().Scroll(1)
			ui.curpad().NoutRefresh()
		case 'd':
			ref := bib.Ref{BookName: "Genesis", ChapterNum: 1, VerseNum: 1, Offset: 10}
			ui.curpad().LoadRef(&ref)
			ui.curpad().NoutRefresh()
		default:
		}
		goncurses.Update()
	}
}
