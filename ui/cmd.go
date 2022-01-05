package ui

import (
	"fmt"
	"time"

	gc "github.com/gbin/goncurses"
)

// CmdBox docs here.
type CmdBox struct {
	WinBox
	cursor XY
}

// NewCmdBox docs here.
func NewCmdBox(box Box) (CmdBox, error) {
	win, err := NewBoxWin(box)
	return CmdBox{win, XY{}}, err
}

// Exec docs here.
// TODO: overflow goes to lines bellow
func (cb *CmdBox) Exec() {
	cb.cursor = XY{0, 0}

	cb.MoveAddChar(cb.cursor.Y, cb.cursor.X, ':')
	cb.cursor = cb.cursor.Move(1, 0)
	cb.Refresh()

	var cmd string

loop:
	for {
		key := cb.GetChar()
		switch key {
		case gc.KEY_RETURN:
			break loop
		case gc.KEY_BACKSPACE:
			// Erase last character.
			cmd = cmd[:len(cmd)-1]
			cb.cursor = cb.cursor.Move(-1, 0)
			cb.MovePrint(cb.cursor.Y, cb.cursor.X, " ")
		default:
			strkey := fmt.Sprintf("%c", key)
			cmd = cmd + strkey
			cb.Print(strkey)
		}
		cb.Refresh()
	}

	cb.Erase()
	cb.Refresh()
	cb.MovePrint(0, 0, cmd)
	cb.Refresh()
	time.Sleep(1 * time.Second)
	cb.Erase()
	cb.Refresh()
}
