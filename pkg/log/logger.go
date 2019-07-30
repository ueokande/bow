package log

import (
	"sync"

	"github.com/fatih/color"
)

// Logger is a color logger with the host name
type Logger struct {
	host    string
	nohosts bool

	msgw  *color.Color
	hostw *color.Color
	m     *sync.Mutex
}

// Println prints the str as a single line with the host name
func (l *Logger) Println(str string) {
	l.m.Lock()
	defer l.m.Unlock()

	if l.nohosts {
		l.msgw.Println(str)
	} else {
		l.hostw.Add(color.Bold).Printf(l.host)
		l.msgw.Println("|", str)
	}
}
