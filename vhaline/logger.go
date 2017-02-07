package vhaline

import (
	"fmt"
)

// dlog is for delux debug logging, automatically including
// pid, Id, and timestamp. shows up at -debug level.
func (m *Replica) dlog(format string, args ...interface{}) {

	m.Cfg.verbmutex.Lock()
	defer m.Cfg.verbmutex.Unlock()

	if m.Cfg.Verbosity >= DEBUG {
		m.Logger.Output(2,
			fmt.Sprintf(
				fmt.Sprintf("[pid %v] replica (%s) %s", m.Pid, m.Me.ShortId(), format), args...))
	}
}

// -info/INFO log level; these also show up at DEBUG level
func (m *Replica) ilog(format string, args ...interface{}) {

	m.Cfg.verbmutex.Lock()
	defer m.Cfg.verbmutex.Unlock()

	if m.Cfg.Verbosity >= INFO {
		m.Logger.Output(2,
			fmt.Sprintf(
				fmt.Sprintf("[pid %v] replica (%s) %s", m.Pid, m.Me.ShortId(), format), args...))
	}
}

// always/all log
func (m *Replica) alog(format string, args ...interface{}) {

	m.Cfg.verbmutex.Lock()
	defer m.Cfg.verbmutex.Unlock()

	m.Logger.Output(2,
		fmt.Sprintf(
			fmt.Sprintf("[pid %v] replica (%s) %s", m.Pid, m.Me.ShortId(), format), args...))
}
