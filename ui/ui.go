package ui

import (
	"github.com/gbin/goncurses"
	"github.com/gpontesss/bib/bib"
)

// UI docs here.
type UI struct {
	Versions []bib.Version
	stdscr   *goncurses.Window
	wins     map[*bib.Version]*goncurses.Window
	curvsri  int
}

// Init docs here.
func (ui *UI) Init() error {
	var err error
	if ui.stdscr, err = goncurses.Init(); err != nil {
		return err
	}
	goncurses.Cursor(1)   // Shows cursor.
	goncurses.Echo(false) // Does not echo typing.
	goncurses.SlkAttributeOff(goncurses.FO_WRAP)

	rows, cols := ui.stdscr.MaxYX()
	winheight := rows
	winwidth := cols / len(ui.Versions)

	ui.wins = map[*bib.Version]*goncurses.Window{}
	for i := range ui.Versions {
		win, err := goncurses.NewWindow(
			winheight, winwidth,
			0, winwidth*i,
		)
		if err != nil {
			return err
		}
		win.ScrollOk(true)
		ui.wins[&ui.Versions[i]] = win
	}

	// by default, first version is selected.
	ui.curvsri = 0
	return nil
}

func (ui *UI) nextvsr()                  { ui.curvsri++ }
func (ui *UI) curvsr() *bib.Version      { return &ui.Versions[ui.curvsri] }
func (ui *UI) curwin() *goncurses.Window { return ui.wins[ui.curvsr()] }

// End docs here.
func (ui *UI) End() {
	for _, win := range ui.wins {
		win.Delete()
	}
	goncurses.End()
}

func (ui *UI) Loop() {
	for {
		switch ui.stdscr.GetChar() {
		case 'q':
			return
		case 'k':
			ui.curwin().Scroll(-1)
			ui.curwin().NoutRefresh()
		case 'j':
			ui.curwin().Scroll(1)
			ui.curwin().NoutRefresh()
		case 'd':
			// initially writes a bunch to windows
			for vsr, win := range ui.wins {
				for i := 0; i < 20; i++ {
					win.Println(&vsr.Books[0].Chapters[0].Verses[i])
					win.NoutRefresh()
				}
			}
		default:
		}
		goncurses.Update()
	}
}
