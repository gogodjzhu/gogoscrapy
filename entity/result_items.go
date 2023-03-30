package entity

type IResultItems interface {
	Get(key string) interface{}
	Put(key string, value interface{})
	IsSkip() bool
	All() map[string]interface{}
}

type ResultItems struct {
	Items map[string]interface{}
	Skip  bool
}

func NewResultItems(skip bool) *ResultItems {
	return &ResultItems{
		Items: make(map[string]interface{}),
		Skip:  skip,
	}
}

func (this *ResultItems) Get(key string) interface{} {
	return this.Items[key]
}

func (this *ResultItems) Put(key string, value interface{}) {
	this.Items[key] = value
}

func (this *ResultItems) IsSkip() bool {
	return this.Skip
}

func (this *ResultItems) All() map[string]interface{} {
	return this.Items
}
