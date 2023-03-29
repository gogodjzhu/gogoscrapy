package entity

type Closeable interface {
	Close() error
	IsClose() bool
}
