//go:build windows
// +build windows

package tea

import "golang.org/x/sys/windows"

// drainInput discards any pending console input events to prevent terminal
// responses from appearing as garbage characters after the program exits.
func (p *Program) drainInput() {
	if p.ttyInput == nil {
		return
	}
	_ = windows.FlushConsoleInputBuffer(windows.Handle(p.ttyInput.Fd()))
}
