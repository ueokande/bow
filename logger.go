package main

import (
	"fmt"
	"sync"
)

type Logger struct {
	host string

	m *sync.Mutex
}

func (l *Logger) Println(str string) {
	l.m.Lock()
	defer l.m.Unlock()

	fmt.Printf("[%s] %s\n", l.host, str)
}

type LoggerFactory struct {
	m sync.Mutex
}

func (f *LoggerFactory) NewLogger(host string) *Logger {
	return &Logger{
		host: host,
		m:    &f.m,
	}
}
