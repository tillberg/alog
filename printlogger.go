package alog

type PrintLogger interface {
	Printf(format string, a ...interface{})
	Println(a ...interface{})
	Write(buf []byte) (int, error)
}
