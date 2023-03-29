package src

type IApp interface {
	Start()
	Shutdown()
	IsShutdown() bool
}
