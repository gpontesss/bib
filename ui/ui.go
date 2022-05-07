package ui

import (
	"os"
	"os/signal"
	"syscall"

	gc "github.com/gbin/goncurses"
	"github.com/gpontesss/bib/bib"
	"golang.org/x/term"
)

// UI docs here.
type UI struct {
	WinBox
	vsrmenu  VersionMenu
	cmdbox   CmdBox
	padding  uint
	curpadi  int // by default, first pad is selected. (unitiated)
	Versions []*bib.Version
	pads     []VersionPad
}

// Init docs here.
func (ui *UI) Init() (err error) {
	stdsrc, err := gc.Init()
	if err != nil {
		return err
	}
	// root UI is anchored at (0,0)
	maxheight, maxwidth := stdsrc.MaxYX()
	ui.WinBox = CastWinBox(Box{XY{0, 0}, uint(maxheight), uint(maxwidth)}, stdsrc)
	ui.padding = 1

	ui.pads = make([]VersionPad, len(ui.Versions))
	for i := range ui.pads {
		vsr := ui.Versions[i]
		if ui.pads[i], err = NewVersionPad(vsr, MinBox(), ui.padding); err != nil {
			return
		}
	}
	if ui.cmdbox, err = NewCmdBox(MinBox()); err != nil {
		return
	}
	if ui.vsrmenu, err = NewVersionMenu(ui.Box, ui.Versions...); err != nil {
		return
	}

	ui.Resize(ui.height, ui.width)

	gc.StartColor() // Allows colors.
	gc.Cursor(1)    // Shows cursor.
	gc.Echo(false)  // Does not echo typing.

	gc.InitPair(1, gc.C_WHITE, 0)         // Verse highlighting
	gc.InitPair(2, gc.C_WHITE, gc.C_BLUE) // Header

	return nil
}

func (ui *UI) nextpad() { ui.curpadi = (ui.curpadi + 1) % len(ui.pads) }
func (ui *UI) prevpad() {
	ui.curpadi = ((ui.curpadi-1)%len(ui.pads) + len(ui.pads)) % len(ui.pads)
}
func (ui *UI) curpad() *VersionPad { return &ui.pads[ui.curpadi] }

// End docs here.
func (ui *UI) End() {
	defer gc.End()

	for padi := range ui.pads {
		(&ui.pads[padi]).Delete()
	}
	ui.cmdbox.Delete()
	ui.WinBox.Delete()
}

// Refresh docs here.
func (ui *UI) Refresh(all bool) {
	if all {
		for i := range ui.pads {
			ui.pads[i].NoutRefresh()
		}
	} else {
		ui.curpad().NoutRefresh()
	}
	gc.Update()
}

// Resize docs here.
func (ui *UI) Resize(height, width uint) {
	// TODO: handle error.
	_ = gc.ResizeTerm(int(ui.height), int(ui.width))
	ui.WinBox.Resize(height, width)

	// Gets rid of previously painted columns that were part of pads, but no
	// longer are.
	ui.WinBox.Erase()
	ui.WinBox.NoutRefresh()

	// ui.vsrmenu.ResizeBox(ui.Box)

	boxiter := ui.VertDiv(uint(len(ui.Versions)))
	for boxiter.Next() {
		ui.pads[boxiter.Index()].Resize(boxiter.Value(), ui.padding)
	}
	ui.cmdbox.ResizeBox(Box{ui.SW().Move(0, -1), 1, ui.width})
}

// IncrPadding docs here.
func (ui *UI) IncrPadding(padding int) {
	ui.padding = uint(max(0, int(ui.padding)+padding))
	ui.Resize(ui.height, ui.width)
}

// HandleKey docs here.
// TODO: delegate key handling to nested windows to better segregate logic.
func (ui *UI) HandleKey(key gc.Key) bool {
	curpad := ui.curpad()
	switch key {
	case 'q': // Quits.
		return true
	case ':':
		ui.cmdbox.Exec()
		ui.Refresh(true)
	case 'g': // Goes to top of text.
		curpad.GotoCursor(0, uint(curpad.cursor.X))
	case 'G': // Goes to bottom of text.
		curpad.GotoCursor(curpad.pad.bufheight, uint(curpad.cursor.X))
	case '_':
		curpad.GotoCursor(uint(curpad.cursor.Y), 0)
	case '$':
		curpad.GotoCursor(uint(curpad.cursor.Y), curpad.pad.bufwidth)
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
		ui.curpad().GotoCursor(uint(ui.curpad().cursor.Y), uint(ui.curpad().cursor.X))
	case 'H': // Changes pad focus to the one on the left.
		ui.prevpad()
		ui.curpad().GotoCursor(uint(ui.curpad().cursor.Y), uint(ui.curpad().cursor.X))
	case gc.KEY_TAB:
		if vsr, err := ui.vsrmenu.Select(); err == nil && vsr != nil {
			curpad.SetVersion(vsr)
			// Refreshes reference to show updated version's text.
			curpad.LoadRef(&curpad.refloaded)
		}
		// always refreshes all for removing menu window "shadow".
		ui.Refresh(true)
		// and moves cursor to where it should be, in the active pad.
		ui.curpad().GotoCursor(uint(ui.curpad().cursor.Y), uint(ui.curpad().cursor.X))
	}
	return false
}

// Loop docs here.
func (ui *UI) Loop() (err error) {
	// Initially loads reference.
	var ref bib.Ref
	if ref, err = bib.ParseRef("John 1"); err != nil {
		return err
	}
	for i := range ui.pads {
		(&ui.pads[i]).LoadRef(&ref)
	}

	// Initializes cursor in right position
	ui.curpad().MoveCursor(0, 0)

	winchchan := make(chan os.Signal, 1)
	signal.Notify(winchchan, syscall.SIGWINCH)
	defer signal.Stop(winchchan)

	// handles asynchronously terminal resizing.
	go func() {
		// TODO: maybe use some lock to prevents weird screen updates.
		for range winchchan {
			width, height, err := term.GetSize(0)
			if err != nil {
				panic(err)
			}
			// safe cast for terminal dimensions cannot be negative
			ui.Resize(uint(height), uint(width))
			ui.Refresh(true)
		}
	}()

	for {
		ui.Refresh(false)
		// TODO: better exiting handler.
		if ui.HandleKey(ui.curpad().GetChar()) {
			return nil
		}
	}
}

// AsyncLoop docs here.
func (ui *UI) AsyncLoop() <-chan error {
	loopend := make(chan error)
	go func() { loopend <- ui.Loop() }()
	return loopend
}

// Close should be deferred before initiating the UI, so that in all cases it is
// ended before any error logging or such, to avoid terminal glitches.
func (ui *UI) Close() {
	err := recover()
	// forces UI to be ended before logging, to avoid terminal bugging
	// glitch.
	ui.End()
	if err != nil {
		panic(err)
	}
}
