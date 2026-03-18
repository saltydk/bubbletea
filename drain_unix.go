//go:build linux || solaris || aix
// +build linux solaris aix

package tea

import "golang.org/x/sys/unix"

// drainInput discards any pending input on the TTY (e.g. terminal responses to
// DECRQM mode queries or kitty keyboard queries) that arrived after the input
// reader was cancelled during shutdown. Without this, those responses appear as
// garbage characters in the user's shell after the program exits.
//
// Terminal responses may arrive in bursts (e.g. three separate responses to
// mode 2026, mode 2027, and kitty keyboard queries) so we loop: flush whatever
// is in the buffer, poll for more, repeat until nothing arrives within the
// timeout. On local terminals the first flush handles everything and poll
// returns immediately. Over SSH the loop accommodates round-trip latency.
func (p *Program) drainInput() {
	if p.ttyInput == nil {
		return
	}
	fd := int(p.ttyInput.Fd())
	fds := []unix.PollFd{{Fd: int32(fd), Events: unix.POLLIN}}

	for {
		// Discard any data in the input buffer.
		_ = unix.IoctlSetInt(fd, unix.TCFLSH, 0)

		// Wait for more data. The timeout accommodates SSH round-trip
		// latency. For local terminals this returns immediately.
		n, _ := unix.Poll(fds, 200)
		if n <= 0 {
			return
		}
	}
}
