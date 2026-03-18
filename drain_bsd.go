//go:build darwin || dragonfly || freebsd || netbsd || openbsd
// +build darwin dragonfly freebsd netbsd openbsd

package tea

import "golang.org/x/sys/unix"

// drainInput discards any pending input on the TTY (e.g. terminal responses to
// DECRQM mode queries or kitty keyboard queries) that arrived after the input
// reader was cancelled during shutdown. Without this, those responses appear as
// garbage characters in the user's shell after the program exits.
func (p *Program) drainInput() {
	if p.ttyInput == nil {
		return
	}
	// TIOCFLUSH with FREAD (1) discards data received but not yet read.
	_ = unix.IoctlSetPointerInt(int(p.ttyInput.Fd()), unix.TIOCFLUSH, 1)
}
