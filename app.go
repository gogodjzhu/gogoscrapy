package gogoscrapy

type IApp interface {
	Start()
	Shutdown()
	IsShutdown() bool
}
