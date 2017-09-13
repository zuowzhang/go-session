package session

type Log interface {
	D(string, ...interface{})
	I(string, ...interface{})
	W(string, ...interface{})
	E(string, ...interface{})
}

type logProxy struct {
	impl Log
}

func (proxy *logProxy)D(format string, args ...interface{}) {
	if proxy.impl != nil {
		proxy.impl.D(format, args...)
	}
}

func (proxy *logProxy)I(format string, args ...interface{}) {
	if proxy.impl != nil {
		proxy.impl.I(format, args...)
	}
}

func (proxy *logProxy)W(format string, args ...interface{}) {
	if proxy.impl != nil {
		proxy.impl.W(format, args...)
	}
}

func (proxy *logProxy)E(format string, args ...interface{}) {
	if proxy.impl != nil {
		proxy.impl.E(format, args...)
	}
}