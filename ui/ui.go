package ui

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	gc "github.com/gbin/goncurses"
	"github.com/gpontesss/bib/bib"
	"golang.org/x/term"
)

// UI docs here.
type UI struct {
	Versions []*bib.Version
	// TODO: can't it be ephemeral?
	vsrmenu   VersionMenu
	stdscr    *gc.Window
	padding   uint
	box       Box
	winchchan chan os.Signal
	keychan   chan gc.Key
	// by default, first pad is selected. (unitiated)
	curpadi int
	pads    []VersionPad
}

// Init docs here.
func (ui *UI) Init() (err error) {
	if ui.stdscr, err = gc.Init(); err != nil {
		return err
	}
	defer func() {
		if val := recover(); val != nil {
			err = fmt.Errorf("%v", val)
		}
	}()

	gc.StartColor() // Allows colors.
	gc.Cursor(1)    // Shows cursor.
	gc.Echo(false)  // Does not echo typing.

	gc.InitPair(1, gc.C_WHITE, 0)         // Verse highlighting
	gc.InitPair(2, gc.C_WHITE, gc.C_BLUE) // Header

	// root UI is anchored at (0,0)
	ui.box = Box{XY{0, 0}, 0, 0}

	maxheight, maxwidth := ui.stdscr.MaxYX()
	// cast is safe for term size can never be negative.
	ui.box = ui.box.Resize(uint(maxheight), uint(maxwidth))
	ui.padding = 1

	ui.pads = make([]VersionPad, len(ui.Versions))

	boxiter := ui.box.VertDiv(uint(len(ui.pads)))
	for boxiter.Next() {
		i := boxiter.Index()
		ui.pads[i], err = NewVersionPad(
			ui.Versions[i], boxiter.Value(), ui.padding)
		if err != nil {
			return err
		}
	}

	// TODO: resize it too.
	if ui.vsrmenu, err = NewVersionMenu(
		int(ui.box.height/4), int(ui.box.width/4),
		int(ui.box.height/2), int(ui.box.width/2),
		ui.Versions...,
	); err != nil {
		return err
	}

	ui.winchchan = make(chan os.Signal, 1)
	signal.Notify(ui.winchchan, syscall.SIGWINCH)

	ui.keychan = make(chan gc.Key)
	go func() {
		for {
			ui.keychan <- ui.curpad().GetChar()
		}
	}()

	return nil
}

func (ui *UI) nextpad()            { ui.curpadi = (ui.curpadi + 1) % len(ui.pads) }
func (ui *UI) prevpad()            { ui.curpadi = (ui.curpadi - 1) % len(ui.pads) }
func (ui *UI) curpad() *VersionPad { return &ui.pads[ui.curpadi] }

// End docs here.
func (ui *UI) End() {
	// TODO: include version menu.
	for padi := range ui.pads {
		(&ui.pads[padi]).Delete()
	}
	gc.End()
}

// Refresh docs here.
func (ui *UI) Refresh(all bool) {
	if all {
		for i := range ui.pads {
			(&ui.pads[i]).NoutRefresh()
		}
	} else {
		ui.curpad().NoutRefresh()
	}
	gc.Update()
}

// Resize docs here.
func (ui *UI) Resize(height, width uint) {
	ui.box = ui.box.Resize(height, width)

	gc.ResizeTerm(int(ui.box.height), int(ui.box.width))
	ui.stdscr.Resize(int(ui.box.height), int(ui.box.width))

	// Gets rid of previously painted columns that were part of pads, but no
	// longer are.
	ui.stdscr.Erase()
	ui.stdscr.NoutRefresh()

	boxiter := ui.box.VertDiv(uint(len(ui.Versions)))
	for boxiter.Next() {
		ui.pads[boxiter.Index()].Resize(boxiter.Value(), ui.padding)
	}
}

// IncrPadding docs here.
func (ui *UI) IncrPadding(padding int) {
	ui.padding = uint(max(0, int(ui.padding)+padding))
	ui.Resize(ui.box.height, ui.box.width)
}

