package http

type stdLogFunc func(...interface{})

func (l stdLogFunc) Print(args ...interface{}) {
	l(args...)
}
