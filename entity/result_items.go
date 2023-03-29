package entity

type IResultItems interface {
	Get(key string) interface{}
	Put(key string, value interface{})
	IsSkip() bool
	All() map[string]interface{}
}

type ResultItems struct {
	items map[string]interface{}
	skip  bool
}

func NewResultItems(skip bool) *ResultItems {
	return &ResultItems{
		items: make(map[string]interface{}),
		skip:  skip,
	}
}

func (this *ResultItems) Get(key string) interface{} {
	return this.items[key]
}

func (this *ResultItems) Put(key string, value interface{}) {
	this.items[key] = value
}

func (this *ResultItems) IsSkip() bool {
	return this.skip
}

func (this *ResultItems) All() map[string]interface{} {
	return this.items
}