// HandleKey docs here.
func (ui *UI) HandleKey(key gc.Key) bool {
	curpad := ui.curpad()
	switch key {
	case 'q': // Quits.
		return true
	case 'g': // Goes to top of text.
		curpad.GotoCursor(curpad.miny(), curpad.minx())
	case 'G': // Goes to bottom of text.
		curpad.GotoCursor(curpad.maxy(), curpad.minx())
	case '_':
		curpad.GotoCursor(curpad.cursor.Y, curpad.minx())
	case '$':
		curpad.GotoCursor(curpad.cursor.Y, curpad.maxx())
	case '(':
		ui.IncrPadding(-1)
	case ')':
		ui.IncrPadding(1)
	case 'k': // Moves cursor up.
		curpad.MoveCursor(-1, 0)
	case 'K': // Scrolls 1 row up.
		curpad.Scroll(-1)
	case 'j': // Moves cursor down.
		curpad.MoveCursor(1, 0)
	case 'J': // Scrolls 1 row down.
		curpad.Scroll(1)
	case 'h': // Moves cursor left.
		curpad.MoveCursor(0, -1)
	case 'l': // Moves cursor right.
		curpad.MoveCursor(0, 1)
	case 'u': // Moves cursor half-page up.
		curpad.Scroll(-int(curpad.height) / 2)
	case 'd': // Moves cursor half-page down.
		curpad.Scroll(int(curpad.height) / 2)
	case 'n': // Advances chapter.
		if next := curpad.RefLoaded().Chapter(curpad.vsr).Next(); next != nil {
			ref := next.Ref()
			curpad.LoadRef(&ref)
		}
	case 'N': // Advances book.
		if next := curpad.RefLoaded().Book(curpad.vsr).Next(); next != nil {
			ref := next.FirstChapter().Ref()
			curpad.LoadRef(&ref)
		}
	case 'p': // Retrocedes chapter.
		if prev := curpad.RefLoaded().Chapter(curpad.vsr).Previous(); prev != nil {
			ref := prev.Ref()
			curpad.LoadRef(&ref)
		}
	case 'P': // Retrocedes book.
		if prev := curpad.RefLoaded().Book(curpad.vsr).Previous(); prev != nil {
			ref := prev.FirstChapter().Ref()
			curpad.LoadRef(&ref)
		}
	case 'L': // Changes pad focus to the one on the right.
		ui.nextpad()
		ui.curpad().GotoCursor(ui.curpad().cursor.Y, ui.curpad().cursor.X)
	case 'H': // Changes pad focus to the one on the left.
		ui.prevpad()
		ui.curpad().GotoCursor(ui.curpad().cursor.Y, ui.curpad().cursor.X)
	case gc.KEY_TAB:
		if vsr, err := ui.vsrmenu.Select(); err == nil && vsr != nil {
			curpad.SetVersion(vsr)
			// Refreshes reference to show updated version's text.
			curpad.LoadRef(&curpad.refloaded)
		}
		// always refreshes all for removing menu window "shadow".
		ui.Refresh(true)
		// and moves cursor to where it should be, in the active pad.
		ui.curpad().GotoCursor(ui.curpad().cursor.Y, ui.curpad().cursor.X)
	}
	return false
}

// Loop docs here.
func (ui *UI) Loop() error {
	// Initially loads reference.
	if ref, err := bib.ParseRef("John 1"); err != nil {
		return err
	} else {
		for i := range ui.pads {
			(&ui.pads[i]).LoadRef(&ref)
		}
	}

	// Initializes cursor in right position
	curpad := ui.curpad()
	curpad.MoveCursor(int(curpad.miny()), int(curpad.minx()))

	for {
		ui.Refresh(false)
		select {
		case <-ui.winchchan:
			width, height, err := term.GetSize(0)
			if err != nil {
				panic(err)
			}
			// safe cast for terminal dimensions cannot be negative
			ui.Resize(uint(height), uint(width))
		case key := <-ui.keychan:
			// TODO: better exiting handler.
			if ui.HandleKey(key) {
				return nil
			}
		}
	}
}

// AsyncLoop docs here.
func (ui *UI) AsyncLoop() <-chan error {
	loopend := make(chan error)
	go func() { loopend <- ui.Loop() }()
	return loopend
}
