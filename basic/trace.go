package basic

import "fmt"

type LogTrace struct {
	Trace *SafeSlice[string]
}

func NewLogTrace() *LogTrace {
	return &LogTrace{
		Trace: NewSafeSlice[string](100),
	}
}

func (p *LogTrace) LogErr(format string, a ...any) {
	p.Log("Error", format, a...)
}

func (p *LogTrace) LogWarn(format string, a ...any) {
	p.Log("Warn", format, a...)
}

func (p *LogTrace) LogOk(format string, a ...any) {
	p.Log("Ok", format, a...)
}

func (p *LogTrace) Log(level string, msg string, a ...any) {
	p.Trace.Append(fmt.Sprintf("[%v] ", level) + fmt.Sprintf(msg, a...))
}
