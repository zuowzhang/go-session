package session

type Log interface {
	D(format string, ...interface{})
	I(format string, ...interface{})
	W(format string, ...interface{})
	E(format string, ...interface{})
}

type LogProxy struct {
	impl Log
}

func (proxy *LogProxy)D(format string, args ...interface{}) {
	if proxy.impl != nil {
		proxy.impl.D(format, args...)
	}
}

func (proxy *LogProxy)I(format string, args ...interface{}) {
	if proxy.impl != nil {
		proxy.impl.I(format, args...)
	}
}

func (proxy *LogProxy)W(format string, args ...interface{}) {
	if proxy.impl != nil {
		proxy.impl.W(format, args...)
	}
}

func (proxy *LogProxy)E(format string, args ...interface{}) {
	if proxy.impl != nil {
		proxy.impl.E(format, args...)
	}
}