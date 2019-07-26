package main

import (
	"sync"

	"github.com/fatih/color"
)

type Logger struct {
	host    string
	nohosts bool

	msgw  *color.Color
	hostw *color.Color
	m     *sync.Mutex
}

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

type LoggerFactory struct {
	nohosts bool

	color color.Attribute
	m     sync.Mutex
}

func (f *LoggerFactory) NewLogger(host string) *Logger {
	f.color = f.nextColor()
	return &Logger{
		host:    host,
		nohosts: f.nohosts,
		hostw:   color.New(f.color, color.Bold),
		msgw:    color.New(f.color),
		m:       &f.m,
	}
}

func (f *LoggerFactory) nextColor() color.Attribute {

	// skip the following colors
	// - color.FgWhite
	// - color.FgHiBlack
	// - color.FgHiWhite
	// - color.FgBlack

	switch f.color {
	case color.FgRed:
		return color.FgGreen
	case color.FgGreen:
		return color.FgYellow
	case color.FgYellow:
		return color.FgBlue
	case color.FgBlue:
		return color.FgMagenta
	case color.FgMagenta:
		return color.FgCyan
	case color.FgCyan:
		return color.FgHiRed
	case color.FgHiRed:
		return color.FgHiGreen
	case color.FgHiGreen:
		return color.FgHiYellow
	case color.FgHiYellow:
		return color.FgHiBlue
	case color.FgHiBlue:
		return color.FgHiMagenta
	case color.FgHiMagenta:
		return color.FgHiCyan
	case color.FgHiCyan:
		return color.FgRed
	}
	return color.FgRed
}
